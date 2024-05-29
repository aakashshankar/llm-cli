package highlight

import (
	"github.com/alecthomas/chroma/v2/quick"
	"regexp"
	"strings"
)

func highlightCodeBlock(language, code string) (string, error) {
	var sb strings.Builder
	err := quick.Highlight(&sb, code, language, "terminal", "monokai")
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func RegularHighlight(text string) string {
	re := regexp.MustCompile("(?s)```(\\w+)(.*?)```")

	highlightedText := re.ReplaceAllStringFunc(text, func(codeBlock string) string {
		matches := re.FindStringSubmatch(codeBlock)
		if len(matches) < 3 {
			return codeBlock // Return original if the format is incorrect
		}
		language := matches[1]
		code := matches[2]

		highlighted, err := highlightCodeBlock(language, code)
		if err != nil {
			return codeBlock // Return original if highlighting fails
		}
		return highlighted
	})

	return highlightedText
}

func PartialHighlighter() func(partial string) string {
	var buffer string

	return func(partial string) string {
		// Append the new partial response to the buffer
		buffer += partial

		// Highlight complete code blocks and update the buffer
		highlightedText := RegularHighlight(buffer)

		// Remove processed parts from the buffer
		re := regexp.MustCompile("(?s)```(\\w+)(.*?)```")
		buffer = re.ReplaceAllString(buffer, "")

		// Return the highlighted text
		return highlightedText
	}
}
