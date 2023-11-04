package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

func main() {
	const width = 400
	const height = 200
	const qrSquare = 50

	// Create a new context
	dc := gg.NewContext(width, height)

	// Load your company logo
	logoImage, err := gg.LoadImage("berry-global.png") // Replace with your logo's path
	if err != nil {
		log.Fatalf("Error loading logo: %v", err)
	}

	// Generate a barcode
	rawBarCode, err := code128.Encode("RMP0001")
	if err != nil {
		log.Fatalf("Error generating barcode: %v", err)
	}

	// Scale the barcode to the desired width and height
	scaledBarCode, err := barcode.Scale(rawBarCode, width-20, 50)
	if err != nil {
		log.Fatalf("Error scaling barcode: %v", err)
	}

	// Generate a QR code
	qrCode, err := qrcode.New("RMP0001", qrcode.Highest)
	if err != nil {
		log.Fatalf("Error generating QR code: %v", err)
	}

	qrCodeImage := qrCode.Image(qrSquare) // generates an image.Image

	// Draw the company logo
	dc.DrawImage(logoImage, 10, 10)

	// Convert context to an *image.RGBA to draw the scaledBarCode and QR code
	rgba := dc.Image().(*image.RGBA)

	// Draw the barcode at the bottom-center
	barCodeBounds := scaledBarCode.Bounds()
	barCodeWidth := barCodeBounds.Dx()
	barCodeHeight := barCodeBounds.Dy()
	barCodePoint := image.Point{X: width/2 - barCodeWidth/2, Y: height - barCodeHeight - 10}
	barCodeRect := image.Rectangle{
		Min: barCodePoint,
		Max: barCodePoint.Add(image.Point{X: barCodeWidth, Y: barCodeHeight}),
	}
	draw.Draw(rgba, barCodeRect, scaledBarCode, image.Point{}, draw.Src)

	// Draw the QR code at the top-right
	qrCodePoint := image.Point{X: width - 125 - 10, Y: 10} // Padding of 10 from top and right
	qrCodeRect := image.Rectangle{
		Min: qrCodePoint,
		Max: qrCodePoint.Add(image.Point{X: 125, Y: 125}),
	}
	draw.Draw(rgba, qrCodeRect, qrCodeImage, image.Point{}, draw.Src)

	// Convert back to gg's context to continue drawing with gg
	dc = gg.NewContextForRGBA(rgba)

	// Draw product details
	dc.SetColor(color.Black)
	dc.DrawStringAnchored("Product Name: XYZ", 10, 30, 0, 0.5)
	dc.DrawStringAnchored("Product Code: 1234", 10, 50, 0, 0.5)
	dc.DrawStringAnchored("Price: $99.99", (width/2)-10, 30, 0, 0.5)
	dc.DrawStringAnchored("Expiry Date: 01/01/2030", (width/2)-10, 50, 0, 0.5)

	// Save the sticker image
	if err := dc.SavePNG("sticker.png"); err != nil {
		log.Fatalf("Error saving PNG: %v", err)
	}
}
