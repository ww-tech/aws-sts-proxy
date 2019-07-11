package helpers

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// IDData is the data we care about from the IDTokens
type IDData struct {
	Email string `json:"email"`
}

// TokenData contains all data from the Tokens we care about
type TokenData struct {
	IDToken string `json:"id_token"`
}

// AccessTokenPost is data for the Post to get AccessTokens
type AccessTokenPost struct {
	JSON         bool   `json:"json"`
	RefreshToken string `json:"refresh_token"`
}

// STSSession the object to represent an sts session
type STSSession struct {
	ExpiresAt time.Time         `json:"expires_at"`
	Creds     credentials.Value `json:"creds"`
}
