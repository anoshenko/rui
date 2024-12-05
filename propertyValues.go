package rui

// Constants for various specific properties of a views
const (
	// Visible - default value of the view Visibility property: View is visible
	Visible = 0

	// Invisible - value of the view Visibility property: View is invisible but takes place
	Invisible = 1

	// Gone - value of the view Visibility property: View is invisible and does not take place
	Gone = 2

	// OverflowHidden - value of the view "overflow" property:
	// Content is clipped if necessary to fit the padding box. No scrollbars are provided,
	// and no support for allowing the user to scroll (such as by dragging or using a scroll wheel) is allowed.
	// The content can be scrolled programmatically, so the element is still a scroll container.
	OverflowHidden = 0

	// OverflowVisible - value of the view "overflow" property:
	// Content is not clipped and may be rendered outside the padding box.
	OverflowVisible = 1

	// OverflowScroll - value of the view "overflow" property:
	// Content is clipped if necessary to fit the padding box. Browsers always display scrollbars whether or
	// not any content is actually clipped, preventing scrollbars from appearing or disappearing as content changes.
	OverflowScroll = 2

	// OverflowAuto - value of the view "overflow" property:
	// Depends on the browser user agent. If content fits inside the padding box, it looks the same as OverflowVisible,
	// but still establishes a new block formatting context. Desktop browsers provide scrollbars if content overflows.
	OverflowAuto = 3

	// NoneTextTransform - not transform text
	NoneTextTransform = 0

	// CapitalizeTextTransform - capitalize text
	CapitalizeTextTransform = 1

	// LowerCaseTextTransform - transform text to lower case
	LowerCaseTextTransform = 2

	// UpperCaseTextTransform - transform text to upper case
	UpperCaseTextTransform = 3

	// HorizontalTopToBottom - content flows horizontally from left to right, vertically from top to bottom.
	// The next horizontal line is positioned below the previous line.
	HorizontalTopToBottom = 0

	// HorizontalBottomToTop - content flows horizontally from left to right, vertically from bottom to top.
	// The next horizontal line is positioned above the previous line.
	HorizontalBottomToTop = 1

	// VerticalRightToLeft - content flows vertically from top to bottom, horizontally from right to left.
	// The next vertical line is positioned to the left of the previous line.
	VerticalRightToLeft = 2

	// VerticalLeftToRight - content flows vertically from top to bottom, horizontally from left to right.
	// The next vertical line is positioned to the right of the previous line.
	VerticalLeftToRight = 3

	// MixedTextOrientation - rotates the characters of horizontal scripts 90° clockwise.
	// Lays out the characters of vertical scripts naturally. Default value.
	MixedTextOrientation = 0

	// UprightTextOrientation - lays out the characters of horizontal scripts naturally (upright),
	// as well as the glyphs for vertical scripts. Note that this keyword causes all characters
	// to be considered as left-to-right: the used value of "text-direction" is forced to be "left-to-right".
	UprightTextOrientation = 1

	// SystemTextDirection - direction of a text and other elements defined by system. This is the default value.
	SystemTextDirection = 0

	// LeftToRightDirection - text and other elements go from left to right.
	LeftToRightDirection = 1

	//RightToLeftDirection - text and other elements go from right to left.
	RightToLeftDirection = 2

	// ThinFont - the value of "text-weight" property: the thin (hairline) text weight
	ThinFont = 1

	// ExtraLightFont - the value of "text-weight" property: the extra light (ultra light) text weight
	ExtraLightFont = 2

	// LightFont - the value of "text-weight" property: the light text weight
	LightFont = 3

	// NormalFont - the value of "text-weight" property (default value): the normal text weight
	NormalFont = 4

	// MediumFont - the value of "text-weight" property: the medium text weight
	MediumFont = 5

	// SemiBoldFont - the value of "text-weight" property: the semi bold (demi bold) text weight
	SemiBoldFont = 6

	// BoldFont - the value of "text-weight" property: the bold text weight
	BoldFont = 7

	// ExtraBoldFont - the value of "text-weight" property: the extra bold (ultra bold) text weight
	ExtraBoldFont = 8

	// BlackFont - the value of "text-weight" property: the black (heavy) text weight
	BlackFont = 9

	// TopAlign - top vertical-align for the "vertical-align" property
	TopAlign = 0

	// BottomAlign - bottom vertical-align for the "vertical-align" property
	BottomAlign = 1

	// LeftAlign - the left horizontal-align for the "horizontal-align" property
	LeftAlign = 0

	// RightAlign - the right horizontal-align for the "horizontal-align" property
	RightAlign = 1

	// CenterAlign - the center horizontal/vertical-align for the "horizontal-align"/"vertical-align" property
	CenterAlign = 2

	// StretchAlign - the stretch horizontal/vertical-align for the "horizontal-align"/"vertical-align" property
	StretchAlign = 3

	// JustifyAlign - the justify text align for "text-align" property
	JustifyAlign = 3

	// BaselineAlign - the baseline cell-vertical-align for the "cell-vertical-align" property
	BaselineAlign = 4

	// WhiteSpaceNormal - sequences of white space are collapsed. Newline characters in the source
	// are handled the same as other white space. Lines are broken as necessary to fill line boxes.
	WhiteSpaceNormal = 0

	// WhiteSpaceNowrap - collapses white space as for normal, but suppresses line breaks (text wrapping)
	// within the source.
	WhiteSpaceNowrap = 1

	// WhiteSpacePre - sequences of white space are preserved. Lines are only broken at newline
	// characters in the source and at <br> elements.
	WhiteSpacePre = 2

	// WhiteSpacePreWrap - Sequences of white space are preserved. Lines are broken at newline
	// characters, at <br>, and as necessary to fill line boxes.
	WhiteSpacePreWrap = 3

	// WhiteSpacePreLine - sequences of white space are collapsed. Lines are broken at newline characters,
	// at <br>, and as necessary to fill line boxes.
	WhiteSpacePreLine = 4

	// WhiteSpaceBreakSpaces - the behavior is identical to that of WhiteSpacePreWrap, except that:
	//   - Any sequence of preserved white space always takes up space, including at the end of the line.
	//   - A line breaking opportunity exists after every preserved white space character, including between white space characters.
	//   - Such preserved spaces take up space and do not hang, and thus affect the box’s intrinsic sizes (min-content size and max-content size).
	WhiteSpaceBreakSpaces = 5

	// WordBreakNormal - use the default line break rule.
	WordBreakNormal = 0

	// WordBreakAll - to prevent overflow, word breaks should be inserted between any two characters
	// (excluding Chinese/Japanese/Korean text).
	WordBreakAll = 1

	// WordBreakKeepAll - word breaks should not be used for Chinese/Japanese/Korean (CJK) text.
	// Non-CJK text behavior is the same as for normal.
	WordBreakKeepAll = 2

	// WordBreakWord - when the block boundaries are exceeded, the remaining whole words can be split
	// in an arbitrary place, unless a more suitable place for the line break is found.
	WordBreakWord = 3

	// TextOverflowClip - truncate the text at the limit of the content area, therefore the truncation
	// can happen in the middle of a character.
	TextOverflowClip = 0

	// TextOverflowEllipsis - display an ellipsis ('…', U+2026 HORIZONTAL ELLIPSIS) to represent clipped text.
	// The ellipsis is displayed inside the content area, decreasing the amount of text displayed.
	// If there is not enough space to display the ellipsis, it is clipped.
	TextOverflowEllipsis = 1

	// DefaultSemantics - default value of the view Semantic property
	DefaultSemantics = 0

	// ArticleSemantics - value of the view Semantic property: view represents a self-contained
	// composition in a document, page, application, or site, which is intended to be
	// independently distributable or reusable (e.g., in syndication)
	ArticleSemantics = 1

	// SectionSemantics - value of the view Semantic property: view represents
	// a generic standalone section of a document, which doesn't have a more specific
	// semantic element to represent it.
	SectionSemantics = 2

	// AsideSemantics - value of the view Semantic property: view represents a portion
	// of a document whose content is only indirectly related to the document's main content.
	// Asides are frequently presented as sidebars or call-out boxes.
	AsideSemantics = 3

	// HeaderSemantics - value of the view Semantic property: view represents introductory
	// content, typically a group of introductory or navigational aids. It may contain
	// some heading elements but also a logo, a search form, an author name, and other elements.
	HeaderSemantics = 4

	// MainSemantics - value of the view Semantic property: view represents the dominant content
	// of the application. The main content area consists of content that is directly related
	// to or expands upon the central topic of a document, or the central functionality of an application.
	MainSemantics = 5

	// FooterSemantics - value of the view Semantic property: view represents a footer for its
	// nearest sectioning content or sectioning root element. A footer view typically contains
	// information about the author of the section, copyright data or links to related documents.
	FooterSemantics = 6

	// NavigationSemantics - value of the view Semantic property: view represents a section of
	// a page whose purpose is to provide navigation links, either within the current document
	// or to other documents. Common examples of navigation sections are menus, tables of contents,
	// and indexes.
	NavigationSemantics = 7

	// FigureSemantics - value of the view Semantic property: view represents self-contained content,
	// potentially with an optional caption, which is specified using the FigureCaptionSemantics view.
	FigureSemantics = 8

	// FigureCaptionSemantics - value of the view Semantic property: view represents a caption or
	// legend describing the rest of the contents of its parent FigureSemantics view.
	FigureCaptionSemantics = 9

	// ButtonSemantics - value of the view Semantic property: view a clickable button
	ButtonSemantics = 10

	// ParagraphSemantics - value of the view Semantic property: view represents a paragraph.
	// Paragraphs are usually represented in visual media as blocks of text separated
	// from adjacent blocks by blank lines and/or first-line indentation
	ParagraphSemantics = 11

	// H1Semantics - value of the view Semantic property: view represent of first level section headings.
	// H1Semantics is the highest section level and H6Semantics is the lowest.
	H1Semantics = 12

	// H2Semantics - value of the view Semantic property: view represent of second level section headings.
	// H1Semantics is the highest section level and H6Semantics is the lowest.
	H2Semantics = 13

	// H3Semantics - value of the view Semantic property: view represent of third level section headings.
	// H1Semantics is the highest section level and H6Semantics is the lowest.
	H3Semantics = 14

	// H4Semantics - value of the view Semantic property: view represent of fourth level section headings.
	// H1Semantics is the highest section level and H6Semantics is the lowest.
	H4Semantics = 15

	// H5Semantics - value of the view Semantic property: view represent of fifth level section headings.
	// H1Semantics is the highest section level and H6Semantics is the lowest.
	H5Semantics = 16

	// H6Semantics - value of the view Semantic property: view represent of sixth level section headings.
	// H1Semantics is the highest section level and H6Semantics is the lowest.
	H6Semantics = 17

	// BlockquoteSemantics - value of the view Semantic property: view indicates that
	// the enclosed text is an extended quotation.
	BlockquoteSemantics = 18

	// CodeSemantics - value of the view Semantic property: view displays its contents styled
	// in a fashion intended to indicate that the text is a short fragment of computer code
	CodeSemantics = 19

	// NoneFloat - value of the view "float" property: the View must not float.
	NoneFloat = 0

	// LeftFloat - value of the view "float" property: the View must float on the left side of its containing block.
	LeftFloat = 1

	// RightFloat - value of the view "float" property: the View must float on the right side of its containing block.
	RightFloat = 2

	// NoneResize - value of the view "resize" property: the View The offers no user-controllable method for resizing it.
	NoneResize = 0

	// BothResize - value of the view "resize" property: the View displays a mechanism for allowing
	// the user to resize it, which may be resized both horizontally and vertically.
	BothResize = 1

	// HorizontalResize - value of the view "resize" property: the View displays a mechanism for allowing
	// the user to resize it in the horizontal direction.
	HorizontalResize = 2

	// VerticalResize - value of the view "resize" property: the View displays a mechanism for allowing
	// the user to resize it in the vertical direction.
	VerticalResize = 3

	// RowAutoFlow - value of the "grid-auto-flow" property of the GridLayout:
	// Views are placed by filling each row in turn, adding new rows as necessary.
	RowAutoFlow = 0

	// ColumnAutoFlow - value of the "grid-auto-flow" property of the GridLayout:
	// Views are placed by filling each column in turn, adding new columns as necessary.
	ColumnAutoFlow = 1

	// RowDenseAutoFlow - value of the "grid-auto-flow" property of the GridLayout:
	// Views are placed by filling each row, adding new rows as necessary.
	// "dense" packing algorithm attempts to fill in holes earlier in the grid, if smaller items come up later.
	// This may cause views to appear out-of-order, when doing so would fill in holes left by larger views.
	RowDenseAutoFlow = 2

	// ColumnDenseAutoFlow - value of the "grid-auto-flow" property of the GridLayout:
	// Views are placed by filling each column, adding new columns as necessary.
	// "dense" packing algorithm attempts to fill in holes earlier in the grid, if smaller items come up later.
	// This may cause views to appear out-of-order, when doing so would fill in holes left by larger views.
	ColumnDenseAutoFlow = 3

	// BlendNormal - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the top color, regardless of what the bottom color is.
	// The effect is like two opaque pieces of paper overlapping.
	BlendNormal = 0

	// BlendMultiply - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of multiplying the top and bottom colors.
	// A black layer leads to a black final layer, and a white layer leads to no change.
	// The effect is like two images printed on transparent film overlapping.
	BlendMultiply = 1

	// BlendScreen - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of inverting the colors, multiplying them, and inverting that value.
	// A black layer leads to no change, and a white layer leads to a white final layer.
	// The effect is like two images shone onto a projection screen.
	BlendScreen = 2

	// BlendOverlay - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of multiply if the bottom color is darker, or screen if the bottom color is lighter.
	// This blend mode is equivalent to hard-light but with the layers swapped.
	BlendOverlay = 3

	// BlendDarken - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is composed of the darkest values of each color channel.
	BlendDarken = 4

	// BlendLighten - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is composed of the lightest values of each color channel.
	BlendLighten = 5

	// BlendColorDodge - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of dividing the bottom color by the inverse of the top color.
	// A black foreground leads to no change. A foreground with the inverse color of the backdrop leads to a fully lit color.
	// This blend mode is similar to screen, but the foreground need only be as light as the inverse of the backdrop to create a fully lit color.
	BlendColorDodge = 6

	// BlendColorBurn - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of inverting the bottom color, dividing the value by the top color, and inverting that value.
	// A white foreground leads to no change. A foreground with the inverse color of the backdrop leads to a black final image.
	// This blend mode is similar to multiply, but the foreground need only be as dark as the inverse of the backdrop to make the final image black.
	BlendColorBurn = 7

	// BlendHardLight - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of multiply if the top color is darker, or screen if the top color is lighter.
	// This blend mode is equivalent to overlay but with the layers swapped. The effect is similar to shining a harsh spotlight on the backdrop.
	BlendHardLight = 8

	// BlendSoftLight - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is similar to hard-light, but softer. This blend mode behaves similar to hard-light.
	// The effect is similar to shining a diffused spotlight on the backdrop*.*
	BlendSoftLight = 9

	// BlendDifference - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is the result of subtracting the darker of the two colors from the lighter one.
	// A black layer has no effect, while a white layer inverts the other layer's color.
	BlendDifference = 10

	// BlendExclusion - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color is similar to difference, but with less contrast.
	// As with difference, a black layer has no effect, while a white layer inverts the other layer's color.
	BlendExclusion = 11

	// BlendHue - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color has the hue of the top color, while using the saturation and luminosity of the bottom color.
	BlendHue = 12

	// BlendSaturation - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color has the saturation of the top color, while using the hue and luminosity of the bottom color.
	// A pure gray backdrop, having no saturation, will have no effect.
	BlendSaturation = 13

	// BlendColor - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color has the hue and saturation of the top color, while using the luminosity of the bottom color.
	// The effect preserves gray levels and can be used to colorize the foreground.
	BlendColor = 14

	// BlendLuminosity - value of the "mix-blend-mode" and "background-blend-mode" property:
	// The final color has the luminosity of the top color, while using the hue and saturation of the bottom color.
	// This blend mode is equivalent to BlendColor, but with the layers swapped.
	BlendLuminosity = 15

	// ColumnFillBalance - value of the "column-fill" property: content is equally divided between columns.
	ColumnFillBalance = 0

	// ColumnFillAuto - value of the "column-fill" property:
	// Columns are filled sequentially. Content takes up only the room it needs, possibly resulting in some columns remaining empty.
	ColumnFillAuto = 1

	// TextWrapOn - value of the "text-wrap" property:
	// text is wrapped across lines at appropriate characters (for example spaces,
	// in languages like English that use space separators) to minimize overflow.
	TextWrapOn = 0

	// TextWrapOff - value of the "text-wrap" property: text does not wrap across lines.
	// It will overflow its containing element rather than breaking onto a new line.
	TextWrapOff = 1

	// TextWrapBalance - value of the "text-wrap" property: text is wrapped in a way
	// that best balances the number of characters on each line, enhancing layout quality
	// and legibility. Because counting characters and balancing them across multiple lines
	// is computationally expensive, this value is only supported for blocks of text
	// spanning a limited number of lines (six or less for Chromium and ten or less for Firefox).
	TextWrapBalance = 2
)
