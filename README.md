# RUI library

The RUI (Remote User Interface) library is designed to create web applications in the go language.

The peculiarity of the library is that all data processing is carried out on the server,
and the browser is used as a thin client. WebSocket is used for client-server communication.

## Hello world

	type helloWorldSession struct {
	}

	func (content *helloWorldSession) CreateRootView(session rui.Session) rui.View {
		return rui.NewTextView(session, rui.Params {
			rui.Text : "Hello world!!!",
		})
	}

	func createHelloWorldSession(session rui.Session) rui.SessionContent {
		return new(helloWorldSession)
	}

	func main() {
		rui.StartApp("localhost:8000", createHelloWorldSession, rui.AppParams{
			Title:      "Hello world",
			Icon:       "icon.svg",
		})
	}

In the main function, the StartApp function is called. It creates a rui app and runs its main loop.
The StartApp function has 3 parameters:
1) IP address where the application will be available (in our example it is "localhost:8000")
2) The function creates a structure that implements the SessionContent interface
3) Additional optional parameters (in our example, this is the title and the icon file name)

The SessionContent interface is declared as:

	type SessionContent interface {
		CreateRootView(session rui.Session) rui.View
	}

A new instance of the helloWorldSession structure is created for each new session,

The CreateRootView function of the SessionContent interface creates a root element.
When the user accesses the application by typing the address "localhost: 8000" in the browser, a new session is created.
A new instance of the helloWorldSession structure is created for it, and at the end the CreateRootView function is called.
The createRootView function returns a representation of a text that is created using the NewTextView function.

If you want the application to be visible outside your computer, then change the address in the Start function:

	rui.StartApp(rui.GetLocalIP() + ":80", ...

## Used data types

### SizeUnit

The SizeUnit structure is used to set various sizes of interface elements such as width, height, padding, font size, etc.
SizeUnit is declared as

	type SizeUnit struct {
		Type  SizeUnitType
		Value float64
	}

where Type is the type of size; Value is the size.

The Type can take the following values:

| Value    | Constant       | Description                                                                    |
|:--------:|----------------|--------------------------------------------------------------------------------|
| 0        | Auto           | default value. The Value field is ignored                                      |
| 1        | SizeInPixel    | the Value field specifies the size in pixels.                                  |
| 2        | SizeInEM       | the Value field specifies the size in em units. 1em is equal to the base font size, which is set in the browser settings |
| 3        | SizeInEX       | the Value field specifies the size in ex units.                                |
| 4        | SizeInPercent  | the Value field specifies the size as a percentage of the parent element size. |
| 5        | SizeInPt       | the Value field specifies the size in pt units (1pt = 1/72").                  |
| 6        | SizeInPc       | the Value field specifies the size in pc units (1pc = 12pt).                   |
| 7        | SizeInInch     | the Value field specifies the size in inches.                                  |
| 8        | SizeInMM       | the Value field specifies the size in millimeters.                             |
| 9        | SizeInCM       | the Value field defines the size in centimeters.                               |
| 10       | SizeInFraction | the Value field specifies the size in parts. Used only for sizing cells of the GridLayout. |

For a more visual and simple setting of variables of the SizeUnit type, the functions below can be used.

| Function       | Equivalent definition                              |
|----------------|----------------------------------------------------|
| rui.AutoSize() | rui.SizeUnit{ Type: rui.Auto, Value: 0 }           |
| rui.Px(n)      | rui.SizeUnit{ Type: rui.SizeInPixel, Value: n }    |
| rui.Em(n)      | rui.SizeUnit{ Type: rui.SizeInEM, Value: n }       |
| rui.Ex(n)      | rui.SizeUnit{ Type: rui.SizeInEX, Value: n }       |
| rui.Percent(n) | rui.SizeUnit{ Type: rui.SizeInPercent, Value: n }  |
| rui.Pt(n)      | rui.SizeUnit{ Type: rui.SizeInPt, Value: n }       |
| rui.Pc(n)      | rui.SizeUnit{ Type: rui.SizeInPc, Value: n }       |
| rui.Inch(n)    | rui.SizeUnit{ Type: rui.SizeInInch, Value: n }     |
| rui.Mm(n)      | rui.SizeUnit{ Type: rui.SizeInMM, Value: n }       |
| rui.Cm(n)      | rui.SizeUnit{ Type: rui.SizeInCM, Value: n }       |
| rui.Fr(n)      | rui.SizeUnit{ Type: rui.SizeInFraction, Value: n } |

Variables of the SizeUnit type have a textual representation (why you need it will be described below).
The textual representation consists of a number (equal to the value of the Value field) followed by
a suffix defining the type. An exception is a value of type Auto, which has the representation “auto”.
The suffixes are listed in the following table:

| Suffix | Type           |
|:------:|----------------|
| px     | SizeInPixel    |
| em     | SizeInEM       |
| ex     | SizeInEX       |
| %      | SizeInPercent  |
| pt     | SizeInPt       |
| pc     | SizeInPc       |
| in     | SizeInInch     |
| mm     | SizeInMM       |
| cm     | SizeInCM       |
| fr     | SizeInFraction |

Examples: auto, 50%, 32px, 1.5in, 0.8em

To convert the textual representation to the SizeUnit structure, is used the function:

	func StringToSizeUnit(value string) (SizeUnit, bool)

You can get a textual representation of the structure using the String() function of SizeUnit structure

### Color

The Color type describes a 32-bit ARGB color:

	type Color uint32

The Color type has three types of text representations:

1) #AARRGGBB, #RRGGBB, #ARGB, #RGB

where A, R, G, B is a hexadecimal digit describing the corresponding component. If the alpha channel is not specified,
then it is considered equal to FF. If the color component is specified by one digit, then it is doubled.
For example, “# 48AD” is equivalent to “# 4488AADD”

2) argb(A, R, G, B), rgb(R, G, B)

where A, R, G, B is the representation of the color component. The component can be west as a float number in the range [0 … 1],
or as an integer in the range [0 … 255], or as a percentage from 0% to 100%.

Examples:
	
	“argb(255, 128, 96, 0)”
	“rgb(1.0, .5, .8)”
	“rgb(0%, 50%, 25%)”
	“argb(50%, 128, .5, 100%)”

The String function is used to convert a Color to a string.
To convert a string to Color, is used the function:

	func StringToColor(value string) (Color, bool)

3) The name of the color. The RUI library defines the following colors

| Name                  | Color     |
|-----------------------|-----------|
| black                 | #ff000000 |
| silver                | #ffc0c0c0 |
| gray                  | #ff808080 |
| white                 | #ffffffff |
| maroon                | #ff800000 |
| red                   | #ffff0000 |
| purple                | #ff800080 |
| fuchsia               | #ffff00ff |
| green                 | #ff008000 |
| lime                  | #ff00ff00 |
| olive                 | #ff808000 |
| yellow                | #ffffff00 |
| navy                  | #ff000080 |
| blue                  | #ff0000ff |
| teal                  | #ff008080 |
| aqua                  | #ff00ffff |
| orange                | #ffffa500 |
| aliceblue             | #fff0f8ff |
| antiquewhite          | #fffaebd7 |
| aquamarine            | #ff7fffd4 |
| azure                 | #fff0ffff |
| beige                 | #fff5f5dc |
| bisque                | #ffffe4c4 |
| blanchedalmond        | #ffffebcd |
| blueviolet            | #ff8a2be2 |
| brown                 | #ffa52a2a |
| burlywood             | #ffdeb887 |
| cadetblue             | #ff5f9ea0 |
| chartreuse            | #ff7fff00 |
| chocolate             | #ffd2691e |
| coral                 | #ffff7f50 |
| cornflowerblue        | #ff6495ed |
| cornsilk              | #fffff8dc |
| crimson               | #ffdc143c |
| cyan                  | #ff00ffff |
| darkblue              | #ff00008b |
| darkcyan              | #ff008b8b |
| darkgoldenrod         | #ffb8860b |
| darkgray              | #ffa9a9a9 |
| darkgreen             | #ff006400 |
| darkgrey              | #ffa9a9a9 |
| darkkhaki             | #ffbdb76b |
| darkmagenta           | #ff8b008b |
| darkolivegreen        | #ff556b2f |
| darkorange            | #ffff8c00 |
| darkorchid            | #ff9932cc |
| darkred               | #ff8b0000 |
| darksalmon            | #ffe9967a |
| darkseagreen          | #ff8fbc8f |
| darkslateblue         | #ff483d8b |
| darkslategray         | #ff2f4f4f |
| darkslategrey         | #ff2f4f4f |
| darkturquoise         | #ff00ced1 |
| darkviolet            | #ff9400d3 |
| deeppink              | #ffff1493 |
| deepskyblue           | #ff00bfff |
| dimgray               | #ff696969 |
| dimgrey               | #ff696969 |
| dodgerblue            | #ff1e90ff |
| firebrick             | #ffb22222 |
| floralwhite           | #fffffaf0 |
| forestgreen           | #ff228b22 |
| gainsboro             | #ffdcdcdc |
| ghostwhite            | #fff8f8ff |
| gold                  | #ffffd700 |
| goldenrod             | #ffdaa520 |
| greenyellow           | #ffadff2f |
| grey                  | #ff808080 |
| honeydew              | #fff0fff0 |
| hotpink               | #ffff69b4 |
| indianred             | #ffcd5c5c |
| indigo                | #ff4b0082 |
| ivory                 | #fffffff0 |
| khaki                 | #fff0e68c |
| lavender              | #ffe6e6fa |
| lavenderblush         | #fffff0f5 |
| lawngreen             | #ff7cfc00 |
| lemonchiffon          | #fffffacd |
| lightblue             | #ffadd8e6 |
| lightcoral            | #fff08080 |
| lightcyan             | #ffe0ffff |
| lightgoldenrodyellow  | #fffafad2 |
| lightgray             | #ffd3d3d3 |
| lightgreen            | #ff90ee90 |
| lightgrey             | #ffd3d3d3 |
| lightpink             | #ffffb6c1 |
| lightsalmon           | #ffffa07a |
| lightseagreen         | #ff20b2aa |
| lightskyblue          | #ff87cefa |
| lightslategray        | #ff778899 |
| lightslategrey        | #ff778899 |
| lightsteelblue        | #ffb0c4de |
| lightyellow           | #ffffffe0 |
| limegreen             | #ff32cd32 |
| linen                 | #fffaf0e6 |
| magenta               | #ffff00ff |
| mediumaquamarine      | #ff66cdaa |
| mediumblue            | #ff0000cd |
| mediumorchid          | #ffba55d3 |
| mediumpurple          | #ff9370db |
| mediumseagreen        | #ff3cb371 |
| mediumslateblue       | #ff7b68ee |
| mediumspringgreen     | #ff00fa9a |
| mediumturquoise       | #ff48d1cc |
| mediumvioletred       | #ffc71585 |
| midnightblue          | #ff191970 |
| mintcream             | #fff5fffa |
| mistyrose             | #ffffe4e1 |
| moccasin              | #ffffe4b5 |
| navajowhite           | #ffffdead |
| oldlace               | #fffdf5e6 |
| olivedrab             | #ff6b8e23 |
| orangered             | #ffff4500 |
| orchid                | #ffda70d6 |
| palegoldenrod         | #ffeee8aa |
| palegreen             | #ff98fb98 |
| paleturquoise         | #ffafeeee |
| palevioletred         | #ffdb7093 |
| papayawhip            | #ffffefd5 |
| peachpuff             | #ffffdab9 |
| peru                  | #ffcd853f |
| pink                  | #ffffc0cb |
| plum                  | #ffdda0dd |
| powderblue            | #ffb0e0e6 |
| rosybrown             | #ffbc8f8f |
| royalblue             | #ff4169e1 |
| saddlebrown           | #ff8b4513 |
| salmon                | #fffa8072 |
| sandybrown            | #fff4a460 |
| seagreen              | #ff2e8b57 |
| seashell              | #fffff5ee |
| sienna                | #ffa0522d |
| skyblue               | #ff87ceeb |
| slateblue             | #ff6a5acd |
| slategray             | #ff708090 |
| slategrey             | #ff708090 |
| snow                  | #fffffafa |
| springgreen           | #ff00ff7f |
| steelblue             | #ff4682b4 |
| tan                   | #ffd2b48c |
| thistle               | #ffd8bfd8 |
| tomato                | #ffff6347 |
| turquoise             | #ff40e0d0 |
| violet                | #ffee82ee |
| wheat                 | #fff5deb3 |
| whitesmoke            | #fff5f5f5 |
| yellowgreen           | #ff9acd32 |

### AngleUnit

The AngleUnit type is used to set angular values. AngleUnit is declared as

	type AngleUnit struct {
		Type  AngleUnitType
		Value float64
	}

where Type is the type of angular value; Value is the angular value

The Type can take the following values:

* Radian (0) - the Value field defines the angular value in radians.
* PiRadian (1) - the Value field defines the angular value in radians multiplied by π.
* Degree (2) - the Value field defines the angular value in degrees.
* Gradian (3) - the Value field defines the angular value in grades (gradians).
* Turn (4) - the Value field defines the angular value in turns (1 turn == 360°).

For a more visual and simple setting of variables of the AngleUnit type, the functions below can be used.

| Function     | Equivalent definition                         |
|--------------|-----------------------------------------------|
| rui.Rad(n)   | rui.AngleUnit{ Type: rui.Radian, Value: n }   |
| rui.PiRad(n) | rui.AngleUnit{ Type: rui.PiRadian, Value: n } |
| rui.Deg(n)   | rui.AngleUnit{ Type: rui.Degree, Value: n }   |
| rui.Grad(n)  | rui.AngleUnit{ Type: rui.Gradian, Value: n }  |

Variables of type AngleUnit have a textual representation consisting of a number (equal to the value of the Value field)
followed by a suffix defining the type. The suffixes are listed in the following table:

| Suffix  | Type     |
|:-------:|----------|
| deg     | Degree   |
| °       | Degree   |
| rad     | Radian   |
| π       | PiRadian |
| pi      | PiRadian |
| grad    | Gradian  |
| turn    | Turn     |

Examples: “45deg”, “90°”, “3.14rad”, “2π”, “0.5pi”

The String function is used to convert AngleUnit to a string.
To convert a string to AngleUnit is used the function:

	func StringToAngleUnit(value string) (AngleUnit, bool)

## View

View is an interface for accessing an element of "View" type. View is a rectangular area of the screen.
All interface elements extend the View interface, i.e. View is the base element for all other elements in the library.

View has a number of properties like height, width, color, text parameters, etc. Each property has a text name.
The Properties interface is used to read and write the property value (View implements this interface):

	type Properties interface {
		Get(tag string) interface{}
		Set(tag string, value interface{}) bool
		Remove(tag string)
		Clear()
		AllTags() []string
	}

The Get function returns the value of the property, or nil if the property is not set.

The Set function sets the value of a property. If the property value is set successfully, then
the function returns true, if not, then false and a description of the error that occurred is written to the log.

The Remove function removes property value, equivalent to Set(nil)

To simplify setting / reading properties, there are also two global functions Get and Set:

	func Get(rootView View, viewID, tag string) interface{}
	func Set(rootView View, viewID, tag string, value interface{}) bool

These functions get/set the value of the child View

### Tracking property changes

You can set a function to track the change of absolutely any View property (there are no exceptions).
To set up a change listener, the View interface contains a function:

	SetChangeListener(tag string, listener func(View, string))

where the first parameter is the name of the tracked property, and the second is the function 
that will be called every time the property value changes.

For example

	view.SetChangeListener(rui.BackgroundColor, listener func(view View, tag string) {
		// The background color changed
	})

### Events

When interacting with the application, various events arise: clicks, resizing, changing input data, etc.

Event listeners are designed to respond to events. A listener is a function that is called every time an event occurs.
Each event can have multiple listeners. Let's analyze the listeners using the example of the "edit-text-changed"
text change event in the "EditView" editor.

The event listener is a function of the form

	func(<View>[, <parameters>])

where the first argument is the View in which the event occurred. Further there are additional parameters of the event.

For "edit-text-changed", the main listener will look like this:

	func(EditView, string)

where the second argument is the new text value

If you do not plan to use the first argument, you can omit it. This will be an additional listener

	func(string)

In order to assign a listener, you must assign it to a property with the event name

	view.Set(rui.EditTextChanged, func(edit EditView, newText string) {
		// do something
	})

or

	view.Set(rui.EditTextChanged, func(newText string) {
		// do something
	})

Each event can have multiple listeners. In this regard, five data types can be used as listeners:

* func(< View >[, < parameters >])
* func([< parameters>])
* []func(< View >[, < parameters >])
* []func([< parameters >])
* []interface{} which only contains func(< View >[, < parameters >]) and func([< parameters >])

After being assigned to a property, all these types are converted to an array of []func(< View >, [< parameters >]).
Accordingly, the Get function always returns an array of []func(< View >, [< parameters >]).
If there are no listeners, this array will be empty.

For the "edit-text-changed" event, this

* func(editor EditView, newText string)
* func(newText string)
* []func(editor EditView, newText string)
* []func(newText string)
* []interface{} содержащий только func(editor EditView, newText string) и func(newText string)

And the "edit-text-changed" property always stores and returns []func(EditView, string).

In what follows, when describing specific events, only the format of the main listener will be presented.

### "id" property

The "id" property is an optional textual identifier for the View. With it, you can find the child View.
To do this, use the ViewByID function

	func ViewByID(rootView View, id string) View

This function looks for a child View with id. The search starts from rootView.
If View is not found, the function returns nil and an error message is written to the log.

Usually id is set when the View is created and is not changed later.
But this is an optional condition. You can change the id at any time.

The Set function is used to set a new value for id. For example

	view.Set(rui.ID, "myView")
	view.Set("id", "myView")

There are two ways to get the id. The first is using the Get function:

	if value := view.Get(rui.ID); value != nil {
		id = value.(string)
	}
	
And the second one is using the ID() function:

	id = view.ID()

### "width", "height", "min-width", "min-height", "max-width", "max-height" properties

These properties are set:

| Property     | Constant      | Description                |
|--------------|---------------|----------------------------|
| "width"      | rui.Width     | The width of View          |
| "height"     | rui.Height    | The height of View         |
| "min-width"  | rui.MinWidth  | The minimum width of View  |
| "min-height" | rui.MinHeight | The minimum height of View |
| "max-width"  | rui.MaxWidth  | The maximum width of View  |
| "max-height" | rui.MaxHeight | The maximum height of View |

These properties are of type SizeUnit.
If the "width" / "height" value is not set or is set to Auto, then the height/width of the View
is determined by its content and limited to the minimum and maximum height/width.
As the value of these properties, you can set the SizeUnit structure, the textual representation of the SizeUnit,
or the name of the constant (about the constants below):

	view.Set("width", rui.Px(8))
	view.Set(rui.MaxHeight, "80%")
	view.Set(rui.Height, "@viewHeight")

After getting the value with the Get function, you must typecast:

	if value := view.Get(rui.Width); value != nil {
		switch value.(type) {
			case string:
				text := value.(string)
				// TODO

			case SizeUnit:	
				size := value.(SizeUnit)
				// TODO
		}
	}

This is quite cumbersome, therefore for each property there is a global function of the same name with the Get prefix,
which performs the given cast, gets the value of the constant, if necessary, and returns it.
All functions of this type have two arguments: View and subviewID string.
The first argument is the root View, the second is the ID of the child View.
If the ID of the child View is passed as "", then the value of the root View is returned.
For the properties "width", "height", "min-width", "min-height", "max-width", "max-height" these are functions:

	func GetWidth(view View, subviewID string) SizeUnit
	func GetHeight(view View, subviewID string) SizeUnit
	func GetMinWidth(view View, subviewID string) SizeUnit
	func GetMinHeight(view View, subviewID string) SizeUnit
	func GetMaxWidth(view View, subviewID string) SizeUnit
	func GetMaxHeight(view View, subviewID string) SizeUnit

### "margin" and "padding" properties

The "margin" property determines the outer margins from this View to its neighbors.
The "padding" property sets the padding from the border of the View to the content.
The values of the "margin" and "padding" properties are stored as the BoundsProperty interface,
which implements the Properties interface (see above). BoundsProperty has 4 SizeUnit properties:

| Property  | Constant    | Description       |
|-----------|-------------|-------------------|
| "top"     | rui.Top     | Top padding       |
| "right"   | rui.Right   | Right padding     |
| "bottom"  | rui.Bottom  | Bottom padding    |
| "left"    | rui.Left    | Дуае padding      |

The NewBoundsProperty function is used to create the BoundsProperty interface. Example

	view.Set(rui.Margin, NewBoundsProperty(rui.Params {
		rui.Top:  rui.Px(8),
		rui.Left: "@topMargin",
		"right":   "1.5em",
		"bottom":  rui.Inch(0.3),
	})))

Accordingly, if you request the "margin" or "padding" property using the Get method, the BoundsProperty interface will return:

	if value := view.Get(rui.Margin); value != nil {
		margin := value.(BoundsProperty)
	}

BoundsProperty using the "Bounds (session Session) Bounds" function of the BoundsProperty interface
can be converted to a more convenient Bounds structure:

	type Bounds struct {
		Top, Right, Bottom, Left SizeUnit
	}

Global functions can also be used for this:

	func GetMargin(view View, subviewID string) Bounds
	func GetPadding(view View, subviewID string) Bounds

The textual representation of the BoundsProperty is as follows:

	"_{ top = <top padding>, right = <right padding>, bottom = <bottom padding>, left = <left padding> }"

The value of the "margin" and "padding" properties can be passed to the Set method:
* BoundsProperty interface or its textual representation;
* Bounds structure;
* SizeUnit or the name of a constant of type SizeUnit, in which case this value is set to all indents. Those.

	view.Set(rui.Margin, rui.Px(8))

equivalent to

	view.Set(rui.Margin, rui.Bounds{Top: rui.Px(8), Right: rui.Px(8), Bottom: rui.Px(8), Left: rui.Px(8)})

Since the value of the "margin" and "padding" property is always stored as the BoundsProperty interface,
if you read the "margin" or "padding" property set by the Bounds or SizeUnit with the Get function,
then you get the BoundsProperty, not the Bounds or SizeUnit.

The "margin" and "padding" properties are used to set four margins at once.
The following properties are used to set individual paddings:

| Property         | Constant          | Description        |
|------------------|-------------------|--------------------|
| "margin-top"     | rui.MarginTop     | The top margin     |
| "margin-right"   | rui.MarginRight   | The right margin   |
| "margin-bottom"  | rui.MarginBottom  | The bottom margin  |
| "margin-left"    | rui.MarginLeft    | The left margin    |
| "padding-top"    | rui.PaddingTop    | The top padding    |
| "padding-right"  | rui.PaddingRight  | The right padding  |
| "padding-bottom" | rui.PaddingBottom | The bottom padding |
| "padding-left"   | rui.PaddingLeft   | The left padding   |

Example

	view.Set(rui.Margin, rui.Px(8))
	view.Set(rui.TopMargin, rui.Px(12))

