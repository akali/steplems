package providers

import (
	"github.com/google/wire"
	"github.com/kkdai/youtube/v2"
)

func ProvideYoutubeClient() youtube.Client {
	return youtube.Client{}
}

var YoutubeClientProvider = wire.NewSet(ProvideYoutubeClient)
