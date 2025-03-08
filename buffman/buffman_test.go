package buffman

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/mse99/buffman/config"
	"github.com/mse99/buffman/repos"
)

func createTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

func createTestDB(t *testing.T) *sql.DB {
	db, err := repos.ConnectToDB(ctx, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestDispatching(t *testing.T) {
	config.PollInterval = time.Millisecond * 100
	config.LoginInterval = time.Millisecond * 100

	t.Run("InitialAuthFail", func(t *testing.T) {
		loginServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized!"))
		})

		config.FmaLoginURL = loginServer.URL

		db := createTestDB(t)
		err := StartDispatchToFMA(ctx, db)

		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("FmaTokenHydration", func(t *testing.T) {
		var (
			loginCount = 0
			loginLock  = sync.Mutex{}
		)

		loginServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			loginLock.Lock()
			defer loginLock.Unlock()

			loginCount++

			var loginBody struct {
				Username string
				Password string
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&loginBody)

			if decodeErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body!"))
			} else if loginBody.Username == config.FmaUsername && loginBody.Password == config.FmaPassword {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ "result": { "token": "FMA-token" } }`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized!"))
			}
		})

		dispatchServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		config.FmaDispatchURL = dispatchServer.URL
		config.FmaLoginURL = loginServer.URL
		config.FmaUsername = "admin"
		config.FmaPassword = "admin"

		db := createTestDB(t)

		err := StartDispatchToFMA(ctx, db)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Millisecond * 150)

		loginLock.Lock()
		defer loginLock.Unlock()

		if loginCount != 2 {
			t.Errorf("expected login count to be 2 but got %d", loginCount)
		}
	})

	t.Run("Dispatch", func(t *testing.T) {
		var (
			lock     = sync.Mutex{}
			payloads = []string{}
		)

		loginServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			lock.Lock()
			defer lock.Unlock()

			var loginBody struct {
				Username string
				Password string
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&loginBody)

			if decodeErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body!"))
			} else if loginBody.Username == config.FmaUsername && loginBody.Password == config.FmaPassword {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ "result": { "token": "FMA-token" } }`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized!"))
			}
		})

		dispatchServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			lock.Lock()
			defer lock.Unlock()

			payloadBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body"))
				return
			}
			payloads = append(payloads, string(payloadBytes))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		config.FmaDispatchURL = dispatchServer.URL
		config.FmaLoginURL = loginServer.URL
		config.FmaUsername = "admin"
		config.FmaPassword = "admin"

		db := createTestDB(t)

		err := StartDispatchToFMA(ctx, db)
		if err != nil {
			t.Error(err)
		}

		queueErr := QueueRequest(ctx, db, `{ "x_id": 123 }`)
		if queueErr != nil {
			t.Error(queueErr)
		}
		time.Sleep(time.Millisecond * 150)

		lock.Lock()
		defer lock.Unlock()
		p1 := payloads[0]

		if p1 != `{ "x_id": 123 }` {
			t.Errorf("unexpected payload %s", p1)
		}

		requests, loadErr := loadUnfinishedRequests(ctx, db)
		if loadErr != nil {
			t.Error(loadErr)
		} else if len(requests) != 0 {
			t.Error("did not process all requests")
		}
	})

	t.Run("ShouldRetryOnDispatchFailure", func(t *testing.T) {
		var (
			lock     = sync.Mutex{}
			payloads = []string{}
		)

		loginServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			var loginBody struct {
				Username string
				Password string
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&loginBody)

			if decodeErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body!"))
			} else if loginBody.Username == config.FmaUsername && loginBody.Password == config.FmaPassword {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ "result": { "token": "FMA-token" } }`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized!"))
			}
		})

		dispatchServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			lock.Lock()
			defer lock.Unlock()

			payloadBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body"))
				return
			}
			payloads = append(payloads, string(payloadBytes))

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error something bad happened try again later"))
		})

		config.FmaDispatchURL = dispatchServer.URL
		config.FmaLoginURL = loginServer.URL
		config.FmaUsername = "admin"
		config.FmaPassword = "admin"

		db := createTestDB(t)

		err := StartDispatchToFMA(ctx, db)
		if err != nil {
			t.Error(err)
		}

		queueErr := QueueRequest(ctx, db, `{ "x_id": 123 }`)
		if queueErr != nil {
			t.Error(queueErr)
		}
		time.Sleep(time.Millisecond * 250)

		lock.Lock()
		defer lock.Unlock()

		p1 := payloads[0]
		p2 := payloads[1]

		if p1 != `{ "x_id": 123 }` {
			t.Errorf("unexpected payload %s", p1)
		}

		if p2 != `{ "x_id": 123 }` {
			t.Errorf("unexpected payload %s", p2)
		}

		requests, loadErr := loadUnfinishedRequests(ctx, db)
		if loadErr != nil {
			t.Error(loadErr)
		} else if len(requests) != 1 {
			t.Error("did not process all requests")
		}
	})

	t.Run("DispatchOnInterval", func(t *testing.T) {
		var (
			lock     = sync.Mutex{}
			payloads = []string{}
		)

		loginServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			lock.Lock()
			defer lock.Unlock()

			var loginBody struct {
				Username string
				Password string
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&loginBody)

			if decodeErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body!"))
			} else if loginBody.Username == config.FmaUsername && loginBody.Password == config.FmaPassword {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ "result": { "token": "FMA-token" } }`))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized!"))
			}
		})

		dispatchServer := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			lock.Lock()
			defer lock.Unlock()

			payloadBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid request body"))
				return
			}
			payloads = append(payloads, string(payloadBytes))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		config.FmaDispatchURL = dispatchServer.URL
		config.FmaLoginURL = loginServer.URL
		config.FmaUsername = "admin"
		config.FmaPassword = "admin"

		db := createTestDB(t)

		err := StartDispatchToFMA(ctx, db)
		if err != nil {
			t.Error(err)
		}

		_, err1 := insertRequest(ctx, db, Request{
			Payload:   "FOO",
			CreatedOn: time.Now().Add(-time.Second * 5),
		})
		if err1 != nil {
			t.Error(err1)
		}

		_, err2 := insertRequest(ctx, db, Request{
			Payload:   "BAR",
			CreatedOn: time.Now().Add(-time.Second * 3),
		})
		if err2 != nil {
			t.Error(err2)
		}

		_, err3 := insertRequest(ctx, db, Request{
			Payload:   "BAZ",
			CreatedOn: time.Now().Add(-time.Second),
		})
		if err3 != nil {
			t.Error(err3)
		}
		time.Sleep(time.Millisecond * 200)

		lock.Lock()
		defer lock.Unlock()

		if !reflect.DeepEqual(payloads, []string{"FOO", "BAR", "BAZ"}) {
			t.Errorf("invalid payloads order %v", payloads)
		}

		reqs, reqsErr := loadUnfinishedRequests(ctx, db)
		if reqsErr != nil {
			t.Error(reqsErr)
		}
		if len(reqs) != 0 {
			t.Error("requests were not removed from queue")
		}
	})
}
