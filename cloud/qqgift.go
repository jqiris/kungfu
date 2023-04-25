package cloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

type RoleInfo struct {
	RoleID        string `json:"roleid"`
	RoleName      string `json:"roleName"`
	CreateTime    int64  `json:"createTime"`
	Level         int    `json:"level"`
	Duration      int    `json:"duration"`
	Power         int    `json:"power"`
	LastLoginTime int64  `json:"lastLoginTime"`
}

type ServerInfo struct {
	ServerID   string     `json:"serverid"`
	ServerName string     `json:"serverName"`
	StartTime  int64      `json:"startTime"`
	Status     int        `json:"status"`
	ActArray   []RoleInfo `json:"actArray"`
}

type ServerList struct {
	Ret             int          `json:"ret"`
	Msg             string       `json:"msg"`
	RecentServer    []ServerInfo `json:"recentServer"`
	RecommendServer []ServerInfo `json:"recommendServer"`
	AllServer       []ServerInfo `json:"allServer"`
}

type RoleList struct {
	Ret            int          `json:"ret"`
	Msg            string       `json:"msg"`
	ServerRoleInfo []ServerInfo `json:"serverRoleInfo"`
}

type GiftParam struct {
	AppID     int    `json:"appid" form:"appid"`
	PF        string `json:"pf,omitempty" form:"pf"`
	OpenID    string `json:"openid" form:"openid"`
	ServerID  string `json:"serverid" form:"serverid"`
	RoleID    string `json:"roleid" form:"roleid"`
	GiftID    string `json:"giftid" form:"giftid"`
	BillNo    string `json:"billno" form:"billno"`
	Timestamp uint64 `json:"ts" form:"ts"`
	Signature string `json:"sig" form:"sig"`
}

type GiftResult struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type IQQGiftClient interface {
	QQServerList(appId int, pf string) ServerList
	QQRoleList(appId int, pf, openId, serverId string) RoleList
	SendQQGift(param GiftParam, uri string) GiftResult
}

type QQGiftClient struct {
	client IQQGiftClient
}

func NewQQGiftClient(client IQQGiftClient) *QQGiftClient {
	return &QQGiftClient{
		client: client,
	}
}

func (c *QQGiftClient) QQServerList(appId int, pf string) ServerList {
	return c.client.QQServerList(appId, pf)
}

func (c *QQGiftClient) QQRoleList(appId int, pf, openId, serverId string) RoleList {
	return c.client.QQRoleList(appId, pf, openId, serverId)
}

func (c *QQGiftClient) SendQQGift(param GiftParam, uri string) GiftResult {
	return c.client.SendQQGift(param, uri)
}

func QQGiftSign(method string, uri string, params url.Values, appKey string) string {
	// Step 1: Construct source string
	var keys []string
	for k := range params {
		if k != "sig" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var sourceString string
	sourceString += strings.ToUpper(method) + "&"
	sourceString += url.QueryEscape(uri) + "&"
	var paramStrings []string
	for _, k := range keys {
		paramStrings = append(paramStrings, k+"="+url.QueryEscape(params.Get(k)))
	}
	sourceString += url.QueryEscape(strings.Join(paramStrings, "&"))

	// Step 2: Construct key
	key := appKey + "&"

	// Step 3: Generate signature
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(sourceString))
	sig := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sig
}
