package buffman

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/mse99/buffman/repos"
)

var ctx = context.Background()

func connectToTestingDB(t *testing.T) *sql.DB {
	db, err := repos.ConnectToDB(ctx, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func Test_Request(t *testing.T) {
	t.Parallel()

	t.Run("insertNewRequest should return a request with filled id field", func(t *testing.T) {
		db := connectToTestingDB(t)

		now := time.Now()

		req, err := insertRequest(ctx, db, Request{
			Payload:   "Hello world",
			CreatedOn: now,
		})

		if err != nil {
			t.Fatal(err)
		} else if req.Id == 0 {
			t.Error("did not auto increment id")
		} else if req.Payload != "Hello world" {
			t.Errorf("wrong payload stored %v", req.Payload)
		} else if req.CreatedOn.UnixMilli() != now.UnixMilli() {
			t.Errorf("wrong time stored %v", req.CreatedOn)
		}
	})

	t.Run("loading unfinished requests with empty DB", func(t *testing.T) {
		db := connectToTestingDB(t)

		requests, err := loadUnfinishedRequests(ctx, db)
		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(requests, []Request{}) {
			t.Errorf("expected requests to be [] but gotten %v", requests)
		}
	})

	t.Run("inserting and loading unfinished requests", func(t *testing.T) {
		db := connectToTestingDB(t)

		now := time.Now()

		req1, err := insertRequest(ctx, db, Request{
			Payload:   "r1",
			CreatedOn: now,
		})
		if err != nil {
			t.Fatal(err)
		}

		req2, err := insertRequest(ctx, db, Request{
			Payload:   "r1",
			CreatedOn: now.Add(time.Second * 5),
		})
		if err != nil {
			t.Fatal(err)
		}

		req3, err := insertRequest(ctx, db, Request{
			Payload:   "r1",
			CreatedOn: now.Add(time.Second * 10),
		})
		if err != nil {
			t.Fatal(err)
		}

		requests, err := loadUnfinishedRequests(ctx, db)

		if err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual([]Request{req1, req2, req3}, requests) {
			t.Errorf("gotten wrong unfinished requests slice %v", requests)
		}
	})

	t.Run("inserting then deleting and loading unfinished requests", func(t *testing.T) {
		db := connectToTestingDB(t)

		now := time.Now()

		req1, err := insertRequest(ctx, db, Request{
			Payload:   "r1",
			CreatedOn: now,
		})
		if err != nil {
			t.Fatal(err)
		}

		req2, err := insertRequest(ctx, db, Request{
			Payload:   "r1",
			CreatedOn: now.Add(time.Second * 5),
		})
		if err != nil {
			t.Fatal(err)
		}

		req3, err := insertRequest(ctx, db, Request{
			Payload:   "r1",
			CreatedOn: now.Add(time.Second * 10),
		})
		if err != nil {
			t.Fatal(err)
		}

		err = deleteRequestByID(ctx, db, req3.Id)
		if err != nil {
			t.Fatal(err)
		}

		requests, err := loadUnfinishedRequests(ctx, db)

		if err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual([]Request{req1, req2}, requests) {
			t.Errorf("gotten wrong unfinished requests slice %v", requests)
		}
	})
}
