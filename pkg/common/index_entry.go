package common

import (
	"bytes"
	"encoding/binary"
)

const DocumentIdSize = 16

type DocumentID [DocumentIdSize]byte

func (d DocumentID) Equal(other DocumentID) bool {
	return bytes.Equal(d[:], other[:])
}

type IndexEntry struct {
	batchId uint16
	offset  uint32
}

func NewIndexEntry(batchId uint16, offset uint32) IndexEntry {
	return IndexEntry{batchId, offset}
}

func (ie *IndexEntry) Encode() []byte {
	return append(binary.LittleEndian.AppendUint16(nil, ie.batchId), binary.LittleEndian.AppendUint32(nil, ie.offset)...)
}
