package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const scopes = "user-top-read"
 
var (
	auth *spotifyauth.Authenticator
	user *spotify.PrivateUser
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func InitSpotifyAuth(redirectURL string) {
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserTopRead),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
	)
}

func getAuthUrl() string {
	return auth.AuthURL(state)
}

func authentication() {
	client := <-ch
	u, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	user = u
	fmt.Println("You are logged in as: ", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	client := spotify.New(auth.Client(r.Context(), token))
	ch <- client
	fmt.Fprintf(w, "Login complete, please close the window")
}
