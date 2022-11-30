package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	filename = flag.String("f", "manifest.txt", "name of manifest file")
	output   = flag.String("o", "background.png", "name of output file")

	imgwidth  = flag.Int("iw", 1200*2, "width of output image")
	imgheight = flag.Int("ih", 630*2, "height of output image")

	dotsize = flag.Int("ds", 16, "size of each dot")
)

func main() {
	flag.Parse()

	// Make an empty image with all white pixels.
	img := image.NewRGBA(image.Rect(0, 0, *imgwidth, *imgheight))
	drawDot(img, 0, 0, *imgwidth, *imgheight, color.RGBA{255, 255, 255, 255})

	f, err := os.Open(*filename)
	catch(err)
	defer f.Close()

	// Read and loop over each line in the verdict manifest file.
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	var px, py int
	for sc.Scan() {
		switch sc.Text() {
		case "Accepted":
			// Draw a green square for Accepted submissions.
			drawDot(img, px, py, px+*dotsize, py+*dotsize, color.RGBA{0x26, 0xc2, 0x81, 255})
		default:
			// Draw a red square for all other submissions.
			drawDot(img, px, py, px+*dotsize, py+*dotsize, color.RGBA{0xe3, 0x5b, 0x5a, 255})
		}

		// Determine the position of the next square.
		px += *dotsize
		if px > *imgwidth {
			px = 0
			py += *dotsize
		}
	}
	catch(sc.Err())

	// Save the image as PNG.
	outf, err := os.Create(*output)
	catch(err)
	defer outf.Close()
	err = png.Encode(outf, img)
	catch(err)
}

func drawDot(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	// Set all pixels within (x0, y0) and (x1, y1) to the colour c, except the right-most column and bottom-most row of pixels.
	for y := y0; y < y1-1; y++ {
		for x := x0; x < x1-1; x++ {
			img.Set(x, y, c)
		}
	}
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
