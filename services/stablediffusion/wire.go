package stablediffusion

import "github.com/google/wire"

var StableDiffusionSet = wire.NewSet(NewStableDiffusionService)