package main

import (
	"fmt"
	"ipmanlk/ani2mal/anilist"
	"ipmanlk/ani2mal/config"
	"ipmanlk/ani2mal/mal"
)

func main() {
	fmt.Println("Hello, World!")

	anilistData, _ := anilist.GetUserData("CrystalBullet")

	token := config.GetAppConfig().GetMalConfig().TokenRes.AccessToken

	malData, _ := mal.GetUserData(token)

	mal.SyncData(token, anilistData, malData)

	// jsonData, _ := json.MarshalIndent(res, "", " ")

	// os.WriteFile("lala.json", jsonData, 0644)
}
