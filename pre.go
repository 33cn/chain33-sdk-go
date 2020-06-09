package sdk

import (
	"crypto/rand"
	"errors"
	"fmt"
	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"

	"golang.org/x/crypto/blake2b"
	"math/big"
)

var baseN = secp256k1.S256().Params().N

type KFrag struct {
	Random    string
	Value     string
	PrecurPub string
}

type ReKeyFrag struct {
	ReKeyR    string
	ReKeyU    string
	Random    string
	PrecurPub string
}

type EccPoit struct {
	x *big.Int
	y *big.Int
}

func NewEccPoint(pubStr string) (*EccPoit, error) {
	reKeyRByte, err := types.FromHex(pubStr)
	if err != nil {
		fmt.Errorf("get reKeyRByte err", err)
		return nil, err
	}
	reKeyR := crypto.PublicFromByte(reKeyRByte)
	return &EccPoit{x: reKeyR.X, y: reKeyR.Y}, nil
}

func (p *EccPoit) Add(v *EccPoit) *EccPoit {
	if v == nil {
		return p
	}

	p.x, p.y = secp256k1.S256().Add(p.x, p.y, v.x, v.y)

	return p
}

func (p *EccPoit) MulInt(i *big.Int) *EccPoit {
	p.x, p.y = secp256k1.S256().ScalarMult(p.x, p.y, i.Bytes())

	return p
}

func (p *EccPoit) ToPublicKey() *secp256k1.PublicKey {
	return &secp256k1.PublicKey{
		X: p.x,
		Y: p.y,
		Curve: secp256k1.S256(),
	}
}

func hashToModInt(digest []byte) *big.Int {
	sum := new(big.Int).SetBytes(digest)
	one := big.NewInt(1)
	order_minus_1 := big.NewInt(0)
	order_minus_1.Sub(baseN, one)

	bigNum := big.NewInt(0)
	bigNum.Mod(sum, order_minus_1).Add(bigNum, one)

	return bigNum
}

func makeShamirPolyCoeff(threshold int) []*big.Int {
	coeffs := make([]*big.Int, threshold-1)
	for i,_ := range coeffs {
		coeffs[i] = new(big.Int).SetBytes(crypto.GeneratePrivateKey())
	}

	return coeffs
}

// p0*x^2 + p1 * x + p2
func hornerPolyEval(poly []*big.Int, x *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Add(result, poly[0])
	for i := 1; i < len(poly); i++ {
		result = result.Mul(result, x).Add(result, poly[i])
	}
	return result.Mod(result, baseN)
}

func calcPart(a *big.Int, b *big.Int) *big.Int {
	p := big.NewInt(0)
	p.Sub(a, b).Mod(p, baseN)

	res := big.NewInt(0)
	res.Mul(a, p.ModInverse(p, baseN)).Mod(res, baseN)

	return res
}

