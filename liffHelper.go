package main

import (
	"fmt"
	"log"

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
func AddLIFFApp(viewType linebot.LIFFViewType, link string) (liffID *linebot.LIFFIDResponse, err error) {
	log.Printf("LIFF URL: %v", link)
	preview := linebot.View{
		Type: viewType,
		URL:  link,
	}
	return bot.AddLIFF(preview).Do()
}
