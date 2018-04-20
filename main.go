package main

import (
	"github.com/evgeniy-scherbina/ecdsa/point"
	"github.com/evgeniy-scherbina/ecdsa/keys"
	"fmt"
)

// z - набор данных
// n - порядок
// G - базовую точку
// d - закрытый ключ
func main() {
	keyPair := keys.NewECDSAKeyPair(2, &point.TestEllipticCurveParams)
	sig := keyPair.Sign(17, &point.TestEllipticCurveParams)
	fmt.Println(sig)
	ok := keyPair.Verify(17, &point.TestEllipticCurveParams)
	fmt.Println("VERIFY:", ok)

	// privateKey := 2
	// publicKey := point.TestEllipticCurveParams.BasePoint.AddSelf()
}