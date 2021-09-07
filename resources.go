package rui

import (
	"embed"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	imageDir   = "images"
	themeDir   = "themes"
	viewDir    = "views"
	rawDir     = "raw"
	stringsDir = "strings"
)

type scaledImage struct {
	path  string
	scale float64
}

type imagePath struct {
	path string
	fs   *embed.FS
}

type resourceManager struct {
	embedFS      []*embed.FS
	themes       map[string]*theme
	images       map[string]imagePath
	imageSrcSets map[string][]scaledImage
	path         string
}

var resources = resourceManager{
	embedFS:      []*embed.FS{},
	themes:       map[string]*theme{},
	images:       map[string]imagePath{},
	imageSrcSets: map[string][]scaledImage{},
}

func AddEmbedResources(fs *embed.FS) {
	resources.embedFS = append(resources.embedFS, fs)
	rootDirs := embedRootDirs(fs)
	for _, dir := range rootDirs {
		switch dir {
		case imageDir:
			scanEmbedImagesDir(fs, dir, "")

		case themeDir:
			scanEmbedThemesDir(fs, dir)

		case stringsDir:
			scanEmbedStringsDir(fs, dir)

		case viewDir, rawDir:
			// do nothing

		default:
			if files, err := fs.ReadDir(dir); err == nil {
				for _, file := range files {
					if file.IsDir() {
						switch file.Name() {
						case imageDir:
							scanEmbedImagesDir(fs, dir+"/"+imageDir, "")

						case themeDir:
							scanEmbedThemesDir(fs, dir+"/"+themeDir)

						case stringsDir:
							scanEmbedStringsDir(fs, dir+"/"+stringsDir)

						case viewDir, rawDir:
							// do nothing
						}
					}
				}
			}
		}
	}
}

func embedRootDirs(fs *embed.FS) []string {
	result := []string{}
	if files, err := fs.ReadDir("."); err == nil {
		for _, file := range files {
			if file.IsDir() {
				result = append(result, file.Name())
			}
		}
	}
	return result
}

func scanEmbedThemesDir(fs *embed.FS, dir string) {
	if files, err := fs.ReadDir(dir); err == nil {
		for _, file := range files {
			name := file.Name()
			path := dir + "/" + name
			if file.IsDir() {
				scanEmbedThemesDir(fs, path)
			} else if strings.ToLower(filepath.Ext(name)) == ".rui" {
				if data, err := fs.ReadFile(path); err == nil {
					RegisterThemeText(string(data))
				}
			}
		}
	}
}

func scanEmbedImagesDir(fs *embed.FS, dir, prefix string) {
	if files, err := fs.ReadDir(dir); err == nil {
		for _, file := range files {
			name := file.Name()
			path := dir + "/" + name
			if file.IsDir() {
				scanEmbedImagesDir(fs, path, prefix+name+"/")
			} else {
				ext := strings.ToLower(filepath.Ext(name))
				switch ext {
				case ".png", ".jpg", ".jpeg", ".svg":
					registerImage(fs, path, prefix+name)
				}
			}
		}
	}
}

func invalidImageFileFormat(filename string) {
	ErrorLog(`Invalid image file name parameters: "` + filename +
		`". Image file name format: name[@x-param].ext (examples: icon.png, icon@1.5x.png)`)
}

func registerImage(fs *embed.FS, path, filename string) {
	resources.images[filename] = imagePath{fs: fs, path: path}

	start := strings.LastIndex(filename, "@")
	if start < 0 {
		return
	}

	ext := strings.LastIndex(filename, ".")
	if start > ext || filename[ext-1] != 'x' {
		invalidImageFileFormat(path)
		return
	}

	if scale, err := strconv.ParseFloat(filename[start+1:ext-1], 32); err == nil {
		key := filename[:start] + filename[ext:]
		images, ok := resources.imageSrcSets[key]
		if ok {
			for _, image := range images {
				if image.scale == scale {
					return
				}
			}
		} else {
			images = []scaledImage{}
		}
		resources.imageSrcSets[key] = append(images, scaledImage{path: filename, scale: scale})
	} else {
		invalidImageFileFormat(path)
		return
	}
}

func scanImagesDirectory(path, filePrefix string) {
	if files, err := ioutil.ReadDir(path); err == nil {
		for _, file := range files {
			filename := file.Name()
			if filename[0] != '.' {
				newPath := path + `/` + filename
				if !file.IsDir() {
					registerImage(nil, newPath, filePrefix+filename)
				} else {
					scanImagesDirectory(newPath, filePrefix+filename+"/")
				}
			}
		}
	} else {
		ErrorLog(err.Error())
	}
}

func scanThemesDir(path string) {
	if files, err := ioutil.ReadDir(path); err == nil {
		for _, file := range files {
			filename := file.Name()
			if filename[0] != '.' {
				newPath := path + `/` + filename
				if file.IsDir() {
					scanThemesDir(newPath)
				} else if strings.ToLower(filepath.Ext(newPath)) == ".rui" {
					if data, err := ioutil.ReadFile(newPath); err == nil {
						RegisterThemeText(string(data))
					} else {
						ErrorLog(err.Error())
					}
				}
			}
		}
	} else {
		ErrorLog(err.Error())
	}
}

