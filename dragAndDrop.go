package rui

import (
	"encoding/base64"
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

	// DragStartEvent is the constant for "drag-start-event" property tag.
	//
	// Used by View.
	// Fired when the user starts dragging an element or text selection.
	//
	// General listener format:
	//
	DragStartEvent PropertyName = "drag-start-event"

	// DragEndEvent is the constant for "drag-end-event" property tag.
	//
	// Used by View.
	// Fired when a drag operation ends (by releasing a mouse button or hitting the escape key).
	//
	// General listener format:
	//
	DragEndEvent PropertyName = "drag-end-event"

	// DragEnterEvent is the constant for "drag-enter-event" property tag.
	//
	// Used by View.
	// Fired when a dragged element or text selection enters a valid drop target.
	//
	// General listener format:
	//
	DragEnterEvent PropertyName = "drag-enter-event"

	// DragLeaveEvent is the constant for "drag-leave-event" property tag.
	//
	// Used by View.
	// Fired when a dragged element or text selection leaves a valid drop target.
	//
	// General listener format:
	//
	DragLeaveEvent PropertyName = "drag-leave-event"

	// DragOverEvent is the constant for "drag-over-event" property tag.
	//
	// Used by View.
	// Fired when an element or text selection is being dragged over a valid drop target (every few hundred milliseconds).
	//
	// General listener format:
	//
	DragOverEvent PropertyName = "drag-over-event"

	// DropEvent is the constant for "drop-event" property tag.
	//
	// Used by View.
	// Fired when an element or text selection is dropped on a valid drop target.
	//
	// General listener format:
	//
	DropEvent PropertyName = "drop-event"
)

// MouseEvent represent a mouse event
type DragAndDropEvent struct {
	MouseEvent
	Data   map[string]string
	Target View
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
}

func handleDragAndDropEvents(view View, tag PropertyName, data DataObject) {
	listeners := getOneArgEventListeners[View, DragAndDropEvent](view, nil, tag)
	if len(listeners) > 0 {
		var event DragAndDropEvent
		event.init(view.Session(), data)

		for _, listener := range listeners {
			listener(view, event)
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

	viewEventsHtml[DragAndDropEvent](view, []PropertyName{DragEndEvent, DragEnterEvent, DragLeaveEvent}, buffer)

	session := view.Session()
	if img, ok := stringProperty(view, DragImage, session); ok && img != "" {
		img = strings.Trim(img, " \t")
		if img[0] == '@' {
			if img, ok = session.ImageConstant(img[1:]); ok {
				buffer.WriteString(` data-drag-image="`)
				buffer.WriteString(img)
				buffer.WriteString(`" `)
			}
		} else {
			buffer.WriteString(` data-drag-image="`)
			buffer.WriteString(img)
			buffer.WriteString(`" `)
		}
	}

	if f, ok := floatTextProperty(view, DragImageXOffset, session, 0); ok {
		buffer.WriteString(` data-drag-image-x="`)
		buffer.WriteString(f)
		buffer.WriteString(`" `)
	}

	if f, ok := floatTextProperty(view, DragImageYOffset, session, 0); ok {
		buffer.WriteString(` data-drag-image-y="`)
		buffer.WriteString(f)
		buffer.WriteString(`" `)
	}
}

// GetDragStartEventListeners returns the "drag-start-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDragStartEventListeners(view View, subviewID ...string) []func(View, DragAndDropEvent) {
	return getOneArgEventListeners[View, DragAndDropEvent](view, subviewID, DragStartEvent)
}

// GetDragEndEventListeners returns the "drag-end-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDragEndEventListeners(view View, subviewID ...string) []func(View, DragAndDropEvent) {
	return getOneArgEventListeners[View, DragAndDropEvent](view, subviewID, DragEndEvent)
}

// GetDragEnterEventListeners returns the "drag-enter-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDragEnterEventListeners(view View, subviewID ...string) []func(View, DragAndDropEvent) {
	return getOneArgEventListeners[View, DragAndDropEvent](view, subviewID, DragEnterEvent)
}

// GetDragLeaveEventListeners returns the "drag-leave-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDragLeaveEventListeners(view View, subviewID ...string) []func(View, DragAndDropEvent) {
	return getOneArgEventListeners[View, DragAndDropEvent](view, subviewID, DragLeaveEvent)
}

// GetDragOverEventListeners returns the "drag-over-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDragOverEventListeners(view View, subviewID ...string) []func(View, DragAndDropEvent) {
	return getOneArgEventListeners[View, DragAndDropEvent](view, subviewID, DragOverEvent)
}

// GetDropEventListeners returns the "drag-start-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropEventListeners(view View, subviewID ...string) []func(View, DragAndDropEvent) {
	return getOneArgEventListeners[View, DragAndDropEvent](view, subviewID, DropEvent)
}

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