equivalent to

	view.Set(rui.Margin, rui.Bounds{Top: rui.Px(12), Right: rui.Px(8), Bottom: rui.Px(8), Left: rui.Px(8)})

### "border" property

The "border" property defines a border around the View. The frame line is described by three attributes:
line style, thickness and color.

The value of the "border" property is stored as the BorderProperty interface,
which implements the Properties interface (see above). BorderProperty can contain the following properties:

| Property       | Constant    | Type     | Description              |
|----------------|-------------|----------|--------------------------|
| "left-style"   | LeftStyle   | int      | Left border line style   |
| "right-style"  | RightStyle  | int      | Right border line style  |
| "top-style"    | TopStyle    | int      | Top  border line style   |
| "bottom-style" | BottomStyle | int      | Bottom border line style |
| "left-width"   | LeftWidth   | SizeUnit | Left border line width   |
| "right-width"  | RightWidth  | SizeUnit | Right border line width  |
| "top-width"    | TopWidth    | SizeUnit | Top border line width    |
| "bottom-width" | BottomWidth | SizeUnit | Bottom border line width |
| "left-color"   | LeftColor   | Color    | Left border line color   |
| "right-color"  | RightColor  | Color    | Right border line color  |
| "top-color"    | TopColor    | Color    | Top border line color    |
| "bottom-color" | BottomColor | Color    | Bottom border line color |

Line style can take the following values:

| Value | Constant   | Name     | Description         |
|:-----:|------------|----------|---------------------|
| 0     | NoneLine   | "none"   | No frame            |
| 1     | SolidLine  | "solid"  | Solid line          |
| 2     | DashedLine | "dashed" | Dashed line         |
| 3     | DottedLine | "dotted" | Dotted line         |
| 4     | DoubleLine | "double" | Double solid line   |

All other style values are ignored.

The NewBorder function is used to create the BorderProperty interface.

If all the lines of the frame are the same, then the following properties can be used to set the style, thickness and color:

| Property  | Constant | Type     | Description         |
|-----------|----------|----------|---------------------|
| "style"   | Style    | int      | Border line style   |
| "width"   | Width    | SizeUnit | Border line width   |
| "color"   | Color    | Color    | Border line color   |

Example

	view.Set(rui.Border, NewBorder(rui.Params{
		rui.LeftStyle:   rui.SolidBorder,
		rui.RightStyle:  rui.SolidBorder,
		rui.TopStyle:    rui.SolidBorder,
		rui.BottomStyle: rui.SolidBorder,
		rui.LeftWidth:   rui.Px(1),
		rui.RightWidth:  rui.Px(1),
		rui.TopWidth:    rui.Px(1),
		rui.BottomWidth: rui.Px(1),
		rui.LeftColor:   rui.Black,
		rui.RightColor:  rui.Black,
		rui.TopColor:    rui.Black,
		rui.BottomColor: rui.Black,
	}))

equivalent to

	view.Set(rui.Border, NewBorder(rui.Params{
		rui.Style: rui.SolidBorder,
		rui.Width: rui.Px(1),
		rui.ColorProperty: rui.Black,
	}))

The BorderProperty interface can be converted to a ViewBorders structure using the Border function.
When converted, all text constants are replaced with real values. ViewBorders is described as

	 type ViewBorders struct {
		Top, Right, Bottom, Left ViewBorder
	}

where the ViewBorder structure is described as

	type ViewBorder struct {
		Style int
		Color Color
		Width SizeUnit
	}

The ViewBorders structure can be passed as a parameter to the Set function when setting the value of the "border" property.
This converts the ViewBorders to BorderProperty. Therefore, when the property is read,
the Get function will return the BorderProperty interface, not the ViewBorders structure.
You can get the ViewBorders structure without additional transformations using the global function

	func GetBorder(view View, subviewID string) ViewBorders

Besides the auxiliary properties "style", "width" and "color" there are 4 more: "left", "right", "top" and "bottom".
As a value, these properties can only take the ViewBorder structure and allow you to set all the attributes of the line of the side of the same name.

You can also set individual frame attributes using the Set function of the View interface.
For this, the following properties are used

| Property              | Constant          | Type       | Description              |
|-----------------------|-------------------|------------|--------------------------|
| "border-left-style"   | BorderLeftStyle   | int        | Left border line style   |
| "border-right-style"  | BorderRightStyle  | int        | Right border line style  |
| "border-top-style"    | BorderTopStyle    | int        | Top  border line style   |
| "border-bottom-style" | BorderBottomStyle | int        | Bottom border line style |
| "border-left-width"   | BorderLeftWidth   | SizeUnit   | Left border line width   |
| "border-right-width"  | BorderRightWidth  | SizeUnit   | Right border line width  |
| "border-top-width"    | BorderTopWidth    | SizeUnit   | Top border line width    |
| "border-bottom-width" | BorderBottomWidth | SizeUnit   | Bottom border line width |
| "border-left-color"   | BorderLeftColor   | Color      | Left border line color   |
| "border-right-color"  | BorderRightColor  | Color      | Right border line color  |
| "border-top-color"    | BorderTopColor    | Color      | Top border line color    |
| "border-bottom-color" | BorderBottomColor | Color      | Bottom border line color |
| "border-style"        | BorderStyle       | int        | Border line style        |
| "border-width"        | BorderWidth       | SizeUnit   | Border line width        |
| "border-color"        | BorderColor       | Color      | Border line color        |
| "border-left"         | BorderLeft        | ViewBorder | Left border line         |
| "border-right"        | BorderRight       | ViewBorder | Right border line        |
| "border-top"          | BorderTop         | ViewBorder | Top  border line         |
| "border-bottom"       | BorderBottom      | ViewBorder | Bottom border line       |

Example

	view.Set(rui.BorderStyle, rui.SolidBorder)
	view.Set(rui.BorderWidth, rui.Px(1))
	view.Set(rui.BorderColor, rui.Black)

equivalent to

	view.Set(rui.Border, NewBorder(rui.Params{
		rui.Style: rui.SolidBorder,
		rui.Width: rui.Px(1),
		rui.ColorProperty: rui.Black,
	}))

### "radius" property

The "radius" property sets the elliptical corner radius of the View. Radii are specified by the RadiusProperty
interface that implements the Properties interface (see above).
For this, the following properties of the SizeUnit type are used:

| Property         | Constant     | Description                         |
|------------------|--------------|-------------------------------------|
| "top-left-x"     | TopLeftX     | x-radius of the top left corner     |
| "top-left-y"     | TopLeftY     | y-radius of the top left corner     |
| "top-right-x"    | TopRightX    | x-radius of the top right corner    |
| "top-right-y"    | TopRightY    | y-radius of the top right corner    |
| "bottom-left-x"  | BottomLeftX  | x-radius of the bottom left corner  |
| "bottom-left-y"  | BottomLeftY  | y-radius of the bottom left corner  |
| "bottom-right-x" | BottomRightX | x-radius of the bottom right corner |
| "bottom-right-y" | BottomRightY | y-radius of the bottom right corner |

If the x- and y-radii are the same, then you can use the auxiliary properties

| Property       | Constant    | Description                |
|----------------|-------------|----------------------------|
| "top-left"     | TopLeft     | top left corner radius     |
| "top-right"    | TopRight    | top right corner radius    |
| "bottom-left"  | BottomLeft  | bottom left corner radius  |
| "bottom-right" | BottomRight | bottom right corner radius |

To set all radii to the same values, use the "x" and "y" properties

The RadiusProperty interface is created using the NewRadiusProperty function. Example

	view.Set(rui.Radius, NewRadiusProperty(rui.Params{
		rui.X: rui.Px(16),
		rui.Y: rui.Px(8),
		rui.TopLeft: rui.Px(0),
		rui.BottomRight: rui.Px(0),
	}))

equivalent to

	view.Set(rui.Radius, NewRadiusProperty(rui.Params{
		rui.TopRightX: rui.Px(16),
		rui.TopRightY: rui.Px(8),
		rui.BottomLeftX: rui.Px(16),
		rui.BottomLeftY: rui.Px(8),
		rui.TopLeftX: rui.Px(0),
		rui.TopLeftX: rui.Px(0),
		rui.BottomRightX: rui.Px(0),
		rui.BottomRightY: rui.Px(0),
	}))

If all radii are the same, then the given SizeUnit value can be directly assigned to the "radius" property

	view.Set(rui.Radius, rui.Px(4))

RadiusProperty has a textual representation of the following form:

	_{ <radius id> = <SizeUnit text> [/ <SizeUnit text>] [, <radius id> = <SizeUnit text> [/ <SizeUnit text>]] … }

where <radius id> can take the following values: "x", "y", "top-left", "top-left-x", "top-left-y", "top-right",
"top-right-x", "top-right-y", "bottom-left", "bottom-left-x", "bottom-left-y", "bottom-right", "bottom-right-x", "bottom-right-y".

Values like "<SizeUnit text> / <SizeUnit text>" can only be assigned to the "top-left", "top-right", "bottom-left" and "bottom-right" properties.

Examples:

	_{ x = 4px, y = 4px, top-left = 8px, bottom-right = 8px }

equivalent to

	_{ top-left = 8px, top-right = 4px, bottom-left = 4px, bottom-right = 8px }

or

	_{ top-left = 8px / 8px, top-right = 4px / 4px, bottom-left = 4px / 4px, bottom-right = 8px / 8px }
	
or

	_{ top-left-x = 8px, top-left-y = 8px, top-right-x = 4px, top-right-y = 4px,
		bottom-left-x = 4px, bottom-left-y = 4px, bottom-right-x = 8px, bottom-right-y = 8px }

The RadiusProperty interface can be converted to a BoxRadius structure using the BoxRadius function.
When converted, all text constants are replaced with real values. BoxRadius is described as

	type BoxRadius struct {
		TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY SizeUnit
	}

The BoxRadius structure can be passed as a parameter to the Set function by setting the value of the "radius" property.
This converts BoxRadius to RadiusProperty. Therefore, when the property is read, 
the Get function will return the RadiusProperty interface, not the BoxRadius structure. 
You can get the BoxRadius structure without additional transformations using the global function

	func GetRadius(view View, subviewID string) BoxRadius

You can also set individual radii using the Set function of the View interface.
For this, the following properties are used

| Property                | Constant           | Description                         |
|-------------------------|--------------------|-------------------------------------|
| "radius-x"              | RadiusX            | All x-radii                       |
| "radius-y"              | RadiusY            | All y-radii                       |
| "radius-top-left-x"     | RadiusTopLeftX     | x-radius of the top left corner     |
| "radius-top-left-y"     | RadiusTopLeftY     | y-radius of the top left corner     |
| "radius-top-right-x"    | RadiusTopRightX    | x-radius of the top right corner    |
| "radius-top-right-y"    | RadiusTopRightY    | y-radius of the top right corner    |
| "radius-bottom-left-x"  | RadiusBottomLeftX  | x-radius of the bottom left corner  |
| "radius-bottom-left-y"  | RadiusBottomLeftY  | y-radius of the bottom left corner  |
| "radius-bottom-right-x" | RadiusBottomRightX | x-radius of the bottom right corner |
| "radius-bottom-right-y" | RadiusBottomRightY | y-radius of the bottom right corner |
| "radius-top-left"       | RadiusTopLeft      | top left corner radius              |
| "radius-top-right"      | RadiusTopRight     | top right corner radius             |
| "radius-bottom-left"    | RadiusBottomLeft   | bottom left corner radius           |
| "radius-bottom-right"   | RadiusBottomRight  | bottom right corner radius          |

Example

	view.Set(rui.RadiusX, rui.Px(4))
	view.Set(rui.RadiusY, rui.Px(32))

equivalent to

	view.Set(rui.Border, NewRadiusProperty(rui.Params{
		rui.X: rui.Px(4),
		rui.Y: rui.Px(32),
	}))

### "shadow" property

The "shadow" property allows you to set shadows for the View. There may be several shadows. 
The shadow is described using the ViewShadow interface extending the Properties interface (see above). 
The shadow has the following properties:

| Property        | Constant      | Type     | Description                                                           |
|-----------------|---------------|----------|-----------------------------------------------------------------------|
| "color"         | ColorProperty | Color    | Shadow color                                                          |
| "inset"         | Inset         | bool     | true - the shadow inside the View, false - outside                    |
| "x-offset"      | XOffset       | SizeUnit | Offset the shadow along the X axis                                    |
| "y-offset"      | YOffset       | SizeUnit | Offset the shadow along the Y axis                                    |
| "blur"          | BlurRadius    | float    | Shadow blur radius. The value must be >= 0                            |
| "spread-radius" | SpreadRadius  | float    | Increase the shadow. Value > 0 increases shadow, < 0 decreases shadow |

Three functions are used to create a ViewShadow:

	func NewViewShadow(offsetX, offsetY, blurRadius, spread-radius SizeUnit, color Color) ViewShadow
	func NewInsetViewShadow(offsetX, offsetY, blurRadius, spread-radius SizeUnit, color Color) ViewShadow
	func NewShadowWithParams(params Params) ViewShadow

The NewViewShadow function creates an outer shadow (Inset == false), 
NewInsetViewShadow - an inner one (Inset == true).
The NewShadowWithParams function is used when constants must be used as parameters. 
For example:

	shadow := NewShadowWithParams(rui.Params{
		rui.ColorProperty : "@shadowColor",
		rui.BlurRadius : 8.0,
		rui.Dilation : 16.0,
	})

ViewShadow, ViewShadow array, and ViewShadow textual representation can be assigned as a value to the "shadow" property.

The ViewShadow text representation has the following format:

	_{ color = <color> [, x-offset = <offset>] [, y-offset = <offset>] [, blur = <radius>]
		[, spread-radius = <increase>] [, inset = <type>] }

You can get the value of "shadow" property using the function

	func GetViewShadows(view View, subviewID string) []ViewShadow

If no shadow is specified, then this function will return an empty array

### "background-color" property

Constant: rui.BackgroundColor. Get function: GetBackgroundColor() Color

The "background-color" property sets the background color. Valid values are Color, an integer, the textual representation of Color, 
and a constant name starting with '@'. An integer must encode the color in the AARRGGBB format

In addition to color, images and gradients can also be used as backgrounds (see below).
In this case, "background-color" is used for transparent areas of images.

### "background-clip" property

The "background-clip" property determines how the background color and / or background image will be displayed below the box borders.

If no background image or color is specified, this property will have a visual effect only 
if the border has transparent areas or partially opaque areas; otherwise, the border hides the difference.

The property can take the following values:

| Value | Constant       | Name          | Description                                    |
|:-----:|----------------|---------------|------------------------------------------------|
| 0     | BorderBoxClip  | "border-box"  | The background extends to the outer edge of the border (but below the border in z-order).         |
| 1     | PaddingBoxClip | "padding-box" | The background extends to the outer edge of the padding. No background is drawn below the border. |
| 2     | ContentBoxClip | "content-box" | The background is painted inside (clipped) of the content box. |

### "background" property

In addition to color, pictures and / or gradient fills can also be specified as the background of the View.
The property "background" is used for this. The background can contain multiple images and gradients.
Each background element is described by the BackgroundElement interface. 
BackgroundElement can be of three types: linear gradient, radial gradient, and image.

#### Linear gradient

A linear gradient is created using the function

	func NewBackgroundLinearGradient(params Params) BackgroundElement

The linear gradient has the following options:

* Direction ("direction") - defines the direction of the gradient line (the line along which the color changes).
Optional parameter. The default direction is from bottom to top. It can be either AngleUnit 
(the angle of inclination of the line relative to the vertical) or one of the following int values:

| Value | Constant              | Name              | Description                                   |
|:-----:|-----------------------|-------------------|-----------------------------------------------|
| 0     | ToTopGradient         | "to-top"          | Line goes from bottom to top (default)        |
| 1     | ToRightTopGradient    | "to-right-top"    | From bottom left to top right                 |
| 2     | ToRightGradient       | "to-right"        | From left to right                            |
| 3     | ToRightBottomGradient | "to-right-bottom" | From top left to bottom right                 |
| 4     | ToBottomGradient      | "to-bottom"       | From top to bottom                            |
| 5     | ToLeftBottomGradient  | "to-left-bottom"  | From the upper right corner to the lower left |
| 6     | ToLeftGradient        | "to-left"         | From right to left                            |
| 7     | ToLeftTopGradient     | "to-left-top"     | From the bottom right corner to the top left  |

* Gradient ("gradient") - array of gradient key points (required parameter). 
Each point is described by a BackgroundGradientPoint structure, which has two fields: Pos of type SizeUnit and Color. 
Pos defines the position of the point relative to the start of the gradient line. The array must have at least 2 points.
You can also pass a Color array as the gradient value. In this case, the points are evenly distributed along the gradient line.
You can also use an array of []interface{} as an array of cue points.
The elements of this array can be BackgroundGradientPoint, Color, BackgroundGradientPoint or Color text representation, and the name of the constant

* Repeat ("repeat") - a boolean value that determines whether the gradient will repeat after the last key point. 
Optional parameter. The default is false (do not repeat)

The linear gradient text representation is as follows:

	linear-gradient { gradient = <value> [, direction = <value>] [, repeat = <value>] }

#### Radial gradient

A radial gradient is created using the function

	func NewBackgroundRadialGradient(params Params) BackgroundElement

The radial gradient has the following parameters:

* Gradient ("gradient") - array of gradient key points (required parameter). Identical to the linear gradient parameter of the same name.

* Repeat ("repeat") - a boolean value that determines whether the gradient will repeat after the last key point. 
Optional parameter. The default is false (do not repeat)

* RadialGradientShape ("radial-gradient-shape") or Shape ("shape") - defines the shape of the gradient.
It can take one of two int values:

| Value | Constant        | Name      | Description                                  |
|:-----:|-----------------|-----------|----------------------------------------------|
| 0     | EllipseGradient | "ellipse" | The shape is an axis-aligned ellipse         |
| 1     | CircleGradient  | "circle"  | The shape is a circle with a constant radius |

Optional parameter. The default is EllipseGradient

* RadialGradientRadius ("radial-gradient-radius") or Radius ("radius") - sets the radius of the gradient.
Can be either SizeUnit or one of the following int values:

| Value | Constant               | Name              | Description                                |
|:-----:|------------------------|-------------------|--------------------------------------------|
| 0     | ClosestSideGradient    | "closest-side"    | The final shape of the gradient corresponds to the side of the rectangle closest to its center (for circles), or both vertical and horizontal sides closest to the center (for ellipses) |
| 1     | ClosestCornerGradient  | "closest-corner"  | The final shape of the gradient is defined so that it exactly matches the closest corner of the window from its center |
| 2     | FarthestSideGradient   | "farthest-side"   | Similar to ClosestSideGradient, except that the size of the shape is determined by the farthest side from its center (or vertical and horizontal sides) |
| 3     | FarthestCornerGradient | "farthest-corner" | The final shape of the gradient is defined so that it exactly matches the farthest corner of the rectangle from its center |

Optional parameter. The default is ClosestSideGradient

* CenterX ("center-x"), CenterY ("center-y") - sets the center of the gradient relative to the upper left corner of the View. Takes in a SizeUnit value. Optional parameter.
The default value is "50%", i.e. the center of the gradient is the center of the View.

The linear gradient text representation is as follows:

	radial-gradient { gradient = <Value> [, repeat = <Value>] [, shape = <Value>]
		[, radius = <Value>][, center-x = <Value>][, center-y = <Value>]}

#### Image

The image has the following parameters:

* Source ("src") - Specifies the URL of the image

* Fit ("fit") - an optional parameter that determines the scaling of the image.
Can be one of the following Int values:

| Constant   | Value | Name      | Description                                                                                             |
|------------|:-----:|-----------|---------------------------------------------------------------------------------------------------------|
| NoneFit    | 0     | "none"    | No scaling (default). The dimensions of the image are determined by the Width and Height parameters.    |
| ContainFit | 1     | "contain" | The image is scaled proportionally so that its width or height is equal to the width or height of the background area. Image can be cropped to width or height |
| CoverFit   | 2     | "cover"   | The image is scaled with the same proportions so that the whole picture fits inside the background area |

* Width ("width"), Height (height) - optional SizeUnit parameters that specify the height and width of the image. 
Used only if Fit is NoneFit. The default is Auto (original size). The percentage value sets the size relative 
to the height and width of the background area, respectively

* Attachment - 

* Repeat (repeat) - an optional parameter specifying the repetition of the image.
Can be one of the following int values:

| Constant    | Value | Name        | Description                                     |
|-------------|:-----:|-------------|-------------------------------------------------|
| NoRepeat    | 0     | "no-repeat" | Image does not repeat (default)                 |
| RepeatXY    | 1     | "repeat"    | The image repeats horizontally and vertically   |
| RepeatX     | 2     | "repeat-x"  | The image repeats only horizontally             |
| RepeatY     | 3     | "repeat-y"  | Image repeats vertically only                   |
| RepeatRound | 4     | "round"     | The image is repeated so that an integer number of images fit into the background area; if this fails, then the background images are scaled |
| RepeatSpace | 5     | "space"     | The image is repeated as many times as necessary to fill the background area; if this fails, an empty space is added between the pictures    |

* ImageHorizontalAlign,

* ImageVerticalAlign,

### "clip" property

The "clip" property (Clip constant) of the ClipShape type specifies the crop area.
There are 4 types of crop areas

#### inset

Rectangular cropping area. Created with the function:

	func InsetClip(top, right, bottom, left SizeUnit, radius RadiusProperty) ClipShape

where top, right, bottom, left are the distance from respectively the top, right, bottom and left borders of the View 
to the cropping border of the same name; radius - sets the radii of the corners of the cropping area 
(see the description of the RadiusProperty type above). If there should be no rounding of corners, then nil must be passed as radius

The textual description of the rectangular cropping area is in the following format

	inset{ top = <top value>, right = <right value>, bottom = <bottom value>, left = <left value>,
		[radius = <RadiusProperty text>] }
	}

#### circle

Round cropping area. Created with the function:

	func CircleClip(x, y, radius SizeUnit) ClipShape

where x, y - coordinates of the center of the circle; radius - radius

The textual description of the circular cropping area is in the following format

	circle{ x = <x value>, y = <y value>, radius = <radius value> }

#### ellipse

Elliptical cropping area. Created with the function:

	func EllipseClip(x, y, rx, ry SizeUnit) ClipShape

where x, y - coordinates of the center of the ellipse; rх - radius of the ellipse along the X axis; ry is the radius of the ellipse along the Y axis.

The textual description of the elliptical clipping region is in the following format

	ellipse{ x = <x value>, y = <y value>, radius-x = <x radius value>, radius-y = <y radius value> }

#### polygon

Polygonal cropping area. Created using functions:

	func PolygonClip(points []interface{}) ClipShape
	func PolygonPointsClip(points []SizeUnit) ClipShape

