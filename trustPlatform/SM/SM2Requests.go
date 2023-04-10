package SM

import "math/big"

type PubKey struct {
	X big.Int `json:"x"`
	Y big.Int `json:"y"`
}

type SignaTure struct {
	R big.Int `json:"r"`
	S big.Int `json:"s"`
}
