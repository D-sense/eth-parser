package server

import (
	"encoding/json"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"trustwallet/business/parser"
)

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int

	// Subscribe add address to observer
	Subscribe(address string) bool

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []parser.Transaction
}

type TransactionsResponse struct {
	Transaction []parser.Transaction `json:"transactions"`
}

type CurrentBlockResponse struct {
	CurrentBlock int `json:"current_block"`
}

type SubscribeAddressResponse struct {
	Result bool `json:"result"`
}

// GetCurrentBlock encrypts a string using the Caesar Cipher.
func (h Handler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := h.Parser.GetCurrentBlock()

	blk := CurrentBlockResponse{
		CurrentBlock: block,
	}

	output, err := json.Marshal(blk)
	if err != nil {
		h.Log.Errorw("marshalling response", "data", blk, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"status": 500, "message":"internal error"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
	return
}

// Subscribe decrypts a string using the Caesar Cipher.
func (h Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	// unmarshalling data into struct
	address := param(r, "address")

	if address == "unknown" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"status": 400, "message":"address is invalid"}`)
		return
	}

	sub := h.Parser.Subscribe(address)

	ok := SubscribeAddressResponse{
		Result: sub,
	}

	output, err := json.Marshal(ok)
	if err != nil {
		h.Log.Errorw("marshalling response", "data", ok, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"status": 500, "message":"internal error"}`)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(output))
	return
}

// GetTransactions decrypts a string using the Caesar Cipher.
func (h Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := param(r, "address")

	if address == "unknown" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"status": 400, "message":"address is invalid"}`)
		return
	}

	txs := h.Parser.GetTransactions(address)

	tx := TransactionsResponse{
		Transaction: txs,
	}
	output, err := json.Marshal(tx)
	if err != nil {
		h.Log.Errorw("marshalling response", "data", tx, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"status": 500, "message":"internal error"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
	return
}

// param returns the web call parameters from the request.
func param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}
