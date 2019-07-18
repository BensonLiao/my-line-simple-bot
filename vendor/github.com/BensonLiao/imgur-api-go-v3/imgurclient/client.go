package imgurclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// Client for imgur api
type Client struct {
	ClientID        string
	ClientSecret    string
	AccessToken     string
	ExpiresIn       int64
	TokenType       string
	RefreshToken    string
	AccountUsername string
	AccountID       int64
	http.Client
}

// APIBase const, base of imgur api endpoint
const APIBase = "https://api.imgur.com/3"

// Auth const, imgur api endpoint to get authorized token
const Auth = "https://api.imgur.com/oauth2/authorize"

// Token const, imgur api endpoint to get refreshed token
const Token = "https://api.imgur.com/oauth2/token"

// New func, initialize imgur client
func New(key, secret, accessToken, refreshToken string) *Client {
	return &Client{
		ClientID:     key,
		ClientSecret: secret,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// GetAuthorizationURL func
func (cl *Client) GetAuthorizationURL(authType string) string {
	return fmt.Sprintf(
		"%s?client_id=%s&response_type=%s",
		Auth,
		cl.ClientID,
		authType)
}

// Authorize func
func (cl *Client) Authorize(pin, authType string) (AuthResponse, error) {
	ir := AuthResponse{}
	v := url.Values{}
	v.Set("client_id", cl.ClientID)
	v.Set("client_secret", cl.ClientSecret)
	v.Set("grant_type", authType)
	v.Set("pin", pin)
	response, err := cl.PostForm(Token, v)
	if response.StatusCode == 200 {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &ir)
	} else {
		err = fmt.Errorf(fmt.Sprintf(
			"ImgurClient#Authorize: Status code: %d, authtype: %s",
			response.StatusCode, authType))
	}
	return ir, err
}

// Refresh func
func (cl *Client) Refresh() error {
	ir := AuthResponse{}
	vals := url.Values{}
	vals.Add("refresh_token", cl.RefreshToken)
	vals.Add("client_id", cl.ClientID)
	vals.Add("client_secret", cl.ClientSecret)
	vals.Add("grant_type", "refresh_token")
	response, err := cl.PostForm(Token, vals)
	if response.StatusCode == 200 {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &ir)
		cl.AccessToken = ir.AccessToken
		fmt.Printf("%v\n", ir)
		fmt.Println(cl)
	} else {
		err = fmt.Errorf(fmt.Sprintf(
			"ImgurClient#Authorize: Status code: %d, authtype: refresh_token",
			response.StatusCode))
	}
	return err
}

// PrepareAuthRequest func, return a http request with Authorization header
func (cl *Client) PrepareAuthRequest(method, url string) (*http.Request, error) {
	fmt.Printf("req url: %s\n", url)
	req, err := http.NewRequest(method, url, nil)
	// fmt.Printf("client access token: %s\n", cl.AccessToken)
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cl.AccessToken))
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", cl.ClientID))
	return req, err
}

// NewFileUploadRequest func,
// Creates a new file upload http request with optional extra params
func (cl *Client) NewFileUploadRequest(
	uri string,
	params map[string]string,
	fileParam,
	path string,
) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParam, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}

// NewImgContentUploadRequest func,
// Creates a new image content upload http request with optional extra params
func (cl *Client) NewImgContentUploadRequest(
	uri string,
	params map[string]string,
	imgContent []byte,
	fieldName string,
) (*http.Request, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormField(fieldName)
	if err != nil {
		return nil, err
	}
	part.Write(imgContent)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}
