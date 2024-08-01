//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/float8
//

package float8

import (
	"bytes"
	"testing"

	"github.com/chewxy/math32"
	"github.com/kshard/float8/internal/math8"
)

func norm(x float32) float32 {
	// Note: It would be expected that ToFloat8(ToFloat32(x)) = x
	//       but due to noticeable error, it is not a case on small numbers
	//       small epsilon makes number to be approximate
	if x < 0 {
		return x - 1e-6
	}

	return x + 1e-6
}

func TestToFloat8(t *testing.T) {
	for expected, f32 := range f8tof32 {
		val := ToFloat8(norm(f32))
		if val != uint8(expected) {
			t.Errorf("0x%02x got=0x%02x f32=%f", expected, val, f32)
		}
	}
}

func TestToSlice8(t *testing.T) {
	f32s := make([]float32, len(f8tof32))
	expected := make([]Float8, len(f8tof32))

	for f8, f32 := range f8tof32 {
		expected = append(expected, Float8(f8))
		f32s = append(f32s, norm(f32))
	}

	f8s := ToSlice8(f32s)
	if !bytes.Equal(f8s, expected) {
		t.Errorf("got=%v expected=%v", f8s, expected)
	}
}

func TestToFloat32(t *testing.T) {
	for a := 0; a < 0x100; a++ {
		c := ToFloat32(uint8(a))
		e := math8.ToFloat32(uint8(a))
		if math32.Abs(c-e) > 1e-6 {
			t.Errorf("0x%02x wanted=%f, got=%f", a, e, c)
		}
	}
}

func TestAdd(t *testing.T) {
	for a := 0; a < 0x100; a++ {
		for b := 0; b < 0x100; b++ {
			c := Add(uint8(a), uint8(b))
			e := math8.Add(uint8(a), uint8(b))
			if c != e {
				t.Errorf("0x%02x + 0x%02x wanted=0x%02x, got=0x%02x", a, b, e, c)
			}
		}
	}
}

func TestSub(t *testing.T) {
	for a := 0; a < 0x100; a++ {
		for b := 0; b < 0x100; b++ {
			c := Sub(uint8(a), uint8(b))
			e := math8.Sub(uint8(a), uint8(b))
			if c != e {
				t.Errorf("0x%02x + 0x%02x wanted=0x%02x, got=0x%02x", a, b, e, c)
			}
		}
	}
}

func TestMul(t *testing.T) {
	for a := 0; a < 0x100; a++ {
		for b := 0; b < 0x100; b++ {
			c := Mul(uint8(a), uint8(b))
			e := math8.Mul(uint8(a), uint8(b))
			if c != e {
				t.Errorf("0x%02x + 0x%02x wanted=0x%02x, got=0x%02x", a, b, e, c)
			}
		}
	}
}

func TestDiv(t *testing.T) {
	for a := 0; a < 0x100; a++ {
		for b := 0; b < 0x100; b++ {
			c := Div(uint8(a), uint8(b))
			e := math8.Div(uint8(a), uint8(b))
			if c != e {
				t.Errorf("0x%02x + 0x%02x wanted=0x%02x, got=0x%02x", a, b, e, c)
			}
		}
	}
}

var (
	f8   uint8
	f32  float32
	f32s = f8tof32[:]
	f8s  []uint8
)

func BenchmarkToFloat8(b *testing.B) {
	for i := b.N; i > 0; i-- {
		f8 = ToFloat8(f8tof32[i%0x100])
	}
}

func BenchmarkToFloat32(b *testing.B) {
	for i := b.N; i > 0; i-- {
		f32 = ToFloat32(uint8(i % 0x100))
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := b.N; i > 0; i-- {
		v := uint8(i % 0x100)
		f8 = Add(v, v)
	}
}

func BenchmarkMul(b *testing.B) {
	for i := b.N; i > 0; i-- {
		v := uint8(i % 0x100)
		f8 = Mul(v, v)
	}
}

func BenchmarkToSlice8(b *testing.B) {
	for i := b.N; i > 0; i-- {
		f8s = ToSlice8(f32s)
	}
}
