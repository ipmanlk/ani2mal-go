package main

import (
	"fmt"
	"ipmanlk/ani2mal/anilist"
	"ipmanlk/ani2mal/mal"
)

func main() {
	fmt.Println("Hello, World!")

	anilistCode, _ := anilist.GetAccessCode()
	anilistData, err := anilist.GetUserData("CrystalBullet", &anilistCode)

	if err != nil {
		panic(err)
	}

	malCode, _ := mal.GetAccessCode()
	malData, err := mal.GetUserData(malCode)

	if err != nil {
		panic(err)
	}

	mal.SyncData(malCode, anilistData, malData)
	// jsonData, _ := json.MarshalIndent(res, "", " ")

	// os.WriteFile("lala.json", jsonData, 0644)
}
