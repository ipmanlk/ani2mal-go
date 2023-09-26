package anilist

import (
	"fmt"
	"ipmanlk/ani2mal/config"
	"ipmanlk/ani2mal/models"
	"ipmanlk/ani2mal/utils"
)

func performAuth() {
	fmt.Print("Enter Anilist Username: ")
	username := utils.GetStrInput()

	config.GetAppConfig().SaveAnilistConfig(
		&models.AnilistConfig{
			Username: username,
		},
	)
	
	fmt.Println("Anilist configuration has been saved")
}
