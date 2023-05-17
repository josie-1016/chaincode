package rangproof

import (
	"fmt"
	"math/big"
)

func TwoVectorPCommitWithGens(G, H []ECPoint, a, b []*big.Int) ECPoint {
	if len(G) != len(H) || len(G) != len(a) || len(a) != len(b) {
		fmt.Println("TwoVectorPCommitWithGens: Uh oh! Arrays not of the same length")
		fmt.Printf("len(G): %d\n", len(G))
		fmt.Printf("len(H): %d\n", len(H))
		fmt.Printf("len(a): %d\n", len(a))
		fmt.Printf("len(b): %d\n", len(b))
	}

	commitment := Zero()

	for i := 0; i < len(G); i++ {
		modA := new(big.Int).Mod(a[i], EC.N)
		modB := new(big.Int).Mod(b[i], EC.N)

		commitment = commitment.Add(G[i].Mult(modA)).Add(H[i].Mult(modB))
	}

	return commitment
}
