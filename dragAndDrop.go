package rui

import (
	"encoding/base64"
	"fmt"
	"maps"
	"strings"
)

const (
	// DragData is the constant for "drag-data" property tag.
	//
	// Used by View:
	//
	// Supported types: map[string]string.
	DragData PropertyName = "drag-data"

	// DragImage is the constant for "drag-image" property tag.
	//
	// Used by View:
	// An url of image to use for the drag feedback image.
	//
	// Supported type: string.
	DragImage PropertyName = "drag-image"

	// DragImageXOffset is the constant for "drag-image-x-offset" property tag.
	//
	// Used by View:
	// The horizontal offset in pixels within the drag feedback image.
	//
	// Supported types: float, int, string.
	DragImageXOffset PropertyName = "drag-image-x-offset"

	// DragImageYOffset is the constant for "drag-image-y-offset" property tag.
	//
	// Used by View.
	// The vertical offset in pixels within the drag feedback image.
	//
	// Supported types: float, int, string.
	DragImageYOffset PropertyName = "drag-image-y-offset"

	// DropEffect is the constant for "drag-effect" property tag.
	//
	// Used by View.
	// Controls the feedback (typically visual) the user is given during a drag and drop operation.
	// It will affect which cursor is displayed while dragging. For example, when the user hovers over a target drop element,
	// the browser's cursor may indicate which type of operation will occur.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (DropEffectUndefined) or "undefined" - The property value is not defined (default value).
	//   - 1 (DropEffectCopy) or "copy" - A copy of the source item may be made at the new location.
	//   - 2 (DropEffectMove) or "move" - An item may be moved to a new location.
	//   - 4 (DropEffectLink) or "link" - A link may be established to the source at the new location.
	DropEffect PropertyName = "drag-effect"

	// DropEffectAllowed is the constant for "drop-effect-allowed" property tag.
	//
	// Used by View.
	// Specifies the effect that is allowed for a drag operation.
	// The copy operation is used to indicate that the data being dragged will be copied
	// from its present location to the drop location.
	// The move operation is used to indicate that the data being dragged will be moved,
	// and the link operation is used to indicate that some form of relationship
	// or connection will be created between the source and drop locations.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (DropEffectUndefined) or "undefined" - The property value is not defined (default value). Equivalent to DropEffectAll
	//   - 1 (DropEffectCopy) or "copy" - A copy of the source item may be made at the new location.
	//   - 2 (DropEffectMove) or "move" - An item may be moved to a new location.
	//   - 3 (DropEffectLink) or "link" - A link may be established to the source at the new location.
	//   - 4 (DropEffectCopyMove) or "copy|move" - A copy or move operation is permitted.
	//   - 5 (DropEffectCopyLink) or "copy|link" - A copy or link operation is permitted.
	//   - 6 (DropEffectLinkMove) or "link|move" - A link or move operation is permitted.
	//   - 7 (DropEffectAll) or "all" or "copy|move|link" - All operations are permitted.
	DropEffectAllowed PropertyName = "drag-effect-allowed"

	// DragStartEvent is the constant for "drag-start-event" property tag.
	//
	// Used by View.
	// Fired when the user starts dragging an element or text selection.
	//
	// General listener format:
	//  func(view rui.View, event rui.DragAndDropEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - event parameters.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(rui.DragAndDropEvent)
	//  func()
	DragStartEvent PropertyName = "drag-start-event"

	// DragEndEvent is the constant for "drag-end-event" property tag.
	//
	// Used by View.
	// Fired when a drag operation ends (by releasing a mouse button or hitting the escape key).
	//
	// General listener format:
	//  func(view rui.View, event rui.DragAndDropEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - event parameters.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(rui.DragAndDropEvent)
	//  func()
	DragEndEvent PropertyName = "drag-end-event"

	// DragEnterEvent is the constant for "drag-enter-event" property tag.
	//
	// Used by View.
	// Fired when a dragged element or text selection enters a valid drop target.
	//
	// General listener format:
	//  func(view rui.View, event rui.DragAndDropEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - event parameters.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(rui.DragAndDropEvent)
	//  func()
	DragEnterEvent PropertyName = "drag-enter-event"

	// DragLeaveEvent is the constant for "drag-leave-event" property tag.
	//
	// Used by View.
	// Fired when a dragged element or text selection leaves a valid drop target.
	//
	// General listener format:
	//  func(view rui.View, event rui.DragAndDropEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - event parameters.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(rui.DragAndDropEvent)
	//  func()
	DragLeaveEvent PropertyName = "drag-leave-event"

	// DragOverEvent is the constant for "drag-over-event" property tag.
	//
	// Used by View.
	// Fired when an element or text selection is being dragged over a valid drop target (every few hundred milliseconds).
	//
	// General listener format:
	//  func(view rui.View, event rui.DragAndDropEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - event parameters.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(rui.DragAndDropEvent)
	//  func()
	DragOverEvent PropertyName = "drag-over-event"

	// DropEvent is the constant for "drop-event" property tag.
	//
	// Used by View.
	// Fired when an element or text selection is dropped on a valid drop target.
	//
	// General listener format:
	//  func(view rui.View, event rui.DragAndDropEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - event parameters.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(rui.DragAndDropEvent)
	//  func()
	DropEvent PropertyName = "drop-event"

	// DropEffectUndefined - the value of the "drop-effect" and "drop-effect-allowed" properties: the value is not defined (default value).
	DropEffectUndefined = 0

	// DropEffectNone - the value of the DropEffect field of the DragEvent struct: the item may not be dropped.
	DropEffectNone = 0

	// DropEffectCopy - the value of the "drop-effect" and "drop-effect-allowed" properties: a copy of the source item may be made at the new location.
	DropEffectCopy = 1

	// DropEffectMove - the value of the "drop-effect" and "drop-effect-allowed" properties: an item may be moved to a new location.
	DropEffectMove = 2

	// DropEffectLink - the value of the "drop-effect" and "drop-effect-allowed" properties: a link may be established to the source at the new location.
	DropEffectLink = 4

	// DropEffectCopyMove - the value of the "drop-effect-allowed" property: a copy or move operation is permitted.
	DropEffectCopyMove = DropEffectCopy + DropEffectMove

	// DropEffectCopyLink - the value of the "drop-effect-allowed" property: a copy or link operation is permitted.
	DropEffectCopyLink = DropEffectCopy + DropEffectLink

	// DropEffectLinkMove - the value of the "drop-effect-allowed" property: a link or move operation is permitted.
	DropEffectLinkMove = DropEffectLink + DropEffectMove

	// DropEffectAll - the value of the "drop-effect-allowed" property: all operations (copy, move, and link) are permitted (default value).
	DropEffectAll = DropEffectCopy + DropEffectMove + DropEffectLink
)

