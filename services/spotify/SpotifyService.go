package spotify

import (
	"context"
	"errors"
	"fmt"
	"github.com/olehan/kek"
	spotify2 "steplems-bot/persistence/spotify_persistence"
	telegram2 "steplems-bot/persistence/telegram_persistence"
	"steplems-bot/types"
	"strings"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	sapi "github.com/zmb3/spotify/v2"

	"github.com/zmb3/spotify/v2/auth"
)

type SpotifyService struct {
	spotifyUserRepo  *spotify2.UserRepository
	telegramUserRepo *telegram2.UserRepository
	authenticator    *spotifyauth.Authenticator
	authService      *SpotifyAuthService
	port             types.Port
	log              *kek.Logger
}

func NewSpotifyService(port types.Port, authService *SpotifyAuthService, userRepo *spotify2.UserRepository, telegramUserRepo *telegram2.UserRepository, authenticator *spotifyauth.Authenticator, lf *kek.Factory) *SpotifyService {
	return &SpotifyService{
		port:             port,
		authService:      authService,
		spotifyUserRepo:  userRepo,
		telegramUserRepo: telegramUserRepo,
		authenticator:    authenticator,
		log:              lf.NewLogger("SpotifyService")}
}

func (s *SpotifyService) getSpotifyClient(ctx context.Context, telegramUser telegram2.User) (*sapi.Client, error) {
	spotifyUser, err := s.telegramUserRepo.EnsureSpotifyUserExists(telegramUser.TelegramExternalID)
	if err != nil {
		return nil, err
	}
	return s.CreateClient(ctx, spotifyUser)
}

func (s *SpotifyService) getOrCreateSpotifyClient(ctx context.Context, sender types.Sender, update tbot.Update) (*sapi.Client, error) {
	externalTelegramUser := update.SentFrom()

	telegramUser, err := s.telegramUserRepo.GetOrCreate(externalTelegramUser.ID, telegram2.FromExternalTelegramUser(externalTelegramUser, update.FromChat()))

	if err != nil {
		return nil, err
	}

	client, err := s.getSpotifyClient(ctx, telegramUser)
	if err != nil {
		if errors.Is(err, telegram2.NoSpotifyUserFound) {
			// Spotify spotifyUser does not exist yet.
			// Let's create one
			return s.authorizeAndSaveNewUser(ctx, sender, update, telegramUser)
		}
		return nil, err
	}

	return client, nil
}

func (s *SpotifyService) AuthorizeUser(ctx context.Context, sender types.Sender, update tbot.Update) (spotify2.User, error) {
	client, err := s.getOrCreateSpotifyClient(ctx, sender, update)
	if err != nil {
		return spotify2.User{}, err
	}
	privateUser, err := client.CurrentUser(ctx)
	if err != nil {
		return spotify2.User{}, err
	}
	return spotify2.PrivateUserToUser(privateUser), nil
}

func (s *SpotifyService) FindAll() []spotify2.User {
	return s.spotifyUserRepo.FindAll()
}

func (s *SpotifyService) CreateClient(ctx context.Context, user spotify2.User) (*sapi.Client, error) {
	return sapi.New(s.authenticator.Client(ctx, user.OAuthToken())), nil
}

func (s *SpotifyService) authorizeNewUser(ctx context.Context, sender types.Sender, update tbot.Update) (spotify2.User, error) {
	state, ch := s.authService.ExpectAuthorize()

	url := s.authenticator.AuthURL(state)
	msg := tbot.NewMessage(update.FromChat().ID, fmt.Sprintf("Follow link below to authorize: %s", url))
	if _, err := sender.Send(msg); err != nil {
		return spotify2.User{}, err
	}
	client := <-ch

	privateUser, err := client.CurrentUser(ctx)
	if err != nil {
		return spotify2.User{}, err
	}
	user := spotify2.PrivateUserToUser(privateUser)
	token, err := client.Token()
	if err != nil {
		return spotify2.User{}, err
	}
	return user.SetOAuthToken(token), nil
}

type SpotifySongMessage struct {
	*sapi.FullTrack
}

func (s SpotifySongMessage) Thumb() string {
	for _, image := range s.Album.Images {
		if image.URL != "" {
			return image.URL
		}
	}

	for _, image := range s.SimpleTrack.Album.Images {
		if image.URL != "" {
			return image.URL
		}
	}

	return ""
}

func (s SpotifySongMessage) Performer() string {
	var result []string
	for _, artist := range s.Artists {
		result = append(result, artist.Name)
	}
	return strings.Join(result, ", ")
}

func (s SpotifySongMessage) SpotifyLink() string {
	for _, value := range s.ExternalURLs {
		return value
	}
	return string(s.URI)
}

func (s SpotifySongMessage) AudioMessage(chatID int64) tbot.AudioConfig {
	preview := tbot.NewAudio(chatID, tbot.FileURL(s.PreviewURL))
	preview.Thumb = tbot.FileID(s.Thumb())
	preview.Duration = int(s.TimeDuration().Seconds())
	preview.Performer = s.Performer()
	preview.Title = s.Name
	preview.Caption = fmt.Sprintf("*%s - %s*\n\n[link | сілтеме](%s)", s.Performer(), s.Name, s.SpotifyLink())
	preview.ParseMode = tbot.ModeMarkdown
	return preview
}

func (s *SpotifyService) NowPlaying(ctx context.Context, sender types.Sender, update tbot.Update) error {
	externalTelegramUser := update.SentFrom()

	telegramUser, err := s.telegramUserRepo.GetOrCreate(externalTelegramUser.ID, telegram2.FromExternalTelegramUser(externalTelegramUser, update.FromChat()))

	if err != nil {
		return err
	}
	client, err := s.getSpotifyClient(ctx, telegramUser)

	if err != nil {
		return err
	}

	currentlyPlaying, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}

	preview := SpotifySongMessage{currentlyPlaying.Item}
	msg := preview.AudioMessage(update.FromChat().ID)

	_, err = sender.Send(msg)
	return err
}

func (s *SpotifyService) authorizeAndSaveNewUser(ctx context.Context, sender types.Sender, update tbot.Update, telegramUser telegram2.User) (*sapi.Client, error) {
	spotifyUser, err := s.authorizeNewUser(ctx, sender, update)
	if err != nil {
		return nil, err
	}
	spotifyUser, err = s.spotifyUserRepo.Create(spotifyUser)
	if err != nil {
		return nil, err
	}
	if err := s.telegramUserRepo.SaveSpotifyUser(telegramUser, spotifyUser); err != nil {
		return nil, err
	}
	return s.CreateClient(ctx, spotifyUser)
}
