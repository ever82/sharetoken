// Package network provides network configuration utilities including UPnP/NAT support
package network

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/huin/goupnp/dcps/internetgateway2"
)

// NATConfig holds NAT configuration options
// nolint:govet // fieldalignment: struct field order is for readability
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
// nolint:govet // fieldalignment: struct field order is for readability
type NATManager struct {
	config     *NATConfig
	mappings   []PortMapping
	logger     log.Logger
	cancel     context.CancelFunc
	ctx        context.Context
	upnpClient UPnPClient
	upnpMu     sync.RWMutex
}

// UPnPClient interface for UPnP operations
type UPnPClient interface {
	GetExternalIPAddress() (string, error)
	AddPortMapping(protocol string, externalPort, internalPort uint16, internalIP string, enabled bool, description string, lease uint32) error
	DeletePortMapping(protocol string, externalPort uint16) error
}

// PortMapping represents a single port mapping
// nolint:govet // fieldalignment: struct field order is for readability
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
			// Start keepalive goroutine
			go nm.keepalive()
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

	// Discover UPnP IGD devices
	client, err := nm.discoverUPnP()
	if err != nil {
		return fmt.Errorf("failed to discover UPnP device: %w", err)
	}

	nm.upnpMu.Lock()
	nm.upnpClient = client
	nm.upnpMu.Unlock()

	// Get local IP
	localIP, err := GetLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP: %w", err)
	}

	// Get external IP
	externalIP, err := client.GetExternalIPAddress()
	if err != nil {
		nm.logger.Error("Failed to get external IP, continuing anyway", "error", err)
		externalIP = "unknown"
	}

	nm.logger.Info("UPnP device discovered",
		"external_ip", externalIP,
		"local_ip", localIP,
	)

	// Add port mapping
	protocol := "TCP"
	externalPort := uint16(nm.config.P2PPort)
	internalPort := uint16(nm.config.P2PPort)
	description := "ShareToken P2P"
	leaseDuration := uint32(3600) // 1 hour lease, will be renewed by keepalive

	err = client.AddPortMapping(
		protocol,
		externalPort,
		internalPort,
		localIP,
		true,
		description,
		leaseDuration,
	)
	if err != nil {
		return fmt.Errorf("failed to add port mapping: %w", err)
	}

	// Record the mapping
	mapping := PortMapping{
		Protocol:     protocol,
		InternalIP:   localIP,
		InternalPort: nm.config.P2PPort,
		ExternalPort: nm.config.P2PPort,
		Description:  description,
	}
	nm.mappings = append(nm.mappings, mapping)

	nm.logger.Info("UPnP port mapping configured",
		"protocol", protocol,
		"external_port", externalPort,
		"internal_port", internalPort,
		"external_ip", externalIP,
	)

	return nil
}

// discoverUPnP discovers UPnP IGD devices
func (nm *NATManager) discoverUPnP() (UPnPClient, error) {
	// Try IGD v2 first
	nm.logger.Debug("Searching for UPnP IGD v2...")
	clients, _, err := internetgateway2.NewWANIPConnection2Clients()
	if err == nil && len(clients) > 0 {
		nm.logger.Info("Found UPnP IGD v2")
		return &igd2Client{client: clients[0]}, nil
	}

	// Try IGD v1
	nm.logger.Debug("Searching for UPnP IGD v1...")
	clients1, _, err := internetgateway1.NewWANIPConnection1Clients()
	if err == nil && len(clients1) > 0 {
		nm.logger.Info("Found UPnP IGD v1")
		return &igd1Client{client: clients1[0]}, nil
	}

	// Try WAN PPP Connection v1 as fallback
	pppClients, _, err := internetgateway1.NewWANPPPConnection1Clients()
	if err == nil && len(pppClients) > 0 {
		nm.logger.Info("Found UPnP WAN PPP Connection")
		return &igd1PPPClient{client: pppClients[0]}, nil
	}

	return nil, fmt.Errorf("no UPnP IGD device found")
}

// cleanupUPnP removes UPnP port mappings
func (nm *NATManager) cleanupUPnP() {
	nm.logger.Info("Cleaning up UPnP mappings")

	nm.upnpMu.RLock()
	client := nm.upnpClient
	nm.upnpMu.RUnlock()

	if client == nil {
		return
	}

	for _, mapping := range nm.mappings {
		err := client.DeletePortMapping(mapping.Protocol, uint16(mapping.ExternalPort))
		if err != nil {
			nm.logger.Error("Failed to delete port mapping",
				"protocol", mapping.Protocol,
				"external_port", mapping.ExternalPort,
				"error", err,
			)
		} else {
			nm.logger.Info("Port mapping removed",
				"protocol", mapping.Protocol,
				"external_port", mapping.ExternalPort,
			)
		}
	}
}

