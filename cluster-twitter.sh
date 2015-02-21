#!/bin/bash

######################################################################
# Runs twitter search clustering given a query.
# e.g.: cluster-twitter.sh -q='amazon'
#
######################################################################

# please change the path depending on where your project directories
# are present
PROJECT_PATH=/Users/thebanstanley/Documents/GO_2/twitter_search_clustering/

# set the right exports
export PATH=$PATH:/usr/local/go/bin:$PROJECT_PATH/bin
export GOPATH=$PROJECT_PATH
export GOBIN=$PROJECT_PATH/bin


# set twitter credentials to access their apis
CONSUMER_KEY="aNks5OghAzGyktYVbmhX4BdfW"
CONSUMER_SECRET="gb5nzctfxQdg2nb8MoKXSJyJKnEsp5Sy65wWvkt4t4EadsTGn2"
ACCESS_TOKEN="3021764810-h1FDUx29E2VrBUziMv9ZfADQPaBpVdcTO2dKhsm"
ACCESS_TOKEN_SECRET="IfcaL9L3xZ7oTdcSkDjZ4PcrxRNXQzHw7ooFgElmuCbh1"

cd $PROJECT_PATH

# install the necessary packages
echo "Installing Packages"
go get github.com/mrjones/oauth
go install tweet_similarity
go install tweet_clusterer
go install twitter_client

go install src/twitter_search.go
echo "Installation done!"

QUERY=`echo $1 | cut -d \= -f 2`
twitter_search --consumerkey $CONSUMER_KEY --consumersecret $CONSUMER_SECRET --accesstoken $ACCESS_TOKEN --accesstokensecret $ACCESS_TOKEN_SECRET --q $QUERY