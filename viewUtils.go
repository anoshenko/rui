package rui

// Get returns a value of the property with name "tag" of the "rootView" subview with "viewID" id value.
// The type of return value depends on the property.
// If the subview don't exists or the property is not set then nil is returned.
func Get(rootView View, viewID, tag string) any {
	var view View
	if viewID != "" {
		view = ViewByID(rootView, viewID)
	} else {
		view = rootView
	}
	if view != nil {
		return view.Get(tag)
	}
	return nil
}

// Set sets the property with name "tag" of the "rootView" subview with "viewID" id by value. Result:
// true - success,
// false - error (incompatible type or invalid format of a string value, see AppLog).
func Set(rootView View, viewID, tag string, value any) bool {
	var view View
	if viewID != "" {
		view = ViewByID(rootView, viewID)
	} else {
		view = rootView
	}
	if view != nil {
		return view.Set(tag, value)
	}
	return false
}

// SetChangeListener sets a listener for changing a subview property value.
// If the second argument (subviewID) is not specified or it is "" then a listener for the first argument (view) is set
func SetChangeListener(view View, viewID, tag string, listener func(View, string)) {
	if viewID != "" {
		view = ViewByID(view, viewID)
	}
	if view != nil {
		view.SetChangeListener(tag, listener)
	}
}

// SetParams sets properties with name "tag" of the "rootView" subview. Result:
// true - all properties were set successful,
// false - error (incompatible type or invalid format of a string value, see AppLog).
func SetParams(rootView View, viewID string, params Params) bool {
	if viewID != "" {
		rootView = ViewByID(rootView, viewID)
	}
	if rootView == nil {
		return false
	}

	session := rootView.Session()
	session.startUpdateScript(rootView.htmlID())
	result := true
	for tag, value := range params {
		result = rootView.Set(tag, value) && result
	}
	session.finishUpdateScript(rootView.htmlID())
	return result
}

// IsDisabled returns "true" if the subview is disabled
// If the second argument (subviewID) is not specified or it is "" then a state of the first argument (view) is returned
func IsDisabled(view View, subviewID ...string) bool {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if disabled, _ := boolProperty(view, Disabled, view.Session()); disabled {
			return true
		}
		if parent := view.Parent(); parent != nil {
			return IsDisabled(parent)
		}
	}
	return false
}

// GetSemantics returns the subview semantics.  Valid semantics values are
// DefaultSemantics (0), ArticleSemantics (1), SectionSemantics (2), AsideSemantics (3),
// HeaderSemantics (4), MainSemantics (5), FooterSemantics (6), NavigationSemantics (7),
// FigureSemantics (8), FigureCaptionSemantics (9), ButtonSemantics (10), ParagraphSemantics (11),
// H1Semantics (12) - H6Semantics (17), BlockquoteSemantics (18), and CodeSemantics (19).
// If the second argument (subviewID) is not specified or it is "" then a semantics of the first argument (view) is returned
func GetSemantics(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, Semantics, DefaultSemantics, false)
}

// GetOpacity returns the subview opacity.
// If the second argument (subviewID) is not specified or it is "" then an opacity of the first argument (view) is returned
func GetOpacity(view View, subviewID ...string) float64 {
	return floatStyledProperty(view, subviewID, Opacity, 1)
}

// GetStyle returns the subview style id.
// If the second argument (subviewID) is not specified or it is "" then a style of the first argument (view) is returned
func GetStyle(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if style, ok := stringProperty(view, Style, view.Session()); ok {
			return style
		}
	}
	return ""
}

// GetDisabledStyle returns the disabled subview style id.
// If the second argument (subviewID) is not specified or it is "" then a style of the first argument (view) is returned
func GetDisabledStyle(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if style, ok := stringProperty(view, StyleDisabled, view.Session()); ok {
			return style
		}
	}
	return ""
}

// GetVisibility returns the subview visibility. One of the following values is returned:
// Visible (0), Invisible (1), or Gone (2)
// If the second argument (subviewID) is not specified or it is "" then a visibility of the first argument (view) is returned
func GetVisibility(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, Visibility, Visible, false)
}

