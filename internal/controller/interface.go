package controller

// DynDNSController is the interface for the dynamic DNS controller
type DynDNSController interface {
	// Start begins the DNS update process and blocks until signal to stop
	Start() error

	// UpdateDNSRecord performs a single DNS update check
	UpdateDNSRecord() error
}
