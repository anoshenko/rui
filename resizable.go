package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants for [Resizable] specific properties and events
const (
	// Side is the constant for "side" property tag.
	//
	// Used by `Resizable`.
	// Determines which side of the container is used to resize. The value of property is an or-combination of values listed.
	// Default value is "all".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `1`(`TopSide`) or "top" - Top frame side.
	// `2`(`RightSide`) or "right" - Right frame side.
	// `4`(`BottomSide`) or "bottom" - Bottom frame side.
	// `8`(`LeftSide`) or "left" - Left frame side.
	// `15`(`AllSides`) or "all" - All frame sides.
	Side = "side"

	// ResizeBorderWidth is the constant for "resize-border-width" property tag.
	//
	// Used by `Resizable`.
	// Specifies the width of the resizing border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	ResizeBorderWidth = "resize-border-width"
)

// Constants for values of [Resizable] "side" property. These constants can be ORed if needed.
const (
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

// Resizable represents a Resizable view
type Resizable interface {
	View
	ParentView
}

type resizableData struct {
	viewData
}

// NewResizable create new Resizable object and return it
func NewResizable(session Session, params Params) Resizable {
	view := new(resizableData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newResizable(session Session) View {
	return new(resizableData)
}

func (resizable *resizableData) init(session Session) {
	resizable.viewData.init(session)
	resizable.tag = "Resizable"
	resizable.systemClass = "ruiGridLayout"
	resizable.set = resizable.setFunc
	resizable.changed = resizable.propertyChanged
}

func (resizable *resizableData) Views() []View {
	if view := resizable.content(); view != nil {
		return []View{view}
	}
	return []View{}
}

func (resizable *resizableData) content() View {
	if value := resizable.getRaw(Content); value != nil {
		if content, ok := value.(View); ok {
			return content
		}
	}
	return nil
}

func (resizable *resizableData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Side:
		return resizableSetSide(resizable, value)

	case ResizeBorderWidth:
		return setSizeProperty(resizable, tag, value)

	case Content:
		var newContent View = nil
		switch value := value.(type) {
		case string:
			newContent = NewTextView(resizable.Session(), Params{Text: value})

		case View:
			newContent = value

		case DataObject:
			if newContent = CreateViewFromObject(resizable.Session(), value); newContent == nil {
				return nil
			}

		default:
			notCompatibleType(tag, value)
			return nil
		}

		resizable.setRaw(Content, newContent)
		return []PropertyName{}

	case CellWidth, CellHeight, GridRowGap, GridColumnGap, CellVerticalAlign, CellHorizontalAlign:
		ErrorLogF(`Not supported "%s" property`, string(tag))
		return nil
	}

	return resizable.viewData.setFunc(tag, value)
}

func (resizable *resizableData) propertyChanged(tag PropertyName) {
	switch tag {
	case Side:
		updateInnerHTML(resizable.htmlID(), resizable.Session())
		fallthrough

	case ResizeBorderWidth:
		htmlID := resizable.htmlID()
		session := resizable.Session()
		column, row := resizableCellSizeCSS(resizable)

		session.updateCSSProperty(htmlID, "grid-template-columns", column)
		session.updateCSSProperty(htmlID, "grid-template-rows", row)

	case Content:
		updateInnerHTML(resizable.htmlID(), resizable.Session())

	default:
		resizable.viewData.propertyChanged(tag)
	}

}

func resizableSide(view View) int {
	if value := view.getRaw(Side); value != nil {
		switch value := value.(type) {
		case string:
			if value, ok := view.Session().resolveConstants(value); ok {
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

func resizableSetSide(properties Properties, value any) []PropertyName {
	switch value := value.(type) {
	case string:
		if n, err := strconv.Atoi(value); err == nil {
			if n >= 1 && n <= AllSides {
				properties.setRaw(Side, n)
				return []PropertyName{Side}
			}
			return nil
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
						return nil
					}
					sides |= n
				} else if n, ok := validValues[val]; ok {
					sides |= n
				} else {
					return nil
				}
			}

			if hasConst {
				value = values[0]
				for i := 1; i < len(values); i++ {
					value += "|" + values[i]
				}
				properties.setRaw(Side, value)
				return []PropertyName{Side}
			}

			if sides >= 1 && sides <= AllSides {
				properties.setRaw(Side, sides)
				return []PropertyName{Side}
			}

		} else if value[0] == '@' {
			properties.setRaw(Side, value)
			return []PropertyName{Side}
		} else if n, ok := validValues[value]; ok {
			properties.setRaw(Side, n)
			return []PropertyName{Side}
		}

	case int:
		if value >= 1 && value <= AllSides {
			properties.setRaw(Side, value)
			return []PropertyName{Side}
		} else {
			ErrorLogF(`Invalid value %d of "side" property`, value)
			return nil
		}

	default:
		if n, ok := isInt(value); ok {
			if n >= 1 && n <= AllSides {
				properties.setRaw(Side, n)
				return []PropertyName{Side}
			} else {
				ErrorLogF(`Invalid value %d of "side" property`, n)
				return nil
			}
		}
	}

	return nil
}

func resizableBorderWidth(view View) SizeUnit {
	result, _ := sizeProperty(view, ResizeBorderWidth, view.Session())
	if result.Type == Auto || result.Value == 0 {
		return Px(4)
	}
	return result
}

func resizableCellSizeCSS(view View) (string, string) {
	w := resizableBorderWidth(view).cssString("4px", view.Session())
	side := resizableSide(view)
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
	column, row := resizableCellSizeCSS(resizable)

	builder.add("grid-template-columns", column)
	builder.add("grid-template-rows", row)

	resizable.viewData.cssStyle(self, builder)
}

func (resizable *resizableData) htmlSubviews(self View, buffer *strings.Builder) {

	side := resizableSide(resizable)
	left := 1
	top := 1
	leftSide := (side & LeftSide) != 0
	rightSide := (side & RightSide) != 0
	w := resizableBorderWidth(resizable).cssString("4px", resizable.Session())

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

	if view := resizable.content(); view != nil {
		view.addToCSSStyle(map[string]string{
			"grid-column-start": strconv.Itoa(left),
			"grid-column-end":   strconv.Itoa(left + 1),
			"grid-row-start":    strconv.Itoa(top),
			"grid-row-end":      strconv.Itoa(top + 1),
		})
		viewHTML(view, buffer)
	}
}
