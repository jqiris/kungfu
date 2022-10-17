package realauth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/utils"
	"github.com/sleagon/chinaid"
)

type RealAuthMgr struct {
	checkUrl    string
	queryUrl    string
	loginOutUrl string
	AppId       string
	BizId       string
	Secret      string
}

func NewRealAuthMgr(appId, bizId, secret string, args ...bool) *RealAuthMgr {
	isTest := false
	if len(args) > 0 {
		isTest = args[0]
	}
	checkUrl := "https://api.wlc.nppa.gov.cn/idcard/authentication/check"
	queryUrl := "http://api2.wlc.nppa.gov.cn/idcard/authentication/query"
	loginOutUrl := "http://api2.wlc.nppa.gov.cn/behavior/collection/loginout"
	if isTest { //是否是测试环境
		checkUrl = "https://wlc.nppa.gov.cn/test/authentication/check"
		queryUrl = "https://wlc.nppa.gov.cn/test/authentication/query"
		loginOutUrl = "https://wlc.nppa.gov.cn/test/collection/loginout"
	}
	mgr := &RealAuthMgr{
		checkUrl:    checkUrl,
		queryUrl:    queryUrl,
		loginOutUrl: loginOutUrl,
		AppId:       appId,
		BizId:       bizId,
		Secret:      secret,
	}
	return mgr
}
func (m *RealAuthMgr) encrypt(plaintext []byte) ([]byte, error) {
	key, err := hex.DecodeString(string(m.Secret))
	if err != nil {
		panic(err)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// logger.Debugf("nonce size: %d", gcm.NonceSize())

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (m *RealAuthMgr) getEncryptData(v interface{}) (*RequestBody, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	// logger.Debugf("body: %s", string(result))

	result, err = m.encrypt(result)
	if err != nil {
		return nil, err
	}
	encode := base64.StdEncoding.EncodeToString(result)
	// logger.Debugf("data: %v", encode)
	return &RequestBody{
		Data: encode,
	}, nil
}

// 获取 header sign 签名
func (m *RealAuthMgr) getHeader(params url.Values, v interface{}) (url.Values, error) {
	header := url.Values{}
	header.Add("appId", m.AppId)
	header.Add("bizId", m.BizId)
	// 填充 timestamps
	t := time.Now()
	header.Add("timestamps", strconv.FormatInt(t.UnixNano()/1000000, 10))

	// 因为 keys 的长度是固定的，所以此处这样写代码比较合理
	keys := make([]string, 0, len(header)+len(params))
	// header 中的 key
	for k := range header {
		keys = append(keys, k)
	}

	// params 中的 key
	for k := range params {
		keys = append(keys, k)
	}

	// 排序
	sort.Strings(keys)

	var requestBuf bytes.Buffer
	requestBuf.WriteString(m.Secret)
	for _, k := range keys {
		vs, ok := header[k]
		if ok {
			// 避免有 sign 的时候签名了数据
			if k == "sign" {
				continue
			}

			for _, v := range vs {
				requestBuf.WriteString(k)
				requestBuf.WriteString(v)
			}
		} else {
			vs, ok := params[k]
			if ok {
				for _, v := range vs {
					requestBuf.WriteString(k)
					requestBuf.WriteString(v)
				}
			}
		}

	}

	// 如果 body 不为 nil
	if v != nil {
		// json 序列化
		result, err := json.Marshal(v)
		if err != nil {
			return header, err
		}
		requestBuf.Write(result)
	}
	data := requestBuf.Bytes()
	// logger.Debugf("待签名字符串:%v", string(data))
	encrypt := fmt.Sprintf("%x", sha256.Sum256(data))
	header.Set("sign", encrypt)
	return header, nil
}

func (m *RealAuthMgr) getResponse(urlValue url.Values, v interface{}) (string, error) {
	var req *http.Request
	var err error
	if v != nil {
		result, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		req, err = http.NewRequest("POST", m.checkUrl, bytes.NewBuffer(result))
		if err != nil {
			return "", err
		}
	} else {
		req, err = http.NewRequest("POST", m.checkUrl, nil)
		if err != nil {
			return "", err
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range urlValue {
		req.Header.Set(k, v[0])
	}
	// logger.Debugf("header: %+v", req.Header)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// logger.Debugf("the raw is:%+v", string(body))
	responseData := &Response{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return "", err
	}

	// logger.Debugf("the response is:%+v", responseData)

	if responseData.ErrCode != 0 {
		return "", fmt.Errorf("%d -> %s", responseData.ErrCode, responseData.ErrMsg)
	}

	result := responseData.Data.Result
	if result.Status == 0 {
		return responseData.Data.Result.Pi, nil
	}

	switch result.Status {
	case 0:
		return responseData.Data.Result.Pi, nil
	case 1:
		return "", ErrNeedQuery
	case 2:
		return "", ErrAuthFailed
	}
	return "", fmt.Errorf("result status error: %d", result.Status)
}

func (m *RealAuthMgr) getResponseCheck(code string, urlValue url.Values, v interface{}) (string, error) {
	var req *http.Request
	var err error
	url := m.checkUrl + "/" + code
	if v != nil {
		result, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(result))
		if err != nil {
			return "", err
		}
	} else {
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return "", err
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range urlValue {
		req.Header.Set(k, v[0])
	}
	// logger.Debugf("header: %+v", req.Header)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// logger.Debugf("the raw is:%+v", string(body))
	return string(body), nil
}

func (m *RealAuthMgr) getReportResponse(urlValue url.Values, v interface{}) (*ReportResponse, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", m.loginOutUrl, bytes.NewBuffer(result))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range urlValue {
		req.Header.Set(k, v[0])
	}
	// logger.Debugf("header: %+v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// logger.Debugf("%s", body)

	responseData := &ReportResponse{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}

	if responseData.ErrCode != 0 {
		return nil, fmt.Errorf("%d -> %s", responseData.ErrCode, responseData.ErrMsg)
	}
	return responseData, nil
}

func (m *RealAuthMgr) getReportCheckResponse(code string, urlValue url.Values, v interface{}) (*ReportResponse, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", m.loginOutUrl+"/"+code, bytes.NewBuffer(result))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range urlValue {
		req.Header.Set(k, v[0])
	}
	fmt.Printf("header: %+v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("the raw is:%+v", string(body))

	responseData := &ReportResponse{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}

	if responseData.ErrCode != 0 {
		return nil, fmt.Errorf("%d -> %s", responseData.ErrCode, responseData.ErrMsg)
	}
	return responseData, nil
}

func (m *RealAuthMgr) getQueryResponse(urlValue url.Values, ai string) (string, error) {
	u := m.queryUrl + fmt.Sprintf("?ai=%s", ai)
	req, _ := http.NewRequest("GET", u, nil)
	// logger.Debugf("query url: %s", u)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range urlValue {
		logger.Debugf("%s -> %s", k, v[0])
		req.Header.Set(k, v[0])
	}
	// logger.Debugf("header: %+v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// logger.Debugf("%s", string(body))

	responseData := &Response{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return "", err
	}

	if responseData.ErrCode != 0 {
		return "", fmt.Errorf("%d -> %s", responseData.ErrCode, responseData.ErrMsg)
	}

	result := responseData.Data.Result
	if result.Status == 0 {
		return responseData.Data.Result.Pi, nil
	}

	switch result.Status {
	case 0:
		return responseData.Data.Result.Pi, nil
	case 1:
		return "", ErrNeedQuery
	case 2:
		return "", fmt.Errorf("")
	}
	return "", fmt.Errorf("result status error: %d", result.Status)
}

func (m *RealAuthMgr) getQueryCheckResponse(code string, urlValue url.Values, ai string) (string, error) {
	u := m.queryUrl + fmt.Sprintf("/%s?ai=%s", code, ai)
	req, _ := http.NewRequest("GET", u, nil)
	// logger.Debugf("query url: %s", u)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for k, v := range urlValue {
		// logger.Debugf("%s -> %s", k, v[0])
		req.Header.Set(k, v[0])
	}
	// logger.Debugf("header: %+v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// logger.Debugf("%s", string(body))

	responseData := &Response{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return "", err
	}

	if responseData.ErrCode != 0 {
		return "", fmt.Errorf("%d -> %s", responseData.ErrCode, responseData.ErrMsg)
	}

	result := responseData.Data.Result
	if result.Status == 0 {
		return responseData.Data.Result.Pi, nil
	}

	switch result.Status {
	case 0:
		return responseData.Data.Result.Pi, nil
	case 1:
		return "", ErrNeedQuery
	case 2:
		return "", fmt.Errorf("")
	}
	return "", fmt.Errorf("result status error: %d", result.Status)
}

func (m *RealAuthMgr) Check(uid int64, name, idNum string) (string, error) {
	ai := utils.Md5(utils.Int64ToString(uid))
	info := &RequestInfo{
		Ai:    ai,
		Name:  name,
		IdNum: idNum,
	}
	body, err := m.getEncryptData(info)
	if err != nil {
		return "", err
	}
	header, err := m.getHeader(nil, body)
	if err != nil {
		return "", err
	}
	pi, err := m.getResponse(header, body)
	if err != nil {
		return "", err
	}
	return pi, nil
}

func (m *RealAuthMgr) CheckTest(ai string, name, idNum, code string) (string, error) {
	info := &RequestInfo{
		Ai:    ai,
		Name:  name,
		IdNum: idNum,
	}
	body, err := m.getEncryptData(info)
	if err != nil {
		return "", err
	}
	header, err := m.getHeader(nil, body)
	if err != nil {
		return "", err
	}
	pi, err := m.getResponseCheck(code, header, body)
	if err != nil {
		return "", err
	}
	return pi, nil
}

func (m *RealAuthMgr) Query(uid int) (string, error) {
	ai := utils.Md5(utils.IntToString(uid))
	param := url.Values{}
	// 设置参数
	param.Add("ai", ai)

	header, err := m.getHeader(param, nil)
	if err != nil {
		return "", err
	}

	pi, err := m.getQueryResponse(header, ai)
	if err != nil {
		return "", err
	}
	return pi, nil
}

func (m *RealAuthMgr) QueryCheck(ai, code string) (string, error) {
	param := url.Values{}
	// 设置参数
	param.Add("ai", ai)

	header, err := m.getHeader(param, nil)
	if err != nil {
		return "", err
	}

	pi, err := m.getQueryCheckResponse(code, header, ai)
	if err != nil {
		return "", err
	}
	return pi, nil
}

func (m *RealAuthMgr) ReportLoginout(item ReportItem) (*ReportResponse, error) {
	report := ReportData{
		Collections: []ReportItem{item},
	}

	// logger.Debugf("%v", report)

	body, err := m.getEncryptData(report)
	if err != nil {
		return nil, err
	}

	header, err := m.getHeader(nil, body)
	if err != nil {
		return nil, err
	}

	res, err := m.getReportResponse(header, body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *RealAuthMgr) ReportLoginoutCheck(item ReportItem, code string) (*ReportResponse, error) {
	report := ReportData{
		Collections: []ReportItem{item},
	}

	// logger.Debugf("%v", report)

	body, err := m.getEncryptData(report)
	if err != nil {
		return nil, err
	}

	header, err := m.getHeader(nil, body)
	if err != nil {
		return nil, err
	}

	res, err := m.getReportCheckResponse(code, header, body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *RealAuthMgr) IdDecode(idcard string) (string, string, string, bool) {
	sex, birthday, address, isAdult := "男", "", "", false
	nowTime := time.Now()
	id := chinaid.IDCard(idcard)
	if result, err := id.Decode(); err == nil {
		if result.Sex == 0 {
			sex = "女"
		}
		startTime := nowTime.AddDate(-18, 0, 0)
		if startTime.Unix() > result.Birthday.Unix() {
			isAdult = true
		}
	} else {
		logger.Errorf("chinaid error idcard:%v, err:%v", idcard, err)
	}
	return sex, birthday, address, isAdult
}
