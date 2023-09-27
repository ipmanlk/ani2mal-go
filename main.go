package main

import (
	"fmt"
	"ipmanlk/ani2mal/anilist"
	"ipmanlk/ani2mal/mal"
)

func main() {
	fmt.Println("Hello, World!")

	anilistCode, _ := anilist.GetAccessCode()
	anilistData, _ := anilist.GetUserData("CrystalBullet", &anilistCode)

	malCode, _ := mal.GetAccessCode()
	malData, _ := mal.GetUserData(malCode)

	mal.SyncData(malCode, anilistData, malData)
	// jsonData, _ := json.MarshalIndent(res, "", " ")

	// os.WriteFile("lala.json", jsonData, 0644)
}
