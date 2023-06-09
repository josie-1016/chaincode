package org

import (
	"encoding/json"
	"fmt"
	"log"
	"trustPlatform/constant"
	"trustPlatform/data"
	"trustPlatform/request"
	"trustPlatform/utils"

	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// ===================================================================================
// org合约初始化函数
// ===================================================================================
func Init(stub shim.ChaincodeStubInterface) {
	log.Println("Org init")
}

// ===================================================================================
// org模块入口函数
// ===================================================================================
func Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("TrustPlatformCC Org Invoke")
	function, args := stub.GetFunctionAndParameters()
	if err := utils.CheckInputNumber(1, args); err != nil {
		return shim.Error(err.Error())
	}
	if function == "/org/createOrgApply" {
		return createOrgApply(stub, args)
	} else if function == "/org/approveOrgApply" {
		return approveOrgApply(stub, args)
	} else if function == "/org/shareSecret" {
		return shareSecret(stub, args)
	} else if function == "/org/getSharedSecret" {
		return getSharedSecret(stub, args)
	} else if function == "/org/submitPartPK" {
		return submitPartPK(stub, args)
	} else if function == "/org/mixPartPK" {
		return mixPartPK(stub, args)
	} else if function == "/org/queryOrgApply" {
		return queryOrgApply(stub, args)
	} else if function == "/org/queryOrg" {
		return queryOrg(stub, args)
	} else if function == "/org/declareAttrApply" {
		return declareAttrApply(stub, args)
	} else if function == "/org/decryptThreshold" {
		return DecryptThreshold(stub, args)
	} else if function == "/org/approveandsubmitThresholdPartOSK" {
		return approveandsubmitThreholdPartOSK(stub, args)
	} else if function == "/org/queryOrgMen" {
		return queryOrgMen(stub, args)
	} else if function == "/org/ThresholdFileApply" {
		return ThreholdFileApply(stub, args)
	} else if function == "/org/queryThresholdApply" {
		return queryThresholdOrgApply(stub, args)
	} else if function == "/org/ThresholdOrgApply" {
		return ThresholdOrgApply(stub, args)
	} else if function == "/org/approveThresholdOrgApply" {
		return approveThresholdOrgApply(stub, args)
	} else if function == "/org/mixPartOSKForThresholdPub" {
		return mixPartOSKForThreholdPub(stub, args)
	} else if function == "/org/thresholdEncrypt" {
		return EncryptThreshold(stub, args)
	} else if function == "/org/queryThresholdFileApply" {
		return QueryThresholdFileApply(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"/org/createOrgApply\" \"/org/approveOrgApply\" " +
		"\"/org/shareSecret\" \"/org/getSharedSecret\" \"/org/submitPartPK\" \"/org/mixPartPK\"" +
		" \"/org/queryOrg\" ")
}

// ===================================================================================
// 申请声明组织新属性
// ===================================================================================
func declareAttrApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("declareAttrApply: apply new org attr start")

	// 反序列化请求，验签
	var requestStr = args[0]
	log.Println(requestStr)
	applyRequest := new(request.DeclareOrgAttrApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), applyRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, applyRequest.Uid, applyRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := applyRequest.Uid
	orgId := applyRequest.OrgId
	attrName := applyRequest.AttrName

	// 检查attr name
	if err := utils.CheckAttr(attrName, orgId); err != nil {
		return shim.Error(err.Error())
	}
	// org 是否存在
	org, err := data.QueryOrgByOid(orgId, stub)
	// stub.GetState("ID:org4"）
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	// 检查attr name是否存在
	attr, _ := data.QueryAttr(attrName, stub)
	// stub.GetState("ATTR:org4:test")
	if attr != nil {
		log.Println("already has this attr:" + attrName)
		return shim.Error("already has this attr:" + attrName)
	}
	// 检查是否存在对相同attrName的Pending apply请求存在，即有效请求
	apply, err := data.QueryDeclareOrgAttrApply(attrName, data.Pending, stub)
	// {"selector":{"docType":"orgApply","attrName":"org4:test","type":1,"status":{"$lt":"2"}}}
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if apply != nil {
		log.Println("already has active createOrgApply for orgId:" + orgId)
		return shim.Error("already has active createOrgApply for orgId:" + orgId)
	}

	// 保存
	newApply := data.NewOrgApply(orgId, uid, org.UidSet, org.T, org.N, constant.DeclareAttr, attrName)
	if err := data.SaveOrgApply(newApply, stub); err != nil {
		//stub.PutState("ORG_APPLY:org4:zhuhai1")
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// 查询组织申请
// ===================================================================================
func queryOrgApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var requestStr = args[0]
	log.Println(requestStr)
	getOrgApplyRequest := new(request.GetOrgApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), getOrgApplyRequest); err != nil {
		return shim.Error(err.Error())
	}

	var apply *data.OrgApply
	var err error
	if getOrgApplyRequest.Type == constant.CreateOrg {
		apply, err = data.QueryCreateOrgApply(getOrgApplyRequest.OrgId, getOrgApplyRequest.Status, stub)
		// {"selector":{"docType":"orgApply","orgId":"org4","type":0,"status":{"$lt":"4"}}}
		log.Println(apply)
	} else {
		apply, err = data.QueryDeclareOrgAttrApply(getOrgApplyRequest.AttrName, getOrgApplyRequest.Status, stub)
	}
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	bytes, err := json.Marshal(apply)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

// ===================================================================================
// 查询组织门限公钥
// ===================================================================================
func queryOrgMen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	oid := args[0]
	oid1 := oid[1 : len(oid)-1]
	org, err := data.QueryOrgByOid(oid1, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	bytes := org.ThreholdPub
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	log.Println("menpub is:", []byte(bytes))
	log.Println("len is:", len([]byte(bytes)))
	return shim.Success([]byte(bytes))
}
func queryOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var requestStr = args[0]
	orgRequest := new(request.GetOrgRequest)
	if err := json.Unmarshal([]byte(requestStr), orgRequest); err != nil {
		return shim.Error(err.Error())
	}
	oid := orgRequest.OrgId
	org, err := data.QueryOrgByOid(oid, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	bytes, err := json.Marshal(org)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

// ===================================================================================
// 整合opk_i或是apk_i
// ===================================================================================
func mixPartPK(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 反序列化请求，验签
	var requestStr = args[0]
	mixRequest := new(request.MixPartPKRequest)
	if err := json.Unmarshal([]byte(requestStr), mixRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, mixRequest.Uid, mixRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	//uid := mixRequest.Id
	orgId := mixRequest.OrgId
	attrName := mixRequest.AttrName
	sceneType := mixRequest.Type

	if sceneType == request.OrgInit {
		apply, err := data.QueryCreateOrgApply(orgId, data.PendingShare, stub)
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		if apply == nil {
			log.Println("no active createOrgApply for orgId:" + orgId)
			return shim.Error("no active createOrgApply for orgId:" + orgId)
		}

		if len(apply.OpkMap) != apply.N {
			log.Println("createOrgApply not ready for mixing with orgId:" + orgId)
			return shim.Error("createOrgApply not ready for mixing with orgId:" + orgId)
		}

		// 调用属性密码合约进行整合得到最终opk
		userNames := make([]string, 0, apply.N)
		partPkList := make([]string, 0, apply.N)
		for key := range apply.UidMap {
			userNames = append(userNames, key)
			partPkList = append(partPkList, apply.OpkMap[key])
		}
		mixOPKRequest := &request.GenerateOPKRequest{
			UserNames:  userNames,
			PartPkList: partPkList,
			N:          apply.N,
			T:          apply.T,
		}
		bytes, err := json.Marshal(mixOPKRequest)
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		mixedPkBytes, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/org/generateOPK", string(bytes)})
		// mixedPkBytes1, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/org/generateMen", string(bytes)})
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}

		apply.Status = data.Success
		if err := data.SaveOrgApply(apply, stub); err != nil {
			return shim.Error(err.Error())
		}
		uidSet := make([]string, 0, apply.N)
		for u := range apply.UidMap {
			uidSet = append(uidSet, u)
		}
		org := data.NewOrg(orgId, apply.T, apply.N, uidSet, string(mixedPkBytes), []byte{})
		if err = data.SaveOrg(org, stub); err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	} else if sceneType == request.OrgAttr {
		apply, err := data.QueryDeclareOrgAttrApply(attrName, data.PendingShare, stub)
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		if apply == nil {
			log.Println("no active declareAttrApply for attr:" + attrName)
			return shim.Error("no active declareAttrApply for attr:" + attrName)
		}

		if len(apply.OpkMap) != apply.N {
			log.Println("declareAttrApply not ready for mixing with attrName:" + attrName)
			return shim.Error("declareAttrApply not ready for mixing with attrName:" + attrName)
		}

		// 调用属性密码合约进行整合得到最终opk
		userNames := make([]string, 0, apply.N)
		partPkList := make([]string, 0, apply.N)
		for key := range apply.UidMap {
			userNames = append(userNames, key)
			partPkList = append(partPkList, apply.OpkMap[key])
		}
		mixAPKRequest := &request.GenerateAPKRequest{
			UserNames:  userNames,
			PartPkList: partPkList,
			N:          apply.N,
			T:          apply.T,
			AttrName:   attrName,
		}
		bytes, err := json.Marshal(mixAPKRequest)
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		mixedPkBytes, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/org/generateAPK", string(bytes)})
		log.Println(string(mixedPkBytes))
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}

		apply.Status = data.Success
		if err := data.SaveOrgApply(apply, stub); err != nil {
			return shim.Error(err.Error())
		}

		org, err := data.QueryOrgByOid(orgId, stub)
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		newAttr := data.NewAttr(orgId, attrName, string(mixedPkBytes))
		if err = data.SaveOrgAttr(org, newAttr, stub); err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	} else {
		return shim.Error("type invalid")
	}
}

// ===================================================================================
// 整合opk_i生成门限私钥
// ===================================================================================
func DecryptThreshold(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 反序列化请求，验签
	var requestStr = args[0]
	mixRequest := new(request.ThreholdMixPartPKRequest)
	if err := json.Unmarshal([]byte(requestStr), mixRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, mixRequest.Uid, mixRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := mixRequest.Uid
	orgId := mixRequest.OrgId
	fileName := mixRequest.FileName

	apply, err := data.QueryThreholdFileApply(orgId, uid, fileName, data.PendingShare, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	if apply == nil {
		log.Println("no active ThreholdFileApply for uid:" + uid)
		return shim.Error("no active ThreholdFileApply for uid:" + uid)
	}

	if len(apply.OpkMap) < apply.T {
		log.Println("ThreholdFileApply not ready for mixing with uId:" + uid)
		return shim.Error("ThreholdFileApply not ready for mixing with uId:" + uid)
	}

	// 调用属性密码合约进行整合得到最终门限私钥
	userNames := make([]string, 0, apply.T)
	partPkList := make([]string, 0, apply.T)
	for key, flag := range apply.UidMap {
		if flag == true {
			userNames = append(userNames, key)
			partPkList = append(partPkList, apply.OpkMap[key])
		}
	}
	mixOPKRequest := &request.GenerateOPKRequest{
		UserNames:  userNames,
		PartPkList: partPkList,
		N:          apply.N,
		T:          apply.T,
	}

	log.Println(mixOPKRequest)
	log.Println("recover pri start")
	bytes, err := json.Marshal(mixOPKRequest)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	mixedPkBytes, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/org/generateMenPriv", string(bytes)})
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	ThresholdDecRequest := &request.ThreholdDecryptRequest{
		Cipher:        []byte(mixRequest.CipherContent),
		ThresholdPriv: mixedPkBytes,
	}

	log.Println(ThresholdDecRequest)
	log.Println("decrypt start")
	bytes1, err := json.Marshal(ThresholdDecRequest)
	log.Println(string(bytes1))
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	mixedPkBytes1, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/common/decryptThreshold", string(bytes1)})
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	log.Println("plaintext is:")
	log.Println(string(mixedPkBytes1))
	return shim.Success([]byte(mixedPkBytes1))
}

// ===================================================================================
// 提交opk_i或是apk_i
// ===================================================================================
func submitPartPK(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 反序列化请求，验签
	var requestStr = args[0]
	partPKRequest := new(request.SubmitPartPKRequest)
	if err := json.Unmarshal([]byte(requestStr), partPKRequest); err != nil {
		return shim.Error(err.Error())
	}

	/* if err := preCheckRequest(requestStr, partPKRequest.Uid, partPKRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := partPKRequest.Uid
	orgId := partPKRequest.OrgId
	attrName := partPKRequest.AttrName
	sceneType := partPKRequest.Type
	parkPk := partPKRequest.PartPK

	var apply *data.OrgApply
	var err error
	if sceneType == request.OrgInit {
		apply, err = data.QueryCreateOrgApply(orgId, data.PendingShare, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		if apply == nil {
			log.Println("no active createOrgApply for orgId:" + orgId)
			return shim.Error("no active createOrgApply for orgId:" + orgId)
		}
	} else if sceneType == request.OrgAttr {
		apply, err = data.QueryDeclareOrgAttrApply(attrName, data.PendingShare, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		if apply == nil {
			log.Println("no active declareAttrApply for attr:" + attrName)
			return shim.Error("no active declareAttrApply for attr:" + attrName)
		}
	} else {
		return shim.Error("type invalid")
	}

	if _, ok := apply.UidMap[uid]; !ok {
		log.Println(uid + "not active in orgId:" + orgId)
		return shim.Error(uid + "not active in orgId:" + orgId)
	}
	apply.OpkMap[uid] = parkPk
	if err := data.SaveOrgApply(apply, stub); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ===================================================================================
// 用户同意解密文件并提交osk_i
// ===================================================================================
func approveandsubmitThreholdPartOSK(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 反序列化请求，验签
	var requestStr = args[0]
	partPKRequest := new(request.SubmitThreholdPartPKRequest)
	if err := json.Unmarshal([]byte(requestStr), partPKRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, partPKRequest.Uid, partPKRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := partPKRequest.Uid
	orgId := partPKRequest.OrgId
	fileName := partPKRequest.FileName
	partPk := partPKRequest.PartPK
	fromUid := partPKRequest.FromUId

	var apply *data.OrgApply
	var err error

	apply, err = data.QueryThreholdFileApply(orgId, fromUid, fileName, data.PendingShare, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if apply == nil {
		log.Println("no active ThreholdFileApply for orgId:" + fileName)
		return shim.Error("no active ThreholdFileApply for orgId:" + fileName)
	}

	apply.UidMap[uid] = true
	apply.OpkMap[uid] = partPk

	num := 0
	for _, val := range apply.UidMap {
		if val {
			num++
		}
	}
	if num >= apply.T {
		apply.Status = data.PendingShare
	}
	if err := data.SaveOrgApply(apply, stub); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ===================================================================================
// 获取别人给自己的secret
// ===================================================================================
func getSharedSecret(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("get secret")
	var requestStr = args[0]
	getShareRequest := new(request.GetShareRequest)
	if err := json.Unmarshal([]byte(requestStr), getShareRequest); err != nil {
		return shim.Error(err.Error())
	}
	orgId := getShareRequest.OrgId
	toUid := getShareRequest.ToUid
	fromUid := getShareRequest.FromUid
	attrName := getShareRequest.AttrName

	var apply *data.OrgApply
	var err error
	if attrName == "" {
		apply, err = data.QueryCreateOrgApply(orgId, data.Pending, stub)
	} else {
		apply, err = data.QueryDeclareOrgAttrApply(attrName, data.Pending, stub)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	if toUid == "" {
		bytes, err := json.Marshal(apply.ShareMap[fromUid])
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(bytes)
	}
	return shim.Success([]byte(apply.ShareMap[fromUid][toUid]))
}

// ===================================================================================
// 申请建立新组织
// ===================================================================================
func createOrgApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("createOrgApply new org start")

	// 反序列化请求，验签
	var requestStr = args[0]
	applyRequest := new(request.CreateOrgApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), applyRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, applyRequest.Uid, applyRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := applyRequest.Uid
	orgId := applyRequest.OrgId
	t := applyRequest.T
	n := applyRequest.N
	uidList := applyRequest.UidList

	// 检查传入参数
	if n != len(uidList) || t > n || n == 1 {
		log.Printf("input params error, t:%d, n:%d\n", t, n)
		return shim.Error("input params error")
	}
	// 检查org id
	if err := utils.CheckId(orgId); err != nil {
		return shim.Error(err.Error())
	}
	// org id是否存在
	if utils.ExistId(orgId, stub) {
		log.Printf("orgId %s already exists\n", orgId)
		return shim.Error("orgId already exists")
	}

	// 检查是否存在对同一个orgId的Pending apply请求存在，即有效请求
	apply, err := data.QueryCreateOrgApply(orgId, data.Pending, stub)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if apply != nil {
		log.Println("already has active createOrgApply for orgId:" + orgId)
		return shim.Error("already has active createOrgApply for orgId:" + orgId)
	}

	// 保存
	newApply := data.NewOrgApply(orgId, uid, uidList, t, n, constant.CreateOrg, "")
	if err := data.SaveOrgApply(newApply, stub); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// 申请建立新门限组织
// ===================================================================================
func ThresholdOrgApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("Threshold apply start")

	// 反序列化请求，验签
	var requestStr = args[0]
	applyRequest := new(request.ThresholdOrgApply)
	if err := json.Unmarshal([]byte(requestStr), applyRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, applyRequest.Uid, applyRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := applyRequest.Uid
	orgId := applyRequest.OrgId

	// org 是否存在
	org, err := data.QueryOrgByOid(orgId, stub)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	//是否存在对同一个门限组织的请求

	apply, err := data.QueryThresholdOrgApply(orgId, data.Pending, stub)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if apply != nil {
		log.Println("already has active ThresholdOrgApply for orgId:" + orgId)
		return shim.Error("already has active ThresholdOrgApply for orgId:" + orgId)
	}

	// 保存
	newApply := data.NewOrgApply(orgId, uid, org.UidSet, org.T, org.N, 3, "")
	newApply.UidMap[uid] = false
	if err := data.SaveOrgApply(newApply, stub); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// 用户确认加入门限组织并提交osk_i
// ===================================================================================
func approveThresholdOrgApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("user approve threshold org apply")

	var requestStr = args[0]
	approveRequest := new(request.ApproveThresholdOrgApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), approveRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, approveRequest.Uid, approveRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	orgId := approveRequest.OrgId
	uid := approveRequest.Uid
	partpk := approveRequest.PartSK

	var apply *data.OrgApply
	var err error

	apply, err = data.QueryThresholdOrgApply(orgId, data.PendingApprove, stub)

	if err != nil {
		return shim.Error(err.Error())
	}
	if apply == nil {
		return shim.Error("no PendingApprove createThresholdOrgApply for org: " + orgId)
	}
	if val, ok := apply.UidMap[uid]; ok {
		if val {
			return shim.Success([]byte("already approve"))
		}
	} else {
		return shim.Error("authority limits")
	}

	apply.UidMap[uid] = true
	apply.OpkMap[uid] = partpk
	allApprove := true
	for _, val := range apply.UidMap {
		if !val {
			allApprove = false
			break
		}
	}
	if allApprove {
		apply.Status = data.PendingShare
	}
	if err := data.SaveOrgApply(apply, stub); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// 生成门限公钥
// ===================================================================================
func mixPartOSKForThreholdPub(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 反序列化请求，验签
	log.Println("mixpartsk for threshold pub key start:")
	var requestStr = args[0]
	mixRequest := new(request.ThreholdMixRequest)
	if err := json.Unmarshal([]byte(requestStr), mixRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, mixRequest.Uid, mixRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	orgId := mixRequest.OrgId

	apply, err := data.QueryThresholdOrgApply(orgId, data.PendingShare, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	if apply == nil {
		log.Println("no active ThreholdOrgApply for orgid:" + orgId)
		return shim.Error("no active ThreholdFileApply for orgid:" + orgId)
	}

	if len(apply.OpkMap) != apply.N {
		log.Println("createThresholdOrgApply not ready for mixing with orgId:" + orgId)
		return shim.Error("createThresholdOrgApply not ready for mixing with orgId:" + orgId)
	}

	// 调用属性密码合约进行整合得到最终门限私钥
	userNames := make([]string, 0, apply.N)
	partPkList := make([]string, 0, apply.N)
	for key, flag := range apply.UidMap {
		if flag == true {
			userNames = append(userNames, key)
			partPkList = append(partPkList, apply.OpkMap[key])
		}
	}
	mixOPKRequest := &request.GenerateOPKRequest{
		UserNames:  userNames,
		PartPkList: partPkList,
		N:          apply.N,
		T:          apply.T,
	}
	bytes, err := json.Marshal(mixOPKRequest)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	mixedPkBytes, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/org/generateMenPub", string(bytes)})
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	log.Println("threshold pub key is :")
	log.Println(string(mixedPkBytes))
	org, err := data.QueryOrgByOid(orgId, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	org.ThreholdPub = mixedPkBytes
	if err = data.SaveOrg(org, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// ===================================================================================
// 申请下载门限文件
// ===================================================================================
func ThreholdFileApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("ThreholdFileApply start")

	// 反序列化请求，验签
	var requestStr = args[0]
	applyRequest := new(request.ApplyThreholdFileRequest)
	if err := json.Unmarshal([]byte(requestStr), applyRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, applyRequest.Uid, applyRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	uid := applyRequest.Uid
	orgId := applyRequest.OrgId
	fileName := applyRequest.FileName

	// org id是否存在
	org, err := data.QueryOrgByOid(orgId, stub)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	// 检查是否存在对同一个orgId的Pending apply请求存在，即有效请求
	apply, err := data.QueryThreholdFileApply(orgId, uid, fileName, data.Pending, stub)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if apply != nil {
		log.Println("already has active ThreholdFileApply for fileName:" + fileName)
		return shim.Error("already has active ThreholdFileApply for fileName:" + fileName)
	}

	// 保存
	newApply := data.NewThresholdOrgApply(orgId, uid, org.UidSet, org.T, org.N, 2, fileName)
	if err := data.SaveOrgApply(newApply, stub); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// 检查请求参数并验签
// ===================================================================================
func preCheckRequest(requestStr string, uid, sign string, stub shim.ChaincodeStubInterface) error {
	requestJson, err := utils.GetRequestParamJson([]byte(requestStr))
	if err != nil {
		return err
	}
	requestUser, err := data.QueryUserByUid(uid, stub)
	if err != nil {
		return err
	}
	if requestUser == nil {
		log.Println("don't have requestUser with uid " + uid)
		return ecode.Error(ecode.RequestErr, "don't have this requestUser")
	}
	if err = utils.VerifySign(string(requestJson), requestUser.PublicKey, sign, uid); err != nil {
		return err
	}
	return nil
}

// ===================================================================================
// 用户确认加入组织/声明新属性
// ===================================================================================
func approveOrgApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("user approve org apply")

	var requestStr = args[0]
	approveRequest := new(request.ApproveOrgApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), approveRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, approveRequest.Uid, approveRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	orgId := approveRequest.OrgId
	uid := approveRequest.Uid
	attrName := approveRequest.AttrName

	var apply *data.OrgApply
	var err error
	if attrName == "" {
		apply, err = data.QueryCreateOrgApply(orgId, data.PendingApprove, stub)
	} else {
		apply, err = data.QueryDeclareOrgAttrApply(attrName, data.PendingApprove, stub)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	if apply == nil {
		return shim.Error("no PendingApprove createOrgApply for org: " + orgId)
	}
	if val, ok := apply.UidMap[uid]; ok {
		if val {
			return shim.Success([]byte("already approve"))
		}
	} else {
		return shim.Error("authority limits")
	}

	apply.UidMap[uid] = true
	allApprove := true
	for _, val := range apply.UidMap {
		if !val {
			allApprove = false
			break
		}
	}
	if allApprove {
		apply.Status = data.PendingShare
	}
	if err := data.SaveOrgApply(apply, stub); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===================================================================================
// 用户确认门限文件解密申请
// ===================================================================================
/*  func approveThreholdfileApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("user approve Threhold apply")

	var requestStr = args[0]
	approveRequest := new(request.ApproveThreholdApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), approveRequest); err != nil {
		return shim.Error(err.Error())
	}
	 if err := preCheckRequest(requestStr, approveRequest.Uid, approveRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	}

	uid := approveRequest.Uid
	fileName := approveRequest.FileName
	orgid := approveRequest.OrgId
	fromUid := approveRequest.FromUid

	var apply *data.OrgApply
	var err error

	apply, err = data.QueryThreholdFileApply(orgid, fromUid, fileName, data.PendingApprove, stub)

	if err != nil {
		return shim.Error(err.Error())
	}
	if apply == nil {
		return shim.Error("no PendingApprove ThreholdApply for file: " + fileName)
	}
	if val, ok := apply.UidMap[uid]; ok {
		if val {
			return shim.Success([]byte("already approve"))
		}
	} else {
		return shim.Error("authority limits")
	}

	apply.UidMap[uid] = true
	allApprove := true
	for _, val := range apply.UidMap {
		if !val {
			allApprove = false
			break
		}
	}
	if allApprove {
		apply.Status = data.PendingShare
	}
	if err := data.SaveOrgApply(apply, stub); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}  */

// ===================================================================================
// 用户交换秘密
// ===================================================================================
func shareSecret(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("user share secret with others")

	var requestStr = args[0]
	shareRequest := new(request.SubmitOrgShareRequest)
	if err := json.Unmarshal([]byte(requestStr), shareRequest); err != nil {
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, shareRequest.Uid, shareRequest.Sign, stub); err != nil {
		return shim.Error(err.Error())
	} */

	orgId := shareRequest.OrgId
	attrName := shareRequest.AttrName
	fromUid := shareRequest.Uid
	toUid := shareRequest.ToUid
	applyType := shareRequest.Type
	share := shareRequest.Share

	var apply *data.OrgApply
	var err error
	if applyType == request.OrgInit {
		// 查询 + 写入
		apply, err = data.QueryCreateOrgApply(orgId, data.Pending, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		if apply == nil {
			log.Println("no Pending createOrgApply for org: " + orgId)
			return shim.Error("no Pending createOrgApply for org: " + orgId)
		}
	} else if applyType == request.OrgAttr {
		// 查询 + 写入
		apply, err = data.QueryDeclareOrgAttrApply(attrName, data.Pending, stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		if apply == nil {
			log.Println("no Pending declareOrgAttrApply for attr: " + attrName)
			return shim.Error("no Pending declareOrgAttrApply for attr: " + attrName)
		}
	} else {
		log.Println("not match type")
		return shim.Error("not match type")
	}
	apply.ShareMap[toUid][fromUid] = share
	if err := data.SaveOrgApply(apply, stub); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ===================================================================================
// 查询是否有组织门限申请
// ===================================================================================
func queryThresholdOrgApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var requestStr = args[0]
	getOrgApplyRequest := new(request.GetThreholdApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), getOrgApplyRequest); err != nil {
		return shim.Error(err.Error())
	}

	var apply *data.OrgApply
	var err error
	apply, err = data.QueryThreholdFileApply(getOrgApplyRequest.OrgId, getOrgApplyRequest.FromUid, getOrgApplyRequest.FileName, getOrgApplyRequest.Status, stub)

	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	bytes, err := json.Marshal(apply)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

// 门限加密文件
func EncryptThreshold(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("Encrypt threshold file start!")
	var requestStr = args[0]
	approveRequest := new(request.ThreholdEncryptApplyRequest)
	if err := json.Unmarshal([]byte(requestStr), approveRequest); err != nil {
		return shim.Error(err.Error())
	}
	org, err := data.QueryOrgByOid(approveRequest.OrgId, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	bytes := org.ThreholdPub
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	log.Println("menpub is:", []byte(bytes))
	log.Println("len is:", len([]byte(bytes)))

	ThresholdEncRequest := &request.ThresholdEncryptRequest{
		PlainContent: approveRequest.PlainContent,
		PubKey:       bytes,
	}
	log.Println([]byte(ThresholdEncRequest.PubKey))
	bytes1, err := json.Marshal(ThresholdEncRequest)
	log.Println(string(bytes1))
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	temp := new(request.ThresholdEncryptRequest)
	json.Unmarshal(bytes1, temp)
	log.Println("11111111111111111111111")
	log.Println(temp.PubKey)
	log.Println([]byte(temp.PubKey))
	mixedPkBytes, err := utils.InvokeOtherChaincode("dabe", stub, []string{"/common/encryptThreshold", string(bytes1)})
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(mixedPkBytes))
}

// 组织内搜索对门限文件得请求
func QueryThresholdFileApply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var requestStr = args[0]
	partPKRequest := new(request.QueryThreholdFileApply)
	if err := json.Unmarshal([]byte(requestStr), partPKRequest); err != nil {
		return shim.Error(err.Error())
	}
	fromUid := partPKRequest.FromUId
	fileName := partPKRequest.FileName
	orgId := partPKRequest.OrgId

	var apply *data.OrgApply
	var err error

	apply, err = data.QueryThreholdFileApply(orgId, fromUid, fileName, data.PendingShare, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if apply == nil {
		s := fmt.Sprintf("user %s didnot apply org:%s for file:%s:\n", fromUid, orgId, fileName)
		log.Printf(s)
		return shim.Error(s)
	}
	return shim.Success([]byte(apply.CreateTime))
}
