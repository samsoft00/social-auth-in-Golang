package tiktok

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ITiktok interface {
	Init(g *gin.Context)
	Callback(g *gin.Context)
}

const (
	NgrokIP = "78d7-102-89-34-115" // <- IP address from ngrok
	// -- Extract from Tiktok developer portal - https://developers.tiktok.com
	AuthURL      = "https://open-api.tiktok.com/platform/oauth/connect/"
	TokenURL     = "https://open-api.tiktok.com/oauth/access_token/"
	ClientKey    = "awda8zjbmkd8jxxx"
	ClientSecret = "09ab93176788fe59d9497b19cdae5xxx"
	GrantType    = "authorization_code"
	Scope        = "user.info.basic"
	ResponseCode = "code"
)

type config struct {
	BaseURL                 *url.URL
	AuthURL                 *url.URL
	RedirectURL             *url.URL
	TokenURL                *url.URL
	ClientKey, ClientSecret string
	Scopes, GrantType       string
	State, ResponseType     string
}

type Controller struct {
	logger *zap.Logger
	config config
	client *http.Client
}

var _ ITiktok = &Controller{}

func NewController(withLogger *zap.Logger) *Controller {
	BaseURL := fmt.Sprintf("https://%s.ngrok.io", NgrokIP)
	RedirectURL := fmt.Sprintf("%s/tiktok/callback", BaseURL)

	return &Controller{logger: withLogger, config: config{
		BaseURL:      parseURL(BaseURL),
		RedirectURL:  parseURL(RedirectURL),
		AuthURL:      parseURL(AuthURL),
		TokenURL:     parseURL(TokenURL),
		ClientKey:    ClientKey,
		ClientSecret: ClientSecret,
		Scopes:       Scope,
		GrantType:    GrantType,
		State:        "vdcm9faw5u",
		ResponseType: ResponseCode,
	},
		client: &http.Client{Timeout: time.Second * 10},
	}
}

func (c Controller) Init(g *gin.Context) {
	redirectURL := c.buildURL()

	c.logger.Info("redirecting to " + redirectURL.String())
	g.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}

func (c Controller) buildURL() *url.URL {
	redirectURL := *c.config.AuthURL

	query := redirectURL.Query()
	query.Set("client_key", c.config.ClientKey)
	query.Set("redirect_uri", c.config.RedirectURL.String())
	query.Set("response_type", c.config.ResponseType)
	query.Set("scope", c.config.Scopes)
	query.Set("state", c.config.State)
	redirectURL.RawQuery = query.Encode()

	return &redirectURL
}

func (c Controller) Callback(g *gin.Context) {
	code := g.Query("code")

	if code == "" {
		ginResponse(g, http.StatusBadRequest, "query parameter code is missing")
		return
	}

	tokenURL := *c.config.TokenURL

	query := tokenURL.Query()
	query.Set("client_key", c.config.ClientKey)
	query.Set("client_secret", c.config.ClientSecret)
	query.Set("grant_type", c.config.GrantType)
	query.Set("code", code)
	tokenURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(g.Request.Context(), "POST", tokenURL.String(), http.NoBody)
	if err != nil {
		ginResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := c.client.Do(req)
	if err != nil {
		ginResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ginResponse(g, http.StatusInternalServerError, err.Error())
		return
	}
	var result RefreshTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		c.logger.Debug(string(body))
		ginResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"open_id":       result.Data.OpenID,
		"access_token":  result.Data.AccessToken,
		"refresh_token": result.Data.RefreshToken,
	})
}

func parseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func ginResponse(g *gin.Context, code int, message string) {
	g.JSON(code, gin.H{"status": false, "message": message})
}
