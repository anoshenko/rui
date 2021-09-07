package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// TableVerticalAlign is the constant for the "table-vertical-align" property tag.
	// The "table-vertical-align" int property sets the vertical alignment of the content inside a table cell.
	// Valid values are LeftAlign (0), RightAlign (1), CenterAlign (2), and BaselineAlign (3, 4)
	TableVerticalAlign = "table-vertical-align"
	// HeadHeight is the constant for the "head-height" property tag.
	// The "head-height" int property sets the number of rows in the table header.
	// The default value is 0 (no header)
	HeadHeight = "head-height"
	// HeadStyle is the constant for the "head-style" property tag.
	// The "head-style" string property sets the header style name
	HeadStyle = "head-style"
	// FootHeight is the constant for the "foot-height" property tag.
	// The "foot-height" int property sets the number of rows in the table footer.
	// The default value is 0 (no footer)
	FootHeight = "foot-height"
	// FootStyle is the constant for the "foot-style" property tag.
	// The "foot-style" string property sets the footer style name
	FootStyle = "foot-style"
	// RowSpan is the constant for the "row-span" property tag.
	// The "row-span" int property sets the number of table row to span.
	// Used only when specifying cell parameters in the implementation of TableCellStyle
	RowSpan = "row-span"
	// ColumnSpan is the constant for the "column-span" property tag.
	// The "column-span" int property sets the number of table column to span.
	// Used only when specifying cell parameters in the implementation of TableCellStyle
	ColumnSpan = "column-span"
	// RowStyle is the constant for the "row-style" property tag.
	// The "row-style" property sets the adapter which specifies styles of each table row.
	// This property can be assigned or by an implementation of TableRowStyle interface, or by an array of Params.
	RowStyle = "row-style"
	// ColumnStyle is the constant for the "column-style" property tag.
	// The "column-style" property sets the adapter which specifies styles of each table column.
	// This property can be assigned or by an implementation of TableColumnStyle interface, or by an array of Params.
	ColumnStyle = "column-style"
	// CellStyle is the constant for the "cell-style" property tag.
	// The "cell-style" property sets the adapter which specifies styles of each table cell.
	// This property can be assigned only by an implementation of TableCellStyle interface.
	CellStyle = "cell-style"
	// CellPadding is the constant for the "cell-padding" property tag.
	// The "cell-padding" Bounds property sets the padding area on all four sides of a table call at once.
	// An element's padding area is the space between its content and its border.
	CellPadding = "cell-padding"
	// CellPaddingLeft is the constant for the "cell-padding-left" property tag.
	// The "cell-padding-left" SizeUnit property sets the width of the padding area to the left of a cell content.
	// An element's padding area is the space between its content and its border.
	CellPaddingLeft = "cell-padding-left"
	// CellPaddingRight is the constant for the "cell-padding-right" property tag.
	// The "cell-padding-right" SizeUnit property sets the width of the padding area to the left of a cell content.
	// An element's padding area is the space between its content and its border.
	CellPaddingRight = "cell-padding-right"
	// CellPaddingTop is the constant for the "cell-padding-top" property tag.
	// The "cell-padding-top" SizeUnit property sets the height of the padding area to the top of a cell content.
	// An element's padding area is the space between its content and its border.
	CellPaddingTop = "cell-padding-top"
	// CellPaddingBottom is the constant for the "cell-padding-bottom" property tag.
	// The "cell-padding-bottom" SizeUnit property sets the height of the padding area to the bottom of a cell content.
	CellPaddingBottom = "cell-padding-bottom"
	// CellBorder is the constant for the "cell-border" property tag.
	// The "cell-border" property sets a table cell's border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorder = "cell-border"
	// CellBorderLeft is the constant for the "cell-border-left" property tag.
	// The "cell-border-left" property sets a view's left border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderLeft = "cell-border-left"
	// CellBorderRight is the constant for the "cell-border-right" property tag.
	// The "cell-border-right" property sets a view's right border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderRight = "cell-border-right"
	// CellBorderTop is the constant for the "cell-border-top" property tag.
	// The "cell-border-top" property sets a view's top border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderTop = "cell-border-top"
	// CellBorderBottom is the constant for the "cell-border-bottom" property tag.
	// The "cell-border-bottom" property sets a view's bottom border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderBottom = "cell-border-bottom"
	// CellBorderStyle is the constant for the "cell-border-style" property tag.
	// The "cell-border-style" int property sets the line style for all four sides of a table cell's border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderStyle = "cell-border-style"
	// CellBorderLeftStyle is the constant for the "cell-border-left-style" property tag.
	// The "cell-border-left-style" int property sets the line style of a table cell's left border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderLeftStyle = "cell-border-left-style"
	// CellBorderRightStyle is the constant for the "cell-border-right-style" property tag.
	// The "cell-border-right-style" int property sets the line style of a table cell's right border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderRightStyle = "cell-border-right-style"
	// CellBorderTopStyle is the constant for the "cell-border-top-style" property tag.
	// The "cell-border-top-style" int property sets the line style of a table cell's top border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderTopStyle = "cell-border-top-style"
	// CellBorderBottomStyle is the constant for the "cell-border-bottom-style" property tag.
	// The "cell-border-bottom-style" int property sets the line style of a table cell's bottom border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderBottomStyle = "cell-border-bottom-style"
	// CellBorderWidth is the constant for the "cell-border-width" property tag.
	// The "cell-border-width" property sets the line width for all four sides of a table cell's border.
	CellBorderWidth = "cell-border-width"
	// CellBorderLeftWidth is the constant for the "cell-border-left-width" property tag.
	// The "cell-border-left-width" SizeUnit property sets the line width of a table cell's left border.
	CellBorderLeftWidth = "cell-border-left-width"
	// CellBorderRightWidth is the constant for the "cell-border-right-width" property tag.
	// The "cell-border-right-width" SizeUnit property sets the line width of a table cell's right border.
	CellBorderRightWidth = "cell-border-right-width"
	// CellBorderTopWidth is the constant for the "cell-border-top-width" property tag.
	// The "cell-border-top-width" SizeUnit property sets the line width of a table cell's top border.
	CellBorderTopWidth = "cell-border-top-width"
	// CellBorderBottomWidth is the constant for the "cell-border-bottom-width" property tag.
	// The "cell-border-bottom-width" SizeUnit property sets the line width of a table cell's bottom border.
	CellBorderBottomWidth = "cell-border-bottom-width"
	// CellBorderColor is the constant for the "cell-border-color" property tag.
	// The "cell-border-color" property sets the line color for all four sides of a table cell's border.
	CellBorderColor = "cell-border-color"
	// CellBorderLeftColor is the constant for the "cell-border-left-color" property tag.
	// The "cell-border-left-color" property sets the line color of a table cell's left border.
	CellBorderLeftColor = "cell-border-left-color"
	// CellBorderRightColor is the constant for the "cell-border-right-color" property tag.
	// The "cell-border-right-color" property sets the line color of a table cell's right border.
	CellBorderRightColor = "cell-border-right-color"
	// CellBorderTopColor is the constant for the "cell-border-top-color" property tag.
	// The "cell-border-top-color" property sets the line color of a table cell's top border.
	CellBorderTopColor = "cell-border-top-color"
	// CellBorderBottomColor is the constant for the "cell-border-bottom-color" property tag.
	// The "cell-border-bottom-color" property sets the line color of a table cell's bottom border.
	CellBorderBottomColor = "cell-border-bottom-color"
)

