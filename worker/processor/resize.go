package processor

import (
	"bytes"
		
	"github.com/disintegration/imaging"
)

func ResizeImage (data []byte, width, height int)([]byte , error) {
	img, err := imaging.Decode(bytes.NewReader(data)) 
	if err != nil {
		return nil , err
	}

	resized := imaging.Resize(img ,width,height,imaging.Lanczos) 

	var buf bytes.Buffer 
	err = imaging.Encode(&buf, resized, imaging.JPEG)
	return buf.Bytes() , err
}