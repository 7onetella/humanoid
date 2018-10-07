package main

// ExecutionPoint represents a single bot execution.
type ExecutionPoint func(request BotRequest) (response BotResponse, err error)

// Middleware is a chainable behavior modifier for executionpoints.
type Middleware func(ExecutionPoint) ExecutionPoint

// BotRequest request struct
type BotRequest struct {
	message   string
	channelID string
	approved  bool
	account   string
}

// BotResponse resposne struct
type BotResponse struct {
	message   string
	channelID string
}
