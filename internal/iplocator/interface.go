package iplocator

// IPLocator defines the interface for getting the current external IP address
type IPLocator interface {
	// GetCurrentIP returns the current external IP address
	GetCurrentIP() (string, error)
}