package imagik

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/patcharp/golib/v2/requests"
	"image"
	"image/png"
	"net/http"
)

const (
	RotateSquare    = 0
	RotatePortrait  = 1
	RotateLandscape = 2
)

type Imagik struct {
	Img      image.Image
	MimeType string
	Size     int
}

func (img *Imagik) LoadFromFile(filename string) error {
	var err error
	img.Img, err = imaging.Open(filename, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}
	return nil
}

func (img *Imagik) LoadFromByte(b []byte) error {
	src, err := imaging.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}
	img.Img = src
	return nil
}

func (img *Imagik) LoadFromUrl(url string, headers map[string]string) error {
	r, err := requests.Get(url, headers, nil, 0)
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		return errors.New(r.Status)
	}
	return img.LoadFromByte(r.Body)
}

func (img *Imagik) Bounds() image.Rectangle {
	return img.Img.Bounds()
}

func (img *Imagik) IsSquare() bool {
	if img.Bounds().Dy() == img.Bounds().Dx() {
		return true
	}
	return false
}

func (img *Imagik) Rotation() int {
	if img.Bounds().Dy() == img.Bounds().Dx() {
		return RotateSquare
	} else if img.Bounds().Dy() > img.Bounds().Dx() {
		return RotatePortrait
	}
	return RotateLandscape
}

func (img *Imagik) Resize(w int, h int) {
	img.Img = imaging.Resize(img.Img, w, h, imaging.Lanczos)
}

func (img *Imagik) CropCenter(w int, h int) {
	img.Img = imaging.CropCenter(img.Img, w, h)
}

func (img *Imagik) Crop(rect image.Rectangle) {
	img.Img = imaging.Crop(img.Img, rect)
}

func (img *Imagik) SquareThumbnailAsPNGByte(size int) ([]byte, error) {
	if !img.IsSquare() {
		// Crop into square
		cropSize := img.Bounds().Dx()
		if img.Rotation() == RotateLandscape {
			cropSize = img.Bounds().Dy()
		}
		img.CropCenter(cropSize, cropSize)
	}
	// Check bounds size if cropped img size larger than size and size != 0
	if img.Bounds().Dx() > size && size != 0 {
		// Resize
		img.Resize(size, size)
	}
	// Return byte as png
	return img.ExportAsPNGByteWithCompLevel(70)
}

func (img *Imagik) ThumbnailAsByte(w int, h int) ([]byte, error) {
	src := imaging.Fill(img.Img, w, h, imaging.Center, imaging.Lanczos)
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, src, imaging.JPEG, imaging.JPEGQuality(70)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (img *Imagik) ThumbnailAsFile(filename string, w int, h int) error {
	src := imaging.Fill(img.Img, w, h, imaging.Center, imaging.Lanczos)
	if err := imaging.Save(src, filename, imaging.JPEGQuality(70)); err != nil {
		return err
	}
	return nil
}

func (img *Imagik) ExportAsDataUrl() (string, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.JPEG, imaging.JPEGQuality(100)); err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", http.DetectContentType(buf.Bytes()), base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

func (img *Imagik) ExportAsPNGDataUrl() (string, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.PNG, imaging.PNGCompressionLevel(100)); err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", http.DetectContentType(buf.Bytes()), base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

func (img *Imagik) ExportAsByte() ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.JPEG, imaging.JPEGQuality(100)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (img *Imagik) ExportAsPNGByte() ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.PNG, imaging.PNGCompressionLevel(100)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (img *Imagik) ExportAsPNGByteWithCompLevel(compLevel png.CompressionLevel) ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.PNG, imaging.PNGCompressionLevel(compLevel)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (img *Imagik) ExportAsFile(filename string) error {
	if err := imaging.Save(img.Img, filename, imaging.JPEGQuality(100)); err != nil {
		return err
	}
	return nil
}

func (img *Imagik) ExportAsPNGFile(filename string) error {
	if err := imaging.Save(img.Img, filename); err != nil {
		return err
	}
	return nil
}
