package rui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Constants for [FilePicker] specific properties and events
const (
	// FileSelectedEvent is the constant for "file-selected-event" property tag.
	//
	// Used by FilePicker.
	// Fired when user selects file(s).
	//
	// General listener format:
	//  func(picker rui.FilePicker, files []rui.FileInfo).
	//
	// where:
	// picker - Interface of a file picker which generated this event,
	// files - Array of description of selected files.
	//
	// Allowed listener formats:
	//  func(picker rui.FilePicker)
	//  func(files []rui.FileInfo)
	//  func()
	FileSelectedEvent PropertyName = "file-selected-event"

	// Accept is the constant for "accept" property tag.
	//
	// Used by FilePicker.
	// Set the list of allowed file extensions or MIME types.
	//
	// Supported types: string, []string.
	//
	// Internal type is string, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - string - may contain single value of multiple separated by comma(,).
	//   - []string - an array of acceptable file extensions or MIME types.
	Accept PropertyName = "accept"

	// Multiple is the constant for "multiple" property tag.
	//
	// Used by FilePicker.
	// Controls whether multiple files can be selected.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - Several files can be selected.
	//   - false, 0, "false", "no", "off", "0" - Only one file can be selected.
	Multiple PropertyName = "multiple"
)

// FileInfo describes a file which selected in the FilePicker view
type FileInfo struct {
	// Name - the file's name.
	Name string

	// LastModified specifying the date and time at which the file was last modified
	LastModified time.Time

	// Size - the size of the file in bytes.
	Size int64

	// MimeType - the file's MIME type.
	MimeType string

	data []byte
}

// FilePicker represents the FilePicker view
type FilePicker interface {
	View
	// Files returns the list of selected files.
	// If there are no files selected then an empty slice is returned (the result is always not nil)
	Files() []FileInfo
	// LoadFile loads the content of the selected file. This function is asynchronous.
	// The "result" function will be called after loading the data.
	LoadFile(file FileInfo, result func(FileInfo, []byte))
}

type filePickerData struct {
	viewData
	files []FileInfo
	//loader map[int]func(FileInfo, []byte)
}

func dataToFileInfo(node DataValue) FileInfo {
	var file FileInfo
	if obj := node.Object(); obj != nil {
		file.Name, _ = obj.PropertyValue("name")
		file.MimeType, _ = obj.PropertyValue("mime-type")

		if size, ok := obj.PropertyValue("size"); ok {
			if n, err := strconv.ParseInt(size, 10, 64); err == nil {
				file.Size = n
			}
		}

		if value, ok := obj.PropertyValue("last-modified"); ok {
			if n, err := strconv.ParseInt(value, 10, 64); err == nil {
				file.LastModified = time.UnixMilli(n)
			}
		}
	}
	return file
}

func (file FileInfo) key() string {
	return fmt.Sprintf("%s:%d", file.Name, int(file.Size))
}

// NewFilePicker create new FilePicker object and return it
func NewFilePicker(session Session, params Params) FilePicker {
	view := new(filePickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newFilePicker(session Session) View {
	return new(filePickerData) // NewFilePicker(session, nil)
}

func (picker *filePickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "FilePicker"
	picker.hasHtmlDisabled = true
	picker.files = []FileInfo{}
	//picker.loader = map[int]func(FileInfo, []byte){}
	picker.set = picker.setFunc
	picker.changed = picker.propertyChanged

}

func (picker *filePickerData) Focusable() bool {
	return true
}

func (picker *filePickerData) Files() []FileInfo {
	return picker.files
}

func (picker *filePickerData) LoadFile(file FileInfo, result func(FileInfo, []byte)) {
	if result != nil {
		for i, info := range picker.files {
			if info.Name == file.Name && info.Size == file.Size && info.LastModified == file.LastModified {
				if info.data != nil {
					result(info, info.data)
				} else {
					picker.fileLoader[info.key()] = func(file FileInfo, data []byte) {
						picker.files[i].data = data
						result(file, data)
					}
					picker.Session().callFunc("loadSelectedFile", picker.htmlID(), i)
				}
				return
			}
		}

		picker.viewData.LoadFile(file, result)
	}
}

func (picker *filePickerData) setFunc(tag PropertyName, value any) []PropertyName {

	switch tag {
	case FileSelectedEvent:
		return setOneArgEventListener[FilePicker, []FileInfo](picker, tag, value)

	case Accept:
		switch value := value.(type) {
		case string:
			return setStringPropertyValue(picker, Accept, strings.Trim(value, " \t\n"))

		case []string:
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)
			for _, val := range value {
				val = strings.Trim(val, " \t\n")
				if val != "" {
					if buffer.Len() > 0 {
						buffer.WriteRune(',')
					}
					buffer.WriteString(val)
				}
			}
			return setStringPropertyValue(picker, Accept, buffer.String())
		}
		notCompatibleType(tag, value)
		return nil
	}

	return picker.viewData.setFunc(tag, value)
}

func (picker *filePickerData) propertyChanged(tag PropertyName) {
	switch tag {
	case Accept:
		session := picker.Session()
		if css := acceptPropertyCSS(picker); css != "" {
			session.updateProperty(picker.htmlID(), "accept", css)
		} else {
			session.removeProperty(picker.htmlID(), "accept")
		}

	default:
		picker.viewData.propertyChanged(tag)
	}
}

func (picker *filePickerData) htmlTag() string {
	return "input"
}

func acceptPropertyCSS(view View) string {
	accept, ok := stringProperty(view, Accept, view.Session())
	if !ok {
		if value := valueFromStyle(view, Accept); value != nil {
			accept, ok = value.(string)
		}
	}

	if ok {
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)
		for _, value := range strings.Split(accept, ",") {
			if value = strings.Trim(value, " \t\n"); value != "" {
				if buffer.Len() > 0 {
					buffer.WriteString(", ")
				}
				if value[0] != '.' && !strings.Contains(value, "/") {
					buffer.WriteRune('.')
				}
				buffer.WriteString(value)
			}
		}
		return buffer.String()
	}
	return ""
}

