package epdcolor

import (
	"image"
	"testing"
)

func TestWBImage(t *testing.T) {
	img := NewWBImage(image.Rectangle{image.Point{0, 0}, image.Point{8, 1}})
	img.Set(0, 0, WBBlack)
	img.Set(1, 0, WBWhite)
	img.Set(2, 0, WBBlack)
	img.Set(3, 0, WBWhite)
	if c := img.At(0, 0); c != WBBlack {
		t.Errorf("it should be black. got %v", c)
	}
	if c := img.At(1, 0); c != WBWhite {
		t.Errorf("it should be white. got %v", c)
	}
	if c := img.At(2, 0); c != WBBlack {
		t.Errorf("it should be black. got %v", c)
	}
	if c := img.At(3, 0); c != WBWhite {
		t.Errorf("it should be White. got %v", c)
	}
}
