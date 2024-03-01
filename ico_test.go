// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/ico

package ico

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"
)

func Test_encode(t *testing.T) {
	bs, err := os.ReadFile("test.jpeg")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(bs))
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	Encode(&buff, img, &Options{[][2]uint8{{32, 32}}})
	fmt.Println(len(buff.Bytes()))
}
