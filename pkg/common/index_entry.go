package common

import "encoding/binary"

type entry interface {
	IsDeleted() bool
	Decode([]byte)
	Encode() []byte
	Size() uint32
}

type IndexEntry struct {
	batchId uint16
	offset  uint32
}

func NewIndexEntry(batchId uint16, offset uint32) IndexEntry {
	return IndexEntry{batchId, offset}
}

func (ie *IndexEntry) IsDeleted() bool {
	return ie.batchId == 0
}

func (ie *IndexEntry) Decode(buff []byte) {
	ie.batchId = binary.LittleEndian.Uint16(buff[:2])
	ie.offset = binary.LittleEndian.Uint32(buff[2:])
}

func (ie *IndexEntry) Encode() []byte {
	return append(binary.LittleEndian.AppendUint16(nil, ie.batchId), binary.LittleEndian.AppendUint32(nil, ie.offset)...)
}

func (ie *IndexEntry) Size() uint32 {
	return 6 // (16 + 32) / 8 = 6 bytes
}
