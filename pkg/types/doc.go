// Package types provides common type definitions for Cosmos SDK modules.
//
// This package contains shared interface definitions that are used across
// multiple modules:
//
//   - AccountKeeper interface for account operations
//   - BankKeeper interface for balance and coin operations
//   - ParamSubspace interface for parameter management
//
// These interfaces should be used instead of duplicating interface definitions
// in each module's types/expected_keepers.go file.
//
// Example usage:
//
//	import pkgtypes "sharetoken/pkg/types"
//
//	type Keeper struct {
//	    accountKeeper pkgtypes.AccountKeeper
//	    bankKeeper    pkgtypes.BankKeeper
//	}
//
package types