// MouseEvent represent a mouse event
type DragAndDropEvent struct {
	MouseEvent
	Data          map[string]string
	Files         []FileInfo
	Target        View
	EffectAllowed int
	DropEffect    int
}

func (event *DragAndDropEvent) init(session Session, data DataObject) {
	event.MouseEvent.init(data)

	event.Data = map[string]string{}
	if value, ok := data.PropertyValue("data"); ok {
		data := strings.Split(value, ";")
		for _, line := range data {
			pair := strings.Split(line, ":")
			if len(pair) == 2 {
				mime, err := base64.StdEncoding.DecodeString(pair[0])
				if err != nil {
					ErrorLog(err.Error())
				} else {
					val, err := base64.StdEncoding.DecodeString(pair[1])
					if err == nil {
						event.Data[string(mime)] = string(val)
					} else {
						ErrorLog(err.Error())
					}
				}
			}
		}
	}

	if targetId, ok := data.PropertyValue("target"); ok {
		event.Target = session.viewByHTMLID(targetId)
	}

	if effect, ok := data.PropertyValue("effect-allowed"); ok {
		for i, value := range []string{"undefined", "copy", "move", "copyMove", "link", "copyLink", "linkMove", "all"} {
			if value == effect {
				event.EffectAllowed = i
				break
			}
		}
	}

	if effect, ok := data.PropertyValue("drop-effect"); ok && effect != "" {
		for i, value := range []string{"none", "copy", "move", "", "link"} {
			if value == effect {
				event.DropEffect = i
				break
			}
		}
	}

	event.Files = parseFilesTag(data)
}

