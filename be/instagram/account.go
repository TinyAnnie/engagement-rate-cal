package instagram

import "encoding/json"

// An Account describes an Instagram account info.
type Account struct {
	Biography       string
	ExternalURL     string
	Followers       uint32
	Follows         uint32
	FullName        string
	ID              string
	Private         bool
	Verified        bool
	MediaCount      uint32
	ProfilePicURL   string
	ProfilePicURLhd string
	Username        string
}

func getFromAccountPage(data []byte) (Account, error) {
	var accountJSON struct {
		User struct {
			Biography   string `json:"biography"`
			ExternalURL string `json:"external_url"`
			FollowedBy  struct {
				Count int `json:"count"`
			} `json:"followed_by"`
			Follows struct {
				Count int `json:"count"`
			} `json:"follows"`
			FullName        string `json:"full_name"`
			ID              string `json:"id"`
			IsPrivate       bool   `json:"is_private"`
			IsVerified      bool   `json:"is_verified"`
			ProfilePicURL   string `json:"profile_pic_url"`
			ProfilePicURLHd string `json:"profile_pic_url_hd"`
			Username        string `json:"username"`
			Media           struct {
				Count int `json:"count"`
			} `json:"media"`
		} `json:"user"`
	}

	err := json.Unmarshal(data, &accountJSON)
	if err != nil {
		return Account{}, err
	}

	account := Account{}
	account.Biography = accountJSON.User.Biography
	account.ExternalURL = accountJSON.User.ExternalURL
	account.FullName = accountJSON.User.FullName
	account.ID = accountJSON.User.ID
	account.Private = accountJSON.User.IsPrivate
	account.Verified = accountJSON.User.IsVerified
	account.ProfilePicURL = accountJSON.User.ProfilePicURL
	account.ProfilePicURLhd = accountJSON.User.ProfilePicURLHd
	account.Username = accountJSON.User.Username
	account.Followers = uint32(accountJSON.User.FollowedBy.Count)
	account.Follows = uint32(accountJSON.User.Follows.Count)
	account.MediaCount = uint32(accountJSON.User.Media.Count)

	return account, nil
}