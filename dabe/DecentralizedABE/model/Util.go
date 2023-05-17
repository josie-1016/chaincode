package model

import (
	"fmt"
	"strconv"
	"strings"
	"log"
)

/* <Policy Parser SECTION */
func ParsePolicyStringToTree(s *string) (*PolicyNode, *AccessStruct) {
	ss := *s
	AS := NewAccessStruct()
	AS.ParsePolicyStringtoMap(&ss)

	*s = strings.Replace(*s, "AND", "&&", -1)
	*s = strings.Replace(*s, "OR", "||", -1)
	var num int = len(strings.Split(*s," "))
	*s = strings.Replace(*s, " ", "", -1)
	MainPolicy, ID := ParsePolicyString1(AS, s, 0, len(*s)-1,num)
	if ID == 0 {
	} //non sense
	return MainPolicy, AS
}

//传入的Policy是有左括号的
func ParsePolicyString1(A *AccessStruct, s *string, startPos int, stopPos int,num int) (*PolicyNode, int) {
	//leftPos := startPos+1+strings.Index((*s)[startPos+1:stopPos], "(")
	log.Println(num)
	this := NewPolicyNode("ThreshHold", 0) //当前节点是阈值节点，即代表AND门OR门

	A.A = append(A.A, make([]int, 2, 2))
	//_A := &(A.A[A.CurrentPointer])
	ID := A.CurrentPointer
	A.A[ID][0] = 0
	A.A[ID][1] = 0
	A.CurrentPointer++
	policy_children := make([]*PolicyNode, 0)

	var i int = startPos + 1 //startPos表示的是起始左括号
	var n int = 0
	var _n int = 0
	var leftPos int
	var rightPos int
	var trueChild string = ""

	for i <= stopPos {
		// 在当前String[startPos:stopPos]中找到第一个左括号(
		leftPos = strings.Index((*s)[i:stopPos], "(")
		// 存在左括号，说明括号内可以继续递归解析，例如(A AND B AND (C OR D))
		if leftPos != -1 {
			//左括号之前的是同层的policy，不包含括号闭合区间了
			//例子中是（C OR D）
			trueChild += (*s)[i : i+leftPos] // A AND B
			rightPos = LookForMyRightBraket(s, i+leftPos)
			tmpPolicy, tmpID := ParsePolicyString2(A, s, i+leftPos, rightPos) //继续递归（C OR D）
			policy_children = append(policy_children, tmpPolicy)
			A.A[ID] = append(A.A[ID], tmpID)
			n++
			i = rightPos + 1 // 从这一段括号后继续寻找新的括号闭合区间
		} else {
			//如果没有新的括号闭合区间，说明剩下的都是不带括号的policy，即可以认为是同层policy
			trueChild += (*s)[i:stopPos]
			break
		}
	}

	var childAttr []string
	//注意：同层policy中只会出现AND或者OR，不会两者都有
	// 存在AND门
	// 假设现在递归A AND B
	if strings.Index(trueChild, "&&") != -1 {
		childAttr = strings.Split(trueChild, "&&")

		for v := range childAttr {
			if childAttr[v] != "" {
				policy_children = append(policy_children, NewPolicyNode(childAttr[v], 1).SetMax(1).SetMin(1))
				A.A[ID] = append(A.A[ID], -A.PolicyMap[childAttr[v]])
				A.LeafID--
				n++
			}
		}
		_n = n
		this.SetOperation(1)
		// 存在OR门
		//假设现在递归：C OR D
	} else if strings.Index(trueChild, "||") != -1 || num == 1{
		fmt.Printf("in OR Gate\n")
		if num==1 {
			childAttr = strings.Split(trueChild," ")
		} else {
			childAttr = strings.Split(trueChild, "||")
		}
		//childAttr包含了policy的属性
		for v := range childAttr {
			if childAttr[v] != "" {
				//加入包含属性的叶子结点policyNode
				policy_children = append(policy_children, NewPolicyNode(childAttr[v], 1).SetMax(1).SetMin(1))
				// 加入属性对应idx的相反值
				A.A[ID] = append(A.A[ID], -A.PolicyMap[childAttr[v]])
				A.LeafID--
				n++
			}
		}
		_n = 1
		this.SetOperation(2) //当前门为OR门
	}

	if n == 0 {
		fmt.Printf("Error:: bad description. \n")
	} else {
		this.SetChildren(policy_children) //设置这个门节点的子节点
		this.SetMax(n) // 说明当前子节点个数
		this.SetMin(_n) // 说明当前最少需要满足几个子节点，1代表OR，n代表AND
		A.A[ID][0] = n
		A.A[ID][1] = _n
	}

	return this, ID
}

