package mal

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const malApiUrl = "https://api.myanimelist.net/v2"

// Media status for each MAL API status
var mediaStatuses = map[string]models.MediaStatus{
	"plan_to_watch": models.MediaStatusPlanning,
	"on_hold":       models.MediaStatusPaused,
	"watching":      models.MediaStatusCurrent,
	"dropped":       models.MediaStatusDropped,
	"completed":     models.MediaStatusCompleted,
	"plan_to_read":  models.MediaStatusPlanning,
	"reading":       models.MediaStatusCurrent,
}

// MAL API status for each Media status
var malStatuses = map[models.MediaStatus]string{
	models.MediaStatusPlanning:  "plan_to_watch",
	models.MediaStatusCurrent:   "watching",
	models.MediaStatusCompleted: "completed",
	models.MediaStatusPaused:    "on_hold",
	models.MediaStatusDropped:   "dropped",
}

func GetUserData(bearerToken string) (*models.SourceData, error) {
	malAnime, err := getList(models.MAL_ANIME_LIST, bearerToken)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch MAL Anime List",
			Err:     err,
		}
	}

	malManga, err := getList(models.MAL_MANGA_LIST, bearerToken)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch MAL Manga List",
			Err:     err,
		}
	}

	stats := models.SourceStats{}
	entriesMap := make(map[int]models.Media)
	formattedAnime := formatListResponse(malAnime, models.MAL_ANIME_LIST, &stats, entriesMap)
	formattedManga := formatListResponse(malManga, models.MAL_MANGA_LIST, &stats, entriesMap)

	return &models.SourceData{
		Stats:    stats,
		MediaMap: entriesMap,
		Anime:    formattedAnime,
		Manga:    formattedManga,
	}, nil
}

func UpdateAnime(bearerToken string, entry models.Media) error {
	requestUrl := fmt.Sprintf("%s/anime/%d/my_list_status", malApiUrl, entry.ID)

	data := url.Values{}
	data.Set("status", getMalStatus(entry.Status, models.MediaTypeAnime))
	data.Set("num_watched_episodes", strconv.Itoa(entry.Progress))
	data.Set("score", strconv.Itoa(entry.Score))

	return sendPutRequest(requestUrl, bearerToken, data)
}

func DeleteAnime(bearerToken string, entry models.Media) error {
	url := fmt.Sprintf("%s/anime/%d/my_list_status", malApiUrl, entry.ID)
	return sendDeleteRequest(url, bearerToken)
}

func UpdateManga(bearerToken string, entry models.Media) error {
	data := url.Values{}
	data.Set("status", getMalStatus(entry.Status, models.MediaTypeManga))
	data.Set("num_chapters_read", strconv.Itoa(entry.Progress))
	data.Set("score", strconv.Itoa(entry.Score))
	requestUrl := fmt.Sprintf("%s/manga/%d/my_list_status", malApiUrl, entry.ID)
	return sendPutRequest(requestUrl, bearerToken, data)
}

func DeleteManga(bearerToken string, entry models.Media) error {
	url := fmt.Sprintf("%s/manga/%d/my_list_status", malApiUrl, entry.ID)
	return sendDeleteRequest(url, bearerToken)
}

func getList(malListType models.MalListType, bearerToken string) (*models.MalListRes, error) {
	listType := "animelist"

	if malListType == models.MAL_MANGA_LIST {
		listType = "mangalist"
	}

	baseURL := fmt.Sprintf("%s/users/@me/%s", malApiUrl, listType)
	url := baseURL + "?fields=list_status,num_episodes,num_chapters&limit=1000&nsfw=true"

	var allMedia []models.MalDatum

	// Loop to fetch all pages
	for url != "" {
		res, err := sendGetRequest(url, bearerToken)

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

	return &combinedList, nil
}

func sendGetRequest(url string, bearerToken string) (*http.Response, error) {
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

func sendPutRequest(url string, bearerToken string, data url.Values) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to update MAL entry, status code: %d", res.StatusCode)
	}

	return nil
}

func sendDeleteRequest(url string, bearerToken string) error {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to delete MAL entry, status code: %d", res.StatusCode)
	}

	return nil
}

func formatListResponse(list *models.MalListRes, listType models.MalListType, stats *models.SourceStats, entriesMap map[int]models.Media) []models.Media {
	formattedList := make([]models.Media, len(list.Data))

	for i, item := range list.Data {
		mediaType := models.MediaTypeAnime
		progress := item.ListStatus.NumEpisodesWatched
		length := item.Node.NumEpisodes
		repeat := item.ListStatus.IsRewatching
		status := mediaStatuses[item.ListStatus.Status]

		if listType == models.MAL_MANGA_LIST {
			mediaType = models.MediaTypeManga
			progress = item.ListStatus.NumChaptersRead
			length = item.Node.NumChapters
			repeat = item.ListStatus.IsRereading
		}

		media := models.Media{
			ID:       item.Node.ID,
			Title:    item.Node.Title,
			Progress: progress,
			Score:    item.ListStatus.Score,
			Status:   status,
			Repeat:   repeat,
			Type:     mediaType,
			Length:   length,
		}

		formattedList[i] = media
		entriesMap[media.ID] = media

		// update stats reference data
		switch status {
		case models.MediaStatusPlanning:
			stats.Planning += 1
		case models.MediaStatusPaused:
			stats.Paused += 1
		case models.MediaStatusCurrent:
			stats.Current += 1
		case models.MediaStatusDropped:
			stats.Dropped += 1
		case models.MediaStatusCompleted:
			stats.Completed += 1
		}
	}

	return formattedList
}

func getMalStatus(status models.MediaStatus, mediaType models.MediaType) string {
	malStatus := malStatuses[status]

	if mediaType == models.MediaTypeManga {
		if status == models.MediaStatusPlanning {
			malStatus = "plan_to_read"
		}
		if status == models.MediaStatusCurrent {
			malStatus = "reading"
		}
	}

	return malStatus
}
