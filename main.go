package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2"
)

type VideoOccurence struct {
	Timestamp time.Time
	YTID      string
}

func main() {
	voChan := findVideos()
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("twitwyrm").C("VideoOccurences")
	session.SetMode(mgo.Monotonic, true)

	for vo := range voChan {
		err = c.Insert(vo)
		if err != nil {
			fmt.Println(err)
		}
	}

	http.HandleFunc("/tweets", handler)
	http.ListenAndServe(":8080", nil)
}

func findVideos() chan VideoOccurence {
	voChan := make(chan VideoOccurence)

	go func() {
		anaconda.SetConsumerKey("tDAaN4yvVg5wFVlfkxeEUOzHm")
		anaconda.SetConsumerSecret("fRtI78PU9T67xNM5dylEO8146qEjtF5qVEHHZIBLuoJ6BHHjm0")
		api := anaconda.NewTwitterApi("332989840-mBSNofbvaHW2eyENXlAmI52GLpyGQwVvRav5jdsl",
			"y6N3Ty0d99JDARs4FPvHmYjtocUrpQlRcN2uQ9ljTWeBu")

		values := url.Values{}
		values.Add("track", "youtube com,youtu be")
		stream := api.PublicStreamFilter(values)

		for item := range stream.C {
			tweet, ok := item.(anaconda.Tweet)
			if ok {
				urls := tweet.Entities.Urls
				for _, u := range urls {
					parsed, err := url.Parse(u.Expanded_url)
					if err != nil {
						fmt.Println(err)
						continue
					}
					vo := VideoOccurence{}
					switch parsed.Host {
					case "youtu.be":
						videoId := parsed.Path

						if len(videoId) == 11 { // All youtube IDs are length 11
							vo.YTID = videoId
							vo.Timestamp = time.Now()
							voChan <- vo
						}
					case "youtube.com":

						if videoIdSlice, ok := parsed.Query()["v"]; ok {
							videoId := videoIdSlice[0]
							if len(videoId) == 11 { // All youtube IDs are length 11
								vo.YTID = videoId
								vo.Timestamp = time.Now()
								voChan <- vo
							}
						}

					}
				}
			}
		}
	}()
	return voChan
}

func writeResults() {

}

func handler(w http.ResponseWriter, r *http.Request) {

	go findVideos()

	w.Write([]byte(fmt.Sprintf("Started with params: %s", r.URL.Query())))
}
