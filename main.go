package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/config"
	"ipmanlk/ani2mal/media"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	// res, err := media.GetAnilistEntries("CrystalBullet")

	// if err != nil {
	// 	panic(err)
	// }

	token := config.GetAppConfig().GetMalConfig().TokenRes.AccessToken

	res, _ := media.GetMalEntries(token)

	jsonData, _ := json.MarshalIndent(res, "", " ")

	os.WriteFile("lala.json", jsonData, 0644)
}
