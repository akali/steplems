package types

import (
	"fmt"
	"strings"
)

type YoutubeMessage struct {
	Title, Link, Path string
}

func (m *YoutubeMessage) SanitizeTitle() {
	replaceChars := []string{"*", "_", "~", "`", "|", "[", "]"}
	for _, old := range replaceChars {
		m.Title = strings.Replace(m.Title, old, fmt.Sprintf("\\%s", old), -1)
	}
}

func (m *YoutubeMessage) FormCaption() string {
	m.SanitizeTitle()

	return fmt.Sprintf("*%s*\n\n[link | сілтеме](%s)", m.Title, m.Link)
}
