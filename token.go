package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type fmaToken struct {
	sync.RWMutex

	lastValue string
	ctx       context.Context
}

func (tk *fmaToken) get() string {
	tk.RLock()
	defer tk.RUnlock()

	return tk.lastValue
}

func (tk *fmaToken) refresh() {
	tk.Lock()
	defer tk.Unlock()

	nextValue, err := fetchApiTokenFromFma(tk.ctx)
	if err != nil {
		log.Println("error while refreshing token", err)
		return
	}
	tk.lastValue = nextValue
}

func (tk *fmaToken) waitAndRefresh() {
	ticker := time.NewTicker(loginInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tk.ctx.Done():
			log.Println("stopping token refresh")
			return

		case <-ticker.C:
			log.Println("refreshing token")
			tk.refresh()
		}
	}
}

func newFmaToken(ctx context.Context) (*fmaToken, error) {
	lastValue, err := fetchApiTokenFromFma(ctx)
	if err != nil {
		return &fmaToken{}, err
	}

	token := fmaToken{
		ctx:       ctx,
		lastValue: lastValue,
	}

	return &token, nil
}

func fetchApiTokenFromFma(ctx context.Context) (string, error) {
	body := strings.NewReader(
		fmt.Sprintf(`{ "username": "%s", "password": "%s" }`, fmaUsername, fmaPassword),
	)

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, fmaLoginURL, body)
	if reqErr != nil {
		return "", reqErr
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-app", "operator-dashboard")

	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		return "", resErr
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received none 200 status code %v", res.StatusCode)
	}

	var responseBody struct {
		Result struct {
			Token string `json:"token"`
		} `json:"result"`
	}
	decodeErr := json.NewDecoder(res.Body).Decode(&responseBody)
	if decodeErr != nil {
		return "", fmt.Errorf("error while reading response body: %w", decodeErr)
	}

	return responseBody.Result.Token, nil
}
