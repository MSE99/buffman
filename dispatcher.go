package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var processRequestsNow = make(chan bool)

type requestProcessingOpts struct {
	db           *sql.DB
	pollInterval time.Duration
}

func processStoredRequests(ctx context.Context, opts requestProcessingOpts) {
	db := opts.db
	intr := opts.pollInterval

	timer := time.NewTicker(intr)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			loadAndDispatch(ctx, db)
		case <-processRequestsNow:
			loadAndDispatch(ctx, db)
		}
	}
}

func loadAndDispatch(ctx context.Context, db *sql.DB) {
	requests, err := loadUnfinishedRequests(ctx, db)
	if err != nil {
		log.Println("error while loading requests", err)
	}

	for _, req := range requests {
		err := dispatchRequest(ctx, req)

		if err != nil {
			log.Println("error while dispatching request", err)
			return
		}

		deleteRequestByID(ctx, db, req.Id)
	}
}

func dispatchRequest(ctx context.Context, req request) error {
	httpReq, httpReqErr := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"TODO",
		strings.NewReader(req.Payload),
	)
	if httpReqErr != nil {
		return httpReqErr
	}

	httpReq.Header.Add("Content-Type", "application/json")

	res, resErr := http.DefaultClient.Do(httpReq)
	if resErr != nil {
		return resErr
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received none 200 status code: %d", res.StatusCode)
	}

	return nil
}
