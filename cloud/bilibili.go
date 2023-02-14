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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/utils"
)

type BilibiliConfig struct {
	Version    string   `json:"version" mapstructure:"version"`
	GameId     int      `json:"game_id" mapstructure:"game_id"`
	MerchantId int      `json:"merchant_id" mapstructure:"merchant_id"`
	SecretKey  string   `json:"secret_key" mapstructure:"secret_key"`
	ApiList    []string `json:"api_list" mapstructure:"api_list"`
}

type BilibiliCommonResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

type BilibiliLoginResp struct {
	BilibiliCommonResp
	OpenId int    `json:"open_id"`
	Uname  string `json:"uname"`
}

type BilibiliAuthResp struct {
	BilibiliCommonResp
	Data struct {
		Uid        int    `json:"uid"`
		AuthStatus int    `json:"auth_status"`
		AgeRange   string `json:"age_range"`
	} `json:"data"`
}
type BilibiliLogoutReq struct {
	VoucherNo string `json:"voucher_no"`
	Uid       int    `json:"uid"`
	GameId    int    `json:"game_id"`
	Sign      string `json:"sign"`
}

type BilibiliClient struct {
	config BilibiliConfig
	apiKey int
	apiLen int
	lock   sync.RWMutex
}

func NewBilibiliClient(conf BilibiliConfig) *BilibiliClient {
	return &BilibiliClient{
		config: conf,
		apiKey: 0,
		apiLen: len(conf.ApiList),
		lock:   sync.RWMutex{},
	}
}

func (b *BilibiliClient) GetApiKey() int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.apiKey
}

func (b *BilibiliClient) SetApiKey(key int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	logger.Infof("BilibiliClient SetApiKey: %d", key)
	b.apiKey = key
}

func (b *BilibiliClient) GetSign(params map[string]any) string {
	//将key排序
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//格式化，拼接元素
	ss := []string{}
	for _, k := range keys {
		if k == "sign" || k == "item_name" || k == "item_desc" {
			continue
		}
		ss = append(ss, fmt.Sprintf("%v", params[k]))
	}
	ss = append(ss, b.config.SecretKey)
	content := strings.Join(ss, "")
	return utils.Md5(content)
}

func (b *BilibiliClient) Req(api string, tryTimes, apiKey int, params map[string]any) ([]byte, error) {
	if tryTimes == b.apiLen {
		logger.Errorf("bilibili全部请求失败,tryTimes:%v,apiKey:%v,params:%v", tryTimes, apiKey, params)
		return nil, fmt.Errorf("请求失败")
	}
	apiKey = apiKey % b.apiLen
	apiUrl := b.config.ApiList[apiKey] + api + "?"
	if tryTimes == 0 {
		params["game_id"] = b.config.GameId
		params["merchant_id"] = b.config.MerchantId
		params["version"] = b.config.Version
		params["timestamp"] = time.Now().UnixMilli()
		sign := b.GetSign(params)
		params["sign"] = sign
	}
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 GameServer",
		"Content-Type": "application/x-www-form-urlencoded",
	}
	var apiParams []string
	for k, v := range params {
		apiParams = append(apiParams, fmt.Sprintf("%v=%v", k, v))
	}
	apiUrl += strings.Join(apiParams, "&")
	req, err := http.NewRequest(http.MethodPost, apiUrl, nil)
	if err != nil {
		logger.Errorf("bilibili请求失败 NewRequest,tryTimes:%v,apiKey:%v,params:%v,err:%v", tryTimes, apiKey, params, err)
		return b.Req(api, tryTimes+1, apiKey+1, params)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("bilibili请求失败 Do,tryTimes:%v,apiKey:%v,params:%v,err:%v", tryTimes, apiKey, params, err)
		return b.Req(api, tryTimes+1, apiKey+1, params)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("bilibili请求失败 ReadAll,tryTimes:%v,apiKey:%v,params:%v,err:%v", tryTimes, apiKey, params, err)
		return nil, err
	}
	if tryTimes > 0 {
		b.SetApiKey(apiKey)
	}
	return body, nil
}

func (b *BilibiliClient) Login(accessKey, uid string) (*BilibiliLoginResp, error) {
	params := map[string]any{
		"access_key": accessKey,
		"uid":        uid,
	}
	api := "/api/server/session.verify"
	body, err := b.Req(api, 0, b.GetApiKey(), params)
	if err != nil {
		return nil, err
	}
	resp := &BilibiliLoginResp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BilibiliClient) Auth(uid string) (*BilibiliAuthResp, error) {
	params := map[string]any{
		"uid": uid,
	}
	api := "/api/server/user/auth/info"
	body, err := b.Req(api, 0, b.GetApiKey(), params)
	if err != nil {
		return nil, err
	}
	resp := &BilibiliAuthResp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
