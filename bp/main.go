package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//"github.com/josie-1016/bullet/rangproof"
	"bullet/rangproof"
	"log"
	"strings"
)

// 主模块
type BPCC struct {
	//RangeProof *RangeProof.RangeProof
}

func init() {
	log.SetPrefix("[BPCC]")
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(BPCC))
	if err != nil {
		log.Printf("Error starting BPCC chaincode: %s", err)
	}
}

// ===================================================================================
// 初始化函数，如果有需要初始化的模块请在自己模块中定义init函数，在此调用，不要直接写在主模块中
// ===================================================================================
func (b *BPCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("BPCC Init")
	// example
	//user.Init(stub)
	//d.Dabe = new(DecentralizedABE.DABE)
	//d.Dabe.GlobalSetup()
	rangproof.NewECPrimeGroupKey(64)
	return shim.Success([]byte("Init success"))
}

// ===================================================================================
// 调用合约入口，按模块划分，例如我们有用户、组织、数据分享三大模块（可以添加），按照路径层级传参
// 例如 用户声明新属性 /user/attribute/add args: userName, attributePublicKey
// user为用户模块，attribute为要操作的用户属性，add代表新建/声明，两个args按指定顺序传入
// ===================================================================================
func (b *BPCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("\nBPCC Invoke")
	//if d.Dabe == nil {
	//	d.Dabe = new(DecentralizedABE.DABE)
	//	d.Dabe.GlobalSetup()
	//}
	function, args := stub.GetFunctionAndParameters()
	if strings.HasPrefix(function, "/create") {
		return b.create(stub, args)
	} else if strings.HasPrefix(function, "/verify") {
		return b.verify(stub, args)
	} else if strings.HasPrefix(function, "/createBatch") {
		return b.createBatch(stub, args)
	} else if strings.HasPrefix(function, "/commit") {
		return b.commit(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"/create\" ")
}
