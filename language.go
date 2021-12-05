package main

import (
	"errors"
	"fmt"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/smacker/go-tree-sitter/yaml"
)

type LanguageInfo struct {
	Name string
	Lang *sitter.Language
}

// Map of file extensions to their languages:
var LANGUAGE_MAP = map[string]LanguageInfo{
	"go":   LanguageInfo{"go", golang.GetLanguage()},
	"js":   LanguageInfo{"javascript", javascript.GetLanguage()},
	"py":   LanguageInfo{"python", python.GetLanguage()},
	"rb":   LanguageInfo{"ruby", ruby.GetLanguage()},
	"yaml": LanguageInfo{"yaml", yaml.GetLanguage()},
	"yml":  LanguageInfo{"yaml", yaml.GetLanguage()},
}

// given a filename or path, return the language to use for parsing it,
// or an error if we don't know
func LangFromFilename(x string) (*sitter.Language, error) {
	ext := filepath.Ext(x)
	if ext == "" {
		// There was no extension
		return nil, errors.New(fmt.Sprintf("No extension to detect language for: %v", x))
	}
	ext = ext[1:] // Ext() includes the dot, this removes it
	return LangFromString(ext)
}

// Given a string corresponding to a language file extension, return the language
// to use for parsing it, or an error if we don't know
func LangFromString(x string) (*sitter.Language, error) {
	info, found := LANGUAGE_MAP[x]
	if !found {
		return nil, errors.New(fmt.Sprintf("Unknown language for %v", x))
	}
	return info.Lang, nil
}

func LangIsSupported(x string) bool {
	_, found := LANGUAGE_MAP[x]
	return found
}
