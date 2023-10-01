package mal

import (
	"fmt"
	"ipmanlk/ani2mal/models"
	"log"
)

func SyncData(malBearerToken string, anilistData, malData *models.SourceData) {
	addedMedia := make([]models.Media, 0)
	removedMedia := make([]models.Media, 0)
	updatedMedia := make([]models.Media, 0)

	updatedStuff := make([]map[string]models.Media, 0)

	for malId, anilistMedia := range anilistData.MediaMap {
		// entry exist in both media maps
		if _, ok := malData.MediaMap[malId]; ok {
			// check if entry is the same
			if isMediaEqual(anilistMedia, malData.MediaMap[malId]) {
				continue
			}
			// otherwise entry is modified
			updatedMedia = append(updatedMedia, anilistMedia)

			updatedStuff = append(updatedStuff, map[string]models.Media{
				"mal":     malData.MediaMap[malId],
				"anilist": anilistMedia,
			})

			continue

		} else {
			// entry does not exist in mal
			addedMedia = append(addedMedia, anilistMedia)
		}
	}

	// removed media should be checked against anilistData
	for malId, malMedia := range malData.MediaMap {
		if _, ok := anilistData.MediaMap[malId]; !ok {
			// entry does not exist in anilist
			removedMedia = append(removedMedia, malMedia)
		}
	}

	log.Printf("Added Media: %d", len(addedMedia))
	log.Printf("Removed Media: %d", len(removedMedia))
	log.Printf("Updated Media: %d", len(updatedMedia))

	// Sync data
	log.Printf("Syncing: Added Media")

	for _, media := range append(addedMedia, updatedMedia...) {
		if media.Type == models.MediaTypeAnime {
			err := UpdateAnime(malBearerToken, media)
			if err != nil {
				fmt.Printf("Failed to update anime %v\n", err)
				continue
			}
		} else {
			err := UpdateManga(malBearerToken, media)
			if err != nil {
				fmt.Printf("Failed to update manga %v\n", err)
				continue
			}
			fmt.Printf("Updated: %s\n", media.Title)
		}
	}

	for _, media := range removedMedia {
		if media.Type == models.MediaTypeAnime {
			err := DeleteAnime(malBearerToken, media)
			if err != nil {
				fmt.Printf("Failed to delete anime %v\n", err)
				continue
			}
		} else {
			err := DeleteManga(malBearerToken, media)
			if err != nil {
				fmt.Printf("Failed to delete manga %v\n", err)
				continue
			}
			fmt.Printf("Deleted: %s\n", media.Title)
		}
	}
}

// TODO: do something about repeat property
func isMediaEqual(media1, media2 models.Media) bool {
	idMatch := media1.ID == media2.ID
	completeMatch := media1.Status == "completed" && media2.Status == "completed"
	statusMatch := media1.Status == media2.Status
	scoreMatch := media1.Score == media2.Score
	progressMatch := media1.Progress == media2.Progress
	lengthMismatch := media1.Length != media2.Length

	return idMatch &&
		((completeMatch && scoreMatch) ||
			(statusMatch && scoreMatch && lengthMismatch) ||
			(progressMatch && scoreMatch && statusMatch))
}
