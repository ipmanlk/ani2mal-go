package auth

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

func AuthMal() {
	fmt.Print("Enter Client ID: ")
	clientId := utils.GetStrInput()

	fmt.Print("Enter Client Secret: ")
	clientSecret := utils.GetStrInput()

	// Generate a code verifier and code challenge
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	// Use the getAuthenticationURL function to retrieve the login URL
	loginURL := getAuthenticationURL(clientId, codeVerifier)
	fmt.Printf("Login URL: %s\n", loginURL)

	fmt.Print("Enter the code from the login URL: ")
	code := utils.GetStrInput()

	token, err := getAccessToken(clientId, clientSecret, code, codeVerifier)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	// Get App config
	appConfig := config.GetAppConfig()

	appConfig.SaveMalConfig(&models.MalConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AccessToken:  token,
	})

	fmt.Println("Authentication successful. Access token has been saved.")
}

// getAuthenticationURL retrieves the authentication URL with code_challenge
func getAuthenticationURL(clientId, codeChallenge string) string {
	return fmt.Sprintf("https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s", clientId, codeChallenge)
}

// getAccessToken exchanges the code for an access token
func getAccessToken(clientId, clientSecret, authorizationCode, codeVerifier string) (string, error) {
	tokenEndpoint := "https://myanimelist.net/v1/oauth2/token"

	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("code", authorizationCode)
	data.Set("code_verifier", codeVerifier)
	data.Set("grant_type", "authorization_code")

	// Send a POST request to obtain the token.
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", &models.AppError{
			Message: "Failed to request the access token",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &models.AppError{
			Message: "Failed to read the access token response body",
			Err:     err,
		}
	}

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		return "", &models.AppError{
			Message: "Access token request failed " + fmt.Sprintf("Error: %s", body),
		}
	}

	// Parse the JSON response to get the access token
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", &models.AppError{
			Message: "Failed to parse the access token response",
			Err:     err,
		}
	}

	return tokenResponse.AccessToken, nil
}

// GenerateCodeVerifier generates a code verifier for OAuth2 PKCE
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
