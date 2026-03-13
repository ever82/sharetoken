// Package store provides utilities for working with Cosmos SDK KV stores.
//
// This package contains functions for building and parsing store keys:
//
//   - Key building with various data types (string, uint64, int64, bytes)
//   - Composite key building for multi-component keys
//   - Key parsing to extract values from store keys
//   - Prefix range helpers for iterators
//
// All numeric encodings use big-endian format for consistent ordering.
//
// Example usage:
//
//	import "sharetoken/pkg/store"
//
//	// Building keys
//	key := store.BuildKey([]byte{0x01}, "myid")
//	key := store.BuildCompositeKey([]byte{0x02}, "owner", "id")
//	key := store.BuildKeyWithUint64([]byte{0x03}, 12345)
//
//	// Parsing keys
//	id := store.ParseKey(key, 1)  // skip 1-byte prefix
//	value := store.ParseUint64Key(key, 1)
//
//	// Prefix iteration
//	start, end := store.PrefixRange([]byte{0x01})
//	iterator := sdk.KVStorePrefixIterator(store, []byte{0x01})
package store
