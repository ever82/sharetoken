// Package keeper provides common keeper utilities for Cosmos SDK modules.
//
// This package contains reusable components to reduce boilerplate code
// across different module keepers:
//
//   - CRUD operations for KV store (crud.go)
//   - Logger helpers (logger.go)
//   - Message server utilities (msgserver.go)
//
// Example usage:
//
//	import "sharetoken/pkg/keeper"
//
//	// Using generic CRUD keeper
//	type MyItem struct {
//	    ID   string
//	    Name string
//	}
//
//	func (i MyItem) GetID() string { return i.ID }
//
//	crudKeeper := keeper.NewCRUDKeeper[MyItem](cdc, storeKey, []byte{0x01})
//	crudKeeper.Set(ctx, item)
//	item, found := crudKeeper.Get(ctx, "id")
//
//	// Using logger helper
//	logger := keeper.NewLoggerFunc("mymodule")(ctx)
//
//	// Using message server utilities
//	sdkCtx := keeper.UnwrapContext(ctx)
package keeper
