package common

import (
	"encoding/json"
	"log"
	"strings"
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
// common初始化函数
// ===================================================================================
func Init(stub shim.ChaincodeStubInterface) {
	log.Println("Common init")
}

// ===================================================================================
// common模块入口函数
// ===================================================================================
func Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("TrustPlatformCC Common Invoke")
	function, args := stub.GetFunctionAndParameters()
	if err := utils.CheckInputNumber(1, args); err != nil {
		return shim.Error(err.Error())
	}
	if strings.HasPrefix(function, "/common/getAttr") {
		return getAttr(stub, args)
	} else if strings.HasPrefix(function, "/common/shareMessage") {
		return shareMessage(stub, args)
	} else if strings.HasPrefix(function, "/common/shareThreholdMessage") {
		return shareMessageThreshold(stub, args)
	} else if strings.HasPrefix(function, "/common/getMessage") {
		return getMessage(stub, args)
	} else if strings.HasPrefix(function, "/common/getThresholdMessage") {
		return getThreholdMessage(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"/common/getAttr\" ")
}

func getAttr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	attrName := args[0]
	attr, err := data.QueryAttrBytes(attrName, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	log.Println(string(attr))
	return shim.Success(attr)
}

// ===================================================================================
// 分享信息
// ===================================================================================
func shareMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("enter shareMessage")
	// 反序列化请求，验签
	var requestStr = args[0]
	log.Println(requestStr)
	shareRequest := new(request.ShareMessageRequest)
	if err := json.Unmarshal([]byte(requestStr), shareRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	/* if err := preCheckRequest(requestStr, shareRequest.Uid, shareRequest.Sign, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	} */
	if len(shareRequest.Tags) > 10 {
		return shim.Error("too much tags")
	}

	var message *data.SharedMessage
	if shareRequest.Org == "" {
		message = data.NewSharedMessage(shareRequest.Uid, shareRequest.Content, shareRequest.Tags, shareRequest.Timestamp, shareRequest.FileName, shareRequest.Ip, shareRequest.Location, shareRequest.Policy)
	} else {
		log.Println("Saved by Org")
		message = data.NewSharedMessage(shareRequest.Org, shareRequest.Content, shareRequest.Tags, shareRequest.Timestamp, shareRequest.FileName, shareRequest.Ip, shareRequest.Location, shareRequest.Policy)
	}
	if err := data.SaveSharedMessage(message, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ===================================================================================
// 分享信息
// ===================================================================================
func shareMessageThreshold(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("enter shareMessageThreshold")
	// 反序列化请求，验签
	var requestStr = args[0]
	log.Println(requestStr)
	shareRequest := new(request.ShareMessageRequest)
	if err := json.Unmarshal([]byte(requestStr), shareRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	message := data.NewThreholdSharedMessage(shareRequest.FileName, shareRequest.Uid, shareRequest.Timestamp)
	if err := data.SaveThreholdSharedMessage(message, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ===================================================================================
// 获得信息
// ===================================================================================
func getMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("get message")
	// 反序列化请求
	var requestStr = args[0]
	log.Println(requestStr)
	getRequest := new(request.GetSharedMessageRequest)
	if err := json.Unmarshal([]byte(requestStr), getRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	result, err := data.GetSharedMessage(getRequest.FromUid, getRequest.Tag, getRequest.PageSize, getRequest.Bookmark, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

// ===================================================================================
// 获得门限文件信息
// ===================================================================================
func getThreholdMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("get Threhold message")
	// 反序列化请求
	var requestStr = args[0]
	log.Println(requestStr)
	getRequest := new(request.GetThresholdSharedMessageRequest)
	if err := json.Unmarshal([]byte(requestStr), getRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	result, err := data.GetThreholdSharedMessage(getRequest.OrgName, getRequest.FileName, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

// ===================================================================================
// 检查请求参数并验签
// ===================================================================================
func preCheckRequest(requestStr string, uid, sign string, stub shim.ChaincodeStubInterface) error {
	requestJson, err := utils.GetRequestParamJson([]byte(requestStr))
	if err != nil {
		log.Println(err)
		return err
	}
	requestUser, err := data.QueryUserByUid(uid, stub)
	if err != nil {
		log.Println(err)
		return err
	}
	if requestUser == nil {
		log.Println("don't have requestUser with uid " + uid)
		return ecode.Error(ecode.RequestErr, "don't have this requestUser")
	}
	if err = utils.VerifySign(string(requestJson), requestUser.PublicKey, sign, uid); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
