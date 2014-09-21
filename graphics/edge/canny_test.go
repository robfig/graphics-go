// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edge

import (
	"image"
	"testing"

	_ "image/png"
)

func TestCannyGopher(t *testing.T) {
	src, err := loadImage("../../testdata/gopher.png")
	if err != nil {
		t.Fatal(err)
	}

	dst := image.NewGray(src.Bounds())
	if err := Canny(dst, src); err != nil {
		t.Fatalf("%d: %v", err)
	}

	cmp, err := loadImage("../../testdata/gopher-canny.png")
	if err != nil {
		t.Fatal(err)
	}
	err = imageWithinTolerance(dst, cmp, 0x101)
	if err != nil {
		t.Fatal(err)
	}
}
