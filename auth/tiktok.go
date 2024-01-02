package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	TIK_TOK_AUTH_ENDPOINT     = "https://www.tiktok.com/v2/auth/authorize/"
	TIK_TOK_AUTH_RESPONSETYPE = "code"
)

type TikTokAuthRequest struct {
	clientKey   string
	redirectUri string
	scope       string
	//'state' should be used to hold a unique identifier to validate incoming redirects.
	state string
}

func TikTokAuthHandler(redirectURI string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("tiktok redirect URL set to: %s", redirectURI)
		//Generate session UUID
		sessionID := generateUuid()

		//Get authorization page URL for our App Data
		authReqData := &TikTokAuthRequest{
			clientKey:   os.Getenv("TT_CK"),
			redirectUri: redirectURI,
			scope:       "user.info.basic,video.publish,video.upload",
			state:       fmt.Sprint(sessionID),
		}

		tikTokAuthURL, err := authReqData.createAuthPageUrl()
		if err != nil {
			log.Err(err).Msg("error generating tiktok auth url")
			w.WriteHeader(http.StatusUnprocessableEntity)
		}

		//Verify our header is set to correctly
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		//Redirect to the TikTokAuth url
		http.Redirect(w, r, tikTokAuthURL.String(), http.StatusTemporaryRedirect)
	}
}

func (tta *TikTokAuthRequest) createAuthPageUrl() (*url.URL, error) {
	urlVals := url.Values{}
	urlVals.Set("client_key", tta.clientKey)
	urlVals.Set("redirect_uri", tta.redirectUri)
	urlVals.Set("response_type", TIK_TOK_AUTH_RESPONSETYPE)
	urlVals.Set("scope", tta.scope)
	urlVals.Set("state", tta.state)

	return url.Parse(fmt.Sprintf("%s%s", TIK_TOK_AUTH_ENDPOINT, urlVals.Encode()))
}

func generateUuid() int {
	return int(uuid.Must(uuid.NewRandom()).ID())
}
