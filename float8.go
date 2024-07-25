//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/float8
//

// Package float8 implement minifloat (https://en.wikipedia.org/wiki/Minifloat)
// compatible with IEEE 754 and FP8 E4M3
// The number is defined as ±mantissa × 2^exponent
package float8

import (
	"math"
)

const (
	signMask     = 0b10000000 // 0x80
	exponentMask = 0b01111000 // 0x78
	mantissaMask = 0b00000111 // 0x07
	mantissaLen  = 3

	// See https://en.wikipedia.org/wiki/Exponent_bias
	//
	// bias = 2^(|exponent|-1) - 1
	// high = 2^|exponent| - 1
	exponentBias = 7
	exponentHi   = 15
	exponentLo   = -7

	// In a floating-point number representation, the mantissa (or significand)
	// represents the precision bits of the number. For an 8-bit minifloat with
	// 3 bits for the mantissa, these bits represent fractional values that
	// need to be converted to a floating-point format. These bits need to be
	// scaled to represent a fractional value between [1, 2). The bias normalize value.
	//
	// 2^|mantissa|
	// mantissaBias = 8.0

	// exponent base
	// base = 2

	//
	float32Bias = 127
)

const (
	Infinity = 0x7f | mantissaMask
)

// Float8 data type
type Float8 = uint8

// Convert float32 to float8
func ToFloat8(f32 float32) Float8 {
	if f32 == 0.0 {
		return 0x00
	}

	bits := math.Float32bits(f32)
	sign := uint8((bits >> 31) & 0x01)   // Extract sign (1 bit)
	exponent := int((bits >> 23) & 0xFF) // Extract exponent (8 bits)

	// Extract mantissa (23 bits) and add the implicit leading 1
	mantissa := int(bits & 0x7FFFFF)
	if exponent != 0 {
		mantissa |= 0x800000
	}

	// Adjust exponent from float32 bias (127) to minifloat bias (7)
	exponent = exponent - float32Bias + exponentBias

	// Handle overflow and underflow
	if exponent > exponentHi {
		return Infinity
	}
	if exponent < 0 {
		return 0x00
	}

	// Normalize mantissa to fit into 3 bits
	shift := 20 // Shift to convert 23-bit mantissa to 3-bit
	mantissa = (mantissa >> shift) & mantissaMask

	return (sign << 7) | (uint8(exponent) << 3) | uint8(mantissa)
}

// Convert float8 to float32
func ToFloat32(f8 Float8) float32 { return f8tof32[f8] }

// Add float8(s)
func Add(a, b Float8) Float8 { return add[int(a)<<8|int(b)] }

// Subtract float8(s)
func Sub(a, b Float8) Float8 { return sub[int(a)<<8|int(b)] }

// Multiply float8(s)
func Mul(a, b Float8) Float8 { return mul[int(a)<<8|int(b)] }

// Divide float8(s)
func Div(a, b Float8) Float8 { return div[int(a)<<8|int(b)] }
