package images

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandleImageProcessingErrorsOnBadMethod(t *testing.T) {
	imgBuffer := new(bytes.Buffer)

	req := httptest.NewRequest("GET", "localhost:8080", imgBuffer)
	w := httptest.NewRecorder()

	HandleImageRequest(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if strings.TrimSpace(string(body)) != ErrorMethodNotSupported {
		t.Errorf("Body was %s, expected %s", string(body), ErrorMethodNotSupported)
	}
}

func TestHandleImageProcessingErrorsOnEmptyFile(t *testing.T) {
	imgBuffer := new(bytes.Buffer)

	req := httptest.NewRequest("POST", "localhost:8080", imgBuffer)
	w := httptest.NewRecorder()

	HandleImageRequest(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(body) == ErrorDecodingImage {
		t.Errorf("Body was %s, expected %s", string(body), ErrorDecodingImage)
	}
}

func TestHandleImageProcessingErrorsOnBadFile(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	imgBuffer := new(bytes.Buffer)
	err := png.Encode(imgBuffer, img)
	if err != nil {
		t.Errorf("Error encoding image: %v", err)
	}

	req := httptest.NewRequest("POST", "localhost:8080", imgBuffer)
	w := httptest.NewRecorder()

	HandleImageRequest(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(body) == ErrorDecodingImage {
		t.Errorf("Body was %s, expected %s", string(body), ErrorDecodingImage)
	}
}

func TestHandleImageProcessingHandlesRealImage(t *testing.T) {
	testImg, err := os.Open("./test_images/test_image.jpeg")
	if err != nil {
		t.Errorf("Error pulling test image: %v", err)
	}

	req := httptest.NewRequest("POST", "localhost:8080", testImg)
	w := httptest.NewRecorder()

	HandleImageRequest(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but was %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	expectedImg, err := os.ReadFile("./test_images/resized_image.png")
	if err != nil {
		t.Errorf("Error pulling expected image: %v", err)
	}
	if !bytes.Equal(body, expectedImg) {
		t.Error("Images differ")
	}
}

func TestResizeImage(t *testing.T) {
	tests := []struct {
		image     image.Image
		expectedX int
		expectedY int
	}{
		// Returns same image if image is already within the limits
		{
			image:     image.NewRGBA(image.Rect(0, 0, 100, 100)),
			expectedX: 100,
			expectedY: 100,
		},
		// Returns image with a minimum pixel of 1 for height/length
		{
			image:     image.NewRGBA(image.Rect(0, 0, 1, 1000)),
			expectedX: 1,
			expectedY: 256,
		},
		{
			image:     image.NewRGBA(image.Rect(0, 0, 1000, 1)),
			expectedX: 256,
			expectedY: 1,
		},
		// Returns resized image while maintaing aspect ratio
		{
			image:     image.NewRGBA(image.Rect(0, 0, 500, 1000)),
			expectedX: 128,
			expectedY: 256,
		},
		{
			image:     image.NewRGBA(image.Rect(0, 0, 1000, 500)),
			expectedX: 256,
			expectedY: 128,
		},
		{
			image:     image.NewRGBA(image.Rect(0, 0, 300, 300)),
			expectedX: 256,
			expectedY: 256,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("resizeImage=%d", i), func(t *testing.T) {
			resizedImage := resizeImage(test.image)
			if resizedImage.Bounds().Dx() != test.expectedX ||
				resizedImage.Bounds().Dy() != test.expectedY {
				t.Errorf(
					"Bounds differed. Received %d, %d. Expected %d, %d.",
					resizedImage.Bounds().Dx(),
					resizedImage.Bounds().Dy(),
					test.expectedX,
					test.expectedY,
				)
			}
		})
	}
}
