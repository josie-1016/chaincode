package main

import (
	"github.com/josie-1016/bullet/rangproof"
	//"bullet/rangproof"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
	"math/big"
	"strconv"
)

func (b *BPCC) commit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("create commit")
	vstr := args[0]
	log.Println(vstr)
	v, _ := strconv.ParseInt(vstr, 10, 64)
	//生成承诺和[0,2^64]的范围证明
	ran, com, open, err := rangproof.RPProve(big.NewInt(v), nil)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	log.Println("commit1:", com)
	res := new(CommitResponse)
	res.Proof = ran
	res.Commit1 = com
	res.Open = open
	bytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

func (b *BPCC) verifyCommits(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//批量证明检查
	log.Println("verify Commits")
	requestBytes := args[0]
	log.Println(requestBytes)
	request := new(VerifyCommitsRequest)
	if err := rangproof.Deserialize2Struct([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	log.Println(request.Commit1)
	log.Println(request.Commits[0])
	com, err := rangproof.PedersenBatchAddCommit(request.Commits)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if request.Commit1.Comm.Equal(com.Comm) {
		log.Println("commits right")
	} else {
		log.Println("commits wrong")
		return shim.Error("commits wrong")
	}

	r, _ := strconv.ParseInt(request.Range, 10, 64)
	comv, _ := rangproof.PedersenSubNum(request.Commit1, big.NewInt(r))
	//测试是否与SubProof方法计算出的承诺值相同
	if request.Commit2.Comm.Equal(comv.Comm) {
		log.Println("commit2 right")
	} else {
		log.Println("commit2 wrong")
		log.Println(comv.Comm)
		log.Println(request.Commit2.Comm)
		return shim.Error("commit2 wrong")
	}
	//验证
	if rangproof.RPVerify(request.Proof, comv) {
		return shim.Success(nil)
	}
	return shim.Error("verify fail")

}
func (b *BPCC) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("create bulletproof")
	vstr := args[0]
	rstr := args[1]
	openstr := args[2]
	log.Println("open", openstr)
	v, _ := strconv.ParseInt(vstr, 10, 64)
	r, _ := strconv.ParseInt(rstr, 10, 64)
	open, _ := new(big.Int).SetString(openstr, 10)
	log.Println("open", open)
	_, com1, _, err := rangproof.RPProve(big.NewInt(v), open)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	log.Println("commit1:", com1)
	ran2, com2, _, err := rangproof.SubProof(big.NewInt(v), big.NewInt(r), com1, open)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	res := new(RangeProofResponse)
	res.Commit1 = com1
	res.Commit2 = com2
	res.Proof = ran2
	//bytes, err := rangproof.Serialize2Bytes(res)
	bytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

func (b *BPCC) verify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("verify bulletproof")
	requestBytes := args[0]
	log.Println(requestBytes)
	request := new(VerifyRequest)
	if err := rangproof.Deserialize2Struct([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	log.Println(request)
	if request.ProofPre != nil {
		log.Println("check proofpre")
		ret := rangproof.RPVerify(request.ProofPre, request.Commit1)
		if !ret {
			return shim.Error("proofpre wrong")
		}
	}
	r, _ := strconv.ParseInt(request.Range, 10, 64)
	comv, _ := rangproof.PedersenSubNum(request.Commit1, big.NewInt(r))
	//测试是否与SubProof方法计算出的承诺值相同
	if request.Commit2.Comm.Equal(comv.Comm) {
		log.Println("commit2 right")
	} else {
		log.Println("commit2 wrong")
		log.Println(comv.Comm)
		log.Println(request.Commit2.Comm)
		return shim.Error("commit2 wrong")
	}
	//验证
	if rangproof.RPVerify(request.Proof, comv) {
		return shim.Success(nil)
	}
	return shim.Error("verify fail")
	//return shim.Success(nil)
}
