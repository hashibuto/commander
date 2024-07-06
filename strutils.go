package commander

import "strings"

func PadRight(text string, width int) string {
	if len(text) > width {
		text = text[:width-3] + "..."
	}

	return text + strings.Repeat(" ", width-len(text))
}
