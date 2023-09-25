package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/models"
	"net/http"
	"time"
)

const malApiUrl = "https://api.myanimelist.net/v2"

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

func UpdateAnime(bearerToken string, entry models.Media) error {
	url := fmt.Sprintf("%s/anime/%d/my_list_status", malApiUrl, entry.ID)
	return sendMalPutRequest(url, bearerToken, map[string]interface{}{
		"status":               getMalStatus(entry.Status, "ANIME"),
		"num_watched_episodes": entry.Progress,
		"score":                entry.Score,
	})
}

func DeleteAnime(bearerToken string, entry models.Media) error {
	url := fmt.Sprintf("%s/anime/%d/my_list_status", malApiUrl, entry.ID)
	return sendMalDeleteRequest(url, bearerToken)
}

func UpdateManga(bearerToken string, entry models.Media) error {
	url := fmt.Sprintf("%s/manga/%d/my_list_status", malApiUrl, entry.ID)
	return sendMalPutRequest(url, bearerToken, map[string]interface{}{
		"status":            getMalStatus(entry.Status, "MANGA"),
		"num_chapters_read": entry.Progress,
		"score":             entry.Score,
	})
}

func DeleteManga(bearerToken string, entry models.Media) error {
	url := fmt.Sprintf("%s/manga/%d/my_list_status", malApiUrl, entry.ID)
	return sendMalDeleteRequest(url, bearerToken)
}

func getMalList(listType string, bearerToken string) (*[]models.Media, error) {
	baseURL := fmt.Sprintf("%s/users/@me/%s", malApiUrl, listType)
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

func sendMalPutRequest(url string, bearerToken string, data map[string]interface{}) error {
	client := &http.Client{}
	body, _ := json.Marshal(data)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to update MAL entry, status code: %d", res.StatusCode)
	}

	return nil
}

func sendMalDeleteRequest(url string, bearerToken string) error {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to delete MAL entry, status code: %d", res.StatusCode)
	}

	return nil
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

func getMalStatus(status string, mediaType string) string {
	malStatuses := map[string]string{
		"planning":  "plan_to_watch",
		"current":   "watching",
		"completed": "completed",
		"paused":    "on_hold",
		"dropped":   "dropped",
	}

	if mediaType == "MANGA" {
		malStatuses["planning"] = "plan_to_read"
		malStatuses["current"] = "reading"
	}

	return malStatuses[status]
}
