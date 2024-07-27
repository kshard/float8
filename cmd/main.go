//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/float8
//

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kshard/float8/internal/math8"
)

func main() {
	fmt.Printf("==> code book for float32\n")
	if err := f8tof32(); err != nil {
		panic(err)
	}

	for name, f := range map[string]func(uint8, uint8) uint8{
		"add": math8.Add,
		"sub": math8.Sub,
		"mul": math8.Mul,
		"div": math8.Div,
	} {
		fmt.Printf("==> code book for %s\n", name)
		if err := codebook(name, f); err != nil {
			panic(err)
		}
	}
}

func f8tof32() error {
	fd, err := os.Create("../float32.go")
	if err != nil {
		return err
	}
	defer fd.Close()

	seq := make([]string, 0x100)
	for f8 := 0; f8 < 0x100; f8++ {
		seq[f8] = fmt.Sprintf("%f", math8.ToFloat32(uint8(f8)))
	}

	tpl := `// DO NOT EDIT! Use cmd to regenerate it.
package float8

//
// The code book for translating float8 to float32
//

var f8tof32 = [0x100]float32{%s}
	`

	_, err = fd.WriteString(fmt.Sprintf(tpl, strings.Join(seq, ",")))
	if err != nil {
		return err
	}

	return nil
}

func codebook(name string, f func(uint8, uint8) uint8) error {
	fd, err := os.Create(fmt.Sprintf("../%s.go", name))
	if err != nil {
		return err
	}
	defer fd.Close()

	seq := make([]string, 0x100*0x100)
	for a := 0; a < 0x100; a++ {
		for b := 0; b < 0x100; b++ {
			seq[a<<8|b] = fmt.Sprintf("0x%x", f(uint8(a), uint8(b)))
		}
	}

	tpl := `// DO NOT EDIT! Use cmd to regenerate it.
package float8

//
// The code book for translating float8 to float32
//

var %s = [0x10000]uint8{%s}
	`

	_, err = fd.WriteString(fmt.Sprintf(tpl, name, strings.Join(seq, ",")))
	if err != nil {
		return err
	}

	return nil
}
