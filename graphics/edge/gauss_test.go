// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edge

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"

	_ "image/png"
)

func sprintBox(box []byte, width, height int) string {
	buf := bytes.NewBuffer(nil)
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Fprintf(buf, " 0x%02x,", box[i])
			i++
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func delta(u0, u1 uint32) int {
	d := int(u0) - int(u1)
	if d < 0 {
		return -d
	}
	return d
}

func loadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

func withinTolerance(c0, c1 color.Color, tol int) bool {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r := delta(r0, r1)
	g := delta(g0, g1)
	b := delta(b0, b1)
	a := delta(a0, a1)
	return r <= tol && g <= tol && b <= tol && a <= tol
}

func imageWithinTolerance(m0, m1 image.Image, tol int) error {
	b0 := m0.Bounds()
	b1 := m1.Bounds()
	if !b0.Eq(b1) {
		return errors.New(fmt.Sprintf("got bounds %v want %v", b0, b1))
	}

	for y := b0.Min.Y; y < b0.Max.Y; y++ {
		for x := b0.Min.X; x < b0.Max.X; x++ {
			c0 := m0.At(x, y)
			c1 := m1.At(x, y)
			if !withinTolerance(c0, c1, tol) {
				e := fmt.Sprintf("got %v want %v at (%d, %d)", c0, c1, x, y)
				return errors.New(e)
			}
		}
	}
	return nil
}

func TestLoG(t *testing.T) {
	src, err := loadImage("../../testdata/gopher.png")
	if err != nil {
		t.Error(err)
		return
	}

	dst := image.NewGray(src.Bounds())
	LaplacianOfGaussian(dst, src)

	cmp, err := loadImage("../../testdata/gopher-log.png")
	if err != nil {
		t.Error(err)
		return
	}
	err = imageWithinTolerance(dst, cmp, 0x101)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDoG(t *testing.T) {
	src, err := loadImage("../../testdata/gopher.png")
	if err != nil {
		t.Error(err)
		return
	}

	dst := image.NewGray(src.Bounds())
	DifferenceOfGaussians(dst, src, 0.9, 1.2)

	cmp, err := loadImage("../../testdata/gopher-dog.png")
	if err != nil {
		t.Error(err)
		return
	}
	err = imageWithinTolerance(dst, cmp, 0x101)
	if err != nil {
		t.Error(err)
		return
	}
}
