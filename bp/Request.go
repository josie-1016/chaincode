package main

import (
	"bullet/rangproof"
	"math/big"
)

type RangeProofResponse struct {
	Commit1 *rangproof.PedersenCommit
	Commit2 *rangproof.PedersenCommit
	Proof   *rangproof.RangeProof
}
type CommitResponse struct {
	Proof   *rangproof.RangeProof
	Commit1 *rangproof.PedersenCommit
	Open    *big.Int
}
type VerifyRequest struct {
	Range    string
	Commit1  *rangproof.PedersenCommit
	Commit2  *rangproof.PedersenCommit
	Proof    *rangproof.RangeProof
	ProofPre *rangproof.RangeProof
}

type CreateBatchRequest struct {
	Range   string
	Commits []*rangproof.PedersenCommit
	Open    string
	Value   string
}
