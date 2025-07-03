package rui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	imageDir   = "images"
	themeDir   = "themes"
	viewDir    = "views"
	popupDir   = "popups"
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
	themes       map[string]Theme
	images       map[string]imagePath
	imageSrcSets map[string][]scaledImage
	path         string
}

var resources = resourceManager{
	embedFS:      []*embed.FS{},
	themes:       map[string]Theme{},
	images:       map[string]imagePath{},
	imageSrcSets: map[string][]scaledImage{},
}

// AddEmbedResources adds embedded resources to the list of application resources
func AddEmbedResources(fs *embed.FS) {
	resources.embedFS = append(resources.embedFS, fs)
	rootDirs := resources.embedRootDirs(fs)
	for _, dir := range rootDirs {
		switch dir {
		case imageDir:
			resources.scanEmbedImagesDir(fs, dir, "")

		case themeDir:
			resources.scanEmbedThemesDir(fs, dir)

		case stringsDir:
			resources.scanEmbedStringsDir(fs, dir)

		case viewDir, rawDir:
			// do nothing

		default:
			if files, err := fs.ReadDir(dir); err == nil {
				for _, file := range files {
					if file.IsDir() {
						switch file.Name() {
						case imageDir:
							resources.scanEmbedImagesDir(fs, dir+"/"+imageDir, "")

						case themeDir:
							resources.scanEmbedThemesDir(fs, dir+"/"+themeDir)

						case stringsDir:
							resources.scanEmbedStringsDir(fs, dir+"/"+stringsDir)

						case viewDir, rawDir:
							// do nothing
						}
					}
				}
			}
		}
	}
}

