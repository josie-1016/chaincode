package rangproof

import (
	"crypto/rand"
	"log"
	"math/big"
)

type PedersenCommit struct {
	//G ECPoint // G value for commitments of a single value
	//H ECPoint // H value for commitments of a single value
	Comm ECPoint
}

func PerdersenCommit(v *big.Int, gamma *big.Int) (*PedersenCommit, *big.Int, error) {
	commit := new(PedersenCommit)
	if gamma == nil {
		gamma, _ = rand.Int(rand.Reader, EC.N)
	}
	commit.Comm = EC.G.Mult(v).Add(EC.H.Mult(gamma)) //承诺V
	return commit, gamma, nil
}

func PedersenSubNum(p *PedersenCommit, v *big.Int) (*PedersenCommit, error) {
	commit := new(PedersenCommit)
	commit.Comm = p.Comm.Add(EC.G.Mult(v).Neg())
	return commit, nil
}

func PedersenAddNum(p *PedersenCommit, v *big.Int) (*PedersenCommit, error) {
	commit := new(PedersenCommit)
	commit.Comm = p.Comm.Add(EC.G.Mult(v))
	return commit, nil
}
func PedersenBatchAddCommit(commits []*PedersenCommit) (*PedersenCommit, error) {
	commit := new(PedersenCommit)
	for i := range commits {
		if i == 0 {
			commit.Comm = commits[i].Comm
		} else {
			commit.Comm = commit.Comm.Add(commits[i].Comm)
		}
		log.Println(commit.Comm)
	}
	return commit, nil
}
func PedersenAddCommit(p *PedersenCommit, com *PedersenCommit) (*PedersenCommit, error) {
	commit := new(PedersenCommit)
	commit.Comm = p.Comm.Add(com.Comm)
	return commit, nil
}

func PedersenSubCommit(p *PedersenCommit, com *PedersenCommit) (*PedersenCommit, error) {
	commit := new(PedersenCommit)
	commit.Comm = p.Comm.Add(com.Comm.Neg())
	return commit, nil
}
