package helpers

import (
	"math/rand"
	"os"
	"os/exec"
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

// ClearTerm runs the "clear" command, clearing out the current
// terminal window where Termify is running.
func ClearTerm() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
