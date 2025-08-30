package common

import "encoding/binary"

const DocumentIdSize = 16

type DocumentID [DocumentIdSize]byte

func (d DocumentID) Compare(other DocumentID) bool {
	for i := range DocumentIdSize {
		if d[i] == other[i] {
			return false
		}
	}
	return true
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
