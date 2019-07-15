package imgurclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Account type
type Account struct {
	ID             int     `json:"id"`
	URL            string  `json:"url"`
	Bio            string  `json:"bio"`
	Reputation     float64 `json:"reputation"`
	CreatedSeconds int64   `json:"created"`
	ProExpiration  bool    `json:"pro_expiration"`
}

// AccountResponse type
type AccountResponse struct {
	Data    Account `json:"data"`
	Status  int     `json:"status"`
	Success bool    `json:"success"`
}

// GetAccount func, send username and return Account
func (cl *Client) GetAccount(username string) (AccountResponse, error) {
	ar := AccountResponse{}
	if username == "" {
		username = "me"
		// "me" Only works when cl.ClientID is valid and  imgur will search
		// for cl.ClientID's registered account
	}
	request, _ := cl.prepareRequest("GET", "account/"+username)
	response, err := cl.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("account body: %s", string(body))
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return AccountResponse{}, err
	}
	return ar, err
}
