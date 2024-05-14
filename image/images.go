package image

import (
	"bytes"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

var (
	Webp = "webp"
	Jpeg = "jpeg"
	Png  = "png"
)

type OptFunc func(image *Image) error

type Image struct {
	filename string
	ext      string
	file     io.Reader
}

func WithWidthAndHeight(w, h uint) OptFunc {
	return func(image *Image) error {
		if w == 0 && h == 0 {
			return nil
		}

		err := image.ResizeImage(w, h)
		if err != nil {
			return err
		}
		return nil
	}
}

func Format(format string) OptFunc {
	return func(image *Image) error {
		switch format {
		case Webp:
			err := image.ConvertToWebp()
			return err
		case Jpeg:
			//TODO soon to be don
		case Png:
			//TODO soon to be done
		}
		return nil
	}
}

func NewImage(filename string, file io.Reader, opts ...OptFunc) (*Image, error) {
	filename = path.Base(filename)

	image := &Image{
		filename: filename[:len(filename)-len(filepath.Ext(filename))],
		file:     file,
		ext:      path.Ext(filename),
	}

	for _, fn := range opts {

		err := fn(image)
		if err != nil {
			return nil, err
		}
	}
	return image, nil
}

func (i *Image) Save(dir ...string) error {
	dirpath := ""
	if len(dir) > 0 {
		dirpath = dir[0]
	}
	if dirpath != "" {
		err := os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	log.Println(path.Join(dirpath, i.filename+i.ext))
	file, err := os.OpenFile(path.Join(dirpath, i.filename+i.ext), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, i.file)
	return err

}

func (i *Image) ResizeImage(w uint, h uint) error {

	if i.ext == ".jpg" || i.ext == ".jpeg" {

		image, err := jpeg.Decode(i.file)
		if err != nil {
			return err
		}

		var newImage = resize.Resize(w, h, image, resize.Lanczos3)

		newFile := bytes.NewBuffer(make([]byte, 0))

		err = jpeg.Encode(newFile, newImage, nil)
		if err != nil {
			return err
		}

		i.file = newFile

	} else if i.ext == ".png" {

		image, err := png.Decode(i.file)
		if err != nil {
			return err
		}

		newImage := resize.Resize(w, h, image, resize.Lanczos3)

		newFile := bytes.NewBuffer(make([]byte, 0))

		err = png.Encode(newFile, newImage)
		if err != nil {
			return err
		}

		i.file = newFile

	}

	return nil
}

func (i *Image) ConvertToWebp() error {

	ext := filepath.Ext(i.filename)
	if ext == ".jpg" || ext == ".jpeg" {
		image, err := jpeg.Decode(i.file)
		if err != nil {
			return err
		}

		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		if err != nil {
			log.Fatalln(err)
		}

		newFile := bytes.NewBuffer(make([]byte, 0))

		if err := webp.Encode(newFile, image, options); err != nil {
			log.Fatalln(err)
		}
		i.file = newFile
		i.ext = ".webp"

		return nil
	} else if ext == ".png" {
		image, err := png.Decode(i.file)
		if err != nil {
			return err
		}

		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		if err != nil {
			log.Fatalln(err)
		}

		newFile := bytes.NewBuffer(make([]byte, 0))

		if err := webp.Encode(newFile, image, options); err != nil {
			log.Fatalln(err)
		}
		i.file = newFile
		i.ext = ".webp"
		return nil
	}
	i.ext = ".webp"
	return nil

}
