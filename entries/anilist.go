package entries

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/models"
	"net/http"
	"strings"
	"time"
)

type graphQLRequest struct {
	Query string `json:"query"`
}

func GetList(username string, mediaType string) (*[]models.Media, error) {
	req, err := getRequestOptions(username, mediaType)

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

	return formatListRes(&anilistRes, mediaType), nil
}

func formatListRes(res *models.AnilistRes, mediaType string) *[]models.Media {
	mType := strings.ToLower(mediaType)
	formattedList := make([]models.Media, 0)

	for _, list := range res.Data.MediaListCollection.Lists {
		if list.IsCustomList {
			continue
		}

		status := strings.ToLower(list.Name)

		if status == "watching" || status == "reading" {
			status = "current"
		}

		for _, i := range list.Entries {
			if i.Media.IDMal == nil {
				continue
			}

			repeat := false
			if i.Repeat == 1 {
				repeat = true
			}

			formattedList = append(formattedList, models.Media{
				ID:       *i.Media.IDMal,
				Progress: i.Progress,
				Score:    i.Score,
				Status:   status,
				Repeat:   repeat,
				Type:     mType,
				Length:   getLength(i.Media.Chapters, i.Media.Episodes),
			})
		}
	}

	return &formattedList
}

func getRequestOptions(username string, mediaType string) (*http.Request, error) {
	query := fmt.Sprintf(`{
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
  }`, username, mediaType)

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

	return req, nil
}

func getLength(chapters, episodes *int) int {
	if chapters != nil {
		return *chapters
	} else if episodes != nil {
		return *episodes
	}
	return 0
}

func roundScore(score float64) int {
	return int(score + 0.5)
}
