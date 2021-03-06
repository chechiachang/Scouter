package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/chechiachang/scouter"
	"golang.org/x/oauth2"
)

func main() {
	githubApiToken := flag.String("token", "", "github api token (string)")
	flag.Parse()

	if *githubApiToken == "" {
		panic("Github api token is empty.")
	}

	log.Println("Starting crawler...")

	// Prepare github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubApiToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	if err := searchUsers(tc); err != nil {
		log.Fatal(err)
	}

}

func searchUsers(tc *http.Client) error {
	log.Println("Starting fetch github user with search api...")

	layout := "2006-01-01T00:00:00"
	// set fetching with time range from start time to now
	endTime := time.Now()
	startTime, err := time.Parse(layout, "2008-01-01T00:00:00")
	if err != nil {
		return err
	}

	total, err := scouter.CountGithubUsers(tc, "location:taiwan")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Total: ", total)

	// set fetch batch time interval
	startCursor := startTime
	endCursor := startCursor.AddDate(0, 1, 0) // interval: 1 month
	sort := "joined"
	order := "asc"

	runtime.GOMAXPROCS(1)

	for endCursor.Before(endTime) {
		query := "location:Taiwan created:" + startCursor.Format(layout) + ".." + endCursor.Format(layout)

		// First fetch
		r, err := scouter.SearchGithubUsers(tc, 1, query, sort, order)
		time.Sleep(2 * time.Second) // Github search API max rate
		if err != nil {
			return err
		}
		log.Println("Fetching ", query, ". Found records:", r.GetTotal())

		// paging fetch if result.Total > searchMaxPerPage
		if *r.Total > scouter.SearchMaxPerPage {

			log.Fatal("Pagesize exceed ", scouter.SearchMaxPerPage, ". Some data may not be fetched")

		} else {

			for _, user := range r.Users {
				u := scouter.User{
					ID:   user.GetID(),
					User: &user,
				}
				if err := scouter.UpsertUser(u); err != nil {
					return err
				}
			}
		}

		// Move cursor forward 1 month
		startCursor = startCursor.AddDate(0, 1, 0)
		endCursor = endCursor.AddDate(0, 1, 0) // interval: 1 month

	}
	return nil
}
