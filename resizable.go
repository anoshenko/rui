package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Side is the constant for the "side" property tag.
	// The "side" int property determines which side of the container is used to resize.
	// The value of property is or-combination of TopSide (1), RightSide (2), BottomSide (4), and LeftSide (8)
	Side = "side"
	// ResizeBorderWidth is the constant for the "resize-border-width" property tag.
	// The "ResizeBorderWidth" SizeUnit property determines the width of the resizing border
	ResizeBorderWidth = "resize-border-width"
	// CellVerticalAlign is the constant for the "cell-vertical-align" property tag.
	CellVerticalAlign = "cell-vertical-align"
	// CellHorizontalAlign is the constant for the "cell-horizontal-align" property tag.
	CellHorizontalAlign = "cell-horizontal-align"

	// TopSide is value of the "side" property: the top side is used to resize
	TopSide = 1
	// RightSide is value of the "side" property: the right side is used to resize
	RightSide = 2
	// BottomSide is value of the "side" property: the bottom side is used to resize
	BottomSide = 4
	// LeftSide is value of the "side" property: the left side is used to resize
	LeftSide = 8
	// AllSides is value of the "side" property: all sides is used to resize
	AllSides = TopSide | RightSide | BottomSide | LeftSide
)

// Resizable - grid-container of View
type Resizable interface {
	View
	ParentView
}

type resizableData struct {
	viewData
	content []View
}

