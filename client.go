package qywx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	corpId      string
	corpSecret  string
	agentId     string
	baseURL     string
	redisClient redis.UniversalClient
}

func NewClient(corpId, corpSecret, agentId string, proxy *string, redisClient redis.UniversalClient) *Client {
	workWX := &Client{
		corpId:      corpId,
		corpSecret:  corpSecret,
		agentId:     agentId,
		redisClient: redisClient,
	}
	if proxy != nil {
		workWX.baseURL = *proxy
	} else {
		workWX.baseURL = "https://qyapi.weixin.qq.com"
	}
	return workWX
}

func (c *Client) getAccessToken(force bool) (string, error) {
	k := fmt.Sprintf("qywx:accessToken:%s:%s", c.corpId, c.agentId)
	if force {
		c.redisClient.Del(context.Background(), k)
	} else {
		accessToken := c.redisClient.Get(context.Background(), k).Val()
		if accessToken != "" {
			return accessToken, nil
		}
	}
	u, _ := url.JoinPath(c.baseURL, "/cgi-bin/gettoken")
	resp, err := HttpRequest(u, "GET", nil, map[string]string{
		"corpid":     c.corpId,
		"corpsecret": c.corpSecret,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("获取 accessToken 异常: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("获取 accessToken 异常: %v", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取 accessToken 响应异常: %v", err)
	}
	var atResp accessTokenResp
	err = json.Unmarshal(data, &atResp)
	if err != nil {
		return "", fmt.Errorf("解析 accessToken 响应异常: %v", err)
	}
	if atResp.ErrCode != 0 {
		return "", fmt.Errorf("获取 accessToken 异常，错误码: %d，错误信息: %s", atResp.ErrCode, atResp.ErrMsg)
	}
	c.redisClient.Set(context.Background(), k, atResp.AccessToken, time.Duration(atResp.ExpiresIn)*time.Second)
	return atResp.AccessToken, nil
}

// 向企业微信发送请求
func (c *Client) send(data any, retry bool) error {
	messageURL, _ := url.JoinPath(c.baseURL, "/cgi-bin/message/send")
	token, err := c.getAccessToken(false)
	if err != nil {
		return fmt.Errorf("获取 accessToken 异常: %v", err)
	}
	res, err := HttpRequest(messageURL, "POST", nil, map[string]string{"access_token": token}, data)
	if err != nil {
		return fmt.Errorf("发送消息异常: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("读取响应体异常: %v", err)
		}
		var sendResp sendMessageResp
		err = json.Unmarshal(body, &sendResp)
		if err != nil {
			return fmt.Errorf("解析响应体异常: %v", err)
		}
		if sendResp.ErrCode == 0 {
			return nil
		} else {
			if sendResp.ErrCode == 42001 {
				// 处理 AccessToken 过期的情况
				_, err := c.getAccessToken(true)
				if err != nil {
					return err
				}
				// 重试一次
				if retry {
					return c.send(data, false)
				} else {
					return fmt.Errorf("accessToken 已过期: %v", sendResp.ErrMsg)
				}
			}
			return fmt.Errorf("发送消息异常: %v", sendResp.ErrMsg)
		}
	} else {
		return fmt.Errorf("错误码：%d，错误原因：%s", res.StatusCode, res.Status)
	}
}
