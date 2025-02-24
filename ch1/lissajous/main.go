// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//!-main
// Packages not needed by version in book.

//!+main

var green = color.RGBA{0, 255, 0, 0xff}
var blue = color.RGBA{0, 0, 255, 0xff}
var red = color.RGBA{255, 0, 0, 0xff}
var palette = []color.Color{color.Black, green, color.White, blue, red}

const (
	whiteIndex = 2 // third color in palette
	blackIndex = 0 // first color in palette
	greenIndex = 1 // second color in palette
	blueIndex = 3 // third color in palette
)

func main() {
	//!-main
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			cycles, ok := r.URL.Query()["cycles"]
			if !ok || len(cycles[0]) < 1 {
					fmt.Println("something wrong with query param.")
					lissajous(w, 5)
			} else {
					val, err :=  strconv.Atoi(cycles[0])
					if err!=nil {
							fmt.Fprintf(os.Stdout, "Parse Int failed, %v", err)
							lissajous(w, 5)
					}
					lissajous(w,val)
			}																							
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	//!+main
	lissajous(os.Stdout, 5)
}

func lissajous(out io.Writer, inputCycles int) {
	const (
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	var cycles float64= float64(inputCycles)     // number of complete x oscillator revolutions
	rand.Seed(time.Now().Unix())
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(rand.Intn(len(palette))))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