// GetOverflow returns a value of the subview "overflow" property. Returns one of next values:
// OverflowHidden (0), OverflowVisible (1), OverflowScroll (2), OverflowAuto (3)
// If the second argument (subviewID) is not specified or it is "" then a value of the first argument (view) is returned
func GetOverflow(view View, subviewID ...string) int {
	defaultOverflow := OverflowHidden
	view2 := view
	if len(subviewID) > 0 && subviewID[0] != "" {
		view2 = ViewByID(view, subviewID[0])
	}
	if view2 != nil {
		switch view.(type) {
		case EditView:
			defaultOverflow = OverflowAuto

		case ListView:
			defaultOverflow = OverflowAuto
		}
	}
	return enumStyledProperty(view, subviewID, Overflow, defaultOverflow, false)
}

// GetZIndex returns the subview z-order.
// If the second argument (subviewID) is not specified or it is "" then a z-order of the first argument (view) is returned
func GetZIndex(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, ZIndex, 0)
}

// GetWidth returns the subview width.
// If the second argument (subviewID) is not specified or it is "" then a width of the first argument (view) is returned
func GetWidth(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Width, false)
}

// GetHeight returns the subview height.
// If the second argument (subviewID) is not specified or it is "" then a height of the first argument (view) is returned
func GetHeight(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Height, false)
}

// GetMinWidth returns a minimal subview width.
// If the second argument (subviewID) is not specified or it is "" then a minimal width of the first argument (view) is returned
func GetMinWidth(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, MinWidth, false)
}

// GetMinHeight returns a minimal subview height.
// If the second argument (subviewID) is not specified or it is "" then a minimal height of the first argument (view) is returned
func GetMinHeight(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, MinHeight, false)
}

// GetMaxWidth returns a maximal subview width.
// If the second argument (subviewID) is not specified or it is "" then a maximal width of the first argument (view) is returned
func GetMaxWidth(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, MaxWidth, false)
}

// GetMaxHeight returns a maximal subview height.
// If the second argument (subviewID) is not specified or it is "" then a maximal height of the first argument (view) is returned
func GetMaxHeight(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, MaxHeight, false)
}

// GetResize returns the "resize" property value if the subview. One of the following values is returned:
// NoneResize (0), BothResize (1), HorizontalResize (2), or VerticalResize (3)
// If the second argument (subviewID) is not specified or it is "" then a value of the first argument (view) is returned
func GetResize(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, Resize, NoneResize, false)
}

// GetLeft returns a left position of the subview in an AbsoluteLayout container.
// If a parent view is not an AbsoluteLayout container then this value is ignored.
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetLeft(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Left, false)
}

// GetRight returns a right position of the subview in an AbsoluteLayout container.
// If a parent view is not an AbsoluteLayout container then this value is ignored.
// If the second argument (subviewID) is not specified or it is "" then a right position of the first argument (view) is returned
func GetRight(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Right, false)
}

// GetTop returns a top position of the subview in an AbsoluteLayout container.
// If a parent view is not an AbsoluteLayout container then this value is ignored.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetTop(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Top, false)
}

// GetBottom returns a top position of the subview in an AbsoluteLayout container.
// If a parent view is not an AbsoluteLayout container then this value is ignored.
// If the second argument (subviewID) is not specified or it is "" then a bottom position of the first argument (view) is returned
func GetBottom(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Bottom, false)
}

// Margin returns the subview margin.
// If the second argument (subviewID) is not specified or it is "" then a margin of the first argument (view) is returned
func GetMargin(view View, subviewID ...string) Bounds {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	var bounds Bounds
	if view != nil {
		bounds.setFromProperties(Margin, MarginTop, MarginRight, MarginBottom, MarginLeft, view, view.Session())
	}
	return bounds
}

// GetPadding returns the subview padding.
// If the second argument (subviewID) is not specified or it is "" then a padding of the first argument (view) is returned
func GetPadding(view View, subviewID ...string) Bounds {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	var bounds Bounds
	if view != nil {
		bounds.setFromProperties(Padding, PaddingTop, PaddingRight, PaddingBottom, PaddingLeft, view, view.Session())
	}
	return bounds
}

