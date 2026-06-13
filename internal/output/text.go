package output

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var tagPattern = regexp.MustCompile(`<[^>]+>`)

func StripHTML(input string) string {
	input = tagPattern.ReplaceAllString(input, "")
	input = html.UnescapeString(input)
	return strings.TrimSpace(strings.Join(strings.Fields(input), " "))
}

func AnyID(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case float64:
		return fmt.Sprintf("%.0f", x)
	case int:
		return fmt.Sprintf("%d", x)
	case int64:
		return fmt.Sprintf("%d", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}

func Truncate(input string, limit int) string {
	input = strings.TrimSpace(input)
	if limit <= 0 || len([]rune(input)) <= limit {
		return input
	}
	runes := []rune(input)
	return string(runes[:limit-1]) + "..."
}
