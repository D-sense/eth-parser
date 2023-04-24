## Ethereum Blockchain Parser
This project implements an Ethereum blockchain parser in Go that allows you to query transactions for subscribed addresses. The main goal of this project is to provide a solution for users who are not able to receive push notifications for incoming/outgoing transactions.

### Problem
The Ethereum blockchain is a decentralized system, and users can subscribe to specific addresses to receive notifications about incoming and outgoing transactions. However, not all users are able to receive push notifications due to various reasons such as unreliable network connection, limited device capabilities, or other technical issues.

To solve this problem, we have implemented a parser interface that allows users to subscribe to specific Ethereum addresses and query the blockchain for transactions related to those addresses.

### Limitations
- This project uses the Go programming language.
- No external libraries are used.
- Ethereum JSONRPC is used to interact with the Ethereum blockchain.
- Memory storage is used to store all data, and the system is designed to be easily extendable to support other storage options in the future.

### Interface
The public interface for the Ethereum blockchain parser is defined by the Parser interface, which includes the following methods:

```go
  type Parser interface {
        // Get the current block number
        GetCurrentBlock() int

	// Add an address to the list of observers
	Subscribe(address string) bool

	// Get a list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
  }
```

### Usage
To use this Ethereum blockchain parser, you can integrate it into your code or use it via a command-line interface or REST API.
If you want to integrate the parser into your code, you can import the parser package and use the methods defined in the Parser interface.

```go
import (
"github.com/yourusername/eth-parser/parser"
)

func main() {
// create a new parser instance
p := parser.NewParser()

    // subscribe to an address
    p.Subscribe("0x123abc")

    // get transactions for an address
    transactions := p.GetTransactions("0x123abc")

    // process transactions
    for _, tx := range transactions {
        // do something with the transaction
    }
}
```
Alternatively, you can use the Ethereum blockchain parser via a command-line interface or REST API. In this case, you can use the cli or api packages, respectively.

### Installation
To use this Ethereum blockchain parser, you need to have Go installed on your system. Then, you can install the parser by running the following command:

```
go get github.com/d-sense/eth-parser
```

### Testing
To run the unit tests for this project, you can use the following command:

```
go test -v ./...
```

### Conclusion
In conclusion, the Ethereum blockchain parser implemented in this project provides a simple and efficient solution for users who are not able to receive push notifications for incoming/outgoing transactions. The system is designed to be easily extendable and can be integrated into your code or used via a command-line interface or REST API.





