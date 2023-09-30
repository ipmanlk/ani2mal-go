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

	for malId, anilistMedia := range anilistData.MediaMap {
		// entry exist in both media maps
		if _, ok := malData.MediaMap[malId]; ok {
			// check if entry is the same
			if isMediaEqual(anilistMedia, malData.MediaMap[malId]) {
				continue
			}
			// otherwise entry is modified
			updatedMedia = append(updatedMedia, anilistMedia)

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
		if media.Type == "anime" {
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
		if media.Type == "anime" {
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
	fmt.Printf("%s\n", media1.Title)

	if media1.ID != media2.ID {
		fmt.Printf("ID mismatch %d %d\n", media1.ID, media2.ID)
	}

	if media1.Progress != media2.Progress {
		fmt.Printf("Progress mismatch %d %d\n", media1.Progress, media2.Progress)
	}

	if media1.Score != media2.Score {
		fmt.Printf("Score mismatch %d %d\n", media1.Score, media2.Score)
	}

	if media1.Status != media2.Status {
		fmt.Printf("Status mismatch %s %s\n", media1.Status, media2.Status)
	}

	fmt.Print("\n===================\n")

	return media1.ID == media2.ID &&
		media1.Progress == media2.Progress &&
		media1.Score == media2.Score &&
		media1.Status == media2.Status
}
