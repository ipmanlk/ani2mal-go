package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/api"
	"ipmanlk/ani2mal/config"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	// res, err := media.GetAnilistEntries("CrystalBullet")

	// if err != nil {
	// 	panic(err)
	// }

	token := config.GetAppConfig().GetMalConfig().TokenRes.AccessToken

	res, _ := api.GetMalEntries(token)

	jsonData, _ := json.MarshalIndent(res, "", " ")

	os.WriteFile("lala.json", jsonData, 0644)
}
