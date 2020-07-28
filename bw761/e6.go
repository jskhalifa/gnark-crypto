// Copyright 2020 ConsenSys AG
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

package bw761

import "math/big"

// E6 is a degree-three finite field extension of fp2
type E6 struct {
	B0, B1, B2 E2
}

// Equal returns true if z equals x, fasle otherwise
func (z *E6) Equal(x *E6) bool {
	return z.B0.Equal(&x.B0) && z.B1.Equal(&x.B1) && z.B2.Equal(&x.B2)
}

// SetString sets a E6 elmt from stringf
func (z *E6) SetString(s1, s2, s3, s4, s5, s6 string) *E6 {
	z.B0.SetString(s1, s2)
	z.B1.SetString(s3, s4)
	z.B2.SetString(s5, s6)
	return z
}

// Set Sets a E6 elmt form another E6 elmt
func (z *E6) Set(x *E6) *E6 {
	z.B0 = x.B0
	z.B1 = x.B1
	z.B2 = x.B2
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *E6) SetOne() *E6 {
	z.B0.A0.SetOne()
	z.B0.A1.SetZero()
	z.B1.A0.SetZero()
	z.B1.A1.SetZero()
	z.B2.A0.SetZero()
	z.B2.A1.SetZero()
	return z
}

// SetRandom set z to a random elmt
func (z *E6) SetRandom() *E6 {
	z.B0.SetRandom()
	z.B1.SetRandom()
	z.B2.SetRandom()
	return z
}

// ToMont converts to Mont form
func (z *E6) ToMont() *E6 {
	z.B0.ToMont()
	z.B1.ToMont()
	z.B2.ToMont()
	return z
}

// FromMont converts from Mont form
func (z *E6) FromMont() *E6 {
	z.B0.FromMont()
	z.B1.FromMont()
	z.B2.FromMont()
	return z
}

// Add adds two elements of E6
func (z *E6) Add(x, y *E6) *E6 {
	z.B0.Add(&x.B0, &y.B0)
	z.B1.Add(&x.B1, &y.B1)
	z.B2.Add(&x.B2, &y.B2)
	return z
}

// Neg negates the E6 number
func (z *E6) Neg(x *E6) *E6 {
	z.B0.Neg(&x.B0)
	z.B1.Neg(&x.B1)
	z.B2.Neg(&x.B2)
	return z
}

// Sub two elements of E6
func (z *E6) Sub(x, y *E6) *E6 {
	z.B0.Sub(&x.B0, &y.B0)
	z.B1.Sub(&x.B1, &y.B1)
	z.B2.Sub(&x.B2, &y.B2)
	return z
}

// Double doubles an element in E6
func (z *E6) Double(x *E6) *E6 {
	z.B0.Double(&x.B0)
	z.B1.Double(&x.B1)
	z.B2.Double(&x.B2)
	return z
}

// String puts E6 elmt in string form
func (z *E6) String() string {
	return (z.B0.String() + "+(" + z.B1.String() + ")*v+(" + z.B2.String() + ")*v**2")
}

// Mul sets z to the E6-product of x,y, returns z
func (z *E6) Mul(x, y *E6) *E6 {
	// Algorithm 13 from https://eprint.iacr.org/2010/354.pdf
	var t0, t1, t2, c0, c1, c2, tmp E2
	t0.Mul(&x.B0, &y.B0)
	t1.Mul(&x.B1, &y.B1)
	t2.Mul(&x.B2, &y.B2)

	c0.Add(&x.B1, &x.B2)
	tmp.Add(&y.B1, &y.B2)
	c0.Mul(&c0, &tmp).Sub(&c0, &t1).Sub(&c0, &t2).MulByNonResidue(&c0).Add(&c0, &t0)

	c1.Add(&x.B0, &x.B1)
	tmp.Add(&y.B0, &y.B1)
	c1.Mul(&c1, &tmp).Sub(&c1, &t0).Sub(&c1, &t1)
	tmp.MulByNonResidue(&t2)
	c1.Add(&c1, &tmp)

	tmp.Add(&x.B0, &x.B2)
	c2.Add(&y.B0, &y.B2).Mul(&c2, &tmp).Sub(&c2, &t0).Sub(&c2, &t2).Add(&c2, &t1)

	z.B0.Set(&c0)
	z.B1.Set(&c1)
	z.B2.Set(&c2)

	return z
}

