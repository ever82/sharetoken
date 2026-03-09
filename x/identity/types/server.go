package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SDKContext alias

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// RegisterIdentity registers a new identity
	RegisterIdentity(ctx sdk.Context, msg *MsgRegisterIdentity) (*MsgRegisterIdentityResponse, error)
	// VerifyIdentity verifies an identity
	VerifyIdentity(ctx sdk.Context, msg *MsgVerifyIdentity) (*MsgVerifyIdentityResponse, error)
	// UpdateLimitConfig updates limit configuration
	UpdateLimitConfig(ctx sdk.Context, msg *MsgUpdateLimitConfig) (*MsgUpdateLimitConfigResponse, error)
	// ResetDailyLimits resets daily limits
	ResetDailyLimits(ctx sdk.Context, msg *MsgResetDailyLimits) (*MsgResetDailyLimitsResponse, error)
	// UpdateParams updates module parameters
	UpdateParams(ctx sdk.Context, msg *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// RegisterMsgServer registers the message server
func RegisterMsgServer(server interface{}, srv MsgServer) {
	// Implementation depends on the server type
	// This is a placeholder
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Params queries the parameters
	Params(ctx sdk.Context, req *QueryParamsRequest) (*QueryParamsResponse, error)
	// Identity queries an identity
	Identity(ctx sdk.Context, req *QueryIdentityRequest) (*QueryIdentityResponse, error)
	// Identities queries all identities
	Identities(ctx sdk.Context, req *QueryIdentitiesRequest) (*QueryIdentitiesResponse, error)
	// LimitConfig queries a limit config
	LimitConfig(ctx sdk.Context, req *QueryLimitConfigRequest) (*QueryLimitConfigResponse, error)
	// IsVerified checks if an address is verified
	IsVerified(ctx sdk.Context, req *QueryIsVerifiedRequest) (*QueryIsVerifiedResponse, error)
}

// RegisterQueryServer registers the query server
func RegisterQueryServer(server interface{}, srv QueryServer) {
	// Implementation depends on the server type
	// This is a placeholder
}

// Response types for Msg service

// MsgRegisterIdentityResponse is the response for RegisterIdentity
type MsgRegisterIdentityResponse struct {
	MerkleRoot string `json:"merkle_root"`
}

// MsgVerifyIdentityResponse is the response for VerifyIdentity
type MsgVerifyIdentityResponse struct {
	IsVerified        bool   `json:"is_verified"`
	UpdatedMerkleRoot string `json:"updated_merkle_root"`
}

// MsgUpdateLimitConfigResponse is the response for UpdateLimitConfig
type MsgUpdateLimitConfigResponse struct{}

// MsgResetDailyLimitsResponse is the response for ResetDailyLimits
type MsgResetDailyLimitsResponse struct {
	ResetCount uint64 `json:"reset_count"`
}

// MsgUpdateParams is the message for updating params
type MsgUpdateParams struct {
	Authority string `json:"authority"`
	Params    Params `json:"params"`
}

// MsgUpdateParamsResponse is the response for UpdateParams
type MsgUpdateParamsResponse struct{}

// Query request types

// QueryParamsRequest is the request for Params query
type QueryParamsRequest struct{}

// QueryParamsResponse is the response for Params query
type QueryParamsResponse struct {
	Params Params `json:"params"`
}

// QueryIdentityRequest is the request for Identity query
type QueryIdentityRequest struct {
	Address string `json:"address"`
}

// QueryIdentityResponse is the response for Identity query
type QueryIdentityResponse struct {
	Identity Identity `json:"identity"`
}

// QueryIdentitiesRequest is the request for Identities query
type QueryIdentitiesRequest struct {
	// Pagination could be added here
}

// QueryIdentitiesResponse is the response for Identities query
type QueryIdentitiesResponse struct {
	Identities []Identity `json:"identities"`
}

// QueryLimitConfigRequest is the request for LimitConfig query
type QueryLimitConfigRequest struct {
	Address string `json:"address"`
}

// QueryLimitConfigResponse is the response for LimitConfig query
type QueryLimitConfigResponse struct {
	LimitConfig LimitConfig `json:"limit_config"`
}

// QueryIsVerifiedRequest is the request for IsVerified query
type QueryIsVerifiedRequest struct {
	Address string `json:"address"`
}

// QueryIsVerifiedResponse is the response for IsVerified query
type QueryIsVerifiedResponse struct {
	IsVerified bool `json:"is_verified"`
}
