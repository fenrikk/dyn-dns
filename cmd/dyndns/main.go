package main

import (
	"log"
	"os"
	"time"

	"github.com/fenrikk/dyn-dns/internal/controller"
	"github.com/fenrikk/dyn-dns/internal/dnsprovider"
	"github.com/fenrikk/dyn-dns/internal/iplocator"
)

func main() {
	// Create logger
	logger := log.New(os.Stdout, "[DYNDNS] ", log.LstdFlags)
	logger.Println("Initializing dynamic DNS service")

	// Create IP locator
	ipLocator := iplocator.NewIPifyLocator()

	// Create DNS provider
	dnsProvider, err := dnsprovider.NewRoute53Provider(dnsprovider.Route53ProviderOptions{
		HostedZoneID: os.Getenv("AWS_HOSTED_ZONE_ID"),
		RecordName:   os.Getenv("DNS_RECORD_NAME"),
		TTL:          300,
	})
	if err != nil {
		logger.Fatalf("Failed to create DNS provider: %v", err)
	}

	// Get check interval from environment variable
	var interval time.Duration
	interval_env := os.Getenv("CHECK_INTERVAL")
	if interval_env != "" {
		var err error
		interval, err = time.ParseDuration(interval_env)
		if err != nil {
			logger.Fatalf("Failed to parse CHECK_INTERVAL: %v", err)
		}
	}

	// Create and start controller
	dyndnsController := controller.NewController(
		dnsProvider,
		ipLocator,
		controller.ControllerOptions{
			CheckInterval: interval,
			Logger:        logger,
		},
	)

	// Start controller
	if err := dyndnsController.Start(); err != nil {
		logger.Fatalf("Controller error: %v", err)
	}
}
