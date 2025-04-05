package processor

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	// "image/jpeg"

	"github.com/disintegration/imaging"
)

func ResizeImage (data []byte, width, height int)([]byte , error) {

	log.Printf("Input: %d bytes, starts with: %x", len(data), data[:4])

	if len(data) == 0 {
		return nil ,errors.New("empty input data")
	}

	if width <= 0 || height <= 0 {
		return nil, errors.New("invalid dimensions")
	}

	img, err := imaging.Decode(bytes.NewReader(data)) 
	if err != nil {
		log.Printf("Decode failed. First 16 bytes: %x", data[:16])
		return nil, fmt.Errorf("decode error: %v", err) 
	}

	bounds := img.Bounds()
    log.Printf("Original dimensions: %dx%d", bounds.Dx(), bounds.Dy())


	resized := imaging.Resize(img ,width,height,imaging.Lanczos) 
	// if err != nil {
	// 	img,err = jpeg.Decode(bytes.NewReader(data))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("decode failed:%v",err)
	// 	}
	// }

	var buf bytes.Buffer 
	err = imaging.Encode(&buf, resized, imaging.JPEG)
	if err != nil {
        return nil, fmt.Errorf("encode failed: %v", err)
    }
	
	if buf.Len() < 100 {
		return nil ,errors.New("suspiciously small o/p")
	}
	return buf.Bytes() , nil
}