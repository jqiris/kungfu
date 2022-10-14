package realauth

import "fmt"

var (
	ErrNeedQuery  = fmt.Errorf("need query result")
	ErrAuthFailed = fmt.Errorf("auth failed")
)

type RequestInfo struct {
	Ai    string `json:"ai"`
	Name  string `json:"name"`
	IdNum string `json:"idNum"`
}

type RequestBody struct {
	Data string `json:"data"`
}

// ReportItem 上下线上报的项目
type ReportItem struct {
	No int    `json:"no"` // 批量模式中的索引
	Si string `json:"si"` // 游戏内部会话标识
	Bt int    `json:"bt"` // 用户行为类型 0: 下线 1: 上线
	Ot int64  `json:"ot"` // 行为发生时间戳，秒
	Ct int    `json:"ct"` // 上报类型 0: 已认证通过类型 2:游客用户
	Di string `json:"di"` // 设备标识 由游戏运营单位生成，游客用户下必填
	Pi string `json:"pi"` // 已通过实名认证用户的唯一标识，已认证通过用户必填
}

// ReportData 上报的数据
type ReportData struct {
	Collections []ReportItem `json:"collections"`
}

// ReportResponse 上报上下线返回的数据
type ReportResponse struct {
	ErrCode int                `json:"errcode"`
	ErrMsg  string             `json:"errmsg"`
	Data    ReportResponseData `json:"data"`
}

type ReportResponseData struct {
	Result []ReportResponseData `json:"result"`
}

type ReportResultData struct {
	No      int    `json:"no"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// ResponseData  check 和 query 返回数据
type ResponseData struct {
	Result ResultData `json:"result"`
}

type ResultData struct {
	Status int    `json:"status"`
	Pi     string `json:"pi"`
}

type Response struct {
	ErrCode int          `json:"errcode"`
	ErrMsg  string       `json:"errmsg"`
	Data    ResponseData `json:"data"`
}
