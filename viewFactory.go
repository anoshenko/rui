package rui

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

var systemViewCreators = map[string]func(Session) View{
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

// ViewCreateListener is the listener interface of a view create event
type ViewCreateListener interface {
	// OnCreate is a function of binding object that is called by the CreateViewFromText, CreateViewFromResources,
	// and CreateViewFromObject functions after the creation of a view
	OnCreate(view View)
}

var viewCreate map[string]func(Session) View = nil

func viewCreators() map[string]func(Session) View {
	if viewCreate == nil {
		viewCreate = map[string]func(Session) View{}
		for tag, fn := range systemViewCreators {
			viewCreate[strings.ToLower(tag)] = fn
		}
	}
	return viewCreate
}

// RegisterViewCreator register function of creating view
func RegisterViewCreator(tag string, creator func(Session) View) bool {
	loTag := strings.ToLower(tag)
	for name := range systemViewCreators {
		if name == loTag {
			ErrorLog(`It is forbidden to override the function of ` + tag + ` creating`)
			return false
		}
	}

	viewCreators()[loTag] = creator
	return true
}

// CreateViewFromObject create new View and initialize it by DataObject data. Parameters:
//   - session - the session to which the view will be attached (should not be nil);
//   - object - data describing View;
//   - binding - object assigned to the Binding property (may be nil).
//
// If the function fails, it returns nil and an error message is written to the log.
func CreateViewFromObject(session Session, object DataObject, binding any) View {
	if session == nil {
		ErrorLog(`Session must not be nil`)
		return nil
	}

	tag := object.Tag()
	creator, ok := viewCreators()[strings.ToLower(tag)]
	if !ok {
		ErrorLog(`Unknown view type "` + tag + `"`)
		return nil
	}

	if !session.ignoreViewUpdates() {
		session.setIgnoreViewUpdates(true)
		defer session.setIgnoreViewUpdates(false)
	}

	view := creator(session)
	view.init(session)
	if customView, ok := view.(CustomView); ok {
		if !InitCustomView(customView, tag, session, nil) {
			return nil
		}
	}
	parseProperties(view, object)
	if binding != nil {
		view.setRaw(Binding, binding)
		if listener, ok := binding.(ViewCreateListener); ok {
			listener.OnCreate(view)
		}
	}
	return view
}

// CreateViewFromText create new View and initialize it by content of text. Parameters:
//   - session - the session to which the view will be attached (should not be nil);
//   - text - text describing View;
//   - binding - object assigned to the Binding property (optional parameter).
//
// If the function fails, it returns nil and an error message is written to the log.
func CreateViewFromText(session Session, text string, binding ...any) View {
	data, err := ParseDataText(text)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	var b any = nil
	if len(binding) > 0 {
		b = binding[0]
	}
	return CreateViewFromObject(session, data, b)
}

// CreateViewFromResources create new View and initialize it by the content of
// the resource file from "views" directory. Parameters:
//   - session - the session to which the view will be attached (should not be nil);
//   - name - file name in the "views" directory of the application resources (it is not necessary to specify the .rui extension, it is added automatically);
//   - binding - object assigned to the Binding property (optional parameter).
//
// If the function fails, it returns nil and an error message is written to the log.
func CreateViewFromResources(session Session, name string, binding ...any) View {
	if strings.ToLower(filepath.Ext(name)) != ".rui" {
		name += ".rui"
	}

	var b any = nil
	if len(binding) > 0 {
		b = binding[0]
	}

	createEmbed := func(fs *embed.FS, path string) View {
		if data, err := fs.ReadFile(path); err == nil {
			data, err := ParseDataText(string(data))
			if err == nil {
				return CreateViewFromObject(session, data, b)
			}
			ErrorLog(err.Error())
		}
		return nil
	}

	for _, fs := range resources.embedFS {
		rootDirs := resources.embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, rawDir:
				// do nothing

			case viewDir:
				if result := createEmbed(fs, dir+"/"+name); result != nil {
					return result
				}

			default:
				if result := createEmbed(fs, dir+"/"+viewDir+"/"+name); result != nil {
					return result
				}
			}
		}
	}

	if resources.path != "" {
		if data, err := os.ReadFile(resources.path + viewDir + "/" + name); err == nil {
			data, err := ParseDataText(string(data))
			if err != nil {
				ErrorLog(err.Error())
			} else {
				return CreateViewFromObject(session, data, b)
			}
		}
	}

	return nil
}