func stringToDropEffect(text string) (int, bool) {
	text = strings.Trim(text, " \t\n")
	if n, ok := enumStringToInt(text, []string{"", "copy", "move", "", "link"}, false); ok {
		switch n {
		case DropEffectUndefined, DropEffectCopy, DropEffectMove, DropEffectLink:
			return n, true
		}
	}
	return 0, false
}

func (view *viewData) setDropEffect(value any) []PropertyName {
	if !setSimpleProperty(view, DropEffect, value) {
		if text, ok := value.(string); ok {

			if n, ok := stringToDropEffect(text); ok {
				if n == DropEffectUndefined {
					view.setRaw(DropEffect, nil)
				} else {
					view.setRaw(DropEffect, n)
				}
			} else {
				invalidPropertyValue(DropEffect, value)
				return nil
			}

		} else if i, ok := isInt(value); ok {

			switch i {
			case DropEffectUndefined:
				view.setRaw(DropEffect, nil)

			case DropEffectCopy, DropEffectMove, DropEffectLink:
				view.setRaw(DropEffect, i)

			default:
				invalidPropertyValue(DropEffect, value)
				return nil
			}

		} else {

			notCompatibleType(DropEffect, value)
			return nil
		}

	}

	return []PropertyName{DropEffect}
}

func stringToDropEffectAllowed(text string) (int, bool) {
	if strings.Contains(text, "|") {
		elements := strings.Split(text, "|")
		result := 0
		for _, element := range elements {
			if n, ok := stringToDropEffect(element); ok && n != DropEffectUndefined {
				result |= n
			} else {
				return 0, false
			}
		}
		return result, true
	}

	text = strings.Trim(text, " \t\n")
	if text != "" {
		if n, ok := enumStringToInt(text, []string{"undefined", "copy", "move", "", "link", "", "", "all"}, false); ok {
			return n, true
		}
	}
	return 0, false
}

func (view *viewData) setDropEffectAllowed(value any) []PropertyName {
	if !setSimpleProperty(view, DropEffectAllowed, value) {
		if text, ok := value.(string); ok {

			if n, ok := stringToDropEffectAllowed(text); ok {
				if n == DropEffectUndefined {
					view.setRaw(DropEffectAllowed, nil)
				} else {
					view.setRaw(DropEffectAllowed, n)
				}
			} else {
				invalidPropertyValue(DropEffectAllowed, value)
				return nil
			}

		} else {
			n, ok := isInt(value)
			if !ok {
				notCompatibleType(DropEffectAllowed, value)
				return nil
			}

			if n == DropEffectUndefined {
				view.setRaw(DropEffectAllowed, nil)
			} else if n > DropEffectUndefined && n <= DropEffectAll {
				view.setRaw(DropEffectAllowed, n)
			} else {
				notCompatibleType(DropEffectAllowed, value)
				return nil
			}
		}
	}

	return []PropertyName{DropEffectAllowed}
}

func handleDragAndDropEvents(view View, tag PropertyName, data DataObject) {
	listeners := getOneArgEventListeners[View, DragAndDropEvent](view, nil, tag)
	if len(listeners) > 0 {
		var event DragAndDropEvent
		event.init(view.Session(), data)

		for _, listener := range listeners {
			listener.Run(view, event)
		}
	}
}

func base64DragData(view View) string {
	if value := view.getRaw(DragData); value != nil {
		if data, ok := value.(map[string]string); ok && len(data) > 0 {
			buf := allocStringBuilder()
			defer freeStringBuilder(buf)

			for mime, value := range data {
				if buf.Len() > 0 {
					buf.WriteRune(';')
				}
				buf.WriteString(base64.StdEncoding.EncodeToString([]byte(mime)))
				buf.WriteRune(':')
				buf.WriteString(base64.StdEncoding.EncodeToString([]byte(value)))
			}

			return buf.String()
		}
	}
	return ""
}

