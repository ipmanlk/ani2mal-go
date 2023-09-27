package models

// Configuration file format
type AnilistConfig struct {
	Username     string   `json:"username"`
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	TokenRes     TokenRes `json:"token_res"`
}

// Response from Anilist Lists (Anime, Manga)
type AnilistRes struct {
	Data AnilistResData `json:"data"`
}

type AnilistResData struct {
	MediaListCollection AnilistMediaListCollection `json:"MediaListCollection"`
}

type AnilistMediaListCollection struct {
	Lists []AnilistList `json:"lists"`
}

type AnilistList struct {
	Entries              []AnilistEntry `json:"entries"`
	Name                 string         `json:"name"`
	IsCustomList         bool           `json:"isCustomList"`
	IsSplitCompletedList bool           `json:"isSplitCompletedList"`
	Status               string         `json:"status"`
}

type AnilistEntry struct {
	ID       int          `json:"id"`
	Status   string       `json:"status"`
	Score    float64      `json:"score"`
	Progress int          `json:"progress"`
	Notes    *string      `json:"notes"`
	Repeat   int          `json:"repeat"`
	Media    AnilistMedia `json:"media"`
}

type AnilistMedia struct {
	Chapters *int         `json:"chapters"`
	Volumes  *int         `json:"volumes"`
	IDMal    *int         `json:"idMal"`
	Episodes *int         `json:"episodes"`
	Title    AnilistTitle `json:"title"`
}

type AnilistTitle struct {
	Romaji string `json:"romaji"`
}
