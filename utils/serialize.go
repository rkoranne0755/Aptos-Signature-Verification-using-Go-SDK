package utils

import (
	"bytes"
	"encoding/binary"

	"github.com/aptos-labs/aptos-go-sdk"
)

// BCSSerializer is a helper struct for serializing data into BCS (Binary Canonical Serialization) format.
type BCSSerializer struct {
	buffer bytes.Buffer
}

// MoveMessageStruct represents a structure to be serialized in the Move programming language.
type MoveMessageStruct struct {
	From   aptos.AccountAddress
	Amount uint64
	Nonce  uint8
}

// Serialize converts the MoveMessageStruct into a byte array using the BCSSerializer.
func (m *MoveMessageStruct) Serialize() []byte {
	serializer := NewBCSSerializer()
	serializer.SerializeAccountAddress(m.From)
	serializer.SerializeU64(uint64(m.Amount))
	serializer.SerializeU8(uint8(m.Nonce))
	return serializer.GetBytes()
}

// SerializeAccountAddress serializes an account address into BCS format.
func (s *BCSSerializer) SerializeAccountAddress(addr aptos.AccountAddress) {
	s.buffer.Write(addr[:])
}

// SerializeU64 serializes a 64-bit unsigned integer into BCS format.
func (s *BCSSerializer) SerializeU64(value uint64) {
	binary.Write(&s.buffer, binary.LittleEndian, uint64(value))
}

// SerializeU8 serializes an 8-bit unsigned integer into BCS format.
func (s *BCSSerializer) SerializeU8(value uint8) {
	binary.Write(&s.buffer, binary.LittleEndian, uint8(value))
}

// NewBCSSerializer creates and returns a new instance of BCSSerializer.
func NewBCSSerializer() *BCSSerializer {
	return &BCSSerializer{}
}

// GetBytes returns the serialized byte array from the BCSSerializer.
func (s *BCSSerializer) GetBytes() []byte {
	return s.buffer.Bytes()
}
