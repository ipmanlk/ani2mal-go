package models

// Configuration file format
type MalConfig struct {
	ClientId     string      `json:"client_id"`
	ClientSecret string      `json:"client_secret"`
	TokenRes     MalTokenRes `json:"token_res"`
}

// Response from MAL auth
type MalTokenRes struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Response from MAL Lists (Anime, Manga)
type MalListRes struct {
	Data   []MalDatum    `json:"data"`
	Paging MalListPaging `json:"paging,omitempty"`
}

type MalDatum struct {
	Node       MalNode       `json:"node"`
	ListStatus MalListStatus `json:"list_status"`
}

type MalListStatus struct {
	Status             string `json:"status"`
	Score              int    `json:"score"`
	NumEpisodesWatched int    `json:"num_episodes_watched"`
	NumChaptersRead    int    `json:"num_chapters_read"`
	IsRewatching       bool   `json:"is_rewatching"`
	UpdatedAt          string `json:"updated_at"`
	IsRereading        bool   `json:"is_rereading"`
}

type MalNode struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	MainPicture MalNodeMainPicture `json:"main_picture"`
	NumEpisodes int                `json:"num_episodes"`
	NumChapters int                `json:"num_chapters"`
}

type MalNodeMainPicture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type MalListPaging struct {
	Next string `json:"next"`
}
