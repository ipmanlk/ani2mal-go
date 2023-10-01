package anilist

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/models"
	"math"
	"net/http"
	"strings"
	"time"
)

type graphQLRequest struct {
	Query string `json:"query"`
}

func GetUserData(username string, bearerToken *string) (*models.SourceData, error) {
	anilistAnime, err := getList(username, models.MediaTypeAnime, bearerToken)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch Anilist Anime List",
			Err:     err,
		}
	}

	anilistManga, err := getList(username, models.MediaTypeManga, bearerToken)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch Anilist Manga List",
			Err:     err,
		}
	}

	stats := models.SourceStats{}
	entriesMap := make(map[int]models.Media)
	formattedAnime := formatListResponse(anilistAnime, models.MediaTypeAnime, &stats, entriesMap)
	formattedManga := formatListResponse(anilistManga, models.MediaTypeManga, &stats, entriesMap)

	return &models.SourceData{
		Stats:    stats,
		MediaMap: entriesMap,
		Anime:    formattedAnime,
		Manga:    formattedManga,
	}, nil
}

func getList(username string, mediaType models.MediaType, bearerToken *string) (*models.AnilistRes, error) {
	query := getGraphQuery(username, mediaType)

	requestBody := graphQLRequest{
		Query: query,
	}

	reqBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://graphql.anilist.co", strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if bearerToken != nil {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to construct Anilist request",
			Err:     err,
		}
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to contact Anilist API",
			Err:     err,
		}
	}
	defer res.Body.Close()

	var anilistRes models.AnilistRes

	err = json.NewDecoder(res.Body).Decode(&anilistRes)
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to Parse Anilist Response",
			Err:     err,
		}
	}

	return &anilistRes, nil
}

func formatListResponse(res *models.AnilistRes, mediaType models.MediaType, stats *models.SourceStats, entriesMap map[int]models.Media) []models.Media {
	formattedList := make([]models.Media, 0)

	for _, list := range res.Data.MediaListCollection.Lists {
		if list.IsCustomList {
			continue
		}

		var status models.MediaStatus

		switch strings.ToLower(list.Name) {
		case "planning":
			status = models.MediaStatusPlanning
		case "paused":
			status = models.MediaStatusPaused
		case "watching", "reading":
			status = models.MediaStatusCurrent
		case "dropped":
			status = models.MediaStatusDropped
		case "completed":
			status = models.MediaStatusCompleted
		}

		if status == "watching" || status == "reading" {
			status = models.MediaStatusCurrent
		}

		for _, i := range list.Entries {
			if i.Media.IDMal == nil {
				continue
			}

			repeat := false
			if i.Repeat == 1 {
				repeat = true
			}

			media := models.Media{
				ID:       *i.Media.IDMal,
				Title:    i.Media.Title.Romaji,
				Progress: i.Progress,
				Score:    int(math.Round(i.Score)),
				Status:   status,
				Repeat:   repeat,
				Type:     mediaType,
				Length:   getMediaLength(&i.Media),
			}

			formattedList = append(formattedList, media)
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
	}

	return formattedList
}

func getGraphQuery(username string, mediaType models.MediaType) string {
	anilistMediaType := "ANIME"

	if mediaType == models.MediaTypeManga {
		anilistMediaType = "MANGA"
	}

	return fmt.Sprintf(`{
      MediaListCollection(userName: "%s", type: %s) {
        lists {
          entries {
            id
            status
            score(format: POINT_10)
            progress
            notes
            repeat
            media {
              chapters
              volumes
              idMal
              episodes
              title { romaji }
            }
          }
          name
          isCustomList
          isSplitCompletedList
          status
        }
      }
  }`, username, anilistMediaType)
}

func getMediaLength(media *models.AnilistMedia) int {
	if media.Chapters != nil {
		return *media.Chapters
	} else if media.Episodes != nil {
		return *media.Episodes
	}
	return 0
}
