package trgm

import (
	"fmt"
	"strings"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyz0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

var CharToCode [128]int

func init() {
	for i := range CharToCode {
		CharToCode[i] = -1
	}
	for i, ch := range Alphabet {
		CharToCode[int(ch)] = i
	}
}

func triToInt(trigram string) uint32 {
	c0 := CharToCode[trigram[0]]
	c1 := CharToCode[trigram[1]]
	c2 := CharToCode[trigram[2]]

	if c0 == -1 || c1 == -1 || c2 == -1 {
		panic("triToInt: bad pattern received: " + trigram)
	}

	return uint32(c0*68*68 + c1*68 + c2) // 8836 = 94*94
}

func intToTri(id uint32) string {
	if id >= TriMaxCount {
		panic(fmt.Sprintf("triToInt: bad id received: %d", id))
	}
	c0 := id / (68 * 68)
	c1 := (id / 68) % 68
	c2 := id % 68

	return string([]byte{
		Alphabet[c0],
		Alphabet[c1],
		Alphabet[c2],
	})
}

func extractTrigrams(s string) []uint32 {
	s = strings.ToLower(s)
	results := make([]uint32, len(s)-2)
	for i := range len(s) - 2 {
		tri := s[i : i+3]
		results[i] = triToInt(tri)
	}
	return results
}
