package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type fmaToken struct {
	sync.RWMutex

	lastValue       string
	refreshInterval time.Duration
	ctx             context.Context
}

func (tk *fmaToken) get() string {
	tk.RLock()
	defer tk.RUnlock()

	return tk.lastValue
}

func (tk *fmaToken) refresh(conf *config) {
	tk.Lock()
	defer tk.Unlock()

	nextValue, err := fetchApiTokenFromFma(tk.ctx, conf)
	if err != nil {
		log.Println("error while refreshing token", err)
		return
	}
	tk.lastValue = nextValue
}

func (tk *fmaToken) waitAndRefresh(conf *config) {
	ticker := time.NewTicker(tk.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tk.ctx.Done():
			return

		case <-ticker.C:
			tk.refresh(conf)
		}
	}
}

func newFmaToken(ctx context.Context, conf *config) (*fmaToken, error) {
	lastValue, err := fetchApiTokenFromFma(ctx, conf)
	if err != nil {
		return &fmaToken{}, err
	}

	token := fmaToken{
		ctx:             ctx,
		lastValue:       lastValue,
		refreshInterval: time.Minute * 15,
	}

	return &token, nil
}

func fetchApiTokenFromFma(ctx context.Context, conf *config) (string, error) {
	body := strings.NewReader(
		fmt.Sprintf(`{ "username": "%s", "password": "%s" }`, conf.FmaUsername, conf.FmaPassword),
	)

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, conf.FmaURL, body)
	if reqErr != nil {
		return "", reqErr
	}

	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		return "", resErr
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received none 200 status code %v", res.StatusCode)
	}

	tokenBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", fmt.Errorf("error while reading response body: %w", readErr)
	}

	return string(tokenBytes), nil
}
