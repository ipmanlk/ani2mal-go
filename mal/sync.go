package mal

import (
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
			if anilistMedia == malData.MediaMap[malId] {
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
		if _, ok := anilistData.MediaMap[malId]; ok {
			continue
		}
		// entry does not exist in anilist
		removedMedia = append(removedMedia, malMedia)
	}

	log.Printf("Added Media: %d", len(addedMedia))
	log.Printf("Removed Media: %d", len(removedMedia))
	log.Printf("Updated Media: %d", len(updatedMedia))
}
