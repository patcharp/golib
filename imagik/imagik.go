package imagik

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/patcharp/golib/requests"
	"image"
	"net/http"
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

func (img *Imagik) Resize(w int, h int) {
	img.Img = imaging.Resize(img.Img, w, h, imaging.Lanczos)
}

func (img *Imagik) Crop(w int, h int) {
	img.Img = imaging.CropCenter(img.Img, w, h)
}

func (img *Imagik) ThumbnailAsByte(w int, h int) ([]byte, error) {
	src := imaging.Fill(img.Img, w, h, imaging.Center, imaging.Lanczos)
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, src, imaging.JPEG, imaging.JPEGQuality(100)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (img *Imagik) ThumbnailAsFile(filename string, w int, h int) error {
	src := imaging.Fill(img.Img, w, h, imaging.Center, imaging.Lanczos)
	if err := imaging.Save(src, filename, imaging.JPEGQuality(100)); err != nil {
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

func (img *Imagik) ExportAsByte() ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.JPEG, imaging.JPEGQuality(100)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (img *Imagik) ExportAsBytePNG() ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img.Img, imaging.PNG, imaging.PNGCompressionLevel(100)); err != nil {
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
