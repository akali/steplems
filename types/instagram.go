package types

import (
	"fmt"
	"strings"
)

type GoInstaConfigPath string

type InstaCachePath string

type InstagramMessage struct {
	Caption  string
	Username string
	Link     string
	Path     string
}

func (m *InstagramMessage) SanitizeTitle() {
	replaceChars := []string{"*", "_", "~", "`", "|", "[", "]"}
	for _, old := range replaceChars {
		m.Caption = strings.Replace(m.Caption, old, fmt.Sprintf("\\%s", old), -1)
	}
	m.Caption = strings.TrimSpace(m.Caption)
}

func (m *InstagramMessage) FormCaption() string {
	m.SanitizeTitle()

	return fmt.Sprintf("*Shared by [%s](https://instagram.com/%s)*\n\n%s\n\n[link | сілтеме](%s)", m.Username, m.Username, m.Caption, m.Link)
}