func (resources *resourceManager) embedRootDirs(fs *embed.FS) []string {
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

func (resources *resourceManager) scanEmbedThemesDir(fs *embed.FS, dir string) {
	if files, err := fs.ReadDir(dir); err == nil {
		for _, file := range files {
			name := file.Name()
			path := dir + "/" + name
			if file.IsDir() {
				resources.scanEmbedThemesDir(fs, path)
			} else if strings.ToLower(filepath.Ext(name)) == ".rui" {
				if data, err := fs.ReadFile(path); err == nil {
					resources.registerThemeText(string(data))
				}
			}
		}
	}
}

func (resources *resourceManager) scanEmbedImagesDir(fs *embed.FS, dir, prefix string) {
	if files, err := fs.ReadDir(dir); err == nil {
		for _, file := range files {
			name := file.Name()
			path := dir + "/" + name
			if file.IsDir() {
				resources.scanEmbedImagesDir(fs, path, prefix+name+"/")
			} else {
				ext := strings.ToLower(filepath.Ext(name))
				switch ext {
				case ".png", ".jpg", ".jpeg", ".svg", ".gif", ".bmp", ".webp":
					resources.registerImage(fs, path, prefix+name)
				}
			}
		}
	}
}

func invalidImageFileFormat(filename string) {
	ErrorLog(`Invalid image file name parameters: "` + filename +
		`". Image file name format: name[@x-param].ext (examples: icon.png, icon@1.5x.png)`)
}

func (resources *resourceManager) registerImage(fs *embed.FS, path, filename string) {
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

func (resources *resourceManager) scanImagesDirectory(path, filePrefix string) {
	if files, err := os.ReadDir(path); err == nil {
		for _, file := range files {
			filename := file.Name()
			if filename[0] != '.' {
				newPath := path + `/` + filename
				if !file.IsDir() {
					resources.registerImage(nil, newPath, filePrefix+filename)
				} else {
					resources.scanImagesDirectory(newPath, filePrefix+filename+"/")
				}
			}
		}
	} else {
		ErrorLog(err.Error())
	}
}

func (resources *resourceManager) scanThemesDir(path string) {
	if files, err := os.ReadDir(path); err == nil {
		for _, file := range files {
			filename := file.Name()
			if filename[0] != '.' {
				newPath := path + `/` + filename
				if file.IsDir() {
					resources.scanThemesDir(newPath)
				} else if strings.ToLower(filepath.Ext(newPath)) == ".rui" {
					if data, err := os.ReadFile(newPath); err == nil {
						resources.registerThemeText(string(data))
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

	resources.scanImagesDirectory(resources.path+imageDir, "")
	resources.scanThemesDir(resources.path + themeDir)
	resources.scanStringsDir(resources.path + stringsDir)
}

func (resources *resourceManager) scanDefaultResourcePath() {
	if exe, err := os.Executable(); err == nil {
		path := filepath.Dir(exe) + "/resources/"
		resources.scanImagesDirectory(path+imageDir, "")
		resources.scanThemesDir(path + themeDir)
		resources.scanStringsDir(path + stringsDir)
	}
}

func (resources *resourceManager) registerThemeText(text string) bool {
	theme, ok := CreateThemeFromText(text)
	if !ok {
		return false
	}

	name := theme.Name()
	if name == "" {
		defaultTheme.Append(theme)
	} else if t, ok := resources.themes[name]; ok {
		t.Append(theme)
	} else {
		resources.themes[name] = theme
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
		for _, dir := range resources.embedRootDirs(fs) {
			if serveEmbed(fs, dir+"/"+filename) {
				return true
			}
			if subDirs, err := fs.ReadDir(dir); err == nil {
				for _, subdir := range subDirs {
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

// ReadRawResource returns the contents of the raw resource with the specified name
func ReadRawResource(filename string) []byte {
	for _, fs := range resources.embedFS {
		rootDirs := resources.embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, viewDir, stringsDir:
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

	if resources.path != "" {
		if data, err := os.ReadFile(resources.path + rawDir + "/" + filename); err == nil {
			return data
		}
	}

	if exe, err := os.Executable(); err == nil {
		if data, err := os.ReadFile(filepath.Dir(exe) + "/resources/" + rawDir + "/" + filename); err == nil {
			return data
		}
	}

	ErrorLogF(`The "%s" raw file don't found`, filename)
	return nil
}

// OpenRawResource returns the contents of the raw resource with the specified name
func OpenRawResource(filename string) fs.File {
	for _, fs := range resources.embedFS {
		rootDirs := resources.embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, viewDir, stringsDir:
				// do nothing

			case rawDir:
				if file, err := fs.Open(dir + "/" + filename); err == nil {
					return file
				}

			default:
				if file, err := fs.Open(dir + "/" + rawDir + "/" + filename); err == nil {
					return file
				}
			}
		}
	}

	if resources.path != "" {
		if file, err := os.Open(resources.path + rawDir + "/" + filename); err == nil {
			return file
		}
	}

	if exe, err := os.Executable(); err == nil {
		if file, err := os.Open(filepath.Dir(exe) + "/resources/" + rawDir + "/" + filename); err == nil {
			return file
		}
	}

	ErrorLogF(`The "%s" raw file don't found`, filename)
	return nil
}

// AllRawResources returns the list of all raw resources
func AllRawResources() []string {
	result := []string{}

	for _, fs := range resources.embedFS {
		rootDirs := resources.embedRootDirs(fs)
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
		if files, err := os.ReadDir(resources.path + rawDir); err == nil {
			for _, file := range files {
				result = append(result, file.Name())
			}
		}
	}

	if exe, err := os.Executable(); err == nil {
		if files, err := os.ReadDir(filepath.Dir(exe) + "/resources/" + rawDir); err == nil {
			for _, file := range files {
				result = append(result, file.Name())
			}
		}
	}

	return result
}

// AllImageResources returns the list of all image resources
func AllImageResources() []string {
	result := make([]string, 0, len(resources.images))
	for image := range resources.images {
		result = append(result, image)
	}
	sort.Strings(result)
	return result
}

// AddTheme adds theme to application
func AddTheme(theme Theme) {
	if theme != nil {
		name := theme.Name()
		if name == "" {
			defaultTheme.Append(theme)
		} else if t, ok := resources.themes[name]; ok {
			t.Append(theme)
		} else {
			resources.themes[name] = theme
		}
	}
}
