package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ipmanlk/ani2mal/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func AuthMal() {
	fmt.Print("Enter Client ID: ")
	clientId := utils.GetStrInput()

	fmt.Print("Enter Client Secret: ")
	clientSecret := utils.GetStrInput()

	// Generate a code verifier and code challenge
	codeVerifier := generateCodeVerifier()

	// Use the getAuthenticationURL function to retrieve the login URL
	loginURL := getAuthenticationURL(clientId, codeVerifier)
	fmt.Printf("Login URL: %s\n", loginURL)

	fmt.Print("Enter the code from the login URL: ")
	code := utils.GetStrInput()

	token, err := getToken(clientId, clientSecret, code, codeVerifier)
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
func getAuthenticationURL(clientId, codeChallenge string) string {
	return fmt.Sprintf("https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s", clientId, codeChallenge)
}

// getToken exchanges the code for an access token
func getToken(clientId, clientSecret, authorizationCode,codeVerifier  string) (string, error) {
	tokenEndpoint := "https://myanimelist.net/v1/oauth2/token"

	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("code", authorizationCode)
	data.Set("code_verifier", codeVerifier)
	data.Set("grant_type", "authorization_code")

	// Send a POST request to obtain the token.
	resp, err := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Unable to request token:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read response body:", err)
		os.Exit(1)
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
