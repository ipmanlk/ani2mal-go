package models

type AppError struct {
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

// general format to store anime / manga
// ID refers the the MAL ID
type Media struct {
	ID       int     `json:"id"`
	Length   int     `json:"length,omitempty"`
	Progress int     `json:"progress"`
	Score    float64 `json:"score"`
	Status   string  `json:"status"`
	Repeat   bool     `json:"repeat,omitempty"`
	Type     string  `json:"type"`
}

type SourceEntries struct {
	Anime *[]Media `json:"anime"`
	Manga *[]Media `json:"manga"`
}
