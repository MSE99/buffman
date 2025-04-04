package buffman

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mse99/buffman/config"
)

var processRequestsNow = make(chan struct{})

type requestProcessingOpts struct {
	db *sql.DB
	tk *fmaToken
}

func processStoredRequests(ctx context.Context, opts requestProcessingOpts) {
	intr := config.PollInterval

	timer := time.NewTicker(intr)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down request polling")
			return
		case <-timer.C:
			log.Println("timed poll for stored requests")
			loadAndDispatch(ctx, opts)
		case <-processRequestsNow:
			log.Println("polling because of a poll signal")
			loadAndDispatch(ctx, opts)
		}
	}
}

func loadAndDispatch(ctx context.Context, opts requestProcessingOpts) {
	requests, err := loadUnfinishedRequests(ctx, opts.db)
	if err != nil {
		log.Println("error while loading requests", err)
		return
	}

	for _, req := range requests {
		err := dispatchRequest(ctx, req, opts)

		if err != nil {
			log.Println("error while dispatching request", err)

			if !config.ContinueOnError {
				log.Println("stopping dispatch")
				return
			}
		}

		deleteRequestByID(ctx, opts.db, req.Id)
	}
}

func dispatchRequest(ctx context.Context, req Request, opts requestProcessingOpts) error {
	httpReq, httpReqErr := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		config.FmaDispatchURL,
		strings.NewReader(req.Payload),
	)
	if httpReqErr != nil {
		return httpReqErr
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add(
		"Authorization",
		fmt.Sprintf(`Bearer %s`, opts.tk.get()),
	)

	res, resErr := http.DefaultClient.Do(httpReq)
	if resErr != nil {
		return resErr
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received none 200 status code: %d", res.StatusCode)
	}

	return nil
}

func QueueRequest(ctx context.Context, db *sql.DB, payload string) error {
	if len(strings.Trim(payload, " ")) == 0 {
		return errors.New("request payload cannot be empty")
	}

	_, err := insertRequest(ctx, db, Request{
		Payload:   payload,
		CreatedOn: time.Now(),
	})

	if err != nil {
		return err
	}

	select {
	case processRequestsNow <- struct{}{}:
		return nil
	case <-ctx.Done():
		return nil
	}
}
