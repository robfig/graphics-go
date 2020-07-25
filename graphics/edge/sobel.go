// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edge

import (
	"errors"
	"image"
	"image/draw"
	"math"

	"github.com/robfig/graphics-go/graphics/convolve"
)

var (
	sobelX = &convolve.SeparableKernel{
		X: []float64{-1, 0, +1},
		Y: []float64{1, 2, 1},
	}
	sobelY = &convolve.SeparableKernel{
		X: []float64{1, 2, 1},
		Y: []float64{-1, 0, +1},
	}
	scharrX = &convolve.SeparableKernel{
		X: []float64{-1, 0, +1},
		Y: []float64{3, 10, 3},
	}
	scharrY = &convolve.SeparableKernel{
		X: []float64{3, 10, 3},
		Y: []float64{-1, 0, +1},
	}
	prewittX = &convolve.SeparableKernel{
		X: []float64{-1, 0, +1},
		Y: []float64{1, 1, 1},
	}
	prewittY = &convolve.SeparableKernel{
		X: []float64{1, 1, 1},
		Y: []float64{-1, 0, +1},
	}
)

func diffOp(mag, dir *image.Gray, src image.Image, opX, opY *convolve.SeparableKernel) error {
	if src == nil {
		return errors.New("graphics: src is nil")
	}
	b := src.Bounds()

	srcg, ok := src.(*image.Gray)
	if !ok {
		srcg = image.NewGray(b)
		draw.Draw(srcg, b, src, b.Min, draw.Src)
	}

	mx := image.NewGray(b)
	if err := convolve.Convolve(mx, srcg, opX); err != nil {
		return err
	}

	my := image.NewGray(b)
	if err := convolve.Convolve(my, srcg, opY); err != nil {
		return err
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			off := (y-mx.Rect.Min.Y)*mx.Stride + (x-mx.Rect.Min.X)*1
			cx := float64(mx.Pix[off])
			cy := float64(my.Pix[off])

			if mag != nil {
				off = (y-mag.Rect.Min.Y)*mag.Stride + (x-mag.Rect.Min.X)*1
				mag.Pix[off] = uint8(math.Sqrt(cx*cx + cy*cy))
			}
			if dir != nil {
				off = (y-dir.Rect.Min.Y)*dir.Stride + (x-dir.Rect.Min.X)*1
				angle := math.Atan(cy / cx)
				// Round the angle to 0, 45, 90, or 135 degrees.
				angle = math.Mod(angle, 2*math.Pi)
				var degree uint8
				if angle <= math.Pi/8 {
					degree = 0
				} else if angle <= math.Pi*3/8 {
					degree = 45
				} else if angle <= math.Pi*5/8 {
					degree = 90
				} else if angle <= math.Pi*7/8 {
					degree = 135
				} else {
					degree = 0
				}
				dir.Pix[off] = degree
			}
		}
	}
	return nil
}

// Sobel returns the magnitude and direction of the Sobel operator.
// dir pixels hold the rounded direction value either 0, 45, 90, or 135.
func Sobel(mag, dir *image.Gray, src image.Image) error {
	return diffOp(mag, dir, src, sobelX, sobelY)
}

// Scharr returns the magnitude and direction of the Scharr operator.
// This is very similar to Sobel, with less angular error.
func Scharr(mag, dir *image.Gray, src image.Image) error {
	return diffOp(mag, dir, src, scharrX, scharrY)
}

// Prewitt returns the magnitude and direction of the Prewitt operator.
func Prewitt(mag, dir *image.Gray, src image.Image) error {
	return diffOp(mag, dir, src, prewittX, prewittY)
}
