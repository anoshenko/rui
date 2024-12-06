package rui

type PropertyName string

// Constants for various properties and events of Views'.
const (
	// ID is the constant for "id" property tag.
	//
	// # Used by View, Animation.
	//
	// Usage in View:
	// Optional textual identifier for the view. Used to reference view from source code if needed.
	//
	// Supported types: string.
	//
	// # Usage in Animation:
	//
	// Specifies the animation identifier. Used only for animation script.
	//
	// Supported types: string.
	ID PropertyName = "id"

	// Style is the constant for "style" property tag.
	//
	// Used by ColumnSeparatorProperty, View, BorderProperty, OutlineProperty.
	//
	// # Usage in ColumnSeparatorProperty:
	//
	// Line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The separator will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a separator.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a separator.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a separator.
	//   - 4 (DoubleLine) or "double" - Double line as a separator.
	//
	// # Usage in View:
	//
	// Sets the name of the style that is applied to the view when the "disabled" property is set to false or "style-disabled"
	// property is not defined.
	//
	// Supported types: string.
	//
	// # Usage in BorderProperty:
	//
	// Border line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The border will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a border.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//   - 4 (DoubleLine) or "double" - Double line as a border.
	//
	// # Usage in OutlineProperty:
	//
	// Outline line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The outline will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as an outline.
	//   - 2 (DashedLine) or "dashed" - Dashed line as an outline.
	//   - 3 (DottedLine) or "dotted" - Dotted line as an outline.
	//   - 4 (DoubleLine) or "double" - Double line as an outline.
	Style PropertyName = "style"

	// StyleDisabled is the constant for "style-disabled" property tag.
	//
	// Used by View.
	// Sets the name of the style that is applied to the view when the "disabled" property is set to true.
	//
	// Supported types: string.
	StyleDisabled PropertyName = "style-disabled"

	// Disabled is the constant for "disabled" property tag.
	//
	// Used by ViewsContainer.
	// Controls whether the view can receive focus and which style to use. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - View can't receive focus and "style-disabled" style will be used by the view.
	//   - false, 0, "false", "no", "off", or "0" - View can receive focus and "style" style will be used by the view.
	Disabled PropertyName = "disabled"

	// Focusable is the constant for "focusable" property tag.
	//
	// Used by View.
	// Controls whether view can receive focus.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - View can have a focus.
	//   - false, 0, "false", "no", "off", or "0" - View can't have a focus.
	Focusable PropertyName = "focusable"

	// Semantics is the constant for "semantics" property tag.
	//
	// Used by View.
	// Defines the semantic meaning of the view. This property may have no visible effect, but it allows search engines to
	// understand the structure of your application. It also helps to voice the interface to systems for people with
	// disabilities.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (DefaultSemantics) or "default" - Default semantics.
	//   - 1 (ArticleSemantics) or "article" - Article semantics.
	//   - 2 (SectionSemantics) or "section" - Section semantics.
	//   - 3 (AsideSemantics) or "aside" - Aside semantics.
	//   - 4 (HeaderSemantics) or "header" - Header semantics.
	//   - 5 (MainSemantics) or "main" - Main semantics.
	//   - 6 (FooterSemantics) or "footer" - Footer semantics.
	//   - 7 (NavigationSemantics) or "navigation" - Navigation semantics.
	//   - 8 (FigureSemantics) or "figure" - Figure semantics.
	//   - 9 (FigureCaptionSemantics) or "figure-caption" - Figure caption semantics.
	//   - 10 (ButtonSemantics) or "button" - Button semantics.
	//   - 11 (ParagraphSemantics) or "p" - Paragraph semantics.
	//   - 12 (H1Semantics) or "h1" - Heading level 1 semantics.
	//   - 13 (H2Semantics) or "h2" - Heading level 2 semantics.
	//   - 14 (H3Semantics) or "h3" - Heading level 3 semantics.
	//   - 15 (H4Semantics) or "h4" - Heading level 4 semantics.
	//   - 16 (H5Semantics) or "h5" - Heading level 5 semantics.
	//   - 17 (H6Semantics) or "h6" - Heading level 6 semantics.
	//   - 18 (BlockquoteSemantics) or "blockquote" - Blockquote semantics.
	//   - 19 (CodeSemantics) or "code" - Code semantics.
	Semantics PropertyName = "semantics"

	// Visibility is the constant for "visibility" property tag.
	//
	// Used by View.
	// Specifies the visibility of the view.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (Visible) or "visible" - The view is visible.
	//   - 1 (Invisible) or "invisible" - The view is invisible but takes up space.
	//   - 2 (Gone) or "gone" - The view is invisible and does not take up space.
	Visibility PropertyName = "visibility"

	// ZIndex is the constant for "z-index" property tag.
	//
	// Used by View.
	// Sets the z-order of a positioned view. Overlapping views with a larger z-index cover those with a smaller one.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - negative value - Views with lower value will be behind views with higher value.
	//   - not negative value - Views with higher value will be on top of views with lower value.
	ZIndex PropertyName = "z-index"

	// Opacity is the constant for "opacity" property tag.
	//
	// Used by View, ViewFilter.
	//
	// # Usage in View:
	//
	// In [1..0] range sets the opacity of view. Opacity is the degree to which content behind the view is hidden, and is the
	// opposite of transparency.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	//
	// # Usage in ViewFilter:
	//
	// Opacity is the degree to which content behind the view is hidden, and is the opposite of transparency. Value is in
	// range 0% to 100%, where 0% is fully transparent.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Opacity PropertyName = "opacity"

	// Overflow is the constant for "overflow" property tag.
	//
	// Used by View.
	// Set the desired behavior for an element's overflow i.e. when an element's content is too big to fit in its block
	// formatting context in both directions.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (OverflowHidden) or "hidden" - The overflow is clipped, and the rest of the content will be invisible.
	//   - 1 (OverflowVisible) or "visible" - The overflow is not clipped. The content renders outside the element's box.
	//   - 2 (OverflowScroll) or "scroll" - The overflow is clipped, and a scrollbar is added to see the rest of the content.
	//   - 3 (OverflowAuto) or "auto" - Similar to OverflowScroll, but it adds scrollbars only when necessary.
	Overflow PropertyName = "overflow"

	// Row is the constant for "row" property tag.
	//
	// Used by View.
	// Row of the view inside the container like GridLayout.
	//
	// Supported types: Range, int, string.
	//
	// Internal type is Range, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - int - set single value(index).
	//   - string - can contain single integer value(index) or a range of integer values(indices), examples: "0", "0:3".
	Row PropertyName = "row"

	// Column is the constant for "column" property tag.
	//
	// Used by View.
	// Column of the view inside the container like GridLayout.
	//
	// Supported types: Range, int, string.
	//
	// Internal type is Range, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - int - set single value(index).
	//   - string - can contain single integer value(index) or a range of integer values(indices), examples: "0", "0:3".
	Column PropertyName = "column"

	// Left is the constant for "left" property tag.
	//
	// Used by View, BoundsProperty, ClipShapeProperty.
	//
	// # Usage in View:
	//
	// Offset from left border of the container. Used only for views placed in an AbsoluteLayout.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in BoundsProperty:
	//
	// Left bound value.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in ClipShapeProperty:
	//
	// Specifies the left border position of inset clip shape.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Left PropertyName = "left"

	// Right is the constant for "right" property tag.
	//
	// Used by View, BoundsProperty, ClipShapeProperty.
	//
	// # Usage in View:
	//
	// Offset from right border of the container. Used only for views placed in an AbsoluteLayout.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in BoundsProperty:
	//
	// Right bound value.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in ClipShapeProperty:
	//
	// Specifies the right border position of inset clip shape.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Right PropertyName = "right"

	// Top is the constant for "top" property tag.
	//
	// Used by View, BoundsProperty, ClipShapeProperty.
	//
	// # Usage in View:
	//
	// Offset from top border of the container. Used only for views placed in an AbsoluteLayout.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in BoundsProperty:
	//
	// Top bound value.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in ClipShapeProperty:
	//
	// Specifies the top border position of inset clip shape.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Top PropertyName = "top"

	// Bottom is the constant for "bottom" property tag.
	//
	// Used by View, BoundsProperty, ClipShapeProperty.
	//
	// # Usage in View:
	//
	// Offset from bottom border of the container. Used only for views placed in an AbsoluteLayout.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in BoundsProperty:
	//
	// Bottom bound value.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in ClipShapeProperty:
	//
	// Specifies the bottom border position of inset clip shape.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Bottom PropertyName = "bottom"

	// Width is the constant for "width" property tag.
	//
	// Used by ColumnSeparatorProperty, View, BorderProperty, OutlineProperty.
	//
	// # Usage in ColumnSeparatorProperty:
	//
	// Line width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in View:
	//
	// Set a view's width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in BorderProperty:
	//
	// Border line width.
	//
	// Supported types: SizeUnit, string.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in OutlineProperty:
	//
	// Outline line width.
	//
	// Supported types: SizeUnit, string.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Width PropertyName = "width"

	// Height is the constant for "height" property tag.
	//
	// Used by View.
	// Set a view's height.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Height PropertyName = "height"

	// MinWidth is the constant for "min-width" property tag.
	//
	// Used by View.
	// Set a view's minimal width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MinWidth PropertyName = "min-width"

	// MinHeight is the constant for "min-height" property tag.
	//
	// Used by View.
	// Set a view's minimal height.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MinHeight PropertyName = "min-height"

	// MaxWidth is the constant for "max-width" property tag.
	//
	// Used by View.
	// Set a view's maximal width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MaxWidth PropertyName = "max-width"

	// MaxHeight is the constant for "max-height" property tag.
	//
	// Used by View.
	// Set a view's maximal height.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MaxHeight PropertyName = "max-height"

	// Margin is the constant for "margin" property tag.
	//
	// Used by View.
	// Set the margin area on all four sides of an element.
	//
	// Supported types: BoundsProperty, Bounds, SizeUnit, float32, float64, int, string.
	//
	// Internal type could be BoundsProperty or SizeUnit depending on whether single value or multiple values has been set, other types converted to them during assignment.
	// See BoundsProperty, Bounds, SizeUnit for more information.
	//
	// Conversion rules:
	//   - BoundsProperty - stored as is, no conversion performed.
	//   - Bounds - new BoundsProperty will be created and corresponding values for top, right, bottom and left border will be set.
	//   - SizeUnit - stored as is and the same value will be used for all borders.
	//   - float - new SizeUnit will be created and the same value(in pixels) will be used for all borders.
	//   - int - new SizeUnit will be created and the same value(in pixels) will be used for all borders.
	//   - string - can contain one or four SizeUnit separated with comma(,). In case one value will be provided a new SizeUnit will be created and the same value will be used for all borders. If four values will be provided then they will be set respectively for top, right, bottom and left border.
	Margin PropertyName = "margin"

	// MarginLeft is the constant for "margin-left" property tag.
	//
	// Used by View.
	// Set the margin area on the left of a view. A positive value places it farther from its neighbors, while a negative
	// value places it closer.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MarginLeft PropertyName = "margin-left"

	// MarginRight is the constant for "margin-right" property tag.
	//
	// Used by View.
	// Set the margin area on the right of a view. A positive value places it farther from its neighbors, while a negative
	// value places it closer.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MarginRight PropertyName = "margin-right"

	// MarginTop is the constant for "margin-top" property tag.
	//
	// Used by View.
	// Set the margin area on the top of a view. A positive value places it farther from its neighbors, while a negative value
	// places it closer.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MarginTop PropertyName = "margin-top"

	// MarginBottom is the constant for "margin-bottom" property tag.
	//
	// Used by View.
	// Set the margin area on the bottom of a view. A positive value places it farther from its neighbors, while a negative
	// value places it closer.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	MarginBottom PropertyName = "margin-bottom"

	// Padding is the constant for "padding" property tag.
	//
	// Used by View.
	// Sets the padding area on all four sides of a view at once. An element's padding area is the space between its content
	// and its border.
	//
	// Supported types: BoundsProperty, Bounds, SizeUnit, float32, float64, int, string.
	//
	// Internal type could be BoundsProperty or SizeUnit depending on whether single value or multiple values has been set, other types converted to them during assignment.
	// See BoundsProperty, Bounds, SizeUnit for more information.
	//
	// Conversion rules:
	//   - BoundsProperty - stored as is, no conversion performed.
	//   - Bounds - new BoundsProperty will be created and corresponding values for top, right, bottom and left border will be set.
	//   - SizeUnit - stored as is and the same value will be used for all borders.
	//   - float - new SizeUnit will be created and the same value(in pixels) will be used for all borders.
	//   - int - new SizeUnit will be created and the same value(in pixels) will be used for all borders.
	//   - string - can contain one or four SizeUnit separated with comma(,). In case one value will be provided a new SizeUnit will be created and the same value will be used for all borders. If four values will be provided then they will be set respectively for top, right, bottom and left border.
	Padding PropertyName = "padding"

	// PaddingLeft is the constant for "padding-left" property tag.
	//
	// Used by View.
	// Set the width of the padding area to the left of a view.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	PaddingLeft PropertyName = "padding-left"

	// PaddingRight is the constant for "padding-right" property tag.
	//
	// Used by View.
	// Set the width of the padding area to the right of a view.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	PaddingRight PropertyName = "padding-right"

	// PaddingTop is the constant for "padding-top" property tag.
	//
	// Used by View.
	// Set the height of the padding area to the top of a view.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	PaddingTop PropertyName = "padding-top"

	// PaddingBottom is the constant for "padding-bottom" property tag.
	//
	// Used by View.
	// Set the height of the padding area to the bottom of a view.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	PaddingBottom PropertyName = "padding-bottom"

	// AccentColor is the constant for "accent-color" property tag.
	//
	// Used by View.
	// Sets the accent color for UI controls generated by some elements.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	AccentColor PropertyName = "accent-color"

	// BackgroundColor is the constant for "background-color" property tag.
	//
	// Used by View.
	// Set the background color of a view.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BackgroundColor PropertyName = "background-color"

	// Background is the constant for "background" property tag.
	//
	// Used by View.
	// Set one or more background images and/or gradients for the view.
	//
	// Supported types: BackgroundElement, []BackgroundElement, string.
	//
	// Internal type is []BackgroundElement, other types converted to it during assignment.
	// See BackgroundElement description for more details.
	//
	// Conversion rules:
	//   - string - must contain text representation of background element(s) like in resource files.
	Background PropertyName = "background"

	// Mask is the constant for "mask" property tag.
	//
	// Used by View.
	// Set one or more images and/or gradients as the view mask.
	// As mask is used only alpha channel of images and/or gradients.
	//
	// Supported types: BackgroundElement, []BackgroundElement, string.
	//
	// Internal type is []BackgroundElement, other types converted to it during assignment.
	// See BackgroundElement description for more details.
	//
	// Conversion rules:
	//   - string - must contain text representation of background element(s) like in resource files.
	Mask PropertyName = "mask"

	// Cursor is the constant for "cursor" property tag.
	//
	// Used by View.
	// Sets the type of mouse cursor, if any, to show when the mouse pointer is over the view.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 or "auto" - Auto cursor.
	//   - 1 or "default" - Default cursor.
	//   - 2 or "none" - None cursor.
	//   - 3 or "context-menu" - Context menu cursor.
	//   - 4 or "help" - Help cursor.
	//   - 5 or "pointer" - Pointer cursor.
	//   - 6 or "progress" - Progress cursor.
	//   - 7 or "wait" - Wait cursor.
	//   - 8 or "cell" - Cell cursor.
	//   - 9 or "crosshair" - Crosshair cursor.
	//   - 10 or "text" - Text cursor.
	//   - 11 or "vertical-text" - Vertical text cursor.
	//   - 12 or "alias" - Alias cursor.
	//   - 13 or "copy" - Copy cursor.
	//   - 14 or "move" - Move cursor.
	//   - 15 or "no-drop" - No drop cursor.
	//   - 16 or "not-allowed" - Not allowed cursor.
	//   - 17 or "e-resize" - Resize cursor.
	//   - 18 or "n-resize" - Resize cursor.
	//   - 19 or "ne-resize" - Resize cursor.
	//   - 20 or "nw-resize" - Resize cursor.
	//   - 21 or "s-resize" - Resize cursor.
	//   - 22 or "se-resize" - Resize cursor.
	//   - 23 or "sw-resize" - Resize cursor.
	//   - 24 or "w-resize" - Resize cursor.
	//   - 25 or "ew-resize" - Resize cursor.
	//   - 26 or "ns-resize" - Resize cursor.
	//   - 27 or "nesw-resize" - Resize cursor.
	//   - 28 or "nwse-resize" - Resize cursor.
	//   - 29 or "col-resize" - Col resize cursor.
	//   - 30 or "row-resize" - Row resize cursor.
	//   - 31 or "all-scroll" - All scroll cursor.
	//   - 32 or "zoom-in" - Zoom in cursor.
	//   - 33 or "zoom-out" - Zoom out cursor.
	//   - 34 or "grab" - Grab cursor.
	//   - 35 or "grabbing" - Grabbing cursor.
	Cursor PropertyName = "cursor"

	// Border is the constant for "border" property tag.
	//
	// Used by View.
	// Set a view's border. It sets the values of a border width, style, and color.
	//
	// Supported types: BorderProperty, ViewBorder, ViewBorders.
	//
	// Internal type is BorderProperty, other types converted to it during assignment.
	// See BorderProperty, ViewBorder, ViewBorders description for more details.
	//
	// Conversion rules:
	//   - ViewBorder - style, width and color applied to all borders and stored in internal implementation of BorderProperty.
	//   - ViewBorders - style, width and color of each border like top, right, bottom and left applied to related borders, stored in internal implementation of BorderProperty.
	Border PropertyName = "border"

	// BorderLeft is the constant for "border-left" property tag.
	//
	// Used by View.
	// Set a view's left border. It sets the values of a border width, style, and color.
	//
	// Supported types: ViewBorder, BorderProperty, string.
	//
	// Internal type is BorderProperty, other types converted to it during assignment.
	// See ViewBorder, BorderProperty description for more details.
	BorderLeft PropertyName = "border-left"

	// BorderRight is the constant for "border-right" property tag.
	//
	// Used by View.
	// Set a view's right border. It sets the values of a border width, style, and color.
	//
	// Supported types: ViewBorder, BorderProperty, string.
	//
	// Internal type is BorderProperty, other types converted to it during assignment.
	// See ViewBorder, BorderProperty description for more details.
	BorderRight PropertyName = "border-right"

	// BorderTop is the constant for "border-top" property tag.
	//
	// Used by View.
	// Set a view's top border. It sets the values of a border width, style, and color.
	//
	// Supported types: ViewBorder, BorderProperty, string.
	//
	// Internal type is BorderProperty, other types converted to it during assignment.
	// See ViewBorder, BorderProperty description for more details.
	BorderTop PropertyName = "border-top"

	// BorderBottom is the constant for "border-bottom" property tag.
	//
	// Used by View.
	// Set a view's bottom border. It sets the values of a border width, style, and color.
	//
	// Supported types: ViewBorder, BorderProperty, string.
	//
	// Internal type is BorderProperty, other types converted to it during assignment.
	// See ViewBorder, BorderProperty description for more details.
	BorderBottom PropertyName = "border-bottom"

	// BorderStyle is the constant for "border-style" property tag.
	//
	// Used by View.
	// Set the line style for all four sides of a view's border.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The border will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a border.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//   - 4 (DoubleLine) or "double" - Double line as a border.
	BorderStyle PropertyName = "border-style"

	// BorderLeftStyle is the constant for "border-left-style" property tag.
	//
	// Used by View.
	// Set the line style of a view's left border.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The border will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a border.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//   - 4 (DoubleLine) or "double" - Double line as a border.
	BorderLeftStyle PropertyName = "border-left-style"

	// BorderRightStyle is the constant for "border-right-style" property tag.
	//
	// Used by View.
	// Set the line style of a view's right border.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The border will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a border.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//   - 4 (DoubleLine) or "double" - Double line as a border.
	BorderRightStyle PropertyName = "border-right-style"

	// BorderTopStyle is the constant for "border-top-style" property tag.
	//
	// Used by View.
	// Set the line style of a view's top border.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The border will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a border.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//   - 4 (DoubleLine) or "double" - Double line as a border.
	BorderTopStyle PropertyName = "border-top-style"

	// BorderBottomStyle is the constant for "border-bottom-style" property tag.
	//
	// Used by View.
	// Sets the line style of a view's bottom border.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The border will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as a border.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//   - 4 (DoubleLine) or "double" - Double line as a border.
	BorderBottomStyle PropertyName = "border-bottom-style"

	// BorderWidth is the constant for "border-width" property tag.
	//
	// Used by View.
	// Set the line width for all four sides of a view's border.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	BorderWidth PropertyName = "border-width"

	// BorderLeftWidth is the constant for "border-left-width" property tag.
	//
	// Used by View.
	// Set the line width of a view's left border.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	BorderLeftWidth PropertyName = "border-left-width"

	// BorderRightWidth is the constant for "border-right-width" property tag.
	//
	// Used by View.
	// Set the line width of a view's right border.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	BorderRightWidth PropertyName = "border-right-width"

	// BorderTopWidth is the constant for "border-top-width" property tag.
	//
	// Used by View.
	// Set the line width of a view's top border.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	BorderTopWidth PropertyName = "border-top-width"

	// BorderBottomWidth is the constant for "border-bottom-width" property tag.
	//
	// Used by View.
	// Set the line width of a view's bottom border.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	BorderBottomWidth PropertyName = "border-bottom-width"

	// BorderColor is the constant for "border-color" property tag.
	//
	// Used by View.
	// Set the line color for all four sides of a view's border.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BorderColor PropertyName = "border-color"

	// BorderLeftColor is the constant for "border-left-color" property tag.
	//
	// Used by View.
	// Set the line color of a view's left border.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BorderLeftColor PropertyName = "border-left-color"

	// BorderRightColor is the constant for "border-right-color" property tag.
	//
	// Used by View.
	// Set the line color of a view's right border.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BorderRightColor PropertyName = "border-right-color"

	// BorderTopColor is the constant for "border-top-color" property tag.
	//
	// Used by View.
	// Set the line color of a view's top border.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BorderTopColor PropertyName = "border-top-color"

	// BorderBottomColor is the constant for "border-bottom-color" property tag.
	//
	// Used by View.
	// Set the line color of a view's bottom border.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BorderBottomColor PropertyName = "border-bottom-color"

	// Outline is the constant for "outline" property tag.
	//
	// Used by View.
	// Set a view's outline. It sets the values of an outline width, style, and color.
	//
	// Supported types: OutlineProperty, ViewOutline, ViewBorder.
	//
	// Internal type is OutlineProperty, other types converted to it during assignment.
	// See OutlineProperty, ViewOutline and ViewBorder description for more details.
	Outline PropertyName = "outline"

	// OutlineStyle is the constant for "outline-style" property tag.
	//
	// Used by View.
	// Set the style of an view's outline.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneLine) or "none" - The outline will not be drawn.
	//   - 1 (SolidLine) or "solid" - Solid line as an outline.
	//   - 2 (DashedLine) or "dashed" - Dashed line as an outline.
	//   - 3 (DottedLine) or "dotted" - Dotted line as an outline.
	//   - 4 (DoubleLine) or "double" - Double line as an outline.
	OutlineStyle PropertyName = "outline-style"

	// OutlineColor is the constant for "outline-color" property tag.
	//
	// Used by View.
	// Set the color of an view's outline.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	OutlineColor PropertyName = "outline-color"

	// OutlineWidth is the constant for "outline-width" property tag.
	//
	// Used by View.
	// Set the width of an view's outline.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	OutlineWidth PropertyName = "outline-width"

	// OutlineOffset is the constant for "outline-offset" property tag.
	//
	// Used by View.
	// Set the amount of space between an outline and the edge or border of a view.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	OutlineOffset PropertyName = "outline-offset"

	// Shadow is the constant for "shadow" property tag.
	//
	// Used by View.
	// Adds shadow effects around a view's frame. A shadow is described by X and Y offsets relative to the element, blur,
	// spread radius and color.
	//
	// Supported types: ShadowProperty, []ShadowProperty, string.
	//
	// Internal type is []ShadowProperty, other types converted to it during assignment.
	// See ShadowProperty description for more details.
	//
	// Conversion rules:
	//   - []ShadowProperty - stored as is. no conversion performed.
	//   - ShadowProperty - converted to []ShadowProperty during assignment.
	//   - string - must contain a string representation of ShadowProperty
	Shadow PropertyName = "shadow"

	// FontName is the constant for "font-name" property tag.
	//
	// Used by View.
	// Specifies a prioritized list of one or more font family names and/or generic family names for the view. Values are
	// separated by commas to indicate that they are alternatives. This is an inherited property, i.e. if it is not defined,
	// then the value of the parent view is used.
	//
	// Supported types: string.
	FontName PropertyName = "font-name"

	// TextColor is the constant for "text-color" property tag.
	//
	// Used by View.
	// Set the foreground color value of a view's text and text decorations. This is an inherited property, i.e. if it is not
	// defined, then the value of the parent view is used.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	TextColor PropertyName = "text-color"

	// TextSize is the constant for "text-size" property tag.
	//
	// Used by View.
	// Set the size of the font. This is an inherited property, i.e. if it is not defined, then the value of the parent view
	// is used.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	TextSize PropertyName = "text-size"

	// Italic is the constant for "italic" property tag.
	//
	// Used by View.
	// Controls whether the text is displayed in italics. This is an inherited property, i.e. if it is not defined, then the
	// value of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Text is displayed in italics.
	//   - false, 0, "false", "no", "off", or "0" - Normal text.
	Italic PropertyName = "italic"

	// SmallCaps is the constant for "small-caps" property tag.
	//
	// Used by View.
	// Controls whether to use small caps characters while displaying the text. This is an inherited property, i.e. if it is
	// not defined, then the value of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Text displayed using small caps.
	//   - false, 0, "false", "no", "off", or "0" - Normal text display.
	SmallCaps PropertyName = "small-caps"

	// Strikethrough is the constant for "strikethrough" property tag.
	//
	// Used by View.
	// Controls whether to draw line over the text. This is an inherited property, i.e. if it is not defined, then the value
	// of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Draw line over the text.
	//   - false, 0, "false", "no", "off", or "0" - Normal text display.
	Strikethrough PropertyName = "strikethrough"

	// Overline is the constant for "overline" property tag.
	//
	// Used by View.
	// Controls whether the line needs to be displayed on top of the text. This is an inherited property, i.e. if it is not
	// defined, then the value of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Overline text.
	//   - false, 0, "false", "no", "off", or "0" - No overline.
	Overline PropertyName = "overline"

	// Underline is the constant for "underline" property tag.
	//
	// Used by View.
	// Controls whether to draw line below the text, This is an inherited property, i.e. if it is not defined, then the value
	// of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Draw line below the text.
	//   - false, 0, "false", "no", "off", or "0" - Normal text display.
	Underline PropertyName = "underline"

	// TextLineThickness is the constant for "text-line-thickness" property tag.
	//
	// Used by View.
	// Set the stroke thickness of the decoration line that is used on text in an element, such as a strikethrough, underline,
	// or overline. This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	TextLineThickness PropertyName = "text-line-thickness"

	// TextLineStyle is the constant for "text-line-style" property tag.
	//
	// Used by View.
	// Set the style of the lines specified by "strikethrough", "overline" and "underline" properties. This is an inherited
	// property, i.e. if it is not defined, then the value of the parent view is used.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 1 (SolidLine) or "solid" - Solid line as a text line.
	//   - 2 (DashedLine) or "dashed" - Dashed line as a text line.
	//   - 3 (DottedLine) or "dotted" - Dotted line as a text line.
	//   - 4 (DoubleLine) or "double" - Double line as a text line.
	//   - 5 (WavyLine) or "wavy" - Wavy line as a text line.
	TextLineStyle PropertyName = "text-line-style"

	// TextLineColor is the constant for "text-line-color" property tag.
	//
	// Used by View.
	// Sets the color of the lines specified by "strikethrough", "overline" and "underline" properties. This is an inherited
	// property, i.e. if it is not defined, then the value of the parent view is used.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	TextLineColor PropertyName = "text-line-color"

	// TextWeight is the constant for "text-weight" property tag.
	//
	// Used by View.
	// Sets weight of the text.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 1 (ThinFont) or "thin" - Thin font.
	//   - 2 (ExtraLightFont) or "extra-light" - Extra light font.
	//   - 3 (LightFont) or "light" - Light font.
	//   - 4 (NormalFont) or "normal" - Normal font.
	//   - 5 (MediumFont) or "medium" - Medium font.
	//   - 6 (SemiBoldFont) or "semi-bold" - Semi-bold font.
	//   - 7 (BoldFont) or "bold" - Bold font.
	//   - 8 (ExtraBoldFont) or "extra-bold" - Extra bold font.
	//   - 9 (BlackFont) or "black" - Black font.
	TextWeight PropertyName = "text-weight"

	// TextAlign is the constant for "text-align" property tag.
	//
	// Used by TableView, View.
	//
	// Usage in TableView:
	// Sets the horizontal alignment of the content inside a table cell.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (JustifyAlign) or "justify" - Justify alignment.
	//
	// Usage in View:
	// Alignment of the text in view. This is an inherited property, i.e. if it is not defined, then the value of the parent
	// view is used.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (JustifyAlign) or "justify" - Justify alignment.
	TextAlign PropertyName = "text-align"

	// TextIndent is the constant for "text-indent" property tag.
	//
	// Used by View.
	// Determines the size of the indent(empty space) before the first line of text. This is an inherited property, i.e. if it
	// is not defined, then the value of the parent view is used.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	TextIndent PropertyName = "text-indent"

	// TextShadow is the constant for "text-shadow" property tag.
	//
	// Used by View.
	// Specify shadow for the text.
	//
	// Supported types: ShadowProperty, []ShadowProperty, string.
	//
	// Internal type is []ShadowProperty, other types converted to it during assignment.
	// See ShadowProperty description for more details.
	//
	// Conversion rules:
	//   - []ShadowProperty - stored as is. no conversion performed.
	//   - ShadowProperty - converted to []ShadowProperty during assignment.
	//   - string - must contain a string representation of ShadowProperty
	TextShadow PropertyName = "text-shadow"

	// TextWrap is the constant for "text-wrap" property tag.
	//
	// Used by View.
	// Controls how text inside the view is wrapped. Default value is "wrap".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TextWrapOn) or "wrap" - Text is wrapped across lines at appropriate characters (for example spaces, in languages like English that use space separators) to minimize overflow.
	//   - 1 (TextWrapOff) or "nowrap" - Text does not wrap across lines. It will overflow its containing element rather than breaking onto a new line.
	//   - 2 (TextWrapBalance) or "balance" - Text is wrapped in a way that best balances the number of characters on each line, enhancing layout quality and legibility. Because counting characters and balancing them across multiple lines is computationally expensive, this value is only supported for blocks of text spanning a limited number of lines (six or less for Chromium and ten or less for Firefox).
	TextWrap PropertyName = "text-wrap"

	// TabSize is the constant for "tab-size" property tag.
	//
	// Used by View.
	// Set the width of tab characters (U+0009) in spaces. This is an inherited property, i.e. if it is not defined, then the
	// value of the parent view is used. Default value is 8.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - greater than 0 - Number of spaces in tab character.
	//   - 0 or negative - ignored.
	TabSize PropertyName = "tab-size"

	// LetterSpacing is the constant for "letter-spacing" property tag.
	//
	// Used by View.
	// Set the horizontal spacing behavior between text characters. This value is added to the natural spacing between
	// characters while rendering the text. Positive values of letter-spacing causes characters to spread farther apart, while
	// negative values of letter-spacing bring characters closer together. This is an inherited property, i.e. if it is not
	// defined, then the value of the parent view is used.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	LetterSpacing PropertyName = "letter-spacing"

	// WordSpacing is the constant for "word-spacing" property tag.
	//
	// Used by View.
	// Set the length of space between words and between tags. This is an inherited property, i.e. if it is not defined, then
	// the value of the parent view is used.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	WordSpacing PropertyName = "word-spacing"

	// LineHeight is the constant for "line-height" property tag.
	//
	// Used by View.
	// Set the height of a line box. It's commonly used to set the distance between lines of text. This is an inherited
	// property, i.e. if it is not defined, then the value of the parent view is used.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	LineHeight PropertyName = "line-height"

	// WhiteSpace is the constant for "white-space" property tag.
	//
	// Used by View.
	// Sets how white space inside an element is handled. This is an inherited property, i.e. if it is not defined, then the
	// value of the parent view is used.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (WhiteSpaceNormal) or "normal" - Sequences of spaces are concatenated into one space.
	//    Newlines in the source are treated as a single space. Applying this value optionally splits lines to fill inline boxes.
	//   - 1 (WhiteSpaceNowrap) or "nowrap" - Concatenates sequences of spaces into one space, like a normal value, but does not wrap lines(text wrapping) within the text.
	//   - 2 (WhiteSpacePre) or "pre" - Sequences of spaces are saved as they are specified in the source.
	//    Lines are wrapped only where newlines are specified in the source and where "br" elements are specified in the source.
	//   - 3 (WhiteSpacePreWrap) or "pre-wrap" - Sequences of spaces are saved as they are indicated in the source.
	//    Lines are wrapped only where newlines are specified in the source and there, where "br" elements are specified in the source, and optionally to fill inline boxes.
	//   - 4 (WhiteSpacePreLine) or "pre-line" - Sequences of spaces are concatenated into one space. Lines are split on newlines, on "br" elements, and optionally to fill inline boxes.
	//   - 5 (WhiteSpaceBreakSpaces) or "break-spaces" - The behavior is identical to WhiteSpacePreWrap with the some differences.
	//
	// Differences WhiteSpaceBreakSpaces (5) from WhiteSpacePreWrap(3):
	//  1. Sequences of spaces are preserved as specified in the source, including spaces at the end of lines.
	//  2. Lines are wrapped on any spaces, including in the middle of a sequence of spaces.
	//  3. Spaces take up space and do not hang at the ends of lines, which means they affect the internal dimensions (min-content and max-content).
	WhiteSpace PropertyName = "white-space"

	// WordBreak is the constant for "word-break" property tag.
	//
	// Used by View.
	// Set whether line breaks appear wherever the text would otherwise overflow its content box. This is an inherited
	// property, i.e. if it is not defined, then the value of the parent view is used. Default value is "normal".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (WordBreakNormal) or "normal" - Default behavior for linefeed placement.
	//   - 1 (WordBreakAll) or "break-all" - If the block boundaries are exceeded, a line break will be inserted between any two characters(except for Chinese/Japanese/Korean text).
	//   - 2 (WordBreakKeepAll) or "keep-all" - Line break will not be used in Chinese/Japanese/ Korean text. For text in other languages, the default behavior(normal) will be applied.
	//   - 3 (WordBreakWord) or "break-word" - When the block boundaries are exceeded, the remaining whole words can be broken in an arbitrary place, if a more suitable place for line break is not found.
	WordBreak PropertyName = "word-break"

	// TextTransform is the constant for "text-transform" property tag.
	//
	// Used by View.
	// Specifies how to capitalize an element's text. It can be used to make text appear in all-uppercase or all-lowercase, or
	// with each word capitalized. This is an inherited property, i.e. if it is not defined, then the value of the parent view
	// is used.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneTextTransform) or "none" - Original case of characters.
	//   - 1 (CapitalizeTextTransform) or "capitalize" - Every word starts with a capital letter.
	//   - 2 (LowerCaseTextTransform) or "lowercase" - All characters are lowercase.
	//   - 3 (UpperCaseTextTransform) or "uppercase" - All characters are uppercase.
	TextTransform PropertyName = "text-transform"

	// TextDirection is the constant for "text-direction" property tag.
	//
	// Used by View.
	// Sets the direction of text, table columns, and horizontal overflow. This is an inherited property, i.e. if it is not
	// defined, then the value of the parent view is used. Default value is "system".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (SystemTextDirection) or "system" - Use the system text direction.
	//   - 1 (LeftToRightDirection) or "left-to-right" - For languages written from left to right (like English and most other languages).
	//   - 2 (RightToLeftDirection) or "right-to-left" - For languages written from right to left (like Hebrew or Arabic).
	TextDirection PropertyName = "text-direction"

	// WritingMode is the constant for "writing-mode" property tag.
	//
	// Used by View.
	// Set whether lines of text are laid out horizontally or vertically, as well as the direction in which blocks progress.
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used. Default value is
	// "horizontal-top-to-bottom".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (HorizontalTopToBottom) or "horizontal-top-to-bottom" - Horizontal lines are displayed from top to bottom.
	//   - 1 (HorizontalBottomToTop) or "horizontal-bottom-to-top" - Horizontal lines are displayed from bottom to top.
	//   - 2 (VerticalRightToLeft) or "vertical-right-to-left" - Vertical lines are output from right to left.
	//   - 3 (VerticalLeftToRight) or "vertical-left-to-right" - Vertical lines are output from left to right.
	WritingMode PropertyName = "writing-mode"

	// VerticalTextOrientation is the constant for "vertical-text-orientation" property tag.
	//
	// Used by View.
	// Set the orientation of the text characters in a line. It only affects text in vertical mode ("writing-mode" property).
	// This is an inherited property, i.e. if it is not defined, then the value of the parent view is used.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (MixedTextOrientation) or "mixed" - Symbols rotated 90° clockwise.
	//   - 1 (UprightTextOrientation) or "upright" - Symbols are arranged normally(vertically).
	VerticalTextOrientation PropertyName = "vertical-text-orientation"

	// TextOverflow is the constant for "text-overflow" property tag.
	//
	// Used by TextView.
	// Sets how hidden overflow content is signaled to users. Default value is "clip".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TextOverflowClip) or "clip" - Text is clipped at the border.
	//   - 1 (TextOverflowEllipsis) or "ellipsis" - At the end of the visible part of the text "…" is displayed.
	TextOverflow PropertyName = "text-overflow"

	// Hint is the constant for "hint" property tag.
	//
	// Used by EditView.
	// Sets a hint to the user of what can be entered in the control.
	//
	// Supported types: string.
	Hint PropertyName = "hint"

	// MaxLength is the constant for "max-length" property tag.
	//
	// Used by EditView.
	// Sets the maximum number of characters that the user can enter.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - positive value - Maximum number of characters.
	//   - 0 or negative value - The maximum number of characters is not limited.
	MaxLength PropertyName = "max-length"

	// ReadOnly is the constant for "readonly" property tag.
	//
	// Used by EditView.
	// Controls whether the user can modify value or not. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - User not able to modify the value.
	//   - false, 0, "false", "no", "off", or "0" - Value can be modified.
	ReadOnly PropertyName = "readonly"

	// Content is the constant for "content" property tag.
	//
	// Used by ViewsContainer, GridLayout, ListLayout, Resizable, StackLayout, SvgImageView, TableView.
	//
	// # Usage in ViewsContainer:
	//
	// An array of child views.
	//
	// Supported types: View, []View, string, []string, []any containing elements of View, string.
	//
	// Internal type is []View, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - View - converted to []View containing one element.
	//   - []View - nil-elements are prohibited, if the array contains nil, then the property will not be set, and the Set function will return false and an error message will be written to the log.
	//   - string - if the string is a text representation of the View, then the corresponding view is created, otherwise a TextView is created, to which the given string is passed as a text. Then a []View is created containing the resulting view.
	//   - []string - each element of an array is converted to View as described above.
	//   - []any - this array must contain only View and a string. Each string element is converted to a view as described above. If array contains invalid values, the "content" property will not be set, and the Set function will return false and an error message will be written to the log.
	//
	// # Usage in GridLayout:
	//
	// Defines an array of child views or can be an implementation of GridAdapter interface.
	//
	// Supported types: []View, GridAdapter, View, string, []string.
	//
	// Internal type is either []View or GridAdapter, other types converted to []View during assignment.
	//
	// Conversion rules:
	//   - View - view which describe one cell, converted to []View.
	//   - []View - describe several cells, stored as is.
	//   - string - text representation of the view which describe one cell, converted to []View.
	//   - []string - an array of text representation of the views which describe several cells, converted to []View.
	//   - GridAdapter - interface which describe several cells, see GridAdapter description for more details.
	//
	// # Usage in ListLayout:
	//
	// Defines an array of child views or can be an implementation of ListAdapter interface.
	//
	// Supported types: []View, ListAdapter, View, string, []string.
	//
	// Internal type is either []View or ListAdapter, other types converted to []View during assignment.
	//
	// Conversion rules:
	//   - View - view which describe one item, converted to []View.
	//   - []View - describe several items, stored as is.
	//   - string - text representation of the view which describe one item, converted to []View.
	//   - []string - an array of text representation of the views which describe several items, converted to []View.
	//   - ListAdapter - interface which describe several items, see ListAdapter description for more details.
	//
	// # Usage in Resizable:
	//
	// Content view to make it resizable or text in this case TextView will be created.
	//
	// Supported types: View, string.
	//
	// Internal type is View, other types converted to it during assignment.
	//
	// Usage in StackLayout:
	// An array of child views.
	//
	// Supported types: View, []View, string, []string, []any containing elements of View, string.
	//
	// Internal type is []View, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - View - converted to []View containing one element.
	//   - []View - nil-elements are prohibited, if the array contains nil, then the property will not be set, and the Set function will return false and an error message will be written to the log.
	//   - string - if the string is a text representation of the View, then the corresponding view is created, otherwise a TextView is created, to which the given string is passed as a text. Then a []View is created containing the resulting view.
	//   - []string - each element of an array is converted to View as described above.
	//   - []any - this array must contain only View and a string. Each string element is converted to a view as described above. If array contains invalid values, the "content" property will not be set, and the Set function will return false and an error message will be written to the log.
	//
	// # Usage in SvgImageView:
	//
	// Image to display. Could be the image file name in the images folder of the resources, image URL or content of the svg image.
	//
	// Supported types: string.
	//
	// # Usage in TableView:
	//
	// Defines the content of the table.
	//
	// Supported types: TableAdapter, [][]string, [][]any.
	//
	// Internal type is TableAdapter, other types converted to it during assignment.
	// See TableAdapter description for more details.
	Content PropertyName = "content"

	// Items is the constant for "items" property tag.
	//
	// Used by DropDownList, ListView, Popup.
	//
	// # Usage in DropDownList:
	//
	// Array of data elements.
	//
	// Supported types: []string, string, []fmt.Stringer, []Color, []SizeUnit, []AngleUnit, []any containing
	// elements of string, fmt.Stringer, bool, rune, float32, float64, int, int8…int64, uint, uint8…uint64.
	//
	// Internal type is []string, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - string - contain single item.
	//   - []string - an array of items.
	//   - []fmt.Stringer - an array of objects convertible to string.
	//   - []Color - An array of color values which will be converted to a string array.
	//   - []SizeUnit - an array of size unit values which will be converted to a string array.
	//   - []any - this array must contain only types which were listed in Types section.
	//
	// # Usage in ListView:
	//
	// List content. Main value is an implementation of ListAdapter interface.
	//
	// Supported types: ListAdapter, []View, []string, []any containing elements of View, string, fmt.Stringer, float and int.
	//
	// Internal type is either []View or ListAdapter, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - ListAdapter - interface which provides an access to list items and other information, stored as is.
	//   - []View - an array of list items, each in a form of some view-based element. Stored as is.
	//   - []string - an array of text. Converted into an internal implementation of ListAdapter, each list item will be an instance of TextView.
	//   - []any - an array of items of arbitrary type, where types like string, fmt.Stringer, float and int will be converted to TextView. View type will remain unchanged. All values after conversion will be wrapped by internal implementation of ListAdapter.
	//
	// # Usage in Popup:
	//
	// Array of menu items.
	//
	// Supported types: ListAdapter, []string.
	//
	// Internal type is ListAdapter internal implementation, other types converted to it during assignment.
	Items PropertyName = "items"

	// DisabledItems is the constant for "disabled-items" property tag.
	//
	// Used by DropDownList.
	// An array of disabled(non selectable) items indices.
	//
	// Supported types: []int, string, []string, []any containing  elements of string or int.
	//
	// Internal type is []int, other types converted to it during assignment.
	// Rules of conversion.
	//   - []int - Array of indices.
	//   - string - Single index value or multiple index values separated by comma(,).
	//   - []string - Array of indices in text format.
	//   - []any - Array of strings or integer values.
	DisabledItems PropertyName = "disabled-items"

	// ItemSeparators is the constant for "item-separators" property tag.
	//
	// Used by DropDownList.
	// An array of indices of DropDownList items after which a separator should be added.
	//
	// Supported types: []int, string, []string, []any containing  elements of string or int.
	//
	// Internal type is []int, other types converted to it during assignment.
	// Rules of conversion.
	//   - []int - Array of indices.
	//   - string - Single index value or multiple index values separated by comma(,).
	//   - []string - Array of indices in text format.
	//   - []any - Array of strings or integer values.
	ItemSeparators PropertyName = "item-separators"

	// Current is the constant for "current" property tag.
	//
	// Used by DropDownList, ListView, TableView, TabsLayout.
	//
	// # Usage in DropDownList:
	//
	// Current selected item.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - negative value - No item has been selected.
	//   - not negative value - Index of selected item.
	//
	// # Usage in ListView:
	//
	// Set or get index of selected item.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - negative value - No item has been selected.
	//   - not negative value - Index of selected item.
	//
	// # Usage in TableView:
	//
	// Sets the coordinates of the selected cell/row.
	//
	// Supported types: CellIndex, int, string.
	//
	// Internal type is CellIndex, other types converted to it during assignment.
	// See CellIndex description for more details.
	//
	// Conversion rules:
	//   - int - specify index of current table row, current column index will be set to -1.
	//   - string - can be one integer value which specify current row or pair of integer values separated by comma(,). When two values provided then first value specify current row index and second one specify column index.
	//
	// # Usage in TabsLayout:
	//
	// Defines index of the current active child view.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - negative value - No visible tab.
	//   - not negative value - Index of visible tab.
	Current PropertyName = "current"

	// Type is the constant for "type" property tag.
	//
	// Used by EditView, NumberPicker.
	//
	// Usage in EditView:
	// Same as "edit-view-type" [EditViewType].
	//
	// Usage in NumberPicker:
	// Same as "number-picker-type" [NumberPickerType].
	Type PropertyName = "type"

	// Pattern is the constant for "pattern" property tag.
	//
	// Used by EditView.
	// Same as "edit-view-pattern" [EditViewPattern].
	Pattern PropertyName = "pattern"

	// GridAutoFlow is the constant for "grid-auto-flow" property tag.
	//
	// Used by GridLayout.
	// Controls how to place child controls if Row and Column properties were not set for children views.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (RowAutoFlow) or "row" - Views are placed by filling each row in turn, adding new rows as necessary.
	//   - 1 (ColumnAutoFlow) or "column" - Views are placed by filling each column in turn, adding new columns as necessary.
	//   - 2 (RowDenseAutoFlow) or "row-dense" - Views are placed by filling each row, adding new rows as necessary. "dense" packing algorithm attempts to fill in holes earlier in the grid, if smaller items come up later. This may cause views to appear out-of-order, when doing so would fill in holes left by larger views.
	//   - 3 (ColumnDenseAutoFlow) or "column-dense" - Views are placed by filling each column, adding new columns as necessary. "dense" packing algorithm attempts to fill in holes earlier in the grid, if smaller items come up later. This may cause views to appear out-of-order, when doing so would fill in holes left by larger views.
	GridAutoFlow PropertyName = "grid-auto-flow"

	// CellWidth is the constant for "cell-width" property tag.
	//
	// Used by GridLayout:
	// Set a fixed width of GridLayout cells regardless of the size of the child elements. Each element in the array
	// determines the size of the corresponding column. By default, the sizes of the cells are calculated based on the sizes
	// of the child views placed in them.
	//
	// Supported types: SizeUnit, []SizeUnit, SizeFunc, string, []string, []any containing elements of string or
	// SizeUnit.
	//
	// Internal type is either SizeUnit or []SizeUnit, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - SizeUnit, SizeFunc - stored as is and all cells are set to have the same width.
	//   - []SizeUnit - stored as is and each column of the grid layout has width which is specified in an array.
	//   - string - containing textual representations of SizeUnit (or SizeUnit constants), may contain several values separated by comma(,). Each column of the grid layout has width which is specified in an array.
	//   - []string - each element must be a textual representation of a SizeUnit (or a SizeUnit constant). Each column of the grid layout has width which is specified in an array.
	// If the number of elements in an array is less than the number of columns used, then the missing elements are set to have Auto size.
	//
	// The values can use SizeUnit type SizeInFraction. This type means 1 part. The part is calculated as follows: the size of all cells that are not of type SizeInFraction is subtracted from the size of the container, and then the remaining size is divided by the number of parts. The SizeUnit value of type SizeInFraction can be either integer or fractional.
	CellWidth PropertyName = "cell-width"

	// CellHeight is the constant for "cell-height" property tag.
	//
	// Used by GridLayout:
	// Set a fixed height of GridLayout cells regardless of the size of the child elements. Each element in the array
	// determines the size of the corresponding row. By default, the sizes of the cells are calculated based on the sizes of
	// the child views placed in them.
	//
	// Supported types: SizeUnit, []SizeUnit, SizeFunc, string, []string, []any containing elements of string or
	// SizeUnit.
	//
	// Internal type is either SizeUnit or []SizeUnit, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - SizeUnit, SizeFunc - stored as is and all cells are set to have the same height.
	//   - []SizeUnit - stored as is and each row of the grid layout has height which is specified in an array.
	//   - string - containing textual representations of SizeUnit (or SizeUnit constants), may contain several values separated by comma(,). Each row of the grid layout has height which is specified in an array.
	//   - []string - each element must be a textual representation of a SizeUnit (or a SizeUnit constant). Each row of the grid layout has height which is specified in an array.
	// If the number of elements in an array is less than the number of rows used, then the missing elements are set to have Auto size.
	//
	// The values can use SizeUnit type SizeInFraction. This type means 1 part. The part is calculated as follows: the size of all cells that are not of type SizeInFraction is subtracted from the size of the container, and then the remaining size is divided by the number of parts. The SizeUnit value of type SizeInFraction can be either integer or fractional.
	CellHeight PropertyName = "cell-height"

	// GridRowGap is the constant for "grid-row-gap" property tag.
	//
	// Used by GridLayout.
	// Space between rows.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	GridRowGap PropertyName = "grid-row-gap"

	// GridColumnGap is the constant for "grid-column-gap" property tag.
	//
	// Used by GridLayout.
	// Space between columns.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	GridColumnGap PropertyName = "grid-column-gap"

	// Source is the constant for "src" property tag.
	//
	// Used by AudioPlayer, ImageView, VideoPlayer.
	//
	// # Usage in AudioPlayer
	//
	// Specifies the location of the media file(s). Since different browsers support different file formats and codecs, it is
	// recommended to specify multiple sources in different formats. The player chooses the most suitable one from the list of
	// sources. Setting mime types makes this process easier for the browser.
	//
	// Supported types: string, MediaSource, []MediaSource.
	//
	// Internal type is []MediaSource, other types converted to it during assignment.
	//
	// # Usage in ImageView
	//
	// Set either the name of the image in the "images" folder of the resources, or the URL of the image or inline-image. An
	// inline-image is the content of an image file encoded in base64 format.
	//
	// Supported types: string.
	//
	// # Usage in VideoPlayer
	//
	// Specifies the location of the media file(s). Since different browsers support different file formats and codecs, it is
	// recommended to specify multiple sources in different formats. The player chooses the most suitable one from the list of
	// sources. Setting mime types makes this process easier for the browser.
	//
	// Supported types: string, MediaSource, []MediaSource.
	//
	// Internal type is []MediaSource, other types converted to it during assignment.
	Source PropertyName = "src"

	// SrcSet is the constant for "srcset" property tag.
	//
	// Used by ImageView.
	// String which identifies one or more image candidate strings, separated using comma(,) each specifying image resources
	// to use under given screen density. This property is only used if building an application for js/wasm platform.
	//
	// Supported types: string.
	SrcSet PropertyName = "srcset"

	// Fit is the constant for "fit" property tag.
	//
	// Used by ImageView, BackgroundElement.
	//
	// # Usage in ImageView
	//
	// Defines the image scaling parameters.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneFit) or "none" - The image is not resized.
	//   - 1 (ContainFit) or "contain" - The image is scaled to maintain its aspect ratio while fitting within the element’s content box. The entire object is made to fill the box, while preserving its aspect ratio, so the object will be "letterboxed" if its aspect ratio does not match the aspect ratio of the box.
	//   - 2 (CoverFit) or "cover" - The image is sized to maintain its aspect ratio while filling the element’s entire content box. If the object's aspect ratio does not match the aspect ratio of its box, then the object will be clipped to fit.
	//   - 3 (FillFit) or "fill" - The image to fill the element’s content box. The entire object will completely fill the box. If the object's aspect ratio does not match the aspect ratio of its box, then the object will be stretched to fit.
	//   - 4 (ScaleDownFit) or "scale-down" - The image is sized as if NoneFit or ContainFit were specified, whichever would result in a smaller concrete object size.
	//
	// # Usage in BackgroundElement
	//
	// Used for image background only. Defines the image scaling parameters.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneFit) or "none" - The image is not resized.
	//   - 1 (ContainFit) or "contain" - The image is scaled to maintain its aspect ratio while fitting within the element’s content box. The entire object is made to fill the box, while preserving its aspect ratio, so the object will be "letterboxed" if its aspect ratio does not match the aspect ratio of the box.
	//   - 2 (CoverFit) or "cover" - The image is sized to maintain its aspect ratio while filling the element’s entire content box. If the object's aspect ratio does not match the aspect ratio of its box, then the object will be clipped to fit.
	Fit           PropertyName = "fit"
	backgroundFit PropertyName = "background-fit"

	// Repeat is the constant for "repeat" property tag.
	//
	// Used by BackgroundElement.
	// Used for image background only. Specifying the repetition of the image. Default value is "no-repeat".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoRepeat) or "no-repeat" - Image does not repeat.
	//   - 1 (RepeatXY) or "repeat" - Image repeat horizontally and vertically.
	//   - 2 (RepeatX) or "repeat-x" - Image repeat only horizontally.
	//   - 3 (RepeatY) or "repeat-y" - Image repeat only vertically.
	//   - 4 (RepeatRound) or "round" - Image is repeated so that an integer number of images fit into the background area. If this fails, then the background images are scaled.
	//   - 5 (RepeatSpace) or "space" - Image is repeated as many times as necessary to fill the background area. If this fails, an empty space is added between the pictures.
	Repeat PropertyName = "repeat"

	// Attachment is the constant for "attachment" property tag.
	//
	// Used by BackgroundElement.
	// Used for image background only. Sets whether a background image's position is fixed within the viewport or scrolls with
	// its containing block.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (ScrollAttachment) or "scroll" - The background image will scroll with the page.
	//   - 1 (FixedAttachment) or "fixed" - The background image will not scroll with the page.
	//   - 2 (LocalAttachment) or "local" - The background image will scroll with the element's contents.
	Attachment PropertyName = "attachment"

	// BackgroundClip is the constant for "background-clip" property tag.
	//
	// Used by View.
	// Determines how the background color and/or background image will be displayed below the box borders.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (BorderBoxClip) or "border-box" - The background extends to the outer edge of the border(but below the border in z-order).
	//   - 1 (PaddingBoxClip) or "padding-box" - The background extends to the outer edge of the padding. No background is drawn below the border.
	//   - 2 (ContentBoxClip) or "content-box" - The background is painted inside(clipped) of the content box.
	BackgroundClip PropertyName = "background-clip"

	// BackgroundOrigin is the constant for "background-origin" property tag.
	//
	// Used by View.
	// Determines the background's origin: from the border start, inside the border, or inside the padding.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (BorderBox) or "border-box" - The background is positioned relative to the border box.
	//   - 1 (PaddingBox) or "padding-box" - The background is positioned relative to the padding box.
	//   - 2 (ContentBox) or "content-box" - The background is positioned relative to the content box.
	BackgroundOrigin PropertyName = "background-origin"

	// MaskClip is the constant for "mask-clip" property tag.
	//
	// Used by View.
	// Determines how image/gradient masks will be used below the box borders.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (BorderBox) or "border-box" - The mask extends to the outer edge of the border.
	//   - 1 (PaddingBox) or "padding-box" - The mask extends to the outer edge of the padding.
	//   - 2 (ContentBox) or "content-box" - The mask is used inside(clipped) of the content box.
	MaskClip PropertyName = "mask-clip"

	// MaskOrigin is the constant for "mask-origin" property tag.
	//
	// Used by View.
	// Determines the mask's origin: from the border start, inside the border, or inside the padding.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (BorderBox) or "border-box" - The mask is positioned relative to the border box.
	//   - 1 (PaddingBox) or "padding-box" - The mask is positioned relative to the padding box.
	//   - 2 (ContentBox) or "content-box" - The mask is positioned relative to the content box.
	MaskOrigin PropertyName = "mask-origin"

	// Gradient is the constant for "gradient" property tag.
	//
	// Used by BackgroundElement.
	// Describe gradient stop points. This is a mandatory property while describing background gradients.
	//
	// Supported types: []BackgroundGradientPoint, []BackgroundGradientAngle, []GradientPoint, []Color, string.
	//
	// Internal type is []BackgroundGradientPoint or []BackgroundGradientAngle, other types converted to it during assignment.
	// See BackgroundGradientPoint, []BackgroundGradientAngle, []GradientPoint description for more details.
	//
	// Conversion rules:
	//   - []BackgroundGradientPoint - stored as is, no conversion performed. Used to set gradient stop points for linear and radial gradients.
	//   - []BackgroundGradientAngle - stored as is, no conversion performed. Used to set gradient stop points for conic gradient.
	//   - []GradientPoint - converted to []BackgroundGradientPoint. Used to set gradient stop points for linear and radial gradients. Since GradientPoint contains values from 0 to 1.0 they will be converted to precent values.
	//   - []Color - converted to []BackgroundGradientPoint. Used for setting gradient stop points which are uniformly distributed across gradient diretion.
	//   - string - string representation of stop points or color values. Format: "color1 pos1,color2 pos2"... . Position of stop points can be described either in SizeUnit or AngleUnit string representations. Examples: "white 0deg, black 90deg, gray 360deg", "white 0%, black 100%".
	Gradient PropertyName = "gradient"

	// Direction is the constant for "direction" property tag.
	//
	// Used by BackgroundElement.
	// Used for linear gradient only. Defines the direction of the gradient line. Default is 4 (ToBottomGradient) or
	// "to-bottom".
	//
	// Supported types: AngleUnit, int, string.
	//
	// See AngleUnit description for more details.
	//
	// Values:
	//   - 0 (ToTopGradient) or "to-top" - Line goes from bottom to top.
	//   - 1 (ToRightTopGradient) or "to-right-top" - From bottom left to top right.
	//   - 2 (ToRightGradient) or "to-right" - From left to right.
	//   - 3 (ToRightBottomGradient) or "to-right-bottom" - From top left to bottom right.
	//   - 4 (ToBottomGradient) or "to-bottom" - From top to bottom.
	//   - 5 (ToLeftBottomGradient) or "to-left-bottom" - From the upper right corner to the lower left.
	//   - 6 (ToLeftGradient) or "to-left" - From right to left.
	//   - 7 (ToLeftTopGradient) or "to-left-top" - From the bottom right corner to the top left.
	Direction PropertyName = "direction"

	// Repeating is the constant for "repeating" property tag.
	//
	// Used by BackgroundElement.
	// Define whether stop points needs to be repeated after the last one.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Gradient will repeat after the last key point.
	//   - false, 0, "false", "no", "off", or "0" - No repetition of gradient stop points. Value of the last point used will be extrapolated.
	Repeating PropertyName = "repeating"

	// From is the constant for "from" property tag.
	//
	// Used by BackgroundElement.
	// Used for conic gradient only. Start angle position of the gradient.
	//
	// Supported types: AngleUnit, string, float, int.
	//
	// Internal type is AngleUnit, other types converted to it during assignment.
	// See AngleUnit description for more details.
	From PropertyName = "from"

	// RadialGradientRadius is the constant for "radial-gradient-radius" property tag.
	//
	// Used by BackgroundElement.
	// Define radius of the radial gradient.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	RadialGradientRadius PropertyName = "radial-gradient-radius"

	// RadialGradientShape is the constant for "radial-gradient-shape" property tag.
	//
	// Used by BackgroundElement.
	// Define shape of the radial gradient. The default is 0 (EllipseGradient or "ellipse".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - EllipseGradient (0) or "ellipse" - The shape is an axis-aligned ellipse.
	//   - CircleGradient (1) or "circle" - The shape is a circle with a constant radius.
	RadialGradientShape PropertyName = "radial-gradient-shape"

	// Shape is the constant for "shape" property tag.
	//
	// Used by BackgroundElement.
	// Same as "radial-gradient-shape".
	Shape PropertyName = "shape"

	// CenterX is the constant for "center-x" property tag.
	//
	// Used by BackgroundElement.
	// Used for conic and radial gradients only. Center X point of the gradient.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	CenterX PropertyName = "center-x"

	// CenterY is the constant for "center-y" property tag.
	//
	// Used by BackgroundElement.
	// Used for conic and radial gradients only. Center Y point of the gradient.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	CenterY PropertyName = "center-y"

	// AltText is the constant for "alt-text" property tag.
	//
	// Used by ImageView.
	// Set a description of the image.
	//
	// Supported types: string.
	AltText PropertyName = "alt-text"

	altTag PropertyName = "alt"

	// AvoidBreak is the constant for "avoid-break" property tag.
	//
	// Used by ColumnLayout.
	// Controls how region breaks should behave inside a generated box.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Avoid any break from being inserted within the principal box.
	//   - false, 0, "false", "no", "off", or "0" - Allow, but does not force, any break to be inserted within the principal box.
	AvoidBreak PropertyName = "avoid-break"

	// ItemWidth is the constant for "item-width" property tag.
	//
	// Used by ListView.
	// Fixed width of list elements.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	ItemWidth PropertyName = "item-width"

	// ItemHeight is the constant for "item-height" property tag.
	//
	// Used by ListView.
	// Fixed height of list elements.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	ItemHeight PropertyName = "item-height"

	// ListWrap is the constant for "list-wrap" property tag.
	//
	// Used by ListLayout, ListView.
	//
	// Usage in ListLayout:
	// Defines the position of elements in case of reaching the border of the container.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (ListWrapOff) or "off" - The column or row of elements continues and goes beyond the bounds of the visible area.
	//   - 1 (ListWrapOn) or "on" - Starts a new column or row of elements as necessary. The new column is positioned towards the end.
	//   - 2 (ListWrapReverse) or "reverse" - Starts a new column or row of elements as necessary. The new column is positioned towards the beginning.
	//
	// Usage in ListView:
	// Defines the position of elements in case of reaching the border of the container.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (ListWrapOff) or "off" - The column or row of elements continues and goes beyond the bounds of the visible area.
	//   - 1 (ListWrapOn) or "on" - Starts a new column or row of elements as necessary. The new column is positioned towards the end.
	//   - 2 (ListWrapReverse) or "reverse" - Starts a new column or row of elements as necessary. The new column is positioned towards the beginning.
	ListWrap PropertyName = "list-wrap"

	// EditWrap is the constant for "edit-wrap" property tag.
	//
	// Used by EditView.
	// Controls whether the text will wrap around when edit view border has been reached. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Text wrapped to the next line.
	//   - false, 0, "false", "no", "off", or "0" - Do not wrap text. Horizontal scrolling will appear if necessary.
	EditWrap PropertyName = "edit-wrap"

	// CaretColor is the constant for "caret-color" property tag.
	//
	// Used by EditView.
	//
	// Sets the color of the insertion caret, the visible marker where the next character typed will be inserted. This is
	// sometimes referred to as the text input cursor.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	CaretColor PropertyName = "caret-color"

	// Min is the constant for "min" property tag.
	//
	// Used by DatePicker, NumberPicker, TimePicker.
	//
	// Usage in DatePicker:
	// Same as "date-picker-min" [DatePickerMin].
	//
	// Usage in NumberPicker:
	// Same as "number-picker-min" [NumberPickerMin].
	//
	// Usage in TimePicker:
	// Same as "time-picker-min" [TimePickerMin].
	Min PropertyName = "min"

	// Max is the constant for "max" property tag.
	//
	// Used by DatePicker, NumberPicker, ProgressBar, TimePicker.
	//
	// Usage in DatePicker:
	// Same as "date-picker-max" [DatePickerMax].
	//
	// Usage in NumberPicker:
	// Same as "number-picker-max" [NumberPickerMax].
	//
	// Usage in ProgressBar:
	// Same as "progress-max" [ProgressBarMax].
	//
	// Usage in TimePicker:
	// Same as "time-picker-max" [TimePickerMax].
	Max PropertyName = "max"

	// Step is the constant for "step" property tag.
	//
	// Used by DatePicker, NumberPicker, TimePicker.
	//
	// Usage in DatePicker:
	// Same as "date-picker-step" [DatePickerStep].
	//
	// Usage in NumberPicker:
	// Same as "number-picker-step" [NumberPickerStep].
	//
	// Usage in TimePicker:
	// Same as "time-picker-step" [TimePickerStep].
	Step PropertyName = "step"

	// Value is the constant for "value" property tag.
	//
	// Used by DatePicker, NumberPicker, ProgressBar, TimePicker.
	//
	// Usage in DatePicker:
	// Same as "date-picker-value" [DatePickerValue].
	//
	// Usage in NumberPicker:
	// Same as "number-picker-value" [NumberPickerValue].
	//
	// Usage in ProgressBar:
	// Same as "progress-value" [ProgressBarValue].
	//
	// Usage in TimePicker:
	// Same as "time-picker-value" [TimePickerValue].
	Value PropertyName = "value"

	// Orientation is the constant for "orientation" property tag.
	//
	// Used by ListLayout, ListView.
	//
	// Specifies how the children will be positioned relative to each other.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopDownOrientation) or "up-down" - Child elements are arranged in a column from top to bottom.
	//   - 1 (StartToEndOrientation) or "start-to-end" - Child elements are laid out in a row from beginning to end.
	//   - 2 (BottomUpOrientation) or "bottom-up" - Child elements are arranged in a column from bottom to top.
	//   - 3 (EndToStartOrientation) or "end-to-start" - Child elements are laid out in a line from end to beginning.
	Orientation PropertyName = "orientation"

	// Gap is the constant for "gap" property tag.
	//
	// Used by GridLayout, ListLayout, ListView, TableView.
	//
	// # Usage in GridLayout
	//
	// Specify both "grid-column-gap" and "grid-row-gap".
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in ListLayout and ListView
	//
	// Specify both "list-column-gap" and "list-row-gap".
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	//
	// # Usage in TableView
	//
	// Define the gap between rows and columns of a table.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	Gap PropertyName = "gap"

	// ListRowGap is the constant for "list-row-gap" property tag.
	//
	// Used by ListLayout, ListView.
	//
	// Set the distance between the rows of the ListLayout. Default value 0px.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	ListRowGap PropertyName = "list-row-gap"

	// ListColumnGap is the constant for "list-column-gap" property tag.
	//
	// Used by ListLayout, ListView.
	//
	// Set the distance between the columns of the ListLayout. Default value 0px.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	ListColumnGap PropertyName = "list-column-gap"

	// Text is the constant for "text" property tag.
	//
	// Used by EditView, TextView.
	//
	// Usage in EditView:
	// Edit view text.
	//
	// Supported types: string.
	//
	// Usage in TextView:
	// Text to display.
	//
	// Supported types: string.
	Text PropertyName = "text"

	// VerticalAlign is the constant for "vertical-align" property tag.
	//
	// Used by Checkbox, ListLayout, ListView, Popup, SvgImageView.
	//
	// # Usage in Checkbox
	//
	// Sets the vertical alignment of the content inside a block element.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Content aligned to top side of the content area.
	//   - 1 (BottomAlign) or "bottom" - Content aligned to bottom side of the content area.
	//   - 2 (CenterAlign) or "center" - Content aligned in the center of the content area.
	//   - 3 (StretchAlign) or "stretch" - Content relaxed to fill all content area.
	//
	// # Usage in ListLayout and ListView
	//
	// Sets the vertical alignment of the content inside a block element.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Height alignment.
	//
	// # Usage in Popup
	//
	// Vertical alignment of the popup on the screen.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Height alignment.
	//
	// # Usage in SvgImageView
	//
	// Sets the vertical alignment of the image relative to its bounds.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	VerticalAlign PropertyName = "vertical-align"

	// HorizontalAlign is the constant for "horizontal-align" property tag.
	//
	// Used by Checkbox, ListLayout, ListView, Popup, SvgImageView.
	//
	// # Usage in Checkbox
	//
	// Sets the horizontal alignment of the content inside a block element.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Content aligned to left side of the content area.
	//   - 1 (RightAlign) or "right" - Content aligned to right side of the content area.
	//   - 2 (CenterAlign) or "center" - Content aligned in the center of the content area.
	//   - 3 (StretchAlign) or "stretch" - Content relaxed to fill all content area.
	//
	// # Usage in ListLayout and ListView
	//
	// Sets the horizontal alignment of the content inside a block element.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Width alignment.
	//
	// # Usage in Popup
	//
	// Horizontal alignment of the popup on the screen.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Width alignment.
	//
	// # Usage in SvgImageView
	//
	// Sets the horizontal alignment of the image relative to its bounds.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	HorizontalAlign PropertyName = "horizontal-align"

	// ImageVerticalAlign is the constant for "image-vertical-align" property tag.
	//
	// Used by ImageView.
	// Sets the vertical alignment of the image relative to its bounds.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	ImageVerticalAlign PropertyName = "image-vertical-align"

	// ImageHorizontalAlign is the constant for "image-horizontal-align" property tag.
	//
	// Used by ImageView.
	// Sets the horizontal alignment of the image relative to its bounds.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	ImageHorizontalAlign PropertyName = "image-horizontal-align"

	// Checked is the constant for "checked" property tag.
	//
	// Used by Checkbox, ListView.
	//
	// # Usage in Checkbox
	//
	// Current state of the checkbox.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Checkbox is checked.
	//   - false, 0, "false", "no", "off", or "0" - Checkbox is unchecked.
	//
	// # Usage in ListView
	//
	// Set or get the list of checked items. Stores array of indices of checked items.
	//
	// Supported types: []int, int, string.
	//
	// Internal type is []int, other types converted to it during assignment.
	//
	// Conversion rules:
	//   - []int - contains indices of selected list items. Stored as is.
	//   - int - contains index of one selected list item, converted to []int.
	//   - string - contains one or several indices of selected list items separated by comma(,).
	Checked PropertyName = "checked"

	// ItemVerticalAlign is the constant for "item-vertical-align" property tag.
	//
	// Used by ListView.
	// Sets the vertical alignment of the contents of the list items.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Height alignment.
	ItemVerticalAlign PropertyName = "item-vertical-align"

	// ItemHorizontalAlign is the constant for "item-horizontal-align" property tag.
	//
	// Used by ListView.
	// Sets the horizontal alignment of the contents of the list items.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Height alignment.
	ItemHorizontalAlign PropertyName = "item-horizontal-align"

	// ItemCheckbox is the constant for "checkbox" property tag.
	//
	// Used by ListView.
	// Style of checkbox used to mark items in a list. Default value is "none".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneCheckbox) or "none" - There is no checkbox.
	//   - 1 (SingleCheckbox) or "single" - A checkbox that allows you to mark only one item, example: ◉.
	//   - 2 (MultipleCheckbox) or "multiple" - A checkbox that allows you to mark several items, example: ☑.
	ItemCheckbox PropertyName = "checkbox"

	// CheckboxHorizontalAlign is the constant for "checkbox-horizontal-align" property tag.
	//
	// Used by Checkbox, ListView.
	//
	// # Usage in Checkbox
	//
	// Horizontal alignment of checkbox inside the checkbox container.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Checkbox on the left edge, content on the right.
	//   - 1 (RightAlign) or "right" - Checkbox on the right edge, content on the left.
	//   - 2 (CenterAlign) or "center" - Center horizontally. Content below or above.
	//
	// # Usage in ListView
	//
	// Checkbox horizontal alignment(if enabled by "checkbox" property).
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Checkbox on the left edge, content on the right.
	//   - 1 (RightAlign) or "right" - Checkbox on the right edge, content on the left.
	//   - 2 (CenterAlign) or "center" - Center horizontally. Content below or above.
	CheckboxHorizontalAlign PropertyName = "checkbox-horizontal-align"

	// CheckboxVerticalAlign is the constant for "checkbox-vertical-align" property tag.
	//
	// Used by Checkbox, ListView.
	//
	// # Usage in Checkbox
	//
	// Vertical alignment of checkbox inside the checkbox container.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Checkbox on the top, content on the bottom.
	//   - 1 (BottomAlign) or "bottom" - Checkbox on the bottom, content on the top.
	//   - 2 (CenterAlign) or "center" - Checkbox on the top, content on the bottom.
	//
	// # Usage in ListView
	//
	// Checkbox vertical alignment(if enabled by "checkbox" property).
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign) or "top" - Top alignment.
	//   - 1 (BottomAlign) or "bottom" - Bottom alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	CheckboxVerticalAlign PropertyName = "checkbox-vertical-align"

	// NotTranslate is the constant for "not-translate" property tag.
	//
	// Used by TextView, View.
	//
	// Controls whether the text set for the text view require translation. This is an inherited property, i.e. if it is not
	// defined, then the value of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - No need to lookup for text translation in resources.
	//   - false, 0, "false", "no", "off", or "0" - Lookup for text translation.
	NotTranslate PropertyName = "not-translate"

	// Filter is the constant for "filter" property tag.
	//
	// Used by View.
	// Applies graphical effects to a view, such as blurring, color shifting, changing brightness/contrast, etc.
	//
	// Supported types: ViewFilter.
	//
	// See ViewFilter description for more details.
	Filter PropertyName = "filter"

	// BackdropFilter is the constant for "backdrop-filter" property tag.
	//
	// Used by View.
	// Applies graphical effects to the area behind a view, such as blurring, color shifting, changing brightness/contrast,
	// etc.
	//
	// Supported types: ViewFilter.
	//
	// See ViewFilter description for more details.
	BackdropFilter PropertyName = "backdrop-filter"

	// Clip is the constant for "clip" property tag.
	//
	// Used by View.
	// Creates a clipping region that sets what part of a view should be shown.
	//
	// Supported types: ClipShapeProperty, string.
	//
	// Internal type is ClipShapeProperty, other types converted to it during assignment.
	// See ClipShapeProperty description for more details.
	Clip PropertyName = "clip"

	// Points is the constant for "points" property tag.
	//
	// Used by ClipShapeProperty.
	// Points which describe polygon clip area. Values are in a sequence of pair like: x1, y1, x2, y2 ...
	//
	// Supported types: []SizeUnit, string.
	Points PropertyName = "points"

	// ShapeOutside is the constant for "shape-outside" property tag.
	//
	// Used by View.
	// __WARNING__ Currently not supported. Property defines a shape(which may be non-rectangular) around which adjacent
	// inline content should wrap. By default, inline content wraps around its margin box. Property provides a way to
	// customize this wrapping, making it possible to wrap text around complex objects rather than simple boxes.
	//
	// Supported types: ClipShapeProperty, string.
	//
	// Internal type is ClipShapeProperty, other types converted to it during assignment.
	// See ClipShapeProperty description for more details.
	ShapeOutside PropertyName = "shape-outside"

	// Float is the constant for "float" property tag.
	//
	// Used by View.
	// Places a view on the left or right side of its container, allowing text and inline views to wrap around it.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneFloat) or "none" - Text and other views inside the container will not wrap around this view.
	//   - 1 (LeftFloat) or "left" - Text and other views inside the container will wrap around this view on the right side.
	//   - 2 (RightFloat) or "right" - Text and other views inside the container will wrap around this view on the left side.
	Float PropertyName = "float"

	// UserData is the constant for "user-data" property tag.
	//
	// Used by View.
	// Can contain any user data.
	//
	// Supported types: any.
	UserData PropertyName = "user-data"

	// Resize is the constant for "resize" property tag.
	//
	// Used by View.
	// Sets whether view is resizable, and if so, in which directions. Default value is "none".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneResize) or "none" - View cannot be resized.
	//   - 1 (BothResize) or "both" - The View displays a mechanism for allowing the user to resize it, which may be resized both horizontally and vertically.
	//   - 2 (HorizontalResize) or "horizontal" - The View displays a mechanism for allowing the user to resize it in the horizontal direction.
	//   - 3 (VerticalResize) or "vertical" - The View displays a mechanism for allowing the user to resize it in the vertical direction.
	Resize PropertyName = "resize"

	// UserSelect is the constant for "user-select" property tag.
	//
	// Used by View.
	// Controls whether the user can select the text. This is an inherited property, i.e. if it is not defined, then the value
	// of the parent view is used. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - User can select the text.
	//   - false, 0, "false", "no", "off", or "0" - Text is not selectable.
	UserSelect PropertyName = "user-select"

	// Order is the constant for "Order" property tag.
	//
	// Used by View.
	//
	// Set the order to layout an item in a ViewsContainer container. Items in a container are sorted by
	// ascending order value and then by their addition to container order.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - negative value - Views with lower value will be at the beginning.
	//   - not negative value - Views with higher value will be at the end.
	Order PropertyName = "Order"

	// BackgroundBlendMode is the constant for "background-blend-mode" property tag.
	//
	// Used by View.
	// Sets how view's background images should blend with each other and with the view's background color.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (BlendNormal) or "normal" - The final color is the top color, regardless of what the bottom color is.
	//    The effect is like two opaque pieces of paper overlapping.
	//   - 1 (BlendMultiply) or "multiply" - The final color is the result of multiplying the top and bottom colors.
	//    A black layer leads to a black final layer, and a white layer leads to no change.
	//    The effect is like two images printed on transparent film overlapping.
	//   - 2 (BlendScreen) or "screen" - The final color is the result of inverting the colors, multiplying them,
	//    and inverting that value. A black layer leads to no change, and a white layer leads to a white final layer.
	//    The effect is like two images shone onto a projection screen.
	//   - 3 (BlendOverlay) or "overlay" - The final color is the result of multiply if the bottom color is darker,
	//    or screen if the bottom color is lighter. This blend mode is equivalent to hard-light but with the layers swapped.
	//   - 4 (BlendDarken) or "darken" - The final color is composed of the darkest values of each color channel.
	//   - 5 (BlendLighten) or "lighten" - The final color is composed of the lightest values of each color channel.
	//   - 6 (BlendColorDodge) or "color-dodge" - The final color is the result of dividing the bottom color by the inverse of the top color.
	//    A black foreground leads to no change. A foreground with the inverse color of the backdrop leads to a fully lit color.
	//    This blend mode is similar to screen, but the foreground need only be as light as the inverse of the backdrop to create a fully lit color.
	//   - 7 (BlendColorBurn) or "color-burn" - The final color is the result of inverting the bottom color, dividing the value by the top color,
	//    and inverting that value. A white foreground leads to no change. A foreground with the inverse color of the backdrop leads to a black final image.
	//    This blend mode is similar to multiply, but the foreground need only be as dark as the inverse of the backdrop to make the final image black.
	//   - 8 (BlendHardLight) or "hard-light" - The final color is the result of multiply if the top color is darker, or screen if the top color is lighter.
	//    This blend mode is equivalent to overlay but with the layers swapped. The effect is similar to shining a harsh spotlight on the backdrop.
	//   - 9 (BlendSoftLight) or "soft-light" - The final color is similar to hard-light, but softer. This blend mode behaves similar to hard-light.
	//    The effect is similar to shining a diffused spotlight on the backdrop.
	//   - 10 (BlendDifference) or "difference" - The final color is the result of subtracting the darker of the two colors from the lighter one.
	//    A black layer has no effect, while a white layer inverts the other layer's color.
	//   - 11 (BlendExclusion) or "exclusion" - The final color is similar to difference, but with less contrast.
	//    As with difference, a black layer has no effect, while a white layer inverts the other layer's color.
	//   - 12 (BlendHue) or "hue" - The final color has the hue of the top color, while using the saturation and luminosity of the bottom color.
	//   - 13 (BlendSaturation) or "saturation" - The final color has the saturation of the top color, while using the hue and luminosity of the bottom color.
	//    A pure gray backdrop, having no saturation, will have no effect.
	//   - 14 (BlendColor) or "color" - The final color has the hue and saturation of the top color, while using the luminosity of the bottom color.
	//    The effect preserves gray levels and can be used to colorize the foreground.
	//   - 15 (BlendLuminosity) or "luminosity" - The final color has the luminosity of the top color, while using the hue and saturation of the bottom color.
	//    This blend mode is equivalent to BlendColor, but with the layers swapped.
	BackgroundBlendMode PropertyName = "background-blend-mode"

	// MixBlendMode is the constant for "mix-blend-mode" property tag.
	//
	// Used by View.
	// Sets how view's content should blend with the content of the view's parent and the view's background.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (BlendNormal) or "normal" - The final color is the top color, regardless of what the bottom color is.
	//    The effect is like two opaque pieces of paper overlapping.
	//   - 1 (BlendMultiply) or "multiply" - The final color is the result of multiplying the top and bottom colors.
	//    A black layer leads to a black final layer, and a white layer leads to no change.
	//    The effect is like two images printed on transparent film overlapping.
	//   - 2 (BlendScreen) or "screen" - The final color is the result of inverting the colors, multiplying them,
	//    and inverting that value. A black layer leads to no change, and a white layer leads to a white final layer.
	//    The effect is like two images shone onto a projection screen.
	//   - 3 (BlendOverlay) or "overlay" - The final color is the result of multiply if the bottom color is darker,
	//    or screen if the bottom color is lighter. This blend mode is equivalent to hard-light but with the layers swapped.
	//   - 4 (BlendDarken) or "darken" - The final color is composed of the darkest values of each color channel.
	//   - 5 (BlendLighten) or "lighten" - The final color is composed of the lightest values of each color channel.
	//   - 6 (BlendColorDodge) or "color-dodge" - The final color is the result of dividing the bottom color by the inverse of the top color.
	//    A black foreground leads to no change. A foreground with the inverse color of the backdrop leads to a fully lit color.
	//    This blend mode is similar to screen, but the foreground need only be as light as the inverse of the backdrop to create a fully lit color.
	//   - 7 (BlendColorBurn) or "color-burn" - The final color is the result of inverting the bottom color, dividing the value by the top color,
	//    and inverting that value. A white foreground leads to no change. A foreground with the inverse color of the backdrop leads to a black final image.
	//    This blend mode is similar to multiply, but the foreground need only be as dark as the inverse of the backdrop to make the final image black.
	//   - 8 (BlendHardLight) or "hard-light" - The final color is the result of multiply if the top color is darker, or screen if the top color is lighter.
	//    This blend mode is equivalent to overlay but with the layers swapped. The effect is similar to shining a harsh spotlight on the backdrop.
	//   - 9 (BlendSoftLight) or "soft-light" - The final color is similar to hard-light, but softer. This blend mode behaves similar to hard-light.
	//    The effect is similar to shining a diffused spotlight on the backdrop.
	//   - 10 (BlendDifference) or "difference" - The final color is the result of subtracting the darker of the two colors from the lighter one.
	//    A black layer has no effect, while a white layer inverts the other layer's color.
	//   - 11 (BlendExclusion) or "exclusion" - The final color is similar to difference, but with less contrast.
	//    As with difference, a black layer has no effect, while a white layer inverts the other layer's color.
	//   - 12 (BlendHue) or "hue" - The final color has the hue of the top color, while using the saturation and luminosity of the bottom color.
	//   - 13 (BlendSaturation) or "saturation" - The final color has the saturation of the top color, while using the hue and luminosity of the bottom color.
	//    A pure gray backdrop, having no saturation, will have no effect.
	//   - 14 (BlendColor) or "color" - The final color has the hue and saturation of the top color, while using the luminosity of the bottom color.
	//    The effect preserves gray levels and can be used to colorize the foreground.
	//   - 15 (BlendLuminosity) or "luminosity" - The final color has the luminosity of the top color, while using the hue and saturation of the bottom color.
	//    This blend mode is equivalent to BlendColor, but with the layers swapped.
	MixBlendMode PropertyName = "mix-blend-mode"

	// TabIndex is the constant for "tabindex" property tag.
	//
	// Used by View.
	// Indicates that view can be focused, and where it participates in sequential keyboard navigation(usually with the Tab
	// key).
	//
	// Supported types: int, string.
	//
	// Values:
	//   - negative value - View can be selected with the mouse or touch, but does not participate in sequential navigation.
	//   - 0 - View can be selected and reached using sequential navigation, the order of navigation is determined by the browser(usually in order of addition).
	//   - positive value - View will be reached(and selected) using sequential navigation, and navigation is performed by ascending "tabindex" value.
	TabIndex PropertyName = "tabindex"

	// Tooltip is the constant for "tooltip" property tag.
	//
	// Used by View.
	// Specifies the tooltip text. Tooltip pops up when hovering the mouse cursor over the view. HTML tags are supported when
	// formatting tooltip text.
	//
	// Supported types: string.
	Tooltip PropertyName = "tooltip"
)
