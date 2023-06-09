package request

import (
	"trustPlatform/constant"
	"trustPlatform/data"
)

type BaseRequest struct {
	Uid       string `json:"uid"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

// 用户初始化
type UserInitRequest struct {
	BaseRequest
	PublicKey string `json:"publicKey"`
	UPK       string `json:"upk"`
}

// 用户声明新属性
type UserAnnounceAttrRequest struct {
	BaseRequest
	AttrName string `json:"attrName"`
	APK      string `json:"apk"`
}

// 用户申请生成门限组织
type ThresholdOrgApply struct {
	BaseRequest
	OrgId string `json:"orgId"`
}

// 发起生成新组织请求
type CreateOrgApplyRequest struct {
	BaseRequest
	T       int      `json:"t"`
	N       int      `json:"n"`
	UidList []string `json:"uidList"`
	OrgId   string   `json:"orgId"`
}

// 发起生成组织属性请求
type DeclareOrgAttrApplyRequest struct {
	BaseRequest
	OrgId    string `json:"orgId"`
	AttrName string `json:"attrName"`
}

// 发起门限解密申请
type ApplyThreholdFileRequest struct {
	BaseRequest
	OrgId    string `json:"orgId"`
	FileName string `json:"fileName"`
}

// 同意加入新组织/同意声明属性
type ApproveOrgApplyRequest struct {
	BaseRequest
	OrgId    string `json:"orgId"`
	AttrName string `json:"attrName"`
}

// 同意加入门限组织
type ApproveThresholdOrgApplyRequest struct {
	BaseRequest
	OrgId  string `json:"orgId"`
	PartSK string `json:"partSk"`
}

// 同意用户申请门限解密
type ApproveThreholdApplyRequest struct {
	BaseRequest
	OrgId    string `json:"orgId"`
	FileName string `json:"fileName"`
	FromUid  string `json:"fromUid"`
}

type SceneType int8

const (
	OrgInit SceneType = iota //0
	OrgAttr
)

// 交换秘密
type SubmitOrgShareRequest struct {
	BaseRequest
	OrgId    string    `json:"orgId"`
	AttrName string    `json:"attrName"`
	Type     SceneType `json:"type"`
	ToUid    string    `json:"toUid"`
	Share    string    `json:"share"`
}

// 提交部分公钥参数
type SubmitPartPKRequest struct {
	BaseRequest
	OrgId    string    `json:"orgId"`
	AttrName string    `json:"attrName"`
	Type     SceneType `json:"type"`
	PartPK   string    `json:"partPk"`
}

// 提交门限公钥参数
type SubmitThreholdPartPKRequest struct {
	BaseRequest
	OrgId    string `json:"orgId"`
	FileName string `json:"fileName"`
	PartPK   string `json:"partPk"`
	FromUId  string `json:"fromUid"`
}

type QueryThreholdFileApply struct {
	BaseRequest
	OrgId    string `json:"orgId"`
	FileName string `json:"fileName"`
	FromUId  string `json:"fromUid"`
}

// 申请属性参数
type ApplyAttrRequest struct {
	BaseRequest
	ToUid    string `json:"toUid"`
	ToOrgId  string `json:"toOrgId"`
	IsPublic bool   `json:"isPublic"`
	AttrName string `json:"attrName"`
	Remark   string `json:"remark"`
}

// 审批属性申请
type ApproveAttrApplyRequest struct {
	BaseRequest
	FromUid  string `json:"fromUid"`
	ToOrgId  string `json:"toOrgId"`
	AttrName string `json:"attrName"`
	Remark   string `json:"remark"`
	Secret   string `json:"secret"`
	Agree    bool   `json:"agree"`
}

// 上传密文申请
type ShareMessageRequest struct {
	BaseRequest
	Tags []string `json:"tags"`
	// 加密内容
	Content  string `json:"content"`
	FileName string `json:"fileName"`
	Ip       string `json:"ip"`
	Location string `json:"location"`
	Policy   string `json:"policy"`
}

// 上传门限密文
type ThreholdShareMessageRequest struct {
	BaseRequest
	FileName string `json:"fileName"`
	OrgName  string `json:"orgName"`
}

// 整合请求
type MixPartPKRequest struct {
	BaseRequest
	OrgId    string    `json:"orgId"`
	AttrName string    `json:"attrName"`
	Type     SceneType `json:"type"`
}

// 整合请求
type ThreholdMixPartPKRequest struct {
	BaseRequest
	OrgId         string `json:"orgId"`
	FileName      string `json:"fileName"`
	CipherContent string `json:"cipherContent"`
}

// 整合请求
type ThreholdMixRequest struct {
	BaseRequest
	OrgId string `json:"orgId"`
}

// get share
type GetShareRequest struct {
	OrgId    string `json:"orgId"`
	AttrName string `json:"attrName"`
	ToUid    string `json:"toUid"`
	FromUid  string `json:"fromUid"`
}

// get user，任选一
type GetUserRequest struct {
	Uid       string `json:"uid"`
	PublicKey string `json:"publicKey"`
}

// 查询属性申请
type GetAttrApplyRequest struct {
	ToUid   string                   `json:"toUid"`
	ToOrgId string                   `json:"toOrgId"`
	FromUid string                   `json:"fromUid"`
	Status  constant.AttrApplyStatus `json:"status"`
}

// 查询分享信息
type GetSharedMessageRequest struct {
	FromUid  string `json:"fromUid"`
	Tag      string `json:"tag"`
	PageSize int    `json:"pageSize"`
	Bookmark string `json:"bookmark"`
}

// 查询分享信息
type GetThresholdSharedMessageRequest struct {
	OrgName  string `json:"orgName"`
	FileName string `json:"fileName"`
}

// 查询分享信息
type GetOrgApplyRequest struct {
	OrgId    string                `json:"orgId"`
	AttrName string                `json:"attrName"`
	Status   data.ApplyStatus      `json:"status"`
	Type     constant.OrgApplyType `json:"type"`
}

// 查询分享信息
type GetThreholdApplyRequest struct {
	OrgId    string           `json:"orgId"`
	FileName string           `json:"fileName"`
	Status   data.ApplyStatus `json:"status"`
	FromUid  string           `json:"fromUid"`
}

// 查询分享信息
type ThreholdEncryptApplyRequest struct {
	PlainContent string `json:"plainContent"`
	OrgId        string `json:"orgId"`
}
