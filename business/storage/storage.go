// Package storage provides an interface and implementation for storing and retrieving Ethereum blockchain transaction data.
//
// The parse.Storage interface defines methods for subscribing to Ethereum addresses, retrieving subscribers, adding transactions to the storage, and retrieving transactions for a given address.
//
// The implementation of the Storage interface uses an in-memory data store and is thread-safe.
//
// Example usage:
//
//	// Create a new storage instance
//	storage := NewMemoryStorage()
//
//	// Subscribe to an Ethereum address
//	storage.Subscribe("0x123abc")
//
//	// Add a transaction to the storage
//	tx := Transaction{From: "0x456def", To: "0x123abc", Value: 1.23}
//	storage.AddTransaction("0x123abc", tx)
//
//	// Get transactions for an address
//	transactions := storage.GetTransactions("0x123abc")
//	fmt.Println(transactions)
package storage

import (
	"sync"
	"trustwallet/business/parser"
)

// MemoryStorage is a simple in-memory storage for storing subscribed addresses and transactions.
type MemoryStorage struct {
	sync.RWMutex
	subscriptions map[string]bool
	transactions  map[string][]parser.Transaction
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		subscriptions: make(map[string]bool),
		transactions:  make(map[string][]parser.Transaction),
	}
}

func (ms *MemoryStorage) Subscribe(address string) bool {
	ms.Lock()
	defer ms.Unlock()
	if _, ok := ms.subscriptions[address]; ok {
		return false
	}
	ms.subscriptions[address] = true
	return true
}

func (ms *MemoryStorage) Subscribers() []string {
	addresses := make([]string, 0)
	for addr, _ := range ms.subscriptions {
		addresses = append(addresses, addr)
	}

	return addresses
}

func (ms *MemoryStorage) AddTransaction(address string, tx parser.Transaction) {
	ms.Lock()
	defer ms.Unlock()
	ms.transactions[address] = append(ms.transactions[address], tx)
}

func (ms *MemoryStorage) GetTransactions(address string) []parser.Transaction {
	ms.RLock()
	defer ms.RUnlock()
	return ms.transactions[address]
}
