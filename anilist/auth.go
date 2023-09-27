package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"ipmanlk/ani2mal/config"
	"ipmanlk/ani2mal/models"
	"ipmanlk/ani2mal/utils"
	"log"
	"net/http"
	"time"
)

type accessTokenReqData struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Code         string `json:"code"`
}

type refreshTokenReqData struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func PerformAuth() {
	fmt.Print("Enter Anilist Username: ")
	username := utils.GetStrInput()

	fmt.Print("Enter Client ID: ")
	clientId := utils.GetStrInput()

	fmt.Print("Enter Client Secret: ")
	clientSecret := utils.GetStrInput()

	loginURL := getAuthenticationURL(clientId)
	fmt.Printf("Login URL: %s\n", loginURL)

	fmt.Print("Enter the code from the login URL: ")
	code := utils.GetStrInput()

	res, err := getAccessTokenRes(clientId, clientSecret, code)

	if err != nil {
		log.Fatal(err.Error())
	}

	appConfig := config.GetAppConfig()

	appConfig.SaveAnilistConfig(&models.AnilistConfig{
		Username:     username,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		TokenRes:     *res,
	})

	fmt.Println("Authentication successful. Access token has been saved.")
}

func GetAccessCode() (string, error) {
	anilistConfig := config.GetAppConfig().GetAnilistConfig()

	// check if token is expired or will expire soon
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Duration(anilistConfig.TokenRes.ExpiresIn) * time.Second)
	expirationBuffer := 20 * time.Minute

	if expirationTime.After(currentTime.Add(expirationBuffer)) {
		// token is not expired
		return anilistConfig.TokenRes.AccessToken, nil
	}

	// token is expired and new one should be requested
	res, err := getRefreshTokenRes(anilistConfig.ClientId, anilistConfig.ClientSecret, anilistConfig.TokenRes.RefreshToken)
	if err != nil {
		return "", err
	}

	anilistConfig.TokenRes = *res
	config.GetAppConfig().SaveAnilistConfig(anilistConfig)

	return res.AccessToken, nil
}

func getAuthenticationURL(clientId string) string {
	return fmt.Sprintf("https://anilist.co/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", clientId, "http://localhost:3000")
}

// exchanges the auth code for an access token
func getAccessTokenRes(clientId, clientSecret, authorizationCode string) (*models.TokenRes, error) {
	data := accessTokenReqData{
		GrantType:    "authorization_code",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURI:  "http://localhost:3000",
		Code:         authorizationCode,
	}

	return sendTokenRequest(data)
}

// request a new access token using refresh token
func getRefreshTokenRes(clientId, clientSecret, refreshToken string) (*models.TokenRes, error) {
	data := refreshTokenReqData{
		GrantType:    "refresh_token",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}
	return sendTokenRequest(data)
}

func sendTokenRequest(data any) (*models.TokenRes, error) {
	tokenEndpoint := "https://anilist.co/api/v2/oauth/token"

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, &models.AppError{
			Message: "Error marshalling anilist token request body",
			Err:     err,
		}
	}

	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, &models.AppError{
			Message: "Error creating anilist token request",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, &models.AppError{
			Message: "Error making anilist token request",
			Err:     err,
		}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to read the token response body",
			Err:     err,
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, &models.AppError{
			Message: "Access token request failed " + fmt.Sprintf("Error: %s", body),
		}
	}

	tokenRes := models.TokenRes{}
	err = json.Unmarshal(body, &tokenRes)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to parse the token response",
		}
	}

	return &tokenRes, nil
}
