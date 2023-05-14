package sm2common

import (
	"encoding/json"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
	"strings"
	"trustPlatform/data"
	"trustPlatform/request"
	"trustPlatform/utils"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// ===================================================================================
// common初始化函数
// ===================================================================================
func Init(stub shim.ChaincodeStubInterface) {
	log.Println("SM2Common init")
}

// ===================================================================================
// sm2common模块入口函数
// ===================================================================================
func Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("TrustPlatformCC SM2Common Invoke")
	function, args := stub.GetFunctionAndParameters()
	if err := utils.CheckInputNumber(1, args); err != nil {
		return shim.Error(err.Error())
	}
	if strings.HasPrefix(function, "/sm2common/shareMessage") {
		return shareMessage(stub, args)
	} else if strings.HasPrefix(function, "/sm2common/getMessage") {
		return getMessage(stub, args)
	}

	return shim.Error("Invalid invoke function name. ")
}

// ===================================================================================
// 分享信息
// ===================================================================================
func shareMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("enter SM2 shareMessage")
	// 反序列化请求，验签
	var requestStr = args[0]
	log.Println(requestStr)
	shareRequest := new(request.ShareSM2MessageRequest)
	if err := json.Unmarshal([]byte(requestStr), shareRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	if err := preCheckRequest(requestStr, shareRequest.Uid, shareRequest.Sign, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	message := data.NewSharedSM2Message(shareRequest.Uid, shareRequest.Content, shareRequest.Timestamp, shareRequest.FileName, shareRequest.ToName)
	if err := data.SaveSharedSM2Message(message, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ===================================================================================
// 获得信息
// ===================================================================================
func getMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("get SM2 message")
	// 反序列化请求
	var requestStr = args[0]
	log.Println(requestStr)
	getRequest := new(request.GetSharedSM2MessageRequest)
	if err := json.Unmarshal([]byte(requestStr), getRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	result, err := data.GetSharedSM2Message(getRequest.FromUid, getRequest.ToName, getRequest.PageSize, getRequest.Bookmark, stub)
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
