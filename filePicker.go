package rui

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// FileSelectedEvent is the constant for "file-selected-event" property tag.
	// The "file-selected-event" is fired when user selects file(s) in the FilePicker.
	FileSelectedEvent = "file-selected-event"
	// Accept is the constant for "accept" property tag.
	// The "accept" property of the FilePicker sets the list of allowed file extensions or MIME types.
	Accept = "accept"
	// Multiple is the constant for "multiple" property tag.
	// The "multiple" bool property of the FilePicker sets whether multiple files can be selected
	Multiple = "multiple"
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
}

// FilePicker - the control view for the files selecting
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
	files                 []FileInfo
	fileSelectedListeners []func(FilePicker, []FileInfo)
	loader                map[int]func(FileInfo, []byte)
}

func (file *FileInfo) initBy(node DataValue) {
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
}

// NewFilePicker create new FilePicker object and return it
func NewFilePicker(session Session, params Params) FilePicker {
	view := new(filePickerData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newFilePicker(session Session) View {
	return NewFilePicker(session, nil)
}

func (picker *filePickerData) Init(session Session) {
	picker.viewData.Init(session)
	picker.tag = "FilePicker"
	picker.files = []FileInfo{}
	picker.loader = map[int]func(FileInfo, []byte){}
	picker.fileSelectedListeners = []func(FilePicker, []FileInfo){}
}

func (picker *filePickerData) String() string {
	return getViewString(picker)
}

func (picker *filePickerData) Focusable() bool {
	return true
}

func (picker *filePickerData) Files() []FileInfo {
	return picker.files
}

func (picker *filePickerData) LoadFile(file FileInfo, result func(FileInfo, []byte)) {
	if result == nil {
		return
	}

	for i, info := range picker.files {
		if info.Name == file.Name && info.Size == file.Size && info.LastModified == file.LastModified {
			picker.loader[i] = result
			picker.Session().runScript(fmt.Sprintf(`loadSelectedFile("%s", %d)`, picker.htmlID(), i))
			return
		}
	}
}

func (picker *filePickerData) Remove(tag string) {
	picker.remove(strings.ToLower(tag))
}

func (picker *filePickerData) remove(tag string) {
	switch tag {
	case FileSelectedEvent:
		if len(picker.fileSelectedListeners) > 0 {
			picker.fileSelectedListeners = []func(FilePicker, []FileInfo){}
			picker.propertyChangedEvent(tag)
		}

	case Accept:
		delete(picker.properties, tag)
		if picker.created {
			removeProperty(picker.htmlID(), "accept", picker.Session())
		}
		picker.propertyChangedEvent(tag)

	default:
		picker.viewData.remove(tag)
	}
}

func (picker *filePickerData) Set(tag string, value any) bool {
	return picker.set(strings.ToLower(tag), value)
}

func (picker *filePickerData) set(tag string, value any) bool {
	if value == nil {
		picker.remove(tag)
		return true
	}

	switch tag {
	case FileSelectedEvent:
		listeners, ok := valueToEventListeners[FilePicker, []FileInfo](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(FilePicker, []FileInfo){}
		}
		picker.fileSelectedListeners = listeners
		picker.propertyChangedEvent(tag)
		return true

	case Accept:
		switch value := value.(type) {
		case string:
			value = strings.Trim(value, " \t\n")
			if value == "" {
				picker.remove(Accept)
			} else {
				picker.properties[Accept] = value
			}

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
			if buffer.Len() == 0 {
				picker.remove(Accept)
			} else {
				picker.properties[Accept] = buffer.String()
			}

		default:
			notCompatibleType(tag, value)
			return false
		}

		if picker.created {
			if css := picker.acceptCSS(); css != "" {
				updateProperty(picker.htmlID(), "accept", css, picker.Session())
			} else {
				removeProperty(picker.htmlID(), "accept", picker.Session())
			}
		}
		picker.propertyChangedEvent(tag)
		return true

	default:
		return picker.viewData.set(tag, value)
	}
}

func (picker *filePickerData) htmlTag() string {
	return "input"
}

func (picker *filePickerData) acceptCSS() string {
	accept, ok := stringProperty(picker, Accept, picker.Session())
	if !ok {
		if value := valueFromStyle(picker, Accept); value != nil {
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

	if accept := picker.acceptCSS(); accept != "" {
		buffer.WriteString(` accept="`)
		buffer.WriteString(accept)
		buffer.WriteRune('"')
	}

	buffer.WriteString(` type="file"`)
	if IsMultipleFilePicker(picker, "") {
		buffer.WriteString(` multiple`)
	}

	buffer.WriteString(` oninput="fileSelectedEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}
}

func (picker *filePickerData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	if IsDisabled(self, "") {
		buffer.WriteString(` disabled`)
	}
	picker.viewData.htmlDisabledProperties(self, buffer)
}

func (picker *filePickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "fileSelected":
		if node := data.PropertyWithTag("files"); node != nil && node.Type() == ArrayNode {
			count := node.ArraySize()
			files := make([]FileInfo, count)
			for i := 0; i < count; i++ {
				if value := node.ArrayElement(i); value != nil {
					files[i].initBy(value)
				}
			}
			picker.files = files

			for _, listener := range picker.fileSelectedListeners {
				listener(picker, files)
			}
		}
		return true

	case "fileLoaded":
		if index, ok := dataIntProperty(data, "index"); ok {
			if result, ok := picker.loader[index]; ok {
				var file FileInfo
				file.initBy(data)

				var fileData []byte = nil
				if base64Data, ok := data.PropertyValue("data"); ok {
					if index := strings.LastIndex(base64Data, ","); index >= 0 {
						base64Data = base64Data[index+1:]
					}
					decode, err := base64.StdEncoding.DecodeString(base64Data)
					if err == nil {
						fileData = decode
					} else {
						ErrorLog(err.Error())
					}
				}

				result(file, fileData)
				delete(picker.loader, index)
			}
		}
		return true

	case "fileLoadingError":
		if error, ok := data.PropertyValue("error"); ok {
			ErrorLog(error)
		}
		if index, ok := dataIntProperty(data, "index"); ok {
			if result, ok := picker.loader[index]; ok {
				if index >= 0 && index < len(picker.files) {
					result(picker.files[index], nil)
				} else {
					result(FileInfo{}, nil)
				}
				delete(picker.loader, index)
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

// GetFilePickerFiles returns the list of FilePicker selected files
// If there are no files selected then an empty slice is returned (the result is always not nil)
// If the second argument (subviewID) is "" then selected files of the first argument (view) is returned
func GetFilePickerFiles(view View, subviewID string) []FileInfo {
	if picker := FilePickerByID(view, subviewID); picker != nil {
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
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func IsMultipleFilePicker(view View, subviewID string) bool {
	return boolStyledProperty(view, subviewID, Multiple, false)
}

// GetFilePickerAccept returns sets the list of allowed file extensions or MIME types.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetFilePickerAccept(view View, subviewID string) []string {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		accept, ok := stringProperty(view, Accept, view.Session())
		if !ok {
			if value := valueFromStyle(view, Accept); value != nil {
				accept, ok = value.(string)
			}
		}
		if ok {
			result := strings.Split(accept, ",")
			for i := 0; i < len(result); i++ {
				result[i] = strings.Trim(result[i], " \t\n")
			}
			return result
		}
	}
	return []string{}
}

// GetFileSelectedListeners returns the "file-selected-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetFileSelectedListeners(view View, subviewID string) []func(FilePicker, []FileInfo) {
	return getEventListeners[FilePicker, []FileInfo](view, subviewID, FileSelectedEvent)
}
