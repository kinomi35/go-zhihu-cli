package output

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var (
	tagPattern         = regexp.MustCompile(`<[^>]+>`)
	blockTagPattern    = regexp.MustCompile(`(?i)</?(p|div|br|li|ul|ol|blockquote|h[1-6]|pre|section|article)[^>]*>`)
	inlineSpacePattern = regexp.MustCompile(`[ \t\r\f\v]+`)
	blankLinesPattern  = regexp.MustCompile(`\n{3,}`)
)

func StripHTML(input string) string {
	input = tagPattern.ReplaceAllString(input, "")
	input = html.UnescapeString(input)
	return strings.TrimSpace(strings.Join(strings.Fields(input), " "))
}

func StripHTMLPreserveLines(input string) string {
	input = blockTagPattern.ReplaceAllString(input, "\n")
	input = tagPattern.ReplaceAllString(input, "")
	input = html.UnescapeString(input)
	input = strings.ReplaceAll(input, "\u00a0", " ")

	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(inlineSpacePattern.ReplaceAllString(line, " "))
	}

	input = strings.Join(lines, "\n")
	input = blankLinesPattern.ReplaceAllString(input, "\n\n")
	return strings.TrimSpace(input)
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
