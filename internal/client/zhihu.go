package client

import (
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) Me() (Me, error) {
	var me Me
	err := c.GetJSON(c.endpoints.MeURL, nil, &me)
	return me, err
}

func (c *Client) Feed(limit int) (FeedResponse, error) {
	params := url.Values{}
	params.Set("page_number", "1")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("action", "down")

	var feed FeedResponse
	err := c.GetJSON(c.endpoints.FeedRecommendURL, params, &feed)
	return feed, err
}

func (c *Client) Answer(answerID string) (Answer, error) {
	params := url.Values{}
	params.Set("include", "content,voteup_count,comment_count,created_time,updated_time,author,question")

	var answer Answer
	err := c.GetJSON(fmt.Sprintf(c.endpoints.AnswerURLTemplate, answerID), params, &answer)
	return answer, err
}

func (c *Client) AnswerComments(answerID string, limit int) (CommentsResponse, error) {
	params := url.Values{}
	params.Set("offset", "0")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("order", "normal")
	params.Set("status", "open")

	var comments CommentsResponse
	err := c.GetJSON(fmt.Sprintf(c.endpoints.AnswerCommentsURLTemplate, answerID), params, &comments)
	return comments, err
}
