package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

// GetDefaultBotMessage func
func GetDefaultBotMessage(userMessage string) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "Hello,",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "World!",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Size: linebot.FlexTextSizeTypeLg,
					Text: "您剛剛說:[" + userMessage + "]",
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypePrimary,
					Action: linebot.NewURIAction("See my bot actions", lineLIFFURL),
				},
			},
		},
	}
}

// GetImageUploadBotMessage func
func GetImageUploadBotMessage(liffID, deleteHash string) *linebot.BubbleContainer {
	lineLIFFURLForImgur := lineLIFFURLBase + liffID
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "上傳成功! 您可以直接複製上面連結或...",
					Wrap: true,
				},
				&linebot.SeparatorComponent{
					Type: linebot.FlexComponentTypeSeparator,
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypePrimary,
					Action: linebot.NewURIAction("預覽上傳的圖片", lineLIFFURLForImgur),
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Action: linebot.NewPostbackAction(
						"刪除上傳的圖片",
						"deletehash="+deleteHash+"&liffid="+liffID,
						"刪除上傳的圖片",
						""),
				},
			},
		},
	}
}

// ActionsOnTextMessage func, to do specific action depend on message received and the last bot sent text message.
func ActionsOnTextMessage(message *linebot.TextMessage, replyToken, userID string) {
	log.Println("message.Text" + message.Text)
	switch lastBotMessages {
	case "請問您要找的帳號是?":
		lastBotMessages = ""
		//Get account info and send to reply message
		accountID := "BensonLiao"
		if message.Text != "" {
			accountID = message.Text
		}
		accountRes, err := imgurClient.GetAccount(accountID)
		if err != nil {
			log.Print(err)
		}
		if !accountRes.Success {
			if _, err = bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("抱歉，找不到"+accountID+"~:(")).Do(); err != nil {
				log.Print(err)
			}
			break
		}
		accountData, err := json.MarshalIndent(accountRes.Data, "", "\t")
		if err != nil {
			log.Print(err)
		}
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("找到"+accountID+"了~\n"+string(accountData))).Do(); err != nil {
			log.Print(err)
		}
	default:
		if lastBotMessages != "" {
			if err := ReplyTextMessage(replyToken, replyTextMessage); err != nil {
				log.Print(err)
			}
		}
		contents := GetDefaultBotMessage(message.Text)
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("Hello, World!", contents)).Do(); err != nil {
			log.Print(err)
		}
	}
}

// ActionsOnImageMessage func, to do specific action depend on message received and the last bot sent image message.
func ActionsOnImageMessage(message *linebot.ImageMessage, replyToken, userID string) {
	log.Println("message.ID" + message.ID)
	switch lastBotMessages {
	case "請問您要上傳哪張圖片?":
		lastBotMessages = ""
		res, err := bot.GetMessageContent(message.ID).Do()
		if err != nil {
			log.Print(err)
		}
		body := res.Content
		defer body.Close()
		imgContent, err := ioutil.ReadAll(body)
		if err != nil {
			log.Print(err)
		}
		log.Printf("imgContent's type: %T\n", imgContent)
		imgUploadRes, err := imgurClient.AnonymousUploadByImgMessage(imgContent)
		if err != nil {
			log.Print(err)
		}
		log.Printf("image upload success? %v", imgUploadRes.Success)
		imgUploadData, err := json.MarshalIndent(imgUploadRes.Data, "", "\t")
		if err != nil {
			log.Print(err)
		}
		log.Printf("image upload result: %s", string(imgUploadData))

		// Call Line API to get LIFF URL
		preview := linebot.View{
			Type: linebot.LIFFViewTypeTall,
			URL:  imgUploadRes.Data.Link,
		}
		lineLIFFAddRes, err := bot.AddLIFF(preview).Do()
		if err != nil {
			log.Print(err)
		}
		// Send flex message to show result and furthur actions.
		// But flex message are not allow to be copy by user currently,
		// and if we send another message by push it will be charged
		// when over 500 messages per month under FREE plan.
		// So if you want to save money for your bot,
		// sending multiple messages(up to 5) in 1 reply.
		flexContent := GetImageUploadBotMessage(lineLIFFAddRes.LIFFID, imgUploadRes.Data.DeleteHash)
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(imgUploadRes.Data.Link),
			linebot.NewFlexMessage("您的圖片連結: "+imgUploadRes.Data.Link, flexContent)).Do(); err != nil {
			log.Print(err)
		}
	default:
		if lastBotMessages != "" {
			if err := ReplyTextMessage(replyToken, replyTextMessage); err != nil {
				log.Print(err)
			}
		}
		res, err := bot.GetMessageContent(message.ID).Do()
		if err != nil {
			log.Print(err)
		}
		body := res.Content
		defer body.Close()
		bodyGot, err := ioutil.ReadAll(body)
		if err != nil {
			log.Print(err)
		}
		log.Printf("bodyGot's type: %T\n", bodyGot)
	}
}
