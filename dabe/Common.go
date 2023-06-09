package main

import (
	"encoding/json"
	"log"

	DecentralizedABE "github.com/vangogo/tree/ThresholdABE"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ===================================================================================
// common模块入口函数
// ===================================================================================
func (d *DABECC) CommonInvoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("DABECC Common Invoke")
	function, args := stub.GetFunctionAndParameters()

	if function == "/common/encryptThreshold" {
		return d.encryptMen(stub, args)
	} else if function == "/common/decryptThreshold" {
		return d.decryptMen(stub, args)
	} /* else if function == "/common/decrypt" {
		return d.decrypt(stub, args)
	} */ /* else if function == "/common/encrypt" {
		return d.encrypt(stub, args)
	}  */

	return shim.Error("Invalid invoke function name. Expecting \"/common/encrypt\" \"/common/decrypt\"")
}

/* func (d *DABECC) encrypt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	requestBytes := args[0]
	request := new(EncryptRequest)
	log.Println(requestBytes)
	if err := DecentralizedABE.Deserialize2Struct([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	authorities := make(map[string]DecentralizedABE.Authority)
	for key, value := range request.AuthorityMap {
		authorities[key] = value
	}
	cipher, err := d.Dabe.Encrypt(request.PlainContent, request.Policy, authorities)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	bytes, err := DecentralizedABE.Serialize2Bytes(cipher)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
} */

func (d *DABECC) encryptMen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	requestBytes := args[0]
	request := new(ThreholdEncryptRequest)
	log.Println(requestBytes)
	if err := json.Unmarshal([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	cipher, err := d.Dabe.EncryptMen(request.PlainContent, []byte(request.PubKey))
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	bytes, err := json.Marshal(cipher)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
}

/* func (d *DABECC) decrypt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	requestBytes := args[0]
	request := new(DecryptRequest)
	if err := DecentralizedABE.Deserialize2Struct([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	cipher := new(DecentralizedABE.Cipher)
	if err := DecentralizedABE.Deserialize2Struct([]byte(request.Cipher), cipher); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	log.Printf("%+v\n", request.AttrMap)
	log.Printf("%+v\n", cipher)

	plainText, err := d.Dabe.Decrypt(cipher, request.AttrMap, request.Uid)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	return shim.Success(plainText)
} */

func (d *DABECC) decryptMen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	requestBytes := args[0]
	request := new(ThreholdDecryptRequest)
	if err := json.Unmarshal([]byte(requestBytes), request); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}
	cipher := new(DecentralizedABE.MenCipher)
	if err := json.Unmarshal([]byte(request.Cipher), cipher); err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	log.Printf("%+v\n", request.ThresholdPriv)
	log.Printf("%+v\n", cipher)

	plainText, err := d.Dabe.DecryptMen(cipher, request.ThresholdPriv)
	if err != nil {
		log.Println(err)
		return shim.Error(err.Error())
	}

	return shim.Success(plainText)
}
