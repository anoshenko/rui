package rui

import (
	"strconv"
)

// ImageLoadingStatus defines type of status of the image loading
type ImageLoadingStatus int

// Constants which represent return values of the LoadingStatus function of an [Image] view
const (
	// ImageLoading is the image loading status: in the process of loading
	ImageLoading ImageLoadingStatus = 0
	// ImageReady is the image loading status: the image is loaded successfully
	ImageReady ImageLoadingStatus = 1
	// ImageLoadingError is the image loading status: an error occurred while loading
	ImageLoadingError ImageLoadingStatus = 2
)

// Image defines the image that is used for drawing operations on the Canvas.
type Image interface {
	// URL returns the url of the image
	URL() string

	// LoadingStatus returns the status of the image loading:
	//  - ImageLoading (0) - in the process of loading;
	//  - ImageReady (1) - the image is loaded successfully;
	//  - ImageLoadingError (2) - an error occurred while loading.
	LoadingStatus() ImageLoadingStatus

	// LoadingError: if LoadingStatus() == ImageLoadingError then returns the error text, "" otherwise
	LoadingError() string
	setLoadingError(err string)

	// Width returns the width of the image in pixels. While LoadingStatus() != ImageReady returns 0
	Width() float64

	// Height returns the height of the image in pixels. While LoadingStatus() != ImageReady returns 0
	Height() float64
}

type imageData struct {
	url           string
	loadingStatus ImageLoadingStatus
	loadingError  string
	width, height float64
	listener      func(Image)
}

type imageManager struct {
	images map[string]*imageData
}

func (image *imageData) URL() string {
	return image.url
}

func (image *imageData) LoadingStatus() ImageLoadingStatus {
	return image.loadingStatus
}

func (image *imageData) LoadingError() string {
	return image.loadingError
}

func (image *imageData) setLoadingError(err string) {
	image.loadingError = err
}

func (image *imageData) Width() float64 {
	return image.width
}

func (image *imageData) Height() float64 {
	return image.height
}

func (manager *imageManager) loadImage(url string, onLoaded func(Image), session Session) Image {
	if manager.images == nil {
		manager.images = make(map[string]*imageData)
	}

	if image, ok := manager.images[url]; ok && image.loadingStatus == ImageReady {
		return image
	}

	image := new(imageData)
	image.url = url
	image.listener = onLoaded
	image.loadingStatus = ImageLoading
	manager.images[url] = image

	session.callFunc("loadImage", url)
	session.sendResponse()
	return image
}

func (manager *imageManager) imageLoaded(obj DataObject) {
	if manager.images == nil {
		manager.images = make(map[string]*imageData)
		return
	}

	if url, ok := obj.PropertyValue("url"); ok {
		if image, ok := manager.images[url]; ok {
			image.loadingStatus = ImageReady
			if width, ok := obj.PropertyValue("width"); ok {
				if w, err := strconv.ParseFloat(width, 64); err == nil {
					image.width = w
				}
			}
			if height, ok := obj.PropertyValue("height"); ok {
				if h, err := strconv.ParseFloat(height, 64); err == nil {
					image.height = h
				}
			}
			if image.listener != nil {
				image.listener(image)
			}
		}
	}
}

func (manager *imageManager) imageLoadError(obj DataObject) {
	if manager.images == nil {
		manager.images = make(map[string]*imageData)
		return
	}

	if url, ok := obj.PropertyValue("url"); ok {
		if image, ok := manager.images[url]; ok {
			delete(manager.images, url)

			text, _ := obj.PropertyValue("message")
			image.setLoadingError(text)

			if image.listener != nil {
				image.listener(image)
			}
		}
	}
}

// LoadImage starts the async image loading by url
func LoadImage(url string, onLoaded func(Image), session Session) Image {
	if ok, constName := isConstantName(url); ok {
		if image, ok := session.ImageConstant(constName); ok {
			url = image
		}
	}
	return session.imageManager().loadImage(url, onLoaded, session)
}
