package controller

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fenrikk/dyn-dns/internal/dnsprovider"
	"github.com/fenrikk/dyn-dns/internal/iplocator"
	"github.com/fenrikk/dyn-dns/internal/trigger"
)

const DEFAULT_CHECK_INTERVAL = 5

// DefaultController implements the DynDNSController interface
type DefaultController struct {
	dnsProvider dnsprovider.DNSProvider
	ipLocator   iplocator.IPLocator
	logger      *log.Logger
	interval    time.Duration
}

// ControllerOptions contains configuration for the controller
type ControllerOptions struct {
	CheckInterval time.Duration
	Logger        *log.Logger
}

// NewController creates a new DefaultController
func NewController(dnsProvider dnsprovider.DNSProvider, ipLocator iplocator.IPLocator, opts ControllerOptions) *DefaultController {
	// Set default values if not provided
	interval := opts.CheckInterval
	if interval == 0 {
		interval = DEFAULT_CHECK_INTERVAL * time.Minute
	}

	logger := opts.Logger
	if logger == nil {
		logger = log.New(log.Writer(), "[DYNDNS] ", log.LstdFlags)
	}

	return &DefaultController{
		dnsProvider: dnsProvider,
		ipLocator:   ipLocator,
		logger:      logger,
		interval:    interval,
	}
}

// Start begins the DNS update process and blocks until signal to stop
func (c *DefaultController) Start() error {
	c.logger.Printf("Starting dynamic DNS controller with check interval of %s", c.interval)

	// Create time trigger
	timeTrigger := trigger.NewTimeTrigger(c.interval)

	// Run initial check immediately
	if err := c.UpdateDNSRecord(); err != nil {
		c.logger.Printf("Initial DNS check failed: %v", err)
	}

	// Start the trigger with our update handler
	if err := timeTrigger.Start(c.UpdateDNSRecord); err != nil {
		return fmt.Errorf("failed to start trigger: %w", err)
	}

	c.logger.Println("DNS update service running. Press Ctrl+C to stop.")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	c.logger.Println("Shutdown signal received, stopping service...")

	// Stop the trigger
	if err := timeTrigger.Stop(); err != nil {
		return fmt.Errorf("failed to stop trigger: %w", err)
	}

	c.logger.Println("Service stopped successfully")
	return nil
}

// UpdateDNSRecord performs a single DNS update check
func (c *DefaultController) UpdateDNSRecord() error {
	// Get current IP address
	ip, err := c.ipLocator.GetCurrentIP()
	if err != nil {
		c.logger.Printf("Failed to get current IP: %v", err)
		return fmt.Errorf("failed to get current IP: %w", err)
	}

	// Update DNS record if needed
	updated, err := c.dnsProvider.UpdateDNSRecord(ip)
	if err != nil {
		c.logger.Printf("Failed to update DNS record: %v", err)
		return fmt.Errorf("failed to update DNS record: %w", err)
	}

	if updated {
		c.logger.Printf("DNS record updated to: %s", ip)
	} else {
		c.logger.Printf("DNS record is up-to-date (IP: %s)", ip)
	}

	return nil
}
