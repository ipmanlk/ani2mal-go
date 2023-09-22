package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"ipmanlk/ani2mal/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

func AuthMal() {
	fmt.Print("Enter Client ID: ")
	clientId := utils.GetStrInput()

	fmt.Print("Enter Client Secret: ")
	clientSecret := utils.GetStrInput()

	// Use the getAuthenticationURL function to retrieve the login URL
	loginURL := getAuthenticationURL(clientId)
	fmt.Printf("Login URL: %s\n", loginURL)

	fmt.Print("Enter the code from the login URL: ")
	code := utils.GetStrInput()

	token, err := getToken(clientId, clientSecret, code)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}

	err = saveTokenToFile(token, "token.txt")
	if err != nil {
		fmt.Printf("Error saving token: %v\n", err)
		return
	}

	fmt.Println("Authentication successful. Token saved to token.txt")
}

// getAuthenticationURL retrieves the authentication URL with code_challenge
func getAuthenticationURL(clientId string) string {
	codeVerifier := generateCodeVerifier()
	codeChallenge := codeVerifier

	return fmt.Sprintf("https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s", clientId, codeChallenge)
}

// getToken exchanges the code for an access token
func getToken(clientId, clientSecret, code string) (string, error) {
	tokenEndpoint := "https://myanimelist.net/v1/oauth2/token"

	// Prepare the request body
	requestBody := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code", clientId, clientSecret, code)

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error: %s", body)
	}

	// Parse the JSON response to get the access token
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// saveTokenToFile saves the token to a file.
func saveTokenToFile(token, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return err
	}

	return nil
}

// GenerateCodeVerifier generates a code verifier for OAuth2 PKCE
func generateCodeVerifier() string {
	// Generate a 32-byte (256-bit) random value
	verifierBytes := make([]byte, 32)
	_, err := rand.Read(verifierBytes)
	if err != nil {
		log.Fatalln("Error generating random string: ", err)
	}

	// Encode the random bytes as a URL-safe base64 string
	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
	return codeVerifier
}
