package main

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	anaconda.SetConsumerKey("tDAaN4yvVg5wFVlfkxeEUOzHm")
	anaconda.SetConsumerSecret("fRtI78PU9T67xNM5dylEO8146qEjtF5qVEHHZIBLuoJ6BHHjm0")
	api := anaconda.NewTwitterApi("332989840-mBSNofbvaHW2eyENXlAmI52GLpyGQwVvRav5jdsl",
		"y6N3Ty0d99JDARs4FPvHmYjtocUrpQlRcN2uQ9ljTWeBu")

	values := url.Values{}
	values.Add("track", "youtube com,youtu be")
	stream := api.PublicStreamFilter(values)
	for item := range stream.C {
		fmt.Println(item)

	}
}
