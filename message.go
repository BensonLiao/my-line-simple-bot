package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// GetFormatTime func
func GetFormatTime(t time.Time) string {
	return fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

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
					Size: linebot.FlexTextSizeTypeLg,
					Text: "您剛剛說:[" + userMessage + "]",
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypePrimary,
					Action: linebot.NewURIAction("See my bot actions", lineLIFFURLBotActions),
				},
			},
		},
	}
}

// GetDefaultLIFFWithNoteBotMessage func
func GetDefaultLIFFWithNoteBotMessage(liffID, actionLabel, note string) *linebot.BubbleContainer {
	lineLIFFURLForImgurUser := lineLIFFURLBase + liffID
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypePrimary,
					Action: linebot.NewURIAction(actionLabel, lineLIFFURLForImgurUser),
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: note,
					Wrap: true,
				},
			},
		},
	}
}

// GetDefaulLIFFBotMessage func
func GetDefaulLIFFBotMessage(liffID, actionLabel string) *linebot.BubbleContainer {
	lineLIFFURLForImgurUser := lineLIFFURLBase + liffID
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypePrimary,
					Action: linebot.NewURIAction(actionLabel, lineLIFFURLForImgurUser),
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
	log.Println("message.Text: " + message.Text)
	switch lastBotMessages {
	case "請問您要找的帳號是?":
		lastBotMessages = ""
		//Get account info and send to reply message
		accountID := "BensonLiao"
		if message.Text != "" {
			accountID = message.Text
		}
		imgurRes, err := imgurClient.GetAccount(accountID)
		if err != nil {
			log.Print(err)
		}
		if !imgurRes.Success {
			if _, err = bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("抱歉，找不到"+accountID+"~:(")).Do(); err != nil {
				log.Print(err)
			}
			break
		}
		accountData, err := json.MarshalIndent(imgurRes.Data, "", "\t")
		if err != nil {
			log.Print(err)
		}
		log.Println(string(accountData))
		// Call Line API to get LIFF App
		imgurUserLink := "https://" + accountID + imgurUserURLBase
		lineLIFFAddRes, err := AddLIFFApp(linebot.LIFFViewTypeTall, imgurUserLink)
		if err != nil {
			log.Print(err)
		}
		expiredAt := time.Now().Add(lineLIFFAppDuration).Local()
		expireAtTime := GetFormatTime(expiredAt)
		log.Printf("LIFF app expired at: %v", expireAtTime)
		note := "*預覽連結將於 " + expireAtTime + " 失效"
		flexContent := GetDefaultLIFFWithNoteBotMessage(
			lineLIFFAddRes.LIFFID,
			"看看他/她是誰?",
			note)
		HowToAccessOldWeb := "*imgur新版網頁如用手機瀏覽此頁面會強制回首頁，解法請參考 https://help.imgur.com/hc/en-us/articles/115002122443-Accessing-the-old-version-of-mobile-web"
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("找到"+accountID+"了~"),
			linebot.NewFlexMessage("帳號連結: "+imgurUserLink, flexContent),
			linebot.NewTextMessage(HowToAccessOldWeb)).Do(); err != nil {
			log.Print(err)
		}
		<-time.After(lineLIFFAppDuration) // Timer expired
		if err := DeleteLIFFApp(lineLIFFAddRes.LIFFID); err != nil {
			log.Print(err)
		}
	case "請問您要找的分類是?":
		lastBotMessages = ""
		//Get galleries from subreddit tag and send to reply message
		tag := "cats"
		if message.Text != "" {
			tag = message.Text
		}
		imgurRes, err := imgurClient.GetSubredditGalleries(tag)
		if err != nil {
			log.Print(err)
		}
		if !imgurRes.Success {
			if _, err = bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("抱歉，找不到這個分類~:(")).Do(); err != nil {
				log.Print(err)
			}
			break
		}
		accountData, err := json.MarshalIndent(imgurRes.Data, "", "\t")
		if err != nil {
			log.Print(err)
		}
		log.Println(string(accountData))
		// Call Line API to get LIFF App or update it
		imgurGalleriesLink := imgurGalleriesURLBase + tag
		lineLIFFID, err := GetLIFFAppID(imgurGalleriesLink)
		if err != nil {
			log.Print(err)
		}
		if lineLIFFID == "" {
			lineLIFFAddRes, err := AddLIFFApp(linebot.LIFFViewTypeTall, imgurGalleriesLink)
			if err != nil {
				log.Print(err)
			}
			lineLIFFID = lineLIFFAddRes.LIFFID
		} else if _, err := UpdateLIFFApp(
			linebot.LIFFViewTypeTall,
			lineLIFFID,
			imgurGalleriesLink); err != nil {
			log.Print(err)
		}
		flexContent := GetDefaulLIFFBotMessage(lineLIFFID, "預覽分類圖庫")
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("找到這個分類了~"),
			linebot.NewFlexMessage("分類圖庫連結: "+imgurGalleriesLink, flexContent)).Do(); err != nil {
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
	log.Println("message.ID: " + message.ID)
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
		imgurRes, err := imgurClient.AnonymousUploadByImgMessage(imgContent)
		if err != nil {
			log.Print(err)
		}
		log.Printf("image upload success? %v", imgurRes.Success)
		imgUploadData, err := json.MarshalIndent(imgurRes.Data, "", "\t")
		if err != nil {
			log.Print(err)
		}
		log.Printf("image upload result: %s", string(imgUploadData))

		// Call Line API to get LIFF URL
		lineLIFFAddRes, err := AddLIFFApp(linebot.LIFFViewTypeTall, imgurRes.Data.Link)
		if err != nil {
			log.Print(err)
		}
		// Send flex message to show result and furthur actions.
		// But flex message are not allow to be copy by user currently,
		// and if we send another message by push it will be charged
		// when over 500 messages per month under FREE plan.
		// So if you want to save money for your bot,
		// sending multiple messages(up to 5) in 1 reply.
		flexContent := GetImageUploadBotMessage(lineLIFFAddRes.LIFFID, imgurRes.Data.DeleteHash)
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(imgurRes.Data.Link),
			linebot.NewFlexMessage("您的圖片連結: "+imgurRes.Data.Link, flexContent)).Do(); err != nil {
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
