package data

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
	"trustPlatform/constant"
	"trustPlatform/utils"
)

func init() {

	log.SetFlags(log.Ltime | log.Lshortfile)
}

// ===================================================================================
// 查找文件
// ===================================================================================
func GetSharedSM2Message(fromUid, toName string, pageSize int, bookmark string, stub shim.ChaincodeStubInterface) (result []byte, err error) {
	log.Printf("query shared message from %s to %s\n", fromUid, toName)
	if fromUid == "" && toName == "" {
		log.Println("cannot query all message")
		return nil, ecode.Error(ecode.RequestErr, "cannot query all message")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"$or\":[{\"uid\":\"%s\"},{\"toName\":\"%s\"}]}}", constant.SharedSM2Message, fromUid, toName)
	log.Println(queryString)
	result, err = utils.GetBytesFromDB2(stub, queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	return
}
