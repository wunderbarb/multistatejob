// V0.10.0
// Author: Diehl E.
// (C) Dec 2023

package go_test

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Rng is a randomly seeded random number generator that can be used for tests.
// The random number generator is not cryptographically safe.
var Rng *rand.Rand

// init initializes the random number generator.
func init() {
	Rng = rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404  It is not crypto secure. OK for test

}

// AlphaNumType represents the kind of characters
// that will be generated.
type AlphaNumType int

const (
	// All requests all the characters from the character set
	// [A..Za..z0.9][ @.!$&_+-:;*?#/\\,()[]{}<>%"]
	All AlphaNumType = iota
	// AllCVS requests the same characters as 'All' at the exception to [;,].  It is to be
	// used when applying to CVS files.
	AllCVS
	// AlphaNum requests only characters that are alphanumerical with space included.
	AlphaNum
	// AlphaNumNoSpace requests only characters that are alphanumerical without space.
	AlphaNumNoSpace
	// Alpha requests only characters that are alphabetical with space included.
	Alpha
	// AlphaNoSpace requests only characters that are alphabetical without space.
	AlphaNoSpace
	// Caps requests only upper characters without space.
	Caps
	// Small requests only minor characters without space.
	Small
	emailBody
)

// RandomEIDR generates a new EIDR.  The check digit is not computed.
// It is concurrent safe.
func RandomEIDR() string {
	const eidrContentID = "10.5240/"
	b := randomSlice(8)
	return fmt.Sprintf(eidrContentID+"%x%x-%x%x-%x%x-%x%x-c", b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7])
}

// RandomEmail returns a random email address.  If ext is not empty, it is used it as the extension
func RandomEmail(ext ...string) string {
	s := RandomAlphaString(12, emailBody) + RandomAlphaString(1, AlphaNumNoSpace) + "@" +
		RandomAlphaString(6, AlphaNumNoSpace) + RandomAlphaString(1, AlphaNumNoSpace) + "."
	if len(ext) == 0 {
		return s + RandomAlphaString(3, AlphaNoSpace)
	}

	return s + strings.TrimPrefix(ext[0], ".")
}

// RandomID returns a random 16-character, alphanumeric, ID.
func RandomID() string {
	const sizeID = 16
	return RandomAlphaString(sizeID, AlphaNumNoSpace)
}

// RandomName returns a random string with size characters.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// The character set is [A..Z][a..z].
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// go-crypto package with GenerateNewKey
func RandomName(size int) string {
	return RandomAlphaString(size, AlphaNoSpace)
}

// RandomString returns a random string with size characters.
// If size is null, then the length of the string is random in the range
// 1 to 256 characters.
//
// The character set is [A..Za..z0..9][ @.!$&_+-:;*?#/\\,()[]{}<>%]
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// go-crypto package with GenerateNewKey
func RandomString(size int) string {

	return RandomAlphaString(size, All)
}

// RandomAlphaString generates a size-character random string which character
// set depends on the value of t.  if t is not a proper value, the returned value
// is the empty string.
// If size is zero or negative, then the length of the string is random in the range
// 1 to 256 characters.
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// go-crypto package with GenerateNewKey
func RandomAlphaString(size int, t AlphaNumType) string {

	s := randomAlphaString(size, t)

	return s
}

// randomAlphaString generates a size-character random string which character
// set depends on the value of t.  if t is not a proper value, the returned value
// is the empty string.
// If size is zero or negative, then the length of the string is random in the range
// 1 to 256 characters.
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// github.com/TechDev-SPE/go-crypto package with GenerateNewKey
func randomAlphaString(size int, t AlphaNumType) string {
	const size0 = 256 // max number of bytes for random set.
	const (
		sc = "abcdefghijklmnopqrstuvwxyz"
		ca = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		nu = "01234567890"
	)

	conv := map[AlphaNumType][]byte{
		All:             []byte(ca + sc + nu + " " + "@.!$&_+-:;*?#/\\,()[]{}<>%\""),
		AllCVS:          []byte(ca + sc + nu + " " + "@.!$&_+-:*?#/()[]{}<>%"),
		AlphaNum:        []byte(ca + sc + nu + " "),
		AlphaNumNoSpace: []byte(ca + sc + nu),
		Alpha:           []byte(ca + sc + " "),
		AlphaNoSpace:    []byte(ca + sc),
		Caps:            []byte(ca),
		Small:           []byte(sc),
		emailBody:       []byte(ca + sc + nu + "._-"),
	}
	if size <= 0 {
		size = Rng.Intn(size0) + 1
	}
	buffer := make([]byte, size)
	choice, ok := conv[t]
	if !ok {
		return ""
	}
	choiceSize := len(choice)
	for i := 0; i < size; i++ {
		// generates the characters
		buffer[i] = choice[Rng.Intn(choiceSize)]
	}
	return string(buffer)
}

// randomSlice returns a random slice with size bytes.
// If size is zero or negative, then the number of bytes in the slice is random in the range
// 1 to 256 characters.
//
// CAUTION: the randomness is not cryptographically secure, thus it should
// not be used for generating keys.  Secure keys are generated using
// github.com/TechDev-SPE/go-crypto package with GenerateNewKey
func randomSlice(size int) []byte {
	const size0 = 256 // max number of bytes for random set.
	if size <= 0 {
		size = Rng.Intn(size0) + 1
	}
	buffer := make([]byte, size)
	_, _ = Rng.Read(buffer)
	return buffer
}
