package qywx

import "fmt"

func (c *Client) SendList(items []Article, userId string) error {
	if userId == "" {
		userId = "@all"
	}
	var articles []Article
	for index, item := range items {
		item.Title = fmt.Sprintf("%d. %s", index+1, item.Title)
		articles = append(articles, item)
	}
	return c.send(imageMessage{
		ToUser:  userId,
		MsgType: "news",
		AgentID: c.agentId,
		News: news{
			Articles: articles,
		},
	}, true)
}
