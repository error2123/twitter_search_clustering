**What the repo does?**
Given a query, the program hits the twitter search api and returns a clustered list of tweets. 

**How does the clustering work?**
* The current version of clustering is a greedy one-pass algorithm.
* Tweets are time ordered as they get parsed from the twitter search result.
* As each tweet is processed according to the following logic:
 * if the tweet's distance is beyond a threshold from all existing clusters, it starts a new cluster
 * lse, gets assigned to the closest cluster
* Cosine similiarity is used for the distance measure.


**How to run it?**
* Make sure you have GO installed.
* open cluster-twitter.sh and set the porject path
* run: cluster-twitter.sh -q='amazon'