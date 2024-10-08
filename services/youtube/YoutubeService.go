package youtube

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"steplems-bot/types"

	"github.com/avast/retry-go/v4"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"github.com/hashicorp/go-multierror"
	"github.com/kkdai/youtube/v2"
)

type YoutubeService struct {
	client  youtube.Client
	pattern *regexp.Regexp
	logger  zerolog.Logger
}

func NewYoutubeService(client youtube.Client, logger zerolog.Logger) *YoutubeService {
	return &YoutubeService{
		client:  client,
		pattern: regexp.MustCompile(ytLinkRegex),
		logger:  logger,
	}
}

var YoutubeServiceProvider = wire.NewSet(NewYoutubeService)

type VideoType string
type QualityType string

const (
	RETRY_TIMES = 5
)

const (
	HD     QualityType = "hd"
	HD720  QualityType = "hd720"
	HD1080 QualityType = "hd1080"
	MEDIUM QualityType = "medium"
)

const (
	MP4  VideoType = "video/mp4"
	MKV  VideoType = "video/mkv"
	WEBM VideoType = "video/webm"
)

var (
	ytLinkRegex    = "(((?:https?:)?\\/\\/)?((?:www|m)\\.)?((?:youtube\\.com))(\\/(shorts\\/))([\\w\\-]+)(\\S+)?)"
	allowedQuality = []QualityType{HD, HD720, HD1080, MEDIUM}
	allowedType    = []VideoType{MP4, MKV, WEBM}
)

func (ys *YoutubeService) Match(text string) []string {
	return ys.pattern.FindAllString(text, -1)
}

func (ys *YoutubeService) IsYoutubeMessage(update tbot.Update) bool {
	if update.Message == nil || update.Message.From.IsBot {
		return false
	}
	links := ys.Match(update.Message.Text)
	return len(links) != 0
}

func (ys *YoutubeService) chooseFormat(formats youtube.FormatList) *youtube.Format {
	formats = formats.WithAudioChannels()

	for i := range formats {
		for _, t := range allowedType {
			if !strings.Contains(formats[i].MimeType, string(t)) {
				continue
			}
			for _, q := range allowedQuality {
				if formats[i].Quality == string(q) || formats[i].QualityLabel == string(q) {
					return &formats[i]
				}
			}
		}
	}

	return &formats[0]
}

func (ys *YoutubeService) downloadPerLink(
	v *youtube.Video,
	format *youtube.Format,
	folder string,
) (string, error) {
	stream, _, err := ys.client.GetStream(v, format)
	if err != nil {
		return "", err
	}

	a := strings.Split(format.MimeType, "/")
	fileExt := strings.Split(a[1], ";")[0]

	filename := path.Join(folder, v.ID) + "." + fileExt

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, stream)
	if err != nil {
		return "", err
	}

	ys.logger.Info().Str("id", v.ID).Str("path", filename).Msg("downloaded short")

	return filename, nil
}

func (ys *YoutubeService) Download(links []string, folder string) ([]types.YoutubeMessage, error) {
	msgs := make([]types.YoutubeMessage, 0)

	for _, l := range links {
		v, err := ys.client.GetVideo(l)
		if err != nil {
			ys.logger.Err(err).Msg("can't get metadata for link")

			return nil, err
		}

		chosenFormat := ys.chooseFormat(v.Formats)
		filename, err := retry.DoWithData(func() (string, error) {
			return ys.downloadPerLink(v, chosenFormat, folder)
		}, retry.Attempts(RETRY_TIMES), retry.DelayType(retry.BackOffDelay))

		if err != nil {
			return nil, err
		}

		msgs = append(msgs, types.YoutubeMessage{
			Link:  l,
			Title: v.Title,
			Path:  filename})
	}

	return msgs, nil
}

func (ys *YoutubeService) MessageUpdate(message *tbot.Message) (tbot.VideoConfig, error) {
	links := ys.Match(message.Text)

	if len(links) == 0 || message.From.IsBot {
		return tbot.VideoConfig{}, nil
	}

	ys.logger.Info().Int("length", len(links)).Stringer("user", message.From).Msg("detected youtube short links of {{length}} length from {{user}}")

	folder, err := os.MkdirTemp("/tmp", "yt*")
	if err != nil {
		return tbot.VideoConfig{}, err
	}

	yms, err := ys.Download(links, folder)
	if err != nil {
		ys.logger.Err(err).Msg("failed to download")
		return tbot.VideoConfig{}, fmt.Errorf("failed to process video: %s", err.Error())
	}
	var filesErrs *multierror.Error
	for _, ym := range yms {
		v := tbot.NewVideo(message.Chat.ID, tbot.FilePath(ym.Path))
		v.Caption = ym.FormCaption()
		v.ParseMode = tbot.ModeMarkdown
		v.ReplyToMessageID = message.MessageID

		return v, nil
	}
	return tbot.VideoConfig{}, filesErrs.ErrorOrNil()
}
