package models

type MalTokenRes struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MalConfig struct {
	ClientId     string      `json:"client_id"`
	ClientSecret string      `json:"client_secret"`
	TokenRes     MalTokenRes `json:"token_res"`
}
