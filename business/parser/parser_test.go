package parser

import (
	"reflect"
	"testing"
)

// Define a mock implementation of the Parser interface
type MockParser struct{}

func (m *MockParser) GetCurrentBlock() int {
	return 0
}

func (m *MockParser) Subscribe(address string) bool {
	return true
}

func (m *MockParser) GetTransactions(address string) []Transaction {
	return []Transaction{
		{From: "0x123", To: "0x456", Value: "1.23", Status: "success"},
		{From: "0x789", To: "0xabc", Value: "4.56", Status: "pending"},
	}
}

// Define a test for the GetTransactions method
func TestGetTransactions(t *testing.T) {
	// Create a mock parser
	parser := &MockParser{}

	// Call the GetTransactions method with a mock address
	address := "0x123"
	transactions := parser.GetTransactions(address)

	// Check that the correct transactions are returned
	expectedTransactions := []Transaction{
		{From: "0x123", To: "0x456", Value: "1.23", Status: "success"},
		{From: "0x789", To: "0xabc", Value: "4.56", Status: "pending"},
	}
	if !reflect.DeepEqual(transactions, expectedTransactions) {
		t.Errorf("GetTransactions returned %+v, expected %+v", transactions, expectedTransactions)
	}
}