// NewResizable create new Resizable object and return it
func NewResizable(session Session, params Params) Resizable {
	view := new(resizableData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newResizable(session Session) View {
	return NewResizable(session, nil)
}

func (resizable *resizableData) init(session Session) {
	resizable.viewData.init(session)
	resizable.tag = "Resizable"
	resizable.systemClass = "ruiGridLayout"
	resizable.content = []View{}
}

func (resizable *resizableData) String() string {
	return getViewString(resizable)
}

func (resizable *resizableData) Views() []View {
	return resizable.content
}

func (resizable *resizableData) Remove(tag string) {
	resizable.remove(strings.ToLower(tag))
}

func (resizable *resizableData) remove(tag string) {
	switch tag {
	case Side:
		oldSide := resizable.getSide()
		delete(resizable.properties, Side)
		if oldSide != resizable.getSide() {
			if resizable.created {
				updateInnerHTML(resizable.htmlID(), resizable.Session())
				resizable.updateResizeBorderWidth()
			}
			resizable.propertyChangedEvent(tag)
		}

	case ResizeBorderWidth:
		w := resizable.resizeBorderWidth()
		delete(resizable.properties, ResizeBorderWidth)
		if !w.Equal(resizable.resizeBorderWidth()) {
			resizable.updateResizeBorderWidth()
			resizable.propertyChangedEvent(tag)
		}

	case Content:
		if len(resizable.content) > 0 {
			resizable.content = []View{}
			if resizable.created {
				updateInnerHTML(resizable.htmlID(), resizable.Session())
			}
			resizable.propertyChangedEvent(tag)
		}

	default:
		resizable.viewData.remove(tag)
	}
}

func (resizable *resizableData) Set(tag string, value any) bool {
	return resizable.set(strings.ToLower(tag), value)
}

func (resizable *resizableData) set(tag string, value any) bool {
	if value == nil {
		resizable.remove(tag)
		return true
	}

	switch tag {
	case Side:
		oldSide := resizable.getSide()
		if !resizable.setSide(value) {
			notCompatibleType(tag, value)
			return false
		}
		if oldSide != resizable.getSide() {
			if resizable.created {
				updateInnerHTML(resizable.htmlID(), resizable.Session())
				resizable.updateResizeBorderWidth()
			}
			resizable.propertyChangedEvent(tag)
		}
		return true

	case ResizeBorderWidth:
		w := resizable.resizeBorderWidth()
		ok := resizable.setSizeProperty(tag, value)
		if ok && !w.Equal(resizable.resizeBorderWidth()) {
			resizable.updateResizeBorderWidth()
			resizable.propertyChangedEvent(tag)
		}
		return ok

	case Content:
		var newContent View = nil
		switch value := value.(type) {
		case string:
			newContent = NewTextView(resizable.Session(), Params{Text: value})

		case View:
			newContent = value

		case DataObject:
			if view := CreateViewFromObject(resizable.Session(), value); view != nil {
				newContent = view
			} else {
				return false
			}

		default:
			notCompatibleType(tag, value)
			return false
		}

		if len(resizable.content) == 0 {
			resizable.content = []View{newContent}
		} else {
			resizable.content[0] = newContent
		}
		if resizable.created {
			updateInnerHTML(resizable.htmlID(), resizable.Session())
		}
		resizable.propertyChangedEvent(tag)
		return true

	case CellWidth, CellHeight, GridRowGap, GridColumnGap, CellVerticalAlign, CellHorizontalAlign:
		ErrorLogF(`Not supported "%s" property`, tag)
		return false
	}

	return resizable.viewData.set(tag, value)
}

func (resizable *resizableData) Get(tag string) any {
	return resizable.get(strings.ToLower(tag))
}

func (resizable *resizableData) getSide() int {
	if value := resizable.getRaw(Side); value != nil {
		switch value := value.(type) {
		case string:
			if value, ok := resizable.session.resolveConstants(value); ok {
				validValues := map[string]int{
					"top":    TopSide,
					"right":  RightSide,
					"bottom": BottomSide,
					"left":   LeftSide,
					"all":    AllSides,
				}

				if strings.Contains(value, "|") {
					values := strings.Split(value, "|")
					sides := 0
					for _, val := range values {
						if n, err := strconv.Atoi(val); err == nil {
							if n < 1 || n > AllSides {
								return AllSides
							}
							sides |= n
						} else if n, ok := validValues[val]; ok {
							sides |= n
						} else {
							return AllSides
						}
					}
					return sides

				} else if n, err := strconv.Atoi(value); err == nil {
					if n >= 1 || n <= AllSides {
						return n
					}
				} else if n, ok := validValues[value]; ok {
					return n
				}
			}

		case int:
			if value >= 1 && value <= AllSides {
				return value
			}
		}
	}
	return AllSides
}

func (resizable *resizableData) setSide(value any) bool {
	switch value := value.(type) {
	case string:
		if n, err := strconv.Atoi(value); err == nil {
			if n >= 1 && n <= AllSides {
				resizable.properties[Side] = n
				return true
			}
			return false
		}
		validValues := map[string]int{
			"top":    TopSide,
			"right":  RightSide,
			"bottom": BottomSide,
			"left":   LeftSide,
			"all":    AllSides,
		}
		if strings.Contains(value, "|") {
			values := strings.Split(value, "|")
			sides := 0
			hasConst := false
			for i, val := range values {
				val := strings.Trim(val, " \t\r\n")
				values[i] = val

				if val[0] == '@' {
					hasConst = true
				} else if n, err := strconv.Atoi(val); err == nil {
					if n < 1 || n > AllSides {
						return false
					}
					sides |= n
				} else if n, ok := validValues[val]; ok {
					sides |= n
				} else {
					return false
				}
			}

			if hasConst {
				value = values[0]
				for i := 1; i < len(values); i++ {
					value += "|" + values[i]
				}
				resizable.properties[Side] = value
				return true
			}

			if sides >= 1 && sides <= AllSides {
				resizable.properties[Side] = sides
				return true
			}

		} else if value[0] == '@' {
			resizable.properties[Side] = value
			return true
		} else if n, ok := validValues[value]; ok {
			resizable.properties[Side] = n
			return true
		}

	case int:
		if value >= 1 && value <= AllSides {
			resizable.properties[Side] = value
			return true
		} else {
			ErrorLogF(`Invalid value %d of "side" property`, value)
			return false
		}

	default:
		if n, ok := isInt(value); ok {
			if n >= 1 && n <= AllSides {
				resizable.properties[Side] = n
				return true
			} else {
				ErrorLogF(`Invalid value %d of "side" property`, n)
				return false
			}
		}
	}

	return false
}

func (resizable *resizableData) resizeBorderWidth() SizeUnit {
	result, _ := sizeProperty(resizable, ResizeBorderWidth, resizable.Session())
	if result.Type == Auto || result.Value == 0 {
		return Px(4)
	}
	return result
}

func (resizable *resizableData) updateResizeBorderWidth() {
	if resizable.created {
		htmlID := resizable.htmlID()
		session := resizable.Session()
		column, row := resizable.cellSizeCSS()

		session.updateCSSProperty(htmlID, "grid-template-columns", column)
		session.updateCSSProperty(htmlID, "grid-template-rows", row)
	}
}

func (resizable *resizableData) cellSizeCSS() (string, string) {
	w := resizable.resizeBorderWidth().cssString("4px", resizable.Session())
	side := resizable.getSide()
	column := "1fr"
	row := "1fr"

	if side&LeftSide != 0 {
		if (side & RightSide) != 0 {
			column = w + " 1fr " + w
		} else {
			column = w + " 1fr"
		}
	} else if (side & RightSide) != 0 {
		column = "1fr " + w
	}

	if side&TopSide != 0 {
		if (side & BottomSide) != 0 {
			row = w + " 1fr " + w
		} else {
			row = w + " 1fr"
		}
	} else if (side & BottomSide) != 0 {
		row = "1fr " + w
	}

	return column, row
}

func (resizable *resizableData) cssStyle(self View, builder cssBuilder) {
	column, row := resizable.cellSizeCSS()

	builder.add("grid-template-columns", column)
	builder.add("grid-template-rows", row)

	resizable.viewData.cssStyle(self, builder)
}

func (resizable *resizableData) htmlSubviews(self View, buffer *strings.Builder) {

	side := resizable.getSide()
	left := 1
	top := 1
	leftSide := (side & LeftSide) != 0
	rightSide := (side & RightSide) != 0
	w := resizable.resizeBorderWidth().cssString("4px", resizable.Session())

	if leftSide {
		left = 2
	}

	writePos := func(x1, x2, y1, y2 int) {
		buffer.WriteString(fmt.Sprintf(` grid-column-start: %d; grid-column-end: %d; grid-row-start: %d;  grid-row-end: %d;"></div>`, x1, x2, y1, y2))
	}
	//htmlID := resizable.htmlID()

	if (side & TopSide) != 0 {
		top = 2

		if leftSide {
			buffer.WriteString(`<div onmousedown="startResize(this, -1, -1, event)" style="cursor: nwse-resize; width: `)
			buffer.WriteString(w)
			buffer.WriteString(`; height: `)
			buffer.WriteString(w)
			buffer.WriteString(`;`)
			writePos(1, 2, 1, 2)
		}

		buffer.WriteString(`<div onmousedown="startResize(this, 0, -1, event)" style="cursor: ns-resize; width: 100%; height: `)
		buffer.WriteString(w)
		buffer.WriteString(`;`)
		writePos(left, left+1, 1, 2)

		if rightSide {
			buffer.WriteString(`<div onmousedown="startResize(this, 1, -1, event)" style="cursor: nesw-resize; width: `)
			buffer.WriteString(w)
			buffer.WriteString(`; height: `)
			buffer.WriteString(w)
			buffer.WriteString(`;`)
			writePos(left+1, left+2, 1, 2)
		}
	}

	if leftSide {
		buffer.WriteString(`<div onmousedown="startResize(this, -1, 0, event)" style="cursor: ew-resize; width: `)
		buffer.WriteString(w)
		buffer.WriteString(`; height: 100%;`)
		writePos(1, 2, top, top+1)
	}

	if rightSide {
		buffer.WriteString(`<div onmousedown="startResize(this, 1, 0, event)" style="cursor: ew-resize; width: `)
		buffer.WriteString(w)
		buffer.WriteString(`; height: 100%;`)
		writePos(left+1, left+2, top, top+1)
	}

	if (side & BottomSide) != 0 {
		if leftSide {
			buffer.WriteString(`<div onmousedown="startResize(this, -1, 1, event)" style="cursor: nesw-resize; width: `)
			buffer.WriteString(w)
			buffer.WriteString(`; height: `)
			buffer.WriteString(w)
			buffer.WriteString(`;`)
			writePos(1, 2, top+1, top+2)
		}

		buffer.WriteString(`<div onmousedown="startResize(this, 0, 1, event)" style="cursor: ns-resize; width: 100%; height: `)
		buffer.WriteString(w)
		buffer.WriteString(`;`)
		writePos(left, left+1, top+1, top+2)

		if rightSide {
			buffer.WriteString(`<div onmousedown="startResize(this, 1, 1, event)" style="cursor: nwse-resize; width: `)
			buffer.WriteString(w)
			buffer.WriteString(`; height: `)
			buffer.WriteString(w)
			buffer.WriteString(`;`)
			writePos(left+1, left+2, top+1, top+2)
		}
	}

	if len(resizable.content) > 0 {
		view := resizable.content[0]
		view.addToCSSStyle(map[string]string{
			"grid-column-start": strconv.Itoa(left),
			"grid-column-end":   strconv.Itoa(left + 1),
			"grid-row-start":    strconv.Itoa(top),
			"grid-row-end":      strconv.Itoa(top + 1),
		})
		viewHTML(view, buffer)
	}
}
