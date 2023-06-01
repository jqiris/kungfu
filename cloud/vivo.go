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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/jqiris/kungfu/v2/logger"
)

type VivoLoginResp struct {
	RetCode int `json:"retcode"`
	Data    struct {
		OpenID string `json:"openid"`
	} `json:"data"`
}

type VivoClient struct {
	apiList []string
	apiKey  int
	apiLen  int
	lock    sync.RWMutex
}

func NewVivoClient(list []string) *VivoClient {
	return &VivoClient{
		apiList: list,
		apiKey:  0,
		apiLen:  len(list),
		lock:    sync.RWMutex{},
	}
}

func (b *VivoClient) GetApiKey() int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.apiKey
}

func (b *VivoClient) SetApiKey(key int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	logger.Infof("VivoClient SetApiKey: %d", key)
	b.apiKey = key
}

func (b *VivoClient) Req(api string, tryTimes, apiKey int, params map[string]string) ([]byte, error) {
	if tryTimes == b.apiLen {
		logger.Errorf("vivo全部请求失败,tryTimes:%v,apiKey:%v,params:%v", tryTimes, apiKey, params)
		return nil, fmt.Errorf("请求失败")
	}
	apiKey = apiKey % b.apiLen
	apiUrl := b.apiList[apiKey] + api
	req, err := http.NewRequest(http.MethodPost, apiUrl, nil)
	if err != nil {
		logger.Errorf("vivo请求失败 NewRequest,tryTimes:%v,apiKey:%v,params:%v,err:%v", tryTimes, apiKey, params, err)
		return b.Req(api, tryTimes+1, apiKey+1, params)
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.PostForm = data
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("vivo请求失败 Do,tryTimes:%v,apiKey:%v,params:%v,err:%v", tryTimes, apiKey, params, err)
		return b.Req(api, tryTimes+1, apiKey+1, params)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("vivo请求失败 ReadAll,tryTimes:%v,apiKey:%v,params:%v,err:%v", tryTimes, apiKey, params, err)
		return nil, err
	}
	if tryTimes > 0 {
		b.SetApiKey(apiKey)
	}
	return body, nil
}

func (b *VivoClient) Login(authtoken string) (*VivoLoginResp, error) {
	params := map[string]string{
		"opentoken": authtoken,
	}
	api := "/cp/user/auth"
	body, err := b.Req(api, 0, b.GetApiKey(), params)
	if err != nil {
		return nil, err
	}
	resp := &VivoLoginResp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
