package quaterniond

import (
	"fmt"
	"math"

	"github.com/ungerik/go3d/vec3d"
	"github.com/ungerik/go3d/vec4d"
)

var (
	Zero  = T{}
	Ident = T{0, 0, 0, 1}
)

type T [4]float64

func FromAxisAngle(axis *vec3d.T, angle float64) T {
	angle *= 0.5
	sin := math.Sin(angle)
	q := T{axis[0] * sin, axis[1] * sin, axis[2] * sin, math.Cos(angle)}
	return q.Normalized()
}

func FromXAxisAngle(angle float64) T {
	angle *= 0.5
	return T{math.Sin(angle), 0, 0, math.Cos(angle)}
}

func FromYAxisAngle(angle float64) T {
	angle *= 0.5
	return T{0, math.Sin(angle), 0, math.Cos(angle)}
}

func FromZAxisAngle(angle float64) T {
	angle *= 0.5
	return T{0, 0, math.Sin(angle), math.Cos(angle)}
}

func FromEulerAngles(yHead, xPitch, zRoll float64) T {
	qy := FromYAxisAngle(yHead)
	qx := FromXAxisAngle(xPitch)
	qz := FromZAxisAngle(zRoll)
	return Mul3(&qy, &qx, &qz)
}

func FromVec4(v *vec4d.T) T {
	return T(*v)
}

func (self *T) Vec4() vec4d.T {
	return vec4d.T(*self)
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

func (self *T) AxisAngle() (axis vec3d.T, angle float64) {
	cos := self[3]
	sin := math.Sqrt(1 - cos*cos)
	angle = math.Acos(cos)

	var ooSin float64
	if math.Abs(sin) < 0.0005 {
		ooSin = 1
	} else {
		ooSin = 1 / sin
	}
	axis[0] = self[0] * ooSin
	axis[1] = self[1] * ooSin
	axis[2] = self[2] * ooSin

	return axis, angle
}

func (self *T) Norm() float64 {
	return self[0]*self[0] + self[1]*self[1] + self[2]*self[2] + self[3]*self[3]
}

func (self *T) Normalize() {
	norm := self.Norm()
	if norm != 1 && norm != 0 {
		ool := 1 / math.Sqrt(norm)
		self[0] *= ool
		self[1] *= ool
		self[2] *= ool
		self[3] *= ool
	}
}

func (self *T) Normalized() T {
	norm := self.Norm()
	if norm != 1 && norm != 0 {
		ool := 1 / math.Sqrt(norm)
		return T{
			self[0] * ool,
			self[1] * ool,
			self[2] * ool,
			self[3] * ool,
		}
	} else {
		return *self
	}
}

func (self *T) Negate() {
	self[0] = -self[0]
	self[1] = -self[1]
	self[2] = -self[2]
	self[3] = -self[3]
}

func (self *T) Negated() T {
	return T{-self[0], -self[1], -self[2], -self[3]}
}

func (self *T) Invert() {
	self[0] = -self[0]
	self[1] = -self[1]
	self[2] = -self[2]
}

func (self *T) Inverted() T {
	return T{-self[0], -self[1], -self[2], self[3]}
}

func (self *T) SetShortestRotation(other *T) {
	if !IsShortestRotation(self, other) {
		self.Negate()
	}
}

func IsShortestRotation(a, b *T) bool {
	return Dot(a, b) >= 0
}

func (self *T) IsUnitQuat(tolerance float64) bool {
	norm := self.Norm()
	return norm >= (1.0-tolerance) && norm <= (1.0+tolerance)
}

func (self *T) RotateVec3(v *vec3d.T) {
	qv := T{v[0], v[1], v[2], 0}
	inv := self.Inverted()
	q := Mul3(self, &qv, &inv)
	v[0] = q[0]
	v[1] = q[1]
	v[2] = q[2]
}

func (self *T) RotatedVec3(v *vec3d.T) vec3d.T {
	qv := T{v[0], v[1], v[2], 0}
	inv := self.Inverted()
	q := Mul3(self, &qv, &inv)
	return vec3d.T{q[0], q[1], q[2]}
}

func Dot(a, b *T) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
}

func Mul(a, b *T) T {
	q := T{
		a[3]*b[0] + a[0]*b[3] + a[1]*b[2] - a[2]*b[1],
		a[3]*b[1] + a[1]*b[3] + a[2]*b[0] - a[0]*b[2],
		a[3]*b[2] + a[2]*b[3] + a[0]*b[1] - a[1]*b[0],
		a[3]*b[3] - a[0]*b[0] - a[1]*b[1] - a[2]*b[2],
	}
	return q.Normalized()
}

func Mul3(a, b, c *T) T {
	q := Mul(a, b)
	return Mul(&q, c)
}

func Mul4(a, b, c, d *T) T {
	q := Mul(a, b)
	q = Mul(&q, c)
	return Mul(&q, d)
}

func Slerp(a, b *T, f float64) T {
	d := math.Acos(a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3])
	ooSinD := 1 / math.Sin(d)

	f1 := math.Sin(d*(1-f)) * ooSinD
	f2 := math.Sin(d*f) * ooSinD

	q := T{
		a[0]*f1 + b[0]*f2,
		a[1]*f1 + b[1]*f2,
		a[2]*f1 + b[2]*f2,
		a[3]*f1 + b[3]*f2,
	}

	return q.Normalized()
}

func Vec3Diff(a, b *vec3d.T) T {
	cr := vec3d.Cross(a, b)
	sr := math.Sqrt(2 * (1 + vec3d.Dot(a, b)))
	oosr := 1 / sr

	q := T{cr[0] * oosr, cr[1] * oosr, cr[2] * oosr, sr * 0.5}
	return q.Normalized()
}
