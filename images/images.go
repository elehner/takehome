package images

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"golang.org/x/image/draw"
)

const (
	ErrorMethodNotSupported = "Only POST is supported"
	ErrorDecodingImage      = "Error while extracting image"
	ErrorEncodingImage      = "Error while converting image"
)

// HandleImageRequest directs the request to the appropriate call based
// on the request method.
func HandleImageRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleImageProcessing(w, r)
	default:
		http.Error(w, ErrorMethodNotSupported, http.StatusMethodNotAllowed)
	}
}

func handleImageProcessing(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	jpeg, err := jpeg.Decode(body)
	if err != nil {
		http.Error(w, ErrorDecodingImage, http.StatusBadRequest)
		return
	}

	resizedImage := resizeImage(jpeg)

	newPngBuffer := new(bytes.Buffer)
	if err := png.Encode(newPngBuffer, resizedImage); err != nil {
		http.Error(w, ErrorEncodingImage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")
	w.Write(newPngBuffer.Bytes())
}

func resizeImage(img image.Image) image.Image {
	maxWidth, maxHeight := 256, 256

	//Determine initial bounds
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// image within bounds and does not need to be resized
	if width <= maxWidth && height <= maxHeight {
		return img
	}
	// Setup the aspect ratio
	newWidth, newHeight := maxWidth, maxHeight
	if width < height {
		newWidth = (width * maxWidth) / height
		if newWidth < 1 {
			newWidth = 1
		}
	} else if height < width {
		newHeight = (height * maxHeight) / width
		if newHeight < 1 {
			newHeight = 1
		}
	}

	// Draw the newly sized the image
	scaledImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.NearestNeighbor.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)

	return scaledImg
}
