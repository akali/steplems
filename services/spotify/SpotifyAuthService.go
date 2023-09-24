package spotify

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/olehan/kek"
	sapi "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"steplems-bot/types"
	"sync"
)

type SpotifyAuthService struct {
	port           types.Port
	expectedStates map[string]chan *sapi.Client
	authenticator  *spotifyauth.Authenticator
	log            *kek.Logger
	eMu            sync.Mutex
}

func NewSpotifyAuthService(port types.Port, authenticator *spotifyauth.Authenticator, factory *kek.Factory) *SpotifyAuthService {
	return &SpotifyAuthService{
		port:           port,
		authenticator:  authenticator,
		log:            factory.NewLogger("SpotifyAuthService"),
		expectedStates: make(map[string]chan *sapi.Client),
	}
}

func (s *SpotifyAuthService) ExpectAuthorize() (string, chan *sapi.Client) {
	s.eMu.Lock()
	defer s.eMu.Unlock()
	ch := make(chan *sapi.Client, 1)
	id := uuid.NewString()
	s.expectedStates[id] = ch
	return id, ch
}

func (s *SpotifyAuthService) Serve() {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		for state, update := range s.expectedStates {
			tok, err := s.authenticator.Token(r.Context(), state, r)
			if err != nil {
				continue
			}
			if st := r.FormValue("state"); st != state {
				continue
			}

			// use the token to get an authenticated client
			update <- sapi.New(s.authenticator.Client(r.Context(), tok))
			delete(s.expectedStates, state)
			return
		}
		http.Error(w, "Couldn't get token", http.StatusForbidden)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
