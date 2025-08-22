package bitset

import (
	"math/bits"
)

type Bitset struct {
	data []byte
	size int
}

// CalcDataSize ...
func CalcDataSize(size int) int {
	return (size + 7) / 8
}

func New(size int) *Bitset {
	// Calculate number of bytes needed
	byteSize := CalcDataSize(size)
	return &Bitset{
		data: make([]byte, byteSize),
		size: size,
	}
}

// DataSize returns the size of the underlying byte array
func (b *Bitset) DataSize() int {
	return len(b.data)
}

// Set sets the bit at the given index to 1
func (b *Bitset) Set(index int) {
	if index < 0 || index >= b.size {
		return
	}
	byteIndex := index / 8
	bitIndex := uint(index % 8)
	b.data[byteIndex] |= (1 << bitIndex)
}

// Clear sets the bit at the given index to 0
func (b *Bitset) Clear(index int) {
	if index < 0 || index >= b.size {
		return
	}
	byteIndex := index / 8
	bitIndex := uint(index % 8)
	b.data[byteIndex] &^= (1 << bitIndex)
}

// Test returns true if the bit at the given index is 1
func (b *Bitset) Test(index int) bool {
	if index < 0 || index >= b.size {
		return false
	}
	byteIndex := index / 8
	bitIndex := uint(index % 8)
	return (b.data[byteIndex] & (1 << bitIndex)) != 0
}

// Flip toggles the bit at the given index
func (b *Bitset) Flip(index int) {
	if index < 0 || index >= b.size {
		return
	}
	byteIndex := index / 8
	bitIndex := uint(index % 8)
	b.data[byteIndex] ^= (1 << bitIndex)
}

// CountOnes returns the number of bits set to 1
func (b *Bitset) CountOnes() int {
	count := 0
	for _, byteVal := range b.data {
		count += bits.OnesCount8(byteVal)
	}
	return count
}

// CountOnesRange returns the number of bits set to 1 in a range [start, end)
func (b *Bitset) CountOnesRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.size {
		end = b.size
	}
	if start >= end {
		return 0
	}

	count := 0
	for i := start; i < end; i++ {
		if b.Test(i) {
			count++
		}
	}
	return count
}

// Density returns the ratio of ones to total bits
func (b *Bitset) Density() float64 {
	if b.size == 0 {
		return 0.0
	}
	return float64(b.CountOnes()) / float64(b.size)
}

// Serialize converts the bitset to a byte slice
func (b *Bitset) Serialize() (int, []byte) {
	return b.size, b.data
}

// Deserialize reconstructs a bitset from serialized data
func (b *Bitset) Deserialize(size int, data []byte) {
	b.data = data
	b.size = size
}

// Size returns the total number of bits in the bitset
func (b *Bitset) Size() int {
	return b.size
}

// Bytes returns the underlying byte slice (use with caution)
func (b *Bitset) Bytes() []byte {
	return b.data
}

// String returns a string representation of the bitset
func (b *Bitset) String() string {
	var result string
	for i := 0; i < b.size; i++ {
		if b.Test(i) {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}

// Reset clears all bits in the bitset
func (b *Bitset) Reset() {
	for i := range b.data {
		b.data[i] = 0
	}
}

// Fill sets all bits in the bitset to 1
func (b *Bitset) Fill() {
	for i := range b.data {
		b.data[i] = 0xFF
	}
	// Clear any extra bits in the last byte that exceed size
	extraBits := b.size % 8
	if extraBits > 0 {
		mask := byte((1 << uint(extraBits)) - 1)
		lastIndex := len(b.data) - 1
		b.data[lastIndex] &= mask
	}
}

// All returns true if all bits are set to 1
func (b *Bitset) All() bool {
	// Check full bytes
	for i := 0; i < len(b.data)-1; i++ {
		if b.data[i] != 0xFF {
			return false
		}
	}

	// Check last byte (may be partial)
	if len(b.data) > 0 {
		extraBits := b.size % 8
		if extraBits == 0 {
			return b.data[len(b.data)-1] == 0xFF
		} else {
			mask := byte((1 << uint(extraBits)) - 1)
			return (b.data[len(b.data)-1] & mask) == mask
		}
	}
	return true
}

// Any returns true if any bit is set to 1
func (b *Bitset) Any() bool {
	for _, byteVal := range b.data {
		if byteVal != 0 {
			return true
		}
	}
	return false
}

// None returns true if no bits are set to 1
func (b *Bitset) None() bool {
	return !b.Any()
}

// Clone creates a copy of the bitset
func (b *Bitset) Clone() *Bitset {
	newData := make([]byte, len(b.data))
	copy(newData, b.data)
	return &Bitset{
		data: newData,
		size: b.size,
	}
}

// Equal returns true if two bitsets are identical
func (b *Bitset) Equal(other *Bitset) bool {
	if b.size != other.size {
		return false
	}
	for i := range b.data {
		if b.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

// FromBytes ...
func FromBytes(size int, data []byte) *Bitset {
	b := New(size)
	b.Deserialize(size, data)
	return b
}
