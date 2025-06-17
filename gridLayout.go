package rui

import (
	"fmt"
	"strings"
)

// Constants related to [GridLayout] specific properties and events
const (
	// CellVerticalAlign is the constant for "cell-vertical-align" property tag.
	//
	// Used by GridLayout, SvgImageView.
	//
	// Usage in GridLayout:
	// Sets the default vertical alignment of GridLayout children within the cell they are occupying.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Full height stretch.
	//
	// Usage in SvgImageView:
	// Same as "vertical-align".
	CellVerticalAlign PropertyName = "cell-vertical-align"

	// CellHorizontalAlign is the constant for "cell-horizontal-align" property tag.
	//
	// Used by GridLayout, SvgImageView.
	//
	// Usage in GridLayout:
	// Sets the default horizontal alignment of GridLayout children within the occupied cell.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Full width stretch.
	//
	// Usage in SvgImageView:
	// Same as "horizontal-align".
	CellHorizontalAlign PropertyName = "cell-horizontal-align"

	// CellVerticalSelfAlign is the constant for "cell-vertical-self-align" property tag.
	//
	// Used by GridLayout.
	// Sets the vertical alignment of GridLayout children within the cell they are occupying. The property is set for the
	// child view of GridLayout.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Full height stretch.
	CellVerticalSelfAlign PropertyName = "cell-vertical-self-align"

	// CellHorizontalSelfAlign is the constant for "cell-horizontal-self-align" property tag.
	//
	// Used by GridLayout.
	// Sets the horizontal alignment of GridLayout children within the occupied cell. The property is set for the child view
	// of GridLayout.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Full width stretch.
	CellHorizontalSelfAlign PropertyName = "cell-horizontal-self-align"
)

// GridAdapter is an interface to define [GridLayout] content. [GridLayout] will query interface functions to populate
// its content
type GridAdapter interface {
	// GridColumnCount returns the number of columns in the grid
	GridColumnCount() int

	// GridRowCount returns the number of rows in the grid
	GridRowCount() int

	// GridCellContent creates a View at the given cell
	GridCellContent(row, column int, session Session) View
}

// GridCellColumnSpanAdapter implements the optional method of the [GridAdapter] interface
type GridCellColumnSpanAdapter interface {
	// GridCellColumnSpan returns the number of columns that a cell spans.
	// Values ​​less than 1 are ignored.
	GridCellColumnSpan(row, column int) int
}

// GridCellColumnSpanAdapter implements the optional method of the [GridAdapter] interface
type GridCellRowSpanAdapter interface {
	// GridCellRowSpan returns the number of rows that a cell spans
	// Values ​​less than 1 are ignored.
	GridCellRowSpan(row, column int) int
}

// GridLayout represents a GridLayout view
type GridLayout interface {
	ViewsContainer

	// UpdateContent updates child Views if the "content" property value is set to GridAdapter,
	// otherwise does nothing
	UpdateGridContent()
}

type gridLayoutData struct {
	viewsContainerData
	adapter GridAdapter
}

// NewGridLayout create new GridLayout object and return it
func NewGridLayout(session Session, params Params) GridLayout {
	view := new(gridLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newGridLayout(session Session) View {
	//return NewGridLayout(session, nil)
	return new(gridLayoutData)
}

// Init initialize fields of GridLayout by default values
func (gridLayout *gridLayoutData) init(session Session) {
	gridLayout.viewsContainerData.init(session)
	gridLayout.tag = "GridLayout"
	gridLayout.systemClass = "ruiGridLayout"
	gridLayout.adapter = nil
	gridLayout.normalize = normalizeGridLayoutTag
	gridLayout.get = gridLayout.getFunc
	gridLayout.set = gridLayout.setFunc
	gridLayout.remove = gridLayout.removeFunc
	gridLayout.changed = gridLayout.propertyChanged
}

func setGridCellSize(properties Properties, tag PropertyName, value any) []PropertyName {
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
			properties.setRaw(tag, sizes)
		} else if isConstantName(values[0]) {
			properties.setRaw(tag, values[0])
		} else if size, err := stringToSizeUnit(values[0]); err == nil {
			properties.setRaw(tag, size)
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
			properties.setRaw(tag, value)

		case string:
			if !setValues(strings.Split(value, ",")) {
				return nil
			}

		case []string:
			if !setValues(value) {
				return nil
			}

		case []DataValue:
			count := len(value)
			if count == 0 {
				invalidPropertyValue(tag, value)
				return nil
			}
			values := make([]string, count)
			for i, val := range value {
				if val.IsObject() {
					invalidPropertyValue(tag, value)
					return nil
				}
				values[i] = val.Value()
			}
			if !setValues(values) {
				return nil
			}

		case []any:
			count := len(value)
			if count == 0 {
				invalidPropertyValue(tag, value)
				return nil
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
						return nil
					}

				default:
					invalidPropertyValue(tag, value)
					return nil
				}
			}
			properties.setRaw(tag, sizes)

		default:
			notCompatibleType(tag, value)
			return nil
		}

		return []PropertyName{tag}
	}

	return nil
}

