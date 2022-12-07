package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/webp"
)

var (
	overrideFile bool = false
	removeFile   bool = false
	pngOut       bool = false
)

func webp2jpeg(path_webp string, img_webp image.Image) {
	// init JPEG
	img_jpeg := image.NewNRGBA(img_webp.Bounds())

	// paste WEBP over to new image
	draw.Draw(img_jpeg, img_jpeg.Bounds(), img_webp, img_jpeg.Bounds().Min, draw.Over)

	// create file descriptor
	path_jpeg := strings.TrimSuffix(path_webp, ".webp") + ".jpeg"
	if _, err := os.Stat(path_jpeg); !overrideFile && err == nil {
		// file exists
		fmt.Printf("File exists: %s , Override? [Yes (y) / No (N) / Skip (s)] ", filepath.Base(path_jpeg))

		var answer string
		fmt.Scan(&answer)

		switch strings.ToUpper(answer) {
		case "Y":
			// override... don't do anthing
			break
		case "S":
			// skip, don't do anything
			return
		default:
			// rename
			path_jpeg = fmt.Sprintf("%s-%d.jpeg", strings.TrimSuffix(path_jpeg, ".jpeg"), time.Now().UnixNano()%100000)
		}
	}
	fmt.Printf("> Writing.... %s\n", path_jpeg)

	file, err := os.Create(path_jpeg)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// write JPEG
	opts := &jpeg.Options{
		Quality: 80,
	}

	err = jpeg.Encode(file, img_jpeg, opts)
	if err != nil {
		log.Fatalln(err)
	}

}

func processFile(path string) {
	fmt.Println(path)

	input, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		input.Close()

		// should remove WEBP ?
		if !removeFile {
			return
		}

		fmt.Printf("> Removing... %s\n", path)
		err = os.Remove(path)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	webpImg, err := webp.Decode(input)
	if err != nil {
		log.Fatalln(err)
	}

	if pngOut {
		// TODO: implement
		return
	}

	webp2jpeg(path, webpImg)
}

func processDirectory(dirpath string) {
	// walk WEBPs
	matches, err := filepath.Glob(filepath.Join(dirpath, "*.webp"))
	if err != nil {
		log.Fatalln(err)
	}

	for _, match := range matches {
		processFile(match)
	}
}

func main() {
	flag.BoolVar(&overrideFile, "override", false, "Override dest file if exists")
	flag.BoolVar(&pngOut, "png", false, "Output PNG")
	flag.BoolVar(&removeFile, "remove", false, "Remove file after converting")
	flag.Parse()

	dest := flag.Arg(0)
	if dest == "" {
		dest = "." // assume current dir
	}

	path, err := filepath.Abs(dest)
	if err != nil {
		log.Fatalln(err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}

	if fi.IsDir() {
		processDirectory(path)
	} else {
		processFile(path)
	}
}
