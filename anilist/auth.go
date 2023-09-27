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

func getAuthenticationURL(clientId string) string {
	return fmt.Sprintf("https://anilist.co/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", clientId, "http://localhost:3000")
}

func getAccessTokenRes(clientId, clientSecret, authorizationCode string) (*models.TokenRes, error) {
	data := accessTokenReqData{
		GrantType:    "authorization_code",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURI:  "http://localhost:3000",
		Code:         authorizationCode,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, &models.AppError{
			Message: "Error marshalling anilist access token request body",
			Err:     err,
		}
	}

	req, err := http.NewRequest("POST", "https://anilist.co/api/v2/oauth/token", bytes.NewReader(reqBody))
	if err != nil {
		return nil, &models.AppError{
			Message: "Error creating anilist access token request",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, &models.AppError{
			Message: "Error making anilist access token request",
			Err:     err,
		}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to read the access token response body",
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
			Message: "Failed to parse the access token response",
		}
	}

	return &tokenRes, nil
}
