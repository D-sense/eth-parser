package storage

import (
	"reflect"
	"testing"
	"trustwallet/business/parser"
)

// Define a mock implementation of the Storage interface
type MockStorage struct {
	transactions map[string][]parser.Transaction
	subscribers  map[string]bool
}

func (m *MockStorage) Subscribe(address string) bool {
	if m.subscribers == nil {
		m.subscribers = make(map[string]bool)
	}
	m.subscribers[address] = true
	return true
}

func (m *MockStorage) Subscribers() []string {
	var subscribers []string
	for k := range m.subscribers {
		subscribers = append(subscribers, k)
	}
	return subscribers
}

func (m *MockStorage) AddTransaction(address string, tx parser.Transaction) {
	if m.transactions == nil {
		m.transactions = make(map[string][]parser.Transaction)
	}
	m.transactions[address] = append(m.transactions[address], tx)
}

func (m *MockStorage) GetTransactions(address string) []parser.Transaction {
	if m.transactions == nil {
		return []parser.Transaction{}
	}
	return m.transactions[address]
}

// Define a test for the GetTransactions method
func TestGetTransactions(t *testing.T) {
	// Create a mock storage
	storage := &MockStorage{}

	// Add some transactions for a mock address
	address := "0x123"
	tx1 := parser.Transaction{From: "0x123", To: "0x456", Value: "1.23", Status: "success"}
	tx2 := parser.Transaction{From: "0x789", To: "0xabc", Value: "4.56", Status: "pending"}
	storage.AddTransaction(address, tx1)
	storage.AddTransaction(address, tx2)

	// Call the GetTransactions method with the mock address
	transactions := storage.GetTransactions(address)

	// Check that the correct transactions are returned
	expectedTransactions := []parser.Transaction{tx1, tx2}
	if !reflect.DeepEqual(transactions, expectedTransactions) {
		t.Errorf("GetTransactions returned %+v, expected %+v", transactions, expectedTransactions)
	}
}
