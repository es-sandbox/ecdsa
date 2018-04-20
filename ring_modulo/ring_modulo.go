package ring_modulo

func MakeConstructor(modulo int) func() *ringModulo {
	return func() *ringModulo {
		return New(modulo)
	}
}

type ringModulo struct {
	modulo int
	value int
}

func New(modulo int) *ringModulo {
	return &ringModulo{ modulo, 0 }
}

func (self *ringModulo) Set(value int) *ringModulo {
	self.value = value
	return self.normalizeSelf()
}

func (self *ringModulo) Get() int {
	return self.value
}

func (self *ringModulo) Add(other int) *ringModulo {
	other = self.normalizeValue(other)
	self.value += other
	return self.normalizeSelf()
}

func (self *ringModulo) Sub(other int) *ringModulo {
	other = self.normalizeValue(other)
	self.value -= other
	return self.normalizeSelf()
}

func (self *ringModulo) Mul(other int) *ringModulo {
	other = self.normalizeValue(other)
	self.value *= other
	return self.normalizeSelf()
}

func (self *ringModulo) Div(other int) *ringModulo {
	other = self.normalizeValue(other)
	inversed := InverseByModuloSlow(other, self.modulo)
	return self.Mul(inversed)
}

func (self *ringModulo) normalizeSelf() *ringModulo {
	self.value = self.normalizeValue(self.value)
	return self
}

func (self *ringModulo) normalizeValue(value int) int {
	value %= self.modulo
	if value < 0 {
		value += self.modulo
	}
	return value
}

func InverseByModuloSlow(a int, modulo int) int {
	result := 1
	for i := 0; i < modulo - 2; i++ {
		result = (result * a) % modulo
	}
	return result
}