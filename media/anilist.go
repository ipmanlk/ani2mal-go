package media

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

func GetAnilistEntries(username string) (*models.SourceEntries, error) {
	anilistAnime, err := getAnilistList(username, "ANIME")
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch Anilist Anime List",
			Err:     err,
		}
	}

	anilistManga, err := getAnilistList(username, "MANGA")
	if err != nil {
		return nil, &models.AppError{
			Message: "Failed to fetch Anilist Manga List",
			Err:     err,
		}
	}

	return &models.SourceEntries{
		Anime: anilistAnime,
		Manga: anilistManga,
	}, nil
}

func getAnilistList(username string, mediaType string) (*[]models.Media, error) {
	req, err := getAnilistRequestOptions(username, mediaType)

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

	return formatAnilistListRes(&anilistRes, mediaType), nil
}

func formatAnilistListRes(res *models.AnilistRes, mediaType string) *[]models.Media {
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
				Score:    int(math.Round(i.Score)),
				Status:   status,
				Repeat:   repeat,
				Type:     mType,
				Length:   getAnilistMediaLength(&i.Media),
			})
		}
	}

	return &formattedList
}

func getAnilistRequestOptions(username string, mediaType string) (*http.Request, error) {
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

func getAnilistMediaLength(media *models.AnilistMedia) int {
	if media.Chapters != nil {
		return *media.Chapters
	} else if media.Episodes != nil {
		return *media.Episodes
	}
	return 0
}
