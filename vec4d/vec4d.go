package vec4d

import (
	"fmt"
	"math"

	"github.com/ungerik/go3d/genericd"
	"github.com/ungerik/go3d/vec3d"
)

var (
	Zero = T{}

	UnitXW = T{1, 0, 0, 1}
	UnitYW = T{0, 1, 0, 1}
	UnitZW = T{0, 0, 1, 1}
	UnitW  = T{0, 0, 0, 1}

	Red   = T{1, 0, 0, 1}
	Green = T{0, 1, 0, 1}
	Blue  = T{0, 0, 1, 1}
	Black = T{0, 0, 0, 1}

	MinVal = T{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64, 1}
	MaxVal = T{+math.MaxFloat64, +math.MaxFloat64, +math.MaxFloat64, 1}
)

type T [4]float64

// From copies a T from a generic.T implementation.
func From(other genericd.T) T {
	switch other.Size() {
	case 2:
		return T{other.Get(0, 0), other.Get(0, 1), 0, 1}
	case 3:
		return T{other.Get(0, 0), other.Get(0, 1), other.Get(0, 2), 1}
	case 4:
		return T{other.Get(0, 0), other.Get(0, 1), other.Get(0, 2), other.Get(0, 3)}
	default:
		panic("Unsupported type")
	}
}

func FromVec3(other *vec3d.T) T {
	return T{other[0], other[1], other[2], 1}
}

// Parse parses T from a string. See also String()
func Parse(s string) (r T, err error) {
	_, err = fmt.Sscanf(s, "%f %f %f %f", &r[0], &r[1], &r[2], &r[3])
	return r, err
}

// String formats T as string. See also Parse().
func (self *T) String() string {
	return fmt.Sprintf("%f %f %f %f", self[0], self[1], self[2], self[3])
}

// Rows returns the number of rows of the vector.
func (self *T) Rows() int {
	return 4
}

// Cols returns the number of columns of the vector.
func (self *T) Cols() int {
	return 1
}

// Size returns the number elements of the vector.
func (self *T) Size() int {
	return 4
}

// Slice returns the elements of the vector as slice.
func (self *T) Slice() []float64 {
	return []float64{self[0], self[1], self[2], self[3]}
}

// Get returns one element of the vector.
func (self *T) Get(col, row int) float64 {
	return self[row]
}

// IsZero checks if all elements of the vector are zero.
func (self *T) IsZero() bool {
	return self[0] == 0 && self[1] == 0 && self[2] == 0 && self[3] == 0
}

// Length returns the length of the vector.
// See also LengthSqr and Normalize.
func (self *T) Length() float64 {
	v3 := self.Vec3DividedByW()
	return v3.Length()
}

// Length returns the squared length of the vector.
// See also Length and Normalize.
func (self *T) LengthSqr() float64 {
	v3 := self.Vec3DividedByW()
	return v3.LengthSqr()
}

// Scale multiplies the first 3 element of the vector by f and returns self.
func (self *T) Scale(f float64) *T {
	self[0] *= f
	self[1] *= f
	self[2] *= f
	return self
}

// Scaled returns a copy of self with the first 3 elements multiplies by f.
func (self *T) Scaled(f float64) T {
	return T{self[0] * f, self[1] * f, self[2] * f, self[3]}
}

func (self *T) Invert() {
	self[0] = -self[0]
	self[1] = -self[1]
	self[2] = -self[2]
}

func (self *T) Inverted() T {
	return T{-self[0], -self[1], -self[2], self[3]}
}

func (self *T) Normalize() *T {
	v3 := self.Vec3DividedByW()
	v3.Normalize()
	self[0] = v3[0]
	self[1] = v3[1]
	self[2] = v3[2]
	self[3] = 1
	return self
}

func (self *T) Normalized() T {
	v := *self
	v.Normalize()
	return v
}

func (self *T) Normal() T {
	v3 := self.Vec3()
	n3 := v3.Normal()
	return T{n3[0], n3[1], n3[2], 1}
}

func (self *T) DivideByW() *T {
	oow := 1 / self[3]
	self[0] *= oow
	self[1] *= oow
	self[2] *= oow
	self[3] = 1
	return self
}

func (self *T) DividedByW() T {
	oow := 1 / self[3]
	return T{self[0] * oow, self[1] * oow, self[2] * oow, 1}
}

func (self *T) Vec3DividedByW() vec3d.T {
	oow := 1 / self[3]
	return vec3d.T{self[0] * oow, self[1] * oow, self[2] * oow}
}

func (self *T) Vec3() vec3d.T {
	return vec3d.T{self[0], self[1], self[2]}
}

func (self *T) AssignVec3(v *vec3d.T) *T {
	self[0] = v[0]
	self[1] = v[1]
	self[2] = v[2]
	self[3] = 1
	return self
}

func (self *T) Add(v *T) *T {
	if v[3] == self[3] {
		self[0] += v[0]
		self[1] += v[1]
		self[2] += v[2]
	} else {
		self.DividedByW()
		v3 := v.Vec3DividedByW()
		self[0] += v3[0]
		self[1] += v3[1]
		self[2] += v3[2]
	}
	return self
}

func (self *T) Sub(v *T) *T {
	if v[3] == self[3] {
		self[0] -= v[0]
		self[1] -= v[1]
		self[2] -= v[2]
	} else {
		self.DividedByW()
		v3 := v.Vec3DividedByW()
		self[0] -= v3[0]
		self[1] -= v3[1]
		self[2] -= v3[2]
	}
	return self
}

func Add(a, b *T) T {
	if a[3] == b[3] {
		return T{a[0] + b[0], a[1] + b[1], a[2] + b[2], 1}
	} else {
		a3 := a.Vec3DividedByW()
		b3 := b.Vec3DividedByW()
		return T{a3[0] + b3[0], a3[1] + b3[1], a3[2] + b3[2], 1}
	}
}

func Sub(a, b *T) T {
	if a[3] == b[3] {
		return T{a[0] - b[0], a[1] - b[1], a[2] - b[2], 1}
	} else {
		a3 := a.Vec3DividedByW()
		b3 := b.Vec3DividedByW()
		return T{a3[0] - b3[0], a3[1] - b3[1], a3[2] - b3[2], 1}
	}
}

func Dot(a, b *T) float64 {
	a3 := a.Vec3DividedByW()
	b3 := b.Vec3DividedByW()
	return vec3d.Dot(&a3, &b3)
}

func Dot4(a, b *T) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
}

func Cross(a, b *T) T {
	a3 := a.Vec3DividedByW()
	b3 := b.Vec3DividedByW()
	c3 := vec3d.Cross(&a3, &b3)
	return T{c3[0], c3[1], c3[2], 1}
}

func Angle(a, b *T) float64 {
	return math.Acos(Dot(a, b))
}
