package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"tweet_clusterer"
	"twitter_client"
)

func Usage() {
	fmt.Println("Usage:")
	fmt.Print("  --consumerkey <consumerkey>")
	fmt.Println("  --consumersecret <consumersecret>")
	fmt.Println("  --accesstoken <accesstoken>")
	fmt.Println("  --accesstokensecret <accesstokensecret>")
	fmt.Println("  --q <query_term>")
	fmt.Println("")
	fmt.Println("In order to get your consumerkey and consumersecret, you must register an 'app' at twitter.com:")
	fmt.Println("https://dev.twitter.com/apps/new")
}

func main() {
	var consumerKey *string = flag.String(
		"consumerkey",
		"",
		"Consumer Key from Twitter. See: https://dev.twitter.com/apps/new")

	var consumerSecret *string = flag.String(
		"consumersecret",
		"",
		"Consumer Secret from Twitter. See: https://dev.twitter.com/apps/new")

	var accessToken *string = flag.String(
		"accesstoken",
		"",
		"Generate accesstoken from Twitter. See: https://dev.twitter.com/apps/new")
		
	var accessTokenSecret *string = flag.String(
		"accesstokensecret",
		"",
		"Accesstokensecret from Twitter. See: https://dev.twitter.com/apps/new")
	
	var query *string = flag.String(
		"q",
		"",
		"Query term we want to search by. Example --q superbowl")
	
	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 || 
	   len(*accessToken) == 0 || len(*accessTokenSecret) == 0 || len(*query) == 0 {
		fmt.Println("You must set the --consumerkey, --consumersecret, --accesstoken, --accesstokensecret and --q flags.")
		fmt.Println("---")
		Usage()
		os.Exit(1)
	}

	bits, err := twitter_client.Search_tweets(consumerKey, consumerSecret, accessToken, accessTokenSecret, query)
	if err != nil {
                      log.Fatal(err)
             }
	tweet_list := twitter_client.Get_ordered_tweets(bits)
	
	//fmt.Println(len(tweet_list))
	tweet_clusterer.Print_clusters(tweet_clusterer.Cluster_tweets(tweet_list, 0.1))
}
