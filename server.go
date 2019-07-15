package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BensonLiao/imgur-api-go-v3/imgurclient"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Note. Any secret or token should remove in production and use environment virable
var err error
var bot *linebot.Client
var channelSecret = ""
var channelToken = ""
var replyMessage = ""
var lastBotMessages = ""
var lineLIFFURLBase = "line://app/"
var lineLIFFURL = lineLIFFURLBase + "1646605627-0bPlzl73"
var imgurClient *imgurclient.Client
var imgurClientID = ""
var imgurClientSecret = ""
var imgurUserURL = "https://imgur.com/user/"

// SplitHTTPReqParams func, split http request params and return a key-value map
func SplitHTTPReqParams(params string) map[string]string {
	keyMaps := make(map[string]string)
	for _, keyPair := range strings.Split(params, "&") {
		kpList := strings.Split(keyPair, "=")
		keyMaps[kpList[0]] = kpList[1]
	}
	return keyMaps
}

//Get account info from imgur.com and return jsonfy string
func getImgurAccInfo(accID string) (info string, err error) {
	req, err := http.NewRequest("GET", "https://api.imgur.com/3/account/"+accID, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Client-ID "+imgurClientID)
	//Make sure no hanging request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "抱歉!找不到這位使用者", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Print(string(body))
	return string(body), nil
}

//Upload an image to imgur.com and return upload result
func uploadToImgur(img []uint8) (result string, err error) {
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Client-ID "+imgurClientID)
	//Make sure no hanging request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Print(string(body))
	return string(body), nil
}

// callbackFunc func, http handler for /callback
func callbackFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		replyText := ""
		bot, err = linebot.New(
			channelSecret,
			channelToken,
		)
		imgurClient = imgurclient.New(
			imgurClientID,
			imgurClientSecret,
			"",
			"",
		)
		log.Println(c.Request)
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.String(400, "linebot error: ErrInvalidSignature")
			} else {
				c.String(500, "unknown linebot error")
			}
			return
		}
		for _, event := range events {
			log.Print(event.Source.UserID)
			log.Print(event.Message)
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					log.Print(message.Text)
					switch message.Text {
					case "give me brown":
						fallthrough
					case "給我熊大":
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewStickerMessage("11537", "52002739")).Do(); err != nil {
							log.Print(err)
						}
					case "search imgur account":
						fallthrough
					case "搜尋imgur帳號":
						replyText := "請問您要找的帳號為?"
						lastBotMessages = replyText
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage(replyText)).Do(); err != nil {
							log.Print(err)
						}
					case "upload image to imgur":
						fallthrough
					case "上傳圖片到imgur":
						replyText = "請問您要上傳哪張圖片?"
						lastBotMessages = replyText
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage(replyText)).Do(); err != nil {
							log.Print(err)
						}
					case "刪除上傳的圖片":
						break
					default:
						if lastBotMessages == "" {
							contents := GetDefaultBotMessage(message.Text)
							if _, err = bot.ReplyMessage(
								event.ReplyToken,
								linebot.NewFlexMessage("Hello, World!", contents)).Do(); err != nil {
								log.Print(err)
							}
						} else {
							ActionsOnBotMessage(message.Text, event.ReplyToken, event.Source.UserID)
						}
					}
				case *linebot.ImageMessage:
					if lastBotMessages == "" {
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
					} else {
						ActionsOnBotMessage(message.ID, event.ReplyToken, event.Source.UserID)
					}
				case *linebot.FileMessage:
					log.Print("File Name: " + message.FileName)
				}
			} else if event.Type == linebot.EventTypePostback {
				log.Printf("event.Postback.Data: %s\n", event.Postback.Data)
				postbackParams := SplitHTTPReqParams(event.Postback.Data)
				// Check if deletehash exist and do iamge deletion from imgur
				value, ok := postbackParams["deletehash"]
				if ok {
					imgDeleteRes, err := imgurClient.DeleteAnonymousUploadedImg(value)
					if err != nil {
						log.Print(err)
					}
					if !imgDeleteRes.Success {
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("很抱歉，刪除發生異常:(")).Do(); err != nil {
							log.Print(err)
						}
						break
					} else {
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("圖片刪除成功!")).Do(); err != nil {
							log.Print(err)
						}
					}
				}
				// Delete LIFF App URL from LINE
				value, ok = postbackParams["liffid"]
				if ok {
					liffDeleteRes, err := bot.DeleteLIFF(value).Do()
					if err != nil {
						log.Print(err)
					}
					if liffDeleteRes != nil {
						log.Println("Delete LIFF App URL success!")
					}
				}
			}
		}
	}
}

func main() {
	channelSecret = os.Getenv("LINEBOT_CHANNEL_SECRET")
	if channelSecret == "" {
		log.Fatal("$LINEBOT_CHANNEL_SECRET must be set to enable linebot")
	}
	channelToken = os.Getenv("LINEBOT_CHANNEL_TOKEN")
	imgurClientID = os.Getenv("IMGUR_CLIENT_ID")
	imgurClientSecret = os.Getenv("IMGUR_CLIENT_SECRET")
	if imgurClientSecret == "" {
		log.Fatal("$IMGUR_CLIENT_SECRET must be set to use imgur api")
	}

	bot, err = linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/callback", callbackFunc())

	router.Run(":" + port)
}
