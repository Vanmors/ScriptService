package model

import "context"

type ContextCommand struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}
