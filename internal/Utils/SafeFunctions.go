package utils

import (
	"fmt"
	"strings"
)

func SafeTruncateString(s string, maxLen int) string {
	runes := []rune(s)

	// Case 1: The string is too long. Truncate to maxLen and add "..."
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}

	// Case 2: The string is shorter than or equal to maxLen.
	// A truncated string is exactly (maxLen + 3) characters long.
	// We need to pad this short string with spaces so it matches that exact length!
	totalTargetLength := maxLen + 3
	spacesNeeded := totalTargetLength - len(runes)

	return string(runes) + strings.Repeat(" ", spacesNeeded)
}

// Helper to pad strings with spaces for perfect grid alignment
func SafePadText(text string, width int, alignRight bool) string {
	if alignRight {
		// Pads spaces on the left, pushing the text to the right (Great for numbers/time)
		return fmt.Sprintf("%*s", width, text)
	}
	// Pads spaces on the right, keeping text on the left (Great for names/labels)
	return fmt.Sprintf("%-*s", width, text)
}
