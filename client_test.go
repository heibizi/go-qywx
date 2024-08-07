package qywx

import (
	"github.com/redis/go-redis/v9"
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	addr := os.Getenv("GO_QYWX_REDIS_ADDR")
	corpId := os.Getenv("GO_QYWX_CORP_ID")
	corpSecret := os.Getenv("GO_QYWX_CORP_SECRET")
	agentId := os.Getenv("GO_QYWX_AGENT_ID")
	var redisClient redis.UniversalClient = redis.NewClient(&redis.Options{Addr: addr})
	client = NewClient(corpId, corpSecret, agentId, nil, redisClient)
	m.Run()
}

func TestWorkWX(t *testing.T) {
	err := client.SendMessage("这是一条测试消息", "test", "", "https://www.baidu.com")
	if err != nil {
		t.Error(err)
	}
}

func TestWorkWXImageMessage(t *testing.T) {
	err := client.SendImage("这是一条测试消息", "test",
		"https://wwcdn.weixin.qq.com/node/wework/images/Pic_right@2x.7a03a9d992.png",
		"https://www.baidu.com", "")
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListMessage(t *testing.T) {
	err := client.SendList([]Article{
		{Title: "这是一条测试消息", Description: "", URL: "https://www.baidu.com",
			PicURL: "https://wwcdn.weixin.qq.com/node/wework/images/Pic_right@2x.7a03a9d992.png"},
		{Title: "这是一条测试消息", Description: "", URL: "https://www.baidu.com",
			PicURL: "https://wwcdn.weixin.qq.com/node/wework/images/Pic_right@2x.7a03a9d992.png"},
	}, "")
	if err != nil {
		t.Error(err)
	}
}
