package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
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

type Review struct {
	Event    ReviewEvent `json:"event"`
	Comments []Comment   `json:"comments,omitempty"`
	Body     string      `json:"body"`
}

type Comment struct {
	Body string `json:"body"`
	Path string `json:"path"`
	Line int    `json:"line"`
}

func (c *Client) Review(pr string, review Review) error {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%s/reviews", c.owner, c.repo, pr)

	json, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("could not marshal review: %w", err)
	}

	slog.Debug("generated request for review", "endpoint", endpoint, "body", string(json))

	if err := c.restClient.Post(endpoint, bytes.NewReader(json), nil); err != nil {
		return fmt.Errorf("could not create review: %w", err)
	}

	return nil
}