// TableView - text View
type TableView interface {
	View
	ReloadTableData()
}

type tableViewData struct {
	viewData
}

type tableCellView struct {
	viewData
}

// NewTableView create new TableView object and return it
func NewTableView(session Session, params Params) TableView {
	view := new(tableViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newTableView(session Session) View {
	return NewTableView(session, nil)
}

// Init initialize fields of TableView by default values
func (table *tableViewData) Init(session Session) {
	table.viewData.Init(session)
	table.tag = "TableView"
}

func (table *tableViewData) Get(tag string) interface{} {
	return table.get(strings.ToLower(tag))
}

func (table *tableViewData) Remove(tag string) {
	table.remove(strings.ToLower(tag))
}

func (table *tableViewData) remove(tag string) {
	switch tag {

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft,
		"top-cell-padding", "right-cell-padding", "bottom-cell-padding", "left-cell-padding":
		table.removeBoundsSide(CellPadding, tag)

	case Gap, CellBorder, CellPadding, RowStyle, ColumnStyle, CellStyle,
		HeadHeight, HeadStyle, FootHeight, FootStyle:
		delete(table.properties, tag)

	default:
		table.viewData.remove(tag)
		return
	}

	table.propertyChanged(tag)
}

func (table *tableViewData) Set(tag string, value interface{}) bool {
	return table.set(strings.ToLower(tag), value)
}

func (table *tableViewData) set(tag string, value interface{}) bool {
	if value == nil {
		table.remove(tag)
		return true
	}

	switch tag {
	case Content:
		switch val := value.(type) {
		case TableAdapter:
			table.properties[Content] = value

		case [][]interface{}:
			table.properties[Content] = NewSimpleTableAdapter(val)

		case [][]string:
			table.properties[Content] = NewTextTableAdapter(val)

		default:
			notCompatibleType(tag, value)
			return false
		}

	case CellStyle:
		if style, ok := value.(TableCellStyle); ok {
			table.properties[tag] = style
		} else {
			notCompatibleType(tag, value)
			return false
		}

	case RowStyle:
		if !table.setRowStyle(value) {
			notCompatibleType(tag, value)
			return false
		}

	case ColumnStyle:
		if !table.setColumnStyle(value) {
			notCompatibleType(tag, value)
			return false
		}

	case HeadHeight, FootHeight:
		switch value := value.(type) {
		case string:
			if isConstantName(value) {
				table.properties[tag] = value
			} else if n, err := strconv.Atoi(value); err == nil {
				table.properties[tag] = n
			} else {
				ErrorLog(err.Error())
				notCompatibleType(tag, value)
				return false
			}

		default:
			if n, ok := isInt(value); ok {
				table.properties[tag] = n
			}
		}

	case HeadStyle, FootStyle:
		switch value := value.(type) {
		case string:
			table.properties[tag] = value

		case Params:
			if len(value) > 0 {
				table.properties[tag] = value
			} else {
				delete(table.properties, tag)
			}

		case DataNode:
			switch value.Type() {
			case ObjectNode:
				obj := value.Object()
				params := Params{}
				for k := 0; k < obj.PropertyCount(); k++ {
					if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
						params[prop.Tag()] = prop.Text()
					}
				}
				if len(params) > 0 {
					table.properties[tag] = params
				} else {
					delete(table.properties, tag)
				}

			case TextNode:
				table.properties[tag] = value.Text()

			default:
				notCompatibleType(tag, value)
				return false
			}

		default:
			notCompatibleType(tag, value)
			return false
		}

	case CellPadding:
		if !table.setBounds(tag, value) {
			return false
		}

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft,
		"top-cell-padding", "right-cell-padding", "bottom-cell-padding", "left-cell-padding":
		if !table.setBoundsSide(CellPadding, tag, value) {
			return false
		}

	case Gap:
		if !table.setSizeProperty(Gap, value) {
			return false
		}

	case CellBorder, CellBorderStyle, CellBorderColor, CellBorderWidth,
		CellBorderLeft, CellBorderLeftStyle, CellBorderLeftColor, CellBorderLeftWidth,
		CellBorderRight, CellBorderRightStyle, CellBorderRightColor, CellBorderRightWidth,
		CellBorderTop, CellBorderTopStyle, CellBorderTopColor, CellBorderTopWidth,
		CellBorderBottom, CellBorderBottomStyle, CellBorderBottomColor, CellBorderBottomWidth:
		if !table.viewData.set(tag, value) {
			return false
		}

	default:
		return table.viewData.set(tag, value)
	}

	table.propertyChanged(tag)
	return true
}

func (table *tableViewData) propertyChanged(tag string) {
	switch tag {
	case Content, RowStyle, ColumnStyle, CellStyle, CellPadding, CellBorder,
		HeadHeight, HeadStyle, FootHeight, FootStyle,
		CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft,
		"top-cell-padding", "right-cell-padding", "bottom-cell-padding", "left-cell-padding":
		table.ReloadTableData()

	case Gap:
		htmlID := table.htmlID()
		session := table.Session()
		gap, ok := sizeProperty(table, Gap, session)
		if !ok || gap.Type == Auto || gap.Value <= 0 {
			updateCSSProperty(htmlID, "border-spacing", "0", session)
			updateCSSProperty(htmlID, "border-collapse", "collapse", session)
		} else {
			updateCSSProperty(htmlID, "border-spacing", gap.cssString("0"), session)
			updateCSSProperty(htmlID, "border-collapse", "separate", session)
		}

	}
}

func (table *tableViewData) htmlTag() string {
	return "table"
}

func (table *tableViewData) htmlSubviews(self View, buffer *strings.Builder) {
	content := table.getRaw(Content)
	if content == nil {
		return
	}

	adapter, ok := content.(TableAdapter)
	if !ok {
		return
	}

	rowCount := adapter.RowCount()
	columnCount := adapter.ColumnCount()
	if rowCount == 0 || columnCount == 0 {
		return
	}

	rowStyle := table.getRowStyle()

	var cellStyle1 TableCellStyle = nil
	if style, ok := content.(TableCellStyle); ok {
		cellStyle1 = style
	}

	var cellStyle2 TableCellStyle = nil
	if value := table.getRaw(CellStyle); value != nil {
		if style, ok := value.(TableCellStyle); ok {
			cellStyle2 = style
		}
	}

	session := table.Session()

	if !session.ignoreViewUpdates() {
		session.setIgnoreViewUpdates(true)
		defer session.setIgnoreViewUpdates(false)
	}

	var cssBuilder viewCSSBuilder
	cssBuilder.buffer = allocStringBuilder()
	defer freeStringBuilder(cssBuilder.buffer)

	var view tableCellView
	view.Init(session)

	ignorCells := []struct{ row, column int }{}

	tableCSS := func(startRow, endRow int, cellTag string, cellBorder BorderProperty, cellPadding BoundsProperty) {
		for row := startRow; row < endRow; row++ {

			cssBuilder.buffer.Reset()
			if rowStyle != nil {
				if styles := rowStyle.RowStyle(row); styles != nil {
					view.Clear()
					for tag, value := range styles {
						view.Set(tag, value)
					}
					view.cssStyle(&view, &cssBuilder)
				}
			}

			if cssBuilder.buffer.Len() > 0 {
				buffer.WriteString(`<tr style="`)
				buffer.WriteString(cssBuilder.buffer.String())
				buffer.WriteString(`">`)
			} else {
				buffer.WriteString("<tr>")
			}

			for column := 0; column < columnCount; column++ {
				ignore := false
				for _, cell := range ignorCells {
					if cell.row == row && cell.column == column {
						ignore = true
						break
					}
				}

				if !ignore {
					rowSpan := 0
					columnSpan := 0

					cssBuilder.buffer.Reset()
					view.Clear()

					if cellBorder != nil {
						view.set(Border, cellBorder)
					}

					if cellPadding != nil {
						view.set(Padding, cellPadding)
					}

					appendFrom := func(cellStyle TableCellStyle) {
						if cellStyle != nil {
							if styles := cellStyle.CellStyle(row, column); styles != nil {
								for tag, value := range styles {
									valueToInt := func() int {
										switch value := value.(type) {
										case int:
											return value

										case string:
											if value, ok = session.resolveConstants(value); ok {
												if n, err := strconv.Atoi(value); err == nil {
													return n
												}
											}
										}
										return 0
									}

									switch tag = strings.ToLower(tag); tag {
									case RowSpan:
										rowSpan = valueToInt()

									case ColumnSpan:
										columnSpan = valueToInt()

									default:
										view.set(tag, value)
									}
								}
							}
						}
					}
					appendFrom(cellStyle1)
					appendFrom(cellStyle2)

					if len(view.properties) > 0 {
						view.cssStyle(&view, &cssBuilder)
					}

					buffer.WriteRune('<')
					buffer.WriteString(cellTag)

					if columnSpan > 1 {
						buffer.WriteString(` colspan="`)
						buffer.WriteString(strconv.Itoa(columnSpan))
						buffer.WriteRune('"')
						for c := column + 1; c < column+columnSpan; c++ {
							ignorCells = append(ignorCells, struct {
								row    int
								column int
							}{row: row, column: c})
						}
					}

					if rowSpan > 1 {
						buffer.WriteString(` rowspan="`)
						buffer.WriteString(strconv.Itoa(rowSpan))
						buffer.WriteRune('"')
						if columnSpan < 1 {
							columnSpan = 1
						}
						for r := row + 1; r < row+rowSpan; r++ {
							for c := column; c < column+columnSpan; c++ {
								ignorCells = append(ignorCells, struct {
									row    int
									column int
								}{row: r, column: c})
							}
						}
					}

					if cssBuilder.buffer.Len() > 0 {
						buffer.WriteString(` style="`)
						buffer.WriteString(cssBuilder.buffer.String())
						buffer.WriteRune('"')
					}
					buffer.WriteRune('>')

					switch value := adapter.Cell(row, column).(type) {
					case string:
						buffer.WriteString(value)

					case View:
						viewHTML(value, buffer)

					case Color:
						buffer.WriteString(`<div style="display: inline; height: 1em; background-color: `)
						buffer.WriteString(value.cssString())
						buffer.WriteString(`">&nbsp;&nbsp;&nbsp;&nbsp;</div> `)
						buffer.WriteString(value.String())

					case fmt.Stringer:
						buffer.WriteString(value.String())

					case rune:
						buffer.WriteRune(value)

					case float32:
						buffer.WriteString(fmt.Sprintf("%g", float64(value)))

					case float64:
						buffer.WriteString(fmt.Sprintf("%g", value))

					case bool:
						if value {
							buffer.WriteString(session.checkboxOnImage())
						} else {
							buffer.WriteString(session.checkboxOffImage())
						}

					default:
						if n, ok := isInt(value); ok {
							buffer.WriteString(fmt.Sprintf("%d", n))
						} else {
							buffer.WriteString("<Unsupported value>")
						}
					}

					buffer.WriteString(`</`)
					buffer.WriteString(cellTag)
					buffer.WriteRune('>')
				}
			}

			buffer.WriteString("</tr>")
		}
	}

	if columnStyle := table.getColumnStyle(); columnStyle != nil {
		buffer.WriteString("<colgroup>")
		for column := 0; column < columnCount; column++ {
			cssBuilder.buffer.Reset()
			if styles := columnStyle.ColumnStyle(column); styles != nil {
				view.Clear()
				for tag, value := range styles {
					view.Set(tag, value)
				}
				view.cssStyle(&view, &cssBuilder)
			}

			if cssBuilder.buffer.Len() > 0 {
				buffer.WriteString(`<col style="`)
				buffer.WriteString(cssBuilder.buffer.String())
				buffer.WriteString(`">`)
			} else {
				buffer.WriteString("<col>")
			}
		}
		buffer.WriteString("</colgroup>")
	}

	headHeight, _ := intProperty(table, HeadHeight, table.Session(), 0)
	footHeight, _ := intProperty(table, FootHeight, table.Session(), 0)
	cellBorder := table.getCellBorder()
	cellPadding := table.boundsProperty(CellPadding)
	if cellPadding == nil {
		if style, ok := stringProperty(table, Style, table.Session()); ok {
			if style, ok := table.Session().resolveConstants(style); ok {
				cellPadding = table.cellPaddingFromStyle(style)
			}
		}
	}

	headFootStart := func(htmlTag, styleTag string) (BorderProperty, BoundsProperty) {
		buffer.WriteRune('<')
		buffer.WriteString(htmlTag)
		if value := table.getRaw(styleTag); value != nil {
			switch value := value.(type) {
			case string:
				if style, ok := session.resolveConstants(value); ok {
					buffer.WriteString(` class="`)
					buffer.WriteString(style)
					buffer.WriteString(`">`)
					return table.cellBorderFromStyle(style), table.cellPaddingFromStyle(style)
				}

			case Params:
				cssBuilder.buffer.Reset()
				view.Clear()
				for tag, val := range value {
					view.Set(tag, val)
				}

				var border BorderProperty = nil
				if value := view.Get(CellBorder); value != nil {
					border = value.(BorderProperty)
				}
				var padding BoundsProperty = nil
				if value := view.Get(CellPadding); value != nil {
					switch value := value.(type) {
					case SizeUnit:
						padding = NewBoundsProperty(Params{
							Top:    value,
							Right:  value,
							Bottom: value,
							Left:   value,
						})

					case BoundsProperty:
						padding = value
					}
				}

				view.cssStyle(&view, &cssBuilder)
				if cssBuilder.buffer.Len() > 0 {
					buffer.WriteString(` style="`)
					buffer.WriteString(cssBuilder.buffer.String())
					buffer.WriteString(`"`)
				}
				buffer.WriteRune('>')
				return border, padding
			}
		}
		buffer.WriteRune('>')
		return nil, nil
	}

	if headHeight > 0 {
		headCellBorder := cellBorder
		headCellPadding := cellPadding

		if headHeight > rowCount {
			headHeight = rowCount
		}

		border, padding := headFootStart("thead", HeadStyle)
		if border != nil {
			headCellBorder = border
		}
		if padding != nil {
			headCellPadding = padding
		}
		tableCSS(0, headHeight, "th", headCellBorder, headCellPadding)
		buffer.WriteString("</thead>")
	}

	if footHeight > rowCount-headHeight {
		footHeight = rowCount - headHeight
	}

	if rowCount > footHeight+headHeight {
		buffer.WriteString("<tbody>")
		tableCSS(headHeight, rowCount-footHeight, "td", cellBorder, cellPadding)
		buffer.WriteString("</tbody>")
	}

	if footHeight > 0 {
		footCellBorder := cellBorder
		footCellPadding := cellPadding

		border, padding := headFootStart("tfoot", FootStyle)
		if border != nil {
			footCellBorder = border
		}
		if padding != nil {
			footCellPadding = padding
		}
		tableCSS(rowCount-footHeight, rowCount, "td", footCellBorder, footCellPadding)
		buffer.WriteString("</tfoot>")
	}
}

func (table *tableViewData) cellPaddingFromStyle(style string) BoundsProperty {
	session := table.Session()
	var result BoundsProperty = nil

	if node := session.stylePropertyNode(style, CellPadding); node != nil && node.Type() == ObjectNode {
		for _, tag := range []string{Left, Right, Top, Bottom} {
			if node := node.Object().PropertyWithTag(tag); node != nil && node.Type() == TextNode {
				if result == nil {
					result = NewBoundsProperty(nil)
				}
				result.Set(tag, node.Text())
			}
		}
	}

	for _, tag := range []string{CellPaddingLeft, CellPaddingRight, CellPaddingTop, CellPaddingBottom} {
		if value, ok := session.styleProperty(style, CellPadding); ok {
			if result == nil {
				result = NewBoundsProperty(nil)
			}
			result.Set(tag, value)
		}
	}

	return result
}

func (table *tableViewData) cellBorderFromStyle(style string) BorderProperty {

	border := new(borderProperty)
	border.properties = map[string]interface{}{}

	session := table.Session()
	if node := session.stylePropertyNode(style, CellBorder); node != nil && node.Type() == ObjectNode {
		border.setBorderObject(node.Object())
	}

	for _, tag := range []string{
		CellBorderLeft,
		CellBorderRight,
		CellBorderTop,
		CellBorderBottom,
		CellBorderStyle,
		CellBorderLeftStyle,
		CellBorderRightStyle,
		CellBorderTopStyle,
		CellBorderBottomStyle,
		CellBorderWidth,
		CellBorderLeftWidth,
		CellBorderRightWidth,
		CellBorderTopWidth,
		CellBorderBottomWidth,
		CellBorderColor,
		CellBorderLeftColor,
		CellBorderRightColor,
		CellBorderTopColor,
		CellBorderBottomColor,
	} {
		if value, ok := session.styleProperty(style, tag); ok {
			border.Set(tag, value)
		}
	}

	if len(border.properties) == 0 {
		return nil
	}
	return border
}

func (table *tableViewData) getCellBorder() BorderProperty {
	if value := table.getRaw(CellBorder); value != nil {
		if border, ok := value.(BorderProperty); ok {
			return border
		}
	}

	if style, ok := stringProperty(table, Style, table.Session()); ok {
		if style, ok := table.Session().resolveConstants(style); ok {
			return table.cellBorderFromStyle(style)
		}
	}

	return nil
}

func (table *tableViewData) cssStyle(self View, builder cssBuilder) {
	table.viewData.cssViewStyle(builder, table.Session(), self)

	gap, ok := sizeProperty(table, Gap, table.Session())
	if !ok || gap.Type == Auto || gap.Value <= 0 {
		builder.add("border-spacing", "0")
		builder.add("border-collapse", "collapse")
	} else {
		builder.add("border-spacing", gap.cssString("0"))
		builder.add("border-collapse", "separate")
	}
}

func (table *tableViewData) ReloadTableData() {
	updateInnerHTML(table.htmlID(), table.Session())
}

func (cell *tableCellView) Set(tag string, value interface{}) bool {
	return cell.set(strings.ToLower(tag), value)
}

func (cell *tableCellView) set(tag string, value interface{}) bool {
	switch tag {
	case VerticalAlign:
		tag = TableVerticalAlign
	}
	return cell.viewData.set(tag, value)
}

func (cell *tableCellView) cssStyle(self View, builder cssBuilder) {
	session := cell.Session()
	cell.viewData.cssViewStyle(builder, session, self)

	if value, ok := enumProperty(cell, TableVerticalAlign, session, 0); ok {
		builder.add("vertical-align", enumProperties[TableVerticalAlign].values[value])
	}
}
