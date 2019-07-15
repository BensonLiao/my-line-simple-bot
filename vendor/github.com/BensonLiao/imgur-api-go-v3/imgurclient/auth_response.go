package imgurclient

// AuthResponse type
type AuthResponse struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn       uint64 `json:"expires_in"`
	Scope           string `json:"scope"`
	RefreshToken    string `json:"refresh_token"`
	AccountID       int64  `json:"account_id"`
	AccountUsername string `json:"account_username"`
}
