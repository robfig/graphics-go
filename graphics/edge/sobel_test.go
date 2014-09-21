// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edge

import (
	"fmt"
	"image"
	_ "image/png"
	"testing"
)

type opTest struct {
	name string
	fn   func(mag, dir *image.Gray, src image.Image) error
}

var opTests = []opTest{
	{"sobel", Sobel},
	{"scharr", Scharr},
	{"prewitt", Prewitt},
}

func TestEdgeOps(t *testing.T) {
	src, err := loadImage("../../testdata/gopher-100x150.png")
	if err != nil {
		t.Fatal(err)
	}

	for _, ot := range opTests {
		mag := image.NewGray(src.Bounds())
		dir := image.NewGray(src.Bounds())
		if err := ot.fn(mag, dir, src); err != nil {
			t.Fatalf("%s: %v", ot.name, err)
		}

		magFile := fmt.Sprintf("../../testdata/%s-mag.png", ot.name)
		cmp, err := loadImage(magFile)
		if err != nil {
			t.Fatalf("%s mag: %v", ot.name, err)
		}
		if err = imageWithinTolerance(mag, cmp, 0x101); err != nil {
			t.Fatalf("%s mag: %v", ot.name, err)
		}

		dirFile := fmt.Sprintf("../../testdata/%s-dir.png", ot.name)
		cmp, err = loadImage(dirFile)
		if err != nil {
			t.Fatalf("%s dir: %v", ot.name, err)
		}
		if err = imageWithinTolerance(dir, cmp, 0x101); err != nil {
			t.Fatalf("%s dir: %v", ot.name, err)
		}
	}
}
