package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type VideoOccurence struct {
	Timestamp time.Time
	YTID      string
	Language  string
}

type topVideoResult struct {
	Id    string   `bson:"_id"`
	Count int      `bson:"count"`
	Lang  []string `bson:"langs"`
}

func main() {

	voChan := findVideos()
	go func() {
		session, err := mgo.Dial("localhost:27017")
		defer session.Close()

		c := session.DB("twitwyrm").C("VideoOccurences")
		session.SetMode(mgo.Monotonic, true)

		if err != nil {
			panic(err)
		}

		for vo := range voChan {
			err = c.Insert(vo)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	http.HandleFunc("/toptweets", TopVideos)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
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
					vo := VideoOccurence{
						Timestamp: time.Now(),
						Language:  tweet.Lang,
					}
					switch parsed.Host {
					case "youtu.be":
						videoId := parsed.Path
						if len(videoId) == 12 {
							vo.YTID = videoId[1:]
						}
					case "youtube.com":
						// log.Println(parsed)
						if videoIdSlice, ok := parsed.Query()["v"]; ok {
							videoId := videoIdSlice[0]
							if len(videoId) == 11 { // All youtube IDs are length 11
								vo.YTID = videoId
							}
						}

					}
					if vo.YTID != "" {
						voChan <- vo
					}
				}
			}
		}
	}()
	return voChan
}

func getTopVideos(c *mgo.Collection, lang string) ([]topVideoResult, error) {

	q := []bson.M{bson.M{"$group": bson.M{"_id": "$ytid", "count": bson.M{"$sum": 1}, "langs": bson.M{"$addToSet": "$language"}}},
		bson.M{"$sort": bson.M{"count": -1}}}

	if lang != "" {
		q = append(q, bson.M{"$match": bson.M{"langs": lang}})
	}

	q = append(q, bson.M{"$limit": 10})
	pipe := c.Pipe(q)

	r := topVideoResult{}

	results := make([]topVideoResult, 0)
	iter := pipe.Iter()
	for iter.Next(&r) {
		results = append(results, r)
	}

	err := iter.Err()
	if err != nil {
		return nil, err
	}

	return results, nil

}

func TopVideos(w http.ResponseWriter, r *http.Request) {
	var results []topVideoResult
	session, err := mgo.Dial("localhost:27017")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	c := session.DB("twitwyrm").C("VideoOccurences")
	lang, ok := r.URL.Query()["lang"]
	if ok {
		results, err = getTopVideos(c, lang[0])
	} else {
		results, err = getTopVideos(c, "")
	}
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(results)

}