// GetBorder returns ViewBorders of the subview.
// If the second argument (subviewID) is not specified or it is "" then a ViewBorders of the first argument (view) is returned.
func GetBorder(view View, subviewID ...string) ViewBorders {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if border := getBorder(view, Border); border != nil {
			return border.ViewBorders(view.Session())
		}
	}
	return ViewBorders{}
}

// Radius returns the BoxRadius structure of the subview.
// If the second argument (subviewID) is not specified or it is "" then a BoxRadius of the first argument (view) is returned.
func GetRadius(view View, subviewID ...string) BoxRadius {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return BoxRadius{}
	}
	return getRadius(view, view.Session())
}

// GetOutline returns ViewOutline of the subview.
// If the second argument (subviewID) is not specified or it is "" then a ViewOutline of the first argument (view) is returned.
func GetOutline(view View, subviewID ...string) ViewOutline {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if outline := getOutline(view); outline != nil {
			return outline.ViewOutline(view.Session())
		}
	}
	return ViewOutline{Style: NoneLine, Width: AutoSize(), Color: 0}
}

// GetViewShadows returns shadows of the subview.
// If the second argument (subviewID) is not specified or it is "" then shadows of the first argument (view) is returned.
func GetViewShadows(view View, subviewID ...string) []ViewShadow {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return []ViewShadow{}
	}
	return getShadows(view, Shadow)
}

// GetTextShadows returns text shadows of the subview.
// If the second argument (subviewID) is not specified or it is "" then shadows of the first argument (view) is returned.
func GetTextShadows(view View, subviewID ...string) []ViewShadow {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return []ViewShadow{}
	}
	return getShadows(view, TextShadow)
}

// GetBackgroundColor returns a background color of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetBackgroundColor(view View, subviewID ...string) Color {
	return colorStyledProperty(view, subviewID, BackgroundColor, false)
}

// GetAccentColor returns the accent color for UI controls generated by some elements.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAccentColor(view View, subviewID ...string) Color {
	return colorStyledProperty(view, subviewID, AccentColor, false)
}

// GetFontName returns the subview font.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetFontName(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if font, ok := stringProperty(view, FontName, view.Session()); ok {
			return font
		}
		if value := valueFromStyle(view, FontName); value != nil {
			if font, ok := value.(string); ok {
				return font
			}
		}
		if parent := view.Parent(); parent != nil {
			return GetFontName(parent)
		}
	}
	return ""
}

// GetTextColor returns a text color of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextColor(view View, subviewID ...string) Color {
	return colorStyledProperty(view, subviewID, TextColor, true)
}

// GetTextSize returns a text size of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextSize(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, TextSize, true)
}

// GetTabSize returns the subview width of tab characters (U+0009) in spaces.
// If the second argument (subviewID) is not specified or it is "" then a width of the first argument (view) is returned
func GetTabSize(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, TabSize, 8)
}

// GetTextWeight returns a text weight of the subview. Returns one of next values:
// 1, 2, 3, 4 (normal text), 5, 6, 7 (bold text), 8 and 9
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextWeight(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TextWeight, NormalFont, true)
}

// GetTextAlign returns a text align of the subview. Returns one of next values:
//
//	LeftAlign = 0, RightAlign = 1, CenterAlign = 2, JustifyAlign = 3
//
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TextAlign, LeftAlign, true)
}

// GetTextIndent returns a text indent of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextIndent(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, TextIndent, true)
}

// GetLetterSpacing returns a letter spacing of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetLetterSpacing(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, LetterSpacing, true)
}

// GetWordSpacing returns a word spacing of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetWordSpacing(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, WordSpacing, true)
}

// GetLineHeight returns a height of a text line of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetLineHeight(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, LineHeight, true)
}

// IsItalic returns "true" if a text font of the subview is displayed in italics, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsItalic(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Italic, true)
}

// IsSmallCaps returns "true" if a text font of the subview is displayed in small caps, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsSmallCaps(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, SmallCaps, true)
}

