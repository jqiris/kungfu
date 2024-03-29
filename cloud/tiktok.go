/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package cloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type VerifyUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SdkOpenId string `json:"sdk_open_id"`
		Nickname  string `json:"nickname"`
		AvatarUrl string `json:"avatar_url"`
		AgeType   int32  `json:"age_type"`
	} `json:"data"`
	LogId string `json:"log_id"`
}

type TiktokClient struct {
	appId     int32
	secretKey string
}

func NewTiktokClient(appId int32, secretKey string) *TiktokClient {
	return &TiktokClient{
		appId:     appId,
		secretKey: secretKey,
	}
}

func (t *TiktokClient) GetSign(params map[string]any) string {
	//将key排序
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//格式化，拼接元素
	ss := []string{}
	for i, k := range keys {
		if k == "sign" {
			continue
		}
		if i > 0 {
			ss = append(ss, "&")
		}
		ss = append(ss, fmt.Sprintf("%v=%v", k, params[k]))
	}
	content := strings.Join(ss, "")

	//使用密钥进行Hmac-sha1加密
	mac := hmac.New(sha1.New, []byte(t.secretKey))
	mac.Write([]byte(content))

	//base64编码获得最终的sign
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (t *TiktokClient) VerifyUser(token string) (*VerifyUserResponse, error) {
	targetUrl := "https://gsdk.snssdk.com/gsdk/usdk/account/verify_user"
	params := map[string]any{
		"app_id":       t.appId,
		"access_token": token,
		"ts":           time.Now().Unix(),
	}
	sign := t.GetSign(params)
	params["sign"] = sign
	data := url.Values{}
	for k, v := range params {
		data.Add(k, fmt.Sprintf("%v", v))
	}
	resp, err := http.PostForm(targetUrl, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &VerifyUserResponse{}
	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}
