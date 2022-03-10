// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bls12377

//Note: This only works for simple extensions

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fp"
	"math/big"
)

func g1IsogenyXNumerator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		false,
		[]fp.Element{
			{9381318728011785451, 8795417190580748876, 15171640721257608922, 11815547924113428908, 15499908520243100994, 75408755324413256},
			{12414498063752772717, 9915153185132073893, 5598625970987438951, 3342254783599619135, 3349592178919125510, 9993871847068096},
			{4662210776746950618, 10687085762534440940, 7484820859645808636, 2221301482234255553, 10609677459585442106, 9950135580589350},
		},
		x)
}

func g1IsogenyXDenominator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		true,
		[]fp.Element{
			{12764504107591987636, 2767124593109192342, 3947759810240204190, 13369019134398476541, 13398368715676502040, 39975487388272384},
		},
		x)
}

func g1IsogenyYNumerator(dst *fp.Element, x *fp.Element, y *fp.Element) {
	var _dst fp.Element
	g1EvalPolynomial(&_dst,
		false,
		[]fp.Element{
			{13844135623281082635, 637899392157745290, 5176720401210677272, 4780940929980393029, 13803251044890140836, 51447363642369244},
			{512010462697120695, 609509684909242946, 13763343875136563934, 2839514380057330869, 15407015190976871917, 114223893455203604},
			{14191436515319700132, 6479619458373647736, 9513056055282499867, 15178407828209519654, 12166396751953702822, 75539964123849493},
			{2331105388373475309, 5343542881267220470, 12965782466677680126, 1110650741117127776, 5304838729792721053, 4975067790294675},
		},
		x)

	dst.Mul(&_dst, y)
}

func g1IsogenyYDenominator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		true,
		[]fp.Element{
			{8694832399336342723, 13482963304561246841, 6984108042366343277, 8355250559073919616, 16937021447778317421, 44890599540624877},
			{1100361703846424922, 5005767817281133373, 917019320419705433, 14251746270386956490, 5522097789867984932, 4443041874334878},
			{1400024175356859676, 8301373779327577028, 11843279430720612570, 3213569255776326391, 3301617999610402890, 119926462164817154},
		},
		x)
}

func g1Isogeny(p *G1Affine) {

	den := make([]fp.Element, 2)

	g1IsogenyYDenominator(&den[1], &p.X)
	g1IsogenyXDenominator(&den[0], &p.X)

	g1IsogenyYNumerator(&p.Y, &p.X, &p.Y)
	g1IsogenyXNumerator(&p.X, &p.X)

	den = fp.BatchInvert(den)

	p.X.Mul(&p.X, &den[0])
	p.Y.Mul(&p.Y, &den[1])
}

// g1SqrtRatio computes the square root of u/v and returns 0 iff u/v was indeed a quadratic residue
// if not, we get sqrt(Z * u / v). Recall that Z is non-residue
// The main idea is that since the computation of the square root involves taking large powers of u/v, the inversion of v can be avoided
func g1SqrtRatio(z *fp.Element, u *fp.Element, v *fp.Element) uint64 {

	// Taken from https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/ F.2.1.1. for any field

	tv1 := fp.Element{1558150277696978216, 2889962898601943991, 2027733260451875201, 1491536148589418669, 4860991584141488473, 78837656424301472} //tv1 = c6

	var tv2, tv3, tv4, tv5 fp.Element
	var exp big.Int
	// c4 = 70368744177663 = 2^46 - 1
	// q is odd so c1 is at least 1.
	exp.SetBytes([]byte{63, 255, 255, 255, 255, 255})

	tv2.Exp(*v, &exp)
	tv3.Mul(&tv2, &tv2)
	tv3.Mul(&tv3, v)

	// line 5
	tv5.Mul(u, &tv3)

	// c3 = 1837921289030710838195067919506396475074392872918698035817074744121558668640693829665401097909504529
	exp.SetBytes([]byte{3, 92, 116, 140, 47, 138, 33, 213, 140, 118, 11, 128, 217, 66, 146, 118, 52, 69, 179, 230, 1, 234, 39, 30, 61, 230, 196, 95, 116, 18, 144, 0, 46, 22, 186, 136, 96, 0, 0, 1, 10, 17})
	tv5.Exp(tv5, &exp)
	tv5.Mul(&tv5, &tv2)
	tv2.Mul(&tv5, v)
	tv3.Mul(&tv5, u)

	// line 10
	tv4.Mul(&tv3, &tv2)

	// c5 = 35184372088832
	exp.SetBytes([]byte{32, 0, 0, 0, 0, 0})
	tv5.Exp(tv4, &exp)

	isQNr := g1NotOne(&tv5)

	tv2.Mul(&tv3, &fp.Element{9851577832091164681, 11031556238464483940, 10838960500261981046, 11853033281081062456, 6403263430752281921, 18823706241905246})
	tv5.Mul(&tv4, &tv1)

	// line 15

	tv3.Select(int(isQNr), &tv3, &tv2)
	tv4.Select(int(isQNr), &tv4, &tv5)

	exp.Lsh(big.NewInt(1), 46-2)

	for i := 46; i >= 2; i-- {
		//line 20
		tv5.Exp(tv4, &exp)
		nE1 := g1NotOne(&tv5)

		tv2.Mul(&tv3, &tv1)
		tv1.Mul(&tv1, &tv1)
		tv5.Mul(&tv4, &tv1)

		tv3.Select(int(nE1), &tv3, &tv2)
		tv4.Select(int(nE1), &tv4, &tv5)

		exp.Rsh(&exp, 1)
	}

	*z = tv3
	return isQNr
}