an array of corner points of the polygon is passed as an argument in the following order: x1, y1, x2, y2, …
The elements of the argument to the PolygonClip function can be either text constants, 
or the text representation of SizeUnit, or elements of type SizeUnit.

The textual description of the polygonal cropping area is in the following format

	polygon{ points = "<x1 value>, <y1 value>, <x2 value>, <y2 value>,…" }

### "оpacity" property

The "opacity" property (constant Opacity) of the float64 type sets the transparency of the View. Valid values are from 0 to 1.
Where 1 - View is fully opaque, 0 - fully transparent.

You can get the value of this property using the function

	func GetOpacity(view View, subviewID string) float64

### "z-index" property

The "z-index" property (constant ZIndex) of type int defines the position of the element and its children along the z-axis.
In the case of overlapping elements, this value determines the stacking order. In general, the elements
higher z-indexes overlap elements with lower.

You can get the value of this property using the function

	func GetZIndex(view View, subviewID string) int

### "visibility" property

The "visibility" property (constant Visibility) of type int specifies the visibility of the View. Valid values

| Value | Constant  | Name        | Visibility                                     |
|:-----:|-----------|-------------|------------------------------------------------|
| 0     | Visible   | "visible"   | View is visible. Default value.                |
| 1     | Invisible | "invisible" | View is invisible but takes up space.          |
| 2     | Gone      | "gone"      | View is invisible and does not take up space.  |

You can get the value of this property using the function

	func GetVisibility(view View, subviewID string) int

### "filter" property

The "filter" property (Filter constant) applies graphical effects such as blur and color shift to the View.
Only the ViewFilter interface is used as the value of the "filter" property. 
ViewFilter is created using the function

	func NewViewFilter(params Params) ViewFilter

The argument lists the effects to apply. The following effects are possible:

| Effect        | Constant   | Type               | Description             |
|---------------|------------|--------------------|-------------------------|
| "blur"        | Blur       | float64  0…10000px | Gaussian blur           |
| "brightness"  | Brightness | float64  0…10000%  | Brightness change       |
| "contrast"    | Contrast   | float64  0…10000%  | Contrast change         |
| "drop-shadow" | DropShadow | []ViewShadow       | Adding shadow           |
| "grayscale"   | Grayscale  | float64  0…100%    | Converting to grayscale |
| "hue-rotate"  | HueRotate  | AngleUnit          | Hue rotation            |
| "invert"      | Invert     | float64  0…100%    | Invert colors           |
| "opacity"     | Opacity    | float64  0…100%    | Changing transparency   |
| "saturate"    | Saturate   | float64  0…10000%  | Saturation change       |
| "sepia"       | Sepia      | float64  0…100%    | Conversion to serpia    |

Example

    rui.Set(view, "subview", rui.Filter, rui.NewFilter(rui.Params{
        rui.Brightness: 200,
        rui.Contrast: 150,
    }))

You can get the value of the current filter using the function

	func GetFilter(view View, subviewID string) ViewFilter

### "semantics" property

The "semantics" string property (Semantics constant) defines the semantic meaning of the View.
This property may have no visible effect, but it allows search engines to understand the structure of your application.
It also helps to voice the interface to systems for people with disabilities:

| Value | Name             | Semantics                                           |
|:-----:|------------------|-----------------------------------------------------|
| 0     | "default"        | Unspecified. Default value.                         |
| 1     | "article"        | A stand-alone part of the application intended for independent distribution or reuse.               |
| 2     | "section"        | A stand-alone section that cannot be represented by a more precise semantically element             |
| 3     | "aside"          | A part of a document whose content is only indirectly related to the main content (footnote, label) |
| 4     | "header"         | Application Title                                   |
| 5     | "main"           | Main content (content) of the application           |
| 6     | "footer"         | Footer                                              |
| 7     | "navigation"     | Navigation bar                                      |
| 8     | "figure"         | Image                                               |
| 9     | "figure-caption" | Image Title. Should be inside "figure"              |
| 10    | "button"         | Button                                              |
| 11    | "p"              | Paragraph                                           |
| 12    | "h1"             | Level 1 text heading. Changes the style of the text |
| 13    | "h2"             | Level 2 text heading. Changes the style of the text |
| 14    | "h3"             | Level 3 text heading. Changes the style of the text |
| 15    | "h4"             | Level 4 text heading. Changes the style of the text |
| 16    | "h5"             | Level 5 text heading. Changes the style of the text |
| 17    | "h6"             | Level 6 text heading. Changes the style of the text |
| 18    | "blockquote"     | Quote. Changes the style of the text                |
| 19    | "code"           | Program code. Changes the style of the text         |

### Text properties

All properties listed in this section are inherited, i.e. the property will apply 
not only to the View for which it is set, but also to all Views nested in it.

The following properties are available to customize the text display options:

#### "font-name" property

Property "font-name" (constant FontName) - the text property specifies the name of the font to use.
Multiple fonts can be specified. In this case, they are separated by a space.
Fonts are applied in the order in which they are listed. Those, the first is applied first, 
if it is not available, then the second, third, etc.

You can get the value of this property using the function

	func GetFontName(view View, subviewID string) string

#### "text-color" property

Property "text-color" (constant TextColor) - the Color property determines the color of the text.

You can get the value of this property using the function

	func GetTextColor(view View, subviewID string) Color

#### "text-size" property

Property "text-size" (constant TextSize) - the SizeUnit property determines the size of the font.

You can get the value of this property using the function

	func GetTextSize(view View, subviewID string) SizeUnit

#### "italic" property

The "italic" property (constant Italic) is the bool property. If the value is true, then italics are applied to the text

You can get the value of this property using the function

	func IsItalic(view View, subviewID string) bool
	
#### "small-caps" property

The "small-caps" property (SmallCaps constant) is the bool property. If the value is true, then small-caps is applied to the text.

You can get the value of this property using the function

	func IsSmallCaps(view View, subviewID string) bool

#### "white-space" property

The "white-space" (WhiteSpace constant) int property controls how whitespace is handled within the View. 
The "white-space" property can take the following values:

0 (constant WhiteSpaceNormal, name "normal") - sequences of spaces are concatenated into one space.
Newlines in the source are treated as a single space. Applying this value optionally splits lines to fill inline boxes.

1 (constant WhiteSpaceNowrap, name "nowrap") - Concatenates sequences of spaces into one space, 
like a normal value, but does not wrap lines (text wrapping) within the text.

2 (constant WhiteSpacePre, name "pre") - sequences of spaces are saved as they are specified in the source. 
Lines are wrapped only where newlines are specified in the source and where "br" elements are specified in the source.

3 (constant WhiteSpacePreWrap, name "pre-wrap") - sequences of spaces are saved as they are
indicated in the source. Lines are wrapped only where newlines are specified in the source and there,
where "br" elements are specified in the source, and optionally to fill inline boxes.

4 (constant WhiteSpacePreLine, name "pre-line") - sequences of spaces are concatenated into one space.
Lines are split on newlines, on "br" elements, and optionally to fill inline boxes.

5 (constant WhiteSpaceBreakSpaces, name "break-spaces") - the behavior is identical to pre-wrap with the following differences:
* Sequences of spaces are preserved as specified in the source, including spaces at the end of lines.
* Lines are wrapped on any spaces, including in the middle of a sequence of spaces.
* Spaces take up space and do not hang at the ends of lines, which means they affect the internal dimensions (min-content and max-content).

The table below shows the behavior of various values of the "white-space" property.

|                       | New lines | Spaces and Tabs | Text wrapping | End of line spaces | End-of-line other space separators |
|-----------------------|-----------|-----------------|---------------|--------------------|------------------------------------|
| WhiteSpaceNormal      | Collapse  | Collapse        | Wrap          | Remove             | Hang                               |
| WhiteSpaceNowrap      | Collapse  | Collapse        | No wrap       | Remove             | Hang                               |
| WhiteSpacePre         | Preserve  | Preserve        | No wrap       | Preserve           | No wrap                            |
| WhiteSpacePreWrap     | Preserve  | Preserve        | Wrap          | Hang               | Hang                               |
| WhiteSpacePreLine     | Preserve  | Collapse        | Wrap          | Remove             | Hang                               |
| WhiteSpaceBreakSpaces | Preserve  | Preserve        | Wrap          | Wrap               | Wrap                               |

#### "word-break" property

The "word-break" int property (WordBreak constant) determines where the newline will be set if the text exceeds the block boundaries.
The "white-space" property can take the following values:

