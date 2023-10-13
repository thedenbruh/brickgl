package fauxgl

import (
	"bytes"
	"image"
	"math"
)

// Texture interface for texture
type Texture interface {
	Sample(u, v float64) Color
	Texture() image.Image
}

// LoadTexture returns texture from filepath
func LoadTexture(path string) (Texture, error) {
	im, err := LoadImage(path)
	if err != nil {
		return nil, err
	}
	return NewImageTexture(im), nil
}

// TexFromBytes returns fauxgl texture created with given bytes
func TexFromBytes(data []byte) (tex Texture) {
	img, _, _ := image.Decode(bytes.NewReader(data))

	tex = NewImageTexture(img)

	return
}

// ImageTexture struct to hold image
type ImageTexture struct {
	Width  int
	Height int
	Image  image.Image
}

// NewImageTexture image.Image to texture
func NewImageTexture(im image.Image) Texture {
	size := im.Bounds().Max
	return &ImageTexture{size.X, size.Y, im}
}

// Sample get the color of a texture at coordinates
func (t *ImageTexture) Sample(u, v float64) Color {
	v = 1 - v
	u -= math.Floor(u)
	v -= math.Floor(v)
	x := int(u * float64(t.Width))
	y := int(v * float64(t.Height))
	return MakeColor(t.Image.At(x, y))
}

// Texture texture to image.Image
func (t *ImageTexture) Texture() image.Image {
	return t.Image
}
