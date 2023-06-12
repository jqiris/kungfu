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
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type OppoConfig struct {
	ApiUrl    string `json:"api_url" mapstructure:"api_url"`
	AppKey    string `json:"app_key" mapstructure:"app_key"`
	AppSecret string `json:"app_secret"  mapstructure:"app_secret"`
}

type OppoLoginResponse struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	SsoID      int    `json:"ssoid"`
	UserName   string `json:"userName"`
	Email      string `json:"email"`
	Mobile     string `json:"mobileNumber"`
}
type OppoClient struct {
	config OppoConfig
}

func NewOppoClient(config OppoConfig) *OppoClient {
	return &OppoClient{
		config: config,
	}
}

func (c *OppoClient) Login(fileId, token string) (*OppoLoginResponse, error) {
	requestServerUrl := fmt.Sprintf("%s?fileId=%s&token=%s", c.config.ApiUrl, fileId, token)
	appvKey := c.config.AppKey
	appSecret := c.config.AppSecret

	dataParams := url.Values{}
	dataParams.Set("oauthConsumerKey", appvKey)
	dataParams.Set("oauthToken", url.QueryEscape(token))
	dataParams.Set("oauthSignatureMethod", "HMAC-SHA1")
	dataParams.Set("oauthTimestamp", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	dataParams.Set("oauthNonce", strconv.Itoa(int(time.Now().Unix())+rand.Intn(10)))
	dataParams.Set("oauthVersion", "1.0")

	requestString := dataParams.Encode()
	oauthSignature := appSecret + "&"
	sign := c.signatureNew(oauthSignature, requestString)
	body, err := c.oauthPostExecuteNew(sign, requestString, requestServerUrl)
	if err != nil {
		return nil, err
	}
	result := &OppoLoginResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 使用HMAC-SHA1算法生成签名
func (c *OppoClient) signatureNew(oauthSignature, requestString string) string {
	key := []byte(oauthSignature)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(requestString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(signature)
}

// Oauth身份认证请求
func (c *OppoClient) oauthPostExecuteNew(sign, requestString, requestServerUrl string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestServerUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("param", requestString)
	req.Header.Set("oauthsignature", sign)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
