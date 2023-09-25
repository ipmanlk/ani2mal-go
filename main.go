package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/media"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	res, err := media.GetAnilistEntries("CrystalBullet")

	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(res, "", " ")

	os.WriteFile("lala.json", jsonData, 0644)
}
