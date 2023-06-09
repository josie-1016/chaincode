package SM

import (
	"encoding/json"
	"log"
)

func VerifySM2(paramJson, sign, publicKey, uid string) (b bool, err error) {
	pub := new(PubKey)
	if err = json.Unmarshal([]byte(publicKey), pub); err != nil {
		log.Println(err.Error())
		return false, err
	}
	pubKey := new(PublicKey)
	pubKey.X = &pub.X
	pubKey.Y = &pub.Y
	pubKey.Curve = GetSm2P256V1()

	Sign := new(SignaTure)
	if err = json.Unmarshal([]byte(sign), Sign); err != nil {
		log.Println(err.Error())
		return false, err
	}
	signature := new(Sm2Signature)
	signature.R = &Sign.R
	signature.S = &Sign.S

	b, err = Verify(pubKey, []byte(uid), []byte(paramJson), signature)
	return
}
