package instagram

import "github.com/google/wire"

var InstagramServiceProviderSet = wire.NewSet(New)
