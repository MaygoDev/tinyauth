package providers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"tinyauth/internal/constants"

	"github.com/rs/zerolog/log"
)

// Response for the google user endpoint
type DiscordUserInfoResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// The scopes required for the google provider
func DiscordScopes() []string {
	return []string{"email", "identify"}
}

func GetDiscordUser(client *http.Client) (constants.Claims, error) {
	var user constants.Claims

	res, err := client.Get("https://discord.com/api/v10/users/@me")
	if err != nil {
		return user, err
	}
	defer res.Body.Close()

	log.Debug().Msg("Got response from discord")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return user, err
	}

	log.Debug().Msg("Read body from discord")

	var userInfo DiscordUserInfoResponse

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return user, err
	}

	log.Debug().Msg("Parsed user from discord")

	user.PreferredUsername = strings.Split(userInfo.Email, "@")[0]
	user.Name = userInfo.Username
	user.Email = userInfo.Email

	return user, nil
}