// IsStrikethrough returns "true" if a text font of the subview is displayed strikethrough, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsStrikethrough(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Strikethrough, true)
}

// IsOverline returns "true" if a text font of the subview is displayed overlined, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsOverline(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Overline, true)
}

// IsUnderline returns "true" if a text font of the subview is displayed underlined, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsUnderline(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Underline, true)
}

// GetTextLineThickness returns the stroke thickness of the decoration line that
// is used on text in an element, such as a line-through, underline, or overline.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextLineThickness(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, TextLineThickness, true)
}

// GetTextLineStyle returns the stroke style of the decoration line that
// is used on text in an element, such as a line-through, underline, or overline.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextLineStyle(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TextLineStyle, SolidLine, true)
}

// GetTextLineColor returns the stroke color of the decoration line that
// is used on text in an element, such as a line-through, underline, or overline.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextLineColor(view View, subviewID ...string) Color {
	return colorStyledProperty(view, subviewID, TextLineColor, true)
}

// GetTextTransform returns a text transform of the subview. Return one of next values:
// NoneTextTransform (0), CapitalizeTextTransform (1), LowerCaseTextTransform (2) or UpperCaseTextTransform (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextTransform(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TextTransform, NoneTextTransform, true)
}

// GetWritingMode returns whether lines of text are laid out horizontally or vertically, as well as
// the direction in which blocks progress. Valid values are HorizontalTopToBottom (0),
// HorizontalBottomToTop (1), VerticalRightToLeft (2) and VerticalLeftToRight (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetWritingMode(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, WritingMode, HorizontalTopToBottom, true)
}

// GetTextDirection - returns a direction of text, table columns, and horizontal overflow.
// Valid values are SystemTextDirection (0), LeftToRightDirection (1), and RightToLeftDirection (2).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextDirection(view View, subviewID ...string) int {
	if view == nil {
		return SystemTextDirection
	}
	defaultDirection := view.Session().TextDirection()
	return enumStyledProperty(view, subviewID, TextDirection, defaultDirection, true)
}

// GetVerticalTextOrientation returns a orientation of the text characters in a line. It only affects text
// in vertical mode (when "writing-mode" is "vertical-right-to-left" or "vertical-left-to-right").
// Valid values are MixedTextOrientation (0), UprightTextOrientation (1), and SidewaysTextOrientation (2).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetVerticalTextOrientation(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, VerticalTextOrientation, MixedTextOrientation, true)
}

// GetRow returns the range of row numbers of a GridLayout in which the subview is placed.
// If the second argument (subviewID) is not specified or it is "" then a values from the first argument (view) is returned.
func GetRow(view View, subviewID ...string) Range {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		session := view.Session()
		if result, ok := rangeProperty(view, Row, session); ok {
			return result
		}
		if value := valueFromStyle(view, Row); value != nil {
			if result, ok := valueToRange(value, session); ok {
				return result
			}
		}
	}
	return Range{}
}

// GetColumn returns the range of column numbers of a GridLayout in which the subview is placed.
// If the second argument (subviewID) is not specified or it is "" then a values from the first argument (view) is returned.
func GetColumn(view View, subviewID ...string) Range {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		session := view.Session()
		if result, ok := rangeProperty(view, Column, session); ok {
			return result
		}
		if value := valueFromStyle(view, Column); value != nil {
			if result, ok := valueToRange(value, session); ok {
				return result
			}
		}
	}
	return Range{}
}

// GetPerspective returns a distance between the z = 0 plane and the user in order to give a 3D-positioned
// element some perspective. Each 3D element with z > 0 becomes larger; each 3D-element with z < 0 becomes smaller.
// The default value is 0 (no 3D effects).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPerspective(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Perspective, false)
}

// GetPerspectiveOrigin returns a x- and y-coordinate of the position at which the viewer is looking.
// It is used as the vanishing point by the Perspective property. The default value is (50%, 50%).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPerspectiveOrigin(view View, subviewID ...string) (SizeUnit, SizeUnit) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return AutoSize(), AutoSize()
	}
	return getPerspectiveOrigin(view, view.Session())
}