0 (constant WordBreak, name "normal) - default behavior for linefeed placement.

1 (constant WordBreakAll, name "break-all) - if the block boundaries are exceeded, 
a line break will be inserted between any two characters (except for Chinese/Japanese/Korean text).

2 (constant WordBreakKeepAll, name "keep-all) - Line break will not be used in Chinese/Japanese/ Korean text. 
For text in other languages, the default behavior (normal) will be applied.

3 (constant WordBreakWord, name "break-word) - when the block boundaries are exceeded, 
the remaining whole words can be broken in an arbitrary place, if a more suitable place for line break is not found.

#### "strikethrough", "overline", "underline" properties

These bool properties set decorative lines on the text:

| Property        | Constant      | Decorative line type    |
|-----------------|---------------|-------------------------|
| "strikethrough" | Strikethrough | Strikethrough line text |
| "overline"      | Overline      | Line above the text     |
| "underline"     | Underline     | Line under the text     |

You can get the value of these properties using the functions

	func IsStrikethrough(view View, subviewID string) bool
	func IsOverline(view View, subviewID string) bool
	func IsUnderline(view View, subviewID string) bool

#### "text-line-thickness" property

The "text-line-thickness" SizeUnit property (TextLineThickness constant)  sets the thickness 
of decorative lines on the text set using the "strikethrough", "overline" and "underline" properties.

You can get the value of this property using the function

	GetTextLineThickness(view View, subviewID string) SizeUnit

#### "text-line-style" property

The "text-line-style" int property (constant TextLineStyle) sets the style of decorative lines 
on the text set using the "strikethrough", "overline" and "underline" properties.

Possible values are:

| Value | Constant   | Name     | Description                |
|:-----:|------------|----------|----------------------------|
| 1     | SolidLine  | "solid"  | Solid line (default value) |
| 2     | DashedLine | "dashed" | Dashed line                |
| 3     | DottedLine | "dotted" | Dotted line                |
| 4     | DoubleLine | "double" | Double solid line          |
| 5     | WavyLine   | "wavy"   | Wavy line                  |

You can get the value of this property using the function

	func GetTextLineStyle(view View, subviewID string) int

#### "text-line-color" property

The "text-line-color" Color property (constant TextLineColor) sets the color of decorative lines 
on the text set using the "strikethrough", "overline" and "underline" properties.
If the property is not defined, then the text color specified by the "text-color" property is used for lines.

You can get the value of this property using the function

	func GetTextLineColor(view View, subviewID string) Color

#### "text-weight" property

Свойство "text-weight" (константа TextWeight) - свойство типа int устанавливает начертание шрифта. Допустимые значения:

| Value | Constant       | Common name of the face   |
|:-----:|----------------|---------------------------|
| 1	    | ThinFont       | Thin (Hairline)           |
| 2	    | ExtraLightFont | Extra Light (Ultra Light) |
| 3	    | LightFont      | Light                     |
| 4	    | NormalFont     | Normal. Default value     |
| 5	    | MediumFont     | Medium                    |
| 6	    | SemiBoldFont   | Semi Bold (Demi Bold)     |
| 7	    | BoldFont       | Bold                      |
| 8	    | ExtraBoldFont  | Extra Bold (Ultra Bold)   |
| 9	    | BlackFont      | Black (Heavy)             |

Some fonts are only available in normal or bold style. In this case, the value of this property is ignored.

You can get the value of this property using the function

	func GetTextWeight(view View, subviewID string) int

#### "text-shadow" property

The "text-shadow" property allows you to set shadows for the text. There may be several shadows. 
The shadow is described using the ViewShadow interface (see above, section "The 'shadow' property"). 
For text shadow, only the "color", "x-offset", "y-offset" and "blur" properties are used. 
The "inset" and "spread-radius" properties are ignored (i.e. setting them is not an error, they just have no effect on the text shadow).

To create a ViewShadow for the text shadow, the following functions are used:

	func NewTextShadow(offsetX, offsetY, blurRadius SizeUnit, color Color) ViewShadow
	func NewShadowWithParams(params Params) ViewShadow

The NewShadowWithParams function is used when constants must be used as parameters. For example:

	shadow := NewShadowWithParams(rui.Params{
		rui.ColorProperty : "@shadowColor",
		rui.BlurRadius    : 8.0,
	})

ViewShadow, ViewShadow array, ViewShadow textual representation can be assigned as a value to the "text-shadow" property (see above, section "The 'shadow' property").

You can get the value of this property using the function

	func GetTextShadows(view View, subviewID string) []ViewShadow

If no shadow is specified, then this function will return an empty array

#### "text-align" property

The "text-align" int property (constant TextAlign) sets the alignment of the text. Valid values:

| Value | Constant     | Name      | Value             |
|:-----:|--------------|-----------|-------------------|
| 0     | LeftAlign    | "left"    | Left alignment    |
| 1     | RightAlign   | "right"   | Right alignment   |
| 2     | CenterAlign  | "center"  | Center alignment  |
| 3     | JustifyAlign | "justify" | Justify alignment |

You can get the value of this property using the function

	func GetTextAlign(view View, subviewID string) int

#### "text-indent" property

The "text-indent" (TextIndent constant) SizeUnit property determines the size of the indent (empty space) 
before the first line of text.

You can get the value of this property using the function

	func GetTextIndent(view View, subviewID string) SizeUnit
	
#### "letter-spacing" property

The "Letter-spacing" (LetterSpacing constant) SizeUnit property determines the letter spacing in the text.
The value can be negative, but there can be implementation-specific restrictions.
The user agent can choose not to increase or decrease the letter spacing to align the text.

You can get the value of this property using the function

	func GetLetterSpacing(view View, subviewID string) SizeUnit

#### "word-spacing" property

The "word-spacing" (WordSpacing constant) SizeUnit property determines the length of the space between words.
If the value is specified as a percentage, then it defines the extra spacing as a percentage of the preliminary character width.
Otherwise, it specifies additional spacing in addition to the inner word spacing as defined by the font.

You can get the value of this property using the function

	func GetWordSpacing(view View, subviewID string) SizeUnit

#### "line-height" property

The "line-height" (LineHeight constant) SizeUnit property sets the amount of space between lines.

You can get the value of this property using the function

	func GetLineHeight(view View, subviewID string) SizeUnit

#### "text-transform" property

The "text-transform" (TextTransform constant) int property defines the case of characters. Valid values:

| Value | Constant                | Case conversion                         |
|:-----:|-------------------------|-----------------------------------------|
| 0     | NoneTextTransform       | Original case of characters             |
| 1     | CapitalizeTextTransform | Every word starts with a capital letter |
| 2     | LowerCaseTextTransform  | All characters are lowercase            |
| 3     | UpperCaseTextTransform  | All characters are uppercase            |

You can get the value of this property using the function

	func GetTextTransform(view View, subviewID string) int

#### "text-direction" property

The "text-direction" (TextDirection constant) int property determines the direction of text output. Valid values:

| Value | Constant             | Text output direction                                                   |
|:-----:|----------------------|-------------------------------------------------------------------------|
| 0     | SystemTextDirection  | Systemic direction. Determined by the language of the operating system. |
| 1     | LeftToRightDirection | From left to right. Used for English and most other languages.          |
| 2     | RightToLeftDirection | From right to left. Used for Hebrew, Arabic and some others.            |

You can get the value of this property using the function

	func GetTextDirection(view View, subviewID string) int

#### "writing-mode" property
The "writing-mode" (WritingMode constant) int property defines how the lines of text are arranged 
vertically or horizontally, as well as the direction in which the lines are displayed.
Possible values are:

| Value | Constant              | Description                                                      |
|:-----:|-----------------------|------------------------------------------------------------------|
| 0     | HorizontalTopToBottom | Horizontal lines are displayed from top to bottom. Default value |
| 1     | HorizontalBottomToTop | Horizontal lines are displayed from bottom to top.               |
| 2     | VerticalRightToLeft   | Vertical lines are output from right to left.                    |
| 3     | VerticalLeftToRight   | Vertical lines are output from left to right.                    |

You can get the value of this property using the function

	func GetWritingMode(view View, subviewID string) int

#### "vertical-text-orientation" property

The "vertical-text-orientation" (VerticalTextOrientation constant) int property is used only if "writing-mode" 
is set to VerticalRightToLeft (2) or VerticalLeftToRight (3) and determines the position of the vertical line characters. 
Possible values are:

| Value | Constant               | Value                                        |
|:-----:|------------------------|----------------------------------------------|
| 0     | MixedTextOrientation   | Symbols rotated 90 clockwise. Default value. |
| 1     | UprightTextOrientation | Symbols are arranged normally (vertically).  |

You can get the value of this property using the function

	func GetVerticalTextOrientation(view View, subviewID string) int

### Transformation properties

These properties are used to transform (skew, scale, etc.) the content of the View.

#### "perspective" property

The "perspective" SizeUnit property (Perspective constant) defines the distance between the z = 0 plane and 
the user in order to give the 3D positioned element a perspective effect. Each transformed element with 
z > 0 will become larger, with z < 0, respectively, less.

Elements of the part that are behind the user, i.e. the z-coordinate of these elements is greater than the value of the perspective property, and are not rendered.

The vanishing point is by default located in the center of the element, but it can be moved using the "perspective-origin-x" and "perspective-origin-y" properties.

You can get the value of this property using the function

	func GetPerspective(view View, subviewID string) SizeUnit

#### "perspective-origin-x" and "perspective-origin-y" properties

The "Perspective-origin-x" and "perspective-origin-y" SizeUnit properties (PerspectiveOriginX and PerspectiveOriginY constants)
determine the position from which the viewer is looking. It is used by the perspective property as a vanishing point.

By default, the "perspective-origin-x" and "perspective-origin-y" properties are set to 50%. point to the center of the View.

You can get the value of these properties using the function

	func GetPerspectiveOrigin(view View, subviewID string) (SizeUnit, SizeUnit)

#### "backface-visibility" property

The "backface-visibility" bool property (BackfaceVisible constant) determines whether 
the back face of an element is visible when it is facing the user.

The back surface of an element is a mirror image of its front surface. However, invisible in 2D, 
the back face can be visible when the transformation causes the element to rotate in 3D space.
(This property has no effect on 2D transforms that have no perspective.)

You can get the value of this property using the function

	func GetBackfaceVisible(view View, subviewID string) bool

#### "origin-x", "origin-y", and "origin-z" properties

The "origin-x", "origin-y", and "origin-z" SizeUnit properties (OriginX, OriginY, and OriginZ constants) set the origin for element transformations.

The origin of the transformation is the point around which the transformation takes place. For example, rotation.

The "origin-z" property is ignored if the perspective property is not set.

You can get the value of these properties using the function

	func GetOrigin(view View, subviewID string) (SizeUnit, SizeUnit, SizeUnit)

#### "translate-x", "translate-y", and "translate-z" properties

The "translate-x", "translate-y" and "translate-z" SizeUnit properties (TranslateX, TranslateY, and TranslateZ constants) 
set the offset of the content of the View.

The translate-z property is ignored if the perspective property is not set.

You can get the value of these properties using the function

	func GetTranslate(view View, subviewID string) (SizeUnit, SizeUnit, SizeUnit)

#### "scale-x", "scale-y" and "scale-z" properties

The "scale-x", "scale-y" and "scale-z" float64 properties (ScaleX, ScaleY and ScaleZ constants) set 
the scaling factor along the x, y and z axes, respectively.
The original scale is 1. A value between 0 and 1 is used to zoom out. More than 1 - to increase.
Values less than or equal to 0 are invalid (the Set function will return false)

The "scale-z" property is ignored if the "perspective" property is not set.

You can get the value of these properties using the function

	func GetScale(view View, subviewID string) (float64, float64, float64)

#### "rotate" property

The "rotate" AngleUnit property (Rotate constant) sets the angle of rotation of the content 
around the vector specified by the "rotate-x", "rotate-y" and "rotate-z" properties.

#### "rotate-x", "rotate-y", and "rotate-z" properties

The "rotate-x", "rotate-y" and "rotate-z" float64 properties (constant RotateX, RotateY and RotateZ) set 
the vector around which the rotation is performed by the angle specified by the "rotate" property.
This vector passes through the point specified by the "origin-x", "origin-y" and "origin-z" properties.

The "rotate-z" property is ignored if the "perspective" property is not set.

You can get the value of these properties, as well as the "rotate" property, using the function

	func GetRotate(view View, subviewID string) (float64, float64, float64, AngleUnit)

#### "skew-x" and "skew-y" properties

The "skew-x" and "skew-y" AngleUnit properties (SkewX and SkewY constants) set the skew (skew) of the content,
thus turning it from a rectangle into a parallelogram. The bevel is carried out around the point
specified by the transform-origin-x and transform-origin-y properties.

You can get the value of these properties using the function

	func GetSkew(view View, subviewID string) (AngleUnit, AngleUnit)

### User data

You can save any of your data as "user-data" property (UserData constant)

### Keyboard events

Two kinds of keyboard events can be generated for a View that has received input focus.

| Event            | Constant     | Description                |
|------------------|--------------|----------------------------|
| "key-down-event" | KeyDownEvent | The key has been pressed.  |
| "key-up-event"   | KeyUpEvent   | The key has been released. |

The main event data listener has the following format:

	func(View, KeyEvent)
	
where the second argument describes the parameters of the keys pressed. The KeyEvent structure has the following fields:

| Field     | Type   | Description                                                                                                                               |
|-----------|--------|-------------------------------------------------------------------------------------------------------------------------------------------|
| TimeStamp | uint64 | The time the event was created (in milliseconds). The starting point depends on the browser implementation (EPOCH, browser launch, etc.). |
| Key       | string | The value of the key on which the event occurred. The value is returned taking into account the current language and case. |
| Code      | string | The key code of the represented event. The value is independent of the current language and case.                     |
| Repeat    | bool   | Repeated pressing: the key was pressed until its input began to be automatically repeated.                            |
| CtrlKey   | bool   | The Ctrl key was active when the event occurred.                                                                      |
| ShiftKey  | bool   | The Shift key was active when the event occurred.                                                                     |
| AltKey    | bool   | The Alt (Option or ⌥ in OS X) key was active when the event occurred.                                                 |
| MetaKey   | bool   | The Meta key (for Mac, this is the ⌘ Command key; for Windows, the Windows key ⊞) was active when the event occurred. |

You can also use listeners in the following formats:

* func(KeyEvent)
* func(View)
* func()

You can get lists of listeners for keyboard events using the functions:

	func GetKeyDownListeners(view View, subviewID string) []func(View, KeyEvent)
	func GetKeyUpListeners(view View, subviewID string) []func(View, KeyEvent)

### Focus events

Focus events are fired when a View gains or loses input focus. Accordingly, two events are possible:

| Event              | Constant       | Description                                |
|--------------------|----------------|--------------------------------------------|
| "focus-event"      | FocusEvent     | View receives input focus (becomes active) |
| "lost-focus-event" | LostFocusEvent | View loses input focus (becomes inactive)  |

The main event data listener has the following format:

	func(View).

You can also use a listener in the following format:

	func()

You can get lists of listeners for focus events using the functions:

	func GetFocusListeners(view View, subviewID string) []func(View)
	func GetLostFocusListeners(view View, subviewID string) []func(View)

### Mouse events

Several kinds of mouse events can be generated for the View

| Event                | Constant         | Description                                                            |
|----------------------|------------------|------------------------------------------------------------------------|
| "mouse-down"         | MouseDown        | The mouse button was pressed.                                          |
| "mouse-up"           | MouseUp          | The mouse button has been released.                                    |
| "mouse-move"         | MouseMove        | Mouse cursor moved                                                     |
| "mouse-out"          | MouseOut         | The mouse cursor has moved outside the View, or entered the child View |
| "mouse-over"         | MouseOver        | The mouse cursor has moved within the arrea of View                    |
| "click-event"        | ClickEvent       | There was a mouse click                                                |
| "double-click-event" | DoubleClickEvent | There was a double mouse click                                         |
| "context-menu-event" | ContextMenuEvent | The key for calling the context menu (right mouse button) is pressed   |

The main event data listener has the following format:
	
	func(View, MouseEvent)
	
where the second argument describes the parameters of the mouse event. The MouseEvent structure has the following fields:

| Field     | Type    | Description                                                                             |
|-----------|---------|-----------------------------------------------------------------------------------------|
| TimeStamp | uint64  | The time the event was created (in milliseconds). The starting point depends on the browser implementation (EPOCH, browser launch, etc.). |
| Button    | int     | The number of the mouse button clicked on which triggered the event                     |
| Buttons   | int     | Bitmask showing which mouse buttons were pressed when the event occurred                |
| X         | float64 | The horizontal position of the mouse relative to the origin View                        |
| Y         | float64 | The vertical position of the mouse relative to the origin View                          |
| ClientX   | float64 | Horizontal position of the mouse relative to the upper left corner of the application   |
| ClientY   | float64 | The vertical position of the mouse relative to the upper left corner of the application |
| ScreenX   | float64 | Horizontal position of the mouse relative to the upper left corner of the screen        |
| ScreenY   | float64 | Vertical position of the mouse relative to the upper left corner of the screen          |
| CtrlKey   | bool    | The Ctrl key was active when the event occurred.                                        |
| ShiftKey  | bool    | The Shift key was active when the event occurred.                                       |
| AltKey    | bool    | The Alt (Option or ⌥ in OS X) key was active when the event occurred.                   |
| MetaKey   | bool    | The Meta key (for Mac this is the ⌘ Command key, for Windows is the Windows key ⊞) was active when the event occurred. |

Button field can take the following values

| Value | Constant             | Description |
|:-----:|----------------------|--------------------------------------------------------------------------------------|
| <0    |                      | No buttons pressed                                                                   |
| 0     | PrimaryMouseButton   | Main button. Usually the left mouse button (can be changed in the OS settings)       |
| 1     | AuxiliaryMouseButton | Auxiliary button (wheel or middle mouse button)                                      |
| 2     | SecondaryMouseButton | Secondary button. Usually the right mouse button (can be changed in the OS settings) |
| 3     | MouseButton4         | Fourth mouse button. Usually the browser's Back button                               |
| 4     | MouseButton5         | Fifth mouse button. Usually the browser button Forward                               |

The Button field is a bit mask combining (using OR) the following values

| Value | Constant           | Description      |
|:-----:|--------------------|------------------|
| 1     | PrimaryMouseMask   | Main button      |
| 2     | SecondaryMouseMask | Secondary button |
| 4     | AuxiliaryMouseMask | Auxiliary button |
| 8     | MouseMask4         | Fourth button    |
| 16    | MouseMask5         | Fifth button     |

You can also use listeners in the following formats:

* func(MouseEvent)
* func(View)
* func()

You can get lists of listeners for mouse events using the functions:

	func GetMouseDownListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseUpListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseMoveListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseOverListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseOutListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetClickListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetDoubleClickListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetContextMenuListeners(view View, subviewID string) []func(View, MouseEvent)

### Pointer Events

A pointer is a device-independent representation of input devices (such as a mouse, pen, 
or point of contact on a touch surface). A pointer can point to a specific coordinate 
(or set of coordinates) on a contact surface such as a screen.

All pointers can generate several kinds of events

| Event            | Constant      | Description                                                             |
|------------------|---------------|-------------------------------------------------------------------------|
| "pointer-down"   | PointerDown   | The pointer was pressed.                                                |
| "pointer-up"     | PointerUp     | The pointer was released.                                               |
| "pointer-move"   | PointerMove   | The pointer has been moved                                              |
| "pointer-cancel" | PointerCancel | Pointer events aborted.                                                 |
| "pointer-out"    | PointerOut    | The pointer went out of bounds of the View, or went into the child View |
| "pointer-over"   | PointerOver   | The pointer is within the limits of View                                |

The main event data listener has the following format:

	func(View, PointerEvent)
	
where the second argument describes the parameters of the pointer. PointerEvent structure extends MouseEvent structure
and has the following additional fields:

| Field              | Type    | Description                                                            |
|--------------------|---------|------------------------------------------------------------------------|
| PointerID          | int     | The unique identifier of the pointer that raised the event.            |
| Width              | float64 | The width (X-axis value) in pixels of the pointer's contact geometry.  |
| Height             | float64 | The height (Y-axis value) in pixels of the pointer's contact geometry. |
| Pressure           | float64 | Normalized gauge inlet pressure ranging from 0 to 1, where 0 and 1 represent the minimum and maximum pressure that the hardware is capable of detecting, respectively. |
| TangentialPressure | float64 | Normalized gauge inlet tangential pressure (also known as cylinder pressure or cylinder voltage) ranges from -1 to 1, where 0 is the neutral position of the control. |
| TiltX              | float64 | The planar angle (in degrees, ranging from -90 to 90) between the Y – Z plane and the plane that contains both the pointer (such as a stylus) axis and the Y axis. |
| TiltY              | float64 | The planar angle (in degrees, ranging from -90 to 90) between the X – Z plane and the plane containing both the pointer (such as a stylus) axis and the X axis. |
| Twist              | float64 | Rotation of a pointer (for example, a stylus) clockwise around its main axis in degrees with a value in the range from 0 to 359. |
| PointerType        | string  | the type of device that triggered the event: "mouse", "pen", "touch", etc. |
| IsPrimary          | bool    | a pointer is the primary pointer of this type.                             |

You can also use listeners in the following formats:

* func(PointerEvent)
* func(View)
* func()

You can get lists of pointer event listeners using the functions:

	func GetPointerDownListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerUpListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerMoveListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerCancelListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerOverListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerOutListeners(view View, subviewID string) []func(View, PointerEvent)

### Touch events

These events are used to track multipoint touches. Single touches emulate mouse events.
If you do not need to track multi-point touches, then it is easier to use mouse events

| Event          | Constant    | Description                          |
|----------------|-------------|--------------------------------------|
| "touch-start"  | TouchStart  | The surface touched.                 |
| "touch-end"    | TouchEnd    | Surface touch completed.             |
| "touch-move"   | TouchMove   | One or more touches changed position |
| "touch-cancel" | TouchCancel | The touch is interrupted.            |

The main event data listener has the following format:
	
	func(View, TouchEvent)

where the second argument describes the touch parameters. The TouchEvent structure has the following fields:

| Field     | Type    | Description                                                           |
|-----------|---------|-----------------------------------------------------------------------|
| TimeStamp | uint64  | The time the event was created (in milliseconds). The starting point depends on the browser implementation (EPOCH, browser launch, etc.). |
| Touches   | []Touch | Array of Touch structures, each describing one touch                  |
| CtrlKey   | bool    | The Ctrl key was active when the event occurred.                      |
| ShiftKey  | bool    | The Shift key was active when the event occurred.                     |
| AltKey    | bool    | The Alt (Option or ⌥ in OS X) key was active when the event occurred. |
| MetaKey   | bool    | The Meta key (for Mac, this is the ⌘ Command key; for Windows, the Windows key ⊞) was active when the event occurred. |

The Touch structure describes a single touch and has the following fields

| Field         | Type    | Description                                                                                             |
|---------------|---------|---------------------------------------------------------------------------------------------------------|
| Identifier    | int     | A unique identifier assigned to each touch and does not change until it is completed.                   |
| X             | float64 | The horizontal position of the mouse relative to the origin View                                        |
| Y             | float64 | The vertical position of the mouse relative to the origin View                                          |
| ClientX       | float64 | Horizontal position of the mouse relative to the upper left corner of the application                   |
| ClientY       | float64 | The vertical position of the mouse relative to the upper left corner of the application                 |
| ScreenX       | float64 | Horizontal position of the mouse relative to the upper left corner of the screen                        |
| ScreenY       | float64 | Vertical position of the mouse relative to the upper left corner of the screen                          |
| RadiusX       | float64 | The x-radius of the ellipse, in pixels, that most closely delimits the area of contact with the screen. |
| RadiusY       | float64 | The y-radius of the ellipse, in pixels, that most closely delimits the area of contact with the screen. |
| RotationAngle | float64 | The angle (in degrees) to rotate the ellipse clockwise, described by the radiusX and radiusY parameters, to best cover the contact area between the user and the surface. |
| Force         | float64 | The amount of pressure from 0.0 (no pressure) to 1.0 (maximum pressure) that the user applies to the surface. |

You can also use listeners in the following formats:

* func(TouchEvent)
* func(View)
* func()

You can get lists of listeners for touch events using the functions:

	func GetTouchStartListeners(view View, subviewID string) []func(View, TouchEvent)
	func GetTouchEndListeners(view View, subviewID string) []func(View, TouchEvent)
	func GetTouchMoveListeners(view View, subviewID string) []func(View, TouchEvent)
	func GetTouchCancelListeners(view View, subviewID string) []func(View, TouchEvent)

### Resize-event

The "resize-event" (ResizeEvent constant) is called when the View changes its position and/or size.
The main event data listener has the following format:
	
	func(View, Frame)

where the structure is declared as

	type Frame struct {
		Left, Top, Width, Height float64
	}

Frame elements contain the following data
* Left - the new horizontal offset in pixels relative to the parent View (left position);
* Top - the new vertical offset in pixels relative to the parent View (top position)
* Width - the new width of the visible part of the View in pixels;
* Height - the new height of the visible part of the View in pixels.

You can also use listeners in the following formats:

* func(Frame)
* func(View)
* func()

You can get a list of listeners for this event using the function:

	func GetResizeListeners(view View, subviewID string) []func(View, Frame)

The current position and dimensions of the visible part of the View can be obtained using the View interface function:

	Frame() Frame

or global function

	func GetViewFrame(view View, subviewID string) Frame

### Scroll event

The "scroll-event" (ScrollEvent constant) is raised when the contents of the View are scrolled.
The main event data listener has the following format:
	
	func(View, Frame)

where the Frame elements contain the following data
* Left - the new horizontal shift of the visible area (in pixels);
* Top - the new vertical offset of the visible area (in pixels);
* Width - the total width of the View in pixels;
* Height - the total height of the View in pixels.

You can also use listeners in the following formats:

* func(Frame)
* func(View)
* func()

You can get a list of listeners for this event using the function:

	func GetScrollListeners(view View) []func(View, Frame)

The current position of the viewable area and the overall dimensions of the View can be obtained using the View interface function:

	Scroll() Frame

or global function

	func GetViewScroll(view View, subviewID string) Frame

The following global functions can be used for manual scrolling

	func ScrollViewTo(view View, subviewID string, x, y float64)
	func ScrollViewToStart(view View, subviewID string)
	func ScrollViewToEnd(view View, subviewID string)

which scroll the view, respectively, to the given position, start and end

## ViewsContainer

The ViewsContainer interface, which implements View, describes a container that contains 
several child interface elements (View). ViewsContainer is the base for other containers 
(ListLayout, GridLayout, StackLayout, etc.) and is not used on its own.

In addition to all View properties, this element has only one additional property "content"

### "content" property

The "content" property (constant Сontent) defines an array of child Views. Interface Get function
always returns []View for the given property.

The following 5 data types can be passed as the value of the "content" property:

* View - converted to []View containing one element;

* []View - nil-elements are prohibited, if the array contains nil, then the property will not be set, 
and the Set function will return false and an error message will be written to the log;

* string - if the string is a text representation of the View, then the corresponding View is created, 
otherwise a TextView is created, to which the given string is passed as text.
Next, a []View is created containing the resulting View;

* []string - each element of the array is converted to View as described in the previous paragraph;

* []interface{} - this array must contain only View and string. Each string element is converted to 
a View as described above. If the array contains invalid values, the "content" property will not be set, 
and the Set function will return false and an error message will be written to the log.

You can learn the value of the "content" property using the ViewsContainer interface function

	Views() []View

The following functions of the ViewsContainer interface can be used to edit the "content" property:

	Append(view View)

This function adds an argument to the end of the View list.

	Insert(view View, index uint)

This function inserts an argument at the specified position in the View list. 
If index is greater than the length of the list, then the View is added to the end of the list. 
If index is less than 0, then to the beginning of the list.

	RemoveView(index uint) View

This function removes the View from the given position and returns it. 
If index points outside the bounds of the list, then nothing is removed, and the function returns nil.

## ListLayout

ListLayout is a container that implements the ViewsContainer interface. To create it, use the function

	func NewListLayout(session Session, params Params) ListLayout

Items in this container are arranged as a list. The position of the children can be controlled. 
For this, ListLayout has a number of properties

### "orientation" property

The "orientation" int property (Orientation constant) specifies how the children will be positioned 
relative to each other. The property can take the following values:

| Value | Constant              | Location                                                     |
|:-----:|-----------------------|--------------------------------------------------------------|
| 0     | TopDownOrientation    | Child elements are arranged in a column from top to bottom.  |
| 1     | StartToEndOrientation | Child elements are laid out in a row from beginning to end.  |
| 2     | BottomUpOrientation   | Child elements are arranged in a column from bottom to top.  |
| 3     | EndToStartOrientation | Child elements are laid out in a line from end to beginning. |

The start and end positions for StartToEndOrientation and EndToStartOrientation depend on the value 
of the "text-direction" property. For languages written from right to left (Arabic, Hebrew), 
the beginning is on the right, for other languages - on the left.

### "wrap" property

The "wrap" int property (Wrap constant) defines the position of elements in case of reaching 
the border of the container. There are three options:

* WrapOff (0) - the column / row of elements continues and goes beyond the bounds of the visible area.

* WrapOn (1) - starts a new column / row of items. The new column is positioned towards the end 
(for the position of the beginning and end, see above), the new line is at the bottom.

* WrapReverse (2) - starts a new column / row of elements. The new column is positioned towards the beginning 
(for the position of the beginning and end, see above), the new line is at the top.

### "vertical-align" property

The "vertical-align" property (VerticalAlign constant) of type int sets the vertical
alignment of items in the container. Valid values:

| Value | Constant     | Name      | Alignment        |
|:-----:|--------------|-----------|------------------|
| 0     | TopAlign     | "top"     | Top alignment    |
| 1     | BottomAlign  | "bottom"  | Bottom alignment |
| 2     | CenterAlign  | "center"  | Center alignment |
| 3     | StretchAlign | "stretch" | Height alignment |

### "horizontal-align" property

The "horizontal-align" int property (HorizontalAlign constant) sets the horizontal 
alignment of items in the list. Valid values:

| Value | Constant     | Name      | Alignment        |
|:-----:|--------------|-----------|------------------|
| 0     | LeftAlign    | "left"    | Left alignment   |
| 1     | RightAlign   | "right"   | Right alignment  |
| 2     | CenterAlign  | "center"  | Center alignment |
| 3     | StretchAlign | "stretch" | Width alignment  |

## GridLayout

GridLayout is a container that implements the ViewsContainer interface. To create it, use the function

	func NewGridLayout(session Session, params Params) GridLayout

The container space of this container is split into cells in the form of a table.
All children are located in the cells of the table. A cell is addressed by row and column number. 
Row and column numbers start at 0.

### "column" and "row" properties

The location of the View inside the GridLayout is determined using the "column" and "row" properties.
These properties must be set for each of the child Views.
Child View can span multiple cells within the GridLayout and they can overlap.

The values "column" and "row" can be set by:

* an integer greater than or equal to 0;

* textual representation of an integer greater than or equal to 0 or a constant;

* a Range structure specifying a range of rows / columns:

	type Range struct {
		First, Last int
	}

where First is the number of the first column / row, Last is the number of the last column / row;

* a line of the form "< number of the first column / row >: < number of the last column / row >", 
which is a textual representation of the Range structure

Example

	grid := rui.NewGridLayout(session, rui.Params {
		rui.Content : []View{
			NewView(session, rui.Params {
				rui.ID     : "view1",
				rui.Row    : 0,
				rui.Column : rui.Range{ First: 1, Last: 2 },
			}),
			NewView(session, rui.Params {
				rui.ID     : "view2",
				rui.Row    : "0:2",
				rui.Column : "0",
			}),
		},
	})

In this example, view1 occupies columns 1 and 2 in row 0, and view1 occupies rows 0, 1, and 2 in column 0.

### "cell-width" and "cell-height" properties

By default, the sizes of the cells are calculated based on the sizes of the child Views placed in them.
The "cell-width" and "cell-height" properties (CellWidth and CellHeight constants) allow you to set 
a fixed width and height of cells regardless of the size of the child elements.
These properties are of type []SizeUnit. Each element in the array determines the size of the corresponding column or row.

These properties can be assigned the following data types:

* SizeUnit or textual representation of SizeUnit (or SizeUnit constant). In this case, the corresponding dimensions of all cells are set to the same;

* [] SizeUnit;

* string containing textual representations of SizeUnit (or SizeUnit constants) separated by commas;

* [] string. Each element must be a textual representation of a SizeUnit (or a SizeUnit constant)

* [] interface {}. Each element must either be of type SizeUnit or be a textual representation of SizeUnit (or a SizeUnit constant)

If the number of elements in the "cell-width" and "cell-height" properties is less than the number of columns and rows used, then the missing elements are set to Auto.

The values of the "cell-width" and "cell-height" properties can use the SizeUnit type SizeInFraction.
This type means 1 part. The part is calculated as follows: the size of all cells 
that are not of type SizeInFraction is subtracted from the size of the container, 
and then the remaining size is divided by the number of parts.
The SizeUnit value of type SizeInFraction can be either integer or fractional.

### "grid-row-gap" and "grid-column-gap" properties

The "grid-row-gap" and "grid-column-gap" SizeUnit properties (GridRowGap and GridColumnGap constants) 
allow you to set the distance between the rows and columns of the container, respectively. The default is 0px.

### "cell-vertical-align" property

The "cell-vertical-align" property (constant CellVerticalAlign) of type int sets the vertical alignment of children within the cell they are occupying. Valid values:

| Value | Constant     | Name      | Alignment           |
|:-----:|--------------|-----------|---------------------|
| 0     | TopAlign     | "top"     | Top alignment       |
| 1     | BottomAlign  | "bottom"  | Bottom alignment    |
| 2     | CenterAlign  | "center"  | Center alignment    |
| 3     | StretchAlign | "stretch" | Full height stretch |

The default value is StretchAlign (3)

### "cell-horizontal-align" property

The "cell-horizontal-align" property (constant CellHorizontalAlign) of type int sets the horizontal alignment of children within the occupied cell. Valid values:

| Value | Constant     | Name      | Alignment          | 
|:-----:|--------------|-----------|--------------------|
| 0     | LeftAlign    | "left"    | Left alignment     |
| 1     | RightAlign   | "right"   | Right alignment    |
| 2     | CenterAlign  | "center"  | Center alignment   |
| 3     | StretchAlign | "stretch" | Full width stretch |

The default value is StretchAlign (3)

## ColumnLayout

ColumnLayout is a container that implements the ViewsContainer interface. 
All child Views are arranged in a vertical list aligned to the left or right and split into several columns. 
The alignment depends on the "text-direction" property.

To create the ColumnLayout, use the function

	func NewColumnLayout(session Session, params Params) ColumnLayout

### "column-count" property

The "column-count" int property (ColumnCount constant) sets the number of columns.

If this property is 0 and the "column-width" property is not set, 
then no column splitting is performed and the container is scrolled down.

If the value of this property is greater than 0, then the list is split into columns. 
The column height is equal to the ColumnLayout height, and the width is calculated 
as the ColumnLayout width divided by "column-count". Each next column is located depending 
on the "text-direction" property to the right or left of the previous one, and the container is scrolled horizontally.

You can get the value of this property using the function

	func GetColumnCount(view View, subviewID string) int

### "column-width" property

The "column-width" SizeUnit property (ColumnWidth constant) sets the column width. 
This property is used only if "column-count" is 0, otherwise it is ignored.

IMPORTANT! Percentages cannot be used as the "column-width" value (i.e. if you specify a value in percent, the system will ignore it)

You can get the value of this property using the function

	func GetColumnWidth(view View, subviewID string) SizeUnit

### "column-gap" property

The "column-gap" SizeUnit property (ColumnGap constant) sets the width of the gap between columns.

You can get the value of this property using the function

	func GetColumnGap(view View, subviewID string) SizeUnit

### "column-separator" property

The "column-separator" property (ColumnSeparator constant) allows you to set a line that will be drawn at column breaks. 
The separator line is described by three attributes: line style, thickness, and color.

The value of the "column-separator" property is stored as the ColumnSeparatorProperty interface, 
which implements the Properties interface (see above). ColumnSeparatorProperty can contain the following properties:

| Property | Constant      | Type     | Description    |
|----------|---------------|----------|----------------|
| "style"  | Style         | int      | Line style     |
| "width"  | Width         | SizeUnit | Line thickness |
| "color"  | ColorProperty | Color    | Line color     |

Line style can take the following values:

| Value | Constant   | Name     | Description       |
|:-----:|------------|----------| ------------------|
| 0     | NoneLine   | "none"   | No frame          |
| 1     | SolidLine  | "solid"  | Solid line        |
| 2     | DashedLine | "dashed" | Dashed line       |
| 3     | DottedLine | "dotted" | Dotted line       |
| 4     | DoubleLine | "double" | Double solid line |

All other style values are ignored.

To create the ColumnSeparatorProperty interface, use the function

	func NewColumnSeparator(params Params) ColumnSeparatorProperty

The ColumnSeparatorProperty interface can be converted to a ViewBorder structure using the ViewBorder function. 
When converted, all text constants are replaced with real values. ViewBorder is described as

	type ViewBorder struct {
		Style int
		Color Color
		Width SizeUnit
	}

The ViewBorder structure can be passed as a parameter to the Set function when setting 
the value of the "column-separator" property. This converts the ViewBorder to ColumnSeparatorProperty.
Therefore, when reading the property, the Get function will return the ColumnSeparatorProperty interface, 
not the ViewBorder structure. 

You can get the ViewBorders structure without additional transformations using the global function

	func GetColumnSeparator(view View, subviewID string) ViewBorder

You can also set individual line attributes using the Set function of the View interface.
For this, the following properties are used

| Property                 | Constant             | Type     | Description    |
|--------------------------|----------------------|----------|----------------|
| "column-separator-style" | ColumnSeparatorStyle | int      | Line style     |
| "column-separator-width" | ColumnSeparatorWidth | SizeUnit | Line thickness |
| "column-separator-color" | ColumnSeparatorColor | Color    | Line color     |

For example

	view.Set(rui.ColumnSeparatorStyle, rui.SolidBorder)
	view.Set(rui.ColumnSeparatorWidth, rui.Px(1))
	view.Set(rui.ColumnSeparatorColor, rui.Black)

equivalent to

	view.Set(rui.ColumnSeparator, ColumnSeparatorProperty(rui.Params{
		rui.Style: rui.SolidBorder,
		rui.Width: rui.Px(1),
		rui.ColorProperty: rui.Black,
	}))

### "avoid-break" property

When forming columns, ColumnLayout can break some types of View, so that the beginning 
will be at the end of one column and the end in the next. For example, the TextView, 
the title of the picture and the picture itself are broken, etc.

The "avoid-break" bool property (AvoidBreak constant) avoids this effect.
You must set this property to "true" for a non-breakable View.
Accordingly, the value "false" of this property allows the View to be broken.
The default is "false".

You can get the value of this property using the function

	func GetAvoidBreak(view View, subviewID string) bool

## StackLayout

StackLayout is a container that implements the ViewsContainer interface. 
All child Views are stacked on top of each other and each takes up the entire container space. 
Only one child View (current) is available at a time.

To create a StackLayout, use the function

	func NewStackLayout(session Session, params Params) StackLayout

In addition to the Append, Insert, RemoveView properties and the "content" property of the ViewsContainer, 
the StackLayout container has two other interface functions for manipulating child Views: Push and Pop.

	Push(view View, animation int, onPushFinished func())

This function adds a new View to the container and makes it current. 
It is similar to Append, but the addition is done using an animation effect. 
The animation type is specified by the second argument and can take the following values:

| Value | Constant            | Animation                   |
|:-----:|---------------------|-----------------------------|
| 0     | DefaultAnimation    | Default animation. For the Push function it is EndToStartAnimation, for Pop - StartToEndAnimation |
| 1     | StartToEndAnimation | Animation from beginning to end. The beginning and the end are determined by the direction of the text output |
| 2     | EndToStartAnimation | End-to-Beginning animation. |
| 3     | TopDownAnimation    | Top-down animation.         |
| 4     | BottomUpAnimation   | Bottom up animation.        |

The third argument onPushFinished is the function to be called when the animation ends. It may be nil.

	Pop(animation int, onPopFinished func(View)) bool

This function removes the current View from the container using animation.
The second argument onPopFinished is the function to be called when the animation ends. It may be nil.
The function will return false if the StackLayout is empty and true if the current item has been removed.

 You can get the current (visible) View using the interface function

	Peek() View

You can also get the current View using its index. The "current" property (constant Current) is used to get the index. 
Example

	func peek(layout rui.StackLayout) {
		views := layout.Views()
		if index := rui.GetCurrent(layout, ""); index >= 0 && index < len(views) {
			return views[index]
		} 
		return nil
	}

Of course, this is less convenient than the Peek function. However, the "current" property can be used to track changes to the current View:

	layout.SetChangeListener(rui.Current, func(view rui.View, tag string) {
		// current view changed
	})

In order to make any child View current (visible), the interface functions are used:

	MoveToFront(view View) bool
	MoveToFrontByID(viewID string) bool

This function will return true if successful and false if the child View or 
View with id does not exist and an error message will be written to the log.

You can also use the "current" property to make any child View current (visible).

## TabsLayout

TabsLayout is a container that implements the ViewsContainer interface. All child Views are stacked 
on top of each other and each takes up the entire container space. 
Only one child View (current) is available at a time. Tabs, that are located along one of the sides of the container, 
are used to select the current View.

To create a TabsLayout, use the function

	func NewTabsLayout(session Session, params Params) TabsLayout

A bookmark is created for each View. A bookmark can display a title, an icon, and a close button.

The title is set using the "title" text property (constant Title) of the child View.
The "title" property is optional. If it is not specified, then there will be no text on the tab.

The icon is set using the "icon" text property (constant Icon) of the child View.
As a value, it is assigned the name of the icon file (if the icon is located in the application resources) or url. 
The "icon" property is optional. If it is not specified, then there will be no icon on the tab.

The display of the tab close button is controlled by the "tab-close-button" boolean property (constant TabCloseButton).
"true" enables the display of the close button for the tab. The default is "false".

The "tab-close-button" properties can be set for both the child View and the TabsLayout itself.
Setting the value of the "tab-close-button" property for the TabsLayout enables/disables the display 
of the close button for all tabs at once. The "tab-close-button" value set on the child View 
takes precedence over the value set on the TabsLayout.

The tab close button does not close the tab, but only generates the "tab-close-event" event (constant TabCloseEvent).
The main handler for this event has the format

	func(layout TabsLayout, index int)

where the second element is the index of the child View.

As already mentioned, clicking on the close tab button does not close the tab.
You must close the tab yourself. This is done as follows

	tabsView.Set(rui.TabCloseEvent, func(layout rui.TabsLayout, index int) {
		layout.RemoveView(index)
	})

You can control the current View using the "current" integer property (constant Current).
To programmatically switch tabs, set this property to the index of the new current View.
You can read the value of the "current" property using the function

	func GetCurrent(view View, subviewID string) int

Also, the "current" property can be used to track changes to the current View:

	tabsView.SetChangeListener(rui.Current, func(view rui.View, tag string) {
		// current view changed
	})

Tabs are positioned along one side of the TabsLayout container. The tabs are positioned using 
the "tabs" integer property (the Tabs constant). This property can take on the following values:

| Value    | Constant      | Name         | Placement of tabs                                |
|:--------:|---------------|--------------|--------------------------------------------------|
| 0	       | TopTabs       | "top"        | Top. Default value.                              |
| 1        | BottomTabs    | "bottom"     | Bottom.                                          |
| 2        | LeftTabs      | "left"       | Left. Each tab is rotated 90 ° counterclockwise. |
| 3        | RightTabs     | "right"      | On right. Each tab is rotated 90 ° clockwise.    |
| 4        | LeftListTabs  | "left-list"  | Left. The tabs are displayed as a list.          |
| 5        | RightListTabs | "right-list" | On right. The tabs are displayed as a list.      |
| 6        | HiddenTabs    | "hidden"     | The tabs are hidden.                             |

Why do I need the value HiddenTabs. The point is that TabsLayout implements the ListAdapter interface.
Which makes it easy to implement tabs with a ListView. This is where the HiddenTabs value comes in.

For displaying the current (selected) tab of type TopTabs, BottomTabs, LeftListTabs and RightListTabs, 
the "ruiCurrentTab" style is used, and for tab of type LeftTabs and RightTabs, the "ruiCurrentVerticalTab" style is used.
If you want to customize the display of tabs, you can either override these styles, or assign your own style using 
the "current-tab-style" property (constant CurrentTabStyle).

Accordingly, for an inactive tab, the "ruiTab" and "ruiVerticalTab" styles are used, 
and you can assign your own style using the "tab-style" property (constant TabStyle).

The "ruiTabBar" style is used to display the tab bar, and you can assign your own style 
using the "tab-bar-style" property (constant TabBarStyle).

## AbsoluteLayout

AbsoluteLayout is a container that implements the ViewsContainer interface. 
Child Views can be positioned at arbitrary positions in the container space.

To create an AbsoluteLayout, use the function

	func NewAbsoluteLayout(session Session, params Params) AbsoluteLayout

The child View is positioned using the properties of the SizeUnit type: "left", "right", "top" and "bottom" 
(respectively, the constants Left, Right, Top and Bottom). You can set any of these properties on the child View. 
If neither "left" or "right" is specified, then the child View will be pinned to the left edge of the container. 
If neither top nor bottom is specified, then the child View will be pinned to the top edge of the container.

## DetailsView

DetailsView is a container that implements the ViewsContainer interface.
To create a DetailsView, the function is used

	func NewDetailsView(session Session, params Params) DetailsView

In addition to child Views, this container has a "summary" property (Summary constant).
The value of the "summary" property can be either View or a string of text.

The DetailsView can be in one of two states:

* only the content of the "summary" property is displayed. Child Views are hidden and do not take up screen space

* the content of the "summary" property is displayed first, and below the child Views.
The layout of the child Views is the same as ColumnLayout with "column-count" equal to 0.

DetailsView switches between states by clicking on "summary" view.

For forced switching of the DetailsView states, the bool property "expanded" (Expanded constant) is used. 
Accordingly, the value "true" shows child Views, "false" - hides.

You can get the value of the "expanded" property using the function

	func IsDetailsExpanded(view View, subviewID string) bool

and the value of the "summary" property can be obtained using the function

	func GetDetailsSummary(view View, subviewID string) View

## Resizable

Resizable is a container in which only one View can be placed. Resizable allows you to interactively resize the content View.
To create a Resizable view, the function is used

	func NewResizable(session Session, params Params) Resizable

A frame is created around the content View, and you can drag it to resize.

Resizable does not implement the ViewsContainer interface. Only the Content property is used to control the content View. 
This property can be assigned a value of type View or a string of text. In the second case, a TextView is created.

The frame around the content View can be either from all sides, or only from separate ones.
To set the sides of the frame, use the "side" int property (Side constant).
It can take the following values:

| Value | Constant   | Name     | Frame side               |
|:-----:|------------|----------|--------------------------|
| 1     | TopSide    | "top"    | Top                      |
| 2     | RightSide  | "right"  | Right                    |
| 4     | BottomSide | "bottom" | Bottom                   |
| 8     | LeftSide   | "left"   | Left                     |
| 15    | AllSides   | "all"    | All sides. Default value |

In addition to these values, an or-combination of TopSide, RightSide, BottomSide and LeftSide can also be used.
AllSides is defined as

	AllSides = TopSide | RightSide | BottomSide | LeftSide

To set the border width, use the SizeUnit property "resize-border-width" (ResizeBorderWidth constant).
The default value of "resize-border-width" is 4px.

## TextView

The TextView element, which extends the View interface, is intended for displaying text.

To create a TextView, the function is used:

    func NewTextView(session Session, params Params) TextView

The displayed text is set by the string property "text" (Text constant).
In addition to the Get method, the value of the "text" property can be obtained using the function

    func GetText (view View, subviewID string) string

TextView inherits from View all properties of text parameters ("font-name", "text-size", "text-color", etc.).
In addition to them, the "text-overflow" int property (TextOverflow constant) is added. 
It determines how the text is cut if it goes out of bounds. 
This property of type int can take the following values

| Value | Constant             | Name       | Cropping Text                                               |
|:-----:|----------------------| -----------|-------------------------------------------------------------|
| 0     | TextOverflowClip     | "clip"     | Text is clipped at the border (default)                     |
| 1     | TextOverflowEllipsis | "ellipsis" | At the end of the visible part of the text '…' is displayed |

## EditView

The EditView element is a test editor and extends the View interface.

To create an EditView, the function is used:

	func NewEditView(session Session, params Params) EditView

Several options for editable text are possible. The type of the edited text is set using 
the int property "edit-view-type" (EditViewType constant).
This property can take the following values:

| Value | Constant       | Name        | Editor type                                      |
|:-----:|----------------|-------------|--------------------------------------------------|
| 0     | SingleLineText | "text"      | One-line text editor. Default value              |
| 1     | PasswordText   | "password"  | Password editor. The text is hidden by asterisks |
| 2     | EmailText      | "email"     | Single e-mail editor                             |
| 3     | EmailsText     | "emails"    | Multiple e-mail editor                           |
| 4     | URLText        | "url"       | Internet address input editor                    |
| 5     | PhoneText      | "phone"     | Phone number editor                              |
| 6     | MultiLineText  | "multiline" | Multi-Line Text Editor                           |

To simplify the text of the program, you can use the "type" properties (Type constant) instead of the "edit-view-type".
These property names are synonymous. But when describing the style, "type" cannot be used.

To set/get edited text, use the string property "text" (Text constant)

The maximum length of editable text is set using the "max-length" int property (MaxLength constant).

You can limit the input text using a regular expression. To do this, use the string property 
"edit-view-pattern" (EditViewPattern constant). Instead of "edit-view-pattern", you can use the synonym "pattern" 
(Pattern constant), except for the style description.

To prohibit text editing, use the bool property "readonly" (ReadOnly constant).

To enable / disable the built-in spell checker, use the bool "spellcheck" property (Spellcheck constant). 
Spell checking can only be enabled if the editor type is set to SingleLineText or MultiLineText.

For the editor, you can set a hint that will be shown while the editor is empty.
To do this, use the string property "hint" (Hint constant).

For a multi-line editor, auto-wrap mode can be enabled. The bool property "wrap" (constant Wrap) is used for this. 
If "wrap" is off (default), then horizontal scrolling is used. 
If enabled, the text wraps to a new line when the EditView border is reached.

The following functions can be used to get the values of the properties of an EditView:

	func GetText(view View, subviewID string) string
	func GetHint(view View, subviewID string) string
	func GetMaxLength(view View, subviewID string) int
	func GetEditViewType(view View, subviewID string) int
	func GetEditViewPattern(view View, subviewID string) string
	func IsReadOnly(view View, subviewID string) bool
	func IsEditViewWrap(view View, subviewID string) bool
	func IsSpellcheck(view View, subviewID string) bool

The "edit-text-changed" event (EditTextChangedEvent constant) is used to track changes to the text. 
The main event listener has the following format:

	func(EditView, string)

where the second argument is the new text value

You can get the current list of text change listeners using the function

	func GetTextChangedListeners(view View, subviewID string) []func(EditView, string)

## NumberPicker

The NumberPicker element extends the View interface to enter numbers.

To create a NumberPicker, the function is used:

	func NewNumberPicker(session Session, params Params) NumberPicker

NumberPicker can work in two modes: text editor and slider.
The mode sets the int property "number-picker-type" (NumberPickerType constant).
The "number-picker-type" property can take the following values:

| Value | Constant     | Name     | Editor type                |
|:-----:|--------------|----------|----------------------------|
| 0     | NumberEditor | "editor" | Text editor. Default value |
| 1     | NumberSlider | "slider" | Slider                     |

You can set/get the current value using the "number-picker-value" property (NumberPickerValue constant). 
The following can be passed as a value to the "number-picker-value" property:

* float64
* float32
* int
* int8 … int64
* uint
* uint8 … uint64
* textual representation of any of the above types

All of these types are cast to float64. Accordingly, the Get function always returns a float64 value.
The value of the "number-picker-value" property can also be read using the function:

	func GetNumberPickerValue(view View, subviewID string) float64

The entered values may be subject to restrictions. For this, the following properties are used:

| Property             | Constant         | Restriction       |
|----------------------|------------------|-------------------|
| "number-picker-min"  | NumberPickerMin  | Minimum value     |
| "number-picker-max"  | NumberPickerMax  | Maximum value     |
| "number-picker-step" | NumberPickerStep | Value change step |

Assignments to these properties can be the same value types as "number-picker-value".

By default, if "number-picker-type" is equal to NumberSlider, the minimum value is 0, maximum is 1. 
If "number-picker-type" is equal to NumberEditor, then the entered numbers, by default, are limited only by the range of float64 values.

You can read the values of these properties using the functions:

	func GetNumberPickerMinMax(view View, subviewID string) (float64, float64)
	func GetNumberPickerStep(view View, subviewID string) float64

The "number-changed" event (NumberChangedEvent constant) is used to track the change in the entered value. 
The main event listener has the following format:

	func(picker NumberPicker, newValue float64)

where the second argument is the new value

You can get the current list of value change listeners using the function

	func GetNumberChangedListeners(view View, subviewID string) []func(NumberPicker, float64)

## DatePicker

The DatePicker element extends the View interface to enter dates.

To create DatePicker function is used:

	func NewDatePicker(session Session, params Params) DatePicker

You can set/get the current value using the "date-picker-value" property (the DatePickerValue constant). 
The following can be passed as a value to the "date-picker-value" property:

* time.Time
* constant
* text that can be converted to time.Time by function

	func time.Parse(layout string, value string) (time.Time, error)

The text is converted to time.Time. Accordingly, the Get function always returns a time.Time value.
The value of the "date-picker-value" property can also be read using the function:

	func GetDatePickerValue(view View, subviewID string) time.Time

The dates you enter may be subject to restrictions. For this, the following properties are used:

| Property           | Constant       | Data type | Restriction              |
|--------------------|----------------|-----------|--------------------------|
| "date-picker-min"  | DatePickerMin  | time.Time | Minimum date value       |
| "date-picker-max"  | DatePickerMax  | time.Time | Maximum date value       |
| "date-picker-step" | DatePickerStep | int       | Date change step in days |

You can read the values of these properties using the functions:

	func GetDatePickerMin(view View, subviewID string) (time.Time, bool)
	func GetDatePickerMax(view View, subviewID string) (time.Time, bool)
	func GetDatePickerStep(view View, subviewID string) int

The "date-changed" event (DateChangedEvent constant) is used to track the change in the entered value. 
The main event listener has the following format:

	func(picker DatePicker, newDate time.Time)

where the second argument is the new date value

You can get the current list of date change listeners using the function

	func GetDateChangedListeners(view View, subviewID string) []func(DatePicker, time.Time)

## TimePicker

The TimePicker element extends the View interface and is intended for entering time.

To create a TimePicker, the function is used:

	func NewTimePicker(session Session, params Params) TimePicker

You can set/get the current value using the "time-picker-value" property (TimePickerValue constant). 
The following can be passed as a value to the "time-picker-value" property:

* time.Time
* constant
* text that can be converted to time.Time by function

    func time.Parse(layout string, value string) (time.Time, error)

The text is converted to time.Time. Accordingly, the Get function always returns a time.Time value.
The value of the "time-picker-value" property can also be read using the function:

	func GetTimePickerValue(view View, subviewID string) time.Time

The time entered may be subject to restrictions. For this, the following properties are used:

| Property           | Constant       | Data type | Restriction               |
|--------------------|----------------|-----------|---------------------------|
| "time-picker-min"  | TimePickerMin  | time.Time | Minimum time value        |
| "time-picker-max"  | TimePickerMax  | time.Time | The maximum value of time |
| "time-picker-step" | TimePickerStep | int       | Time step in seconds      |

You can read the values of these properties using the functions:

	func GetTimePickerMin(view View, subviewID string) (time.Time, bool)
	func GetTimePickerMax(view View, subviewID string) (time.Time, bool)
	func GetTimePickerStep(view View, subviewID string) int

The "time-changed" event (TimeChangedEvent constant) is used to track the change in the entered value. 
The main event listener has the following format:

	func(picker TimePicker, newTime time.Time)

where the second argument is the new time value

You can get the current list of date change listeners using the function

	func GetTimeChangedListeners(view View, subviewID string) []func(TimePicker, time.Time)

## ColorPicker

The ColorPicker element extends the View interface and is designed to select a color in RGB format without an alpha channel.

To create a ColorPicker, the function is used:

	func NewColorPicker(session Session, params Params) ColorPicker

You can set/get the current color value using the "color-picker-value" property (ColorPickerValue constant). 
The following can be passed as a value to the "color-picker-value" property:

* Color
* text representation of Color
* constant

The value of the property "color-picker-value" can also be read using the function:

	func GetColorPickerValue(view View, subviewID string) Color

The "color-changed" event (ColorChangedEvent constant) is used to track the change in the selected color. 
The main event listener has the following format:

	func(picker ColorPicker, newColor Color)

where the second argument is the new color value

You can get the current list of date change listeners using the function

	func GetColorChangedListeners(view View, subviewID string) []func(ColorPicker, Color)

## FilePicker

The FilePicker element extends the View interface to select one or more files.

To create a FilePicker, the function is used:

	func NewFilePicker(session Session, params Params) FilePicker

The boolean property "multiple" (constant Multiple) is used to set the mode of selecting multiple files.
The value "true" enables the selection of multiple files, "false" enables the selection of a single file.
The default is "false".

You can restrict the selection to only certain types of files. To do this, use the "accept" property (constant Accept).
This property is assigned a list of allowed file extensions and / or mime-types. 
The value can be specified either as a string (elements are separated by commas) or as an array of strings. Examples

	rui.Set(view, "myFilePicker", rui.Accept, "png, jpg, jpeg")
	rui.Set(view, "myFilePicker", rui.Accept, []string{"png", "jpg", "jpeg"})
	rui.Set(view, "myFilePicker", rui.Accept, "image/*")
	
Two functions of the FilePicker interface are used to access the selected files:

	Files() []FileInfo
	LoadFile(file FileInfo, result func(FileInfo, []byte))

as well as the corresponding global functions

	func GetFilePickerFiles(view View, subviewID string) []FileInfo
	func LoadFilePickerFile(view View, subviewID string, file FileInfo, result func(FileInfo, []byte))

The Files/GetFilePickerFiles functions return a list of the selected files as a slice of FileInfo structures. 
The FileInfo structure is declared as

	type FileInfo struct {
		// Name - the file's name.
		Name string
		// LastModified specifying the date and time at which the file was last modified
		LastModified time.Time
		// Size - the size of the file in bytes.
		Size int64
		// MimeType - the file's MIME type.
		MimeType string
	}

FileInfo contains only information about the file, not the file content. 
The LoadFile/LoadFilePickerFile function allows you to load the contents of one of the selected files. 
The LoadFile function is asynchronous. After loading, the contents of the selected file are passed to the argument-function of the LoadFile. 
Example

	if filePicker := rui.FilePickerByID(view, "myFilePicker"); filePicker != nil {
		if files := filePicker.Files(); len(files) > 0 {
			filePicker.LoadFile(files[0], func(file rui.FileInfo, data []byte) {
				if data != nil {
					// ... 
				}
			})
		}
	}

equivalent to

	if files := rui.GetFilePickerFiles(view, "myFilePicker"); len(files) > 0 {
		rui.LoadFilePickerFile(view, "myFilePicker", files[0], func(file rui.FileInfo, data []byte) {
			if data != nil {
				// ... 
			}
		})
	}

If an error occurs while loading the file, the data value passed to the result function will be nil, 
and the error description will be written to the log

The "file-selected-event" event (constant FileSelectedEvent) is used to track changes in the list of selected files. 
The main event listener has the following format:

	func(picker FilePicker, files []FileInfo))

where the second argument is the new value of the list of selected files.

You can get the current list of listeners of the list of files changing using the function

	func GetFileSelectedListeners(view View, subviewID string) []func(FilePicker, []FileInfo)

## DropDownList

The DropDownList element extends the View interface and is designed to select a value from a drop-down list.

To create a DropDownList, use the function:

	func NewDropDownList(session Session, params Params) DropDownList

The list of possible values is set using the "items" property (Items constant).
The following data types can be passed as a value to the "items" property

* []string
* []fmt.Stringer
* []interface{} containing as elements only: string, fmt.Stringer, bool, rune,
float32, float64, int, int8 … int64, uint, uint8 … uint64.

All of these data types are converted to []string and assigned to the "items" property.
You can read the value of the "items" property using the function

	func GetDropDownItems(view View, subviewID string) []string

The selected value is determined by the int property "current" (Current constant). The default is 0.
You can read the value of this property using the function

	func GetCurrent(view View, subviewID string) int

To track the change of the "current" property, the "drop-down-event" event (DropDownEvent constant) is used. 
The main event listener has the following format:

	func(list DropDownList, newCurrent int)

where the second argument is the index of the selected item

You can get the current list of date change listeners using the function

	func GetDropDownListeners(view View, subviewID string) []func(DropDownList, int)

## ProgressBar

The DropDownList element extends the View interface and is designed to display progress as a fillable bar.

To create a ProgressBar, the function is used:

	func NewProgressBar(session Session, params Params) ProgressBar

ProgressBar has two float64 properties:
* "progress-max" (ProgressBarMax constant) - maximum value (default 1);
* "progress-value" (ProgressBarValue constant) - current value (default 0).

The minimum is always 0.
In addition to float64, float32, int, int8 … int64, uint, uint8 … uint64

You can read the value of these properties using the functions

	func GetProgressBarMax(view View, subviewID string) float64
	func GetProgressBarValue(view View, subviewID string) float64

## Button

The Button element implements a clickable button. This is a CustomView (about it below) based on ListLayout and, 
accordingly, has all the properties of ListLayout. But unlike ListLayout, it can receive input focus.

Content is centered by default.

To create a Button, use the function:

	func NewButton(session Session, params Params) Button

## ListView

The ListView element implements a list.
The ListView is created using the function:

	func NewListView(session Session, params Params) ListView

### The "items" property

List items are set using the "items" property (Items constant). 
The main value of the "items" property is the ListAdapter interface:

	type ListAdapter interface {
		ListSize() int
		ListItem(index int, session Session) View
		IsListItemEnabled(index int) bool
	}

Accordingly, the functions of this interface must return the number of elements, 
the View of the i-th element and the status of the i-th element (allowed/denied).

You can implement this interface yourself or use helper functions:

	func NewTextListAdapter(items []string, params Params) ListAdapter
	func NewViewListAdapter(items []View) ListAdapter

NewTextListAdapter creates an adapter from an array of strings, the second argument 
is the parameters of the TextView used to display the text. 
NewViewListAdapter creates an adapter from the View array.

The "items" property can be assigned the following data types:

* ListAdapter;
* [] View, when assigned, is converted to a ListAdapter using the NewViewListAdapter function;
* [] string, when assigned, is converted to a ListAdapter using the NewTextListAdapter function;
* [] interface {} which can contain elements of type View, string, fmt.Stringer, bool, rune, 
float32, float64, int, int8 ... int64, uint, uint8 ... uint64. 
When assigning, all types except View and string are converted to string, then all string in TextView 
and from the resulting View array using the NewViewListAdapter function, a ListAdapter is obtained.

If the list items change during operation, then after the change, either the ReloadListViewData() 
function of the ListView interface or the global ReloadListViewData(view View, subviewID string) function must be called.
These functions update the displayed list items.

### "Orientation" property

List items can be arranged both vertically (in columns) and horizontally (in rows).
The "orientation" property (Orientation constant) of int type specifies how the list items 
will be positioned relative to each other. The property can take the following values:

| Value | Constant              | Location                                               |
|:-----:|-----------------------|--------------------------------------------------------|
| 0     | TopDownOrientation    | Items are arranged in a column from top to bottom.     |
| 1     | StartToEndOrientation | Elements are laid out on a line from beginning to end. |
| 2     | BottomUpOrientation   | Items are arranged in a column from bottom to top.     |
| 3     | EndToStartOrientation | Elements are arranged in a row from end to beginning.  |

The start and end positions for StartToEndOrientation and EndToStartOrientation depend 
on the value of the "text-direction" property. For languages written from right to left 
(Arabic, Hebrew), the beginning is on the right, for other languages - on the left.

You can get the value of this property using the function

	func GetListOrientation(view View, subviewID string) int

### "wrap" property

The "wrap" int property (Wrap constant) defines the position of elements 
in case of reaching the border of the container. There are three options:

* WrapOff (0) - the column/row of elements continues and goes beyond the bounds of the visible area.

* WrapOn (1) - starts a new column/row of items. The new column is positioned towards the end 
(for the position of the beginning and end, see above), the new line is at the bottom.

* WrapReverse (2) - starts a new column/row of elements. The new column is positioned towards 
the beginning (for the position of the beginning and end, see above), the new line is at the top.

You can get the value of this property using the function

	func GetListWrap(view View, subviewID string) int

### "item-width" and "item-height" properties

By default, the height and width of list items are calculated based on their content.
This leads to the fact that the elements of the vertical list can have different heights, 
and the elements of the horizontal - different widths.

You can set a fixed height and width of the list item. To do this, use the SizeUnit 
properties "item-width" and "item-height"

You can get the values of these properties using the functions

	func GetListItemWidth(view View, subviewID string) SizeUnit
	func GetListItemHeight(view View, subviewID string) SizeUnit

### "item-vertical-align" property

The "item-vertical-align" int property (ItemVerticalAlign constant) sets the vertical alignment 
of the contents of the list items. Valid values:

| Value | Constant     | Name      | Alignment        |
|:-----:|--------------|-----------|------------------|
| 0     | TopAlign     | "top"     | Top alignment    |
| 1     | BottomAlign  | "bottom"  | Bottom alignment |
| 2     | CenterAlign  | "center"  | Center alignment |
| 3     | StretchAlign | "stretch" | Height alignment |

You can get the value of this property using the function

	func GetListItemVerticalAlign(view View, subviewID string) int

### "item-horizontal-align" property

The "item-horizontal-align" int property (ItemHorizontalAlign constant) sets the 
horizontal alignment of the contents of the list items. Valid values:

| Value | Constant     | Name      | Alignment        |
|:-----:|--------------|-----------|------------------|
| 0	    | LeftAlign    | "left"    | Left alignment   |
| 1     | RightAlign   | "right"   | Right alignment  |
| 2     | CenterAlign  | "center"  | Center alignment |
| 3     | StretchAlign | "stretch" | Height alignment |

You can get the value of this property using the function

	GetListItemHorizontalAlign(view View, subviewID string) int

### "current" property

ListView allows you to select list items with the "allowed" status (see ListAdapter).
The item can be selected both interactively and programmatically. 
To do this, use the int property "current" (constant Current). 
The value "current" is less than 0 means that no item is selected

You can get the value of this property using the function

	func GetCurrent(view View, subviewID string) int

### "list-item-style", "current-style", and "current-inactive-style" properties

These three properties are responsible for the background style and text properties of each list item.

| Property                 | Constant             | Style                                                    |
|--------------------------|----------------------|----------------------------------------------------------|
| "list-item-style"        | ListItemStyle        | Unselected element style                                 |
| "current-style"          | CurrentStyle         | The style of the selected item. ListView in focus        |
| "current-inactive-style" | CurrentInactiveStyle | The style of the selected item. ListView is out of focus |

### "checkbox", "checked", "checkbox-horizontal-align", and "checkbox-vertical-align" properties

The "current" property allows you to select one item in the list.
The "checkbox" properties allow you to add a checkbox to each item in the list with which 
you can select several items in the list. The "checkbox" int property (ItemCheckbox constant) 
can take the following values

| Value    | Constant         | Name       | Checkbox view                                      |
|:--------:|------------------|------------|----------------------------------------------------|
| 0	       | NoneCheckbox     | "none"     | There is no checkbox. Default value                |
| 1	       | SingleCheckbox   | "single"   | ◉ A checkbox that allows you to mark only one item |
| 2	       | MultipleCheckbox | "multiple" | ☑ A checkbox that allows you to mark several items |


You can get the value of this property using the function

	func GetListViewCheckbox(view View, subviewID string) int

You can get/set the list of checked items using the "checked" property (Checked constant).
This property is of type []int and stores the indexes of the marked elements.
You can get the value of this property using the function

	func GetListViewCheckedItems(view View, subviewID string) []int

You can check if a specific element is marked using the function

	func IsListViewCheckedItem(view View, subviewID string, index int) bool

By default, the checkbox is located in the upper left corner of the element. 
You can change its position using int properties "checkbox-horizontal-align" and 
"checkbox-vertical-align" (CheckboxHorizontalAlign and CheckboxVerticalAlign constants)

The "checkbox-horizontal-align" int property can take the following values:

| Value | Constant     | Name     | Checkbox location                           |
|:-----:|--------------|----------|---------------------------------------------|
| 0	    | LeftAlign    | "left"   | At the left edge. Content on the right      |
| 1     | RightAlign   | "right"  | At the right edge. Content on the left      |
| 2     | CenterAlign  | "center" | Center horizontally. Content below or above |

The "checkbox-vertical-align" int property can take the following values:

| Value | Constant     | Name     | Checkbox location |
|:-----:|--------------|----------|-------------------|
| 0	    | TopAlign     | "top"    | Top alignment     |
| 1     | BottomAlign  | "bottom" | Bottom alignment  |
| 2     | CenterAlign  | "center" | Center alignment  |

Special case where both "checkbox-horizontal-align" and "checkbox-vertical-align" are CenterAlign (2).
In this case, the checkbox is centered horizontally, the content is below

You can get property values for "checkbox-horizontal-align" and "checkbox-vertical-align" using the functions

	func GetListViewCheckboxHorizontalAlign(view View, subviewID string) int
	func GetListViewCheckboxVerticalAlign(view View, subviewID string) int

### ListView events

There are three specific events for ListView

* "list-item-clicked" (ListItemClickedEvent constant) event occurs when the user clicks on a list item. 
The main listener for this event has the following format: func(ListView, int). 
Where the second argument is the index of the element.

* "list-item-selected" (ListItemSelectedEvent constant) event occurs when the user selects a list item. 
The main listener for this event has the following format: func(ListView, int).
Where the second argument is the index of the element.

* "list-item-checked" (ListItemCheckedEvent constant) event occurs when the user checks/unchecks the checkbox of a list item. 
The main listener for this event has the following format: func(ListView, []int).
Where the second argument is an array of indexes of the tagged items.

You can get lists of listeners for these events using the functions:

	func GetListItemClickedListeners(view View, subviewID string) []func(ListView, int)
	func GetListItemSelectedListeners(view View, subviewID string) []func(ListView, int)
	func GetListItemCheckedListeners(view View, subviewID string) []func(ListView, []int)

## TableView

The TableView element implements a table. To create a TableView, the function is used:

	func NewTableView(session Session, params Params) TableView

### "content" property

The "content" property defines the content of the table. 
To describe the content, you need to implement the TableAdapter interface declared as

	type TableAdapter interface {
		RowCount() int
		ColumnCount() int
		Cell(row, column int) interface{}
	}

where RowCount() and ColumnCount() functions must return the number of rows and columns in the table;
Cell(row, column int) returns the contents of a table cell. The Cell() function can return elements of the following types:

* string
* rune
* float32, float64
* integer values: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
* bool
* rui.Color
* rui.View
* fmt.Stringer
* rui.VerticalTableJoin, rui.HorizontalTableJoin

The "content" property can also be assigned the following data types

* TableAdapter
* [][]interface{}
* [][]string

[][]interface{} and [][]string are converted to a TableAdapter when assigned.

### "cell-style" property

The "cell-style" property (CellStyle constant) is used to customize the appearance of a table cell. 
Only an implementation of the TableCellStyle interface can be assigned to this property.

	type TableCellStyle interface {
		CellStyle(row, column int) Params
	}

This interface contains only one CellStyle function that returns the styling parameters for a given table cell. 
Any properties of the View interface can be used. For example

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		if row == 0 {
			return rui.Params {
				rui.BackgroundColor: rui.Gray,
				rui.Italic:          true,
			}
		}
		return nil
	}

If you don't need to change the appearance of a cell, you can return nil for it.

#### "row-span" and "column-span" properties

In addition to the properties of the View interface, the CellStyle function can return two more properties of type int:
"row-span" (RowSpan constant) and "column-span" (ColumnSpan constant).
These properties are used to combine table cells.

The "row-span" property specifies how many cells to merge vertically, and the "column-span" property - horizontally. For example

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		if row == 0 && column == 0 {
			return rui.Params { rui.RowSpan: 2 }
		}
		if row == 0 && column == 1 {
			return rui.Params { rui.ColumnSpan: 2 }
		}
		return nil
	}

In this case, the table will look like this

	|------+----------------|
	|      |                |
	|      +-------+--------|
	|      |       |        |
	|------+-------+--------|

If [][]interface{} is used as the value of the "content" property, then empty structures are used to merge cells

	type VerticalTableJoin struct {
	}

	type HorizontalTableJoin struct {
	}

These structures attach the cell to the top/left, respectively. The description of the above table will be as follows

	content := [][]interface{} {
		{"", "", rui.HorizontalTableJoin{}},
		{rui.VerticalTableJoin{}, "", ""},
	}

### "row-style" property

The "row-style" property (RowStyle constant) is used to customize the appearance of a table row.
This property can be assigned either an implementation of the TableRowStyle interface or []Params.
TableRowStyle is declared as

	type TableRowStyle interface {
		RowStyle(row int) Params
	}

The RowStyle function returns parameters that apply to the entire row of the table.
The "row-style" property has a lower priority than the "cell-style" property, i.e. 
properties set in "cell-style" will be used instead of those set in "row-style"

### "column-style" property

The "column-style" property (ColumnStyle constant) is used to customize the appearance of a table column.
This property can be assigned either an implementation of the TableColumnStyle interface or []Params.
TableColumnStyle is declared as

	type TableColumnStyle interface {
		ColumnStyle(column int) Params
	}

The ColumnStyle function returns the parameters applied to the entire column of the table.
The "column-style" property has a lower precedence over the "cell-style" and "row-style" properties.

### "head-height" and "head-style" properties

The table can have a header.
The "head-height" int property (constant HeadHeight) indicates how many first rows of the table form the header.
The "head-style" property (constant HeadStyle) sets the style of the heading. The "head-style" property can be
assigned, value of type:

* string - style name;
* []Params - enumeration of header properties.

### "foot-height" and "foot-style" properties

The table can have finalizing lines at the end (footer). For example, the "total" line.
The "foot-height" int property (the FootHeight constant) indicates the number of these footer lines.
The "foot-style" property (constant FootStyle) sets footer style. 
The values for the "foot-style" property are the same as for the "head-style" property.

### "cell-padding" property

The "cell-padding" BoundsProperty property (CellPadding constant) sets the padding from the cell borders to the content. 
This property is equivalent to

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		return rui.Params { rui.Padding: <my padding> }
	}

