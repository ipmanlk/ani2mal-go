package models

type AppError struct {
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

// Response after exchanging auth code
type TokenRes struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// general format to store anime / manga
// ID refers the the MAL ID
type Media struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Length   int    `json:"length,omitempty"`
	Progress int    `json:"progress"`
	Score    int    `json:"score"`
	Status   string `json:"status"`
	Repeat   bool   `json:"repeat,omitempty"`
	Type     string `json:"type"`
}

type SourceStats struct {
	Planning  int `json:"planning"`
	Paused    int `json:"paused"`
	Current   int `json:"current"`
	Dropped   int `json:"dropped"`
	Completed int `json:"completed"`
}

type SourceData struct {
	Stats    SourceStats   `json:"stats"`
	MediaMap map[int]Media `json:"media_map"`
	Anime    []Media       `json:"anime"`
	Manga    []Media       `json:"manga"`
}
