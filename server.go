package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

// Sticker struct
type Sticker struct {
	PackageID string
	StickerID string
}

var brownSaluteSticker = Sticker{
	PackageID: "11537",
	StickerID: "52002739",
}

var replyTextMessage = ""
var lastBotMessages = ""
var lineLIFFURLBase = "line://app/"
var lineLIFFURL = lineLIFFURLBase + "1646605627-0bPlzl73"
var imgurClient *imgurclient.Client
var imgurClientID = ""
var imgurClientSecret = ""
var imgurUserURL = "https://imgur.com/user/"

// GetRandomSticker func
func GetRandomSticker() *Sticker {
	randomSticker := new(Sticker)
	rand.Seed(time.Now().UnixNano())
	randomPackageID := strconv.Itoa(rand.Intn(2) + 11537)
	randomSticker.PackageID = randomPackageID
	switch randomPackageID {
	case "11537":
		randomStickerID := rand.Intn(39) + 52002734
		if randomStickerID > 52002770 {
			randomStickerID += 6
		}
		randomSticker.StickerID = strconv.Itoa(randomStickerID)
	case "11538":
		randomSticker.StickerID = strconv.Itoa(rand.Intn(39) + 51626494)
	case "11539":
		randomSticker.StickerID = strconv.Itoa(rand.Intn(39) + 52114110)
	}
	return randomSticker
}

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

// ReplyTextMessage func
func ReplyTextMessage(replyToken, text string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyStickerMessage func
func ReplyStickerMessage(replyToken, packgeID, stickerID string) error {
	log.Printf("replyToken: %s, packgeID: %s, stickerID: %s", replyToken, packgeID, stickerID)
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewStickerMessage(packgeID, stickerID),
	).Do(); err != nil {
		return fmt.Errorf(
			"[ReplyStickerMessage] error in calling bot.ReplyMessage: %v",
			err)
	}
	return nil
}

// HandleImage func
func HandleImage(message *linebot.ImageMessage, replyToken string, source *linebot.EventSource) error {
	if lastBotMessages == "" {
		res, err := bot.GetMessageContent(message.ID).Do()
		if err != nil {
			return err
		}
		body := res.Content
		defer body.Close()
		bodyGot, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		log.Printf("bodyGot's type: %T\n", bodyGot)
	} else {
		ActionsOnBotMessage(message.ID, replyToken, source.UserID)
	}
	return nil
}

// HandleFile func
func HandleFile(message *linebot.FileMessage, replyToken string) error {
	return ReplyTextMessage(replyToken,
		fmt.Sprintf(
			"File `%s` (%d bytes) received.",
			message.FileName,
			message.FileSize))
}

// HandleLocation func
func HandleLocation(message *linebot.LocationMessage, replyToken string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewLocationMessage(message.Title, message.Address, message.Latitude, message.Longitude),
	).Do(); err != nil {
		return err
	}
	return nil
}

// HandleSticker func
func HandleSticker(message *linebot.StickerMessage, replyToken string) error {
	// There's a linebot api 400 error if reply sticker that same as user sent
	// are not available
	mySticker := GetRandomSticker()
	if err := ReplyStickerMessage(
		replyToken,
		mySticker.PackageID,
		mySticker.StickerID); err != nil {
		return fmt.Errorf(
			"[HandleSticker] error in calling ReplyStickerMessage: %v",
			err)
	}
	return nil
}

// HandleText func
func HandleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "give me brown":
		fallthrough
	case "給我熊大":
		if err := ReplyStickerMessage(
			replyToken,
			brownSaluteSticker.PackageID,
			brownSaluteSticker.StickerID); err != nil {
			return err
		}
	case "search imgur account":
		fallthrough
	case "搜尋imgur帳號":
		replyTextMessage = "請問您要找的帳號為?"
		lastBotMessages = replyTextMessage
		if err := ReplyTextMessage(replyToken, replyTextMessage); err != nil {
			return err
		}
	case "upload image to imgur":
		fallthrough
	case "上傳圖片到imgur":
		replyTextMessage = "請問您要上傳哪張圖片?"
		lastBotMessages = replyTextMessage
		if err := ReplyTextMessage(replyToken, replyTextMessage); err != nil {
			log.Print(err)
		}
	case "刪除上傳的圖片":
		break
	default:
		if lastBotMessages == "" {
			contents := GetDefaultBotMessage(message.Text)
			if _, err = bot.ReplyMessage(
				replyToken,
				linebot.NewFlexMessage("Hello, World!", contents)).Do(); err != nil {
				return err
			}
		} else {
			ActionsOnBotMessage(message.Text, replyToken, source.UserID)
		}
	}
	return nil
}

// Callback func, http handler for /callback
func Callback(c *gin.Context) {
	log.Println("Callback function called")
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
	log.Printf("%+v\n", c.Request)
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
		log.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := HandleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				if err := HandleSticker(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.ImageMessage:
				if err := HandleImage(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			case *linebot.FileMessage:
				if err := HandleFile(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				if err := HandleLocation(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			}
		case linebot.EventTypePostback:
			log.Printf("event.Postback.Data: %s\n", event.Postback.Data)
			postbackParams := SplitHTTPReqParams(event.Postback.Data)
			// Check if deletehash exist and do image deletion from imgur
			value, ok := postbackParams["deletehash"]
			if ok {
				imgDeleteRes, err := imgurClient.DeleteAnonymousUploadedImg(value)
				if err != nil {
					log.Print(err)
				}
				replyTextMessage = "圖片刪除成功!"
				if !imgDeleteRes.Success {
					replyTextMessage = "很抱歉，刪除發生異常:("
				}
				if err := ReplyTextMessage(event.ReplyToken, replyTextMessage); err != nil {
					log.Print(err)
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
		default:
			log.Printf("Unknown event: %v", event)
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
	router.POST("/callback", Callback)

	router.Run(":" + port)
}
