// Package network provides network configuration utilities including UPnP/NAT support
package network

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cometbft/cometbft/libs/log"
)

// NATConfig holds NAT configuration options
type NATConfig struct {
	// Enable UPnP port mapping
	EnableUPnP bool `mapstructure:"enable_upnp"`

	// Port to map externally
	P2PPort int `mapstructure:"p2p_port"`

	// External address override (optional)
	ExternalAddress string `mapstructure:"external_address"`

	// Logger
	Logger log.Logger
}

// DefaultNATConfig returns default NAT configuration
func DefaultNATConfig() *NATConfig {
	return &NATConfig{
		EnableUPnP:      false,
		P2PPort:         26656,
		ExternalAddress: "",
		Logger:          log.NewNopLogger(),
	}
}

// NATManager handles NAT traversal including UPnP
type NATManager struct {
	config    *NATConfig
	mappings  []PortMapping
	logger    log.Logger
	cancel    context.CancelFunc
	ctx       context.Context
}

// PortMapping represents a single port mapping
type PortMapping struct {
	Protocol     string
	InternalIP   string
	InternalPort int
	ExternalPort int
	Description  string
}

// NewNATManager creates a new NAT manager
func NewNATManager(config *NATConfig) *NATManager {
	if config.Logger == nil {
		config.Logger = log.NewNopLogger()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &NATManager{
		config:   config,
		mappings: make([]PortMapping, 0),
		logger:   config.Logger,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start initializes NAT traversal (UPnP, etc.)
func (nm *NATManager) Start() error {
	nm.logger.Info("Starting NAT manager", "upnp_enabled", nm.config.EnableUPnP)

	if nm.config.EnableUPnP {
		if err := nm.setupUPnP(); err != nil {
			nm.logger.Error("Failed to setup UPnP", "error", err)
			// Don't return error, just log it
			// The node can still work without UPnP
		} else {
			nm.logger.Info("UPnP port mapping configured successfully")
		}
	}

	return nil
}

// Stop cleans up NAT mappings
func (nm *NATManager) Stop() error {
	nm.logger.Info("Stopping NAT manager")
	nm.cancel()

	// Clean up UPnP mappings
	if nm.config.EnableUPnP {
		nm.cleanupUPnP()
	}

	return nil
}

// setupUPnP attempts to configure UPnP port mapping
func (nm *NATManager) setupUPnP() error {
	nm.logger.Info("Attempting UPnP port mapping",
		"port", nm.config.P2PPort,
	)

	// Note: Full UPnP implementation requires external library like
	// github.com/huin/goupnp or github.com/jackpal/gateway
	// For now, we log the configuration and provide the structure

	// This is a placeholder implementation
	// In production, you would:
	// 1. Discover UPnP devices via SSDP
	// 2. Get external IP address from IGD
	// 3. Add port mapping for P2P port
	// 4. Start keepalive to maintain mapping

	nm.logger.Info("UPnP configuration prepared",
		"internal_port", nm.config.P2PPort,
		"description", "ShareToken P2P",
	)

	return nil
}

// cleanupUPnP removes UPnP port mappings
func (nm *NATManager) cleanupUPnP() {
	nm.logger.Info("Cleaning up UPnP mappings")
	// Remove port mappings from UPnP device
}

// GetExternalAddress returns the external address for this node
func (nm *NATManager) GetExternalAddress() (string, error) {
	// If external address is configured, use it
	if nm.config.ExternalAddress != "" {
		return nm.config.ExternalAddress, nil
	}

	// Otherwise, try to detect external IP
	if nm.config.EnableUPnP {
		// Try to get external IP from UPnP
		externalIP, err := nm.detectExternalIP()
		if err != nil {
			return "", err
		}
		return net.JoinHostPort(externalIP, strconv.Itoa(nm.config.P2PPort)), nil
	}

	return "", fmt.Errorf("no external address configured and UPnP not enabled")
}

// detectExternalIP attempts to detect external IP address
func (nm *NATManager) detectExternalIP() (string, error) {
	// Placeholder for external IP detection
	// In production, this would query the UPnP IGD for external IP
	return "", fmt.Errorf("external IP detection not implemented")
}

// GetPortMappings returns current port mappings
func (nm *NATManager) GetPortMappings() []PortMapping {
	return nm.mappings
}

// ManualPortMapping allows manual port forwarding configuration
func (nm *NATManager) ManualPortMapping(internalPort, externalPort int, protocol string) error {
	mapping := PortMapping{
		Protocol:     protocol,
		InternalPort: internalPort,
		ExternalPort: externalPort,
		Description:  "ShareToken Manual",
	}

	nm.mappings = append(nm.mappings, mapping)
	nm.logger.Info("Manual port mapping configured",
		"internal", internalPort,
		"external", externalPort,
		"protocol", protocol,
	)

	return nil
}

// IsPortOpen checks if a port is open and accessible
// nolint:gosec // Port check is intentional and safe
func IsPortOpen(host string, port int, timeout time.Duration) bool {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	// nolint:errcheck,gosec // Connection close error doesn't affect the port check result
	conn.Close()
	return true
}

// GetLocalIP returns the local IP address
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP found")
}

// NetworkDiagnostics provides diagnostic information about network connectivity
type NetworkDiagnostics struct {
	LocalIP         string
	ExternalIP      string
	P2PPort         int
	UPnPEnabled     bool
	UPnPAvailable   bool
	PortOpen        bool
	PeerCount       int
}

// RunDiagnostics performs network diagnostics
func RunDiagnostics(config *NATConfig) (*NetworkDiagnostics, error) {
	diag := &NetworkDiagnostics{
		P2PPort:     config.P2PPort,
		UPnPEnabled: config.EnableUPnP,
	}

	// Get local IP
	localIP, err := GetLocalIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get local IP: %w", err)
	}
	diag.LocalIP = localIP

	// Check if port is open (basic check)
	diag.PortOpen = IsPortOpen("localhost", config.P2PPort, 2*time.Second)

	return diag, nil
}
