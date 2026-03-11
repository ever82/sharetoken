package store

import (
	"encoding/binary"
)

// BuildKey builds a store key from a prefix and a string ID.
// This is the most common pattern for building store keys.
func BuildKey(prefix []byte, id string) []byte {
	return append(prefix, []byte(id)...)
}

// BuildCompositeKey builds a store key from a prefix and multiple string components.
// Useful for composite keys like "prefix|owner|id".
func BuildCompositeKey(prefix []byte, components ...string) []byte {
	key := make([]byte, len(prefix))
	copy(key, prefix)
	for _, component := range components {
		key = append(key, []byte(component)...)
	}
	return key
}

// BuildKeyWithUint64 builds a store key from a prefix and a uint64 ID.
// The uint64 is encoded in big-endian format.
func BuildKeyWithUint64(prefix []byte, id uint64) []byte {
	key := make([]byte, len(prefix)+8)
	copy(key, prefix)
	binary.BigEndian.PutUint64(key[len(prefix):], id)
	return key
}

// BuildKeyWithInt64 builds a store key from a prefix and an int64 ID.
// The int64 is encoded in big-endian format (using uint64 representation).
func BuildKeyWithInt64(prefix []byte, id int64) []byte {
	return BuildKeyWithUint64(prefix, uint64(id))
}

// BuildKeyWithUint32 builds a store key from a prefix and a uint32 ID.
// The uint32 is encoded in big-endian format.
func BuildKeyWithUint32(prefix []byte, id uint32) []byte {
	key := make([]byte, len(prefix)+4)
	copy(key, prefix)
	binary.BigEndian.PutUint32(key[len(prefix):], id)
	return key
}

// BuildKeyWithBytes builds a store key from a prefix and a byte slice.
func BuildKeyWithBytes(prefix []byte, id []byte) []byte {
	return append(prefix, id...)
}

// ParseKey extracts the ID from a store key by removing the prefix.
// Returns the ID portion after the prefix.
func ParseKey(key []byte, prefixLen int) string {
	if len(key) <= prefixLen {
		return ""
	}
	return string(key[prefixLen:])
}

// ParseUint64Key extracts a uint64 ID from a store key.
// The uint64 is decoded from big-endian format.
func ParseUint64Key(key []byte, prefixLen int) uint64 {
	if len(key) < prefixLen+8 {
		return 0
	}
	return binary.BigEndian.Uint64(key[prefixLen:])
}

// PrefixRange returns the start and end keys for iterating over a prefix range.
// This is useful for creating iterators over a specific key prefix.
func PrefixRange(prefix []byte) (start []byte, end []byte) {
	start = prefix
	end = PrefixEndBytes(prefix)
	return start, end
}

// PrefixEndBytes returns the end bytes for a prefix range.
// This creates a key that is just after all keys with the given prefix.
func PrefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}

	end := make([]byte, len(prefix))
	copy(end, prefix)

	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < byte(0xff) {
			end[i]++
			return end[:i+1]
		}
	}

	// The input is all 0xff bytes, so there's no prefix end
	return nil
}

// AppendKey appends a string key component to an existing key.
func AppendKey(key []byte, component string) []byte {
	return append(key, []byte(component)...)
}

// AppendKeyBytes appends a byte slice component to an existing key.
func AppendKeyBytes(key []byte, component []byte) []byte {
	return append(key, component...)
}

// AppendUint64 appends a uint64 to an existing key using big-endian encoding.
func AppendUint64(key []byte, value uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	return append(key, buf...)
}
