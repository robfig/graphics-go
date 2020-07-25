// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edge

import (
	"image"
	"image/draw"

	"github.com/robfig/graphics-go/graphics"
	"github.com/robfig/graphics-go/graphics/convolve"
)

// LaplacianOfGaussian approximates a 2D laplacian of gaussian with a convolution kernel.
func LaplacianOfGaussian(dst *image.Gray, src image.Image) {
	srcg, ok := src.(*image.Gray)
	if !ok {
		b := src.Bounds()
		srcg = image.NewGray(b)
		draw.Draw(srcg, b, src, b.Min, draw.Src)
	}

	k, err := convolve.NewKernel([]float64{
		0, 0, 1, 0, 0,
		0, 1, 2, 1, 0,
		1, 2, -16, 2, 1,
		0, 1, 2, 1, 0,
		0, 0, 1, 0, 0,
	})
	if err != nil {
		panic(err) // impossible
	}

	convolve.Convolve(dst, srcg, k)
}

// DifferenceOfGaussians produces the difference of Gaussians sd0 and sd1.
func DifferenceOfGaussians(dst *image.Gray, src image.Image, sd0, sd1 float64) {
	b := src.Bounds()
	srcg := image.NewGray(b)
	m0 := image.NewGray(b)
	m1 := image.NewGray(b)

	draw.Draw(srcg, b, src, b.Min, draw.Src)

	graphics.Blur(m0, srcg, &graphics.BlurOptions{StdDev: sd0})
	graphics.Blur(m1, srcg, &graphics.BlurOptions{StdDev: sd1})

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			off := (y-m0.Rect.Min.Y)*m0.Stride + (x - m0.Rect.Min.X)
			c := m0.Pix[off] - m1.Pix[off]
			if c < 0 {
				c = -c
			}

			doff := (y-dst.Rect.Min.Y)*dst.Stride + (x - dst.Rect.Min.X)
			dst.Pix[doff] = c
		}
	}
}
