package main

// ExecutionPoint represents a single bot execution.
type ExecutionPoint func(request BotRequest) (response BotResponse, err error)

// Middleware is a chainable behavior modifier for executionpoints.
type Middleware func(ExecutionPoint) ExecutionPoint

// BotRequest request struct
type BotRequest struct {
	message  string
	channel  string
	approved bool
}

// BotResponse resposne struct
type BotResponse struct {
	message string
	channel string
}
