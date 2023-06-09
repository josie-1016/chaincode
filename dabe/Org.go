package main

import (
	"log"

	DecentralizedABE "github.com/vangogo/tree/ThresholdABE"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ===================================================================================
// org模块入口函数
// ===================================================================================
func (d *DABECC) OrgInvoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("DABECC Org Invoke")
	function, args := stub.GetFunctionAndParameters()

	if function == "/org/generateOPK" {
		return d.generateOPK(stub, args)
	} else if function == "/org/generateMenPub" {
		return d.generateMen(stub, args)
	} else if function == "/org/generateAPK" {
		return d.generateAPK(stub, args)
	} else if function == "/org/generateMenPriv" {
		return d.generateMenPriv(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"/org/generateOPK\" ")
}

func (d *DABECC) generateOPK(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("org generateOPK")
	requestStr := args[0]
	log.Println(requestStr)
	request := new(GenerateOPKRequest)
	if err := DecentralizedABE.Deserialize2Struct([]byte(requestStr), request); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	t := request.T
	n := request.N

	org, err := d.Dabe.OrgSetup(n, t, "temp", request.UserNames)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if err = org.GenerateOPK(request.UserNames[:t], request.PartPkList[:t], d.Dabe); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	firstEggAlpha := org.EGGAlpha.NewFieldElement().Set(org.EGGAlpha)
	if err = org.GenerateOPK(request.UserNames[n-t:], request.PartPkList[n-t:], d.Dabe); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if firstEggAlpha.Equals(org.EGGAlpha) {
		return shim.Success([]byte(firstEggAlpha.String()))
	}

	return shim.Error("unknown error")
}

func (d *DABECC) generateAPK(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("org generateAPK")
	requestStr := args[0]
	log.Println(requestStr)
	request := new(GenerateAPKRequest)
	if err := DecentralizedABE.Deserialize2Struct([]byte(requestStr), request); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	t := request.T
	n := request.N
	attrName := request.AttrName

	org, err := d.Dabe.OrgSetup(n, t, "temp", request.UserNames)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	if err = org.GenerateNewAttr(request.UserNames[:t], request.PartPkList[:t], attrName, d.Dabe); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	log.Printf("%+v\n", org.APKMap[attrName])
	firstGy := org.APKMap[attrName].Gy.NewFieldElement().Set(org.APKMap[attrName].Gy)
	delete(org.APKMap, attrName)
	if err = org.GenerateNewAttr(request.UserNames[n-t:], request.PartPkList[n-t:], attrName, d.Dabe); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	log.Printf("%+v\n", org.APKMap[attrName])
	if firstGy.Equals(org.APKMap[attrName].Gy) {
		log.Println(firstGy.String())
		return shim.Success([]byte(firstGy.String()))
	}
	return shim.Error("unknown error")
}

func (d *DABECC) generateMen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("org generateMen")
	requestStr := args[0]
	log.Println(requestStr)
	request := new(GenerateOSKRequest)
	if err := DecentralizedABE.Deserialize2Struct([]byte(requestStr), request); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	t := request.T
	n := request.N
	len := len(request.PartPkList)
	org, err := d.Dabe.OrgSetup(n, t, "temp", request.UserNames)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	pub, err := org.ABC(request.UserNames[:t], request.PartPkList[:t], d.Dabe)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	pub1, err := org.ABC(request.UserNames[len-t:], request.PartPkList[len-t:], d.Dabe)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	log.Println(pub)
	log.Println(pub1)
	if string(pub) == string(pub1) {
		return shim.Success(pub)
	}

	return shim.Error("unknown error")
}

func (d *DABECC) generateMenPriv(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("org generateMenPriv")
	requestStr := args[0]
	log.Println(requestStr)
	request := new(GenerateOSKRequest)
	if err := DecentralizedABE.Deserialize2Struct([]byte(requestStr), request); err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	t := request.T
	n := request.N
	org, err := d.Dabe.OrgSetup(n, t, "temp", request.UserNames)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	priv, err := org.DEF(request.UserNames[:t], request.PartPkList[:t], d.Dabe)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(priv)
}
