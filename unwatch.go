package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

const Version string = "0.0.1"

func main() {
	var target *string = flag.String("target", "", "Unwatach repository name(msut)")
	var token *string = flag.String("token", "", "Repository API Token(must)")

	flag.Parse()

	if len(*target) == 0 || len(*token) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// auth
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// unwatch
	opt := &github.ListOptions{PerPage: 100}

	var allRepos []github.Repository

	for {
		repos, resp, err := client.Activity.ListWatched("", opt)

		if err != nil {
			break
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	for _, repo := range allRepos {
		if strings.Contains(*repo.FullName, *target) {
			owner := *repo.Owner.Login
			name := *repo.Name
			_, err := client.Activity.DeleteRepositorySubscription(owner, name)
			if err != nil {
				fmt.Printf("Failed : Unwatch %s/%s\n", owner, name)
			}
			fmt.Printf("Success : Unwatch %s/%s\n", owner, name)
		}
	}
}
