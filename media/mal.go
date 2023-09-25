package media

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/models"
	"net/http"
	"time"
)

func GetMalEntries(bearerToken string) (*models.SourceEntries, error) {
	malAnime, err := getMalList("animelist", bearerToken)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch MAL Anime List",
			Err:     err,
		}
	}

	malManga, err := getMalList("mangalist", bearerToken)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch MAL Manga List",
			Err:     err,
		}
	}

	return &models.SourceEntries{
		Anime: malAnime,
		Manga: malManga,
	}, nil
}

func getMalList(listType string, bearerToken string) (*[]models.Media, error) {
	// Define the MAL API URL
	baseURL := fmt.Sprintf("https://api.myanimelist.net/v2/users/@me/%s", listType)
	url := baseURL + "?fields=list_status,num_episodes,num_chapters&limit=1000&nsfw=true"

	var allMedia []models.MalDatum

	// Loop to fetch all pages
	for url != "" {
		res, err := sendMalGetRequest(url, bearerToken)

		if err != nil {
			return nil, &models.AppError{
				Message: "Failed to fetch MAL list",
				Err:     err,
			}
		}
		defer res.Body.Close()

		var malList models.MalListRes
		err = json.NewDecoder(res.Body).Decode(&malList)
		if err != nil {
			return nil, &models.AppError{
				Message: "Failed to parse MAL list response",
				Err:     err,
			}
		}

		allMedia = append(allMedia, malList.Data...)

		// Check if there is a next page
		url = malList.Paging.Next
	}

	// Create a MalListRes struct to hold the combined results
	var combinedList models.MalListRes
	combinedList.Data = allMedia

	formattedList := formatMalListRes(&combinedList, listType)

	return formattedList, nil
}

func sendMalGetRequest(url string, bearerToken string) (*http.Response, error) {
	timeout := 15 * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func formatMalListRes(list *models.MalListRes, listType string) *[]models.Media {
	malStatuses := map[string]string{
		"plan_to_watch": "planning",
		"on_hold":       "paused",
		"watching":      "current",
		"dropped":       "dropped",
		"completed":     "completed",
		"plan_to_read":  "planning",
		"reading":       "current",
	}

	formattedList := make([]models.Media, len(list.Data))

	for i, item := range list.Data {
		mediaType := "anime"
		progress := item.ListStatus.NumEpisodesWatched
		length := item.Node.NumEpisodes
		repeat := item.ListStatus.IsRewatching

		if listType == "mangalist" {
			mediaType = "manga"
			progress = item.ListStatus.NumChaptersRead
			length = item.Node.NumChapters
			repeat = item.ListStatus.IsRereading
		}

		formattedList[i] = models.Media{
			ID:       item.Node.ID,
			Progress: progress,
			Score:    item.ListStatus.Score,
			Status:   malStatuses[item.ListStatus.Status],
			Repeat:   repeat,
			Type:     mediaType,
			Length:   length,
		}
	}

	return &formattedList
}
