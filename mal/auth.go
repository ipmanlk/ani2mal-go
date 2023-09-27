package mal

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"ipmanlk/ani2mal/config"
	"ipmanlk/ani2mal/models"
	"ipmanlk/ani2mal/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func performAuth() {
	fmt.Print("Enter Client ID: ")
	clientId := utils.GetStrInput()

	fmt.Print("Enter Client Secret: ")
	clientSecret := utils.GetStrInput()

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		log.Fatal(err.Error())
	}

	loginURL := getAuthenticationURL(clientId, codeVerifier)
	fmt.Printf("Login URL: %s\n", loginURL)

	fmt.Print("Enter the code from the login URL: ")
	code := utils.GetStrInput()

	res, err := getAccessTokenRes(clientId, clientSecret, code, codeVerifier)
	if err != nil {
		log.Fatal(err.Error())
	}

	appConfig := config.GetAppConfig()

	appConfig.SaveMalConfig(&models.MalConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		TokenRes:     *res,
	})

	fmt.Println("Authentication successful. Access token has been saved.")
}

func GetMalAcessCode() (string, error) {
	malConfig := config.GetAppConfig().GetMalConfig()

	// check if token is expired or will expire soon
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Duration(malConfig.TokenRes.ExpiresIn) * time.Second)
	expirationBuffer := 20 * time.Minute

	if expirationTime.After(currentTime.Add(expirationBuffer)) {
		// token is not expired
		return malConfig.TokenRes.AccessToken, nil
	}

	// token is expired and new one should be requested
	res, err := getRefreshTokenRes(malConfig.ClientId, malConfig.ClientSecret, malConfig.TokenRes.RefreshToken)
	if err != nil {
		return "", err
	}

	// Save new token info in the Mal config
	malConfig.TokenRes = *res
	config.GetAppConfig().SaveMalConfig(malConfig)

	return res.AccessToken, nil
}

// retrieves the authentication URL with code_challenge
func getAuthenticationURL(clientId, codeChallenge string) string {
	return fmt.Sprintf("https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s", clientId, codeChallenge)
}

// exchanges the auth code for an access token
func getAccessTokenRes(clientId, clientSecret, authorizationCode, codeVerifier string) (*models.TokenRes, error) {
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("code", authorizationCode)
	data.Set("code_verifier", codeVerifier)
	data.Set("grant_type", "authorization_code")

	return sendMalTokenRequest(data)
}

// request a new access token using refresh token
func getRefreshTokenRes(clientId, clientSecret, refreshToken string) (*models.TokenRes, error) {
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	return sendMalTokenRequest(data)
}

func sendMalTokenRequest(data url.Values) (*models.TokenRes, error) {
	tokenEndpoint := "https://myanimelist.net/v1/oauth2/token"

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	res, err := client.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to request the access token",
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

	// Check for errors in the response
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
			Err:     err,
		}
	}

	return &tokenRes, nil
}

// generates a code verifier for OAuth2 PKCE
func generateCodeVerifier() (string, error) {
	// Generate a 32-byte (256-bit) random value
	verifierBytes := make([]byte, 32)
	_, err := rand.Read(verifierBytes)
	if err != nil {
		return "", &models.AppError{
			Message: "Error generating random string",
			Err:     err,
		}
	}
	// Encode the random bytes as a URL-safe base64 string
	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
	return codeVerifier, nil
}
