# v0.20.0

* Added support of binding
* Added "binding" argument to CreateViewFromResources, CreateViewFromText, and CreateViewFromObject functions
* Added CreatePopupFromResources, CreatePopupFromText, and CreatePopupFromObject functions
* Added implementation of Properties interface to Popup
* Changed ParseDataText function return values

# v0.19.0

* Added support of drag-and-drop
* Added LoadFile method to View interface

# v0.18.2

* fixed typo: GetShadowProperties -> GetShadowProperty

# v0.18.0

* Property name type changed from string to PropertyName. 

* Renamed:
    Transform interface -> TransformProperty
    NewTransform function -> NewTransformProperty
    TransformTag constant -> Transform.
    "origin-x" property -> "transform-origin-x"
    "origin-y" property -> "transform-origin-y"
    "origin-z" property -> "transform-origin-z"
    GetOrigin function -> GetTransformOrigin.
    BorderBoxClip constant -> BorderBox
    PaddingBoxClip constant -> PaddingBox
    ContentBoxClip constant -> ContentBox.
    ViewShadow interface -> ShadowProperty
    NewViewShadow function -> NewShadow
    NewInsetViewShadow function -> NewInsetShadow
    NewShadowWithParams function -> NewShadowProperty
    NewColumnSeparator function -> NewColumnSeparatorProperty
    ClipShape interface -> ClipShapeProperty
    InsetClip function -> NewInsetClip
    CircleClip function -> NewCircleClip
    EllipseClip function -> NewEllipseClip
    PolygonClip function -> NewPolygonClip
    PolygonPointsClip function -> NewPolygonPointsClip
    ViewFilter interface -> FilterProperty
    NewViewFilter function -> NewFilterProperty
    Animation interface -> AnimationProperty
    AnimationTag constant -> Animation
    NewAnimation function -> NewAnimationProperty

* Added functions: NewBounds, NewEllipticRadius, NewRadii, NewLinearGradient, NewCircleRadialGradient,
NewEllipseRadialGradient, GetPushTransform, GetPushDuration, GetPushTiming, IsMoveToFrontAnimation,
GetBackground, GetMask, GetBackgroundClip,GetBackgroundOrigin, GetMaskClip, GetMaskOrigin, NewColumnSeparator,
NewClipShapeProperty, NewTransitionAnimation, NewAnimation, IsSummaryMarkerHidden.

* Changed ViewByID functions

* Added SetConicGradientFillStyle and SetConicGradientStrokeStyle methods to Canvas interface.

* Changed Push, Pop, MoveToFront, and MoveToFrontByID methods of StackLayout interface.

* Removed DefaultAnimation, StartToEndAnimation, EndToStartAnimation, TopDownAnimation, and BottomUpAnimation constants.

* Added StackLayout properties: "push-transform", "push-duration", "push-timing", "move-to-front-animation", "push-perspective", 
"push-rotate-x", "push-rotate-y", "push-rotate-z", "push-rotate", "push-skew-x", "push-skew-y",
"push-scale-x", "push-scale-y", "push-scale-z", "push-translate-x", "push-translate-y", "push-translate-z".

* Added "show-opacity", "show-transform", "show-duration", and "show-timing" Popup properties.

* Added "mask", "mask-clip", "mask-origin", and "background-origin" properties.

* Added "hide-summary-marker" DetailsView property.

* Added LineJoin type. Type of constants MiterJoin, RoundJoin, and BevelJoin changed to LineJoin. Type of Canvas.SetLineJoin function argument changed to LineJoin.

* Added LineCap type. Type of constants ButtCap, RoundCap, and SquareCap changed to LineCap. Type of Canvas.SetLineCap function argument changed to LineCap.

# v0.17.3
Added SetParams method to View interface

# v0.17.0
* Added "mod", "rem", "round", "round-up", "round-down", and "round-to-zero" SizeFunc functions
* Added ModSize, RemSize, RoundSize, RoundUpSize, RoundDownSize, and RoundToZeroSize functions
* Added Start, Stop, Pause, and Resume methods to Animation interface
* Added "transform" property and Transform interface
* Added OpenRawResource, GetCheckboxChangedListeners functions
* Added RemoveClientItem method to Session interface
* Added "item-separators" property to DropDownList and GetDropDownItemSeparators function
* Added NewPath and NewPathFromSvg methods to Canvas interface
* Removed NewPath function
* Removed Reset methods from Path interface

# v0.16.0
* Can use ListAdapter as "content" property value of ListLayout
* The IsListItemEnabled method of the ListAdapter interface has been made optional 
* Can use GridAdapter as "content" property value of GridLayout
* Added "text-wrap" property and GetGetTextWrap function
* Bug fixing

# v0.15.0
* Added "data-list" property 
* Bug fixing

