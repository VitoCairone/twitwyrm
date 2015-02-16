package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	http.HandleFunc("/tweets", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	anaconda.SetConsumerKey("tDAaN4yvVg5wFVlfkxeEUOzHm")
	anaconda.SetConsumerSecret("fRtI78PU9T67xNM5dylEO8146qEjtF5qVEHHZIBLuoJ6BHHjm0")
	api := anaconda.NewTwitterApi("332989840-mBSNofbvaHW2eyENXlAmI52GLpyGQwVvRav5jdsl",
		"y6N3Ty0d99JDARs4FPvHmYjtocUrpQlRcN2uQ9ljTWeBu")

	values := url.Values{}
	params := r.URL.Query()
	values.Add("track", "youtube com,youtu be")
	stream := api.PublicStreamFilter(params)

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

					w.Write([]byte(videoId[0]))
				}
			}
		} else {
			fmt.Println("not a tweet")
			fmt.Println(item)
		}

	}
}