// GetBackfaceVisible returns a bool property that sets whether the back face of an element is
// visible when turned towards the user. Values:
// true - the back face is visible when turned towards the user (default value).
// false - the back face is hidden, effectively making the element invisible when turned away from the user.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetBackfaceVisible(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, BackfaceVisible, false)
}

// GetOrigin returns a x-, y-, and z-coordinate of the point around which a view transformation is applied.
// The default value is (50%, 50%, 50%).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetOrigin(view View, subviewID ...string) (SizeUnit, SizeUnit, SizeUnit) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return AutoSize(), AutoSize(), AutoSize()
	}
	return getOrigin(view, view.Session())
}

// GetTranslate returns a x-, y-, and z-axis translation value of a 2D/3D translation
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTranslate(view View, subviewID ...string) (SizeUnit, SizeUnit, SizeUnit) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return AutoSize(), AutoSize(), AutoSize()
	}
	return getTranslate(view, view.Session())
}

// GetSkew returns a angles to use to distort the element along the abscissa (x-axis)
// and the ordinate (y-axis). The default value is 0.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetSkew(view View, subviewID ...string) (AngleUnit, AngleUnit) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return AngleUnit{Value: 0, Type: Radian}, AngleUnit{Value: 0, Type: Radian}
	}
	x, y, _ := getSkew(view, view.Session())
	return x, y
}

// GetScale returns a x-, y-, and z-axis scaling value of a 2D/3D scale. The default value is 1.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetScale(view View, subviewID ...string) (float64, float64, float64) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return 1, 1, 1
	}

	session := view.Session()
	x, _ := floatProperty(view, ScaleX, session, 1)
	y, _ := floatProperty(view, ScaleY, session, 1)
	z, _ := floatProperty(view, ScaleZ, session, 1)
	return x, y, z
}

// GetRotate returns a x-, y, z-coordinate of the vector denoting the axis of rotation, and the angle of the view rotation
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetRotate(view View, subviewID ...string) (float64, float64, float64, AngleUnit) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return 0, 0, 0, AngleUnit{Value: 0, Type: Radian}
	}

	session := view.Session()
	angle, _ := angleProperty(view, Rotate, view.Session())
	rotateX, _ := floatProperty(view, RotateX, session, 1)
	rotateY, _ := floatProperty(view, RotateY, session, 1)
	rotateZ, _ := floatProperty(view, RotateZ, session, 1)
	return rotateX, rotateY, rotateZ, angle
}

// GetAvoidBreak returns "true" if avoids any break from being inserted within the principal box,
// and "false" if allows, but does not force, any break to be inserted within the principal box.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetAvoidBreak(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, AvoidBreak, true)
}

func GetNotTranslate(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, NotTranslate, true)
}

func valueFromStyle(view View, tag string) any {
	session := view.Session()
	getValue := func(styleTag string) any {
		if style, ok := stringProperty(view, styleTag, session); ok {
			if style, ok := session.resolveConstants(style); ok {
				return session.styleProperty(style, tag)
			}
		}
		return nil
	}

	if IsDisabled(view) {
		if value := getValue(StyleDisabled); value != nil {
			return value
		}
	}
	return getValue(Style)
}

func sizeStyledProperty(view View, subviewID []string, tag string, inherit bool) SizeUnit {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value, ok := sizeProperty(view, tag, view.Session()); ok {
			return value
		}
		if value := valueFromStyle(view, tag); value != nil {
			if result, ok := valueToSizeUnit(value, view.Session()); ok {
				return result
			}
		}

		if inherit {
			if parent := view.Parent(); parent != nil {
				return sizeStyledProperty(parent, []string{}, tag, true)
			}
		}
	}
	return AutoSize()
}

func enumStyledProperty(view View, subviewID []string, tag string, defaultValue int, inherit bool) int {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value, ok := enumProperty(view, tag, view.Session(), defaultValue); ok {
			return value
		}
		if value := valueFromStyle(view, tag); value != nil {
			if result, ok := valueToEnum(value, tag, view.Session(), defaultValue); ok {
				return result
			}
		}

		if inherit {
			if parent := view.Parent(); parent != nil {
				return enumStyledProperty(parent, []string{}, tag, defaultValue, true)
			}
		}
	}
	return defaultValue
}

