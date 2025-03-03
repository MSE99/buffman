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
	tk           *fmaToken
}

func processStoredRequests(ctx context.Context, opts requestProcessingOpts) {
	intr := opts.pollInterval

	timer := time.NewTicker(intr)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			loadAndDispatch(ctx, opts)
		case <-processRequestsNow:
			loadAndDispatch(ctx, opts)
		}
	}
}

func loadAndDispatch(ctx context.Context, opts requestProcessingOpts) {
	requests, err := loadUnfinishedRequests(ctx, opts.db)
	if err != nil {
		log.Println("error while loading requests", err)
	}

	for _, req := range requests {
		err := dispatchRequest(ctx, req, opts)

		if err != nil {
			log.Println("error while dispatching request", err)
			return
		}

		deleteRequestByID(ctx, opts.db, req.Id)
	}
}

func dispatchRequest(ctx context.Context, req request, opts requestProcessingOpts) error {
	httpReq, httpReqErr := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmaDispatchURL,
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
