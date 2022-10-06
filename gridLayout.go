package rui

import (
	"fmt"
	"strings"
)

// GridLayout - grid-container of View
type GridLayout interface {
	ViewsContainer
}

type gridLayoutData struct {
	viewsContainerData
}

// NewGridLayout create new GridLayout object and return it
func NewGridLayout(session Session, params Params) GridLayout {
	view := new(gridLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newGridLayout(session Session) View {
	return NewGridLayout(session, nil)
}

// Init initialize fields of GridLayout by default values
func (gridLayout *gridLayoutData) init(session Session) {
	gridLayout.viewsContainerData.init(session)
	gridLayout.tag = "GridLayout"
	gridLayout.systemClass = "ruiGridLayout"
}

func (gridLayout *gridLayoutData) String() string {
	return getViewString(gridLayout)
}

func (style *viewStyle) setGridCellSize(tag string, value any) bool {
	setValues := func(values []string) bool {
		count := len(values)
		if count > 1 {
			sizes := make([]any, count)
			for i, val := range values {
				val = strings.Trim(val, " \t\n\r")
				if isConstantName(val) {
					sizes[i] = val
				} else if fn := parseSizeFunc(val); fn != nil {
					sizes[i] = SizeUnit{Type: SizeFunction, Function: fn}
				} else if size, err := stringToSizeUnit(val); err == nil {
					sizes[i] = size
				} else {
					invalidPropertyValue(tag, value)
					return false
				}
			}
			style.properties[tag] = sizes
		} else if isConstantName(values[0]) {
			style.properties[tag] = values[0]
		} else if size, err := stringToSizeUnit(values[0]); err == nil {
			style.properties[tag] = size
		} else {
			invalidPropertyValue(tag, value)
			return false
		}
		return true
	}

	switch tag {
	case CellWidth, CellHeight:
		switch value := value.(type) {
		case SizeUnit, []SizeUnit:
			style.properties[tag] = value

		case string:
			if !setValues(strings.Split(value, ",")) {
				return false
			}

		case []string:
			if !setValues(value) {
				return false
			}

		case []DataValue:
			count := len(value)
			if count == 0 {
				invalidPropertyValue(tag, value)
				return false
			}
			values := make([]string, count)
			for i, val := range value {
				if val.IsObject() {
					invalidPropertyValue(tag, value)
					return false
				}
				values[i] = val.Value()
			}
			if !setValues(values) {
				return false
			}

		case []any:
			count := len(value)
			if count == 0 {
				invalidPropertyValue(tag, value)
				return false
			}
			sizes := make([]any, count)
			for i, val := range value {
				switch val := val.(type) {
				case SizeUnit:
					sizes[i] = val

				case string:
					if isConstantName(val) {
						sizes[i] = val
					} else if size, err := stringToSizeUnit(val); err == nil {
						sizes[i] = size
					} else {
						invalidPropertyValue(tag, value)
						return false
					}

				default:
					invalidPropertyValue(tag, value)
					return false
				}
			}
			style.properties[tag] = sizes

		default:
			notCompatibleType(tag, value)
			return false
		}

		return true
	}

	return false
}

func (style *viewStyle) gridCellSizesCSS(tag string, session Session) string {
	switch cellSize := gridCellSizes(style, tag, session); len(cellSize) {
	case 0:

	case 1:
		if cellSize[0].Type != Auto {
			return `repeat(auto-fill, ` + cellSize[0].cssString(`auto`, session) + `)`
		}

	default:
		allAuto := true
		allEqual := true
		for i, size := range cellSize {
			if size.Type != Auto {
				allAuto = false
			}
			if i > 0 && !size.Equal(cellSize[0]) {
				allEqual = false
			}
		}
		if !allAuto {
			if allEqual {
				return fmt.Sprintf(`repeat(%d, %s)`, len(cellSize), cellSize[0].cssString(`auto`, session))
			}

			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)
			for _, size := range cellSize {
				buffer.WriteRune(' ')
				buffer.WriteString(size.cssString(`auto`, session))
			}
			return buffer.String()
		}
	}

	return ""
}

func (gridLayout *gridLayoutData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case VerticalAlign:
		return CellVerticalAlign

	case HorizontalAlign:
		return CellHorizontalAlign

	case "row-gap":
		return GridRowGap

	case ColumnGap:
		return GridColumnGap
	}
	return tag
}

func (gridLayout *gridLayoutData) Get(tag string) any {
	return gridLayout.get(gridLayout.normalizeTag(tag))
}

