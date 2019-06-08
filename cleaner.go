package twittercleaner

import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"fmt"
	"errors"
)

type cleaner struct {
	client *anaconda.TwitterApi
}

func New() (*cleaner, error) {
	client, err := client()
	if err != nil {
		return nil, err
	}

	return &cleaner{
		client: client,
	}, nil
}

func (cleaner *cleaner) getTweets() (*[]anaconda.Tweet, error) {

	search, err := cleaner.client.GetSearch("colonial williamsburg since:2019-06-03 until:2019-06-08", nil)
	if err != nil {
		return nil, err
	}

	tweets := []anaconda.Tweet{}
	for ; len(search.Statuses) > 0; {
		tweets = append(tweets, search.Statuses...)
		fmt.Printf("retrieving tweets, %d found so far\n", len(tweets))
		search, err = search.GetNext(cleaner.client)
		fmt.Printf("found %d more tweets\n", len(search.Statuses))
	}
	return &tweets, nil
}

func client() (*anaconda.TwitterApi, error) {
	errs := compoundError{}

	accessToken, err := getEnv("ACCESS_TOKEN")
	errs.Add(err)

	accessTokenSecret, err := getEnv("ACCESS_TOKEN_SECRET")
	errs.Add(err)

	consumerKey, err := getEnv("CONSUMER_KEY")
	errs.Add(err)

	consumerSecretKey, err := getEnv("CONSUMER_SECRET_KEY")
	errs.Add(err)

	err = errs.Error()
	if err != nil {
		return nil, err
	}

	return anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecretKey), nil
}

type compoundError []error

func (c *compoundError) Add(e error) {
	if e != nil {
		*c = append(*c, e)
	}
}

func (c *compoundError) Error() (err error) {
	if len(*c) == 0 {
		return nil
	}

	msg := "Collected errors:\n"
	for i, e := range *c {
		msg += fmt.Sprintf("\tError %d: %s\n", i, e.Error())
	}
	return errors.New(msg)
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		msg := fmt.Sprintf("no value for for environment variable %s", key)
		err := errors.New(msg)
		return "", err
	}
	return value, nil
}