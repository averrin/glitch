package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/image/bmp"
)

var STEAMID *string

const APIKEY = "90B9CAB39D586E26412485E7E92D3613"
const APPSLIST = "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&include_appinfo=1&include_played_free_games=1&format=json"
const HEADERIMAGE = "http://cdn.akamai.steamstatic.com/steam/apps/%v/header.jpg"

type Game struct {
	Appid                    int    `json:"appid"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
	ImgIconURL               string `json:"img_icon_url"`
	ImgLogoURL               string `json:"img_logo_url"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
}

type SteamResponse struct {
	Response struct {
		GameCount int    `json:"game_count"`
		Games     []Game `json:"games"`
	} `json:"response"`
}

func GetImage(appId int) (ok bool, cache bool) {
	url := fmt.Sprintf(HEADERIMAGE, appId)
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

func GetGames() []Game {
	url := fmt.Sprintf(APPSLIST, APIKEY, *STEAMID)
	log.Print(url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer response.Body.Close()
	var r SteamResponse
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)
	return r.Response.Games
}

func fetch() {
	games := GetGames()
	for _, game := range games {
		_, cache := GetImage(game.Appid)
		if cache {
			fmt.Print("_")
			continue
		}
		fmt.Print(".")
	}
}
