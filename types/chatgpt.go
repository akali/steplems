package types

import "fmt"

type ModelStorage struct {
	Model   string
	Backend string
}

func (m ModelStorage) String() string {
	return fmt.Sprintf("model: %s, backend: %s", m.Model, m.Backend)
}