And it was introduced for convenience, so that you do not have to write an adapter to set indents.
The cell-padding property has a lower priority than the "cell-style" property.

"cell-padding" can also be used when setting parameters in the "row-style", "column-style", "foot-style", and "head-style" properties

### "cell-border" property

The "cell-border" property (CellBorder constant) sets the memory for all table cells.
This property is equivalent to

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		return rui.Params { rui.Border: <my padding> }
	}

And it was introduced for convenience, so that it is not necessary to write an adapter for the frame.
The "cell-border" property has a lower precedence over the "cell-style" property.

"cell-border" can also be used when setting parameters in properties
"row-style", "column-style", "foot-style" and "head-style"

### "table-vertical-align" property

The "table-vertical-align" int property (TableVerticalAlign constant) specifies 
the vertical alignment of data within a table cell. Valid values:

| Value | Constant      | Name       | Alignment          |
|:-----:|---------------|------------|--------------------|
| 0	    | TopAlign      | "top"      | Top alignment      |
| 1     | BottomAlign   | "bottom"   | Bottom alignment   |
| 2     | CenterAlign   | "center"   | Center alignment   |
| 3, 4  | BaselineAlign | "baseline" | Baseline alignment |

For horizontal alignment, use the "text-align" property

You can get the value of this property using the function

	func GetTableVerticalAlign(view View, subviewID string) int

