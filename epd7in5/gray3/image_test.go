package gray3

import (
	"image"
	"testing"
)

func TestSetAt(t *testing.T) {
	img := NewImage(image.Rectangle{image.Point{0, 0}, image.Point{4, 1}})
	img.Set(0, 0, Gray)
	img.Set(1, 0, White)
	img.Set(2, 0, Black)
	img.Set(3, 0, Gray)
	if c := img.At(0, 0); c != Gray {
		t.Errorf("it should be gray. got %v", c)
	}
	if c := img.At(1, 0); c != White {
		t.Errorf("it should be white. got %v", c)
	}
	if c := img.At(2, 0); c != Black {
		t.Errorf("it should be black. got %v", c)
	}
	if c := img.At(3, 0); c != Gray {
		t.Errorf("it should be gray. got %v", c)
	}
}
