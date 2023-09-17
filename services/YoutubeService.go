package services

import (
	"fmt"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"github.com/hashicorp/go-multierror"
	"github.com/kkdai/youtube/v2"
	"github.com/olehan/kek"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type YoutubeService struct {
	client     youtube.Client
	log        *kek.Logger
	pattern    *regexp.Regexp
	updateChan YoutubeUpdatesChannel
}

func NewYoutubeService(client youtube.Client, logFactory *kek.Factory) *YoutubeService {
	return &YoutubeService{
		client:     client,
		log:        logFactory.NewLogger("YoutubeService"),
		pattern:    regexp.MustCompile(ytLinkRegex),
		updateChan: make(YoutubeUpdatesChannel),
	}
}

var YoutubeServiceProvider = wire.NewSet(NewYoutubeService)

type YoutubeUpdatesChannel chan TelegramBotAction

func (ys *YoutubeService) Updates() YoutubeUpdatesChannel {
	return ys.updateChan
}

type TelegramBotAction func(api *tbot.BotAPI)

type YoutubeMessage struct {
	Title, Link, Path string
}

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

func (ys *YoutubeService) YoutubeMessage(update tbot.Update) bool {
	if update.Message == nil || update.Message.From.IsBot {
		return false
	}
	links := ys.Match(update.Message.Text)
	return len(links) != 0
}

func (ys *YoutubeService) retryIfErr(f func() error) {
	for i := 1; i <= RETRY_TIMES; i++ {
		if err := f(); err != nil {
			time.Sleep(time.Second * time.Duration(2<<i))
		} else {
			return
		}
	}
}

func (ys *YoutubeService) chooseFormat(formats youtube.FormatList) *youtube.Format {
	formats = formats.WithAudioChannels()

	for i := range formats {
		for _, q := range allowedQuality {
			for _, t := range allowedType {
				if (formats[i].Quality == string(q) || formats[i].QualityLabel == string(q)) &&
					strings.Contains(formats[i].MimeType, string(t)) {

					return &formats[i]
				}
			}
		}
	}

	return &formats[0]
}

func (ys *YoutubeService) downloadPerLinkBackedOff(
	v *youtube.Video,
	format *youtube.Format,
	folder string,
) (s string, err error) {
	ys.retryIfErr(func() error {
		s, err = ys.downloadPerLink(v, format, folder)

		return err
	})

	return
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

	filename := folder + "/" + url.PathEscape(v.ID) + "." + fileExt

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, stream)
	if err != nil {
		return "", err
	}

	ys.log.Succ.PrintTKV(
		"downloaded short by id {{id}} and saved it into {{path}}",
		"id", v.ID, "path", filename)

	return filename, nil
}

func (m *YoutubeMessage) SanitizeTitle() {
	m.Title = strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						strings.Replace(
							strings.Replace(m.Title, "*", "\\*", -1),
							"_", "\\_", -1),
						"~", "\\~", -1),
					"`", "\\`", -1),
				"|", "\\|", -1),
			"[", "\\[", -1),
		"]", "\\]", -1)

}

func (m *YoutubeMessage) FormCaption() string {
	m.SanitizeTitle()

	b := strings.Builder{}

	b.WriteRune('*')
	b.WriteString(m.Title)
	b.WriteRune('*')

	b.WriteString("\n\n[link | сілтеме]")
	b.WriteRune('(')
	b.WriteString(m.Link)
	b.WriteRune(')')
	return b.String()
}

func (ys *YoutubeService) Download(links []string, folder string) ([]YoutubeMessage, error) {
	msgs := make([]YoutubeMessage, 0)

	for _, l := range links {
		v, err := ys.client.GetVideo(l)
		if err != nil {
			ys.log.Note.PrintTKV("can't get metadata for link: {{error}}", "error", err)

			return nil, err
		}

		chosenFormat := ys.chooseFormat(v.Formats)
		filename, err := ys.downloadPerLinkBackedOff(v, chosenFormat, folder)
		if err != nil {
			return nil, err
		}

		// there has to be ffmpeg stuff

		msgs = append(msgs, YoutubeMessage{
			Link:  l,
			Title: v.Title,
			Path:  filename})
	}

	return msgs, nil
}

func (ys *YoutubeService) MessageUpdate(message *tbot.Message) error {
	links := ys.Match(message.Text)

	if len(links) == 0 || message.From.IsBot {
		return nil
	}

	ys.log.Info.PrintTKV(
		"detected youtube short links of {{length}} length from {{user}}",
		"length", len(links), "user", message.From.String())

	folder, err := os.MkdirTemp("/tmp", "yt*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(folder)

	yms, err := ys.Download(links, folder)
	if err != nil {
		ys.log.Error.Println(err.Error())
		// Let's try to reply to message with error message
		v := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("failed to process video: %s", err.Error()))
		v.ReplyToMessageID = message.MessageID

		if err := ys.SendMessage(v); err != nil {
			ys.log.Error.Println("failed to reply to message: ", err.Error())
		}
		return err
	}
	var filesErrs *multierror.Error
	for _, ym := range yms {
		v := tbot.NewVideo(message.Chat.ID, tbot.FilePath(ym.Path))
		v.Caption = ym.FormCaption()
		v.ParseMode = tbot.ModeMarkdown
		v.ReplyToMessageID = message.MessageID

		if err := ys.SendMessage(v); err != nil {
			ys.log.Error.Println(err.Error())
			// Let's try to reply to message with error message
			v := tbot.NewMessage(message.Chat.ID, fmt.Sprintf("failed to process video: %s", err.Error()))
			v.ReplyToMessageID = message.MessageID

			if err := ys.SendMessage(v); err != nil {
				ys.log.Error.Println("failed to reply to message: ", err.Error())
				filesErrs = multierror.Append(filesErrs, err)
			}
		}
	}
	return filesErrs.ErrorOrNil()
}

func (ys *YoutubeService) SendMessage(cfg tbot.Chattable) error {
	errChan := make(chan error)
	ys.updateChan <- TelegramBotAction(func(api *tbot.BotAPI) {
		_, err := api.Send(cfg)
		errChan <- err
	})
	return <-errChan
}
