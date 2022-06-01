package services

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/TinyAnnie/engagement-rate-cal/be/instagram"
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
	"os"
	"strings"
)

func CalEngagementRate(username string) float64 {
	profile, err := instagram.GetAccountByUsername(username)
	followersCount := profile.Followers
	posts, err := instagram.GetAccountMedia(username, 12)
	fmt.Println(err)
	var postLikes uint32 = 0
	for _, post := range posts {
		postLikes += post.LikesCount
	}
	return float64(postLikes / followersCount)
}

//"id": user id, "after": end cursor
const nextPageURL string = `https://www.instagram.com/graphql/query/?query_hash=%s&variables=%s`
const nextPagePayload string = `{"id":"%s","first":12,"after":"%s"}`

var requestID string

type pageInfo struct {
	EndCursor string `json:"end_cursor"`
	NextPage  bool   `json:"has_next_page"`
}

type mainPageData struct {
	Rhxgis    string `json:"rhx_gis"`
	EntryData struct {
		ProfilePage []struct {
			Graphql struct {
				User struct {
					Id        string `json:"id"`
					Followers uint32 `json:"followers"`
					Media     struct {
						Edges []struct {
							Node struct {
								LikesCount    uint32 `json:"likes_count"`
								CommentsCount uint32 `json:"comments_count"`
								ImageURL      string `json:"display_url"`
							} `json:"node"`
						} `json:"edges"`
						PageInfo pageInfo `json:"page_info"`
					} `json:"edge_owner_to_timeline_media"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

type nextPageData struct {
	Data struct {
		User struct {
			Container struct {
				PageInfo pageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						CommentsCount uint32 `json:"comments_count"`
						LikesCount    uint32 `json:"likes_count"`
						ImageURL      string `json:"display_url"`
					}
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		}
	} `json:"data"`
}

func CalEngagementRate2(username string) float64 {
	var followerCount uint32 = 0
	var interactCount uint32 = 0
	var postsLimit = 12

	if len(username) == 0 {
		log.Println("username is missing")
		os.Exit(1)
	}

	var actualUserId string
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referrer", "https://www.instagram.com/"+username)
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		d := c.Clone()
		d.OnResponse(func(r *colly.Response) {
			idStart := bytes.Index(r.Body, []byte(`:n},queryId:"`))
			requestID = string(r.Body[idStart+13 : idStart+45])
		})
		requestIDURL := e.Request.AbsoluteURL(e.ChildAttr(`link[as="script"]`, "href"))
		d.Visit(requestIDURL)

		dat := e.ChildText("body > script:first-of-type")
		jsonData := dat[strings.Index(dat, "{") : len(dat)-1]
		data := &mainPageData{}
		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			log.Fatal(err)
		}
		page := data.EntryData.ProfilePage[0]
		followerCount = page.Graphql.User.Followers
		actualUserId = page.Graphql.User.Id
		for _, obj := range page.Graphql.User.Media.Edges {
			if postsLimit == 0 {
				break
			}
			interactCount += obj.Node.LikesCount + obj.Node.CommentsCount
			postsLimit--
			c.Visit(obj.Node.ImageURL)
		}
		nextPageVars := fmt.Sprintf(nextPagePayload, actualUserId, page.Graphql.User.Media.PageInfo.EndCursor)
		e.Request.Ctx.Put("variables", nextPageVars)
		if page.Graphql.User.Media.PageInfo.NextPage {
			u := fmt.Sprintf(
				nextPageURL,
				requestID,
				url.QueryEscape(nextPageVars),
			)
			log.Println("Next page found", u)
			e.Request.Ctx.Put("gis", data.Rhxgis)
			e.Request.Visit(u)
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			return
		}

		if strings.Index(r.Headers.Get("Content-Type"), "json") == -1 {
			return
		}

		data := &nextPageData{}
		err := json.Unmarshal(r.Body, data)
		if err != nil {
			log.Fatal(err)
		}

		for _, obj := range data.Data.User.Container.Edges {
			if postsLimit == 0 {
				break
			}
			interactCount += obj.Node.LikesCount + obj.Node.CommentsCount
			postsLimit--
			c.Visit(obj.Node.ImageURL)
		}
		if data.Data.User.Container.PageInfo.NextPage {
			nextPageVars := fmt.Sprintf(nextPagePayload, actualUserId, data.Data.User.Container.PageInfo.EndCursor)
			r.Request.Ctx.Put("variables", nextPageVars)
			u := fmt.Sprintf(
				nextPageURL,
				requestID,
				url.QueryEscape(nextPageVars),
			)
			log.Println("Next page found", u)
			r.Request.Visit(u)
		}
	})

	c.Visit("https://instagram.com/" + username)
	return float64(interactCount / followerCount)
}