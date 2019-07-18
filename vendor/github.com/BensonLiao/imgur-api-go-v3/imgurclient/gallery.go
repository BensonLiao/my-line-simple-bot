package imgurclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// GalleryBase const, imgur's gallery related api endpoint
const GalleryBase = APIBase + "/gallery"

// GetSubredditGalleries func, send tag and return Gallery's subreddit
func (cl *Client) GetSubredditGalleries(tag string) (DataListResponse, error) {
	dlr := DataListResponse{}
	request, _ := cl.PrepareAuthRequest("GET", GalleryBase+"/r/"+tag)
	response, err := cl.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("Subreddit galleries body: %s", string(body))
	err = json.Unmarshal(body, &dlr)
	if err != nil {
		return DataListResponse{}, err
	}
	return dlr, err
}
