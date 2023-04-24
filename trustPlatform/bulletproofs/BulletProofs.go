package bulletproofs

import (
	"encoding/json"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
	"trustPlatform/data"
	"trustPlatform/request"
	"trustPlatform/utils"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// ===================================================================================
// bulletproofs合约初始化函数
// ===================================================================================
func Init(stub shim.ChaincodeStubInterface) {
	log.Println("BulletProofs init")
}

// ===================================================================================
// bulletproofs模块入口函数
// ===================================================================================
func Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("TrustPlatformCC Org Invoke")
	function, args := stub.GetFunctionAndParameters()
	if err := utils.CheckInputNumber(1, args); err != nil {
		return shim.Error(err.Error())
	}
	if function == "/zk/create" {
		return createBulletProofs(stub, args)
	} else if function == "/zk/getZKs" {
		return getBulletProofs(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"/org/createOrgApply\" \"/org/approveOrgApply\" " +
		"\"/org/shareSecret\" \"/org/getSharedSecret\" \"/org/submitPartPK\" \"/org/mixPartPK\"" +
		" \"/org/queryOrg\" ")
}

func getBulletProofs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("get bulletproofs")
	// 反序列化请求
	var requestStr = args[0]
	log.Println(requestStr)
	getRequest := new(request.GetBulletProofsRequest)
	if err := json.Unmarshal([]byte(requestStr), getRequest); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	result, err := data.GetBulletProofs(getRequest.Uid, getRequest.Pid, getRequest.Tag, getRequest.PageSize, getRequest.Bookmark, stub)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func createBulletProofs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//如果存proof，要先找到之前存commit的地方，修改proof
	log.Println("create new bulletproofs")
	var requestStr = args[0]
	log.Println(requestStr)
	createRequest := new(request.BulletProofsRequest)
	if err := json.Unmarshal([]byte(requestStr), createRequest); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if err := preCheckRequest(requestStr, createRequest.Uid, createRequest.Sign, stub); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	if len(createRequest.Tags) > 10 {
		return shim.Error("too much tags")
	}
	//TODO:saveCommit也用这个？只有uid，commit1和tags
	if createRequest.Range == "" {
		message := data.NewBulletProofs(createRequest.Uid, createRequest.Pid, createRequest.Tags, createRequest.Timestamp, createRequest.Range, createRequest.Commit1, createRequest.Commit2, createRequest.Open, createRequest.Proof, createRequest.ProofPre)
		if err := data.SaveBulletProofs(message, stub); err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
	} else {
		proof, err := data.QueryBulletProofs(createRequest.Uid, createRequest.Pid, stub)
		if err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
		proof.Range = createRequest.Range
		proof.Proof = createRequest.Proof
		proof.Commit2 = createRequest.Commit2
		if err := data.SaveBulletProofs(proof, stub); err != nil {
			log.Println(err)
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
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
	if err = utils.VerifySign(string(requestJson), requestUser.PublicKey, sign); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