# v0.14.0
* Added the ability to work without creating a WebSocket. Added NoSocket property to AppParams.
* Added SocketAutoClose property to AppParams.
* Added the ability to run a timer on the client side. Added StartTimer and StopTimer methods to Session interface.
* Added "cell-vertical-self-align", and "cell-horizontal-self-align" properties 
* Bug fixing

# v0.13.x
* Added NewHandler function
* Bug fixing

# v0.13.0

* Added SetHotKey function to Session interface
* Added ViewIndex function to ViewsContainer interface
* Added ReloadCell function to TableView interface
* Added ReloadTableViewCell function
* Added "tooltip" property and GetTooltip function
* Added "outline-offset" property and GetOutlineOffset function
* Changed the main event listener format for "drop-down-event", "edit-text-changed", 
"color-changed", "number-changed", "date-changed", and "time-changed" events. 
Old format is "<listener>(<view>, <new value>)", new format is "<listener>(<view>, <new value>, <old value>)"
* Changed FocusView function
* Added support for height and width range in media styles. 
Changed MediaStyle, SetMediaStyle, and MediaStyles functions of Theme interface
* Bug fixing

# v0.12.0

* Added SvgImageView
* Added InlineImageFromResource function

# v0.11.0

* Added "tabindex", "order", "column-fill", "column-span-all", "background-blend-mode", and "mix-blend-mode" properties 
* Added GetTabIndex, GetOrder, GetColumnFill, IsColumnSpanAll, GetBackgroundBlendMode, and GetMixBlendMode functions
* ClientItem, SetClientItem, and RemoveAllClientItems method added to Session interface
* PropertyWithTag method of DataObject interface renamed to PropertyByTag

# v0.10.0

* The Canvas.TextWidth method replaced by Canvas.TextMetrics
* Added support of WebAssembly

# v0.9.0

* Requires go 1.18 or higher
* The "interface{}" type replaced by "any"
* Added SizeFunc interface and Function field to SizeUnit struct
* Added MaxSize, MinSize, SumSize, SubSize, MulSize, DivSize, ClampSize functions
* Added "list-row-gap", "list-column-gap", "accent-color", "tab-size", "overflow", 
"arrow", "arrow-align", "arrow-size", "arrow-width", and "arrow-offset" properties 
* Added "@ruiArrowSize" and "@ruiArrowWidth" constants to the default theme
* Added Transition, Transitions, and SetTransition functions to the ViewStyle interface
* Added GetListRowGap, GetListColumnGap, GetAccentColor, GetTabSize, GetOverflow, IsTimingFunctionValid, and GetTransitions functions
* Changed GetTransition functions
* Added the OpenURL function to the Session interface
* Changed the type of the second argument of all Get functions to "subviewID ...string" (previously "subviewID string")

# v0.8.0

* Added "loaded-event" and "error-event" events to ImageView
* Added NaturalSize and CurrentSource methods to ImageView
* Added "user-select" property and IsUserSelect function
* Renamed "LightGoldenrodYellow" color constant to "LightGoldenRodYellow"

# v0.7.0

* Added "resize", "grid-auto-flow", "caret-color", and "backdrop-filter" properties 
* Added BlurView, BlurViewByID, GetResize, GetGridAutoFlow, GetCaretColor, GetBackdropFilter functions
* The "warp" property for ListView and ListLayout renamed to "list-warp"
* The "warp" property for EditView renamed to "edit-warp"
* Added CertFile and KeyFile optional fields to the AppParams struct.If they are set, then an https connection is created, otherwise http.

# v0.6.0

* Added "user-data" property
* Added "focusable" property
* Added "disabled-items" property to DropDownList
* Added ReloadTableViewData, AllImageResources, NamedColors functions
* Added Theme interface, NewTheme, CreateThemeFromText, and AddTheme functions
* Added image constants to the theme
* Changed BackgroundGradientPoint struct
* Added the background conic gradient

# v0.5.0

* NewApplication function and  Start function of the Application interface were replaced by StartApp function
* Added HasFocus function to the View interface
* Added the UserAgent function to the Session interface
* Added the following properties to TableView: "selection-mode", "allow-selection", "current", "current-style", "current-inactive-style"
* Added the following events to TableView: "table-cell-selected", "table-cell-clicked", "table-row-selected", "table-row-clicked"
* Bug fixing

# v0.4.0

* Added SetTitle and SetTitleColor function to the Session interface
* Added a listener for changing a view property value
* Added the "current" property to StackLayout
* GetDropDownCurrent and GetListViewCurrent functions replaced by the GetCurrent function
* Updated TabsLayout
* Bug fixing

# v0.3.0

* Added FilePicker
* Added DownloadFile and DownloadFileData function to the Session interface
* Updated comments and readme
* Added the FilePicker example to the demo
* Bug fixing

# v0.2.0

* Added "animation" and "transition" properties, Animation interface, animation events
* Renamed ColorProperty constant to ColorTag
* Updated readme
* Added the Animation example to the demo
* Bug fixing

# v0.1.1

* Bug fixing