### "selection-mode" property

The "selection-mode" property (SelectionMode constant) of the int type determines the mode of selection (highlighting) of table elements. Available modes:

* NoneSelection (0). Default mode. In this mode, you cannot select table elements. The table cannot receive input focus.

* CellSelection (1). In this mode, one table cell can be selected (highlighted).
The cell is selected interactively using the mouse or keyboard (using the cursor keys).
In this mode, the table can receive input focus. In this mode, the table generates two types of events: "table-cell-selected" and "table-cell-clicked" (see below).

* RowSelection (2). In this mode, only the entire table row can be selected (highlighted).
In this mode, the table is similar to a ListView. The row is selected interactively 
with the mouse or keyboard (using the cursor keys). In this mode, the table can receive input focus.
In this mode, the table generates two types of events: "table-row-selected" and "table-row-clicked" (see below).

You can get the value of this property using the function

	func GetSelectionMode(view View, subviewID string) int

### "current" property

The "current" property (Current constant) sets the coordinates of the selected cell/row as a structure

	type CellIndex struct {
		Row, Column int
	}

If the cell is not selected, then the values of the Row and Column fields will be less than 0.

In RowSelection mode, the value of the Column field is ignored. Also in this mode, 
the "current" property can be assigned a value of type int (row index).

You can get the value of this property using the function

	func GetTableCurrent(view View, subviewID string) CellIndex

### "allow-selection" property

By default, you can select any cell/row of the table. However, it is often necessary to disable the selection of certain elements. The "selection-mode" property (SelectionMode constant) allows you to set such a rule.

In CellSelection mode, this property is assigned the implementation of the interface

	type TableAllowCellSelection interface {
		AllowCellSelection(row, column int) bool
	}

and in RowSelection mode this property is assigned the implementation of the interface

	type TableAllowRowSelection interface {
		AllowRowSelection(row int) bool
	}

The AllowCellSelection/AllowRowSelection function must return "true" 
if the cell/row can be selected and "false" if the cell/row cannot be selected.

### "table-cell-selected" and "table-cell-clicked" events

The "table-cell-selected" event is fired in CellSelection mode when the user has selected 
a table cell with the mouse or keyboard.

The "table-cell-clicked" event occurs if the user clicks on a table cell (and if it is not selected, 
the "table-cell-selected" event occurs first) or presses the Enter or Space key.

The main listener for these events has the following format:

	func(TableView, int, int)

where the second argument is the cell row index, the third argument is the column index

You can also use a listener in the following format:

	func(int, int)

### "table-row-selected" and "table-row-clicked" events

The "table-row-selected" event is fired in RowSelection mode when the user has selected 
a table row with the mouse or keyboard.

The "table-row-clicked" event occurs if the user clicks on a table row (if it is not selected, 
the "table-row-selected" event fires first) or presses the Enter or Space key.

The main listener for these events has the following format:

	func(TableView, int)

where the second argument is the row index.

You can also use a listener in the following format:

	func(int)

## Custom View

A custom View must implement the CustomView interface, which extends the ViewsContainer and View interfaces. 
A custom View is created based on another, which is named Super View.

To simplify the task, there is already a basic CustomView implementation in the form of a CustomViewData structure.

Let's consider creating a custom View using the built-in Buttom element as an example:

1) declare the Button interface as extending CustomView, and the buttonData structure as extending CustomViewData

	type Button interface {
		rui.CustomView
	}

	type buttonData struct {
		rui.CustomViewData
	}

2) implement the CreateSuperView function

	func (button *buttonData) CreateSuperView(session Session) View {
		return rui.NewListLayout(session, rui.Params{
			rui.Semantics:       rui.ButtonSemantics,
			rui.Style:           "ruiButton",
			rui.StyleDisabled:   "ruiDisabledButton",
			rui.HorizontalAlign: rui.CenterAlign,
			rui.VerticalAlign:   rui.CenterAlign,
			rui.Orientation:     rui.StartToEndOrientation,
		})
	}

3) if necessary, override the methods of the CustomView interface, for Button this is the Focusable() function 
(since the button can receive focus, but ListLayout does not)

	func (button *buttonData) Focusable() bool {
		return true
	}

4) write a function to create a Button:

	func NewButton(session rui.Session, params rui.Params) Button {
		button := new(buttonData)
		rui.InitCustomView(button, "Button", session, params)
		return button
	}

When creating a CustomView, it is mandatory to call the InitCustomView function.
This function initializes the CustomViewData structure. 
The first argument is a pointer to the structure to be initialized, 
the second is the name assigned to your View, the third is the session and the fourth is the parameters

5) registering the item. It is recommended to register in the init method of the package

	rui.RegisterViewCreator("Button", func(session rui.Session) rui.View {
		return NewButton(session, nil)
	})

All! The new element is ready

## CanvasView

CanvasView is an area in which you can draw. To create a CanvasView, the function is used:

	func NewCanvasView(session Session, params Params) CanvasView

CanvasView has only one additional property: "draw-function" (DrawFunction constant).
Using this property, a drawing function is set with the following description

	func(Canvas)

where Canvas is the drawing context with which to draw

The Canvas interface contains a number of functions for customizing styles, text and drawing itself.

All coordinates and sizes are set only in pixels, so SizeUnit is not used when drawing.
float64 used everywhere

### Setting the line style

The following functions of the Canvas interface are used to customize the line color:

* SetSolidColorStrokeStyle(color Color) - the line will be drawn with a solid color

* SetLinearGradientStrokeStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) - 
the line will be drawn using a linear gradient. 
The gradient starts at x0, y0, and color0, and the gradient ends at x1, y1, and color1. 
The []GradientPoint array specifies the intermediate points of the gradient. 
If there are no intermediate points, then nil can be passed as the last parameter

* SetRadialGradientStrokeStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) - 
the line will be drawn using a radial gradient. 
x0, y0, r0, color0 - center coordinates, radius and color of the starting circle.
x1, y1, r1, color1 - center coordinates, radius and color of the end circle. 
The []GradientPoint array specifies intermediate points of the gradient

The GradientPoint structure is described as

	type GradientPoint struct {
		Offset float64
		Color Color
	}

where Offset is a value in the range from 0 to 1 specifies the relative position of the intermediate point, Color is the color of this point.

Line width in pixels is set by the function

	SetLineWidth(width float64)

The type of line ends is set using the function

	SetLineCap(cap int)

where cap can take the following values

| Value | Constant  | View |
|:-----:|-----------|------------------------------------------------------------------------------|
| 0     | ButtCap   | The ends of lines are squared off at the endpoints. Default value.           |
| 1     | RoundCap  | The ends of lines are rounded. The center of the circle is at the end point. |
| 2     | SquareCap | the ends of lines are squared off by adding a box with an equal width and half the height of the line's thickness. |

