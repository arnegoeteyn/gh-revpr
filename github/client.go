package github

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

type PullRequest struct {
	Head struct {
		Ref string `json:"ref"`
	} `json:"head"`
}

type Client struct {
	restClient *api.RESTClient

	owner string
	repo  string
}

func NewClient() (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("could not create GitHub API client: %w", err)
	}

	repo, err := repository.Current()
	if err != nil {
		return nil, fmt.Errorf("could not get current repository: %w", err)
	}

	return &Client{
		restClient: client,
		owner:      repo.Owner,
		repo:       repo.Name,
	}, nil
}

func (c *Client) GetPullRequestBranch(prNumber string) (string, error) {
	repo, err := repository.Current()
	if err != nil {
		return "", fmt.Errorf("could not get current repository: %w", err)
	}

	var pr PullRequest
	path := fmt.Sprintf("repos/%s/%s/pulls/%s", repo.Owner, repo.Name, prNumber)
	if err := c.restClient.Get(path, &pr); err != nil {
		return "", err
	}

	return pr.Head.Ref, nil
}
