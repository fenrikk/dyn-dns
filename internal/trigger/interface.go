package trigger

// Handler defines the function that will be called by the trigger
type Handler func() error

// Trigger defines the interface for the component that initiates the DNS check
type Trigger interface {
	// Start starts the trigger with the specified handler
	Start(handler Handler) error
	
	// stops the trigger
	Stop() error
}