// SetResourcePath set path of the resource directory
func SetResourcePath(path string) {
	resources.path = path
	pathLen := len(path)
	if pathLen > 0 && path[pathLen-1] != '/' {
		resources.path += "/"
	}

	scanImagesDirectory(resources.path+imageDir, "")
	scanThemesDir(resources.path + themeDir)
	scanStringsDir(resources.path + stringsDir)
}

// RegisterThemeText parse text and add result to the theme list
func RegisterThemeText(text string) bool {
	data := ParseDataText(text)
	if data == nil {
		return false
	}

	if !data.IsObject() {
		ErrorLog(`Root element is not object`)
		return false
	}
	if data.Tag() != "theme" {
		ErrorLog(`Invalid the root object tag. Must be "theme"`)
		return false
	}

	if name, ok := data.PropertyValue("name"); ok && name != "" {
		t := resources.themes[name]
		if t == nil {
			t = new(theme)
			t.init()
			resources.themes[name] = t
		}
		t.addData(data)
	} else {
		defaultTheme.addData(data)
	}

	return true
}

func serveResourceFile(filename string, w http.ResponseWriter, r *http.Request) bool {
	serveEmbed := func(fs *embed.FS, path string) bool {
		if file, err := fs.Open(path); err == nil {
			if stat, err := file.Stat(); err == nil {
				http.ServeContent(w, r, filename, stat.ModTime(), file.(io.ReadSeeker))
				return true
			}
		}
		return false
	}

	if image, ok := resources.images[filename]; ok {
		if image.fs != nil {
			if serveEmbed(image.fs, image.path) {
				return true
			}
		} else {
			if _, err := os.Stat(image.path); err == nil {
				http.ServeFile(w, r, image.path)
				return true
			}
		}
	}

	for _, fs := range resources.embedFS {
		if serveEmbed(fs, filename) {
			return true
		}
		for _, dir := range embedRootDirs(fs) {
			if serveEmbed(fs, dir+"/"+filename) {
				return true
			}
			if subdirs, err := fs.ReadDir(dir); err == nil {
				for _, subdir := range subdirs {
					if subdir.IsDir() {
						if serveEmbed(fs, dir+"/"+subdir.Name()+"/"+filename) {
							return true
						}
					}
				}
			}
		}
	}

	serve := func(path, filename string) bool {
		filepath := path + filename
		if _, err := os.Stat(filepath); err == nil {
			http.ServeFile(w, r, filepath)
			return true
		}

		filepath = path + imageDir + "/" + filename
		if _, err := os.Stat(filepath); err == nil {
			http.ServeFile(w, r, filepath)
			return true
		}

		return false
	}

	if resources.path != "" && serve(resources.path, filename) {
		return true
	}

	if exe, err := os.Executable(); err == nil {
		path := filepath.Dir(exe) + "/resources/"
		if serve(path, filename) {
			return true
		}
	}

	return false
}

func ReadRawResource(filename string) []byte {
	for _, fs := range resources.embedFS {
		rootDirs := embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, viewDir:
				// do nothing

			case rawDir:
				if data, err := fs.ReadFile(dir + "/" + filename); err == nil {
					return data
				}

			default:
				if data, err := fs.ReadFile(dir + "/" + rawDir + "/" + filename); err == nil {
					return data
				}
			}
		}
	}

	readFile := func(path string) []byte {
		if data, err := os.ReadFile(resources.path + rawDir + "/" + filename); err == nil {
			return data
		}
		return nil
	}

	if resources.path != "" {
		if data := readFile(resources.path + rawDir + "/" + filename); data != nil {
			return data
		}
	}

	if exe, err := os.Executable(); err == nil {
		if data := readFile(filepath.Dir(exe) + "/resources/" + rawDir + "/" + filename); data != nil {
			return data
		}
	}

	ErrorLogF(`The raw file "%s" don't found`, filename)
	return nil
}

func AllRawResources() []string {
	result := []string{}

	for _, fs := range resources.embedFS {
		rootDirs := embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, viewDir:
				// do nothing

			case rawDir:
				if files, err := fs.ReadDir(rawDir); err == nil {
					for _, file := range files {
						result = append(result, file.Name())
					}
				}

			default:
				if files, err := fs.ReadDir(dir + "/" + rawDir); err == nil {
					for _, file := range files {
						result = append(result, file.Name())
					}
				}
			}
		}
	}

	if resources.path != "" {
		if files, err := ioutil.ReadDir(resources.path + rawDir); err == nil {
			for _, file := range files {
				result = append(result, file.Name())
			}
		}
	}

	if exe, err := os.Executable(); err == nil {
		if files, err := ioutil.ReadDir(filepath.Dir(exe) + "/resources/" + rawDir); err == nil {
			for _, file := range files {
				result = append(result, file.Name())
			}
		}
	}

	return result
}