The shape used to connect two line segments at their intersection is specified by the function

	SetLineJoin(join int)

where join can take the following values

| Value | Constant  | View                                   |
|:-----:|-----------|----------------------------------------|
| 0     | MiterJoin | Connected segments are joined by extending their outside edges to connect at a single point, with the effect of filling an additional lozenge-shaped area. This setting is affected by the miterLimit property |
| 1     | RoundJoin | rounds off the corners of a shape by filling an additional sector of disc centered at the common endpoint of connected segments. The radius for these rounded corners is equal to the line width. |
| 2     | BevelJoin | Fills an additional triangular area between the common endpoint of connected segments, and the separate outside rectangular corners of each segment. |

By default, a solid line is drawn. If you want to draw a broken line, you must first set the pattern using the function

	SetLineDash(dash []float64, offset float64)

where dash []float64 specifies the line pattern in the form of alternating line lengths and gaps. 
The second argument is the offset of the template relative to the beginning of the line.

Example

	canvas.SetLineDash([]float64{16, 8, 4, 8}, 0)

The line is drawn as follows: a 16-pixel segment, then an 8-pixel gap, then a 4-pixel segment, then an 8-pixel gap, then a 16-pixel segment again, and so on.

### Setting the fill style

The following functions of the Canvas interface are used to customize the fill style:

* SetSolidColorFillStyle(color Color) - the shape will be filled with a solid color

* SetLinearGradientFillStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) - 
the shape will be filled with a linear gradient. 
The gradient starts at x0, y0, and color0, and the gradient ends at x1, y1, and color1. 
The []GradientPoint array specifies the intermediate points of the gradient. 
If there are no intermediate points, then nil can be passed as the last parameter

* SetRadialGradientFillStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) - 
the shape will be filled with a radial gradient. 
x0, y0, r0, color0 - center coordinates, radius and color of the starting circle.
x1, y1, r1, color1 - center coordinates, radius and color of the end circle. 
Array []GradientPoint specifies intermediate points of the gradient

### Drawing geometric shapes

#### Rectangle

Three functions can be used to draw rectangles:

	FillRect(x, y, width, height float64)
	StrokeRect(x, y, width, height float64)
	FillAndStrokeRect(x, y, width, height float64)

FillRect draws a filled rectangle.

StrokeRect draws the outline of a rectangle.

FillAndStrokeRect draws a path and fills in the interior.

#### Rounded Rectangle

Similar to the rectangle, there are three drawing functions

	FillRoundedRect(x, y, width, height, r float64)
	StrokeRoundedRect(x, y, width, height, r float64)
	FillAndStrokeRoundedRect(x, y, width, height, r float64)

where r is the radius of the rounding

#### Ellipse

Three functions can also be used to draw ellipses:

	FillEllipse(x, y, radiusX, radiusY, rotation float64)
	StrokeEllipse(x, y, radiusX, radiusY, rotation float64)
	FillAndStrokeEllipse(x, y, radiusX, radiusY, rotation float64)

where x, y is the center of the ellipse, radiusX, radiusY are the radii of the ellipse along the X and Y axes,
rotation - the angle of rotation of the ellipse relative to the center in radians.

#### Path

The Path interface allows you to describe a complex shape. Path is created using the NewPath () function.

Once created, you must describe the shape. For this, the following interface functions can be used:

* MoveTo(x, y float64) - move the current point to the specified coordinates;

* LineTo(x, y float64) - add a line from the current point to the specified one;

* ArcTo(x0, y0, x1, y1, radius float64) - add a circular arc using the specified control points and radius.
If necessary, the arc is automatically connected to the last point of the path with a straight line.
x0, y0 - coordinates of the first control point;
x1, y1 - coordinates of the second control point;
radius - radius of the arc. Must be non-negative.

* Arc(x, y, radius, startAngle, endAngle float64, clockwise bool) - add a circular arc.
x, y - coordinates of the arc center;
radius - radius of the arc. Must be non-negative;
startAngle - The angle, in radians, at which the arc begins, measured clockwise from the positive X-axis.
endAngle - The angle, in radians, at which the arc ends, measured clockwise from the positive X-axis.
clockwise - if true, the arc will be drawn clockwise between the start and end corners, otherwise counterclockwise

* BezierCurveTo(cp0x, cp0y, cp1x, cp1y, x, y float64) - add a cubic Bezier curve from the current point.
cp0x, cp0y - coordinates of the first control point;
cp1x, cp1y - coordinates of the second control point;
x, y - coordinates of the end point.

* QuadraticCurveTo(cpx, cpy, x, y float64) - add a quadratic Bezier curve from the current point.
cpx, cpy - coordinates of the control point;
x, y - coordinates of the end point.

* Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, clockwise bool) - add an elliptical arc.
x, y - coordinates of the center of the ellipse;
radiusX is the radius of the major axis of the ellipse. Must be non-negative;
radiusY - radius of the minor axis of the ellipse. Must be non-negative;
rotation - the rotation of the ellipse, expressed in radians;
startAngle The angle of the start of the ellipse in radians, measured clockwise from the positive x-axis.
endAngle The angle, in radians, at which the ellipse ends, measured clockwise from the positive x-axis.
clockwise - if true, draws the ellipse clockwise, otherwise counterclockwise.

The Close () function is called at the end and connects the start and end points of the shape. Used only for closed shapes.

After the Path is formed, it can be drawn using the following 3 functions

	FillPath(path Path)
	StrokePath(path Path)
	FillAndStrokePath(path Path)

#### Line

To draw a line, use the function

	DrawLine(x0, y0, x1, y1 float64)

### Text

To display text in specified coordinates, two functions are used

	FillText(x, y float64, text string)
	StrokeText(x, y float64, text string)

The StrokeText function draws the outline of the text, FillText draws the text itself.

The horizontal alignment of the text relative to the specified coordinates is set using the function

	SetTextAlign(align int)

where align can be one of the following values:

| Value | Constant    | Alignment |
|:-----:|-------------|--------------------------------------------------------|
| 0     | LeftAlign   | The specified point is the leftmost point of the text  |
| 1     | RightAlign  | The specified point is the rightmost point of the text |
| 2     | CenterAlign | The text is centered on the specified point            |
| 3     | StartAlign  | If the text is displayed from left to right, then the text output is equivalent to LeftAlign, otherwise RightAlign |
| 4     | EndAlign    | If the text is displayed from left to right, then the text output is equivalent to RightAlign, otherwise LeftAlign |
	
The vertical alignment of the text relative to the specified coordinates is set using the function

	SetTextBaseline(baseline int)

where baseline can be one of the following values:

| Value | Constant            | Alignment |
|:-----:|---------------------|--------------------------------------------------|
| 0     | AlphabeticBaseline  | Relatively normal baseline of text               |
| 1     | TopBaseline         | Relative to the top border of the text           |
| 2     | MiddleBaseline      | About the middle of the text                     |
| 3     | BottomBaseline      | To the bottom of the text                        |
| 4     | HangingBaseline     | Relative to the dangling baseline of the text (used in Tibetan and other Indian scripts) |
| 5     | IdeographicBaseline | Relative to the ideographic baseline of the text |

An ideographic baseline is the bottom of a character display if the main character is 
outside the alphabet baseline (Used in Chinese, Japanese, and Korean fonts).

To set the font parameters of the displayed text, use the functions

	SetFont(name string, size SizeUnit)
	SetFontWithParams(name string, size SizeUnit, params FontParams)

where FontParams is defined as

	type FontParams struct {
		// Italic - if true then a font is italic
		Italic bool
		// SmallCaps - if true then a font uses small-caps glyphs
		SmallCaps bool
		// Weight - a font weight. Valid values: 0…9, there
		//   0 - a weight does not specify;
		//   1 - a minimal weight;
		//   4 - a normal weight;
		//   7 - a bold weight;
		//   9 - a maximal weight.
		Weight int
		// LineHeight - the height (relative to the font size of the element itself) of a line box.
		LineHeight SizeUnit
	}

The TextWidth function allows you to find out the width of the displayed text in pixels

	TextWidth(text string, fontName string, fontSize SizeUnit) float64

### Image

Before drawing an image, it must first be loaded. The global function is used for this:

	func LoadImage(url string, onLoaded func(Image), session Session) Image {

The image is loaded asynchronously. After the download is finished, the function passed in the second argument will be called.
If the image was loaded successfully, then the LoadingStatus() function of the Image interface will 
return the value ImageReady (1), if an error occurred while loading, then this function will return ImageLoadingError (2).
The textual description of the error is returned by the LoadingError() function

Unlike an ImageView, loading an Image does not take into account the pixel density. 
It is up to you to decide which image to upload. You can do it like this:

	var url string
	if session.PixelRatio() == 2 {
		url = "image@2x.png"
	} else {
		url = "image.png"
	}

The following functions are used to draw the image:

	DrawImage(x, y float64, image Image)
	DrawImageInRect(x, y, width, height float64, image Image)
	DrawImageFragment(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight float64, image Image)

The DrawImage function displays the image as it is (without scaling): x, y - coordinates of the upper left corner of the image

The DrawImageInRect function displays the image with scaling: x, y are coordinates of the upper left corner of the image,
width, height are width and height of the result

The DrawImageFragment function displays a fragment of the image with scaling: srcX, srcY, srcWidth, srcHeight describe 
the original area of the image, dstX, dstY, dstWidth, dstHeight describe the resulting area.

Image can also be used in fill style

	SetImageFillStyle(image Image, repeat int)

where repeat can take on the following values:

| Value | Constant  | Description                                   |
|:-----:|-----------|-----------------------------------------------|
| 0     | NoRepeat  | Image is not repeated                         |
| 1     | RepeatXY  | Image is repeated vertically and horizontally |
| 2     | RepeatX   | The image is repeated horizontally only       |
| 3     | RepeatY   | The image is repeated vertically only         |

## AudioPlayer, VideoPlayer, MediaPlayer

AudioPlayer and VideoPlayer are elements for audio and video playback.
Both elements implement the MediaPlayer interface. Most of the properties and all events 
of AudioPlayer and VideoPlayer are implemented through the MediaPlayer.

### Свойство "src"

The "src" property (Source constant) specifies one or more media sources. The "src" property can take on the following types:

* string,
* MediaSource,
* []MediaSource.

The MediaSource structure is declared as

	type MediaSource struct {
		Url      string
		MimeType string
	}

where Url is a required parameter, MimeType is an optional mime file type

Since different browsers support different file formats and codecs, it is recommended 
to specify multiple sources in different formats. The player chooses the most suitable 
one from the list of sources. Setting mime types makes this process easier for the browser

### "controls" property

The "controls" bool property (Controls constant) specifies whether UI elements should be 
displayed to control playback of the media resource. The default is false.

If the "controls" property is false for the AudioPlayer, then it will be invisible and will not take up screen space.

### "loop" property

The "loop" bool property (Loop constant). If it set to true, then the media file will start over when it reaches the end. The default is false.

### "muted" property

The "muted" bool property (constant Muted) enables (true) / disables (false) silent mode. The default is false.

### "preload" property

The "preload" int property (constant Preload) defines what data should be preloaded.
Valid values:

| Value | Constant        | Name       | Description                                                                      |
|:-----:|-----------------|------------|----------------------------------------------------------------------------------|
| 0	    | PreloadNone     | "none"     | Media file must not be pre-loaded                                                |
| 1	    | PreloadMetadata | "metadata" | Only metadata is preloaded                                                       |
| 2	    | PreloadAuto     | "auto"     | The entire media file can be downloaded even if the user doesn't have to use it. |

The default value is PreloadAuto (2)

### "poster" property

The "poster" string property (Poster constant) is used only for VideoPlayer.
It sets the url of the image that will be shown until the video is loaded.
If this property is not set, then a black screen will be shown first, and then the first frame (as soon as it is loaded).

### "video-width" and "video-height" properties

The "video-width" (VideoWidth constant) and "video-height" (VideoHeight constant) float64 properties are used only for VideoPlayer.
It defines the width and height of the rendered video in pixels.

If "video-width" and "video-height" are not specified, then the actual dimensions of the video are used, 
while the dimensions of the container in which the video is placed are ignored and the video may overlap 
other interface elements. Therefore, it is recommended to set these values, for example, like this

	rui.Set(view, "videoPlayerContainer", rui.ResizeEvent, func(frame rui.Frame) {
		rui.Set(view, "videoPlayer", rui.VideoWidth, frame.Width)
		rui.Set(view, "videoPlayer", rui.VideoHeight, frame.Height)
	})

If only one of the "video-width" or "video-height" properties is set, then the second is calculated based on the aspect ratio of the video

### Developments

MediaPlayer has two groups of events:

1) has a handler like 

    func(MediaPlayer) 
    
You can also use func(). This group includes the following events:

* "abort-event" (constant AbortEvent) - Fires when the resource is not fully loaded, but not as a result of an error.

* "can-play-event" (CanPlayEvent constant) is fired when the user agent is able to play media 
but judges that there is not enough data loaded to play the media to the end without 
having to stop to further buffer the content.

* "can-play-through-event" (constant CanPlayThroughEvent) is fired when the user agent is able 
to play the media and evaluates that enough data has been loaded to play the media to its end, 
without having to stop to further buffer the content.

* "complete-event" (CompleteEvent constant) -

* "emptied-event" (EmptiedEvent constant) is fired when the media becomes empty; for example when the media is already loaded (or partially loaded)

* "ended-event" (EndedEvent constant) - Fires when playback stops, when the end of media is reached, or if no further data is available.

* "loaded-data-event" (LoadedDataEvent constant) is fired when the first frame of the media has finished loading.

* "loaded-metadata-event" (LoadedMetadataEvent constant) is fired when the metadata has been loaded.

* "loadstart-event" (LoadstartEvent constant) is fired when the browser starts loading the resource.

* "pause-event" (PauseEvent constant) is fired when a pause request is processed and the action pauses, most often when the Pause () method is called.

* "play-event" (PlayEvent constant) is fired when the media file starts playing, for example, as a result of using the Play () method

* "playing-event" (PlayingEvent constant) is fired when playback is about to start after being paused or delayed due to lack of data.

* "progress-event" (ProgressEvent constant) is fired periodically when the browser loads the resource.

* "seeked-event" (SeekedEvent constant) is fired when the playback speed has changed.

* "seeking-event" (SeekingEvent constant) is fired when a seeking operation begins.

* "stalled-event" (StalledEvent constant) is fired when the user agent tries to retrieve media data, but no data arrives unexpectedly.

* "suspend-event" (SuspendEvent constant) is fired when media loading has been suspended.

* "waiting-event" (WaitingEvent constant) is fired when playback is stopped due to a temporary lack of data

2) has a handler like 

    func(MediaPlayer, float64) 
    
You can also use func(float64), func(MediaPlayer) and func().
This group includes events related to changing the parameters of the player. 
The new value of the changed parameter is passed as the second argument.

* "duration-changed-event" (DurationChangedEvent constant) is fired when the duration attribute has been updated.

* "time-updated-event" (TimeUpdatedEvent constant) is fired when the current time has been updated.

* "volume-changed-event" (VolumeChangedEvent constant) is fired when the volume changes.

* "rate-changed-event" (RateChangedEvent constant) is fired when the playback speed has changed.

A separate event that does not belong to these two groups, "player-error-event" (PlayerErrorEvent constant) 
is fired when the resource cannot be loaded due to an error (eg network error).

The handler for this event looks like 

    func(player MediaPlayer, code int, message string) 

You can also use func(int, string), func(MediaPlayer) and func(). 
Where the argument "message" is the error message, "code" is the error code:

| Error code | Constant                      | Value                                                                     |
|:----------:|-------------------------------|---------------------------------------------------------------------------|
| 0	         | PlayerErrorUnknown            | Unknown error                                                             |
| 1	         | PlayerErrorAborted            | Fetching the associated resource was interrupted by a user request.       |
| 2	         | PlayerErrorNetwork            | Some kind of network error has occurred that prevented the media from successfully ejecting, even though it was previously available. |
| 3	         | PlayerErrorDecode             | Although the resource was previously identified as being used, an error occurred while trying to decode the media resource. |
| 4	         | PlayerErrorSourceNotSupported | The associated resource object or media provider was found to be invalid. |

### Methods

MediaPlayer has a number of methods for controlling player parameters:

* Play() starts playback of a media file;

* Pause() pauses playback;

* SetCurrentTime(seconds float64) sets the current playback time in seconds;

* CurrentTime() float64 returns the current playing time in seconds;

* Duration() float64 returns the duration of the media file in seconds;

* SetPlaybackRate(rate float64) sets the playback speed. Normal speed is 1.0;

* PlaybackRate() float64 returns the current playback speed;

* SetVolume(volume float64) sets the volume speed in the range from 0 (silence) to 1 (maximum volume);

* Volume() float64 returns the current volume;

* IsEnded() bool returns true if the end of the media file is reached;

* IsPaused() bool returns true if playback is paused.

For quick access to these methods, there are global functions:

	func MediaPlayerPlay(view View, playerID string)
	func MediaPlayerPause(view View, playerID string)
	func SetMediaPlayerCurrentTime(view View, playerID string, seconds float64)
	func MediaPlayerCurrentTime(view View, playerID string) float64
	func MediaPlayerDuration(view View, playerID string) float64
	func SetMediaPlayerVolume(view View, playerID string, volume float64)
	func MediaPlayerVolume(view View, playerID string) float64
	func SetMediaPlayerPlaybackRate(view View, playerID string, rate float64)
	func MediaPlayerPlaybackRate(view View, playerID string) float64
	func IsMediaPlayerEnded(view View, playerID string) bool
	func IsMediaPlayerPaused(view View, playerID string) bool

where view is the root View, playerID is the id of AudioPlayer or VideoPlayer

## Animation

The library supports two types of animation:

* Animated property value changes (hereinafter "transition animation")
* Script animated change of one or more properties (hereinafter simply "animation script")

### Animation interface

The Animation interface is used to set animation parameters. It extends the Properties interface.
The interface is created using the function:

	func NewAnimation(params Params) Animation

Some of the properties of the Animation interface are used in both types of animation, the rest are used only 
in animation scripts.

Common properties are

| Property          | Constant       | Type    | Default      | Description |
|-------------------|----------------|---------|--------------|----------------------------------------------|
| "duration"        | Duration       | float64 | 1            | Animation duration in seconds                |
| "delay"           | Delay          | float64 | 0            | Delay before animation in seconds            |
| "timing-function" | TimingFunction | string  | "ease"       | The function of changing the animation speed |

Properties used only in animation scripts will be described below.

#### "timing-function" property

The "timing-function" property describes in text the function of changing the speed of the animation.
Functions can be divided into 2 types: simple functions and functions with parameters.

Simple functions

| Function      | Constant        | Description                                                       |
|---------------|-----------------|-------------------------------------------------------------------|
| "ease"        | EaseTiming      | the speed increases towards the middle and slows down at the end. |
| "ease-in"     | EaseInTiming    | the speed is slow at first, but increases in the end.             |
| "ease-out"    | EaseOutTiming   | speed is fast at first, but decreases rapidly. Most of the slow   |
| "ease-in-out" | EaseInOutTiming | the speed is fast at first, but quickly decreases, and at the end it increases again. |
| "linear"      | LinearTiming    | constant speed                                                    |

And there are two functions with parameters:

* "steps(N)" - discrete function, where N is an integer specifying the number of steps. 
You can specify this function either as text or using the function:

	func StepsTiming(stepCount int) string

For example

	animation := rui.NewAnimation(rui.Params{
		rui.TimingFunction: rui.StepsTiming(10),
	})

equivalent to 

	animation := rui.NewAnimation(rui.Params{
		rui.TimingFunction: "steps(10)",
	})
	
* "cubic-bezier(x1, y1, x2, y2)" - time function of a cubic Bezier curve. x1, y1, x2, y2 are of type float64.
x1 and x2 must be in the range [0...1]. You can specify this function either as text or using the function:

	func CubicBezierTiming(x1, y1, x2, y2 float64) string

### Transition animation

Transition animation can be applied to properties of the type: SizeUnit, Color, AngleUnit, float64 and composite properties that contain elements of the listed types (for example, "shadow", "border", etc.).

If you try to animate other types of properties (for example, bool, string), no error will occur, 
there will simply be no animation.

There are two types of transition animations:
* single-fold;
* constant;

A one-time animation is triggered using the SetAnimated function of the View interface. 
This function has the following description:

	SetAnimated(tag string, value interface{}, animation Animation) bool

It assigns a new value to the property, and the change occurs using the specified animation.
For example,

	view.SetAnimated(rui.Width, rui.Px(400), rui.NewAnimation(rui.Params{
		rui.Duration:       0.75,
		rui.TimingFunction: rui.EaseOutTiming,
	}))

There is also a global function for animated one-time change of the property value of the child View

	func SetAnimated(rootView View, viewID, tag string, value interface{}, animation Animation) bool