func gridCellSizesCSS(properties Properties, tag PropertyName, session Session) string {
	switch cellSize := gridCellSizes(properties, tag, session); len(cellSize) {
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

func normalizeGridLayoutTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
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

func (gridLayout *gridLayoutData) getFunc(tag PropertyName) any {
	switch tag {
	case Gap:
		rowGap := GetGridRowGap(gridLayout)
		columnGap := GetGridColumnGap(gridLayout)
		if rowGap.Equal(columnGap) {
			return rowGap
		}
		return AutoSize()

	case Content:
		if gridLayout.adapter != nil {
			return gridLayout.adapter
		}
	}

	return gridLayout.viewsContainerData.getFunc(tag)
}

func (gridLayout *gridLayoutData) removeFunc(tag PropertyName) []PropertyName {
	switch tag {
	case Gap:
		result := []PropertyName{}
		for _, tag := range []PropertyName{GridRowGap, GridColumnGap} {
			if gridLayout.getRaw(tag) != nil {
				gridLayout.setRaw(tag, nil)
				result = append(result, tag)
			}
		}
		return result

	case Content:
		if len(gridLayout.views) > 0 || gridLayout.adapter != nil {
			gridLayout.views = []View{}
			gridLayout.adapter = nil
			return []PropertyName{Content}
		}
		return []PropertyName{}
	}

	return gridLayout.viewsContainerData.removeFunc(tag)
}

func (gridLayout *gridLayoutData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Gap:
		result := gridLayout.setFunc(GridRowGap, value)
		if result != nil {
			if gap := gridLayout.getRaw(GridRowGap); gap != nil {
				gridLayout.setRaw(GridColumnGap, gap)
				result = append(result, GridColumnGap)
			}
		}
		return result

	case Content:
		if adapter, ok := value.(GridAdapter); ok {
			gridLayout.adapter = adapter
			gridLayout.createGridContent()
		} else if gridLayout.setContent(value) {
			gridLayout.adapter = nil
		} else {
			return nil
		}
		return []PropertyName{Content}
	}

	return gridLayout.viewsContainerData.setFunc(tag, value)
}

func (gridLayout *gridLayoutData) propertyChanged(tag PropertyName) {
	switch tag {
	case CellWidth:
		session := gridLayout.Session()
		session.updateCSSProperty(gridLayout.htmlID(), `grid-template-columns`,
			gridCellSizesCSS(gridLayout, CellWidth, session))

	case CellHeight:
		session := gridLayout.Session()
		session.updateCSSProperty(gridLayout.htmlID(), `grid-template-rows`,
			gridCellSizesCSS(gridLayout, CellHeight, session))

	default:
		gridLayout.viewsContainerData.propertyChanged(tag)
	}
}

func (gridLayout *gridLayoutData) createGridContent() bool {
	if gridLayout.adapter == nil {
		return false
	}

	adapter := gridLayout.adapter
	gridLayout.views = []View{}

	session := gridLayout.session
	htmlID := gridLayout.htmlID()
	isDisabled := IsDisabled(gridLayout)

	var columnSpan GridCellColumnSpanAdapter = nil
	if span, ok := adapter.(GridCellColumnSpanAdapter); ok {
		columnSpan = span
	}

	var rowSpan GridCellRowSpanAdapter = nil
	if span, ok := adapter.(GridCellRowSpanAdapter); ok {
		rowSpan = span
	}

	width := adapter.GridColumnCount()
	height := adapter.GridRowCount()
	for column := 0; column < width; column++ {
		for row := 0; row < height; row++ {
			if view := adapter.GridCellContent(row, column, session); view != nil {
				view.setParentID(htmlID)

				columnCount := 1
				if columnSpan != nil {
					columnCount = columnSpan.GridCellColumnSpan(row, column)
				}

				if columnCount > 1 {
					view.Set(Column, Range{First: column, Last: column + columnCount - 1})
				} else {
					view.Set(Column, column)
				}

				rowCount := 1
				if rowSpan != nil {
					rowCount = rowSpan.GridCellRowSpan(row, column)
				}

				if rowCount > 1 {
					view.Set(Row, Range{First: row, Last: row + rowCount - 1})
				} else {
					view.Set(Row, row)
				}

				if isDisabled {
					view.Set(Disabled, true)
				}

				gridLayout.views = append(gridLayout.views, view)
			}
		}
	}

	return true
}

func (gridLayout *gridLayoutData) UpdateGridContent() {
	if gridLayout.createGridContent() {
		if gridLayout.created {
			updateInnerHTML(gridLayout.htmlID(), gridLayout.session)
		}

		if listener, ok := gridLayout.changeListener[Content]; ok {
			listener(gridLayout, Content)
		}
	}
}

func gridCellSizes(properties Properties, tag PropertyName, session Session) []SizeUnit {
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

// GetCellVerticalAlign returns the vertical align of a GridLayout cell content: TopAlign (0), BottomAlign (1), CenterAlign (2), StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetCellVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CellVerticalAlign, StretchAlign, false)
}

// GetCellHorizontalAlign returns the vertical align of a GridLayout cell content: LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetCellHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CellHorizontalAlign, StretchAlign, false)
}

// GetGridAutoFlow returns the value of the  "grid-auto-flow" property
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetGridAutoFlow(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, GridAutoFlow, 0, false)
}

// GetCellWidth returns the width of a GridLayout cell. If the result is an empty array, then the width is not set.
// If the result is a single value array, then the width of all cell is equal.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetCellWidth(view View, subviewID ...string) []SizeUnit {
	if view = getSubview(view, subviewID); view != nil {
		return gridCellSizes(view, CellWidth, view.Session())
	}
	return []SizeUnit{}
}

// GetCellHeight returns the height of a GridLayout cell. If the result is an empty array, then the height is not set.
// If the result is a single value array, then the height of all cell is equal.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetCellHeight(view View, subviewID ...string) []SizeUnit {
	if view = getSubview(view, subviewID); view != nil {
		return gridCellSizes(view, CellHeight, view.Session())
	}
	return []SizeUnit{}
}

// GetGridRowGap returns the gap between GridLayout rows.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetGridRowGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, GridRowGap, false)
}

// GetGridColumnGap returns the gap between GridLayout columns.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetGridColumnGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, GridColumnGap, false)
}