func (picker *filePickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	if accept := acceptPropertyCSS(picker); accept != "" {
		buffer.WriteString(` accept="`)
		buffer.WriteString(accept)
		buffer.WriteRune('"')
	}

	buffer.WriteString(` type="file"`)
	if IsMultipleFilePicker(picker) {
		buffer.WriteString(` multiple`)
	}

	buffer.WriteString(` oninput="fileSelectedEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}
}

func parseFilesTag(data DataObject) []FileInfo {
	if node := data.PropertyByTag("files"); node != nil && node.Type() == ArrayNode {
		count := node.ArraySize()
		files := make([]FileInfo, count)
		for i := range count {
			if value := node.ArrayElement(i); value != nil {
				files[i] = dataToFileInfo(value)
			}
		}
		return files
	}
	return nil
}

func (picker *filePickerData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "fileSelected":
		if files := parseFilesTag(data); files != nil {
			picker.files = files
			for _, listener := range GetFileSelectedListeners(picker) {
				listener(picker, files)
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

// GetFilePickerFiles returns the list of FilePicker selected files
// If there are no files selected then an empty slice is returned (the result is always not nil)
// If the second argument (subviewID) is not specified or it is "" then selected files of the first argument (view) is returned
func GetFilePickerFiles(view View, subviewID ...string) []FileInfo {
	subview := ""
	if len(subviewID) > 0 {
		subview = subviewID[0]
	}

	if picker := FilePickerByID(view, subview); picker != nil {
		return picker.Files()
	}
	return []FileInfo{}
}

// LoadFilePickerFile loads the content of the selected file. This function is asynchronous.
// The "result" function will be called after loading the data.
// If the second argument (subviewID) is "" then the file from the first argument (view) is loaded
func LoadFilePickerFile(view View, subviewID string, file FileInfo, result func(FileInfo, []byte)) {
	if picker := FilePickerByID(view, subviewID); picker != nil {
		picker.LoadFile(file, result)
	}
}

// IsMultipleFilePicker returns "true" if multiple files can be selected in the FilePicker, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsMultipleFilePicker(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Multiple, false)
}

// GetFilePickerAccept returns sets the list of allowed file extensions or MIME types.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetFilePickerAccept(view View, subviewID ...string) []string {
	if view = getSubview(view, subviewID); view != nil {
		accept, ok := stringProperty(view, Accept, view.Session())
		if !ok {
			if value := valueFromStyle(view, Accept); value != nil {
				accept, ok = value.(string)
			}
		}
		if ok {
			result := strings.Split(accept, ",")
			for i := range len(result) {
				result[i] = strings.Trim(result[i], " \t\n")
			}
			return result
		}
	}
	return []string{}
}

// GetFileSelectedListeners returns the "file-selected-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetFileSelectedListeners(view View, subviewID ...string) []func(FilePicker, []FileInfo) {
	return getOneArgEventListeners[FilePicker, []FileInfo](view, subviewID, FileSelectedEvent)
}
