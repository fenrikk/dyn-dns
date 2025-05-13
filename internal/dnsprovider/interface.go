package dnsprovider

type DNSProvider interface {
	// UpdateDNSRecord updates the DNS record with the specified IP address
	// Returns true if the record was updated, and false if no changes are needed
	UpdateDNSRecord(ip string) (bool, error)
	
	// GetCurrentDNSRecord gets the current IP address in the DNS record
	GetCurrentDNSRecord() (string, error)
}