func (gridLayout *gridLayoutData) get(tag string) any {
	if tag == Gap {
		rowGap := GetGridRowGap(gridLayout)
		columnGap := GetGridColumnGap(gridLayout)
		if rowGap.Equal(columnGap) {
			return rowGap
		}
		return AutoSize()
	}

	return gridLayout.viewsContainerData.get(tag)
}

func (gridLayout *gridLayoutData) Remove(tag string) {
	gridLayout.remove(gridLayout.normalizeTag(tag))
}

func (gridLayout *gridLayoutData) remove(tag string) {
	if tag == Gap {
		gridLayout.remove(GridRowGap)
		gridLayout.remove(GridColumnGap)
		gridLayout.propertyChangedEvent(Gap)
		return
	}

	gridLayout.viewsContainerData.remove(tag)
	if gridLayout.created {
		switch tag {
		case CellWidth:
			updateCSSProperty(gridLayout.htmlID(), `grid-template-columns`,
				gridLayout.gridCellSizesCSS(CellWidth, gridLayout.session), gridLayout.session)

		case CellHeight:
			updateCSSProperty(gridLayout.htmlID(), `grid-template-rows`,
				gridLayout.gridCellSizesCSS(CellHeight, gridLayout.session), gridLayout.session)

		}
	}
}

func (gridLayout *gridLayoutData) Set(tag string, value any) bool {
	return gridLayout.set(gridLayout.normalizeTag(tag), value)
}

func (gridLayout *gridLayoutData) set(tag string, value any) bool {
	if value == nil {
		gridLayout.remove(tag)
		return true
	}

	if tag == Gap {
		return gridLayout.set(GridRowGap, value) && gridLayout.set(GridColumnGap, value)
	}

	if gridLayout.viewsContainerData.set(tag, value) {
		if gridLayout.created {
			switch tag {
			case CellWidth:
				updateCSSProperty(gridLayout.htmlID(), `grid-template-columns`,
					gridLayout.gridCellSizesCSS(CellWidth, gridLayout.session), gridLayout.session)

			case CellHeight:
				updateCSSProperty(gridLayout.htmlID(), `grid-template-rows`,
					gridLayout.gridCellSizesCSS(CellHeight, gridLayout.session), gridLayout.session)

			}
		}
		return true
	}

	return false
}

func gridCellSizes(properties Properties, tag string, session Session) []SizeUnit {
	if value := properties.Get(tag); value != nil {
		switch value := value.(type) {
		case []SizeUnit:
			return value

		case SizeUnit:
			return []SizeUnit{value}

		case []any:
			result := make([]SizeUnit, len(value))
			for i, val := range value {
				result[i] = AutoSize()
				switch val := val.(type) {
				case SizeUnit:
					result[i] = val

				case string:
					if text, ok := session.resolveConstants(val); ok {
						result[i], _ = stringToSizeUnit(text)
					}
				}
			}
			return result

		case string:
			if text, ok := session.resolveConstants(value); ok {
				values := strings.Split(text, ",")
				result := make([]SizeUnit, len(values))
				for i, val := range values {
					result[i], _ = stringToSizeUnit(val)
				}
				return result
			}
		}
	}

	return []SizeUnit{}
}

/*
func (gridLayout *gridLayoutData) cssStyle(self View, builder cssBuilder) {
	gridLayout.viewsContainerData.cssStyle(self, builder)
}
*/

// GetCellVerticalAlign returns the vertical align of a GridLayout cell content: TopAlign (0), BottomAlign (1), CenterAlign (2), StretchAlign (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCellVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CellVerticalAlign, StretchAlign, false)
}

// GetCellHorizontalAlign returns the vertical align of a GridLayout cell content: LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCellHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CellHorizontalAlign, StretchAlign, false)
}

// GetGridAutoFlow returns the value of the  "grid-auto-flow" property
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetGridAutoFlow(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, GridAutoFlow, 0, false)
}

// GetCellWidth returns the width of a GridLayout cell. If the result is an empty array, then the width is not set.
// If the result is a single value array, then the width of all cell is equal.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCellWidth(view View, subviewID ...string) []SizeUnit {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		return gridCellSizes(view, CellWidth, view.Session())
	}
	return []SizeUnit{}
}

// GetCellHeight returns the height of a GridLayout cell. If the result is an empty array, then the height is not set.
// If the result is a single value array, then the height of all cell is equal.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCellHeight(view View, subviewID ...string) []SizeUnit {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		return gridCellSizes(view, CellHeight, view.Session())
	}
	return []SizeUnit{}
}

// GetGridRowGap returns the gap between GridLayout rows.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetGridRowGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, GridRowGap, false)
}

// GetGridColumnGap returns the gap between GridLayout columns.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetGridColumnGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, GridColumnGap, false)
}
