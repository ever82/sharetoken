package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSet implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVerificationRequired, &p.VerificationRequired, validateVerificationRequired),
		paramtypes.NewParamSetPair(KeyAllowedProviders, &p.AllowedProviders, validateAllowedProviders),
	}
}

// NewParams creates a new Params instance
func NewParams(verificationRequired bool, allowedProviders []string) Params {
	return Params{
		VerificationRequired: verificationRequired,
		AllowedProviders:     allowedProviders,
	}
}

var (
	KeyVerificationRequired = []byte("VerificationRequired")
	KeyAllowedProviders     = []byte("AllowedProviders")
)

func validateVerificationRequired(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateAllowedProviders(i interface{}) error {
	providers, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(providers) == 0 {
		return fmt.Errorf("allowed providers cannot be empty")
	}
	return nil
}

// Identity represents a user's on-chain identity
type Identity struct {
	Address              string `json:"address"`
	DID                  string `json:"did"`
	VerificationHash     string `json:"verification_hash"`
	VerificationProvider string `json:"verification_provider"`
	RegistrationTime     int64  `json:"registration_time"`
	IsVerified           bool   `json:"is_verified"`
	MerkleRoot           string `json:"merkle_root"`
	MetadataHash         string `json:"metadata_hash"`
}

// NewIdentity creates a new identity
func NewIdentity(address, did string) *Identity {
	return &Identity{
		Address:          address,
		DID:              did,
		RegistrationTime: time.Now().Unix(),
		IsVerified:       false,
	}
}

// ValidateBasic performs basic validation of identity fields
func (i Identity) ValidateBasic() error {
	if i.Address == "" {
		return ErrInvalidAddress
	}

	// Validate address format
	_, err := sdk.AccAddressFromBech32(i.Address)
	if err != nil {
		return ErrInvalidAddress.Wrap(err.Error())
	}

	// Validate DID if provided
	if i.DID != "" && !isValidDID(i.DID) {
		return ErrInvalidDID.Wrap(i.DID)
	}

	return nil
}

// GenerateMerkleRoot generates a merkle root from identity data
func (i *Identity) GenerateMerkleRoot() string {
	// Create data hash
	data := fmt.Sprintf("%s:%s:%s:%s:%d",
		i.Address,
		i.DID,
		i.VerificationHash,
		i.VerificationProvider,
		i.RegistrationTime,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// VerifyIdentityProof verifies an identity proof (simplified)
func (i *Identity) VerifyIdentityProof(proof []byte, leafIndex int) bool {
	// Simplified proof verification
	// In production, this would use a proper merkle tree implementation
	expectedHash := i.GenerateMerkleRoot()
	return hex.EncodeToString(proof) == expectedHash
}

// GetVerificationHash generates a hash of verification data
func GetVerificationHash(provider, providerID, timestamp string) string {
	data := fmt.Sprintf("%s:%s:%s", provider, providerID, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// isValidDID checks if a string is a valid DID format
func isValidDID(did string) bool {
	// Basic DID format: did:method:specific-identifier
	// Example: did:sharetoken:abc123...
	if len(did) < 7 {
		return false
	}
	return did[:4] == "did:"
}

// DefaultVerificationProviders is the list of supported verification providers
var DefaultVerificationProviders = []string{"wechat", "github", "google"}

// IsValidProvider checks if a provider is valid
func IsValidProvider(provider string) bool {
	for _, p := range DefaultVerificationProviders {
		if p == provider {
			return true
		}
	}
	return false
}

// GenesisState defines the identity module's genesis state
type GenesisState struct {
	Params       Params        `json:"params"`
	Identities   []Identity    `json:"identities"`
	LimitConfigs []LimitConfig `json:"limit_configs"`
}

// Params defines the parameters for the identity module
type Params struct {
	VerificationRequired bool               `json:"verification_required"`
	AllowedProviders     []string           `json:"allowed_providers"`
	DefaultLimits        DefaultLimitConfig `json:"default_limits"`
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		VerificationRequired: false,
		AllowedProviders:     DefaultVerificationProviders,
		DefaultLimits:        DefaultDefaultLimitConfig(),
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if len(p.AllowedProviders) == 0 {
		return fmt.Errorf("allowed_providers cannot be empty")
	}
	return nil
}

// String implements stringer interface
func (p Params) String() string {
	return fmt.Sprintf(`
Params:
  Verification Required: %v
  Allowed Providers: %v
`, p.VerificationRequired, p.AllowedProviders)
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		Identities:   []Identity{},
		LimitConfigs: []LimitConfig{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	// Check for duplicate addresses
	addressMap := make(map[string]bool)
	for _, identity := range data.Identities {
		if addressMap[identity.Address] {
			return fmt.Errorf("duplicate identity address: %s", identity.Address)
		}
		addressMap[identity.Address] = true

		if err := identity.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
