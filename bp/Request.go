package main

import "bullet/rangproof"

type RangeProofResponse struct {
	Commit1 *rangproof.PedersenCommit
	Commit2 *rangproof.PedersenCommit
	Proof   *rangproof.RangeProof
}

type VerifyRequest struct {
	Range   string
	Commit1 string
	Commit2 string
	Proof   string
}
