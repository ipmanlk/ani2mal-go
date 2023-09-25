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
		res, err := sendRequest(url, bearerToken)

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

func sendRequest(url string, bearerToken string) (*http.Response, error) {
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

	formattedList := make([]models.Media, 0)

	if listType == "animelist" {
		for _, i := range list.Data {
			formattedList = append(formattedList, models.Media{
				ID:       i.Node.ID,
				Progress: i.ListStatus.NumEpisodesWatched,
				Score:    i.ListStatus.Score,
				Status:   malStatuses[i.ListStatus.Status],
				Repeat:   i.ListStatus.IsRewatching,
				Type:     "anime",
				Length:   i.Node.NumEpisodes,
			})
		}
	} else {
		for _, i := range list.Data {
			formattedList = append(formattedList, models.Media{
				ID:       i.Node.ID,
				Progress: i.ListStatus.NumChaptersRead,
				Score:    i.ListStatus.Score,
				Status:   malStatuses[i.ListStatus.Status],
				Repeat:   i.ListStatus.IsRereading,
				Type:     "manga",
				Length:   i.Node.NumChapters,
			})
		}
	}

	return &formattedList
}