func calcLambdaCoeff(inId *big.Int, selectedIds []*big.Int) *big.Int {
	var ids []*big.Int
	for _, id := range selectedIds {
		if inId.Cmp(id) == 0 {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return big.NewInt(1)
	}

	result := calcPart(ids[0], inId)
	if len(ids) > 1 {
		for _, id_j := range ids[1:] {
			result.Mul(result, calcPart(id_j, inId)).Mod(result, baseN)
		}
	}

	return result
}

func getRandomInt(order *big.Int) *big.Int {
	randInt, err := rand.Int(rand.Reader, order)
	if err != nil {
		panic(err)
	}
	return randInt
}

func GeneratePreEncryptKey(pubOwner []byte) ([]byte, string, string) {
	pubOwnerKey := crypto.PublicFromByte(pubOwner)

	priv_r := crypto.PrivateFromByte(crypto.GeneratePrivateKey())
	priv_u := crypto.PrivateFromByte(crypto.GeneratePrivateKey())

	result := &secp256k1.PublicKey{}
	result.Curve = pubOwnerKey.Curve
	sum := big.NewInt(0)
	sum.Add(priv_u.D, priv_r.D).Mod(sum, baseN)
	result.X, result.Y = secp256k1.S256().ScalarMult(pubOwnerKey.X, pubOwnerKey.Y, sum.Bytes())

	pub_r := types.ToHex((*secp256k1.PublicKey)(&priv_r.PublicKey).SerializeCompressed())
	pub_u := types.ToHex((*secp256k1.PublicKey)(&priv_u.PublicKey).SerializeCompressed())


	return result.SerializeCompressed()[1:], pub_r, pub_u
}

func GenerateKeyFragments(privOwner []byte, pubRecipient []byte, numSplit, threshold int) ([]*KFrag, error) {
	precursor := crypto.PrivateFromByte(crypto.GeneratePrivateKey())
	precurPub := types.ToHex((*secp256k1.PublicKey)(&precursor.PublicKey).SerializeCompressed())

	privOwnerKey := crypto.PrivateFromByte(privOwner)
	pubRecipientKey := crypto.PublicFromByte(pubRecipient)

	dh_Alice_poit_x := types.ECDH(precursor, pubRecipientKey)
	dAliceHash, err := blake2b.New256(precursor.X.Bytes())
	if err != nil {
		fmt.Errorf("Generate precursor Key err", err)
		return nil, err
	}
	dAliceHash.Write(pubRecipientKey.X.Bytes())
	dAliceHash.Write(dh_Alice_poit_x)
	dAlice := dAliceHash.Sum(nil)
	dAliceBN := hashToModInt(dAlice)

	// c0, c1, c2
	f0 := big.NewInt(0)
	f0.Mul(privOwnerKey.D, f0.ModInverse(dAliceBN, baseN)).Mod(f0, baseN)

	kFrags := make([]*KFrag, numSplit)
	if numSplit == 1 {
		id := getRandomInt(baseN)
		kFrags[0] = &KFrag{Random: id.String(), Value: f0.String(), PrecurPub: precurPub}
	} else {
		coeffs := makeShamirPolyCoeff(threshold)
		coeffs = append(coeffs, f0)

		// rk[i] = f2*id^2 + f1*id + f0
		for i, _ := range kFrags {
			id := getRandomInt(baseN)
			dShareHash, err := blake2b.New256(precursor.X.Bytes())
			if err != nil {
				fmt.Errorf("Generate precursor Key err", err)
				return nil, err
			}
			dShareHash.Write(pubRecipientKey.X.Bytes())
			dShareHash.Write(dh_Alice_poit_x)
			dShareHash.Write(id.Bytes())
			share := hashToModInt(dShareHash.Sum(nil))
			rk := hornerPolyEval(coeffs, share)
			kFrags[i] = &KFrag{Random: id.String(), Value: rk.String(), PrecurPub: precurPub}
		}
	}

	return kFrags, nil
}

func AssembleReencryptFragment(privRecipient []byte, reKeyFrags []*ReKeyFrag) ([]byte, error) {
	privRecipientKey := crypto.PrivateFromByte(privRecipient)
	precursor, err := types.FromHex(reKeyFrags[0].PrecurPub)
	if err != nil {
		fmt.Errorf("FromHex", err)
		return nil, err
	}
	precursorPubKey := crypto.PublicFromByte(precursor)
	dh_Bob_poit_x := types.ECDH(privRecipientKey, precursorPubKey)
	dBobHash, err := blake2b.New256(precursorPubKey.X.Bytes())
	if err != nil {
		fmt.Errorf("Generate precursor Key err", err)
		return nil, err
	}
	dBobHash.Write(privRecipientKey.X.Bytes())
	dBobHash.Write(dh_Bob_poit_x)
	dhBob := dBobHash.Sum(nil)
	dhBobBN := hashToModInt(dhBob)

	var share_key *EccPoit
	if len(reKeyFrags) == 1 {
		rPoint, err := NewEccPoint(reKeyFrags[0].ReKeyR)
		if err != nil {
			fmt.Errorf("get reKeyRByte err", err)
			return nil, err
		}
		uPoint, err := NewEccPoint(reKeyFrags[0].ReKeyU)
		if err != nil {
			fmt.Errorf("get reKeyRByte err", err)
			return nil, err
		}

		share_key = rPoint.Add(uPoint).MulInt(dhBobBN)
	} else {
		var eFinal, vFinal *EccPoit

		ids := make([]*big.Int, len(reKeyFrags))
		for x, _ := range ids {
			xs, err := blake2b.New256(precursorPubKey.X.Bytes())
			if err != nil {
				fmt.Errorf("Generate precursor Key err", err)
				return nil, err
			}
			xs.Write(privRecipientKey.X.Bytes())
			xs.Write(dh_Bob_poit_x)
			random, ret := new(big.Int).SetString(reKeyFrags[x].Random, 10)
			if !ret {
				fmt.Errorf("AssembleReencryptFragment.get value int",)
				return nil, errors.New("get big int value from keyFragment failed")
			}
			xs.Write(random.Bytes())
			ids[x] = hashToModInt(xs.Sum(nil))
		}

		for i, id := range ids {
			lambda := calcLambdaCoeff(id, ids)
			if lambda == nil {
				continue
			}
			rPoint, err := NewEccPoint(reKeyFrags[i].ReKeyR)
			if err != nil {
				fmt.Errorf("get reKeyRByte err", err)
				return nil, err
			}
			uPoint, err := NewEccPoint(reKeyFrags[i].ReKeyU)
			if err != nil {
				fmt.Errorf("get reKeyRByte err", err)
				return nil, err
			}
			e := rPoint.MulInt(lambda)
			v := uPoint.MulInt(lambda)
			eFinal = e.Add(eFinal)
			vFinal = v.Add(vFinal)
		}
		share_key = eFinal.Add(vFinal).MulInt(dhBobBN)
	}

	return share_key.ToPublicKey().SerializeCompressed(), nil
}