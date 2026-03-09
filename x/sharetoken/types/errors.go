package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sharetoken module sentinel errors
// nolint:staticcheck // sdkerrors.Register is deprecated but used for compatibility
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
)
