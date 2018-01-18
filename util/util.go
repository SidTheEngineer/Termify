package util

import (
	"math/rand"

	tui "github.com/gizak/termui"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomString returns a random string of length strLen
func GenerateRandomString(strLen int) string {
	b := make([]byte, strLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// IsEmpty checks if byte array(s) is/are empty.
func IsEmpty(arrays ...[]byte) bool {
	for _, arr := range arrays {
		if len(arr) == 0 {
			return true
		}
	}
	return false
}

// IsNil checks if byte array(s) is/are nil.
func IsNil(arrays ...[]byte) bool {
	for _, arr := range arrays {
		if arr == nil {
			return true
		}
	}
	return false
}

// ResetTerminal resets the current ui rows that are being displayed
func ResetTerminal() {
	tui.Body.Rows = tui.Body.Rows[:0]
}
