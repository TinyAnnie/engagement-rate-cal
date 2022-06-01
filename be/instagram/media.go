package instagram

import (
	"encoding/json"
	"strconv"
)

// TypeImage is a string that define image type for media.
const TypeImage = "image"

// TypeVideo is a string that define video type for media.
const TypeVideo = "video"

// TypeCarousel is a string that define carousel (collection of media) type for media.
const TypeCarousel = "carousel"

const (
	graphVideo   = "GraphVideo"
	graphSidecar = "GraphSidecar"

	video    = "video"
	carousel = "carousel"
)

// A Media describes an Instagram media info.
type Media struct {
	Caption       string
	Code          string
	CommentsCount uint32
	Date          uint64
	ID            string
	AD            bool
	LikesCount    uint32
	Type          string
	MediaURL      string
	Owner         Account
	MediaList     []mediaItem
}

type mediaItem struct {
	Type string
	URL  string
	Code string
}

func getFromMediaPage(data []byte) (Media, error) {
	var mediaJSON struct {
		Graphql struct {
			ShortcodeMedia struct {
				Typename           string `json:"__typename"`
				ID                 string `json:"id"`
				Shortcode          string `json:"shortcode"`
				DisplayURL         string `json:"display_url"`
				VideoURL           string `json:"video_url"`
				IsVideo            bool   `json:"is_video"`
				EdgeMediaToCaption struct {
					Edges []struct {
						Node struct {
							Text string `json:"text"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"edge_media_to_caption"`
				EdgeMediaToComment struct {
					Count int `json:"count"`
				} `json:"edge_media_to_comment"`
				TakenAtTimestamp     int `json:"taken_at_timestamp"`
				EdgeMediaPreviewLike struct {
					Count int `json:"count"`
				} `json:"edge_media_preview_like"`
				Owner struct {
					ID            string `json:"id"`
					ProfilePicURL string `json:"profile_pic_url"`
					Username      string `json:"username"`
					FullName      string `json:"full_name"`
					IsPrivate     bool   `json:"is_private"`
				} `json:"owner"`
				IsAd                  bool `json:"is_ad"`
				EdgeSidecarToChildren struct {
					Edges []struct {
						Node struct {
							Typename   string `json:"__typename"`
							ID         string `json:"id"`
							Shortcode  string `json:"shortcode"`
							DisplayURL string `json:"display_url"`
							VideoURL   string `json:"video_url"`
							IsVideo    bool   `json:"is_video"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"edge_sidecar_to_children"`
			} `json:"shortcode_media"`
		} `json:"graphql"`
	}

	err := json.Unmarshal(data, &mediaJSON)
	if err != nil {
		return Media{}, err
	}

	media := Media{}
	media.Code = mediaJSON.Graphql.ShortcodeMedia.Shortcode
	media.ID = mediaJSON.Graphql.ShortcodeMedia.ID
	media.AD = mediaJSON.Graphql.ShortcodeMedia.IsAd
	media.Date = uint64(mediaJSON.Graphql.ShortcodeMedia.TakenAtTimestamp)
	media.CommentsCount = uint32(mediaJSON.Graphql.ShortcodeMedia.EdgeMediaToComment.Count)
	media.LikesCount = uint32(mediaJSON.Graphql.ShortcodeMedia.EdgeMediaPreviewLike.Count)

	if len(mediaJSON.Graphql.ShortcodeMedia.EdgeMediaToCaption.Edges) > 0 {
		media.Caption = mediaJSON.Graphql.ShortcodeMedia.EdgeMediaToCaption.Edges[0].Node.Text
	}

	var mediaType = mediaJSON.Graphql.ShortcodeMedia.Typename
	if mediaType == graphSidecar {
		for _, itemJSON := range mediaJSON.Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges {
			var item mediaItem
			item.Code = itemJSON.Node.Shortcode
			if itemJSON.Node.IsVideo {
				item.URL = itemJSON.Node.VideoURL
				item.Type = TypeVideo
			} else {
				item.URL = itemJSON.Node.DisplayURL
				item.Type = TypeImage
			}
			media.MediaList = append(media.MediaList, item)
		}
		media.Type = TypeCarousel
	} else {
		if mediaType == graphVideo {
			media.Type = TypeVideo
			media.MediaURL = mediaJSON.Graphql.ShortcodeMedia.VideoURL
		} else {
			media.Type = TypeImage
			media.MediaURL = mediaJSON.Graphql.ShortcodeMedia.DisplayURL
		}
		var item mediaItem
		item.Code = media.Code
		item.Type = media.Type
		item.URL = media.MediaURL
		media.MediaList = append(media.MediaList, item)
	}

	media.Owner.ID = mediaJSON.Graphql.ShortcodeMedia.Owner.ID
	media.Owner.ProfilePicURL = mediaJSON.Graphql.ShortcodeMedia.Owner.ProfilePicURL
	media.Owner.Username = mediaJSON.Graphql.ShortcodeMedia.Owner.Username
	media.Owner.FullName = mediaJSON.Graphql.ShortcodeMedia.Owner.FullName
	media.Owner.Private = mediaJSON.Graphql.ShortcodeMedia.Owner.IsPrivate

	return media, nil
}

func getFromAccountMediaList(data []byte) (Media, error) {
	var mediaJSON struct {
		ID   string `json:"id"`
		Code string `json:"code"`
		User struct {
			ID             string `json:"id"`
			FullName       string `json:"full_name"`
			ProfilePicture string `json:"profile_picture"`
			Username       string `json:"username"`
		} `json:"user"`
		Images struct {
			StandardResolution struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"standard_resolution"`
		} `json:"images"`
		CreatedTime string `json:"created_time"`
		Caption     struct {
			Text string `json:"text"`
		} `json:"caption"`
		Likes struct {
			Count float64 `json:"count"`
		} `json:"likes"`
		Comments struct {
			Count float64 `json:"count"`
		} `json:"comments"`
		Type   string `json:"type"`
		Videos struct {
			StandardResolution struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"standard_resolution"`
		} `json:"videos"`
		CarouselMedia []struct {
			Images struct {
				StandardResolution struct {
					URL string `json:"url"`
				} `json:"standard_resolution"`
			} `json:"images"`
			Videos struct {
				StandardResolution struct {
					URL string `json:"url"`
				} `json:"standard_resolution"`
			} `json:"videos"`
			UsersInPhoto []interface{} `json:"users_in_photo"`
			Type         string        `json:"type"`
		} `json:"carousel_media"`
	}

	err := json.Unmarshal(data, &mediaJSON)
	if err != nil {
		return Media{}, err
	}

	media := Media{}
	media.Code = mediaJSON.Code
	media.ID = mediaJSON.ID
	media.Caption = mediaJSON.Caption.Text
	media.LikesCount = uint32(mediaJSON.Likes.Count)
	media.CommentsCount = uint32(mediaJSON.Comments.Count)

	date, err := strconv.ParseUint(mediaJSON.CreatedTime, 10, 64)
	if err == nil {
		media.Date = date
	}

	if mediaJSON.Type == carousel {
		media.Type = TypeCarousel
		for _, itemJSOM := range mediaJSON.CarouselMedia {
			var item mediaItem
			item.Type = itemJSOM.Type
			if item.Type == video {
				item.URL = itemJSOM.Videos.StandardResolution.URL
			} else {
				item.URL = itemJSOM.Images.StandardResolution.URL
			}
			media.MediaList = append(media.MediaList, item)
		}
	} else {
		if mediaJSON.Type == video {
			media.MediaURL = mediaJSON.Videos.StandardResolution.URL
			media.Type = TypeVideo
		} else {
			media.MediaURL = mediaJSON.Images.StandardResolution.URL
			media.Type = TypeImage
		}
		var item mediaItem
		item.Type = media.Type
		item.URL = media.MediaURL
		item.Code = media.Code
		media.MediaList = append(media.MediaList, item)
	}

	media.Owner.Username = mediaJSON.User.Username
	media.Owner.FullName = mediaJSON.User.FullName
	media.Owner.ID = mediaJSON.User.ID
	media.Owner.ProfilePicURL = mediaJSON.User.ProfilePicture

	return media, nil
}