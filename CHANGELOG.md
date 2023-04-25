# v0.13.0

* Added ViewIndex function to ViewsContainer interface
* Added "tooltip" property and GetTooltip function
* Added "outline-offset" property and GetOutlineOffset function
* Changed the main event listener format for "drop-down-event", "edit-text-changed", 
"color-changed", "number-changed", "date-changed", and "time-changed" events. 
Old format is "<listener>(<view>, <new value>)", new format is "<listener>(<view>, <new value>, <old value>)"
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