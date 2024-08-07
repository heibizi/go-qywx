package qywx

type (
	accessTokenResp struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	sendMessageResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	message struct {
		ToUser               string      `json:"touser"`
		MsgType              string      `json:"msgtype"`
		AgentID              string      `json:"agentid"`
		Text                 messageText `json:"text"`
		Safe                 int         `json:"safe"`
		EnableIDTrans        int         `json:"enable_id_trans"`
		EnableDuplicateCheck int         `json:"enable_duplicate_check"`
	}

	messageText struct {
		Content string `json:"content"`
	}

	imageMessage struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		AgentID string `json:"agentid"`
		News    news   `json:"news"`
	}

	news struct {
		Articles []Article `json:"articles"`
	}

	Article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PicURL      string `json:"picurl"`
		URL         string `json:"url"`
	}
)
