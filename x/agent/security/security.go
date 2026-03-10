package security

import (
	"fmt"
)

// SecurityLevel represents the security enforcement level
type SecurityLevel int

const (
	LevelMinimal  SecurityLevel = 1
	LevelStandard SecurityLevel = 2
	LevelHigh     SecurityLevel = 3
	LevelParanoid SecurityLevel = 4
)

// SecurityLayer represents a single security mechanism
type SecurityLayer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// SecurityConfig contains all 16 security layers
type SecurityConfig struct {
	Level  SecurityLevel   `json:"level"`
	Layers []SecurityLayer `json:"layers"`
}

// GetSecurityLayers returns the 16 security layers
func GetSecurityLayers() []SecurityLayer {
	return []SecurityLayer{
		{Name: "WASM Sandbox", Description: "Code runs in isolated WASM sandbox", Enabled: true},
		{Name: "Memory Limit", Description: "Restricts memory usage per agent", Enabled: true},
		{Name: "CPU Limit", Description: "Restricts CPU time per execution", Enabled: true},
		{Name: "Network Isolation", Description: "Controls network access", Enabled: true},
		{Name: "File System Sandbox", Description: "Restricted file system access", Enabled: true},
		{Name: "Input Validation", Description: "Validates all inputs", Enabled: true},
		{Name: "Output Sanitization", Description: "Sanitizes all outputs", Enabled: true},
		{Name: "Timeout Enforcement", Description: "Enforces execution timeouts", Enabled: true},
		{Name: "Rate Limiting", Description: "Limits request rate", Enabled: true},
		{Name: "Audit Logging", Description: "Logs all activities", Enabled: true},
		{Name: "Resource Quota", Description: "Enforces resource quotas", Enabled: true},
		{Name: "Code Signing", Description: "Verifies code signatures", Enabled: true},
		{Name: "Dependency Scan", Description: "Scans dependencies", Enabled: true},
		{Name: "Secret Masking", Description: "Masks secrets in logs", Enabled: true},
		{Name: "Execution Trace", Description: "Traces execution flow", Enabled: true},
		{Name: "Emergency Kill", Description: "Emergency stop capability", Enabled: true},
	}
}

// NewSecurityConfig creates a new security configuration
func NewSecurityConfig(level SecurityLevel) *SecurityConfig {
	layers := GetSecurityLayers()

	// Disable some layers for lower security levels
	if level <= LevelMinimal {
		// Disable layers 10-16 for minimal
		for i := 9; i < len(layers); i++ {
			layers[i].Enabled = false
		}
	}
	if level <= LevelStandard {
		// Disable layers 14-16 for standard
		for i := 13; i < len(layers); i++ {
			layers[i].Enabled = false
		}
	}

	return &SecurityConfig{
		Level:  level,
		Layers: layers,
	}
}

// Validate validates the security configuration
func (sc SecurityConfig) Validate() error {
	if sc.Level < LevelMinimal || sc.Level > LevelParanoid {
		return fmt.Errorf("invalid security level: %d", sc.Level)
	}

	// Check at least minimal layers are enabled
	enabledCount := 0
	for _, layer := range sc.Layers {
		if layer.Enabled {
			enabledCount++
		}
	}

	if enabledCount < 9 {
		return fmt.Errorf("at least 9 security layers must be enabled, got %d", enabledCount)
	}

	return nil
}

// GetEnabledLayers returns enabled security layers
func (sc SecurityConfig) GetEnabledLayers() []SecurityLayer {
	var enabled []SecurityLayer
	for _, layer := range sc.Layers {
		if layer.Enabled {
			enabled = append(enabled, layer)
		}
	}
	return enabled
}
