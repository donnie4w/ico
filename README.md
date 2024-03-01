# ICO Encoding Library

This is a Go-implemented library designed for encoding images into the ICO format.

#### Usage Guide

1. To acquire the library, run:
   ```
   go get github.com/donnie4w/ico
   ```

2. Using the `Encode` Function:

   ```go
   import "github.com/donnie4w/ico"

   // Read image data from a file
   bs, err := os.ReadFile("test.jpeg")
   img, _, err := image.Decode(bytes.NewReader(bs))
   var buff bytes.Buffer

   // Encode an ICO image with one image at its inherent dimensions
   ico.Encode(&buff, img, nil)

   // Encode an ICO image with multiple images at specified sizes
   // Sizes: 16x16 pixels, 32x32 pixels, and the original image's size (three sizes altogether)
   ico.Encode(&buff, img, &ico.Options{[][2]uint8{{16, 16}, {32, 32}}})
   ```