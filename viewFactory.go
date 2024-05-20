package rui

import (
	"os"
	"path/filepath"
	"strings"
)

var viewCreators = map[string]func(Session) View{
	"View":           newView,
	"ColumnLayout":   newColumnLayout,
	"ListLayout":     newListLayout,
	"GridLayout":     newGridLayout,
	"StackLayout":    newStackLayout,
	"TabsLayout":     newTabsLayout,
	"AbsoluteLayout": newAbsoluteLayout,
	"Resizable":      newResizable,
	"DetailsView":    newDetailsView,
	"TextView":       newTextView,
	"Button":         newButton,
	"Checkbox":       newCheckbox,
	"DropDownList":   newDropDownList,
	"ProgressBar":    newProgressBar,
	"NumberPicker":   newNumberPicker,
	"ColorPicker":    newColorPicker,
	"DatePicker":     newDatePicker,
	"TimePicker":     newTimePicker,
	"FilePicker":     newFilePicker,
	"EditView":       newEditView,
	"ListView":       newListView,
	"CanvasView":     newCanvasView,
	"ImageView":      newImageView,
	"SvgImageView":   newSvgImageView,
	"TableView":      newTableView,
	"AudioPlayer":    newAudioPlayer,
	"VideoPlayer":    newVideoPlayer,
}

// RegisterViewCreator register function of creating view
func RegisterViewCreator(tag string, creator func(Session) View) bool {
	builtinViews := []string{
		"View",
		"ViewsContainer",
		"ColumnLayout",
		"ListLayout",
		"GridLayout",
		"StackLayout",
		"TabsLayout",
		"AbsoluteLayout",
		"Resizable",
		"DetailsView",
		"TextView",
		"Button",
		"Checkbox",
		"DropDownList",
		"ProgressBar",
		"NumberPicker",
		"ColorPicker",
		"DatePicker",
		"TimePicker",
		"EditView",
		"ListView",
		"CanvasView",
		"ImageView",
		"TableView",
	}

	for _, name := range builtinViews {
		if name == tag {
			return false
		}
	}

	viewCreators[tag] = creator
	return true
}

// CreateViewFromObject create new View and initialize it by Node data
func CreateViewFromObject(session Session, object DataObject) View {
	tag := object.Tag()

	if creator, ok := viewCreators[tag]; ok {
		if !session.ignoreViewUpdates() {
			session.setIgnoreViewUpdates(true)
			defer session.setIgnoreViewUpdates(false)
		}
		view := creator(session)
		if customView, ok := view.(CustomView); ok {
			if !InitCustomView(customView, tag, session, nil) {
				return nil
			}
		}
		parseProperties(view, object)
		return view
	}

	ErrorLog(`Unknown view type "` + object.Tag() + `"`)
	return nil
}

// CreateViewFromText create new View and initialize it by content of text
func CreateViewFromText(session Session, text string) View {
	if data := ParseDataText(text); data != nil {
		return CreateViewFromObject(session, data)
	}
	return nil
}

// CreateViewFromResources create new View and initialize it by the content of
// the resource file from "views" directory
func CreateViewFromResources(session Session, name string) View {
	if strings.ToLower(filepath.Ext(name)) != ".rui" {
		name += ".rui"
	}

	for _, fs := range resources.embedFS {
		rootDirs := resources.embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, rawDir:
				// do nothing

			case viewDir:
				if data, err := fs.ReadFile(dir + "/" + name); err == nil {
					if data := ParseDataText(string(data)); data != nil {
						return CreateViewFromObject(session, data)
					}
				}

			default:
				if data, err := fs.ReadFile(dir + "/" + viewDir + "/" + name); err == nil {
					if data := ParseDataText(string(data)); data != nil {
						return CreateViewFromObject(session, data)
					}
				}
			}
		}
	}

	if resources.path != "" {
		if data, err := os.ReadFile(resources.path + viewDir + "/" + name); err == nil {
			if data := ParseDataText(string(data)); data != nil {
				return CreateViewFromObject(session, data)
			}
		}
	}

	return nil
}