func dragAndDropHtml(view View, buffer *strings.Builder) {

	if len(getOneArgEventListeners[View, DragAndDropEvent](view, nil, DropEvent)) > 0 {
		buffer.WriteString(`ondragover="dragOverEvent(this, event)" ondrop="dropEvent(this, event)" `)
		if len(getOneArgEventListeners[View, DragAndDropEvent](view, nil, DragOverEvent)) > 0 {
			buffer.WriteString(`data-drag-over="1" `)
		}
	}

	if dragData := base64DragData(view); dragData != "" {
		buffer.WriteString(`draggable="true" data-drag="`)
		buffer.WriteString(dragData)
		buffer.WriteString(`" ondragstart="dragStartEvent(this, event)" `)
	} else if len(getOneArgEventListeners[View, DragAndDropEvent](view, nil, DragStartEvent)) > 0 {
		buffer.WriteString(` ondragstart="dragStartEvent(this, event)" `)
	}

	enterEvent := false
	switch GetDropEffect(view) {
	case DropEffectCopy:
		buffer.WriteString(` data-drop-effect="copy" ondragenter="dragEnterEvent(this, event)"`)
		enterEvent = true

	case DropEffectMove:
		buffer.WriteString(` data-drop-effect="move" ondragenter="dragEnterEvent(this, event)"`)
		enterEvent = true

	case DropEffectLink:
		buffer.WriteString(` data-drop-effect="link" ondragenter="dragEnterEvent(this, event)"`)
		enterEvent = true
	}

	if enterEvent {
		viewEventsHtml[DragAndDropEvent](view, []PropertyName{DragEndEvent, DragLeaveEvent}, buffer)
	} else {
		viewEventsHtml[DragAndDropEvent](view, []PropertyName{DragEndEvent, DragEnterEvent, DragLeaveEvent}, buffer)
	}

	if img := GetDragImage(view); img != "" {
		buffer.WriteString(` data-drag-image="`)
		buffer.WriteString(img)
		buffer.WriteString(`" `)
	}

	if f := GetDragImageXOffset(view); f != 0 {
		buffer.WriteString(` data-drag-image-x="`)
		fmt.Fprintf(buffer, "%g", f)
		buffer.WriteString(`" `)
	}

	if f := GetDragImageYOffset(view); f != 0 {
		buffer.WriteString(` data-drag-image-y="`)
		fmt.Fprintf(buffer, "%g", f)
		buffer.WriteString(`" `)
	}

	effects := []string{"undefined", "copy", "move", "copyMove", "link", "copyLink", "linkMove", "all"}
	if n := GetDropEffectAllowed(view); n > 0 && n < len(effects) {
		buffer.WriteString(` data-drop-effect-allowed="`)
		buffer.WriteString(effects[n])
		buffer.WriteString(`" `)
	}
}

func (view *viewData) LoadFile(file FileInfo, result func(FileInfo, []byte)) {
	if result != nil {
		view.fileLoader[file.key()] = result
		view.Session().callFunc("loadDropFile", view.htmlID(), file.Name, file.Size)
	}
}

// GetDragStartEventListeners returns the "drag-start-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.DragAndDropEvent),
//   - func(rui.View),
//   - func(rui.DragAndDropEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragStartEventListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, DragAndDropEvent](view, subviewID, DragStartEvent)
}

// GetDragEndEventListeners returns the "drag-end-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.DragAndDropEvent),
//   - func(rui.View),
//   - func(rui.DragAndDropEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragEndEventListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, DragAndDropEvent](view, subviewID, DragEndEvent)
}

// GetDragEnterEventListeners returns the "drag-enter-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.DragAndDropEvent),
//   - func(rui.View),
//   - func(rui.DragAndDropEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragEnterEventListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, DragAndDropEvent](view, subviewID, DragEnterEvent)
}

// GetDragLeaveEventListeners returns the "drag-leave-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.DragAndDropEvent),
//   - func(rui.View),
//   - func(rui.DragAndDropEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragLeaveEventListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, DragAndDropEvent](view, subviewID, DragLeaveEvent)
}

// GetDragOverEventListeners returns the "drag-over-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.DragAndDropEvent),
//   - func(rui.View),
//   - func(rui.DragAndDropEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragOverEventListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, DragAndDropEvent](view, subviewID, DragOverEvent)
}