A persistent animation runs every time the property value changes. 
To set the constant animation of the transition, use the "transition" property (the Transition constant). 
As a value, this property is assigned rui.Params, where the property name should be the key, 
and the value should be the Animation interface.
For example,

	view.Set(rui.Transition, rui.Params{
		rui.Height: rui.NewAnimation(rui.Params{
			rui.Duration:       0.75,
			rui.TimingFunction: rui.EaseOutTiming,
		},
		rui.BackgroundColor: rui.NewAnimation(rui.Params{
			rui.Duration:       1.5,
			rui.Delay:          0.5,
			rui.TimingFunction: rui.Linear,
		},
	})

Calling the SetAnimated function does not change the value of the "transition" property.

To get the current list of permanent transition animations, use the function

	func GetTransition(view View, subviewID string) Params

It is recommended to add new transition animations using the function 

	func AddTransition(view View, subviewID, tag string, animation Animation) bool

Calling this function is equivalent to the following code

	transitions := rui.GetTransition(view, subviewID)
	transitions[tag] = animation
	rui.Set(view, subviewID, rui.Transition, transitions)

#### Transition animation events

The transition animation generates the following events

| Event                     | Constant              | Description                                                      |
|---------------------------|-----------------------|------------------------------------------------------------------|
| "transition-run-event"    | TransitionRunEvent    | The transition animation loop has started, i.e. before the delay |
| "transition-start-event"  | TransitionStartEvent  | The transition animation has actually started, i.e. after delay  |
| "transition-end-event"    | TransitionEndEvent    | Transition animation finished                                    |
| "transition-cancel-event" | TransitionCancelEvent | Transition animation interrupted                                 |

The main event listener has the following format:

	func(View, string)

where the second argument is the name of the property.

You can also use a listener in the following format:

	func()
	func(string)
	func(View)

Get lists of listeners for transition animation events using functions:

	func GetTransitionRunListeners(view View, subviewID string) []func(View, string)
	func GetTransitionStartListeners(view View, subviewID string) []func(View, string)
	func GetTransitionEndListeners(view View, subviewID string) []func(View, string)
	func GetTransitionCancelListeners(view View, subviewID string) []func(View, string)

### Animation script

An animation script describes a more complex animation than a transition animation. To do this, additional properties are added to Animation:

#### "property" property

The "property" property (constant PropertyTag) describes property changes. 
[]AnimatedProperty or AnimatedProperty is assigned as a value. The AnimatedProperty structure describes 
the change script of one property. She is described as

	type AnimatedProperty struct {
		Tag       string
		From, To  interface{}
		KeyFrames map[int]interface{}
	}

where Tag is the name of the property, From is the initial value of the property, 
To is the final value of the property, KeyFrames is intermediate property values (keyframes).

The required fields are Tag, From, To. The KeyFrames field is optional, it can be nil.

The KeyFrames field describes key frames. As a key of type int, the percentage of time elapsed 
since the beginning of the animation is used (exactly the beginning of the animation itself, 
the time specified by the "delay" property is excluded).
And the value is the value of the property at a given moment in time. For example

	prop := rui.AnimatedProperty {
		Tag:       rui.Width,
		From:      rui.Px(100),
		To:        rui.Px(200),
		KeyFrames: map[int]interface{
			90: rui.Px(220),
		}
	}

In this example, the "width" property will grow from 100px to 220px 90% of the time. 
In the remaining 10% of the time, it will decrease from 220px to 200px.

The "property" property is assigned to []AnimatedProperty, which means that you can animate several properties at once.

You must set at least one "property" element, otherwise the animation will be ignored.

#### "id" property

The "id" string property (constant ID) specifies the animation identifier. 
Passed as a parameter to the animation event listener. If you do not plan to use event listeners for animation, 
then you do not need to set this property.

#### "iteration-count" property

The "iteration-count" int property (constant IterationCount) specifies the number of animation repetitions.
The default is 1. A value less than zero causes the animation to repeat indefinitely.

#### "animation-direction" property

The "animation-direction" int property (an AnimationDirection constant) specifies whether 
the animation should play forward, backward, or alternately forward and backward between forward and 
backward playback of the sequence. It can take the following values:

| Value    | Constant                  | Description                                                             |
|:--------:|---------------------------|-----------------------------------------------------------------------|
| 0        | NormalAnimation           | The animation plays forward every iteration, that is, when the animation ends, it is immediately reset to its starting position and played again. |
| 1        | ReverseAnimation          | The animation plays backwards, from the last position to the first, and then resets to the final position and plays again. |
| 2        | AlternateAnimation        | The animation changes direction in each cycle, that is, in the first cycle, it starts from the start position, reaches the end position, and in the second cycle, it continues from the end position and reaches the start position, and so on. |
| 3        | AlternateReverseAnimation | The animation starts playing from the end position and reaches the start position, and in the next cycle, continuing from the start position, it goes to the end position. |

#### Animation start

To start the animation script, you must assign the interface created by Animation to the "animation" property 
(the AnimationTag constant). If the View is already displayed on the screen, then the animation starts immediately
(taking into account the specified delay), otherwise the animation starts as soon as the View is displayed 
on the screen.

The "animation" property can be assigned Animation and [] Animation, ie. you can run several animations 
at the same time for one View

Example,

	prop := rui.AnimatedProperty {
		Tag:       rui.Width,
		From:      rui.Px(100),
		To:        rui.Px(200),
		KeyFrames: map[int]interface{
			90: rui.Px(220),
		}
	}
	animation := rui.NewAnimation(rui.Params{
		rui.PropertyTag:    []rui.AnimatedProperty{prop},
		rui.Duration:       2,
		rui.TimingFunction: LinearTiming,
	})
	rui.Set(view, "subview", rui.AnimationTag, animation)

#### "animation-paused" property

The "animation-paused" bool property of View (AnimationPaused constant) allows the animation to be paused.
In order to pause the animation, set this property to "true", and to resume to "false".

Attention. When you assign a value to the "animation" property, the "animation-paused" property is set to false.
 
#### Animation events

The animation script generates the following events

| Event                       | Constant                | Description                            |
|-----------------------------|-------------------------|----------------------------------------|
| "animation-start-event"     | AnimationStartEvent     | Animation started                      |
| "animation-end-event"       | AnimationEndEvent       | Animation finished                     |
| "animation-cancel-event"    | AnimationCancelEvent    | Animation interrupted                  |
| "animation-iteration-event" | AnimationIterationEvent | A new iteration of animation has begun |

Attention! Not all browsers support the "animation-cancel-event" event. This is currently only Safari and Firefox.

The main event data listener has the following format:

	func(View, string)

where the second argument is the id of the animation.

You can also use a listener in the following format:

	func()
	func(string)
	func(View)

Get lists of animation event listeners using functions:

	func GetAnimationStartListeners(view View, subviewID string) []func(View, string)
	func GetAnimationEndListeners(view View, subviewID string) []func(View, string)
	func GetAnimationCancelListeners(view View, subviewID string) []func(View, string)
	func GetAnimationIterationListeners(view View, subviewID string) []func(View, string)

## Session

When a client creates a connection to a server, a Session interface is created for that connection.
This interface is used to interact with the client.
You can get the current Session interface by calling the Session() method of the View interface.

When a session is created, it gets a custom implementation of the SessionContent interface.

	type SessionContent interface {
		CreateRootView(session rui.Session) rui.View
	}

This interface is created by the function passed as a parameter when creating an application by the NewApplication function.

In addition to the mandatory CreateRootView() function, SessionContent can have several optional functions:

	OnStart(session rui.Session)
	OnFinish(session rui.Session)
	OnResume(session rui.Session)
	OnPause(session rui.Session)
	OnDisconnect(session rui.Session)
	OnReconnect(session rui.Session)

Immediately after creating a session, the CreateRootView function is called. After creating the root View, the OnStart function is called (if implemented)

The OnFinish function (if implemented) is called when the user closes the application page in the browser

The OnPause function is called when the application page in the client's browser becomes inactive.
This happens if the user switches to a different browser tab / window, minimizes the browser, or switches to another application.

The OnResume function is called when the application page in the client's browser becomes active. Also, this function is called immediately after OnStart

The OnDisconnect function is called if the server loses connection with the client. This happens either when the connection is broken.

The OnReconnect function is called after the server reconnects with the client.

The Session interface provides the following methods:

* DarkTheme() bool returns true if a dark theme is used. Determined by client-side settings

* TouchScreen() bool  returns true if client supports touch screen

* PixelRatio() float64  returns the size of a logical pixel, i.e. how many physical pixels form a logical. For example, for iPhone, this value will be 2 or 3

* TextDirection() int returns the direction of the letter: LeftToRightDirection (1) or RightToLeftDirection (2)

* Constant(tag string) (string, bool) returns the value of a constant

* Color(tag string) (Color, bool) returns the value of the color constant

* SetCustomTheme(name string) bool sets the theme with the given name as the current one. 
Returns false if no topic with this name was found. Themes named "" are the default theme.

* Language() string returns the current interface language, for example: "en", "ru", "ptBr"

* SetLanguage(lang string) sets the current interface language (see "Support for multiple languages")

* GetString(tag string) (string, bool) returns a textual text value for the current language
(see "Support for multiple languages")

* Content() SessionContent returns the current SessionContent instance

* RootView() View returns the root View of the session

* SetTitle(title string) sets the text of the browser title/tab

* SetTitleColor(color Color) sets the color of the browser navigation bar. Supported only in Safari and Chrome for android

* Get(viewID, tag string) interface{} returns the value of the View property named tag. Equivalent to

	rui.Get(session.RootView(), viewID, tag)

* Set(viewID, tag string, value interface {}) bool sets the value of the View property named tag.

	rui.Set(session.RootView(), viewID, tag, value)

* DownloadFile(path string) downloads (saves) on the client side the file located at the specified path on the server.
It is used when the client needs to transfer a file from the server.

* DownloadFileData(filename string, data [] byte) downloads (saves) on the client side a file 
with a specified name and specified content. Typically used to transfer a file generated in server memory.

## Resource description format

Application resources (themes, views, translations) can be described as text (utf-8). 
This text is placed in a file with the ".rui" extension.

The root element of the resource file must be an object. It has the following format:

	< object name > {
		< object data >
	}

if the object name contains the following characters: '=', '{', '}', '[', ']', ',', '', '\t', '\n', 
'\' ',' "','` ',' / 'and any spaces, then the object name must be enclosed in quotation marks. 
If these characters are not used, then quotation marks are optional.

You can use three types of quotation marks:

* "…" is equivalent to the same string in the go language, i.e. inside you can use escape sequences:
\n, \r, \\, \", \', \0, \t, \x00, \u0000

* '…' is similar to the line "…"

* `…` is equivalent to the same string in the go language, i.e. the text within this line remains as is. Inside
you cannot use the ` character.

Object data is a set of < key > = < value > pairs separated by commas.

The key is a string of text. The design rules are the same as for the object name.

Values can be of 3 types:

* Simple value - a line of text formatted according to the same rules as the name of the object

* An object

* Array of values

An array of values is enclosed in square brackets. Array elements are separated by commas.
Elements can be simple values or objects.

There may be comments in the text. The design rules are the same as in the go language: // and / * ... * /

Example:

	GridLayout {
		id = gridLayout, width = 100%, height = 100%,
		cell-width = "150px, 1fr, 30%", cell-height = "25%, 200px, 1fr",
		content = [
			// Subviews
			TextView { row = 0, column = 0:1,
				text = "View 1", text-align = center, vertical-align = center,
				background-color = #DDFF0000, radius = 8px, padding = 32px,
				border = _{ style = solid, width = 1px, color = #FFA0A0A0 }
			},
			TextView { row = 0:1, column = 2,
				text = "View 2", text-align = center, vertical-align = center,
				background-color = #DD00FF00, radius = 8px, padding = 32px,
				border = _{ style = solid, width = 1px, color = #FFA0A0A0 }
			},
			TextView { row = 1:2, column = 0,
				text = "View 3", text-align = center, vertical-align = center,
				background-color = #DD0000FF, radius = 8px, padding = 32px,
				border = _{ style = solid, width = 1px, color = #FFA0A0A0 }
			},
			TextView { row = 1, column = 1,
				text = "View 4", text-align = center, vertical-align = center,
				background-color = #DDFF00FF, radius = 8px, padding = 32px,
				border = _{ style = solid, width = 1px, color = #FFA0A0A0 }
			},
			TextView { row = 2, column = 1:2,
				text = "View 5", text-align = center, vertical-align = center,
				background-color = #DD00FFFF, radius = 8px, padding = 32px,
				border = _{ style = solid, width = 1px, color = #FFA0A0A0 }
			},
		]
	}

To work with text resources, the DataNode interface is used

	type DataNode interface {
		Tag() string
		Type() int
		Text() string
		Object() DataObject
		ArraySize() int
		ArrayElement(index int) DataValue
		ArrayElements() []DataValue
	}

This element describes the underlying data element.

The Tag method returns the value of the key.

The data type is returned by the Type method. It returns one of 3 values

| Value | Constant   | Data type    |
|:-----:|------------|--------------|
| 0	    | TextNode   | Simple value |
| 1	    | ObjectNode | Object       |
| 2     | ArrayNode  | Array        |

The Text() method is used to get a simple value.
To get an object, use the Object() method.
To get the elements of an array, use the ArraySize, ArrayElement and ArrayElements methods

## Resources

Resources (pictures, themes, translations, etc.) with which the application works should be placed 
in subdirectories within one resource directory. Resources should be located in the following subdirectories:

* images - all images are placed in this subdirectory. Here you can make nested subdirectories.
In this case, they must be included in the file name. For example, "subdir/image1.png"

* themes - application themes are placed in this subdirectory (see below)

* views - View descriptions are placed in this subdirectory

* strings - translations of text resources are placed in this subdirectory (see Support for multiple languages)

* raw - all other resources are placed in this subdirectory: sounds, video, binary data, etc.

The resource directory can either be included in the executable file or located separately.

If the resources need to be included in the executable file, then the name of the directory must be "resources" and it must be connected as follows:

	import (
		"embed"

		"github.com/anoshenko/rui"
	)

	//go:embed resources
	var resources embed.FS

	func main() {
		rui.AddEmbedResources(&resources)
		
		app := rui.NewApplication("Hello world", createHelloWorldSession)
		app.Start("localhost:8000")
	}

If the resources are supplied as a separate directory, then it must be registered 
using the SetResourcePath function before creating the Application:

	func main() {
		rui.SetResourcePath(path)
		
		app := rui.NewApplication("Hello world", createHelloWorldSession)
		app.Start("localhost:8000")
	}

## Images for screens with different pixel densities

If you need to add separate images to the resources for screens with different pixel densities, 
then this is done in the style of iOS, i.e. '@< density >x' is appended to the filename. For example

	image@2x.png
	image@3x.jpg
	image@1.5x.gif

For example, you have images for three densities: image.png, image@2x.png, and image@3x.png.
In this case, you only assign the value "image.png" to the "src" field of the ImageView. 
The library itself will find the rest in the "images" directory and transfer the image to the client with the required density

## Themes

The topic includes three types of data:

* constants
* color constants
* View styles

Themes are designed as a rui file and placed in the themes folder.

The root of the theme is an object named 'theme'. This object can contain the following properties:

* name - an optional text property that specifies the name of the theme. 
If this property is not set or is equal to an empty string, then this is the default theme.

* constants - property object defining constants. The name of the object can be anything. It is recommended to use "_".
An object can have any number of text properties specifying the "constant name" = "value" pair.
This section contains constants of type SizeUnit, AngleUnit, text and numeric. In order to assign a constant to any View property, 
you need to assign the name of the constant to the property by adding the '@' symbol at the beginning.
For example

	theme {
		constants = _{
			defaultPadding = 4px,
			buttonPadding = @defaultPadding,
			angle = 30deg,
		}
	}

	rui.Set(view, "subView", rui.Padding, "@defaultPadding")

* constants:touch is property object defining constants used only for touch screen.
For example, how to make indents larger on a touch screen:

	theme {
		constants = _{
			defaultPadding = 4px,
		},
		constants:touch = _{
			defaultPadding = 12px,
		},
	}

* colors is an object property that defines color constants for a light skin (default theme).
An object can have any number of text properties specifying the "color name" = "color" pair. 
Similar to constants, when assigning, you must add '@' at the beginning of the color name. For example

	theme {
		colors = _{
			textColor = #FF101010,
			borderColor = @textColor,
			backgroundColor = white,
		}
	}

	rui.Set(view, "subView", rui.TextColor, "@textColor")

Color names such as "black", "white", "red", etc. are used without the '@' character. 
However, you can specify color constants with the same names. For example

	theme {
		colors = _{
			red = blue,
		}
	}

	rui.Set(view, "subView", rui.TextColor, "@red") // blue text
	rui.Set(view, "subView", rui.TextColor, "red")  // red text

* colors:dark is an object property that defines color constants for a dark theme

* styles is an array of common styles. Each element of the array must be an object. 
The object name is and is the name of the style. For example,

	theme {
		styles = [
			demoPage {
				width = 100%,
				height = 100%,
				cell-width = "1fr, auto",
			},
			demoPanel {
				width = 100%,
				height = 100%,
				orientation = start-to-end,
			},
		]
	}

To use styles, the View has two text properties "style" (Style constant) and "style-disabled" (StyleDisabled constant). 
The "style" property is assigned the property name that is applied to the View when the "disabled" property is set to false. 
The "style-disabled" property is assigned the property name that is applied to the View when the "disabled" property is set to true. 
If "style-disabled" is not specified, then the "style" property is used in both modes.

Attention! The '@' symbol should NOT be added to the style name. If you add the '@' symbol to the name, 
then the style name will be extracted from the constant of the same name. For example

	theme {
		constants = _{
			@demoPanel = demoPage
		},
		styles = [
			demoPage {
				width = 100%,
				height = 100%,
				cell-width = "1fr, auto",
			},
			demoPanel {
				width = 100%,
				height = 100%,
				orientation = start-to-end,
			},
		]
	}

	rui.Set(view, "subView", rui.Style, "demoPanel")   // style == demoPanel
	rui.Set(view, "subView", rui.Style, "@demoPanel")  // style == demoPage

In addition to general styles, you can add styles for specific work modes. To do this, the following modifiers are added to the name "styles":

* ":portrait" or ":landscape" are respectively styles for portrait or landscape mode of the program.
Attention means the aspect ratio of the program window, not the screen.

* ":width< size >" are styles for a screen whose width does not exceed the specified size in logical pixels.

* ":height< size >" are styles for a screen whose height does not exceed the specified size in logical pixels.

For example

	theme {
		styles = [
			demoPage {
				width = 100%,
				height = 100%,
				cell-width = "1fr, auto",
			},
			demoPage2 {
				row = 0,
				column = 1,
			}
		],
		styles:landscape = [
			demoPage {
				width = 100%,
				height = 100%,
				cell-height = "1fr, auto",
			},
			demoPage2 {
				row = 1,
				column = 0,
			}
		],
		styles:portrait:width320 = [
			sapmplePage {
				width = 100%,
				height = 50%,
			},
		]
	}

## Standard constants and styles

The library defines a number of constants and styles. You can override them in your themes.

System styles that you can override:

| Style name          | Описание                                                            |
|---------------------|---------------------------------------------------------------------|
| ruiApp              | This style is used to set the default text style (font, size, etc.) |
| ruiView             | Default View Style                                                  |
| ruiArticle          | The style to use if the "semantics" property is set to "article"    |
| ruiSection          | The style used if the "semantics" property is set to "section"      |
| ruiAside            | The style used if the "semantics" property is set to "aside"        |
| ruiHeader           | The style used if the "semantics" property is set to "header"       |
| ruiMain             | The style used if the "semantics" property is set to "main"         |
| ruiFooter           | Style used if the "semantics" property is set to "footer"           |
| ruiNavigation       | Style used if property "semantics" is set to "navigation"           |
| ruiFigure           | Style used if property "semantics" is set to "figure"               |
| ruiFigureCaption    | Style used if property "semantics" is set to "figure-caption"       |
| ruiButton           | Style used if property "semantics" is set to "button"               |
| ruiParagraph        | The style used if the "semantics" property is set to "paragraph"    |
| ruiH1               | Style used if property "semantics" is set to "h1"                   |
| ruiH2               | Style used if property "semantics" is set to "h2"                   |
| ruiH3               | Style used if property "semantics" is set to "h3"                   |
| ruiH4               | Style used if property "semantics" is set to "h4"                   |
| ruiH5               | Style used if property "semantics" is set to "h5"                   |
| ruiH6               | Style used if property "semantics" is set to "h6"                   |
| ruiBlockquote       | Style used if the "semantics" property is set to "blockquote"       |
| ruiCode             | Style used if property "semantics" is set to "code"                 |
| ruiTable            | Default TableView style                                             |
| ruiTableHead        | Default TableView header style                                      |
| ruiTableFoot        | Default TableView footer style                                      |
| ruiTableRow         | Default TableView row style                                         |
| ruiTableColumn      | Default TableView column style                                      |
| ruiTableCell        | Default TableView cell style                                        |
| ruiDisabledButton   | Button style if property "disabled" is set to true                  |
| ruiCheckbox         | Checkbox style                                                      |
| ruiListItem         | ListView item style                                                 |
| ruiListItemSelected | Style the selected ListView item when the ListView does not have focus |
| ruiListItemFocused  | Style the selected ListView item when the ListView has focus        |
| ruiPopup            | Popup style                                                         |
| ruiPopupTitle       | Popup title style                                                   |
| ruiMessageText      | Popup text style (Message, Question)                                |
| ruiPopupMenuItem    | Popup menu item style                                               |

System color constants that you can override:

| Color constant name        | Description                                         |
|----------------------------|-----------------------------------------------------|
| ruiBackgroundColor         | Background color                                    |
| ruiTextColor               | Text color                                          |
| ruiDisabledTextColor       | Banned text color                                   |
| ruiHighlightColor          | Backlight color                                     |
| ruiHighlightTextColor      | Highlighted text color                              |
| ruiButtonColor             | Button color                                        |
| ruiButtonActiveColor       | Focus button color                                  |
| ruiButtonTextColor         | Button text color                                   |
| ruiButtonDisabledColor     | Denied button color                                 |
| ruiButtonDisabledTextColor | Disabled button text color                          |
| ruiSelectedColor           | Background color of inactive selected ListView item |
| ruiSelectedTextColor       | Text color of inactive selected ListView item       |
| ruiPopupBackgroundColor    | Popup background color                              |
| ruiPopupTextColor          | Popup text color                                    |
| ruiPopupTitleColor         | Popup title background color                        |
| ruiPopupTitleTextColor     | Popup Title Text Color                              |

Constants that you can override:

| Constant name                | Description |
|------------------------------|------------------------------------------------|
| ruiButtonHorizontalPadding   | Horizontal padding inside the button           |
| ruiButtonVerticalPadding     | Vertical padding inside the button             |
| ruiButtonMargin              | External button access                         |
| ruiButtonRadius              | Button corner radius                           |
| ruiButtonHighlightDilation   | Width of the outer border of the active button |
| ruiButtonHighlightBlur       | Blur the active button frame                   |
| ruiCheckboxGap               | Break between checkbox and content             |
| ruiListItemHorizontalPadding | Horizontal padding inside a ListView item      |
| ruiListItemVerticalPadding   | Vertical padding inside a ListView item        |
| ruiPopupTitleHeight          | Popup title height                             |
| ruiPopupTitlePadding         | Popup title padding                            |
| ruiPopupButtonGap            | Break between popup buttons                    |

## Multi-language support

If you want to add support for several languages to the program, you need to place 
the translation files in the "strings" folder of the resources. 
Translation files must have the "rui" extension and the following format

	strings {
		<language 1> = _{
			<text 1> = <translation 1>,
			<text 2> = <translation 2>,
			…
		},
		<язык 2> = _{
			<text 1> = <translation 1>,
			<text 2> = <translation 2>,
			…
		},
		…
	}

If the translation for each language is placed in a separate file, then the following format can be used

	strings:<language> {
		<text 1> = <translation 1>,
		<text 2> = <translation 2>,
		…
	}

For example, if all translations are in one file strings.rui

	strings {
		ru = _{
			"Yes" = "Да",
			"No" = "Нет",
		},
		de = _{
			"Yes" = "Ja",
			"No" = "Nein",
		},
	}

If in different. ru.rui file:

	strings:ru {
		"Yes" = "Да",
		"No" = "Нет",
	}

de.rui file:

	strings:de {
		"Yes" = "Ja",
		"No" = "Nein",
	}

The translation can also be split into multiple files.

Translations are automatically inserted in all Views.

However, if you are drawing text in a CanvasView, then you must request the translation yourself. 
To do this, there is a method in the Session interface:

	GetString(tag string) (string, bool)

If there is no translation of the given string, then the method will return the original string and false as the second parameter.

You can get the current language using the Language() method of the Session interface. 
The current language is determined by the user's browser settings. 
You can change the session language using the SetLanguage(lang string) method of the Session interface.
