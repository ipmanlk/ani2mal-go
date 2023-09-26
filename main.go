package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/config"
	"ipmanlk/ani2mal/mal"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	// res, err := anilist.GetData("CrystalBullet")

	// if err != nil {
	// 	panic(err)
	// }

	token := config.GetAppConfig().GetMalConfig().TokenRes.AccessToken

	res, _ := mal.GetData(token)

	jsonData, _ := json.MarshalIndent(res, "", " ")

	os.WriteFile("lala.json", jsonData, 0644)
}