// GetDropEventListeners returns the "drag-start-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.DragAndDropEvent),
//   - func(rui.View),
//   - func(rui.DragAndDropEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDropEventListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, DragAndDropEvent](view, subviewID, DropEvent)
}

// GetDropEventListeners returns the "drag-data" data.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragData(view View, subviewID ...string) map[string]string {
	result := map[string]string{}
	if view = getSubview(view, subviewID); view != nil {
		if value := view.getRaw(DragData); value != nil {
			if data, ok := value.(map[string]string); ok {
				maps.Copy(result, data)
			}
		}
	}

	return result
}

// GetDragImage returns the drag feedback image.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragImage(view View, subviewID ...string) string {
	if view = getSubview(view, subviewID); view != nil {
		value := view.getRaw(DragImage)
		if value == nil {
			value = valueFromStyle(view, DragImage)
		}

		if value != nil {
			if img, ok := value.(string); ok {
				img = strings.Trim(img, " \t")
				if img != "" && img[0] == '@' {
					if img, ok = view.Session().ImageConstant(img[1:]); ok {
						return img
					}
				} else {
					return img
				}
			}
		}
	}
	return ""
}

// GetDragImageXOffset returns the horizontal offset in pixels within the drag feedback image.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragImageXOffset(view View, subviewID ...string) float64 {
	return floatStyledProperty(view, subviewID, DragImageXOffset, 0)
}

// GetDragImageYOffset returns the vertical offset in pixels within the drag feedback image.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDragImageYOffset(view View, subviewID ...string) float64 {
	return floatStyledProperty(view, subviewID, DragImageYOffset, 0)
}

// GetDropEffect returns the effect that is allowed for a drag operation.
// Controls the feedback (typically visual) the user is given during a drag and drop operation.
// It will affect which cursor is displayed while dragging.
//
// Returns one of next values:
//   - 0 (DropEffectUndefined) - The value is not defined (all operations are permitted).
//   - 1 (DropEffectCopy) - A copy of the source item may be made at the new location.
//   - 2 (DropEffectMove) - An item may be moved to a new location.
//   - 4 (DropEffectLink) - A link may be established to the source at the new location.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDropEffect(view View, subviewID ...string) int {
	if view = getSubview(view, subviewID); view != nil {
		value := view.getRaw(DropEffect)
		if value == nil {
			value = valueFromStyle(view, DropEffect)
		}

		if value != nil {
			switch value := value.(type) {
			case int:
				return value

			case string:
				if value, ok := view.Session().resolveConstants(value); ok {
					if n, ok := stringToDropEffect(value); ok {
						return n
					}
				}

			default:
				return DropEffectUndefined
			}
		}
	}
	return DropEffectUndefined
}

// GetDropEffectAllowed returns the effect that is allowed for a drag operation.
// The copy operation is used to indicate that the data being dragged will be copied from its present location to the drop location.
// The move operation is used to indicate that the data being dragged will be moved,
// and the link operation is used to indicate that some form of relationship
// or connection will be created between the source and drop locations.
//
// Returns one of next values:
//   - 0 (DropEffectUndefined) - The value is not defined (all operations are permitted).
//   - 1 (DropEffectCopy) - A copy of the source item may be made at the new location.
//   - 2 (DropEffectMove) - An item may be moved to a new location.
//   - 4 (DropEffectLink) - A link may be established to the source at the new location.
//   - 3 (DropEffectCopyMove) - A copy or move operation is permitted.
//   - 5 (DropEffectCopyLink) - A copy or link operation is permitted.
//   - 6 (DropEffectLinkMove) - A link or move operation is permitted.
//   - 7 (DropEffectAll) - All operations are permitted.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDropEffectAllowed(view View, subviewID ...string) int {
	if view = getSubview(view, subviewID); view != nil {
		value := view.getRaw(DropEffectAllowed)
		if value == nil {
			value = valueFromStyle(view, DropEffectAllowed)
		}

		if value != nil {
			switch value := value.(type) {
			case int:
				return value

			case string:
				if value, ok := view.Session().resolveConstants(value); ok {
					if n, ok := stringToDropEffectAllowed(value); ok {
						return n
					}
				}

			default:
				return DropEffectUndefined
			}
		}
	}
	return DropEffectUndefined
}
