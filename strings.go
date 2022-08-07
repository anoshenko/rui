package rui

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

var stringResources = map[string]map[string]string{}

func scanEmbedStringsDir(fs *embed.FS, dir string) {
	if files, err := fs.ReadDir(dir); err == nil {
		for _, file := range files {
			name := file.Name()
			path := dir + "/" + name
			if file.IsDir() {
				scanEmbedStringsDir(fs, path)
			} else if strings.ToLower(filepath.Ext(name)) == ".rui" {
				if data, err := fs.ReadFile(path); err == nil {
					loadStringResources(string(data))
				} else {
					ErrorLog(err.Error())
				}
			}
		}
	}
}

func scanStringsDir(path string) {
	if files, err := os.ReadDir(path); err == nil {
		for _, file := range files {
			filename := file.Name()
			if filename[0] != '.' {
				newPath := path + `/` + filename
				if file.IsDir() {
					scanStringsDir(newPath)
				} else if strings.ToLower(filepath.Ext(newPath)) == ".rui" {
					if data, err := os.ReadFile(newPath); err == nil {
						loadStringResources(string(data))
					} else {
						ErrorLog(err.Error())
					}
				}
			}
		}
	} else {
		DebugLog(err.Error())
	}
}

func loadStringResources(text string) {
	data := ParseDataText(text)
	if data == nil {
		return
	}

	parseStrings := func(obj DataObject, lang string) {
		table, ok := stringResources[lang]
		if !ok {
			table = map[string]string{}
		}

		for i := 0; i < obj.PropertyCount(); i++ {
			if prop := obj.Property(i); prop != nil && prop.Type() == TextNode {
				table[prop.Tag()] = prop.Text()
			}
		}

		stringResources[lang] = table
	}

	tag := data.Tag()
	if tag == "strings" {
		for i := 0; i < data.PropertyCount(); i++ {
			if prop := data.Property(i); prop != nil && prop.Type() == ObjectNode {
				parseStrings(prop.Object(), prop.Tag())
			}
		}

	} else if strings.HasPrefix(tag, "strings:") {
		if lang := tag[8:]; lang != "" {
			parseStrings(data, lang)
		}
	}
}

// GetString returns the text for the language which is defined by "lang" parameter
func GetString(tag, lang string) (string, bool) {
	if table, ok := stringResources[lang]; ok {
		if text, ok := table[tag]; ok {
			return text, true
		}
		DebugLogF(`There is no "%s" string resource`, tag)
	}
	DebugLogF(`There are no "%s" language resources`, lang)
	return tag, false
}

func (session *sessionData) GetString(tag string) (string, bool) {
	getString := func(tag, lang string) (string, bool) {
		if table, ok := stringResources[lang]; ok {
			if text, ok := table[tag]; ok {
				return text, true
			}
			DebugLogF(`There is no "%s" string in "%s" resources`, tag, lang)
		}
		return tag, false
	}

	if session.language != "" {
		if text, ok := getString(tag, session.language); ok {
			return text, true
		}
	}

	if session.languages != nil {
		for _, lang := range session.languages {
			if lang != session.language {
				if text, ok := getString(tag, lang); ok {
					return text, true
				}
			}
		}
	}

	return tag, false
}
