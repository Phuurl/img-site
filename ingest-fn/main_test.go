package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGetImageSize(t *testing.T) {
	var tests = []struct {
		width     int
		height    int
		imageType string
		path      string
	}{
		{500, 500, "image/jpeg", ".test-resources/500x500.jpg"},
		{500, 500, "image/png", ".test-resources/500x500.png"},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%dx%d %s", test.width, test.height, test.imageType)
		t.Run(testName, func(t *testing.T) {
			w, h, imageType, err := getImageSize(test.path)
			if err != nil {
				t.Errorf("err: %v", err.Error())
			} else if w != test.width {
				t.Errorf("width expected %d, got %d", test.width, w)
			} else if h != test.height {
				t.Errorf("height expected %d, got %d", test.height, h)
			} else if imageType != test.imageType {
				t.Errorf("imageType expected %s, got %s", test.imageType, imageType)
			}
		})
	}
}

func TestCreateHtml(t *testing.T) {
	var testValues = []templateValues{
		{
			Domain:    "example.com",
			Id:        "500-501-png",
			ImageFile: "500x501.png",
			Width:     500,
			Height:    501,
			SiteName:  "Test",
			Type:      "image/png",
		},
		{
			Domain:    "example.com",
			Id:        "501-500-jpeg",
			ImageFile: "501x500.jpg",
			Width:     501,
			Height:    500,
			SiteName:  "Test",
			Type:      "image/jpeg",
		},
	}

	for _, test := range testValues {
		t.Run(test.Id, func(t *testing.T) {
			outPath := t.TempDir() + "out.html"
			err := createHtml(outPath, "template.gohtml", test)
			if err != nil {
				t.Errorf("err: %v", err.Error())
			}

			resultHtml, err := os.ReadFile(outPath)
			if err != nil {
				t.Errorf("error reading created file: %v", err.Error())
			}
			expectedHtml, err := os.ReadFile(".test-resources/" + test.Id + ".html")
			if err != nil {
				t.Fatalf("")
			}
			if string(resultHtml) != string(expectedHtml) {
				t.Errorf("html mismatch with expected file .test-resources/%s.html", test.Id)
			}
		})
	}
}
