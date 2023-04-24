// Package parser provides interfaces and implementations for parsing and storing data related to a blockchain.
//
// The EthereumParser interface defines methods for parsing blockchain data and subscribing to updates. The Storage interface
// defines methods for storing and retrieving data related to blockchain transactions.
//
// The package includes mock implementations of the EthereumParser and Storage interfaces for use in testing, as well as
// test cases for these implementations. These tests cover basic functionality and edge cases.
//
// In addition to the mock implementations provided, the package allows for easy creation of custom implementations
// of the EthereumParser and Storage interfaces using any storage mechanism desired.
//
// Example usage:
//
//	// Create a new Parser implementation
//	parser := NewEthereumParser()
//
//	// Subscribe to updates for a particular address
//	parser.Subscribe("0x123456789abcdef")
//
//	// Get the current block number
//	currentBlock := parser.GetCurrentBlock()
//
//	// Get transactions for a particular address
//	transactions := parser.GetTransactions("0x123456789abcdef")
//
//	// Create a new Storage implementation
//	storage := NewMemoryStorage()
//
//	// Add a transaction to the storage
//	tx := Transaction{From: "0x123456789abcdef", To: "0xabcdef123456789", Value: 1.23, Status: "success"}
//	storage.AddTransaction("0x123456789abcdef", tx)
//
//	// Get transactions for a particular address from the storage
//	transactions := storage.GetTransactions("0x123456789abcdef")
package parser

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Storage interface {
	Subscribe(address string) bool
	Subscribers() []string
	AddTransaction(address string, tx Transaction)
	GetTransactions(address string) []Transaction
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash        string   `json:"hash"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Value       string   `json:"value"`
	Status      string   `json:"status"`
	Gas         string   `json:"gas"`
	GasPrice    string   `json:"gasPrice"`
	BlockNumber *big.Int `json:"blockNumber"`
	BlockHash   string   `json:"blockHash"`
}

type BlockResp struct {
	Result struct {
		Transactions struct {
			Hash     string `json:"hash"`
			From     string `json:"from"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Gas      string `json:"gas"`
			GasPrice string `json:"gasPrice"`
		} `json:"transactions"`
	} `json:"result"`
}

// EthereumParser implements the Parser interface
type EthereumParser struct {
	httpClient      *http.Client
	ethNodeURL      string
	storage         Storage
	currentBlock    int
	lastPolledBlock int
	lock            sync.Mutex
	pollingInterval time.Duration
	Log             *zap.SugaredLogger
}

// NewEthereumParser creates a new Ethereum Parser instance
func NewEthereumParser(storage Storage, nodeEndpoint string, pollingInterval time.Duration, logger *zap.SugaredLogger) *EthereumParser {
	client := &EthereumParser{
		httpClient:      &http.Client{},
		ethNodeURL:      nodeEndpoint,
		storage:         storage,
		currentBlock:    0,
		lastPolledBlock: 0,
		lock:            sync.Mutex{},
		pollingInterval: pollingInterval * time.Second,
		Log:             logger,
	}

	// Start polling Ethereum node
	go func() {
		client.pollTransactions()
	}()

	return client
}

// Subscribe Creates an address subscription
func (p *EthereumParser) Subscribe(address string) bool {
	return p.storage.Subscribe(address)
}

// GetCurrentBlock Gets the current block number
func (p *EthereumParser) GetCurrentBlock() int {
	return p.currentBlock
}

// GetTransactions Gets an address's transactions
func (p *EthereumParser) GetTransactions(address string) []Transaction {
	if txs := p.storage.GetTransactions(address); txs != nil {
		return txs
	}

	return []Transaction{}
}

// pollTransactions Pools Ethereum gateway for new updates and updates the local storage
func (p *EthereumParser) pollTransactions() {
	for {
		// Wait for this period
		time.Sleep(p.pollingInterval)

		// Check for new transactions for each subscribed address
		for _, address := range p.storage.Subscribers() {
			// Get the transactions for the address since the last polled block
			err := p.getTransactionsSinceBlock(address)
			if err != nil {
				p.Log.Errorw("marshalling response", "error", err)
			}
		}
	}
}

// GetTransactions returns a list of inbound or outbound transactions for an address
func (p *EthereumParser) getTransactionsSinceBlock(address string) error {
	// make a JSONRPC call to get the latest block number
	reqBody := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)
	resp, err := http.Post(p.ethNodeURL, "application/json", strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var blockNumResp struct {
		Result string `json:"result"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&blockNumResp); err != nil {
		return err
	}
	blockNum := new(big.Int)
	blockNum, _ = blockNum.SetString(blockNumResp.Result[2:], 16)

	// iterate over blocks starting from the last parsed block
	for i := p.currentBlock + 1; i <= int(blockNum.Int64()); i++ {
		// make a JSONRPC call to get block data
		reqBody = fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x%x",true],"id":1}`, i)
		resp, err = http.Post(p.ethNodeURL, "application/json", strings.NewReader(reqBody))
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		var blockResp struct {
			Result struct {
				Transactions []struct {
					Hash        string   `json:"hash"`
					From        string   `json:"from"`
					To          string   `json:"to"`
					Value       string   `json:"value"`
					Gas         string   `json:"gas"`
					GasPrice    string   `json:"gasPrice"`
					BlockNumber *big.Int `json:"blockNumber"`
					BlockHash   string   `json:"blockHash"`
				} `json:"transactions"`
			} `json:"result"`
		}

		if err = json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
			return err
		}

		for _, tx := range blockResp.Result.Transactions {
			if tx.From == address {
				p.storage.AddTransaction(tx.From, Transaction{
					Hash:     tx.Hash,
					From:     tx.From,
					To:       tx.To,
					Value:    tx.Value,
					Gas:      tx.Gas,
					GasPrice: tx.GasPrice,
				})
			} else if tx.To == address {
				p.storage.AddTransaction(tx.To, Transaction{
					Hash:     tx.Hash,
					From:     tx.From,
					To:       tx.To,
					Value:    tx.Value,
					Gas:      tx.Gas,
					GasPrice: tx.GasPrice,
				})
			}
		}
	}

	// update the last parsed block number
	p.currentBlock = int(blockNum.Int64())

	return nil
}
