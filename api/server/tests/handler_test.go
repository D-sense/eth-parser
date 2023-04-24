package tests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"trustwallet/api/server"
	"trustwallet/business/logger"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

// UserTests holds methods for each user subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type HandlerTests struct {
	app http.Handler
}

// Test_Encypter is the entry point for testing Encrypter functions.
func Test_Encypter_Decrypter(t *testing.T) {
	t.Parallel()

	sugaredLogger, err := logger.New("test-server")
	if err != nil {
		log.Fatal("error initializing sugaredLogger")
	}
	shutdown := make(chan os.Signal, 1)
	tests := HandlerTests{
		app: server.APIMux(server.APIMuxConfig{
			Shutdown: shutdown,
			Log:      sugaredLogger,
			Parser:   ethParser,
		}),
	}

	t.Run("currentBlock200", tests.currentBlock200)
	t.Run("subscribeAddress400", tests.subscribeAddress400)
	t.Run("subscribeAddress200", tests.subscribeAddress200)
	t.Run("getTransactions400", tests.getTransactions400)
	t.Run("getTransactions200", tests.getTransactions200)
}

// currentBlock200 get current block number.
func (ht *HandlerTests) currentBlock200(t *testing.T) {
	t.Log("Should return current block number")
	{
		w := ht.helperHttpClient(http.MethodGet, "/current_block", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("%s Should receive a status code of 200 for the response : %v", failed, w.Code)
		}

		// TODO: parse the response body and compare result

		t.Logf("%s Should receive a status code of 200 for the response", success)
	}
}

// subscribeAddress400 subscribe a new address.
func (ht *HandlerTests) subscribeAddress400(t *testing.T) {
	t.Log("Should return 400 for empty address")
	{
		w := ht.helperHttpClient(http.MethodPost, fmt.Sprintf("/subscribe/%v", "unknown"), nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s Should receive a status code of 400 for the response : %v", failed, w.Code)
		}

		// TODO: parse the response body and compare result

		t.Logf("%s Should receive a status code of 400 for the response", success)
	}
}

// subscribeAddress200 subscribe a new address.
func (ht *HandlerTests) subscribeAddress200(t *testing.T) {
	t.Log("Should return 201 for a valid address")
	{
		w := ht.helperHttpClient(http.MethodPost, fmt.Sprintf("/subscribe/%v", "0x123"), nil)
		if w.Code != http.StatusCreated {
			t.Fatalf("%s Should receive a status code of 201 for the response : %v", failed, w.Code)
		}

		// TODO: parse the response body and compare result

		t.Logf("%s Should receive a status code of 201 for the response", success)
	}
}

// getTransactions400 get transactions for an address.
func (ht *HandlerTests) getTransactions400(t *testing.T) {
	t.Log("Should return 400 for empty address")
	{
		w := ht.helperHttpClient(http.MethodGet, fmt.Sprintf("/transactions/%v", "unknown"), nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s Should receive a status code of 400 for the response : %v", failed, w.Code)
		}

		// TODO: parse the response body and compare result

		t.Logf("%s Should receive a status code of 400 for the response", success)
	}
}

// getTransactions400 get transactions for an address.
func (ht *HandlerTests) getTransactions200(t *testing.T) {
	t.Log("Should return 200 for empty address")
	{
		w := ht.helperHttpClient(http.MethodGet, fmt.Sprintf("/transactions/%v", "0x123"), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("%s Should receive a status code of 200 for the response : %v", failed, w.Code)
		}

		// TODO: parse the response body and compare result

		t.Logf("%s Should receive a status code of 200 for the response", success)
	}
}

func (ht *HandlerTests) helperHttpClient(method, url string, body []byte) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	ht.app.ServeHTTP(w, r)
	return w
}
