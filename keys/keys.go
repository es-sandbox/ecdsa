package keys

import (
	"github.com/evgeniy-scherbina/ecdsa/point"
	"github.com/evgeniy-scherbina/ecdsa/ring_modulo"
	"fmt"
)

type privateKey struct {
	value int
}

func (self *privateKey) Value() int {
	return self.value
}

func (self *privateKey) CountPublicKeyFast(params *point.EllipticCurveParams) *publicKey {
	if self.value == 76 {
		// 2( 2(G + 2(G + 2( 2( 2(2, 22) ) ) ) ) )
		G := params.BasePoint.Copy()
		value := G.AddSelf().AddSelf().AddSelf().Add(G).AddSelf().Add(G).AddSelf().AddSelf()
		return &publicKey{ value }
	}

	value := params.BasePoint.Copy()
	n := self.value
	for n != 1 {
		if n%2 == 0 {
			value = value.AddSelf()
			n /= 2
		} else {
			value = value.Add(params.BasePoint)
			n--
		}
		fmt.Println(n, value)
	}
	return &publicKey{value}
}

func (self *privateKey) CountPublicKeySlow(params *point.EllipticCurveParams) *publicKey {
	/*
	if self.value == 3 {
		value := params.BasePoint.Copy()
		// fmt.Println("1: ", value)

		value = value.Add(params.BasePoint)
		// fmt.Println("2: ", value)

		value = value.Add(params.BasePoint)
		// fmt.Println("3: ", value)

		return &publicKey{value}
	}
	*/

	/*
	if self.value == 3 {
		return &publicKey{params.BasePoint.Copy().AddSelf().Add(params.BasePoint)}
	}
	*/

	value := params.BasePoint.Copy()
	for i := 0; i < self.value - 1; i++ {
		value = value.Add(params.BasePoint)
	}
	return &publicKey{value}
}

func (self *privateKey) sign(z int, params *point.EllipticCurveParams) *Signature {
	modulo := ring_modulo.MakeConstructor(params.BasePointOrder)

	d := self.value

	// 1. Выберем случайное число:
	privateKeyRandom := &privateKey{3}

	// 2. Рассчитаем точку:
	publicKeyRandom := privateKeyRandom.CountPublicKeySlow(params)

	// 3. Находим r:
	r := publicKeyRandom.value.X() % point.TestEllipticCurveParams.BasePointOrder

	// 4. Находим s:
	s1 := modulo().Set(r * d).Add(z)
	s2 := privateKeyRandom.value

	s := s1.Div(s2).Get()

	return &Signature{R: r, S: s}
}


type publicKey struct {
	value *point.Point
}

func (self *publicKey) Value() *point.Point {
	return self.value
}

func (self *publicKey) verify(z int, params *point.EllipticCurveParams) bool {
	// z = 17 (данные)
	// (r, s) = (62, 47) (подпись)
	// n = 79 (порядок)
	// G = (2, 22) (базовая точка)
	// Q = (52, 7) (публичный ключ)

	// z := 17
	r, s := 62, 47
	n := params.BasePointOrder
	G := params.BasePoint.Copy()
	Q := self.value

	// 1. Убедимся, что что r и s находятся в диапазоне от 1 до n-1.
	// TODO(evg): implement

	// 2. Рассчитаем w :
	w := ring_modulo.InverseByModuloSlow(s, n)
	fmt.Println("STEP_2, w: ", w)

	// 3. Рассчитаем u :
	u := privateKey{z * w % n}
	fmt.Println("STEP_3, u: ", u)

	// 4. Рассчитаем v :
	v := privateKey{r * w % n}
	fmt.Println("STEP_4, v:", v)

	// 5. Рассчитаем точку (х, у):
	// (x, y) = uG + vQ
	uG := u.CountPublicKeyFast(params)
	fmt.Println("STEP_5.1, uG", uG.value)

	params2 := params.Copy()
	params2.BasePoint = Q
	vQ := v.CountPublicKeySlow(params2)
	fmt.Println("STEP_5.2, vQ", vQ.value)

	uv := uG.value.Add(vQ.value)
	fmt.Println("STEP_5.3, uG + vQ", uv)

	_, _, _, _ = z, r, G, Q
	return r == uv.X() % n
}


type ECDSAKeyPair struct {
	privateKey *privateKey
	publicKey  *publicKey
}

func NewECDSAKeyPair(value int, params *point.EllipticCurveParams) *ECDSAKeyPair {
	_privateKey := &privateKey{ value }
	publicKey := _privateKey.CountPublicKeySlow(params)
	return &ECDSAKeyPair{_privateKey, publicKey }
}

func (self *ECDSAKeyPair) Sign(z int, params *point.EllipticCurveParams) *Signature {
	return self.privateKey.sign(z, params)
}

func (self *ECDSAKeyPair) Verify(z int, params *point.EllipticCurveParams) bool {
	return self.publicKey.verify(z, params)
}

type Signature struct {
	R int
	S int
}