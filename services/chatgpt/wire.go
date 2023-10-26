package chatgpt

import "github.com/google/wire"

var ChatGPTServiceProviderSet = wire.NewSet(New)
