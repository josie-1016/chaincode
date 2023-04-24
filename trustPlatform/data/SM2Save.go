package data

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
	"trustPlatform/constant"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// ===================================================================================
// 保存分享的信息
// ===================================================================================
func SaveSharedSM2Message(message *SharedSM2Message, stub shim.ChaincodeStubInterface) (err error) {
	log.Printf("save shared message from: %s to: %s\n", message.Uid, message.ToName)
	log.Printf("timestamp:%s\n filename:%s\n", message.Timestamp, message.FileName)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	padding := "0"
	var key string
	sth := []byte{1}
	for len(sth) != 0 {
		hash := sha256.Sum256([]byte(message.Content + padding))
		key = constant.SharedSM2MessagePrefix + message.Uid + "to:" + message.ToName + ":" + hex.EncodeToString(hash[:])
		log.Println("key: ", key)
		sth, _ = stub.GetState(key)
		padding += padding
	}

	if err = stub.PutState(key, messageBytes); err != nil {
		return err
	}

	log.Printf("save shared message from: %s to: %s success\n", message.Uid, message.ToName)
	return
}