func boolStyledProperty(view View, subviewID []string, tag string, inherit bool) bool {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value, ok := boolProperty(view, tag, view.Session()); ok {
			return value
		}
		if value := valueFromStyle(view, tag); value != nil {
			if b, ok := valueToBool(value, view.Session()); ok {
				return b
			}
		}

		if inherit {
			if parent := view.Parent(); parent != nil {
				return boolStyledProperty(parent, []string{}, tag, inherit)
			}
		}
	}

	return false
}

func intStyledProperty(view View, subviewID []string, tag string, defaultValue int) int {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value, ok := intProperty(view, tag, view.Session(), defaultValue); ok {
			return value
		}
		if value := valueFromStyle(view, tag); value != nil {
			n, _ := valueToInt(value, view.Session(), defaultValue)
			return n
		}
	}
	return defaultValue
}

func floatStyledProperty(view View, subviewID []string, tag string, defaultValue float64) float64 {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value, ok := floatProperty(view, tag, view.Session(), defaultValue); ok {
			return value
		}
		if value := valueFromStyle(view, tag); value != nil {
			f, _ := valueToFloat(value, view.Session(), defaultValue)
			return f
		}
	}
	return defaultValue
}

func colorStyledProperty(view View, subviewID []string, tag string, inherit bool) Color {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value, ok := colorProperty(view, tag, view.Session()); ok {
			return value
		}
		if value := valueFromStyle(view, tag); value != nil {
			if color, ok := valueToColor(value, view.Session()); ok {
				return color
			}
		}
		if inherit {
			if parent := view.Parent(); parent != nil {
				return colorStyledProperty(parent, []string{}, tag, true)
			}
		}
	}
	return Color(0)
}

// FocusView sets focus on the specified View, if it can be focused.
// The focused View is the View which will receive keyboard events by default.
func FocusView(view View) {
	if view != nil {
		view.Session().runFunc("focus", view.htmlID())
	}
}

// FocusView sets focus on the View with the specified viewID, if it can be focused.
// The focused View is the View which will receive keyboard events by default.
func FocusViewByID(viewID string, session Session) {
	if viewID != "" {
		session.runFunc("focus", viewID)
	}
}

// BlurView removes keyboard focus from the specified View.
func BlurView(view View) {
	if view != nil {
		view.Session().runFunc("blur", view.htmlID())
	}
}

// BlurViewByID removes keyboard focus from the View with the specified viewID.
func BlurViewByID(viewID string, session Session) {
	if viewID != "" {
		session.runFunc("blur", viewID)
	}
}

// GetCurrent returns the index of the selected item (<0 if there is no a selected item) or the current view index (StackLayout, TabsLayout).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCurrent(view View, subviewID ...string) int {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	defaultValue := -1
	if view != nil {
		if result, ok := intProperty(view, Current, view.Session(), defaultValue); ok {
			return result
		} else if view.Tag() != "ListView" {
			defaultValue = 0
		}
	}
	return defaultValue
}

// IsUserSelect returns "true" if the user can select text, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsUserSelect(view View, subviewID ...string) bool {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		value, _ := isUserSelect(view)
		return value
	}

	return false
}

func isUserSelect(view View) (bool, bool) {
	if value, ok := boolProperty(view, UserSelect, view.Session()); ok {
		return value, true
	}
	if value := valueFromStyle(view, UserSelect); value != nil {
		if b, ok := valueToBool(value, view.Session()); ok {
			return b, true
		}
	}

	if parent := view.Parent(); parent != nil {
		if result, ok := isUserSelect(parent); ok {
			return result, true
		}
	}

	switch GetSemantics(view) {
	case ParagraphSemantics, H1Semantics, H2Semantics, H3Semantics, H4Semantics, H5Semantics,
		H6Semantics, BlockquoteSemantics, CodeSemantics:
		return true, false
	}

	if _, ok := view.(TableView); ok {
		return true, false
	}

	return false, false
}
