package server

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Ctx      context.Context
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Parser   Parser
}

// Handler manages the set of user endpoints.
type Handler struct {
	Parser Parser
	Log    *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	mux := httptreemux.NewContextMux()

	// Register endpoints.
	hd := Handler{
		Parser: cfg.Parser,
		Log:    cfg.Log,
	}

	mux.Handle(http.MethodGet, "/current_block", hd.GetCurrentBlock)
	mux.Handle(http.MethodPost, "/subscribe/:address", hd.Subscribe)
	mux.Handle(http.MethodGet, "/transactions/:address", hd.GetTransactions)

	return mux
}
