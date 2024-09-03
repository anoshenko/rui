package rui

const (
	// ID is the constant for the "id" property tag.
	// The "id" property is an optional textual identifier for the View.
	ID = "id"

	// Style is the constant for the "style" property tag.
	// The string "style" property sets the name of the style that is applied to the View when the "disabled" property is set to false
	// or "style-disabled" property is not defined.
	Style = "style"

	// StyleDisabled is the constant for the "style-disabled" property tag.
	// The string "style-disabled" property sets the name of the style that is applied to the View when the "disabled" property is set to true.
	StyleDisabled = "style-disabled"

	// Disabled is the constant for the "disabled" property tag.
	// The bool "disabled" property allows/denies the View to receive focus.
	Disabled = "disabled"

	// Focusable is the constant for the "disabled" property tag.
	// The bool "focusable" determines whether the view will receive focus.
	Focusable = "focusable"

	// Semantics is the constant for the "semantics" property tag.
	// The "semantics" property defines the semantic meaning of the View.
	// This property may have no visible effect, but it allows search engines to understand the structure of your application.
	// It also helps to voice the interface to systems for people with disabilities.
	Semantics = "semantics"

	// Visibility is the constant for the "visibility" property tag.
	// The "visibility" int property specifies the visibility of the View. Valid values are
	// * Visible (0) - the View is visible (default value);
	// * Invisible (1) - the View is invisible but takes up space;
	// * Gone (2) - the View is invisible and does not take up space.
	Visibility = "visibility"

	// ZIndex is the constant for the "z-index" property tag.
	// The int "z-index" property sets the z-order of a positioned view.
	// Overlapping views with a larger z-index cover those with a smaller one.
	ZIndex = "z-index"

	// Opacity is the constant for the "opacity" property tag.
	// The float "opacity" property in [1..0] range sets the opacity of an element.
	// Opacity is the degree to which content behind an element is hidden, and is the opposite of transparency.
	Opacity = "opacity"

	// Overflow is the constant for the "overflow" property tag.
	// The "overflow" int property sets the desired behavior for an element's overflow — i.e.
	// when an element's content is too big to fit in its block formatting context — in both directions.
	// Valid values: OverflowHidden (0), OverflowVisible (1), OverflowScroll (2), OverflowAuto (3)
	Overflow = "overflow"

	// Row is the constant for the "row" property tag.
	Row = "row"

	// Column is the constant for the "column" property tag.
	Column = "column"

	// Left is the constant for the "left" property tag.
	// The "left" SizeUnit property participates in specifying the left border position of a positioned view.
	// Used only for views placed in an AbsoluteLayout.
	Left = "left"

	// Right is the constant for the "right" property tag.
	// The "right" SizeUnit property participates in specifying the right border position of a positioned view.
	// Used only for views placed in an AbsoluteLayout.
	Right = "right"

	// Top is the constant for the "top" property tag.
	// The "top" SizeUnit property participates in specifying the top border position of a positioned view.
	// Used only for views placed in an AbsoluteLayout.
	Top = "top"

	// Bottom is the constant for the "bottom" property tag.
	// The "bottom" SizeUnit property participates in specifying the bottom border position of a positioned view.
	// Used only for views placed in an AbsoluteLayout.
	Bottom = "bottom"

	// Width is the constant for the "width" property tag.
	// The "width" SizeUnit property sets an view's width.
	Width = "width"

	// Height is the constant for the "height" property tag.
	// The "height" SizeUnit property sets an view's height.
	Height = "height"

	// MinWidth is the constant for the "min-width" property tag.
	// The "width" SizeUnit property sets an view's minimal width.
	MinWidth = "min-width"

	// MinHeight is the constant for the "min-height" property tag.
	// The "height" SizeUnit property sets an view's minimal height.
	MinHeight = "min-height"

	// MaxWidth is the constant for the "max-width" property tag.
	// The "width" SizeUnit property sets an view's maximal width.
	MaxWidth = "max-width"

	// MaxHeight is the constant for the "max-height" property tag.
	// The "height" SizeUnit property sets an view's maximal height.
	MaxHeight = "max-height"

	// Margin is the constant for the "margin" property tag.
	// The "margin" property sets the margin area on all four sides of an element.
	// ...
	Margin = "margin"

	// MarginLeft is the constant for the "margin-left" property tag.
	// The "margin-left" SizeUnit property sets the margin area on the left of a view.
	// A positive value places it farther from its neighbors, while a negative value places it closer.
	MarginLeft = "margin-left"

	// MarginRight is the constant for the "margin-right" property tag.
	// The "margin-right" SizeUnit property sets the margin area on the right of a view.
	// A positive value places it farther from its neighbors, while a negative value places it closer.
	MarginRight = "margin-right"

	// MarginTop is the constant for the "margin-top" property tag.
	// The "margin-top" SizeUnit property sets the margin area on the top of a view.
	// A positive value places it farther from its neighbors, while a negative value places it closer.
	MarginTop = "margin-top"

	// MarginBottom is the constant for the "margin-bottom" property tag.
	// The "margin-bottom" SizeUnit property sets the margin area on the bottom of a view.
	// A positive value places it farther from its neighbors, while a negative value places it closer.
	MarginBottom = "margin-bottom"

	// Padding is the constant for the "padding" property tag.
	// The "padding" Bounds property sets the padding area on all four sides of a view at once.
	// An element's padding area is the space between its content and its border.
	Padding = "padding"

	// PaddingLeft is the constant for the "padding-left" property tag.
	// The "padding-left" SizeUnit property sets the width of the padding area to the left of a view.
	PaddingLeft = "padding-left"

	// PaddingRight is the constant for the "padding-right" property tag.
	// The "padding-right" SizeUnit property sets the width of the padding area to the right of a view.
	PaddingRight = "padding-right"

	// PaddingTop is the constant for the "padding-top" property tag.
	// The "padding-top" SizeUnit property sets the height of the padding area to the top of a view.
	PaddingTop = "padding-top"

	// PaddingBottom is the constant for the "padding-bottom" property tag.
	// The "padding-bottom" SizeUnit property sets the height of the padding area to the bottom of a view.
	PaddingBottom = "padding-bottom"

	// AccentColor is the constant for the "accent-color" property tag.
	// The "accent-color" property sets sets the accent color for UI controls generated by some elements.
	AccentColor = "accent-color"

	// BackgroundColor is the constant for the "background-color" property tag.
	// The "background-color" property sets the background color of a view.
	BackgroundColor = "background-color"

	// Background is the constant for the "background" property tag.
	// The "background" property sets one or more background images and/or gradients on a view.
	// ...
	Background = "background"

	// Cursor is the constant for the "cursor" property tag.
	// The "cursor" int property sets the type of mouse cursor, if any, to show when the mouse pointer is over a view
	// Valid values are "auto" (0), "default" (1), "none" (2), "context-menu" (3), "help" (4), "pointer" (5),
	// "progress" (6), "wait" (7), "cell" (8), "crosshair" (9), "text" (10), "vertical-text" (11), "alias" (12),
	// "copy" (13), "move" (14), "no-drop" (15), "not-allowed" (16), "e-resize" (17), "n-resize" (18),
	// "ne-resize" (19), "nw-resize" (20), "s-resize" (21), "se-resize" (22), "sw-resize" (23), "w-resize" (24),
	// "ew-resize" (25), "ns-resize" (26), "nesw-resize" (27), "nwse-resize" (28), "col-resize" (29),
	// "row-resize" (30), "all-scroll" (31), "zoom-in" (32), "zoom-out" (33), "grab" (34), "grabbing" (35).
	Cursor = "cursor"

	// Border is the constant for the "border" property tag.
	// The "border" property sets a view's border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	Border = "border"

	// BorderLeft is the constant for the "border-left" property tag.
	// The "border-left" property sets a view's left border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	BorderLeft = "border-left"

	// BorderRight is the constant for the "border-right" property tag.
	// The "border-right" property sets a view's right border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	BorderRight = "border-right"

	// BorderTop is the constant for the "border-top" property tag.
	// The "border-top" property sets a view's top border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	BorderTop = "border-top"

	// BorderBottom is the constant for the "border-bottom" property tag.
	// The "border-bottom" property sets a view's bottom border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	BorderBottom = "border-bottom"

	// BorderStyle is the constant for the "border-style" property tag.
	// The "border-style" property sets the line style for all four sides of a view's border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	BorderStyle = "border-style"

	// BorderLeftStyle is the constant for the "border-left-style" property tag.
	// The "border-left-style" int property sets the line style of a view's left border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	BorderLeftStyle = "border-left-style"

	// BorderRightStyle is the constant for the "border-right-style" property tag.
	// The "border-right-style" int property sets the line style of a view's right border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	BorderRightStyle = "border-right-style"

	// BorderTopStyle is the constant for the "border-top-style" property tag.
	// The "border-top-style" int property sets the line style of a view's top border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	BorderTopStyle = "border-top-style"

	// BorderBottomStyle is the constant for the "border-bottom-style" property tag.
	// The "border-bottom-style" int property sets the line style of a view's bottom border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	BorderBottomStyle = "border-bottom-style"

	// BorderWidth is the constant for the "border-width" property tag.
	// The "border-width" property sets the line width for all four sides of a view's border.
	BorderWidth = "border-width"

	// BorderLeftWidth is the constant for the "border-left-width" property tag.
	// The "border-left-width" SizeUnit property sets the line width of a view's left border.
	BorderLeftWidth = "border-left-width"

	// BorderRightWidth is the constant for the "border-right-width" property tag.
	// The "border-right-width" SizeUnit property sets the line width of a view's right border.
	BorderRightWidth = "border-right-width"

	// BorderTopWidth is the constant for the "border-top-width" property tag.
	// The "border-top-width" SizeUnit property sets the line width of a view's top border.
	BorderTopWidth = "border-top-width"

	// BorderBottomWidth is the constant for the "border-bottom-width" property tag.
	// The "border-bottom-width" SizeUnit property sets the line width of a view's bottom border.
	BorderBottomWidth = "border-bottom-width"

	// BorderColor is the constant for the "border-color" property tag.
	// The "border-color" property sets the line color for all four sides of a view's border.
	BorderColor = "border-color"

	// BorderLeftColor is the constant for the "border-left-color" property tag.
	// The "border-left-color" property sets the line color of a view's left border.
	BorderLeftColor = "border-left-color"

	// BorderRightColor is the constant for the "border-right-color" property tag.
	// The "border-right-color" property sets the line color of a view's right border.
	BorderRightColor = "border-right-color"

	// BorderTopColor is the constant for the "border-top-color" property tag.
	// The "border-top-color" property sets the line color of a view's top border.
	BorderTopColor = "border-top-color"

	// BorderBottomColor is the constant for the "border-bottom-color" property tag.
	// The "border-bottom-color" property sets the line color of a view's bottom border.
	BorderBottomColor = "border-bottom-color"

	// Outline is the constant for the "outline" property tag.
	// The "border" property sets a view's outline. It sets the values of an outline width, style, and color.
	Outline = "outline"

	// OutlineStyle is the constant for the "outline-style" property tag.
	// The "outline-style" int property sets the style of an view's outline.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	OutlineStyle = "outline-style"

	// OutlineColor is the constant for the "outline-color" property tag.
	// The "outline-color" property sets the color of an view's outline.
	OutlineColor = "outline-color"

	// OutlineWidth is the constant for the "outline-width" property tag.
	// The "outline-width" SizeUnit property sets the width of an view's outline.
	OutlineWidth = "outline-width"

	// OutlineWidth is the constant for the "outline-offset" property tag.
	// The "outline-offset" SizeUnit property sets the amount of space between an outline and the edge or border of an element..
	OutlineOffset = "outline-offset"

	// Shadow is the constant for the "shadow" property tag.
	// The "shadow" property adds shadow effects around a view's frame. A shadow is described
	// by X and Y offsets relative to the element, blur and spread radius, and color.
	// ...
	Shadow = "shadow"

	// FontName is the constant for the "font-name" property tag.
	// The "font-name" string property specifies a prioritized list of one or more font family names and/or
	// generic family names for the selected view. Values are separated by commas to indicate that they are alternatives.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	FontName = "font-name"

	// TextColor is the constant for the "text-color" property tag.
	// The "color" property sets the foreground color value of a view's text and text decorations.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextColor = "text-color"

	// TextSize is the constant for the "text-size" property tag.
	// The "text-size" SizeUnit property sets the size of the font.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextSize = "text-size"

	// Italic is the constant for the "italic" property tag.
	// The "italic" is the bool property. If it is "true" then a text is displayed in italics.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	Italic = "italic"

	// SmallCaps is the constant for the "small-caps" property tag.
	// The "small-caps" is the bool property. If it is "true" then a text is displayed in small caps.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	SmallCaps = "small-caps"

	// Strikethrough is the constant for the "strikethrough" property tag.
	// The "strikethrough" is the bool property. If it is "true" then a text is displayed strikethrough.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	Strikethrough = "strikethrough"

	// Overline is the constant for the "overline" property tag.
	// The "overline" is the bool property. If it is "true" then a text is displayed overlined.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	Overline = "overline"

	// Underline is the constant for the "underline" property tag.
	// The "underline" is the bool property. If it is "true" then a text is displayed underlined.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	Underline = "underline"

	// TextLineThickness is the constant for the "text-line-thickness" property tag.
	// The "text-line-thickness" SizeUnit property sets the stroke thickness of the decoration line that
	// is used on text in an element, such as a line-through, underline, or overline.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextLineThickness = "text-line-thickness"

	// TextLineStyle is the constant for the "text-line-style" property tag.
	// The "text-line-style" int property sets the style of the lines specified by "text-decoration" property.
	// The style applies to all lines that are set with "text-decoration" property.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextLineStyle = "text-line-style"

	// TextLineColor is the constant for the "text-line-color" property tag.
	// The "text-line-color" Color property sets the color of the lines specified by "text-decoration" property.
	// The color applies to all lines that are set with "text-decoration" property.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextLineColor = "text-line-color"

	// TextWeight is the constant for the "text-weight" property tag.
	// Valid values are SolidLine (1), DashedLine (2), DottedLine (3), DoubleLine (4) and WavyLine (5).
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextWeight = "text-weight"

	// TextAlign is the constant for the "text-align" property tag.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextAlign = "text-align"

	// TextIndent is the constant for the "text-indent" property tag.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextIndent = "text-indent"

	// TextShadow is the constant for the "text-shadow" property tag.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextShadow = "text-shadow"

	// TextWrap is the constant for the "text-wrap" property tag.
	// The "text-wrap" int property controls how text inside a View is wrapped.
	// Valid values ​​are:
	// * TextWrapOn / 0 / "wrap" - text is wrapped across lines at appropriate characters (for example spaces, in languages like English that use space separators) to minimize overflow. This is the default value.
	// * TextWrapOff / 1 / "nowrap" - text does not wrap across lines. It will overflow its containing element rather than breaking onto a new line.
	// * TextWrapBalance / 2 / "balance" - text is wrapped in a way that best balances the number of characters on each line, enhancing layout quality and legibility. Because counting characters and balancing them across multiple lines is computationally expensive, this value is only supported for blocks of text spanning a limited number of lines (six or less for Chromium and ten or less for Firefox).
	TextWrap = "text-wrap"

	// TabSize is the constant for the "tab-size" property tag.
	// The "tab-size" int property sets the width of tab characters (U+0009) in spaces.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TabSize = "tab-size"

	// LetterSpacing is the constant for the "letter-spacing" property tag.
	// The "letter-spacing" SizeUnit property sets the horizontal spacing behavior between text characters.
	// This value is added to the natural spacing between characters while rendering the text.
	// Positive values of letter-spacing causes characters to spread farther apart,
	// while negative values of letter-spacing bring characters closer together.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	LetterSpacing = "letter-spacing"

	// WordSpacing is the constant for the "word-spacing" property tag.
	// The "word-spacing" SizeUnit property sets the length of space between words and between tags.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	WordSpacing = "word-spacing"

	// LineHeight is the constant for the "line-height" property tag.
	// The "line-height" SizeUnit property sets the height of a line box.
	// It's commonly used to set the distance between lines of text.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	LineHeight = "line-height"

	// WhiteSpace is the constant for the "white-space" property tag.
	// The "white-space" int property sets how white space inside an element is handled.
	// Valid values are WhiteSpaceNormal (0),  WhiteSpaceNowrap (1), WhiteSpacePre (2),
	// WhiteSpacePreWrap (3), WhiteSpacePreLine (4), WhiteSpaceBreakSpaces (5)
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	WhiteSpace = "white-space"

	// WordBreak is the constant for the "word-break" property tag.
	// The "word-break" int property sets whether line breaks appear wherever the text would otherwise overflow its content box.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	WordBreak = "word-break"

	// TextTransform is the constant for the "text-transform" property tag.
	// The "text-transform" int property specifies how to capitalize an element's text.
	// It can be used to make text appear in all-uppercase or all-lowercase, or with each word capitalized.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextTransform = "text-transform"

	// TextDirection is the constant for the "text-direction" property tag.
	// The "text-direction" int property sets the direction of text, table columns, and horizontal overflow.
	// Use 1 (LeftToRightDirection) for languages written from right to left (like Hebrew or Arabic),
	// and 2 (RightToLeftDirection) for those written from left to right (like English and most other languages).
	// The default value of the property is 0 (SystemTextDirection): use the system text direction.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	TextDirection = "text-direction"

	// WritingMode is the constant for the "writing-mode" property tag.
	// The "writing-mode" int property sets whether lines of text are laid out horizontally or vertically,
	// as well as the direction in which blocks progress
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	WritingMode = "writing-mode"

	// VerticalTextOrientation is the constant for the "vertical-text-orientation" property tag.
	// The "vertical-text-orientation" int property sets the orientation of the text characters in a line.
	// It only affects text in vertical mode ("writing-mode" property).
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	VerticalTextOrientation = "vertical-text-orientation"

	// TextOverflow is the constant for the "text-overflow" property tag.
	// The "text-overflow" int property sets how hidden overflow content is signaled to users.
	// It can be clipped or display an ellipsis ('…'). Valid values are
	TextOverflow = "text-overflow"

	// Hint is the constant for the "hint" property tag.
	// The "hint" string property sets a hint to the user of what can be entered in the control.
	Hint = "hint"

	// MaxLength is the constant for the "max-length" property tag.
	// The "max-length" int property sets the maximum number of characters that the user can enter
	MaxLength = "max-length"

	// ReadOnly is the constant for the "readonly" property tag.
	// This bool property indicates that the user cannot modify the value of the EditView.
	ReadOnly = "readonly"

	// Content is the constant for the "content" property tag.
	Content = "content"

	// Items is the constant for the "items" property tag.
	Items = "items"

	// DisabledItems is the constant for the "disabled-items" property tag.
	// The "disabled-items" []int property specifies an array of disabled(non selectable) items indices of DropDownList.
	DisabledItems = "disabled-items"

	// ItemSeparators is the constant for the "item-separators" property tag.
	// The "item-separators" []int property specifies an array of indices of DropDownList items after which a separator should be added.
	ItemSeparators = "item-separators"

	// Current is the constant for the "current" property tag.
	Current = "current"

	// Type is the constant for the "type" property tag.
	Type = "type"

	// Pattern is the constant for the "pattern" property tag.
	Pattern = "pattern"

	// GridAutoFlow is the constant for the "grid-auto-flow" property tag.
	// The "grid-auto-flow" int property controls how the GridLayout auto-placement algorithm works,
	// specifying exactly how auto-placed items get flowed into the grid.
	// Valid values are RowAutoFlow (0), ColumnAutoFlow (1), RowDenseAutoFlow (2), and ColumnDenseAutoFlow (3)
	GridAutoFlow = "grid-auto-flow"

	// CellWidth is the constant for the "cell-width" property tag.
	// The "cell-width" properties allow to set a fixed width of GridLayout cells regardless of the size of the child elements.
	// These properties are of type []SizeUnit. Each element in the array determines the size of the corresponding column.
	CellWidth = "cell-width"

	// CellHeight is the constant for the "cell-height" property tag.
	// The "cell-height" properties allow to set a fixed height of GridLayout cells regardless of the size of the child elements.
	// These properties are of type []SizeUnit. Each element in the array determines the size of the corresponding row.
	CellHeight = "cell-height"

	// GridRowGap is the constant for the "grid-row-gap" property tag.
	// The "grid-row-gap" SizeUnit properties allow to set the distance between the rows of the GridLayout container.
	// The default is 0px.
	GridRowGap = "grid-row-gap"

	// GridColumnGap is the constant for the "grid-column-gap" property tag.
	// The "grid-column-gap" SizeUnit properties allow to set the distance between the columns of the GridLayout container.
	// The default is 0px.
	GridColumnGap = "grid-column-gap"

	// Source is the constant for the "src" property tag.
	// The "src" property specifies the image to display in the ImageView.
	Source = "src"

	// SrcSet is the constant for the "srcset" property tag.
	// The "srcset" property is a string which identifies one or more image candidate strings, separated using commas (,)
	// each specifying image resources to use under given screen density.
	// This property is only used if you are building an application for js/wasm platform
	SrcSet = "srcset"

	// Fit is the constant for the "fit" property tag.
	Fit           = "fit"
	backgroundFit = "background-fit"

	// Repeat is the constant for the "repeat" property tag.
	Repeat = "repeat"

	// Attachment is the constant for the "attachment" property tag.
	Attachment = "attachment"

	// BackgroundClip is the constant for the "background-clip" property tag.
	BackgroundClip = "background-clip"

	// Gradient is the constant for the "gradient" property tag.
	Gradient = "gradient"

	// Direction is the constant for the "direction" property tag.
	Direction = "direction"

	// Repeating is the constant for the "repeating" property tag.
	Repeating = "repeating"

	// Repeating is the constant for the "repeating" property tag.
	From = "from"

	// RadialGradientRadius is the constant for the "radial-gradient-radius" property tag.
	RadialGradientRadius = "radial-gradient-radius"

	// RadialGradientShape is the constant for the "radial-gradient-shape" property tag.
	RadialGradientShape = "radial-gradient-shape"

	// Shape is the constant for the "shape" property tag. It's a short form of "radial-gradient-shape"
	Shape = "shape"

	// CenterX is the constant for the "center-x" property tag.
	CenterX = "center-x"

	// CenterY is the constant for the "center-x" property tag.
	CenterY = "center-y"

	// AltText is the constant for the "alt-text" property tag.
	AltText = "alt-text"
	altTag  = "alt"

	// AvoidBreak is the constant for the "avoid-break" property tag.
	// The "avoid-break" bool property sets how region breaks should behave inside a generated box.
	// If the property value is "true" then avoids any break from being inserted within the principal box.
	// If the property value is "false" then allows, but does not force, any break to be inserted within
	// the principal box.
	AvoidBreak = "avoid-break"

	// ItemWidth is the constant for the "item-width" property tag.
	ItemWidth = "item-width"

	// ItemHeight is the constant for the "item-height" property tag.
	ItemHeight = "item-height"

	// ListWrap is the constant for the "wrap" property tag.
	ListWrap = "list-wrap"

	// EditWrap is the constant for the "wrap" property tag.
	EditWrap = "edit-wrap"

	// CaretColor is the constant for the "caret-color" property tag.
	// The "caret-color" Color property sets the color of the insertion caret, the visible marker
	// where the next character typed will be inserted. This is sometimes referred to as the text input cursor.
	CaretColor = "caret-color"

	// Min is the constant for the "min" property tag.
	Min = "min"

	// Max is the constant for the "max" property tag.
	Max = "max"

	// Step is the constant for the "step" property tag.
	Step = "step"

	// Value is the constant for the "value" property tag.
	Value = "value"

	// Orientation is the constant for the "orientation" property tag.
	Orientation = "orientation"

	// Gap is t he constant for the "gap" property tag.
	Gap = "gap"

	// ListRowGap is the constant for the "list-row-gap" property tag.
	// The "list-row-gap" SizeUnit properties allow to set the distance between the rows of the ListLayout or ListView.
	// The default is 0px.
	ListRowGap = "list-row-gap"

	// ListColumnGap is the constant for the "list-column-gap" property tag.
	// The "list-column-gap" SizeUnit properties allow to set the distance between the columns of the GridLayout or ListView.
	// The default is 0px.
	ListColumnGap = "list-column-gap"

	// Text is the constant for the "text" property tag.
	Text = "text"

	// VerticalAlign is the constant for the "vertical-align" property tag.
	VerticalAlign = "vertical-align"

	// HorizontalAlign is the constant for the "horizontal-align" property tag.
	// The "horizontal-align" int property sets the horizontal alignment of the content inside a block element
	HorizontalAlign = "horizontal-align"

	// ImageVerticalAlign is the constant for the "image-vertical-align" property tag.
	ImageVerticalAlign = "image-vertical-align"

	// ImageHorizontalAlign is the constant for the "image-horizontal-align" property tag.
	ImageHorizontalAlign = "image-horizontal-align"

	// Checked is the constant for the "checked" property tag.
	Checked = "checked"

	// ItemVerticalAlign is the constant for the "item-vertical-align" property tag.
	ItemVerticalAlign = "item-vertical-align"

	// ItemHorizontalAlign is the constant for the "item-horizontal-align" property tag.
	ItemHorizontalAlign = "item-horizontal-align"

	// ItemCheckbox is the constant for the "checkbox" property tag.
	ItemCheckbox = "checkbox"

	// CheckboxHorizontalAlign is the constant for the "checkbox-horizontal-align" property tag.
	CheckboxHorizontalAlign = "checkbox-horizontal-align"

	// CheckboxVerticalAlign is the constant for the "checkbox-vertical-align" property tag.
	CheckboxVerticalAlign = "checkbox-vertical-align"

	// NotTranslate is the constant for the "not-translate" property tag.
	// This bool property indicates that no need to translate the text.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	NotTranslate = "not-translate"

	// Filter is the constant for the "filter" property tag.
	// The "filter" property applies graphical effects to a View,
	// such as such as blurring, color shifting, changing brightness/contrast, etc.
	Filter = "filter"

	// BackdropFilter is the constant for the "backdrop-filter" property tag.
	// The "backdrop-filter" property applies graphical effects to the area behind a View,
	// such as such as blurring, color shifting, changing brightness/contrast, etc.
	BackdropFilter = "backdrop-filter"

	// Clip is the constant for the "clip" property tag.
	// The "clip" property creates a clipping region that sets what part of a View should be shown.
	Clip = "clip"

	// Points is the constant for the "points" property tag.
	Points = "points"

	// ShapeOutside is the constant for the "shape-outside" property tag.
	// The "shape-outside" property defines a shape (which may be non-rectangular) around which adjacent
	// inline content should wrap. By default, inline content wraps around its margin box;
	// "shape-outside" provides a way to customize this wrapping, making it possible to wrap text around
	// complex objects rather than simple boxes.
	ShapeOutside = "shape-outside"

	// Float is the constant for the "float" property tag.
	// The "float" property places a View on the left or right side of its container,
	// allowing text and inline Views to wrap around it.
	Float = "float"

	// UserData is the constant for the "user-data" property tag.
	// The "user-data" property can contain any user data
	UserData = "user-data"

	// Resize is the constant for the "resize" property tag.
	// The "resize" int property sets whether an element is resizable, and if so, in which directions.
	// Valid values are "none" / NoneResize (0), "both" / BothResize (1),
	// "horizontal" / HorizontalResize (2), and "vertical" / VerticalResize (3)
	Resize = "resize"

	// UserSelect is the constant for the "user-select" property tag.
	// The "user-select" bool property controls whether the user can select text.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	UserSelect = "user-select"

	// Order is the constant for the "Order" property tag.
	// The "Order" int property sets the order to layout an item in a ListLayout or GridLayout container.
	// Items in a container are sorted by ascending order value and then by their source code order.
	Order = "Order"

	// BackgroundBlendMode is the constant for the "background-blend-mode" property tag.
	// The "background-blend-mode" int property sets how an view's background images should blend
	// with each other and with the view's background color.
	// Valid values are "normal" (0), "multiply" (1), "screen" (2), "overlay" (3), "darken" (4), "lighten" (5),
	// "color-dodge" (6), "color-burn" (7), "hard-light" (8), "soft-light" (9), "difference" (10),
	// "exclusion" (11), "hue" (12), "saturation" (13), "color" (14), "luminosity" (15).
	BackgroundBlendMode = "background-blend-mode"

	// MixBlendMode is the constant for the "mix-blend-mode" property tag.
	// The "mix-blend-mode" int property sets how a view's content should blend
	// with the content of the view's parent and the view's background.
	// Valid values are "normal" (0), "multiply" (1), "screen" (2), "overlay" (3), "darken" (4), "lighten" (5),
	// "color-dodge" (6), "color-burn" (7), "hard-light" (8), "soft-light" (9), "difference" (10),
	// "exclusion" (11), "hue" (12), "saturation" (13), "color" (14), "luminosity" (15).
	MixBlendMode = "mix-blend-mode"

	// TabIndex is the constant for the "tabindex" property tag.
	// The "tabindex" int property indicates that View can be focused, and where it participates in sequential keyboard navigation
	// (usually with the Tab key, hence the name).
	// * A negative value means that View is not reachable via sequential keyboard navigation, but could be focused by clicking with the mouse or touching.
	// * tabindex="0" means that View should be focusable in sequential keyboard navigation, after any positive tabindex values and its order is defined in order of its addition.
	// * A positive value means View should be focusable in sequential keyboard navigation, with its order defined by the value of the number.
	TabIndex = "tabindex"

	Tooltip = "tooltip"
)
