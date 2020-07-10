package main

import (
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var repositoriesQuery struct {
	User struct {
		Repositories struct {
			Nodes    []Repository
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"repositories(first: 100, isFork: false, after: $repositoriesCusor)"`
	} `graphql:"user(login: $login)"`
}

type Repository struct {
	Name githubv4.String
}

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_OAUTH_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	repositoriesVar := map[string]interface{}{
		"login":             githubv4.String("mattn"),
		"repositoriesCusor": (*githubv4.String)(nil),
	}
	for {
		err := client.Query(context.Background(), &repositoriesQuery, repositoriesVar)
		if err != nil {
			fmt.Println(err)
		}
		for _, repository := range repositoriesQuery.User.Repositories.Nodes {
			fmt.Println(repository.Name)
		}
		if !repositoriesQuery.User.Repositories.PageInfo.HasNextPage {
			break
		}
		repositoriesVar["repositoriesCusor"] = githubv4.NewString(repositoriesQuery.User.Repositories.PageInfo.EndCursor)
	}
}
