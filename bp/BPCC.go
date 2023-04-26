package main

import (
	"bullet/rangproof"
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
	v, _ := strconv.ParseInt(vstr, 10, 64)
	com, err := rangproof.PerdersenCommit(big.NewInt(v), nil)
	//TODO:error情况
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	//bytes, err := rangproof.Serialize2Bytes(com)
	bytes, err := json.Marshal(com)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

func (b *BPCC) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("create bulletproof")
	vstr := args[0]
	rstr := args[1]
	v, _ := strconv.ParseInt(vstr, 10, 64)
	r, _ := strconv.ParseInt(rstr, 10, 64)
	_, com1, open1 := rangproof.RPProve(big.NewInt(v), nil)
	ran2, com2, _ := rangproof.SubProof(big.NewInt(v), big.NewInt(r), com1, open1)
	//TODO:error情况
	//可以将生成的com1、com2和ran2发送给另一方，证明自己的值大于5

	//TODO:序列化？？
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
	requestBytes := args[0]
	request := new(VerifyRequest)
	//TODO:序列化？？
	if err := json.Unmarshal([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	//if err := DecentralizedABE.Deserialize2Struct([]byte(requestBytes), request); err != nil {
	//	log.Println(err)
	//	return shim.Error(err.Error())
	//}

	com1 := new(rangproof.PedersenCommit)
	//if err := DecentralizedABE.Deserialize2Struct([]byte(request.Commit1), com1); err != nil {
	//	log.Println(err)
	//	return shim.Error(err.Error())
	//}
	if err := json.Unmarshal([]byte(request.Commit1), com1); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	com2 := new(rangproof.PedersenCommit)
	//if err := DecentralizedABE.Deserialize2Struct([]byte(request.Commit2), com2); err != nil {
	//	log.Println(err)
	//	return shim.Error(err.Error())
	//}
	if err := json.Unmarshal([]byte(request.Commit2), com2); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	ran := new(rangproof.RangeProof)
	//if err := DecentralizedABE.Deserialize2Struct([]byte(request.Proof), com2); err != nil {
	//	log.Println(err)
	//	return shim.Error(err.Error())
	//}
	if err := json.Unmarshal([]byte(request.Proof), ran); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	r, _ := strconv.ParseInt(request.Range, 10, 64)
	comv, _ := rangproof.PedersenSubNum(com1, big.NewInt(r))
	//测试是否与SubProof方法计算出的承诺值相同
	if com2.Comm.Equal(comv.Comm) {
		log.Println("commit right")
	} else {
		return shim.Error("commit2 wrong")
	}
	//验证
	if rangproof.RPVerify(ran, comv) {
		return shim.Success(nil)
	}
	//TODO:验证失败
	return shim.Error("verify fail")
}

//
//func (d *DABECC) decrypt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	requestBytes := args[0]
//	request := new(DecryptRequest)
//	if err := DecentralizedABE.Deserialize2Struct([]byte(requestBytes), request); err != nil {
//		log.Println(err)
//		return shim.Error(err.Error())
//	}
//	cipher := new(DecentralizedABE.Cipher)
//	if err := DecentralizedABE.Deserialize2Struct([]byte(request.Cipher), cipher); err != nil {
//		log.Println(err)
//		return shim.Error(err.Error())
//	}
//
//	log.Printf("%+v\n", request.AttrMap)
//	log.Printf("%+v\n", cipher)
//
//	plainText, err := d.Dabe.Decrypt(cipher, request.AttrMap, request.Uid)
//	if err != nil {
//		log.Println(err)
//		return shim.Error(err.Error())
//	}
//
//	return shim.Success(plainText)
//}
