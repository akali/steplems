package providers

import (
	"fmt"
	"github.com/google/wire"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"steplems-bot/types"
)

func ProvideSpotifyClientID() (types.SpotifyClientID, error) {
	return ProvideEnvironmentVariable[types.SpotifyClientID]("SPOTIFY_CLIENT_ID")()
}

func ProvideSpotifyClientSecret() (types.SpotifyClientSecret, error) {
	return ProvideEnvironmentVariable[types.SpotifyClientSecret]("SPOTIFY_CLIENT_SECRET")()
}

func ProvideSpotifyAuth(clientID types.SpotifyClientID, secret types.SpotifyClientSecret, hostname types.Hostname, port types.Port) *spotifyauth.Authenticator {
	var redirectURI = fmt.Sprintf("http://%s:%s/callback", string(hostname), string(port))
	return spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithClientID(string(clientID)),
		spotifyauth.WithClientSecret(string(secret)),
		spotifyauth.WithScopes(
			// ScopeImageUpload seeks permission to upload images to Spotify on your behalf.
			spotifyauth.ScopeImageUpload,
			// ScopePlaylistReadPrivate seeks permission to read
			// a user's private playlists.
			spotifyauth.ScopePlaylistReadPrivate,
			//spotifyauth.ScopePlaylistModifyPublic,
			//spotifyauth.ScopePlaylistModifyPrivate,
			// ScopePlaylistReadCollaborative seeks permission to
			// access a user's collaborative playlists.
			spotifyauth.ScopePlaylistReadCollaborative,
			// ScopeUserFollowModify seeks write/delete access to
			// the list of artists and other users that a user follows.
			spotifyauth.ScopeUserFollowModify,
			// ScopeUserFollowRead seeks read access to the list of
			// artists and other users that a user follows.
			spotifyauth.ScopeUserFollowRead,
			//spotifyauth.ScopeUserLibraryModify,
			// ScopeUserLibraryRead seeks read access to a user's "Your Music" library.
			spotifyauth.ScopeUserLibraryRead,
			// ScopeUserReadPrivate seeks read access to a user's
			// subscription details (type of user account).
			spotifyauth.ScopeUserReadPrivate,
			// ScopeUserReadEmail seeks read access to a user's email address.
			spotifyauth.ScopeUserReadEmail,
			// ScopeUserReadCurrentlyPlaying seeks read access to a user's currently playing track
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			// ScopeUserReadPlaybackState seeks read access to the user's current playback state
			spotifyauth.ScopeUserReadPlaybackState,
			// ScopeUserModifyPlaybackState seeks write access to the user's current playback state
			spotifyauth.ScopeUserModifyPlaybackState,
			// ScopeUserReadRecentlyPlayed allows access to a user's recently-played songs
			spotifyauth.ScopeUserReadRecentlyPlayed,
			// ScopeUserTopRead seeks read access to a user's top tracks and artists
			spotifyauth.ScopeUserTopRead,
			// ScopeStreaming seeks permission to play music and control playback on your other devices.
			spotifyauth.ScopeStreaming,
		))
}

var SpotifyAuthProviders = wire.NewSet(ProvideSpotifyClientID, ProvideSpotifyClientSecret, ProvideSpotifyAuth)
