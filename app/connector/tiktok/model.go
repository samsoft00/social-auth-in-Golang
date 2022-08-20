package tiktok

type TokenResponse struct {
	OpenID          string `json:"open_id"`
	Scope           string `json:"scope,omitempty"`
	AccessToken     string `json:"access_token"`
	ExpiresIn       int64  `json:"expires_in,omitempty"`
	RefreshToken    string `json:"refresh_token"`
	RefreshExpireIn int64  `json:"refresh_expires_in,omitempty"`
}

type RefreshTokenResponse struct {
	Data    TokenResponse `json:"data"`
	Message string        `json:"message"`
}
