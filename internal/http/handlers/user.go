package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"steamshark-api/internal/models"
	"steamshark-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const steamOpenIDURL = "https://steamcommunity.com/openid/login"

var steamIDRegex = regexp.MustCompile(`https://steamcommunity\.com/openid/id/(\d+)`)

type UserHandler struct {
	logger      *zap.Logger
	db          *gorm.DB
	steamAPIKey string
	jwtSecret   string
	jwtExpiry   time.Duration
	baseURL     string
}

func NewUserHandler(logger *zap.Logger, db *gorm.DB, steamAPIKey, jwtSecret, jwtExpiry, baseURL string) *UserHandler {
	expiry, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		expiry = 24 * time.Hour
	}
	return &UserHandler{
		logger:      logger,
		db:          db,
		steamAPIKey: steamAPIKey,
		jwtSecret:   jwtSecret,
		jwtExpiry:   expiry,
		baseURL:     baseURL,
	}
}

/*
Returns the Steam OpenID redirect URL for the frontend.

GET /auth/steam
*/
func (h *UserHandler) SteamLogin(ctx *gin.Context) {
	params := url.Values{
		"openid.ns":         {"http://specs.openid.net/auth/2.0"},
		"openid.mode":       {"checkid_setup"},
		"openid.return_to":  {h.baseURL + "/api/v1/auth/steam/callback"},
		"openid.realm":      {h.baseURL},
		"openid.identity":   {"http://specs.openid.net/auth/2.0/identifier_select"},
		"openid.claimed_id": {"http://specs.openid.net/auth/2.0/identifier_select"},
	}
	utils.Success(ctx, "Steam login URL", gin.H{
		"url": steamOpenIDURL + "?" + params.Encode(),
	})
}

/*
Receives and verifies the Steam OpenID assertion, returns JWT.

GET /auth/steam/callback
*/
func (h *UserHandler) SteamCallback(ctx *gin.Context) {
	if h.jwtSecret == "" {
		utils.Error(ctx, http.StatusServiceUnavailable, "JWT not configured on server")
		return
	}

	params := ctx.Request.URL.Query()

	// Verify assertion with Steam
	steamID, err := h.verifySteamAssertion(params)
	if err != nil {
		h.logger.Error("steam assertion failed: " + err.Error())
		utils.Error(ctx, http.StatusUnauthorized, "Steam authentication failed")
		return
	}

	// Fetch profile from Steam Web API
	profile, err := h.fetchSteamProfile(steamID)
	if err != nil {
		h.logger.Error("failed to fetch steam profile: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "failed to fetch Steam profile")
		return
	}

	// Upsert user
	var user models.User
	result := h.db.WithContext(ctx).Where("steam_id = ?", steamID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = models.User{
			SteamID:     steamID,
			DisplayName: profile.DisplayName,
			AvatarURL:   profile.AvatarURL,
			Role:        "user",
		}
		if err := h.db.WithContext(ctx).Create(&user).Error; err != nil {
			h.logger.Error("failed to create user: " + err.Error())
			utils.Error(ctx, http.StatusInternalServerError, "failed to create user")
			return
		}
	} else if result.Error != nil {
		h.logger.Error("db error on user lookup: " + result.Error.Error())
		utils.Error(ctx, http.StatusInternalServerError, "failed to authenticate")
		return
	} else {
		// Update display name and avatar in case they changed on Steam
		h.db.WithContext(ctx).Model(&user).Updates(map[string]any{
			"display_name": profile.DisplayName,
			"avatar_url":   profile.AvatarURL,
		})
	}

	// Issue JWT
	token, err := h.issueJWT(user)
	if err != nil {
		h.logger.Error("failed to sign JWT: " + err.Error())
		utils.Error(ctx, http.StatusInternalServerError, "failed to issue token")
		return
	}

	utils.Success(ctx, "Authenticated", gin.H{
		"token": token,
		"user": gin.H{
			"id":           user.ID,
			"steam_id":     user.SteamID,
			"display_name": user.DisplayName,
			"avatar_url":   user.AvatarURL,
			"role":         user.Role,
		},
	})
}

/*
Returns the current authenticated user's profile.

GET /auth/me
*/
func (h *UserHandler) Me(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var user models.User
	if err := h.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "user not found")
		return
	}

	utils.Success(ctx, "User profile", gin.H{
		"id":           user.ID,
		"steam_id":     user.SteamID,
		"display_name": user.DisplayName,
		"avatar_url":   user.AvatarURL,
		"role":         user.Role,
		"created_at":   user.CreatedAt.Format(time.RFC3339),
	})
}

/*
Logout — on stateless JWT the client discards the token.
This endpoint exists as a clean contract for the frontend.

POST /auth/logout
*/
func (h *UserHandler) Logout(ctx *gin.Context) {
	utils.Success(ctx, "Logged out", nil)
}

// --- helpers ---

func (h *UserHandler) verifySteamAssertion(params url.Values) (string, error) {
	// Build check_authentication request
	checkParams := url.Values{}
	for k, v := range params {
		checkParams[k] = v
	}
	checkParams.Set("openid.mode", "check_authentication")

	resp, err := http.PostForm(steamOpenIDURL, checkParams)
	if err != nil {
		return "", fmt.Errorf("verification request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "is_valid:true") {
		return "", fmt.Errorf("steam rejected the assertion")
	}

	claimedID := params.Get("openid.claimed_id")
	matches := steamIDRegex.FindStringSubmatch(claimedID)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract steam id from claimed_id")
	}

	return matches[1], nil
}

type steamProfile struct {
	DisplayName string
	AvatarURL   string
}

func (h *UserHandler) fetchSteamProfile(steamID string) (*steamProfile, error) {
	if h.steamAPIKey == "" {
		return &steamProfile{DisplayName: steamID, AvatarURL: ""}, nil
	}

	apiURL := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s",
		h.steamAPIKey, steamID,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			Players []struct {
				PersonaName string `json:"personaname"`
				Avatar      string `json:"avatarfull"`
			} `json:"players"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Response.Players) == 0 {
		return &steamProfile{DisplayName: steamID, AvatarURL: ""}, nil
	}

	p := result.Response.Players[0]
	return &steamProfile{DisplayName: p.PersonaName, AvatarURL: p.Avatar}, nil
}

func (h *UserHandler) issueJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"steam_id": user.SteamID,
		"role":     user.Role,
		"exp":      time.Now().Add(h.jwtExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
