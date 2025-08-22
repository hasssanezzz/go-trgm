package trgm

import (
	"errors"
	"fmt"
	"strings"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyz0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
const TriMaxCount = 69 * 69 * 69

var CharToCode [128]int
var ErrInvalidInput error = errors.New("invalid input")
var AlphaMap = map[rune]struct{}{}

func init() {
	for i := range CharToCode {
		CharToCode[i] = -1
	}
	for i, ch := range Alphabet {
		CharToCode[int(ch)] = i
		AlphaMap[ch] = struct{}{}
	}
}

func isValidString(_ string) error {
	// NOTE: commented out for now
	// should iterate through the ascii table instead of using strings.ToLower

	// for _, c := range input {
	// 	if _, ok := AlphaMap[c]; !ok {
	// 		return ErrInvalidInput
	// 	}
	// }
	return nil
}

func triToInt(trigram string) uint32 {
	c0 := CharToCode[trigram[0]]
	c1 := CharToCode[trigram[1]]
	c2 := CharToCode[trigram[2]]

	if c0 == -1 || c1 == -1 || c2 == -1 {
		panic("triToInt: bad pattern received: " + trigram)
	}

	return uint32(c0*69*69 + c1*69 + c2)
}

func intToTri(id uint32) string {
	if id >= TriMaxCount {
		panic(fmt.Sprintf("triToInt: bad id received: %d", id))
	}
	c0 := id / (69 * 69)
	c1 := (id / 69) % 69
	c2 := id % 69

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

func validateInputAndExtractTrigrams(s string) ([]uint32, error) {
	if err := isValidString(s); err != nil {
		return nil, err
	}

	return extractTrigrams(s), nil
}
