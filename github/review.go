package github

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type pendingReview struct {
	Id int `json:"id"`
}

type ReviewEvent string

const (
	ReviewEventApprove        ReviewEvent = "APPROVE"
	ReviewEventComment        ReviewEvent = "COMMENT"
	ReviewEventRequestChanges ReviewEvent = "REQUEST_CHANGES"
)

func (c *Client) StartReview(pr string) (int, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%s/reviews", c.owner, c.repo, pr)

	var review pendingReview
	if err := c.restClient.Post(endpoint, nil, &review); err != nil {
		return 0, fmt.Errorf("could not start review: %w", err)
	}

	return review.Id, nil
}

type submitReview struct {
	Event string `json:"event"`
	Body  string `json:"body"`
}

func (c *Client) CompleteReview(pr string, reviewId int, event ReviewEvent) error {
	endpoint := fmt.Sprintf(
		"repos/%s/%s/pulls/%s/reviews/%d/events",
		c.owner, c.repo, pr, reviewId)

	review := submitReview{
		Event: string(event),
		Body:  "made with the cool PR tool",
	}

	json, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("could not marshal review: %w", err)
	}

	if err := c.restClient.Post(endpoint, bytes.NewReader(json), nil); err != nil {
		return fmt.Errorf("could not complete review: %w", err)
	}

	return nil
}
