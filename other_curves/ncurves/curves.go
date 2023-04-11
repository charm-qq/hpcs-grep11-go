package ncurves

import (
	"crypto/elliptic"
	"math/big"
)

func P192() elliptic.Curve {
	initP192()
	return secp192r1
}

var secp192r1 *elliptic.CurveParams

func initP192() {
	secp192r1 = &elliptic.CurveParams{Name: "secp192r1"}
	secp192r1.P, _ = new(big.Int).SetString("6277101735386680763835789423207666416083908700390324961279", 10)
	secp192r1.N, _ = new(big.Int).SetString("6277101735386680763835789423176059013767194773182842284081", 10)
	secp192r1.B, _ = new(big.Int).SetString("64210519e59c80e70fa7e9ab72243049feb8deecc146b9b1", 16)
	secp192r1.Gx, _ = new(big.Int).SetString("188da80eb03090f67cbf20eb43a18800f4ff0afd82ff1012", 16)
	secp192r1.Gy, _ = new(big.Int).SetString("07192b95ffc8da78631011ed6b24cdd573f977a11e794811", 16)
	secp192r1.BitSize = 192
}
