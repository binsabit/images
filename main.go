package main

import (
	"github.com/binsabit/images/image"
	"log"
	"os"
)

func GetFiles(filenames []string) chan *os.File {
	fileChan := make(chan *os.File)

	go func() {
		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				log.Println(err)
				continue
			}
			fileChan <- file
		}
		close(fileChan)
	}()

	return fileChan
}

func GetImages(fileChan chan *os.File, opts ...image.OptFunc) chan *image.Image {
	imgChan := make(chan *image.Image, 10)
	go func() {
		for file := range fileChan {
			img, err := image.NewImage(file.Name(), file, opts...)
			if err != nil {
				log.Println(err)
			}
			imgChan <- img
		}
		close(imgChan)
	}()

	return imgChan
}

func main() {
	filename := []string{"./test1.jpeg", "./test2.jpeg"}

	filesChan := GetFiles(filename)

	getImages := GetImages(filesChan, image.Format(image.Webp))

	for img := range getImages {
		err := img.Save("results")
		if err != nil {
			log.Println(err)
		}
	}

}