// keepalive periodically renews UPnP port mappings
func (nm *NATManager) keepalive() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-nm.ctx.Done():
			return
		case <-ticker.C:
			nm.upnpMu.RLock()
			client := nm.upnpClient
			nm.upnpMu.RUnlock()

			if client == nil {
				continue
			}

			// Renew mappings
			for _, mapping := range nm.mappings {
				err := client.AddPortMapping(
					mapping.Protocol,
					uint16(mapping.ExternalPort),
					uint16(mapping.InternalPort),
					mapping.InternalIP,
					true,
					mapping.Description,
					3600, // 1 hour lease
				)
				if err != nil {
					nm.logger.Error("Failed to renew port mapping",
						"protocol", mapping.Protocol,
						"external_port", mapping.ExternalPort,
						"error", err,
					)
				} else {
					nm.logger.Debug("Port mapping renewed",
						"protocol", mapping.Protocol,
						"external_port", mapping.ExternalPort,
					)
				}
			}
		}
	}
}

// GetExternalAddress returns the external address for this node
func (nm *NATManager) GetExternalAddress() (string, error) {
	// If external address is configured, use it
	if nm.config.ExternalAddress != "" {
		return nm.config.ExternalAddress, nil
	}

	// Otherwise, try to detect external IP via UPnP
	nm.upnpMu.RLock()
	client := nm.upnpClient
	nm.upnpMu.RUnlock()

	if client != nil {
		externalIP, err := client.GetExternalIPAddress()
		if err == nil && externalIP != "" && externalIP != "0.0.0.0" {
			return net.JoinHostPort(externalIP, strconv.Itoa(nm.config.P2PPort)), nil
		}
	}

	return "", fmt.Errorf("no external address configured and UPnP not available")
}

// detectExternalIP attempts to detect external IP address
// nolint:unused // Kept for future use
func (nm *NATManager) detectExternalIP() (string, error) {
	nm.upnpMu.RLock()
	client := nm.upnpClient
	nm.upnpMu.RUnlock()

	if client != nil {
		return client.GetExternalIPAddress()
	}

	return "", fmt.Errorf("UPnP client not available")
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
// nolint:govet // fieldalignment: struct field order is for readability
type NetworkDiagnostics struct {
	LocalIP       string
	ExternalIP    string
	P2PPort       int
	UPnPEnabled   bool
	UPnPAvailable bool
	PortOpen      bool
	PeerCount     int
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

	// Try to get external IP if UPnP is enabled
	if config.EnableUPnP {
		// Quick UPnP discovery
		_, err := goupnp.DiscoverDevices(internetgateway1.URN_WANIPConnection_1)
		if err == nil {
			diag.UPnPAvailable = true
		}
	}

	return diag, nil
}

// igd1Client wraps internetgateway1.WANIPConnection1 for UPnPClient interface
type igd1Client struct {
	client *internetgateway1.WANIPConnection1
}

func (c *igd1Client) GetExternalIPAddress() (string, error) {
	return c.client.GetExternalIPAddress()
}

func (c *igd1Client) AddPortMapping(protocol string, externalPort, internalPort uint16, internalIP string, enabled bool, description string, lease uint32) error {
	// IGD v1: NewRemoteHost, ExternalPort, Protocol, InternalPort, InternalClient, Enabled, Description, LeaseDuration
	return c.client.AddPortMapping("", externalPort, protocol, internalPort, internalIP, enabled, description, lease)
}

func (c *igd1Client) DeletePortMapping(protocol string, externalPort uint16) error {
	// IGD v1: NewRemoteHost, ExternalPort, Protocol
	return c.client.DeletePortMapping("", externalPort, protocol)
}

// igd2Client wraps internetgateway2.WANIPConnection2 for UPnPClient interface
type igd2Client struct {
	client *internetgateway2.WANIPConnection2
}

func (c *igd2Client) GetExternalIPAddress() (string, error) {
	return c.client.GetExternalIPAddress()
}

func (c *igd2Client) AddPortMapping(protocol string, externalPort, internalPort uint16, internalIP string, enabled bool, description string, lease uint32) error {
	// IGD v2: NewRemoteHost, ExternalPort, Protocol, InternalPort, InternalClient, Enabled, Description, LeaseDuration
	return c.client.AddPortMapping("", externalPort, protocol, internalPort, internalIP, enabled, description, lease)
}

func (c *igd2Client) DeletePortMapping(protocol string, externalPort uint16) error {
	// IGD v2: NewRemoteHost, ExternalPort, Protocol
	return c.client.DeletePortMapping("", externalPort, protocol)
}

// igd1PPPClient wraps internetgateway1.WANPPPConnection1 for UPnPClient interface
type igd1PPPClient struct {
	client *internetgateway1.WANPPPConnection1
}

func (c *igd1PPPClient) GetExternalIPAddress() (string, error) {
	return c.client.GetExternalIPAddress()
}

func (c *igd1PPPClient) AddPortMapping(protocol string, externalPort, internalPort uint16, internalIP string, enabled bool, description string, lease uint32) error {
	// PPP: NewRemoteHost, ExternalPort, Protocol, InternalPort, InternalClient, Enabled, Description, LeaseDuration
	return c.client.AddPortMapping("", externalPort, protocol, internalPort, internalIP, enabled, description, lease)
}

func (c *igd1PPPClient) DeletePortMapping(protocol string, externalPort uint16) error {
	// PPP: NewRemoteHost, ExternalPort, Protocol
	return c.client.DeletePortMapping("", externalPort, protocol)
}
