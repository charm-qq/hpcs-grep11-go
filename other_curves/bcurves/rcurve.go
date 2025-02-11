package bcurves

import (
	"crypto/elliptic"
	"encoding/asn1"
	"math/big"
)

type rcurve struct {
	twisted elliptic.Curve
	params  elliptic.CurveParams
	z       *big.Int
	zinv    *big.Int
	z2      *big.Int
	z3      *big.Int
	zinv2   *big.Int
	zinv3   *big.Int
}

var (
	// The following variables are standardized elliptic curve definitions
	// refer http://oid-info.com/cgi-bin/display?oid=brainpoolP160t1&submit=Display&action=display
	OIDNamedCurveBP160T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 2}
	OIDNamedCurveBP192T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 4}
	OIDNamedCurveBP224T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 6}
	OIDNamedCurveBP256T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 8}
	OIDNamedCurveBP320T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 10}
	OIDNamedCurveBP384T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 12}
	OIDNamedCurveBP512T1 = asn1.ObjectIdentifier{1, 3, 36, 3, 3, 2, 8, 1, 1, 14}
)

func newrcurve(twisted elliptic.Curve, gx, gy, z *big.Int) *rcurve {
	var curve rcurve

	curve.twisted = twisted
	curve.params = *twisted.Params()
	curve.params.B = nil // FIXME: crypto/elliptic assumes A=-3
	curve.params.Gx = gx
	curve.params.Gy = gy

	two := big.NewInt(2)
	three := big.NewInt(3)

	curve.z = z
	curve.zinv = new(big.Int).ModInverse(z, curve.params.P)
	curve.z2 = new(big.Int).Exp(curve.z, two, curve.params.P)
	curve.z3 = new(big.Int).Exp(curve.z, three, curve.params.P)
	curve.zinv2 = new(big.Int).Exp(curve.zinv, two, curve.params.P)
	curve.zinv3 = new(big.Int).Exp(curve.zinv, three, curve.params.P)

	// auto rename twisted curves to r1 from t1
	name := curve.params.Name
	curve.params.Name = name[:len(name)-2] + "r1"

	return &curve
}

func (curve *rcurve) toTwisted(x, y *big.Int) (*big.Int, *big.Int) {
	var tx, ty big.Int
	tx.Mul(x, curve.z2)
	tx.Mod(&tx, curve.params.P)
	ty.Mul(y, curve.z3)
	ty.Mod(&ty, curve.params.P)
	return &tx, &ty
}

func (curve *rcurve) fromTwisted(tx, ty *big.Int) (*big.Int, *big.Int) {
	var x, y big.Int
	x.Mul(tx, curve.zinv2)
	x.Mod(&x, curve.params.P)
	y.Mul(ty, curve.zinv3)
	y.Mod(&y, curve.params.P)
	return &x, &y
}

func (curve *rcurve) Params() *elliptic.CurveParams {
	return &curve.params
}

func (curve *rcurve) IsOnCurve(x, y *big.Int) bool {
	return curve.twisted.IsOnCurve(curve.toTwisted(x, y))
}

func (curve *rcurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	tx1, ty1 := curve.toTwisted(x1, y1)
	tx2, ty2 := curve.toTwisted(x2, y2)
	return curve.fromTwisted(curve.twisted.Add(tx1, ty1, tx2, ty2))
}

func (curve *rcurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	return curve.fromTwisted(curve.twisted.Double(curve.toTwisted(x1, y1)))
}

func (curve *rcurve) ScalarMult(x1, y1 *big.Int, scalar []byte) (x, y *big.Int) {
	tx1, ty1 := curve.toTwisted(x1, y1)
	return curve.fromTwisted(curve.twisted.ScalarMult(tx1, ty1, scalar))
}

func (curve *rcurve) ScalarBaseMult(scalar []byte) (x, y *big.Int) {
	return curve.fromTwisted(curve.twisted.ScalarBaseMult(scalar))
}