/*
// g1SetZ sets z to [2].
func g1SetZ(z *fp.Element) {
    z.Set( &fp.Element  { 404198066556501712, 11709709805437321058, 4538334656037814244, 17770411857874044427, 11090443381845330384, 79601084644714804 } )
}*/

// g1MulByZ multiplies x by [2] and stores the result in z
func g1MulByZ(z *fp.Element, x *fp.Element) {

	res := *x

	res.Double(&res)

	*z = res
}

// From https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/ Pg 80
func g1SswuMap(u *fp.Element) G1Affine {

	var tv1 fp.Element
	tv1.Square(u)

	//mul tv1 by Z
	g1MulByZ(&tv1, &tv1)

	var tv2 fp.Element
	tv2.Square(&tv1)
	tv2.Add(&tv2, &tv1)

	var tv3 fp.Element
	//Standard doc line 5
	var tv4 fp.Element
	tv4.SetOne()
	tv3.Add(&tv2, &tv4)
	tv3.Mul(&tv3, &fp.Element{11130294635325289193, 6502679372128844082, 15863297759487624914, 16270683149854112145, 3560014356538878812, 27923742146399959})

	tv2NZero := g1NotZero(&tv2)

	// tv4 = Z
	tv4 = fp.Element{404198066556501712, 11709709805437321058, 4538334656037814244, 17770411857874044427, 11090443381845330384, 79601084644714804}

	tv2.Neg(&tv2)
	tv4.Select(int(tv2NZero), &tv4, &tv2)
	tv2 = fp.Element{17252667382019449424, 8408110001211059699, 18415587021986261264, 10797086888535946954, 9462758283094809199, 54995354010328751}
	tv4.Mul(&tv4, &tv2)

	tv2.Square(&tv3)

	var tv6 fp.Element
	//Standard doc line 10
	tv6.Square(&tv4)

	var tv5 fp.Element
	tv5.Mul(&tv6, &fp.Element{17252667382019449424, 8408110001211059699, 18415587021986261264, 10797086888535946954, 9462758283094809199, 54995354010328751})

	tv2.Add(&tv2, &tv5)
	tv2.Mul(&tv2, &tv3)
	tv6.Mul(&tv6, &tv4)

	//Standards doc line 15
	tv5.Mul(&tv6, &fp.Element{11130294635325289193, 6502679372128844082, 15863297759487624914, 16270683149854112145, 3560014356538878812, 27923742146399959})
	tv2.Add(&tv2, &tv5)

	var x fp.Element
	x.Mul(&tv1, &tv3)

	var y1 fp.Element
	gx1NSquare := g1SqrtRatio(&y1, &tv2, &tv6)

	var y fp.Element
	y.Mul(&tv1, u)

	//Standards doc line 20
	y.Mul(&y, &y1)

	x.Select(int(gx1NSquare), &tv3, &x)
	y.Select(int(gx1NSquare), &y1, &y)

	y1.Neg(&y)
	y.Select(int(g1Sgn0(u)^g1Sgn0(&y)), &y, &y1)

	//Standards doc line 25
	x.Div(&x, &tv4)

	return G1Affine{x, y}
}

// EncodeToCurveG1SSWU maps a fp.Element to a point on the curve using the Simplified Shallue and van de Woestijne Ulas map
//https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/#section-6.6.3
func EncodeToCurveG1SSWU(msg, dst []byte) (G1Affine, error) {

	var res G1Affine
	u, err := hashToFp(msg, dst, 1)
	if err != nil {
		return res, err
	}

	res = g1SswuMap(&u[0])

	//this is in an isogenous curve
	g1Isogeny(&res)

	res.ClearCofactor(&res)

	return res, nil
}

// HashToCurveG1SSWU hashes a byte string to the G1 curve. Usable as a random oracle.
// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#section-3
func HashToCurveG1SSWU(msg, dst []byte) (G1Affine, error) {
	u, err := hashToFp(msg, dst, 2*1)
	if err != nil {
		return G1Affine{}, err
	}

	Q0 := g1SswuMap(&u[0])
	Q1 := g1SswuMap(&u[1])

	//TODO: Add in E' first, then apply isogeny
	g1Isogeny(&Q0)
	g1Isogeny(&Q1)

	var _Q0, _Q1 G1Jac
	_Q0.FromAffine(&Q0)
	_Q1.FromAffine(&Q1).AddAssign(&_Q0)

	_Q1.ClearCofactor(&_Q1)

	Q1.FromJacobian(&_Q1)
	return Q1, nil
}

// g1Sgn0 is an algebraic substitute for the notion of sign in ordered fields
// Namely, every non-zero quadratic residue in a finite field of characteristic =/= 2 has exactly two square roots, one of each sign
// Taken from https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/ section 4.1
// The sign of an element is not obviously related to that of its Montgomery form
func g1Sgn0(z *fp.Element) uint64 {

	nonMont := *z
	nonMont.FromMont()

	return nonMont[0] % 2

}

func g1EvalPolynomial(z *fp.Element, monic bool, coefficients []fp.Element, x *fp.Element) {
	dst := coefficients[len(coefficients)-1]

	if monic {
		dst.Add(&dst, x)
	}

	for i := len(coefficients) - 2; i >= 0; i-- {
		dst.Mul(&dst, x)
		dst.Add(&dst, &coefficients[i])
	}

	z.Set(&dst)
}

func g1NotZero(x *fp.Element) uint64 {

	return x[0] | x[1] | x[2] | x[3] | x[4] | x[5]

}

func g1NotOne(x *fp.Element) uint64 {

	var one fp.Element
	return one.SetOne().NotEqual(x)

}
