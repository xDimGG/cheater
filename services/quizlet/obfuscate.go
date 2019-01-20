package quizlet

import (
	"strconv"
	"strings"
)

const offset = 77

func obfuscate(data []byte) string {
	var result strings.Builder

	for i, char := range data {
		if i != 0 {
			result.WriteByte('-')
		}

		result.WriteString(strconv.Itoa(int(char) + (offset % (i + 1))))
	}

	return result.String()
}
