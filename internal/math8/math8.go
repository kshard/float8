//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/float8
//

// Package math8 implements canonical operations using float8 type.
// It implements functionally correct library but slow ops.
package math8

import (
	"math"

	"github.com/chewxy/math32"
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
	mantissaBias = 8.0

	// exponent base
	base = 2

	//
	positiveInf = 0x7f ^ mantissaMask
	negativeInf = 0xff ^ mantissaMask
)

type Float8 = uint8

// Return Float8 value from float32
func ToFloat8(f32 float32) Float8 {
	if f32 == 0 {
		return 0
	}

	// Extract sign, exponent, and mantissa from float32
	sign := uint8(0)
	if f32 < 0.0 {
		sign = 1
		f32 = -f32
	}

	// Handle special cases: infinity and NaN
	if math32.IsInf(f32, 1) {
		return positiveInf
	}
	if math32.IsInf(f32, -1) {
		return negativeInf
	}

	expValue := math32.Floor(math32.Log2(f32))
	if expValue > exponentHi {
		return positiveInf
	}
	if expValue < exponentLo {
		return 0
	}

	exponent := uint8(expValue + exponentBias)
	if exponent > exponentHi {
		exponent = exponentHi
	}

	mantissa := uint8((f32/math32.Pow(base, expValue) - 1.0) * mantissaBias)
	if mantissa > mantissaMask {
		mantissa = mantissaMask
	}

	return (sign << 7) | (exponent << mantissaLen) | (mantissa & mantissaMask)
}

// Return float32 value from Float8
func ToFloat32(f8 Float8) float32 {
	if f8 == 0 {
		return 0.0
	}

	sign := (f8 & signMask) >> 7
	exponent := (f8 & exponentMask) >> mantissaLen
	mantissa := f8 & mantissaMask

	// Calculate the actual exponent value
	exponentValue := int(exponent) - exponentBias

	// Calculate the actual mantissa value
	mantissaValue := 1.0 + float32(mantissa)/mantissaBias

	// Calculate the float32 value
	val := mantissaValue * float32(math.Pow(base, float64(exponentValue)))

	// Apply sign
	if sign == 1 {
		val = -val
	}

	return val
}

// Add two Float8
func Add(a, b Float8) Float8 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	aSign := (a & signMask) >> 7
	bSign := (b & signMask) >> 7

	aExponent := (a & exponentMask) >> mantissaLen
	bExponent := (b & exponentMask) >> mantissaLen

	aMantissa := 1.0 + float32(a&mantissaMask)/mantissaBias
	bMantissa := 1.0 + float32(b&mantissaMask)/mantissaBias

	// Align exponents
	if aExponent > bExponent {
		bMantissa /= float32(math.Pow(base, float64(aExponent-bExponent)))
		bExponent = aExponent
	} else if aExponent < bExponent {
		aMantissa /= float32(math.Pow(base, float64(bExponent-aExponent)))
		aExponent = bExponent
	}

	// Perform addition/subtraction
	var mantissa float32
	var sign uint8
	if aSign == bSign {
		mantissa = aMantissa + bMantissa
		sign = aSign
	} else {
		if aMantissa > bMantissa {
			mantissa = aMantissa - bMantissa
			sign = aSign
		} else {
			mantissa = bMantissa - aMantissa
			sign = bSign
		}
	}

	// Normalize result
	exponent := int(aExponent)
	if mantissa >= 2.0 {
		mantissa /= 2.0
		exponent++
	}
	for mantissa < 1.0 && mantissa != 0 {
		mantissa *= 2.0
		exponent--
	}

	if exponent > exponentHi {
		if sign == 0 {
			return positiveInf
		} else {
			return negativeInf
		}
	}
	if exponent < 0 {
		return 0
	}

	// Reconstruct the minifloat
	result := uint8(sign << 7)
	result |= uint8(exponent << mantissaLen)
	result |= uint8((mantissa-1.0)*mantissaBias) & mantissaMask

	return result
}

// Subtract two Float8
func Sub(a, b Float8) Float8 {
	if a == b {
		return 0
	}

	return Add(a, b^signMask)
}

// Multiply Float8
func Mul(a, b Float8) Float8 {
	if a == 0 || b == 0 {
		return 0
	}

	aSign := (a & signMask) >> 7
	bSign := (b & signMask) >> 7
	sign := aSign ^ bSign

	aExponent := (a & exponentMask) >> mantissaLen
	bExponent := (b & exponentMask) >> mantissaLen
	exponent := int(aExponent) + int(bExponent) - exponentBias

	aMantissa := 1.0 + float32(a&mantissaMask)/mantissaBias
	bMantissa := 1.0 + float32(b&mantissaMask)/mantissaBias
	mantissa := aMantissa * bMantissa

	if mantissa >= 2.0 {
		mantissa /= 2.0
		exponent++
	}

	if exponent > exponentHi {
		if sign == 0 {
			return positiveInf
		} else {
			return negativeInf
		}
	}

	if exponent < 0 {
		return 0
	}

	val := uint8(sign << 7)
	val |= uint8(exponent << mantissaLen)
	val |= uint8((mantissa-1.0)*mantissaBias) & mantissaMask

	return val
}

// Divide float8
func Div(a, b Float8) Float8 {
	if a == 0 {
		return 0
	}

	// Extract components
	aSign := (a & signMask) >> 7
	bSign := (b & signMask) >> 7
	sign := aSign ^ bSign

	if b == 0 {
		if aSign == 0 {
			return positiveInf
		} else {
			return negativeInf
		}
	}

	aExponent := (a & exponentMask) >> mantissaLen
	bExponent := (b & exponentMask) >> mantissaLen
	exponent := int(aExponent) - int(bExponent) + exponentBias

	aMantissa := 1.0 + float32(a&mantissaMask)/mantissaBias
	bMantissa := 1.0 + float32(b&mantissaMask)/mantissaBias
	mantissa := aMantissa / bMantissa

	// Normalize result mantissa
	if mantissa >= 2.0 {
		mantissa /= 2.0
		exponent++
	} else if mantissa < 1.0 && mantissa != 0 {
		mantissa *= 2.0
		exponent--
	}

	if exponent > exponentHi {
		if sign == 0 {
			return positiveInf
		} else {
			return negativeInf
		}
	}
	if exponent < 0 {
		return 0
	}

	// Convert result mantissa to 3-bit format
	mantissaBits := uint8((mantissa - 1.0) * mantissaBias)
	if mantissaBits > mantissaMask {
		mantissaBits = mantissaMask
	}

	// Construct the result minifloat
	result := uint8(sign << 7)
	result |= uint8(exponent << 3)
	result |= mantissaBits & mantissaMask

	return result
}
