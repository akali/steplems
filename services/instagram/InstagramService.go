package instagram

import (
	"fmt"
	"github.com/Davincible/goinsta/v3"
	"github.com/avast/retry-go/v4"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/h2non/filetype"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"steplems-bot/types"
	"strings"
	"time"

	"github.com/gregjones/httpcache/diskcache"
	"github.com/loganrk/go-heap-cache"

	_ "gopkg.in/yaml.v2"
)

const (
	RetryTimes     = 5
	UpdateDuration = 1 * time.Minute
	UploadLimit    = 10
)

const (
	InstagramLinkRegex = `(((?:https?:)?\/\/)?((?:www|m)\.)?((?:instagram\.com))(\/((reel|post)\/))([\w\-]+)(\S+)?)`
)

func LinksMatching(text string) []string {
	pattern := regexp.MustCompile(InstagramLinkRegex)
	return pattern.FindAllString(text, -1)
}

type InstagramService struct {
	client        *goinsta.Instagram
	configPath    types.GoInstaConfigPath
	lastCheckTime time.Time
	profileCache  cache.Cache
	seenCache     *diskcache.Cache
	chatID        int64
}

func New(client *goinsta.Instagram, configPath types.GoInstaConfigPath, cachePath types.InstaCachePath) *InstagramService {
	return &InstagramService{
		client:        client,
		configPath:    configPath,
		seenCache:     diskcache.New(string(cachePath)),
		lastCheckTime: time.Now(),
	}
}

func (is *InstagramService) MessageUpdate(message *tbot.Message) (tbot.VideoConfig, error) {
	links := LinksMatching(message.Text)
	if len(links) == 0 || message.From.IsBot {
		return tbot.VideoConfig{}, nil
	}
	return tbot.VideoConfig{}, fmt.Errorf("unimplemented")
}

func (is *InstagramService) saveConfig() error {
	if err := is.client.Export(string(is.configPath)); err != nil {
		return err
	}
	log.Printf("Exported config to %q\n", is.configPath)
	return nil
}

func (is *InstagramService) getUserByIDCached(id int64) (*goinsta.User, error) {
	if is.profileCache == nil {
		is.profileCache = cache.New(&cache.Config{
			Capacity:       cache.DEFAULT_CAPACITY,
			Expire:         int64(time.Hour / 1e9),
			EvictionPolicy: cache.EVICTION_POLICY_LRU,
		})
	}
	key := fmt.Sprintf("%d", id)
	value, err := is.profileCache.Get(key)
	if err != nil {
		value, err = is.client.Profiles.ByID(id)
		if err != nil {
			return nil, err
		}
		is.profileCache.Set(key, value)
	}
	return value.(*goinsta.User), nil
}

func (is *InstagramService) SetUpdateChatID(id int64) {
	is.chatID = id
	log.Printf("Updated ig update chat id to %d\n", id)
}

func (is *InstagramService) seen(item *goinsta.InboxItem) bool {
	_, ok := is.seenCache.Get(item.ID)
	if !ok {
		return false
	}
	return true
}

func (is *InstagramService) markSeen(item *goinsta.InboxItem) {
	is.seenCache.Set(item.ID, []byte{1})
}

func GetDownloadable(media goinsta.Item) (any, error) {
	var downloadable any = nil

	if len(media.Images.Versions) > 0 {
		downloadable = media.Images.Versions
	}

	if len(media.Videos) > 0 {
		downloadable = media.Videos
	}

	if downloadable == nil {
		return nil, fmt.Errorf("failed to find any downloadable for %q from %q", media.Title, media.User.Username)
	}

	return downloadable, nil
}

func (is *InstagramService) extractMessages(conv *goinsta.Conversation) ([]goinsta.Item, error) {
	items := make(map[string]goinsta.InboxItem)

	startTime := time.Now()

	for conv.Next() {
		itemsAdded := false
		for _, item := range conv.Items {
			ts := time.UnixMicro(item.Timestamp)
			// lastCheckTime < itemTime
			if is.lastCheckTime.Before(ts) || (is.lastCheckTime.Add(-time.Hour).Before(ts) && !is.seen(item)) {
				if _, ok := items[item.ID]; !ok {
					is.markSeen(item)
					itemsAdded = true
					items[item.ID] = *item
					if len(items) >= UploadLimit {
						break
					}
				}
			}
		}
		if !itemsAdded || len(items) >= UploadLimit {
			break
		}
	}

	is.lastCheckTime = startTime

	log.Printf("Received %d new messages in chat %q!\n", len(items), conv.Title)

	var urls []goinsta.Item

	for _, item := range items {
		var media *goinsta.Item

		if item.Clip != nil {
			media = &item.Clip.Media
		}
		if item.Reel != nil {
			media = &item.Reel.Media
		}
		if item.Media != nil {
			media = item.Media
		}
		if item.MediaShare != nil {
			media = item.MediaShare
		}
		if item.VisualMedia != nil {
			media = item.VisualMedia.Media
		}

		if media == nil {
			continue
		}

		sender, err := is.getUserByIDCached(item.UserID)
		if err != nil {
			log.Printf("Failed to get user profile %d in item: %w", item.UserID, err)
		} else {
			media.User = *sender
		}

		urls = append(urls, *media)
	}

	return urls, nil
}

