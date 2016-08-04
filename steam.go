package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	st "github.com/averrin/shodan/modules/steam"
	"golang.org/x/image/bmp"
)

func GetImage(appId int) (ok bool, cache bool) {
	url := fmt.Sprintf(st.HeaderURL, appId)
	fName := fmt.Sprintf("cache/%d.bmp", appId)
	if _, err := os.Stat(fName); err == nil {
		return true, true
	}
	output, err := os.Create(fName)
	if err != nil {
		fmt.Println("Error while creating", fName, "-", err)
		return false, false
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false, false
	}
	defer response.Body.Close()

	img1, err := jpeg.Decode(response.Body)
	if err != nil {
		log.Print(fName, err)
		return
	}
	toimg, _ := os.Create(fName)
	defer toimg.Close()
	bmp.Encode(toimg, img1)
	// _, err = io.Copy(output, response.Body)
	// if err != nil {
	// 	fmt.Println("Error while downloading", url, "-", err)
	// 	return false, false
	// }
	return true, false
}

func fetch() {
	games := steam.GetGames()
	for _, game := range games {
		_, cache := GetImage(game.Appid)
		if cache {
			fmt.Print("_")
			continue
		}
		fmt.Print(".")
	}
}
