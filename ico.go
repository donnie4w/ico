// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/ico
package ico

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"io"
	"math"

	"github.com/disintegration/imaging"
)

type icoHeader struct {
	Reserved  [2]byte
	ImageType [2]byte
	Count     uint16
}

func (h icoHeader) Size() int {
	return 6 // 2 reserved bytes + 2 type bytes + 2 count bytes
}

type icoEntry struct {
	Width        uint8
	Height       uint8
	ColorCount   uint8
	Reserved     uint8
	Planes       uint16 // must be 0 or 1
	BitsPerPixel uint16
	SizeInBytes  uint32 // Size of the image data
	Offset       uint32 // The offset of the image data relative to the beginning of the ICO file
}

func (e icoEntry) Size() int {
	return 16
}

// encoding parameters
type Options struct {
	Thumbnails [][2]uint8  //specifications for multiple images. eg, {64,64} means 64x64 pixel 
}

func Encode(w io.Writer, img image.Image, options *Options) (err error) {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()
	var newWidth, newHeight uint8
	if originalWidth > 255 || originalHeight > 255 {
		ratio := math.Min(float64(255)/float64(originalWidth), float64(255)/float64(originalHeight))
		newWidth = uint8(math.Floor(ratio * float64(originalWidth)))
		newHeight = uint8(math.Floor(ratio * float64(originalHeight)))
	} else {
		newWidth = uint8(originalWidth)
		newHeight = uint8(originalHeight)
	}

	colorCount := uint8(0)
	bitsPerPixel := uint16(32)

	if options == nil {
		options = &Options{Thumbnails: [][2]uint8{}}
	} else if options.Thumbnails == nil {
		options.Thumbnails = [][2]uint8{}
	}

	options.Thumbnails = append(options.Thumbnails, [][2]uint8{{newWidth, newHeight}}...)

	header := icoHeader{[2]byte{0x00, 0x00}, [2]byte{0x01, 0x00}, uint16(len(options.Thumbnails))}

	if err = binary.Write(w, binary.LittleEndian, header); err != nil {
		return
	}

	Offset := uint32(header.Size() + len(options.Thumbnails)*icoEntry{}.Size())
	var dataBuf bytes.Buffer
	for _, thumbnails := range options.Thumbnails {
		imageData := getImageByThumbnails(img, thumbnails)
		dataBuf.Write(imageData)
		entry := icoEntry{
			Width:        thumbnails[0],
			Height:       thumbnails[1],
			ColorCount:   colorCount,
			Reserved:     0,
			Planes:       1,
			BitsPerPixel: bitsPerPixel,
			SizeInBytes:  uint32(len(imageData)),
			Offset:       Offset,
		}
		Offset += uint32(len(imageData))
		if err = binary.Write(w, binary.LittleEndian, entry); err != nil {
			return
		}
	}
	w.Write(dataBuf.Bytes())
	return
}

func getImageByThumbnails(img image.Image, thumbnails [2]uint8) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, imaging.Resize(img, int(thumbnails[0]), int(thumbnails[1]), imaging.MitchellNetravali))
	return buf.Bytes()
}
