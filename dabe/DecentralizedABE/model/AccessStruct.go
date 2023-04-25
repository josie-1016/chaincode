package model

import (
	"github.com/Nik-U/pbc"
	"strconv"
	"strings"
	//	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

/*  AccessStruct [][] */
type AccessStruct struct {
	CurrentPointer int
	LeafID         int
	A              [][]int
	//A_LsssLine []string
	LsssMatrix        [][]int8
	ElementLsssMatrix [][]*pbc.Element
	PolicyMap         map[string]int //1 Attr, 1 index
	PolicyMaps        []string       //1 index, 1 Attr
	PolicyTreePath    [][]int8       //1 index, 1 path
	W                 []*pbc.Element
}

func NewAccessStruct() (*AccessStruct) {
	A := new(AccessStruct)
	A.CurrentPointer = 0
	A.LeafID = -1
	A.A = make([][]int, 0)
	return A
}

func (this *AccessStruct) ParsePolicyStringtoMap(s *string) (map[string]int) {
	this.PolicyMap = make(map[string]int)
	*s = strings.Replace(*s, "AND", ",", -1)
	*s = strings.Replace(*s, "OR", ",", -1)
	*s = strings.Replace(*s, " ", "", -1)
	*s = strings.Replace(*s, "(", "", -1)
	*s = strings.Replace(*s, ")", "", -1)
	// 先提取出所有的属性名形成数组
	str := strings.Split(*s, ",")

	this.PolicyMaps = make([]string, len(str)+1, len(str)+1)

	for v := range str {
		if str[v] != "" {
			// 根据属性出现的顺序赋予index，从1开始
			this.PolicyMap[str[v]] = v + 1
			this.PolicyMaps[v+1] = str[v]
		}
	}
	return this.PolicyMap
}

func (this *AccessStruct) genLsssMatrix() {
	// A_LsssLine表示Lsss矩阵，因为每个String表示矩阵一行
	A_LsssLine := make([]string, len(this.A), len(this.A))

	//x := make([]string,len(this.A),len(this.A))
	//this.PolicyTreePath = make([]string,len(this.PolicyMap),len(this.PolicyMap))

	// 初始向量为(1)
	A_LsssLine[0] = "1"
	//x[0] = "1"
	// Lsss矩阵的行数是属性个数+1 一个属性一行
	this.LsssMatrix = make([][]int8, len(this.PolicyMap)+1, len(this.PolicyMap)+1)

	// 遍历每一行？
	for i := 0; i < len(A_LsssLine); i++ {
		var sign, first int64
		if this.A[i][0] == this.A[i][1] { //表示当前为AND门
			sign = int64(-1) // 通过给第一个child赋值first 给剩下其他child赋值-1 这样所有child加起来的值为0
			first = int64(len(this.A[i]) - 3) //first表示除去第一个child外的child个数
		} else { //表示当前为OR门
			sign = int64(0) // 因为是OR门 因此每个child都赋值为0即可
			first = int64(0)
		}
		if this.A[i][2] > 0 { // 说明这个child是一个ThresholdNode 不是属性
			// 同层第一个Threshold child节点赋值为first，但是不解析为向量，因为后续还需要解析其子节点
			A_LsssLine[this.A[i][2]] = A_LsssLine[i] + "," + strconv.FormatInt(first, 10) 
		} else {// 说明这个child是一个属性叶节点
			// 直接解析成矩阵向量
			tmp := strings.Split(A_LsssLine[i]+","+strconv.FormatInt(first, 10), ",")
			this.LsssMatrix[-this.A[i][2]] = make([]int8, len(tmp), len(tmp))
			for k := 0; k < len(tmp); k++ {
				n, err := strconv.ParseInt(tmp[k], 10, 0)
				this.LsssMatrix[-this.A[i][2]][k] = int8(n)
				if err == nil {
				}
			}
		}
		//遍历后面的child
		for j := 3; j < len(this.A[i]); j++ {
			if this.A[i][j] > 0 { //判断是不是属性叶节点
				A_LsssLine[this.A[i][j]] = A_LsssLine[i] + "," + strconv.FormatInt(sign, 10)
			} else {
				// 分割当前A_LsssLine 转换为Lsss矩阵对应行的向量 行是由属性的idx标记的
				tmp := strings.Split(A_LsssLine[i]+","+strconv.FormatInt(sign, 10), ",")
				this.LsssMatrix[-this.A[i][j]] = make([]int8, len(tmp), len(tmp))
				for k := 0; k < len(tmp); k++ {
					n, err := strconv.ParseInt(tmp[k], 10, 0)
					this.LsssMatrix[-this.A[i][j]][k] = int8(n)
					if err == nil {
					}
				}
			}
			//sign = -sign
		}
	}
}

func (this *AccessStruct) genPolicyTreePath() {
	PolicyTreeLine := make([]string, len(this.A), len(this.A))
	PolicyTreeLine[0] = "1"
	this.PolicyTreePath = make([][]int8, len(this.PolicyMap)+1, len(this.PolicyMap)+1)

	for i := 0; i < len(PolicyTreeLine); i++ {
		var sign int64
		sign = int64(-1)
		// 直接从第一个child开始
		for j := 2; j < len(this.A[i]); j++ {
			if this.A[i][j] > 0 { 
				PolicyTreeLine[this.A[i][j]] = PolicyTreeLine[i] + "," + strconv.FormatInt(sign, 10)
			} else {
				tmp := strings.Split(PolicyTreeLine[i]+","+strconv.FormatInt(sign, 10), ",")
				this.PolicyTreePath[-this.A[i][j]] = make([]int8, len(tmp), len(tmp))
				for k := 0; k < len(tmp); k++ {
					n, err := strconv.ParseInt(tmp[k], 10, 0)
					this.PolicyTreePath[-this.A[i][j]][k] = int8(n)
					if err == nil {
					}
				}
			}
			sign = -sign
		}
	}
}

func (this *AccessStruct) padLsssMatrix() {
	n := len(this.LsssMatrix)
	l := 0
	// l表示Lsss矩阵的最大列数
	for i := 0; i < n; i++ {
		if l < len(this.LsssMatrix[i]) {
			l = len(this.LsssMatrix[i])
		}
	}
	// 如果当前行的列数不足l 后面补0到l长
	for i := 0; i < n; i++ {
		if l > len(this.LsssMatrix[i]) {
			for j := len(this.LsssMatrix[i]); j < l; j++ {
				this.LsssMatrix[i] = append(this.LsssMatrix[i], 0)
			}
		}
	}
	//fmt.Printf("l:: %v\n",l)
}

func (this *AccessStruct) genElementLsssMatrix(a *pbc.Element) {
	// n为Lsss矩阵的行数
	n := len(this.LsssMatrix)
	// l为Lsss矩阵的列数
	l := len(this.LsssMatrix[0])
	// 将Lsss矩阵转换为群元素的Lsss矩阵
	this.ElementLsssMatrix = make([][]*pbc.Element, n, n)
	for i := 0; i < n; i++ {
		this.ElementLsssMatrix[i] = make([]*pbc.Element, l, l)
		for j := 0; j < l; j++ {
			this.ElementLsssMatrix[i][j] = a.NewFieldElement().SetInt32(int32(this.LsssMatrix[i][j]))
		}
	}
}

func (this *AccessStruct) LsssMatrixDotMulVector(row int, u []*pbc.Element) (*pbc.Element) {
	l := len(u)
	result := u[0].NewFieldElement().Set0()
	for i := 0; i < l; i++ {
		tmp := this.ElementLsssMatrix[row][i].NewFieldElement().Set(this.ElementLsssMatrix[row][i])
		tmp.Mul(tmp, u[i])
		result.Add(result, tmp)
	}
	return result
}

func (this *AccessStruct) gen_w(a *pbc.Element) {
	n := len(this.LsssMatrix)
	this.W = make([]*pbc.Element, n, n)
	this.W[1] = a.NewFieldElement().Set1().ThenDiv(a.NewFieldElement().SetInt32(2))
	this.W[2] = a.NewFieldElement().Set1().ThenDiv(a.NewFieldElement().SetInt32(2))
	/*
		x := a.NewFieldElement().SetInt32(8)
		x.ThenMul(this.w[1])*/

	//fmt.Printf("w:: %v\n", this.w[1])
	//fmt.Printf("w:: %v\n", this.w[2])
	//	fmt.Printf("x:: %v\n", x)
}

//// stop point
func (this *AccessStruct) isSatisfy(node *[]int, leaf *[]int, Smap *map[string]int, r int) (bool) {
	var satisfy bool
	// 表示当前行满足的child个数
	var t = 0 
	// r指示了当前的行数
	// 从第一个child节点开始
	for i := 2; i < len(this.A[r]); i++ {
		if this.A[r][i] > 0 { // 说明当前是Threshold节点
			// 递归查看Threshold节点是否满足
			if this.isSatisfy(node, leaf, Smap, this.A[r][i]) {
				t++ // 满足则计数+1
				(*node)[this.A[r][i]] = 1 // 把对应的行设为1 表示满足
			} else {
				(*node)[this.A[r][i]] = 0 // 把对应行设为0 表示不满足
			}
		} else { // 说明当前是属性节点
			// 说明拥有此属性的priKey
			if (*Smap)[this.PolicyMaps[-this.A[r][i]]] == 1 {
				t++ // 满足计数+1
				(*leaf)[-this.A[r][i]] = 1 // 把对应的属性设为1 表示满足
			} else { 
				(*leaf)[-this.A[r][i]] = 0 // 把对应属性设为0 表示不满足
			}
		}
	}
	for i := 2; i < len(this.A[r]); i++ {
		if this.A[r][i] > 0 {
			if (*node)[this.A[r][i]] > 0 { // 说明此Threshold节点已满足
				(*node)[this.A[r][i]] = t // 设为当前行满足的child节点个数
			}
		} else {
			if (*leaf)[-this.A[r][i]] > 0 {
				(*leaf)[-this.A[r][i]] = t
			}
		}
	}
	if t < this.A[r][1] { // 说明当前行满足的节点个数未达到阈值要求
		satisfy = false
	} else {
		satisfy = true
	}
	return satisfy
}

type sendacs struct {
	CurrentPointer int
	LeafID         int
	A              [][]int
	//A_LsssLine []string
	LsssMatrix        [][]int8
	ElementLsssMatrix [][][]byte
	PolicyMap         map[string]int //1 Attr, 1 index
	PolicyMaps        []string       //1 index, 1 Attr
	PolicyTreePath    [][]int8       //1 index, 1 path
	W                 [][]byte
}

func (b *AccessStruct) Serialize() []byte {
	var result bytes.Buffer
	var sacs *sendacs
	sacs = new(sendacs)
	sacs.CurrentPointer = b.CurrentPointer
	sacs.LeafID = b.LeafID
	sacs.LsssMatrix = b.LsssMatrix
	sacs.PolicyMap = b.PolicyMap
	sacs.PolicyMaps = b.PolicyMaps
	sacs.PolicyTreePath = b.PolicyTreePath
	sacs.A = b.A
	for _, e1 := range b.ElementLsssMatrix {
		var temp [][]byte
		for _, e2 := range e1 {
			temp = append(temp, e2.Bytes())
		}
		sacs.ElementLsssMatrix = append(sacs.ElementLsssMatrix, temp)
	}
	for _, w := range b.W {
		sacs.W = append(sacs.W, w.Bytes())
	}

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(sacs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeACS(d []byte, cp *CurveParam) *AccessStruct {
	acs := new(AccessStruct)
	var sacs sendacs

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&sacs)
	if err != nil {
		log.Panic(err)
	}

	acs.CurrentPointer = sacs.CurrentPointer
	acs.LeafID = sacs.LeafID
	acs.LsssMatrix = sacs.LsssMatrix
	acs.PolicyMap = sacs.PolicyMap
	acs.PolicyMaps = sacs.PolicyMaps
	acs.PolicyTreePath = sacs.PolicyTreePath
	acs.A = sacs.A
	for _, e1 := range sacs.ElementLsssMatrix {
		var temp []*pbc.Element
		for _, e2 := range e1 {
			temp = append(temp, cp.GetNewZn().SetBytes(e2))
		}
		acs.ElementLsssMatrix = append(acs.ElementLsssMatrix, temp)
	}
	for _, w := range sacs.W {
		acs.W = append(acs.W, cp.GetNewZn().SetBytes(w))
	}

	return acs
}