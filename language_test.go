package main

import (
	"testing"
)

func TestLangIsSupported(t *testing.T) {
	supported := []string{"py", "go", "rb"}
	for _, x := range supported {
		val := LangIsSupported(x)
		if val != true {
			t.Fatalf("Expected true for supported language %v, got %v", x, val)
		}
	}

	val := LangIsSupported("fake")
	if val != false {
		t.Fatalf("Expected false for unsupported, got %v", val)
	}
}

func TestLangFromFilenameNo(t *testing.T) {
	filenames := [][]string{
		[]string{"sample.go", "Go"},
		[]string{"sample_files/sample.py", "Python"},
		[]string{"/absolute/path/sample.yml", "YAML"},
	}
	for _, x := range filenames {
		lang, _ := LangFromFilename(x[0])
		if lang.Name != x[1] {
			t.Fatalf("Expected %v for %v, got %v", x[1], x[0], lang.Name)
		}
	}
}

func TestLangFromFilenameNoExtension(t *testing.T) {
	_, err := LangFromFilename("noextension")
	if err.Error() != "No extension to detect language for: noextension" {
		t.Fatalf("Expected error for no file extension, got %v", err)
	}
}
