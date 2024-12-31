package handlers

import (
	"github.com/go-go-golems/prompto/pkg/server/state"
)

type Handlers struct {
	state *state.ServerState
}

func NewHandlers(state *state.ServerState) *Handlers {
	return &Handlers{
		state: state,
	}
}
