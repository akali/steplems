package types

import "fmt"

type ModelStorage struct {
	Model      string `json:"model"`
	Backend    string `json:"backend"`
	ImGenModel string `json:"imgenmodel"`
}

func (m ModelStorage) String() string {
	return fmt.Sprintf("model: %s, backend: %s, imgenmodel: %s", m.Model, m.Backend, m.ImGenModel)
}
