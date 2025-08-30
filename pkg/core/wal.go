package core

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/hasssanezzz/go-trgm/pkg/common"
)

type WAL struct {
	filepath string
	file     *os.File
}

func NewWAL(filepath string) (*WAL, error) {
	w := &WAL{
		filepath: filepath,
	}

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %v", filepath, err)
	}
	w.file = file

	return w, nil
}

func (w *WAL) Append(tri uint32, docID common.DocumentID) error {
	record := [20]byte{}
	copy(record[:4], binary.LittleEndian.AppendUint32(nil, tri))
	copy(record[4:], docID[:])
	if _, err := w.file.Write(record[:]); err != nil {
		return err
	}
	return nil
}

func (w *WAL) Parse() (map[uint32][]common.DocumentID, error) {
	file, err := os.Open(w.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r io.Reader = bufio.NewReader(file)

	result := make(map[uint32][]common.DocumentID, 2048)
	for {
		record := [20]byte{}
		if _, err := r.Read(record[:]); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		tri := binary.LittleEndian.Uint32(record[:4])
		if _, ok := result[tri]; !ok {
			result[tri] = make([]common.DocumentID, 0, 512)
		}
		var docID common.DocumentID
		copy(docID[:], record[4:])
		result[tri] = append(result[tri], docID)
	}

	return result, nil
}

func (w *WAL) Clear() error {
	if err := w.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to clear WAL file %q: %v", w.filepath, err)
	}

	if _, err := w.file.Seek(0, 0); err != nil {
		// If seeking fails after truncation, the file state might be inconsistent
		// for the next write, so report the error.
		return fmt.Errorf("failed to seek to beginning of WAL file %q after clearing: %v", w.filepath, err)
	}

	return nil
}