func ParsePolicyString2(A *AccessStruct, s *string, startPos int, stopPos int) (*PolicyNode, int) {
	//leftPos := startPos+1+strings.Index((*s)[startPos+1:stopPos], "(")
	this := NewPolicyNode("ThreshHold", 0) //当前节点是阈值节点，即代表AND门OR门

	A.A = append(A.A, make([]int, 2, 2))
	//_A := &(A.A[A.CurrentPointer])
	ID := A.CurrentPointer
	A.A[ID][0] = 0
	A.A[ID][1] = 0
	A.CurrentPointer++
	policy_children := make([]*PolicyNode, 0)

	var i int = startPos + 1 //startPos表示的是起始左括号
	var n int = 0
	var _n int = 0
	var leftPos int
	var rightPos int
	var trueChild string = ""

	for i <= stopPos {
		// 在当前String[startPos:stopPos]中找到第一个左括号(
		leftPos = strings.Index((*s)[i:stopPos], "(")
		// 存在左括号，说明括号内可以继续递归解析，例如(A AND B AND (C OR D))
		if leftPos != -1 {
			//左括号之前的是同层的policy，不包含括号闭合区间了
			//例子中是（C OR D）
			trueChild += (*s)[i : i+leftPos] // A AND B
			rightPos = LookForMyRightBraket(s, i+leftPos)
			tmpPolicy, tmpID := ParsePolicyString2(A, s, i+leftPos, rightPos) //继续递归（C OR D）
			policy_children = append(policy_children, tmpPolicy)
			A.A[ID] = append(A.A[ID], tmpID)
			n++
			i = rightPos + 1 // 从这一段括号后继续寻找新的括号闭合区间
		} else {
			//如果没有新的括号闭合区间，说明剩下的都是不带括号的policy，即可以认为是同层policy
			trueChild += (*s)[i:stopPos]
			break
		}
	}

	var childAttr []string
	//注意：同层policy中只会出现AND或者OR，不会两者都有
	// 存在AND门
	// 假设现在递归A AND B
	if strings.Index(trueChild, "&&") != -1 {
		childAttr = strings.Split(trueChild, "&&")

		for v := range childAttr {
			if childAttr[v] != "" {
				policy_children = append(policy_children, NewPolicyNode(childAttr[v], 1).SetMax(1).SetMin(1))
				A.A[ID] = append(A.A[ID], -A.PolicyMap[childAttr[v]])
				A.LeafID--
				n++
			}
		}
		_n = n
		this.SetOperation(1)
		// 存在OR门
		//假设现在递归：C OR D
	} else if strings.Index(trueChild, "||") != -1 {
		childAttr = strings.Split(trueChild, "||")
		//childAttr包含了policy的属性
		for v := range childAttr {
			if childAttr[v] != "" {
				//加入包含属性的叶子结点policyNode
				policy_children = append(policy_children, NewPolicyNode(childAttr[v], 1).SetMax(1).SetMin(1))
				// 加入属性对应idx的相反值
				A.A[ID] = append(A.A[ID], -A.PolicyMap[childAttr[v]])
				A.LeafID--
				n++
			}
		}
		_n = 1
		this.SetOperation(2) //当前门为OR门
	}

	if n == 0 {
		fmt.Printf("Error:: bad description. \n")
	} else {
		this.SetChildren(policy_children) //设置这个门节点的子节点
		this.SetMax(n) // 说明当前子节点个数
		this.SetMin(_n) // 说明当前最少需要满足几个子节点，1代表OR，n代表AND
		A.A[ID][0] = n
		A.A[ID][1] = _n
	}

	return this, ID
}

func LookForMyRightBraket(s *string, posL int) int {
	// 找到第一个右括号
	rightPos := posL + strings.Index((*s)[posL:], ")")

	for true {
		// 说明不存在右括号？
		if rightPos < posL {
			return -1
		} else {
			// 检查在posL和rightPos之间是否还存在左括号
			leftPos := posL + 1 + strings.Index((*s)[posL+1:rightPos], "(")
			// 存在左括号，说明现在这个rightPos可能会与这个leftPos匹配
			if leftPos >= posL+1 {
				// 查找现在这个leftPos对应的右括号
				posL = LookForMyRightBraket(s, leftPos)
				// rightPos的查找起始点posL要进行更新
				rightPos = posL + 1 + strings.Index((*s)[posL+1:], ")")
			} else {
				// 不存在多余的左括号，因此找到了匹配的右括号，返回
				return rightPos
			}
		}
	}
	return 0
}

/* Policy Parser SECTION> */

/* <Utility SECTION */
//读取文件需要经常进行错误检查，这个帮助方法可以精简下面的错误检查过程。
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

/* Utility SECTION> */

/* utils */
func CharToString(s string, t int) string {
	var sp string = ""
	for i := 0; i < t; i++ {
		sp += s
	}
	return sp
}

func GetPadding(m int, l int, depth int) string {
	var sp string = ""
	sp += CharToString("*", depth-l)
	if m == 0 {
		sp = CharToString("0", l-1) + sp
	} else {
		sp = CharToString("0", l-2) + strconv.FormatUint(uint64(m), 2) + sp
	}
	return sp[len(sp)-(depth-1):]
}

// 检查属性是否以组织/用户名称为前缀
func CheckAttrName(attrName, authorityName string) bool {
	splitN := strings.SplitN(attrName, ":", 2)
	if len(splitN) != 2 {
		return false
	}
	return authorityName == splitN[0]
}

// 根据属性名称获取组织/用户名，出错返回空字符串
func GetAuthorityNameFromAttrName(attrName string) string {
	splitN := strings.SplitN(attrName, ":", 2)
	if len(splitN) != 2 {
		return ""
	}
	return splitN[0]
}
