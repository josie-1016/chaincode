package main

import (
	DecentralizedABE "github.com/vangogo/tree/ThresholdABE"

	"github.com/Nik-U/pbc"
)

type GenerateOPKRequest struct {
	UserNames  []string
	PartPkList []*pbc.Element `field:"2"`
	N          int
	T          int
}
type GenerateOSKRequest struct {
	UserNames  []string
	PartPkList []*pbc.Element `field:"3"`
	N          int
	T          int
}
type GenerateAPKRequest struct {
	UserNames  []string
	PartPkList []*pbc.Element `field:"0"`
	N          int
	T          int
	AttrName   string
}

type EncryptRequest struct {
	PlainContent string
	Policy       string
	AuthorityMap map[string]*Authority
}

type ThreholdEncryptRequest struct {
	PlainContent string
	PubKey       []byte
}

type Authority struct {
	PK     *pbc.Element `field:"2"`
	APKMap map[string]*DecentralizedABE.APK
}

type DecryptRequest struct {
	Cipher  string
	AttrMap map[string]*pbc.Element `field:"0"`
	Uid     string
}

type ThreholdDecryptRequest struct {
	Cipher        []byte
	ThresholdPriv []byte
}

func (a *Authority) GetPK() *pbc.Element {
	return a.PK
}

func (a *Authority) GetAPKMap() map[string]*DecentralizedABE.APK {
	return a.APKMap
}
