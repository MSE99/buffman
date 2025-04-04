package main

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/mse99/buffman/config"
// )

// func assertNotErr(t *testing.T, err error) {
// 	if err != nil {
// 		t.Errorf("gotten error %v", err)
// 	}
// }

// func createDispatchServer(t *testing.T) (*httptest.Server, func() [][]byte) {
// 	var (
// 		l        sync.Mutex
// 		payloads = [][]byte{}
// 	)

// 	dispatchServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		l.Lock()
// 		defer l.Unlock()

// 		defer r.Body.Close()

// 		bytes, err := io.ReadAll(r.Body)
// 		assertNotErr(t, err)

// 		payloads = append(payloads, bytes)
// 	}))
// 	t.Cleanup(dispatchServer.Close)

// 	return dispatchServer, func() [][]byte {
// 		l.Lock()
// 		defer l.Unlock()

// 		return payloads
// 	}
// }

// func createLoginServer(t *testing.T, username, password string) *httptest.Server {
// 	loginServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()

// 		var body struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 		}

// 		decodeErr := json.NewDecoder(r.Body).Decode(&body)
// 		assertNotErr(t, decodeErr)

// 		if body.Username != username && body.Password != password {
// 			t.Errorf("invalid credentials supplied to login server %s %s", body.Username, body.Password)
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{ "result": { "token": "FMA_TOKEN_FROM_LOGIN" } }`))
// 	}))
// 	t.Cleanup(func() { loginServer.Close() })

// 	return loginServer
// }

// func createTestApp(t *testing.T) (*fiber.App, *sql.DB) {
// 	tCtx, cancel := context.WithCancel(ctx)
// 	t.Cleanup(cancel)

// 	db := connectToTestingDB(t)

// 	app := fiber.New()
// 	t.Cleanup(func() {
// 		app.Shutdown()
// 	})
// 	setupRouter(tCtx, app, db)

// 	token, err := newFmaToken(tCtx)
// 	if err != nil {
// 		t.Fatalf("Error while creating test app %v", err)
// 	}
// 	go token.waitAndRefresh()

// 	go processStoredRequests(tCtx, requestProcessingOpts{
// 		db: db,
// 		tk: token,
// 	})

// 	return app, db
// }

// func TestStatus(t *testing.T) {
// 	db := connectToTestingDB(t)

// 	app := fiber.New()
// 	defer app.Shutdown()

// 	setupRouter(ctx, app, db)

// 	req := httptest.NewRequest(http.MethodGet, "/status", nil)

// 	res, err := app.Test(req)
// 	assertNotErr(t, err)

// 	if res.StatusCode != http.StatusOK {
// 		t.Errorf("server responded with non 200 status %d", res.StatusCode)
// 	}
// }

// func TestRequestProcessing(t *testing.T) {
// 	config.FmaUsername = "admin"
// 	config.FmaPassword = "admin"

// 	t.Run("InvalidSecret", func(t *testing.T) {
// 		loginServer := createLoginServer(t, config.FmaUsername, config.FmaPassword)
// 		dispatchServer, _ := createDispatchServer(t)

// 		config.FmaLoginURL = loginServer.URL
// 		config.FmaDispatchURL = dispatchServer.URL
// 		config.LoginInterval = time.Millisecond * 100
// 		config.PollInterval = time.Millisecond * 100
// 		config.OdooSecret = "Foo-123"

// 		app, _ := createTestApp(t)

// 		req := httptest.NewRequest(http.MethodPost, "/", nil)

// 		res, err := app.Test(req)
// 		assertNotErr(t, err)

// 		if res.StatusCode != http.StatusUnauthorized {
// 			t.Errorf("expected 401 but got %d", res.StatusCode)
// 		}
// 	})

// 	t.Run("InvalidPayload", func(t *testing.T) {
// 		loginServer := createLoginServer(t, config.FmaUsername, config.FmaPassword)
// 		dispatchServer, _ := createDispatchServer(t)

// 		config.FmaLoginURL = loginServer.URL
// 		config.FmaDispatchURL = dispatchServer.URL
// 		config.LoginInterval = time.Millisecond * 100
// 		config.PollInterval = time.Millisecond * 100
// 		config.OdooSecret = "hi!"

// 		app, _ := createTestApp(t)

// 		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?token=%s", config.OdooSecret), nil)

// 		res, err := app.Test(req)
// 		assertNotErr(t, err)

// 		if res.StatusCode != http.StatusBadRequest {
// 			t.Errorf("expected 400 but got %d", res.StatusCode)
// 		}
// 	})

// 	t.Run("ShouldImmediatelyDispatchRequest", func(t *testing.T) {
// 		loginServer := createLoginServer(t, config.FmaUsername, config.FmaPassword)
// 		dispatchServer, getDispatches := createDispatchServer(t)

// 		config.FmaLoginURL = loginServer.URL
// 		config.FmaDispatchURL = dispatchServer.URL
// 		config.LoginInterval = time.Millisecond * 100
// 		config.PollInterval = time.Millisecond * 100
// 		config.OdooSecret = "hi!"

// 		app, _ := createTestApp(t)

// 		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?token=%s", config.OdooSecret), strings.NewReader("FOO IS GREAT"))

// 		res, err := app.Test(req)
// 		assertNotErr(t, err)

// 		if res.StatusCode != http.StatusOK {
// 			t.Errorf("expected 200 but got %d", res.StatusCode)
// 		}
// 		time.Sleep(time.Millisecond * 200)

// 		lastDispatch := getDispatches()[0]

// 		if string(lastDispatch) != "FOO IS GREAT" {
// 			t.Errorf("expected last dispatched value to be FOO IS GREAT but got %s", string(lastDispatch))
// 		}
// 	})

// 	t.Run("ShouldDispatchOlderRequestsFirst", func(t *testing.T) {
// 		loginServer := createLoginServer(t, config.FmaUsername, config.FmaPassword)
// 		dispatchServer, getDispatches := createDispatchServer(t)

// 		config.FmaLoginURL = loginServer.URL
// 		config.FmaDispatchURL = dispatchServer.URL
// 		config.LoginInterval = time.Millisecond * 100
// 		config.PollInterval = time.Millisecond * 100
// 		config.OdooSecret = "hi!"

// 		app, db := createTestApp(t)

// 		_, err := insertRequest(ctx, db, request{
// 			CreatedOn: time.Now().Add(-time.Hour * 3),
// 			Payload:   "FOO",
// 		})
// 		assertNotErr(t, err)
// 		_, err = insertRequest(ctx, db, request{
// 			CreatedOn: time.Now().Add(-time.Hour * 2),
// 			Payload:   "BAR",
// 		})
// 		assertNotErr(t, err)
// 		_, err = insertRequest(ctx, db, request{
// 			CreatedOn: time.Now().Add(-time.Hour),
// 			Payload:   "BAZ",
// 		})
// 		assertNotErr(t, err)

// 		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?token=%s", config.OdooSecret), strings.NewReader("NAZ"))

// 		res, err := app.Test(req)
// 		assertNotErr(t, err)

// 		if res.StatusCode != http.StatusOK {
// 			t.Errorf("expected 200 but got %d", res.StatusCode)
// 		}
// 		time.Sleep(time.Millisecond * 500)

// 		dispatches := getDispatches()

// 		if string(dispatches[0]) != "FOO" {
// 			t.Errorf("expected last dispatched value to be FOO but got %s", string(dispatches[0]))
// 		}
// 		if string(dispatches[1]) != "BAR" {
// 			t.Errorf("expected last dispatched value to be BAR but got %s", string(dispatches[1]))
// 		}
// 		if string(dispatches[2]) != "BAZ" {
// 			t.Errorf("expected last dispatched value to be BAZ but got %s", string(dispatches[2]))
// 		}
// 		if string(dispatches[3]) != "NAZ" {
// 			t.Errorf("expected last dispatched value to be NAZ but got %s", string(dispatches[2]))
// 		}
// 	})

// 	t.Run("ShouldDispatchAfterInterval", func(t *testing.T) {
// 		loginServer := createLoginServer(t, config.FmaUsername, config.FmaPassword)
// 		dispatchServer, getDispatches := createDispatchServer(t)

// 		config.FmaLoginURL = loginServer.URL
// 		config.FmaDispatchURL = dispatchServer.URL
// 		config.OdooSecret = "hi!"
// 		config.PollInterval = time.Millisecond * 350
// 		config.LoginInterval = time.Millisecond * 100

// 		_, db := createTestApp(t)
// 		_, err := insertRequest(ctx, db, request{
// 			CreatedOn: time.Now().Add(-time.Hour * 3),
// 			Payload:   "FOO",
// 		})
// 		assertNotErr(t, err)

// 		time.Sleep(time.Millisecond * 500)

// 		dispatches := getDispatches()

// 		if string(dispatches[0]) != "FOO" {
// 			t.Errorf("expected last dispatched value to be FOO but got %s", string(dispatches[0]))
// 		}
// 	})

// 	t.Run("TokenRefresh", func(t *testing.T) {
// 		var (
// 			l          sync.Mutex
// 			loginCount = 0
// 		)

// 		loginServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			l.Lock()
// 			defer l.Unlock()

// 			defer r.Body.Close()

// 			var body struct {
// 				Username string `json:"username"`
// 				Password string `json:"password"`
// 			}

// 			decodeErr := json.NewDecoder(r.Body).Decode(&body)
// 			assertNotErr(t, decodeErr)

// 			if body.Username != "admin" && body.Password != "admin" {
// 				t.Errorf("invalid credentials supplied to login server %s %s", body.Username, body.Password)
// 			}

// 			w.WriteHeader(http.StatusOK)
// 			w.Write([]byte(fmt.Sprintf(`{ "result": { "token": "FMA_TOKEN_FROM_LOGIN_%d" } }`, loginCount)))

// 			loginCount++
// 		}))

// 		dispatchServer, _ := createDispatchServer(t)

// 		config.FmaLoginURL = loginServer.URL
// 		config.FmaDispatchURL = dispatchServer.URL
// 		config.OdooSecret = "hi!"
// 		config.PollInterval = time.Millisecond * 350
// 		config.LoginInterval = time.Millisecond * 100

// 		createTestApp(t)

// 		time.Sleep(time.Millisecond * 500)

// 		l.Lock()
// 		defer l.Unlock()

// 		if loginCount < 3 {
// 			t.Errorf("Expected fma token to be refreshed at least 3 times but got %d refreshes", loginCount)
// 		}
// 	})
// }
