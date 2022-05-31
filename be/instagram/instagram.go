package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetMediaByURL try to find media by url.
// URL should be like https://www.instagram.com/p/12376OtT5o/
func GetMediaByURL(url string) (Media, error) {
	code := strings.Split(url, "/")[4]
	return GetMediaByCode(code)
}

// GetMediaByCode try to find media by code.
// Code can be find in URL to media, after p/.
// If URL to media is https://www.instagram.com/p/12376OtT5o/,
// then code of the media is 12376OtT5o.
func GetMediaByCode(code string) (Media, error) {
	url := fmt.Sprintf(mediaInfoURL, code)
	data, err := getDataFromURL(url)
	if err != nil {
		return Media{}, err
	}
	media, err := getFromMediaPage(data)
	if err != nil {
		return Media{}, err
	}
	return media, nil
}

// GetTagMedia try to get slice of last tag's media.
// The limit set how much media you need.
func GetTagMedia(tag string, quantity uint16) ([]Media, error) {
	var count uint16
	maxID := ""
	hasNext := true
	medias := []Media{}
	for hasNext && count < quantity {
		url := fmt.Sprintf(tagURL, tag, maxID)
		jsonBody, err := getJSONFromURL(url)
		if err != nil {
			return nil, err
		}
		jsonBody, _ = jsonBody["tag"].(map[string]interface{})
		jsonBody, _ = jsonBody["media"].(map[string]interface{})

		nodes, _ := jsonBody["nodes"].([]interface{})
		for _, node := range nodes {
			if count >= quantity {
				return medias, nil
			}
			count++
			nodeData, err := json.Marshal(node)
			if err == nil {
				media, err := getFromSearchMediaList(nodeData)
				if err == nil {
					medias = append(medias, media)
				}
			}
		}

		jsonBody, _ = jsonBody["page_info"].(map[string]interface{})
		hasNext, _ = jsonBody["has_next_page"].(bool)
		maxID, _ = jsonBody["end_cursor"].(string)
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