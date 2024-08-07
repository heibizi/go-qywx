package qywx

import (
	"fmt"
	"strings"
)

func (c *Client) SendMessage(title, text, userId, url string) error {
	var content string
	if text != "" {
		content = fmt.Sprintf("%s\n%s", title, strings.ReplaceAll(text, "\n\n", "\n"))
	} else {
		content = title
	}
	if url != "" {
		content = fmt.Sprintf("%s\n\n<a href='%s'>查看详情</a>", content, url)
	}
	if userId == "" {
		userId = "@all"
	}
	return c.send(message{
		ToUser:               userId,
		MsgType:              "text",
		AgentID:              c.agentId,
		Text:                 messageText{content},
		Safe:                 0,
		EnableIDTrans:        0,
		EnableDuplicateCheck: 0,
	}, true)
}
