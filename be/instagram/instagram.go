package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetAccountByUsername try to find account by username.
func GetAccountByUsername(username string) (Account, error) {
	url := fmt.Sprintf(accountInfoURL, username)
	data, err := getDataFromURL(url)
	if err != nil {
		fmt.Println("err", err)
		return Account{}, err
	}
	account, err := getFromAccountPage(data)
	if err != nil {
		return account, err
	}
	return account, nil
}

// GetAccountMedia try to get slice of user's media.
// Limit set how much media you need.
func GetAccountMedia(username string, limit uint16) ([]Media, error) {
	var count uint16
	maxID := ""
	available := true
	medias := []Media{}
	for available && count < limit {
		url := fmt.Sprintf(accountMediaURL, username, maxID)
		jsonBody, err := getJSONFromURL(url)
		if err != nil {
			return nil, err
		}
		available, _ = jsonBody["more_available"].(bool)

		items, _ := jsonBody["items"].([]interface{})
		for _, item := range items {
			if count >= limit {
				return medias, nil
			}
			count++
			itemData, err := json.Marshal(item)
			if err == nil {
				media, err := getFromAccountMediaList(itemData)
				if err == nil {
					medias = append(medias, media)
					maxID = media.ID
				}
			}
		}
	}
	return medias, nil
}

func getJSONFromURL(url string) (map[string]interface{}, error) {
	data, err := getDataFromURL(url)
	if err != nil {
		return nil, err
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(data, &jsonBody)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
}

func getDataFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("statusCode != 200")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}