// MulAssign sets z to the E6-product of z,y, returns z
func (z *E6) MulAssign(x *E6) *E6 {
	z.Mul(z, x)
	return z
}

// MulByE2 multiplies x by an elements of E2
func (z *E6) MulByE2(x *E6, y *E2) *E6 {
	var yCopy E2
	yCopy.Set(y)
	z.B0.Mul(&x.B0, &yCopy)
	z.B1.Mul(&x.B1, &yCopy)
	z.B2.Mul(&x.B2, &yCopy)
	return z
}

// Square sets z to the E6-product of x,x, returns z
func (z *E6) Square(x *E6) *E6 {

	// Algorithm 16 from https://eprint.iacr.org/2010/354.pdf
	var c4, c5, c1, c2, c3, c0 E2
	c4.Mul(&x.B0, &x.B1).Double(&c4)
	c5.Square(&x.B2)
	c1.MulByNonResidue(&c5).Add(&c1, &c4)
	c2.Sub(&c4, &c5)
	c3.Square(&x.B0)
	c4.Sub(&x.B0, &x.B1).Add(&c4, &x.B2)
	c5.Mul(&x.B1, &x.B2).Double(&c5)
	c4.Square(&c4)
	c0.MulByNonResidue(&c5).Add(&c0, &c3)
	z.B2.Add(&c2, &c4).Add(&z.B2, &c5).Sub(&z.B2, &c3)
	z.B0.Set(&c0)
	z.B1.Set(&c1)

	return z
}

// CyclotomicSquare https://eprint.iacr.org/2009/565.pdf, 3.2
func (z *E6) CyclotomicSquare(x *E6) *E6 {

	var res, a E6
	var tmp E2

	// A
	res.B0.Square(&x.B0)
	a.B0.Conjugate(&x.B0)

	// B
	res.B2.A0.Set(&x.B1.A1)
	res.B2.A1.MulByNonResidueInv(&x.B1.A0)
	res.B2.Square(&res.B2).Double(&res.B2).Double(&res.B2).Neg(&res.B2)
	a.B2.Conjugate(&x.B2)

	// C
	tmp.Square(&x.B2)
	res.B1.A0.MulByNonResidue(&tmp.A1)
	res.B1.A1.Set(&tmp.A0)
	a.B1.A0.Neg(&x.B1.A0)
	a.B1.A1.Set(&x.B1.A1)

	z.Sub(&res, &a).Double(z).Add(z, &res)

	return z
}

// Inverse an element in E6
func (z *E6) Inverse(x *E6) *E6 {
	// Algorithm 17 from https://eprint.iacr.org/2010/354.pdf
	// step 9 is wrong in the paper it's t1-t4
	var t0, t1, t2, t3, t4, t5, t6, c0, c1, c2, d1, d2 E2
	t0.Square(&x.B0)
	t1.Square(&x.B1)
	t2.Square(&x.B2)
	t3.Mul(&x.B0, &x.B1)
	t4.Mul(&x.B0, &x.B2)
	t5.Mul(&x.B1, &x.B2)
	c0.MulByNonResidue(&t5).Neg(&c0).Add(&c0, &t0)
	c1.MulByNonResidue(&t2).Sub(&c1, &t3)
	c2.Sub(&t1, &t4)
	t6.Mul(&x.B0, &c0)
	d1.Mul(&x.B2, &c1)
	d2.Mul(&x.B1, &c2)
	d1.Add(&d1, &d2).MulByNonResidue(&d1)
	t6.Add(&t6, &d1)
	t6.Inverse(&t6)
	z.B0.Mul(&c0, &t6)
	z.B1.Mul(&c1, &t6)
	z.B2.Mul(&c2, &t6)

	return z
}

// Exp sets z=x**e and returns it
func (z *E6) Exp(x *E6, e big.Int) *E6 {
	var res E6
	res.SetOne()
	b := e.Bytes()
	for i := range b {
		w := b[i]
		mask := byte(0x80)
		for j := 0; j < 8; j++ {
			res.Square(&res)
			if (w&mask)>>(7-j) != 0 {
				res.Mul(&res, x)
			}
			mask = mask >> 1
		}
	}
	z.Set(&res)
	return z
}
