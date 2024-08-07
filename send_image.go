package qywx

func (c *Client) SendImage(title, text, imageUrl, url, userId string) error {
	if imageUrl == "" {
		return c.SendMessage(title, text, userId, url)
	}
	if userId == "" {
		userId = "@all"
	}
	return c.send(imageMessage{
		ToUser:  userId,
		MsgType: "news",
		AgentID: c.agentId,
		News: news{
			Articles: []Article{
				{
					Title:       title,
					Description: text,
					PicURL:      imageUrl,
					URL:         url,
				},
			},
		},
	}, true)
}
