package types

import (
	"fmt"
	"strings"
)

type GoInstaConfigPath string

type InstagramMessage struct {
	Caption string
	Link    string
	Path    string
}

func (m *InstagramMessage) SanitizeTitle() {
	replaceChars := []string{"*", "_", "~", "`", "|", "[", "]"}
	for _, old := range replaceChars {
		m.Caption = strings.Replace(m.Caption, old, fmt.Sprintf("\\%s", old), -1)
	}
}

func (m *InstagramMessage) FormCaption() string {
	m.SanitizeTitle()

	return fmt.Sprintf("%s\n\n[link | сілтеме](%s)", m.Caption, m.Link)
}
