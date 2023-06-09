package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
	"trustPlatform/constant"
	"trustPlatform/utils"
)

func GetBulletProofs(uid, pid, tag string, pageSize int, bookmark string, stub shim.ChaincodeStubInterface) (result []byte, err error) {
	log.Printf("query bulletproofs from %s with pid=%s and tags=%s\n", uid, pid, tag)
	if uid == "" && pid == "" && tag == "" {
		log.Println("cannot query all bulletproofs")
		return nil, ecode.Error(ecode.RequestErr, "cannot query all bulletproofs")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"", constant.BulletProofs)
	if uid != "" {
		queryString += fmt.Sprintf(",\"uid\":\"%s\"", uid)
	}
	if pid != "" {
		queryString += fmt.Sprintf(",\"pid\":\"%s\"", pid)
	}
	if tag != "" {
		queryString += fmt.Sprintf(",\"tags\":{\"$elemMatch\":{\"$eq\":\"%s\"}}", tag)
	}
	queryString += "}}"

	log.Println(queryString)
	if result, err = utils.GetBytesFromDB2(stub, queryString, pageSize, bookmark); err != nil {
		return nil, err
	}
	return
}

func QueryBulletProofsBytes(uid string, pid string, stub shim.ChaincodeStubInterface) (attr []byte, err error) {
	log.Println("query BulletProofs by pid: " + pid + " uid:" + uid)

	attr, err = stub.GetState(constant.BulletProofs + uid + pid)
	if err != nil {
		return nil, err
	}
	return
}

func QueryBulletProofs(uid string, pid string, stub shim.ChaincodeStubInterface) (proof *BulletProofs, err error) {
	bytes, err := QueryBulletProofsBytes(uid, pid, stub)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, nil
	}

	if err = json.Unmarshal(bytes, &proof); err != nil {
		return nil, err
	}
	return
}

func QueryCommits(uid string, pids []string, stub shim.ChaincodeStubInterface) (resBytes []byte, err error) {
	log.Println("query Commits by uid:" + uid)
	var commits []Commit
	for _, pid := range pids {
		bytes, err := stub.GetState(constant.BulletProofs + uid + pid)
		if err != nil {
			return nil, err
		}
		bp := new(BulletProofs)

		if err = json.Unmarshal(bytes, &bp); err != nil {
			return nil, err
		}
		commits = append(commits, bp.Commit1)
	}
	resBytes, err = json.Marshal(commits)
	if err != nil {
		return nil, err
	}
	return
}
