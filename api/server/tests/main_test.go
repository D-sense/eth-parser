package tests

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
	"trustwallet/business/logger"
	"trustwallet/business/parser"
	"trustwallet/business/storage"
)

var ethParser *parser.EthereumParser

func TestMain(m *testing.M) {
	var err error
	log, err := logger.New("MENTSPACE-API")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(log)

	ethParser = parser.NewEthereumParser(storage.NewMemoryStorage(), "3b7ef887e2b244b9b0bd9b2a0c36cdf1", 5, log)

	m.Run()
}
