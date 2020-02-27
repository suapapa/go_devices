package epdcolor

import (
	"image"
	"testing"
)

func TestBWImage(t *testing.T) {
	img := NewBWImage(image.Rectangle{image.Point{0, 0}, image.Point{8, 1}})
	img.Set(0, 0, BWBlack)
	img.Set(1, 0, BWWhite)
	img.Set(2, 0, BWBlack)
	img.Set(3, 0, BWWhite)
	if c := img.At(0, 0); c != BWBlack {
		t.Errorf("it should be black. got %v", c)
	}
	if c := img.At(1, 0); c != BWWhite {
		t.Errorf("it should be white. got %v", c)
	}
	if c := img.At(2, 0); c != BWBlack {
		t.Errorf("it should be black. got %v", c)
	}
	if c := img.At(3, 0); c != BWWhite {
		t.Errorf("it should be White. got %v", c)
	}
}
