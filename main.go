package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type VideoOccurence struct {
	Timestamp time.Time
	YTID      string
}

func main() {
	http.HandleFunc("/tweets", handler)
	http.ListenAndServe(":8080", nil)
}

func findVideos() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("twitwyrm").C("VideoOccurences")
	session.SetMode(mgo.Monotonic, true)

	anaconda.SetConsumerKey("tDAaN4yvVg5wFVlfkxeEUOzHm")
	anaconda.SetConsumerSecret("fRtI78PU9T67xNM5dylEO8146qEjtF5qVEHHZIBLuoJ6BHHjm0")
	api := anaconda.NewTwitterApi("332989840-mBSNofbvaHW2eyENXlAmI52GLpyGQwVvRav5jdsl",
		"y6N3Ty0d99JDARs4FPvHmYjtocUrpQlRcN2uQ9ljTWeBu")

	values := url.Values{}
	// params := r.URL.Query()
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

				if videoId, ok := parsed.Query()["v"]; ok {
					fmt.Println(videoId[0])

					// w.Write([]byte(videoId[0]))

					vo := &VideoOccurence{
						YTID:      videoId[0],
						Timestamp: time.Now(),
					}

					err = c.Insert(vo)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// result := VideoOccurence{}
					// err = c.Find(bson.M{"id": videoId[0]}).One(&result)
					// if err != nil {
					// 	fmt.Println(err)
					// 	fmt.Println("find")
					// 	continue
					// }
					//
					// fmt.Println("VO:", result)
					// return

					// err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
					// 	&Person{"Cla", "+55 53 8402 8510"})
				}
			}
		} else {
			fmt.Println("not a tweet")
			fmt.Println(item)
		}

	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	go findVideos()

	w.Write([]byte(fmt.Sprintf("Started with params: %s", r.URL.Query())))
}
