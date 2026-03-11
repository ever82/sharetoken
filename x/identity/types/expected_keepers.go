package types

import (
	"sharetoken/pkg/types"
)

// AccountKeeper defines the expected account keeper used for simulations
// Deprecated: Use sharetoken/pkg/types.AccountKeeper instead
type AccountKeeper = types.AccountKeeper

// BankKeeper defines the expected interface needed to retrieve account balances
// Deprecated: Use sharetoken/pkg/types.BankKeeper instead
type BankKeeper = types.BankKeeper
