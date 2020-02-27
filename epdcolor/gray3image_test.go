package epdcolor

import (
	"image"
	"testing"
)

func TestGray3Image(t *testing.T) {
	img := NewGray3Image(image.Rectangle{image.Point{0, 0}, image.Point{4, 1}})
	img.Set(0, 0, Gray3Gray)
	img.Set(1, 0, Gray3White)
	img.Set(2, 0, Gray3Black)
	img.Set(3, 0, Gray3Gray)
	if c := img.At(0, 0); c != Gray3Gray {
		t.Errorf("it should be gray. got %v", c)
	}
	if c := img.At(1, 0); c != Gray3White {
		t.Errorf("it should be white. got %v", c)
	}
	if c := img.At(2, 0); c != Gray3Black {
		t.Errorf("it should be black. got %v", c)
	}
	if c := img.At(3, 0); c != Gray3Gray {
		t.Errorf("it should be gray. got %v", c)
	}
}