func (is *InstagramService) perItemDownload(item goinsta.Item, directory string) (string, error) {
	client := http.DefaultClient
	downloadable, err := GetDownloadable(item)
	if err != nil {
		return "", err
	}

	url := goinsta.GetBest(downloadable)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	kind, err := filetype.Match(b)
	if err != nil {
		return "", err
	}

	ext := kind.Extension

	filename := path.Join(directory, item.GetID()) + "." + ext
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(b); err != nil {
		return "", err
	}

	return filename, nil
}

func (is *InstagramService) downloadItems(items []goinsta.Item, directory string) ([]types.InstagramMessage, error) {
	log.Printf("Have %d items to download.", len(items))

	var result []types.InstagramMessage
	for _, item := range items {

		filename, err := retry.DoWithData(func() (string, error) {
			return is.perItemDownload(item, directory)
		}, retry.Attempts(RetryTimes), retry.DelayType(retry.BackOffDelay))

		if err != nil {
			log.Printf("Failed to download item %q from %q", item.GetID(), item.User.Username)
		}

		result = append(result, types.InstagramMessage{
			Username: item.User.Username,
			Caption:  item.Caption.Text,
			Link:     fmt.Sprintf("instagram.com/reel/%s", item.Code),
			Path:     filename,
		})
	}
	return result, nil
}

func (is *InstagramService) checkForUpdates(chatTitle string) ([]goinsta.Item, error) {
	defer is.saveConfig()

	if err := is.client.OpenApp(); err != nil {
		return nil, err
	}

	is.client.Inbox.Reset()

	if err := is.client.Inbox.Sync(); err != nil {
		return nil, err
	}

	for {
		for _, conv := range is.client.Inbox.Conversations {
			if !strings.Contains(conv.Title, chatTitle) {
				continue
			}

			if err := conv.GetItems(); err != nil {
				return nil, err
			}

			return is.extractMessages(conv)
		}
		if !is.client.Inbox.Next() {
			break
		}
	}

	return nil, fmt.Errorf("failed to find any chat with chatTitle %q", chatTitle)
}

func (is *InstagramService) runEach(chatTitle string, sender types.Sender, chatID int64) {
	log.Println("Starting runEach iteration.")

	updates, err := is.checkForUpdates(chatTitle)
	if err != nil {
		msg := fmt.Sprintf("Failed to fetch updates: %w", err)
		log.Println(msg)
		sender.Send(tbot.NewMessage(chatID, msg))
		return
	}
	directory, err := os.MkdirTemp("/tmp", "ig*")
	if err != nil {
		log.Printf("Failed to fetch updates: %w", updates)
		return
	}
	messages, err := is.downloadItems(updates, directory)
	if err != nil {
		msg := fmt.Sprintf("Failed to download items: %w", err)
		log.Println(msg)
		sender.Send(tbot.NewMessage(chatID, msg))
		return
	}
	for _, msg := range messages {
		v := tbot.NewVideo(chatID, tbot.FilePath(msg.Path))
		v.Caption = msg.FormCaption()
		v.ParseMode = tbot.ModeMarkdown

		if _, err := sender.Send(v); err != nil {
			msg := fmt.Sprintf("Failed to send message: %w", err)
			log.Println(msg)
		}
	}
}

func (is *InstagramService) Run(chatTitle string, sender types.Sender, chatID int64, restartChan chan<- struct{}) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Printf("Paniced on error: %v\n", err)
			} else {
				log.Printf("Paniced with: %v\n", r)
			}
			restartChan <- struct{}{}
		}
	}()

	is.chatID = chatID

	log.Println("InstagramService started")

	timer := time.NewTimer(UpdateDuration)

	for {
		select {
		case <-timer.C:
			is.runEach(chatTitle, sender, is.chatID)
			timer.Reset(UpdateDuration)
		}
	}
}
