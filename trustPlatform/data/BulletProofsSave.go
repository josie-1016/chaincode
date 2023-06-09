package data

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
	"trustPlatform/constant"
)

func SaveBulletProofs(proof *BulletProofs, stub shim.ChaincodeStubInterface) (err error) {
	log.Printf("save bulletproofs from: %s with %s\n", proof.Uid, proof.Tags)
	proofBytes, err := json.Marshal(proof)
	if err != nil {
		return err
	}
	if err = stub.PutState(constant.BulletProofs+proof.Uid+proof.Pid, proofBytes); err != nil {
		return err
	}

	log.Printf("save bulletproofs from: %s with %s success\n", proof.Uid, proof.Tags)
	return
}
