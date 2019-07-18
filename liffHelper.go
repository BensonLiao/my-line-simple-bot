package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

// DeleteLIFFApp func
func DeleteLIFFApp(liffID string) error {
	if _, err := bot.DeleteLIFF(liffID).Do(); err != nil {
		return fmt.Errorf(
			"[DeleteLIFF] error in calling bot.DeleteLIFF: %v",
			err)
	}
	log.Println("Delete LIFF App URL success!")
	return nil
}

// AddLIFFApp func
func AddLIFFApp(
	viewType linebot.LIFFViewType,
	link string) (liffID *linebot.LIFFIDResponse, err error) {
	log.Printf("LIFF URL: %v", link)
	preview := linebot.View{
		Type: viewType,
		URL:  link,
	}
	return bot.AddLIFF(preview).Do()
}

// UpdateLIFFApp func
func UpdateLIFFApp(
	viewType linebot.LIFFViewType,
	liffID,
	link string) (status *linebot.BasicResponse, err error) {
	log.Printf("Update LIFF URL: %v", link)
	newView := linebot.View{
		Type: viewType,
		URL:  link,
	}
	return bot.UpdateLIFF(liffID, newView).Do()
}

// GetLIFFAppID func
func GetLIFFAppID(link string) (liffID string, err error) {
	liffID = ""
	res, err := bot.GetLIFF().Do()
	if err != nil {
		return liffID, fmt.Errorf(
			"[GetLIFFAppID] error in calling bot.GetLIFF: %v",
			err)
	}
	for _, app := range res.Apps {
		if strings.Contains(app.View.URL, imgurGalleriesURLBase) {
			liffID = app.LIFFID
		}
	}
	return liffID, nil
}
