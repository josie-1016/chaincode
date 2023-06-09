package data

import (
	"trustPlatform/constant"
)

type BulletProofs struct {
	// couchDB使用的type
	ObjectType string `json:"docType"`

	Uid  string   `json:"uid"`
	Pid  string   `json:"pid"`
	Tags []string `json:"tags"`
	// bulletproofs
	Range      string `json:"range"`
	Commit1    Commit `json:"commit1"`
	Commit2    Commit `json:"commit2"`
	Open       string `json:"open"`
	Timestamp1 string `json:"timestamp1"`
	Timestamp2 string `json:"timestamp2"`
	//Commit1 string `json:"commit1"`
	//Commit2 string `json:"commit2"`
	//RangeProof    string `json:"bulletProof"`
	//ProofFileName string `json:"proofFileName"`
	ProofPre string `json:"proofpre"`
	Proof    string `json:"proof"`
}

type ECPoint struct {
	X string `json:"X"`
	Y string `json:"Y"`
}

type Commit struct {
	Comm ECPoint `json:"Comm"`
}

func NewBulletProofs(uid, pid string, tags []string, timestamp1 string, timestamp2 string, ran string, commit1, commit2 Commit, open string, proof string, proofpre string) *BulletProofs {
	return &BulletProofs{
		ObjectType: constant.BulletProofs,
		Uid:        uid,
		Pid:        pid,
		Tags:       tags,
		Timestamp1: timestamp1,
		Timestamp2: timestamp2,
		Range:      ran,
		Commit1:    commit1,
		Commit2:    commit2,
		Open:       open,
		//ProofFileName: prooffilename,
		ProofPre: proofpre,
		Proof:    proof,
	}
}
