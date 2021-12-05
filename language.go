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
	Name string           // Human-friendly name for this language, eg "Python" or "Go"
	Exts []string         // List of file extensions associated with this language
	Lang *sitter.Language // The thing we need to pass to sitter.Parser.SetLanguage()
	// sitter.Node "type" values that correspond to "mutli-branch" nodes.
	// For example, in Python it would be "elif" and "else"
	BranchTypes []string
}

// Alphabetical ist of languages.  LANGUAGE_MAP is built from this:
var LANGUAGES = []LanguageInfo{
	LanguageInfo{
		Name:        "Go",
		Exts:        []string{"go"},
		Lang:        golang.GetLanguage(),
		BranchTypes: []string{}, // TODO: fill this in after testing
	},
	LanguageInfo{
		Name:        "Javascript",
		Exts:        []string{"js"},
		Lang:        javascript.GetLanguage(),
		BranchTypes: []string{},
	},
	LanguageInfo{
		Name:        "Python",
		Exts:        []string{"py"},
		Lang:        python.GetLanguage(),
		BranchTypes: []string{"elif_clause", "else_clause"},
	},
	LanguageInfo{
		Name:        "Ruby",
		Exts:        []string{"rb"},
		Lang:        ruby.GetLanguage(),
		BranchTypes: []string{"elsif", "else", "case", "when"},
	},
	LanguageInfo{
		Name: "YAML",
		Exts: []string{"yaml", "yml"},
		Lang: yaml.GetLanguage(),
	},
}

// Map of file extensions to their languages:
var LANGUAGE_MAP = map[string]LanguageInfo{}

// LANGUAGE_MAP is built fom LANGUAGES by calling initLangMap() when it's needed.
// Once it's initialized, this flag is set to true:
var langMapInitialized bool = false

// given a filename or path, return the language to use for parsing it,
// or an error if we don't know
func LangFromFilename(x string) (LanguageInfo, error) {
	ext := filepath.Ext(x)
	if ext == "" {
		// There was no extension
		return LanguageInfo{}, errors.New(fmt.Sprintf("No extension to detect language for: %v", x))
	}
	ext = ext[1:] // Ext() includes the dot, this removes it
	return LangFromString(ext)
}

// Given a string corresponding to a language file extension, return the language
// to use for parsing it, or an error if we don't know
func LangFromString(x string) (LanguageInfo, error) {
	initLangMap()
	info, found := LANGUAGE_MAP[x]
	if !found {
		return LanguageInfo{}, errors.New(fmt.Sprintf("Unknown language for %v", x))
	}
	return info, nil
}

func LangIsSupported(x string) bool {
	initLangMap()
	_, found := LANGUAGE_MAP[x]
	return found
}

func initLangMap() {
	if langMapInitialized {
		return
	}

	for _, lang := range LANGUAGES {
		for _, ext := range lang.Exts {
			LANGUAGE_MAP[ext] = lang
		}
	}

	langMapInitialized = true
}

// Return true if the node is an "alternate" type, eg an "else" or "else if" block.
func IsMultiBranchNode(n *sitter.Node, lang LanguageInfo) bool {
	for _, v := range lang.BranchTypes {
		if n.Type() == v {
			return true
		}
	}
	return false
}
