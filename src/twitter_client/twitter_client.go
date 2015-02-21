package twitter_client

import (
	"encoding/json"
	"fmt"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"time"
)

/********************************************************************
* Takes a tweet map and returns a list of tweets ordered by time
* bits: unmarshalled json data
* returns: ordered list of tweets
*********************************************************************/
func Get_ordered_tweets(bits []byte) []string {

	// parse the json result
	tweet_map, time_map := parse_json_results(bits)
	tweet_list := make([]string, 0, len(tweet_map))

	// get the unix time keys so we can sort
	var keys []int
	for k := range time_map {
		keys = append(keys, k)
	}
	// sort the key
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	// read the tweets in sorted order
	for _, value := range keys {
		for _, id := range time_map[value] {
			tweet_list = append(tweet_list, tweet_map[id])
		}
	}
	return tweet_list
}

/********************************************************************
* Parse json results from the twitter search api. Assumes the message
* is a list of json dicts and looks for necessary fields needed for
* clustering
* bits: unmarshalled json data
* returns: a map of tweets(id => tweet) and map of time=> tweet(for)
* ordering data by time
*********************************************************************/
func parse_json_results(bits []byte) (map[string]string, map[int][]string) {

	// map for twitter id to tweets
	//var tweet_dict map[string]string
	tweet_dict := make(map[string]string)

	// map of tweeted_time to tweets
	//
	time_map := make(map[int][]string)

	// an empty interface that reads in the json response
	var f interface{}
	json.Unmarshal(bits, &f)

	m := f.(map[string]interface{})
	for _, v := range m {
		switch vv := v.(type) {
		case string:
			continue
		case int:
			continue
		// if its an interface its a tweet
		case []interface{}:
			for _, u := range vv {
				var id_str string
				var text string
				var created_at int
				for x, y := range u.(map[string]interface{}) {
					if reflect.TypeOf(x).Kind() == reflect.String {
						// look for keys=text,id_str or created_at
						// to build the hashmap
						if x == "text" {
							if str, ok := y.(string); ok {
								text = str
							} else {
								text = ""
							}
						}
						if x == "id_str" {
							if str, ok := y.(string); ok {
								id_str = str
							} else {
								id_str = "error_in_type_conversion"
							}
						}
						if x == "created_at" {
							if str, ok := y.(string); ok {
								created_at = get_unix_time(str)
							} else {
								created_at = 0
							}
						}
					}
				}
				tweet_dict[id_str] = text
				time_map[created_at] = append(time_map[created_at], id_str)

			}
		default:
			continue
		}
	}
	return tweet_dict, time_map
}

/********************************************************************
* Helper function that return unix time given a rubydate
* s: ruby date
* returns: unix time
*********************************************************************/
func get_unix_time(s string) int {
	//template to use for parsing
	const rfc2822 = "Mon Jan 02 15:04:05 -0700 2006"
	t, err := time.Parse(rfc2822, s)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	u := t.Unix()
	return int(u)
}

/********************************************************************
* Hits twitter search apis and returns back json
* consumerkey: from twitter app account
* consumersecret: from twitter app account
* accesstoken: generated in apps.twitter.com
* accesstokensecret: from twitter app account
* returns: ordered list of tweets
*********************************************************************/
func Search_tweets(consumerkey *string, consumersecret *string,
	accesstoken *string, accesstokensecret *string,
	query *string) (bits []byte, err error) {
	c := oauth.NewConsumer(
		*consumerkey,
		*consumersecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	ac_tok := &oauth.AccessToken{*accesstoken, *accesstokensecret, map[string]string{}}

	response, err := c.Get(
		"https://api.twitter.com/1.1/search/tweets.json",
		// restricting the search to English. Otherwise the clusters are
		// dominated by groupings by language
		map[string]string{"q": *query, "count": "100", "lang": "en"},
		ac_tok)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bits, readerr := ioutil.ReadAll(response.Body)
	return bits, readerr
}
