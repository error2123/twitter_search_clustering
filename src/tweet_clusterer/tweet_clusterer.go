package tweet_clusterer

import (
	"bytes"
	"fmt"
	"strings"
	"tweet_similarity"
)

/********************************************************************
* Cluster struct:
* best_result: Holds the best tweet
* first_result: Holds the earliest tweet in that cluster
* rest_results: list of all the other tweets ordered chrono
* logically
* all_appended: Holds a concatenated version of all the tweets
* in the cluster(used to find the closest cluster given an
* unclassified tweet)
********************************************************************/
type Cluster struct {
	best_result  string
	first_result string
	rest_results []string
	all_appended string
}

/********************************************************************
* Helper function to dynamically grows a string array that holds a
* list of tweets for a given cluster
* original: original array
* position: position to insert the new value
* value: the new value to be inserted
* returns: a list with a new value appended
*********************************************************************/
func insert(original []string, position int, value string) []string {
	l := len(original)
	target := original
	if cap(original) == l {
		target = make([]string, l+1, l+10)
		copy(target, original[:position])
	} else {
		target = append(target, "")
	}
	copy(target[position+1:], original[position:])
	target[position] = value
	return target
}

/********************************************************************
* Helper function to efficiently concat strings using bytes.Buffer
* instead of "+"
* s1: string 1
* s2: string 2
* returns: s1 + s2
*********************************************************************/
func string_concat(s1 string, s2 string) string {
	var buffer bytes.Buffer
	buffer.WriteString(s1)
	buffer.WriteString(s2)
	return buffer.String()
}

/********************************************************************
* Remove top stop words(from reuters v1)
* instead of "+"
* s: sentence
* returns: s with no stop words
*********************************************************************/
func remove_stopwords(s string) string {
	stopwords := []string{"a", "an", "and", "are", "as", "at",
		"be", "by", "for", "from", "has", "he",
		"in", "is", "it", "its", "of", "on",
		"that", "the", "to", "was", "were",
		"will", "with", "rt", "http", "&amp;"}
	// surrounding s with spaces to match stopwords in boundaries
	s = " " + s + " "
	for _, k := range stopwords {
		s = strings.Replace(s, " "+k+" ", " ", -1)
	}
	return s
}

/********************************************************************
* Removed non-alphanumeric characters that may be distorting the
* cosine similarity
* s: sentence
* returns: cleaned sentence
*********************************************************************/
func clean_string(s string) string {
	// a few characters which we might want to remove that
	// might make cosine similarity perform better
	s = strings.Replace(s, "#", "", -1)
	s = strings.Replace(s, "'", "", -1)
	s = strings.Replace(s, "/", " ", -1)
	s = strings.Replace(s, "!", "", -1)
	s = strings.Replace(s, ":", "", -1)
	return s
}

//TODO
/*********************************************************************
* Clusters tweets into N clusters in one pass(greedy). It assigns a tweet
* into a new cluster if its similarity with all of the other clusters is
* less than X(configurable). Otherwise it assigns the tweet to the
* cluster which is closest to the tweet.
* tweets: list of tweets(usually ordered by time)
* thresh: threshold below which the tweet is allowed to start a new
* cluster
**********************************************************************/
func Cluster_tweets(tweets []string, thresh float64) map[int]Cluster {

	cluster_list := make(map[int]Cluster)

	// counter for tweet and cluster count
	cluster_cnt := 0
	cnt := 0
	for k, _ := range tweets {
		// first one starts a cluster
		if cnt == 0 {
			cluster_list[cluster_cnt] = start_new_cluster(tweets[k])
			cluster_cnt = cluster_cnt + 1
		} else {
			// variables that track the closest cluster(by similarity)
			closest_cluster := 0
			max_sim := -1000.00
			// loop thru the cluster list to get the closest match
			for c := range cluster_list {
				// compute tweet similarity to all the tweets currently present in that cluster
				sim := tweet_similarity.Cosine(remove_stopwords(clean_string(strings.ToLower(tweets[k]))),
					cluster_list[c].all_appended)
				if sim > max_sim {
					max_sim = sim
					closest_cluster = c
				}

			}
			// cluster assignment
			// if: the maximum similarity is lower than threshold start your own cluster
			// else: assign urself to the closest cluster
			if max_sim < thresh {
				cluster_list[cluster_cnt] = start_new_cluster(tweets[k])
				cluster_cnt = cluster_cnt + 1

			} else {
				cluster := cluster_list[closest_cluster]
				cluster.all_appended = string_concat(cluster.all_appended, tweets[k])
				cluster.rest_results = insert(cluster.rest_results, len(cluster.rest_results), tweets[k])
				cluster_list[closest_cluster] = cluster
			}
		}
		cnt = cnt + 1
	}
	// recomputes best tweets now we know the cluster assignments
	cluster_list = find_best_tweet(cluster_list)
	return cluster_list
}

/********************************************************************
* Initializes a new cluster and returns
* tweet: tweet that starts the new cluster
* returns: a new cluster
*********************************************************************/
func start_new_cluster(tweet string) Cluster {
	cluster := Cluster{}
	cluster.best_result = tweet
	cluster.first_result = tweet
	cluster.all_appended = remove_stopwords(clean_string(strings.ToLower(tweet)))
	return cluster
}

/********************************************************************
* Recomputes the best tweet by finding the tweet with min similarity
* with rest of the tweets(proxy to centroid)
* tweet: tweet that starts the new cluster
* returns: a new cluster
*********************************************************************/
func find_best_tweet(clusters map[int]Cluster) map[int]Cluster {
	for k := range clusters {
		cluster := clusters[k]
		best_sim := tweet_similarity.Cosine(remove_stopwords(clean_string(strings.ToLower(cluster.first_result))),
			cluster.all_appended)
		best_idx := -1
		for idx, t := range cluster.rest_results {
			sim := tweet_similarity.Cosine(remove_stopwords(clean_string(strings.ToLower(t))),
				cluster.all_appended)
			if sim > best_sim {
				best_sim = sim
				best_idx = idx
			}
		}
		if best_idx != -1 {
			cluster.best_result = cluster.rest_results[best_idx]

		}
		clusters[k] = cluster
	}
	return clusters
}

/********************************************************************
* Prints the clusters out in below format
* "Cluster: X
* best result: YYYYY
* first result: ZZZZ
* rest_result: AAAA, BBBB
* clusters: the populated list of clusters
*********************************************************************/
func Print_clusters(clusters map[int]Cluster) {
	for k := range clusters {
		fmt.Println("********************************************************************")
		fmt.Println("\n"+"Cluster: ", k)
		fmt.Println("Best result: " + strings.Replace(clusters[k].best_result, "\n", "", -1))
		fmt.Println("First result: " + strings.Replace(clusters[k].first_result, "\n", "", -1))
		fmt.Println("Rest results:")
		for _, c := range clusters[k].rest_results {
			fmt.Println(strings.Replace(c, "\n", "", -1))
		}
		fmt.Println("********************************************************************")
	}
}
