# Библиотека RUI

Библиотека RUI (Remoute User Interface) предназначена для создания web приложений на языке go. 

Особенностью библиотеки заключается в том, что вся обработка данных осуществляется на сервере,
а браузер используется как тонкий клиент. Для связи клиента и сервера используется WebSocket.

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
		app := rui.NewApplication("Hello world", "icon.svg", createHelloWorldSession)
		app.Start("localhost:8000")
	}

В функции main создается rui приложение и запускается основной цикл.
При создании приложения задаются 3 параметра: имя приложения, имя иконки и функция createHelloWorldSession.
Функция createHelloWorldSession создает структуру реализующую интерфейс SessionContent:

	type SessionContent interface {
		CreateRootView(session rui.Session) rui.View
	}

Для каждой новой сессии создается свой экземпляр структуры.

Функция CreateRootView интерфейса SessionContent создает корневой элемент.
Когда пользователь обращается к приложению набрав в браузере адрес "localhost:8000", то создается новая сессия,
для нее создается новый экземпляр структуры helloWorldSession и в конце вызывается функция CreateRootView.
Функция createRootView возвращает представление строки текста, создаваемое с помощью функции NewTextView.

Если вы хотите чтобы приложение было видно вне вашего компьютера, то поменяйте адрес в функции Start:

	app.Start(rui.GetLocalIP() + ":80")

## Используемые типы данных

### SizeUnit

Структура SizeUnit используется для задания различных размеров элементов интерфейса, таких как ширина, высота, отступы, размер шрифта и т.п.
SizeUnit объявлена как

	type SizeUnit struct {
		Type  SizeUnitType
		Value float64
	}

где Type - тип размера;
Value - размер

Тип может принимать следующие значения:

| Значение | Константа      | Описание                                                                    |
|:--------:|----------------|-----------------------------------------------------------------------------|
| 0        | Auto           | значение по умолчанию. Значение поля Value игнорируется                     |
| 1        | SizeInPixel    | поле Value определяет размер в пикселях.                                    |
| 2        | SizeInEM       | поле Value определяет размер в em единицах. 1em равен базовому размеру шрифта, который задается в настройках браузера |
| 3        | SizeInEX       | поле Value определяет размер в ex единицах.                                 |
| 4        | SizeInPercent  | поле Value определяет размер в процентах от размера родительского элемента. |
| 5        | SizeInPt       | поле Value определяет размер в pt единицах (1pt = 1/72”).                   |
| 6        | SizeInPc       | поле Value определяет размер в pc единицах (1pc = 12pt).                    |
| 7        | SizeInInch     | поле Value определяет размер в дюймах.                                      |
| 8        | SizeInMM       | поле Value определяет размер в миллиметрах.                                 |
| 9        | SizeInCM       | поле Value определяет размер в сантиметрах.                                 |
| 10       | SizeInFraction | поле Value определяет размер в частях. Используется только для задания размеров ячеек в GridLayout. |

Для более наглядного и простого задания переменных типа SizeUnit могут использоваться функции приведенные ниже

| Функция        | Эквивалентное определение                          |
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

Переменные типа SizeUnit имеют текстовое представление (зачем оно нужно будет описано ниже). Текстовое представление состоит из числа (равному значению поля Value) и следующим за ним суффиксом определяющим тип. Исключением является значение типа Auto, которое имеет представление “auto”. Суффиксы перечислены в следующей таблице:

| Суффикс | Тип            |
|:-------:|----------------|
| px      | SizeInPixel    |
| em      | SizeInEM       |
| ex      | SizeInEX       |
| %       | SizeInPercent  |
| pt      | SizeInPt       |
| pc      | SizeInPc       |
| in      | SizeInInch     |
| mm      | SizeInMM       |
| cm      | SizeInCM       |
| fr      | SizeInFraction |

Примеры: auto, 50%, 32px, 1.5in, 0.8em

Чтобы преобразовать текстовое представление в структуру SizeUnit используется функция:

	func StringToSizeUnit(value string) (SizeUnit, bool)

Получить текстовое представление структуры можно свойством String()

### Color

Тип Color описывает 32-битный цвет в формате ARGB:
	
	type Color uint32

Тип Color имеет три типа текстовых представлений:

1) #AARRGGBB, #RRGGBB, #ARGB, #RGB

где A, R, G, B это шестнадцатеричная цифра описывающая соответствующую компоненту. Если альфа канал не задается, то он считается равным FF. Если цветовая компонента задается одной цифрой, то она удваивается. Например “#48AD” эквивалентно “#4488AADD”

2) argb(A, R, G, B), rgb(R, G, B)

где A, R, G, B это представление цветовой компоненты. Компонента может быть запада в виде дробного числа в диапазоне [0…1] или в виде целого числа в диапазоне [0…255] или в виде процентов от 0% до 100%.
Примеры:
	
	“argb(255, 128, 96, 0)”
	“rgb(1.0, .5, .8)”
	“rgb(0%, 50%, 25%)”
	“argb(50%, 128, .5, 100%)”

Для преобразования Color в строку используется метод String.
Для преобразования строки в Color используется функция:

	func StringToColor(value string) (Color, bool)

3) Имя цвета. В библиотеке определены следующие цвета

| Имя                   | Значение  |
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

Тип AngleUnit используется для задания угловых величин.
AngleUnit объявлена как

	type AngleUnit struct {
		Type  AngleUnitType
		Value float64
	}

где Type - тип угловой величины;
Value - угловая величина

Тип может принимать следующие значения:
* Radian (0) - поле Value определяет угловую величину в радианах.
* PiRadian (1) - поле Value определяет угловую величину в радианах умноженных на π.
* Degree (2) - поле Value определяет угловую величину в градусах.
* Gradian (3) - поле Value определяет угловую величину в градах (градианах).
* Turn (4) - поле Value определяет угловую величину в оборотах (1 оборот == 360°).

Для более наглядного и простого задания переменных типа AngleUnit могут использоваться функции приведенные ниже

| Функция       | Эквивалентное определение                       |
|---------------|-------------------------------------------------|
| rui.Rad(n)   | rui.AngleUnit{ Type: rui.Radian, Value: n }   |
| rui.PiRad(n) | rui.AngleUnit{ Type: rui.PiRadian, Value: n } |
| rui.Deg(n)   | rui.AngleUnit{ Type: rui.Degree, Value: n }   |
| rui.Grad(n)  | rui.AngleUnit{ Type: rui.Gradian, Value: n }  |

Переменные типа AngleUnit имеют текстовое представление состоящее из числа (равному значению поля Value) и следующим за ним суффиксом определяющим тип. Суффиксы перечислены в следующей таблице:

| Суффикс | Тип      |
|:-------:|----------|
| deg     | Degree   |
| °       | Degree   |
| rad     | Radian   |
| π       | PiRadian |
| pi      | PiRadian |
| grad    | Gradian  |
| turn    | Turn     |

Примеры: “45deg”, “90°”, “3.14rad”, “2π”, “0.5pi”

Для преобразования AngleUnit в строку используется метод String.
Для преобразования строки в AngleUnit используется функция:

	func StringToAngleUnit(value string) (AngleUnit, bool)

## View

View это интерфейс для доступа к элементу типа "View". View это прямоугольная область экрана.
Все элементы интерфейса расширяют интерфейс View, т.е. View является базовым элементом для всех
других элементов библиотеки.

View имеет ряд свойств, таких как высота, ширина, цвет, параметры текста и т.д. Каждое свойство
имеет текстовое имя. Для чтения и записи значения свойства используются интерфейс Properties
(View реализует данный интерфейс):

	type Properties interface {
		Get(tag string) interface{}
		Set(tag string, value interface{}) bool
		Remove(tag string)
		Clear()
		AllTags() []string
	}

Функция Get возвращает значение свойства или nil если свойство не установлено.

Функция Set устанавливает значение свойства. Если значение свойства установлено успешно, то
функция возвращает true, если нет то false и в лог записывается описание возникшей ошибки.

Функция Remove удаляет значение свойства, эквивалентно Set(nil)

Для упрощения установки/чтения свойств имеются также две глобальные функции Get и Set:

	func Get(rootView View, viewID, tag string) interface{}
	func Set(rootView View, viewID, tag string, value interface{}) bool

Данные функции возвращают/устанавливают значение дочернего View

### События

При взаимодействии с приложением возникаю различные события: клики, изменение размеров,
изменение вводимых данных и т.п.

Для реакции на события предназначены слушатели событий. Слушатель это функция которая вызывается
каждый раз когда возникает событие. У каждого события может быть несколько слушателей. Разберем
слушателей на примере события изменения текста "edit-text-changed" в редакторе "EditView".

Слушателем события является функция вида

	func(<View>[, <параметры>])

где первый аргумент это View в котором произошло событие. Далее идут дополнительные параметры события.

Для "edit-text-changed" основной слушатель будет иметь следующий вид:

	func(EditView, string)

где второй аргумент это новое значение текста

Если вы не планируте использовать первый аргумент, то его можно опустить. Это будет дополнительный слушатель

	func(string)

Для того чтобы назначить слушателя необходимо его присвоить свойству с именем события

	view.Set(rui.EditTextChanged, func(edit EditView, newText string) {
		// do something
	})

или

	view.Set(rui.EditTextChanged, func(newText string) {
		// do something
	})

У каждого события может быть несколько слушателей. В связи с этим в качестве слушателей могут
использоваться пять типов данных

* функция func(< View >[, <параметры>])
* функция func([<параметры>])
* массив функций []func(< View >[, <параметры>])
* массив функций []func([<параметры>])
* []interface{} содержащий только func(< View >[, <параметры>]) и func([<параметры>])

После присваивания свойству все эти типы преобразуются в массив функций []func(<View>, [<параметры>]).
Соответственно функция Get всегда возвращает массив функций []func(<View>, [<параметры>]).
В случае отсутствия слушателей этот массив будет пуст

Для события "edit-text-changed" это

* func(editor EditView, newText string)
* func(newText string)
* []func(editor EditView, newText string)
* []func(newText string)
* []interface{} содержащий только func(editor EditView, newText string) и func(newText string)

А свойство "edit-text-changed" всегда хранит и возвращает []func(EditView, string).

В дальнейшем при описании конкретных событий будет приводиться только формат основного слушателя.

### Свойство "id"

Свойство "id" это необязательный текстовый идентификатор View. С его помощью можно найти
дочерний View. Для этого используется функция ViewByID

	func ViewByID(rootView View, id string) View

Данная функция ищет дочерний View с идентификатором id. Поиск начинается с rootView.
Если View не найден, то функция возвращает nil и в лог записывается сообщение об ошибке.

Обычно id устанавливается при создании View и в дальнейшем не меняется.
Но это необязательное условие. Вы можете поменять id в любой момент.

Для установки нового значения id используется функция Set. Например

	view.Set(rui.ID, "myView")
	view.Set("id", "myView")

Получить id можно двумя способами. Первый - используя функцию Get:
Например

	if value := view.Get(rui.ID); value != nil {
		id = value.(string)
	}
	
И второй - используя функцию ID():

	id = view.ID()


### Свойства "width", "height", "min-width", "min-height", "max-width", "max-height"

Данные свойства устанавливают:

| Свойство     | Константа     | Описание                 |
|--------------|---------------|--------------------------|
| "width"      | rui.Width     | Ширина View              |
| "height"     | rui.Height    | Высота View              |
| "min-width"  | rui.MinWidth  | Минимальная ширина View  |
| "min-height" | rui.MinHeight | Минимальная высота View  |
| "max-width"  | rui.MaxWidth  | Максимальная ширина View |
| "max-height" | rui.MaxHeight | Максимальная высота View |

Данные свойства имеют тип SizeUnit.
Если значение "width"/"height" не установлены или установлены в Auto, то высота/ширина
View определяется его содержимым и ограничено минимальной и максимальной высотой/шириной.
В качестве значения данных свойств можно установить SizeUnit структуру, текстовое представление
SizeUnit или имя константы (о константах ниже):

	view.Set("width", rui.Px(8))
	view.Set(rui.MaxHeight, "80%")
	view.Set(rui.Height, "@viewHeight")

После получения значения функцией Get вы должны выполнить приведение типов:

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

Это довольно громоздко поэтому для каждого свойства существует одноимённая глобальная функция с префиксом Get,
которая выполняет данное приведение типа, получает значение константы, если необходимо, и
возвращает его. Все функции данного типа имеют два аргумента: View и subviewID string.
Первый аргумент это корневой View, второй - ID дочернего View. Если ID дочернего View передать как "",
то возвращается значение корневого View.
Для свойств "width", "height", "min-width", "min-height", "max-width", "max-height" это функции:

	func GetWidth(view View, subviewID string) SizeUnit
	func GetHeight(view View, subviewID string) SizeUnit
	func GetMinWidth(view View, subviewID string) SizeUnit
	func GetMinHeight(view View, subviewID string) SizeUnit
	func GetMaxWidth(view View, subviewID string) SizeUnit
	func GetMaxHeight(view View, subviewID string) SizeUnit

### Свойства "margin" и "padding"

Свойство "margin" определяет внешние отступы от данного View до соседних.
Свойство "padding" устанавливает внутренние отступы от границы View до контента.
Значение свойств "margin" и "padding" хранятся в виде интерфейса BoundsProperty,
реализующего интерфейс Properties (см. выше). BoundsProperty имеет 4 свойства типа SizeUnit:

| Свойство  | Константа    | Описание         |
|-----------|--------------|------------------|
| "top"     | rui.Top     | Верхний отступ    |
| "right"   | rui.Right   | Правый отступ     |
| "bottom"  | rui.Bottom  | Нижний отступ     |
| "left"    | rui.Left    | Левый отступ      |

Для создания интерфейса BoundsProperty используется функция NewBoundsProperty. Пример

	view.Set(rui.Margin, NewBoundsProperty(rui.Params {
		rui.Top:  rui.Px(8),
		rui.Left: "@topMargin",
		"right":   "1.5em",
		"bottom":  rui.Inch(0.3),
	})))

Соотвественно если вы запросите свойство "margin" или "padding" с помощью метода Get,
то вернется интерфейс BoundsProperty:

	if value := view.Get(rui.Margin); value != nil {
		margin := value.(BoundsProperty)
	}

BoundsProperty с помощью функции "Bounds(session Session) Bounds" интерфейса BoundsProperty
может быть преобразовано в более удобную структуру Bounds:

	type Bounds struct {
		Top, Right, Bottom, Left SizeUnit
	}

Для этого используется также могут использоваться глобальные функции:

	func GetMargin(view View, subviewID string) Bounds
	func GetPadding(view View, subviewID string) Bounds

Текстовое представление BoundsProperty имеет следующий вид:

	"_{ top = <верхний отступ>, right = <правый отступ>, bottom = <нижний отступ>, left = <левый отступ> }"

В качестве значения свойств "margin" и "padding" методу Set может быть передано:
* интерфейс BoundsProperty или его текстовое представление;
* структура Bounds;
* SizeUnit или имя константы типа SizeUnit, в этом случай это значение устанавливается во все отступы. Т.е.

	view.Set(rui.Margin, rui.Px(8))

эквивалентно

	view.Set(rui.Margin, rui.Bounds{Top: rui.Px(8), Right: rui.Px(8), Bottom: rui.Px(8), Left: rui.Px(8)})

Так как значение свойства "margin" и "padding" всегда хранятся в виде интерфейса BoundsProperty,
то если вы прочитаете функцией Get свойство "margin" или "padding" установленное значением Bounds
или SizeUnit, то вы получите BoundsProperty, а не Bounds или SizeUnit.

Свойства "margin" и "padding" используются для установки сразу четырех отступов. Для установки
отдельных отступов используются следующие свойства:

| Свойство         | Константа          | Описание                 |
|------------------|--------------------|--------------------------|
| "margin-top"     | rui.MarginTop     | Верхний внешний отступ    |
| "margin-right"   | rui.MarginRight   | Правый внешний отступ     |
| "margin-bottom"  | rui.MarginBottom  | Нижний внешний отступ     |
| "margin-left"    | rui.MarginLeft    | Левый внешний отступ      |
| "padding-top"    | rui.PaddingTop    | Верхний внутренний отступ |
| "padding-right"  | rui.PaddingRight  | Правый внутренний отступ  |
| "padding-bottom" | rui.PaddingBottom | Нижний внутренний отступ  |
| "padding-left"   | rui.PaddingLeft   | Левый внутренний отступ   |

Например

	view.Set(rui.Margin, rui.Px(8))
	view.Set(rui.TopMargin, rui.Px(12))

эквивалентно

	view.Set(rui.Margin, rui.Bounds{Top: rui.Px(12), Right: rui.Px(8), Bottom: rui.Px(8), Left: rui.Px(8)})

### Свойство "border"

Свойство "border" определяет рамку вокруг View. Линия рамки описывается тремя атрибутами:
стиль линии, толщина и цвет.

Значение свойства "border" хранится в виде интерфейса BorderProperty,
реализующего интерфейс Properties (см. выше). BorderProperty может содержать следующие свойства:

| Свойство       | Константа   | Тип      | Описание                    |
|----------------|-------------|----------|-----------------------------|
| "left-style"   | LeftStyle   | int      | Стиль левой линии рамки     |
| "right-style"  | RightStyle  | int      | Стиль правой линии рамки    |
| "top-style"    | TopStyle    | int      | Стиль верхней линии рамки   |
| "bottom-style" | BottomStyle | int      | Стиль нижней линии рамки    |
| "left-width"   | LeftWidth   | SizeUnit | Толщина левой линии рамки   |
| "right-width"  | RightWidth  | SizeUnit | Толщина правой линии рамки  |
| "top-width"    | TopWidth    | SizeUnit | Толщина верхней линии рамки |
| "bottom-width" | BottomWidth | SizeUnit | Толщина нижней линии рамки  |
| "left-color"   | LeftColor   | Color    | Цвет левой линии рамки      |
| "right-color"  | RightColor  | Color    | Цвет правой линии рамки     |
| "top-color"    | TopColor    | Color    | Цвет верхней линии рамки    |
| "bottom-color" | BottomColor | Color    | Цвет нижней линии рамки     |

Стиль линии может принимать следующие значения:

| Значение | Константа  | Имя      | Описание                 |
|:--------:|------------|----------|--------------------------|
| 0        | NoneLine   | "none"   | Нет рамки                |
| 1        | SolidLine  | "solid"  | Сплошная линия           |
| 2        | DashedLine | "dashed" | Пунктирная линия         |
| 3        | DottedLine | "dotted" | Линия состоящая из точек |
| 4        | DoubleLine | "double" | Двойная сплошная линия   |

Все другие значения стиля игнорируются.

Для создания интерфейса BorderProperty используется функция NewBorder.

Если все линии рамки одинаковы, то для задания стиля, толщины и цвета могут использоваться следующие свойства:

| Свойство | Константа     | Тип      | Описание              |
|----------|---------------|----------|-----------------------|
| "style"  | Style         | int      | Стиль линии рамки     |
| "width"  | Width         | SizeUnit | Толщина линии рамки   |
| "color"  | ColorProperty | Color    | Цвет линии рамки      |

Пример

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

эквивалентно

	view.Set(rui.Border, NewBorder(rui.Params{
		rui.Style: rui.SolidBorder,
		rui.Width: rui.Px(1),
		rui.ColorProperty: rui.Black,
	}))

Интерфейс BorderProperty может быть преобразован в структуру ViewBorders с помощью функции Border.
При преобразовании все текстовые константы заменяются реальными значениями. ViewBorders описана как

	 type ViewBorders struct {
		Top, Right, Bottom, Left ViewBorder
	}

где структура ViewBorder описана как

	type ViewBorder struct {
		Style int
		Color Color
		Width SizeUnit
	}

Структура ViewBorders может быть передана в качестве параметра функции Set при установке значения свойства "border".
При этом ViewBorders преобразуется в BorderProperty. Поэтому при чтении свойства функцией Get будет возвращен интерфейс
BorderProperty, а не структура ViewBorders. Получить структуру ViewBorders без дополнительных преобразований можно
с помощью глобальной функции

	func GetBorder(view View, subviewID string) ViewBorders

Кроме вспомогательных свойств "style", "width" и "color" есть еще 4: "left", "right", "top" и "bottom".
В качестве значения эти свойства могут принимать только структуру ViewBorder и позволяю установить все
атрибуты линии одноименной стороны.

Вы также можете устанавливать отдельные атрибуты рамки использую функцию Set интерфейса View.
Для этого используются следующие свойства

| Свойство              | Константа         | Тип        | Описание                    |
|-----------------------|-------------------|------------|-----------------------------|
| "border-left-style"   | BorderLeftStyle   | int        | Стиль левой линии рамки     |
| "border-right-style"  | BorderRightStyle  | int        | Стиль правой линии рамки    |
| "border-top-style"    | BorderTopStyle    | int        | Стиль верхней линии рамки   |
| "border-bottom-style" | BorderBottomStyle | int        | Стиль нижней линии рамки    |
| "border-left-width"   | BorderLeftWidth   | SizeUnit   | Толщина левой линии рамки   |
| "border-right-width"  | BorderRightWidth  | SizeUnit   | Толщина правой линии рамки  |
| "border-top-width"    | BorderTopWidth    | SizeUnit   | Толщина верхней линии рамки |
| "border-bottom-width" | BorderBottomWidth | SizeUnit   | Толщина нижней линии рамки  |
| "border-left-color"   | BorderLeftColor   | Color      | Цвет левой линии рамки      |
| "border-right-color"  | BorderRightColor  | Color      | Цвет правой линии рамки     |
| "border-top-color"    | BorderTopColor    | Color      | Цвет верхней линии рамки    |
| "border-bottom-color" | BorderBottomColor | Color      | Цвет нижней линии рамки     |
| "border-style"        | BorderStyle       | int        | Стиль линии рамки           |
| "border-width"        | BorderWidth       | SizeUnit   | Толщина линии рамки         |
| "border-color"        | BorderColor       | Color      | Цвет линии рамки            |
| "border-left"         | BorderLeft        | ViewBorder | Левая линия рамки           |
| "border-right"        | BorderRight       | ViewBorder | Правая линия рамки          |
| "border-top"          | BorderTop         | ViewBorder | Верхняя линия рамки         |
| "border-bottom"       | BorderBottom      | ViewBorder | Нижняя линия рамки          |

Например

	view.Set(rui.BorderStyle, rui.SolidBorder)
	view.Set(rui.BorderWidth, rui.Px(1))
	view.Set(rui.BorderColor, rui.Black)

эквивалентно

	view.Set(rui.Border, NewBorder(rui.Params{
		rui.Style: rui.SolidBorder,
		rui.Width: rui.Px(1),
		rui.ColorProperty: rui.Black,
	}))

### Свойство "radius"

Свойство "radius" задает эллиптический радиус скругления углов View. Радиусы задаются интерфейсом
RadiusProperty реализующим интерфейс Properties (см. выше).
Для этого используются следующие свойства типа SizeUnit:

| Свойство         | Константа    | Описание                       |
|------------------|--------------|--------------------------------|
| "top-left-x"     | TopLeftX     | x-радиус верхнего левого угла  |
| "top-left-y"     | TopLeftY     | y-радиус верхнего левого угла  |
| "top-right-x"    | TopRightX    | x-радиус верхнего правого угла |
| "top-right-y"    | TopRightY    | y-радиус верхнего правого угла |
| "bottom-left-x"  | BottomLeftX  | x-радиус нижнего левого угла   |
| "bottom-left-y"  | BottomLeftY  | y-радиус нижнего левого угла   |
| "bottom-right-x" | BottomRightX | x-радиус нижнего правого угла  |
| "bottom-right-y" | BottomRightY | y-радиус нижнего правого угла  |

Если x- и y-радиусы одинаковы то можно воспользоваться вспомогательными свойствами

| Свойство       | Константа    | Описание                     |
|----------------|--------------|------------------------------|
| "top-left"     | TopLeft      | радиус верхнего левого угла  |
| "top-right"    | TopRight     | радиус верхнего правого угла |
| "bottom-left"  | BottomLeft   | радиус нижнего левого угла   |
| "bottom-right" | BottomRight  | радиус нижнего правого угла  |

Для установки всех радиусов одинаковыми значениями используются свойства "x" и "y"

Интерфейс RadiusProperty создается с помощью функции NewRadiusProperty. Пример

	view.Set(rui.Radius, NewRadiusProperty(rui.Params{
		rui.X: rui.Px(16),
		rui.Y: rui.Px(8),
		rui.TopLeft: rui.Px(0),
		rui.BottomRight: rui.Px(0),
	}))

эквивалентно

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

Если все радиусы одинаковы, то данное значение типа SizeUnit может быть напрямую присвоено свойству "radius"

	view.Set(rui.Radius, rui.Px(4))

RadiusProperty имеет текстовое представление следующего вида:

	_{ <id радиуса> = <SizeUnit text> [/ <SizeUnit text>] [, <id радиуса> = <SizeUnit text> [/ <SizeUnit text>]] … }

где <id радиуса> может принимать следующие значения:

| Свойство         | Описание                            |
|------------------|-------------------------------------|
| "x"              | Все x-радиусы                       |
| "y"              | Все y-радиусы                       |
| "top-left"       | x- и y-радиус верхнего левого угла  |
| "top-left-x"     | x-радиус верхнего левого угла       |
| "top-left-y"     | y-радиус верхнего левого угла       |
| "top-right"      | x- и y-радиус верхнего правого угла |
| "top-right-x"    | x-радиус верхнего правого угла      |
| "top-right-y"    | y-радиус верхнего правого угла      |
| "bottom-left"    | x- и y-радиус нижнего левого угла   |
| "bottom-left-x"  | x-радиус нижнего левого угла        |
| "bottom-left-y"  | y-радиус нижнего левого угла        |
| "bottom-right"   | x- и y-радиус нижнего правого угла  |
| "bottom-right-x" | x-радиус нижнего правого угла       |
| "bottom-right-y" | y-радиус нижнего правого угла       |

Значения вида "<SizeUnit text> / <SizeUnit text>" можно присваивать только
свойствам "top-left", "top-right", "bottom-left" и "bottom-right".

Примеры:

	_{ x = 4px, y = 4px, top-left = 8px, bottom-right = 8px }

эквивалентно

	_{ top-left = 8px, top-right = 4px, bottom-left = 4px, bottom-right = 8px }

или

	_{ top-left = 8px / 8px, top-right = 4px / 4px, bottom-left = 4px / 4px, bottom-right = 8px / 8px }
	
или

	_{ top-left-x = 8px, top-left-y = 8px, top-right-x = 4px, top-right-y = 4px,
		bottom-left-x = 4px, bottom-left-y = 4px, bottom-right-x = 8px, bottom-right-y = 8px }

Интерфейс RadiusProperty может быть преобразован в структуру BoxRadius с помощью функции BoxRadius.
При преобразовании все текстовые константы заменяются реальными значениями. BoxRadius описана как

	type BoxRadius struct {
		TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY SizeUnit
	}

Структура BoxRadius может быть передана в качестве параметра функции Set при установке значения свойства "radius".
При этом BoxRadius преобразуется в RadiusProperty. Поэтому при чтении свойства функцией Get будет возвращен интерфейс
RadiusProperty, а не структура BoxRadius. Получить структуру BoxRadius без дополнительных преобразований можно
с помощью глобальной функции

	func GetRadius(view View, subviewID string) BoxRadius

Вы также можете устанавливать отдельные радиусы использую функцию Set интерфейса View.
Для этого используются следующие свойства

| Свойство                | Константа          | Описание                            |
|-------------------------|--------------------|-------------------------------------|
| "radius-x"              | RadiusX            | Все x-радиусы                       |
| "radius-y"              | RadiusY            | Все y-радиусы                       |
| "radius-top-left"       | RadiusTopLeft      | x- и y-радиус верхнего левого угла  |
| "radius-top-left-x"     | RadiusTopLeftX     | x-радиус верхнего левого угла       |
| "radius-top-left-y"     | RadiusTopLeftY     | y-радиус верхнего левого угла       |
| "radius-top-right"      | RadiusTopRight     | x- и y-радиус верхнего правого угла |
| "radius-top-right-x"    | RadiusTopRightX    | x-радиус верхнего правого угла      |
| "radius-top-right-y"    | RadiusTopRightY    | y-радиус верхнего правого угла      |
| "radius-bottom-left"    | RadiusBottomLeft   | x- и y-радиус нижнего левого угла   |
| "radius-bottom-left-x"  | RadiusBottomLeftX  | x-радиус нижнего левого угла        |
| "radius-bottom-left-y"  | RadiusBottomLeftY  | y-радиус нижнего левого угла        |
| "radius-bottom-right"   | RadiusBottomRight  | x- и y-радиус нижнего правого угла  |
| "radius-bottom-right-x" | RadiusBottomRightX | x-радиус нижнего правого угла       |
| "radius-bottom-right-y" | RadiusBottomRightY | y-радиус нижнего правого угла       |

Например

	view.Set(rui.RadiusX, rui.Px(4))
	view.Set(rui.RadiusY, rui.Px(32))

эквивалентно

	view.Set(rui.Border, NewRadiusProperty(rui.Params{
		rui.X: rui.Px(4),
		rui.Y: rui.Px(32),
	}))

### Свойство "shadow"

Свойство "shadow" позволяет задать тени для View. Теней может быть несколько. Тень описывается
с помощью интерфейса ViewShadow расширяющего интерфейс Properties (см. выше). У тени имеются следующие свойства:

| Свойство        | Константа     | Тип      | Описание                                                      |
|-----------------|---------------|----------|---------------------------------------------------------------|
| "color"         | ColorProperty | Color    | Цвет тени                                                     |
| "inset"         | Inset         | bool     | true - тень внутри View, false - снаружи                      |
| "x-offset"      | XOffset       | SizeUnit | Смещение тени по оси X                                        |
| "y-offset"      | YOffset       | SizeUnit | Смещение тени по оси Y                                        |
| "blur"          | BlurRadius    | float    | Радиус размытия тени. Значение должно быть >= 0               |
| "spread-radius" | SpreadRadius  | float    | Увеличение тени. Значение > 0 увеличивает тень, < 0 уменьшает |

Для создания ViewShadow используются три функции:

	func NewViewShadow(offsetX, offsetY, blurRadius, spread-radius SizeUnit, color Color) ViewShadow
	func NewInsetViewShadow(offsetX, offsetY, blurRadius, spread-radius SizeUnit, color Color) ViewShadow
	func NewShadowWithParams(params Params) ViewShadow

Функция NewViewShadow создает внешнюю тень (Inset == false), NewInsetViewShadow - внутреннюю
(Inset == true).
Функция NewShadowWithParams используется когда в качестве параметров необходимо использовать
константы. Например:

	shadow := NewShadowWithParams(rui.Params{
		rui.ColorProperty : "@shadowColor",
		rui.BlurRadius : 8.0,
		rui.Dilation : 16.0,
	})

В качестве значения свойству "shadow" может быть присвоено ViewShadow, массив ViewShadow,
текстовое представление ViewShadow.

Текстовое представление ViewShadow имеет следующий формат:

	_{ color = <цвет> [, x-offset = <смещение>] [, y-offset = <смещение>] [, blur = <радиус>]
		[, spread-radius = <увеличени>] [, inset = <тип>] }


Получить значение данного свойства можно с помощью функции

	func GetViewShadows(view View, subviewID string) []ViewShadow

Если тень не задана, то данная функция вернет пустой массив

### Свойство "background-color"

Константа: rui.BackgroundColor. Get функция: BackgroundColor() Color

Свойство "background-color" устанавливает цвет фона. Допустимые значения: Color, целое число, текстовой представление Color и имя константы начинающееся с '@'.
Целое число должно кодировать цвет в формате AARRGGBB

Кроме цвета в качестве фона можно также использовать изображения и градиенты (см. ниже).
В этом случае "background-color" используется для прозрачных участков изображений.

### Свойство "background-clip"

Свойство "background-clip" определяет как цвет фона и/или фоновое изображение будут выводиться под границами блока.

Если фоновое изображение или цвет не заданы, это свойство будет иметь визуальный эффект, только если у границы есть прозрачные области или частично непрозрачные области; в противном случае граница скрывает разницу.

Свойство может принимать следующие значения:

| Значение | Константа      | Имя           | Описание                                       |
|:--------:|----------------|---------------|------------------------------------------------|
| 0        | BorderBoxClip  | "border-box"  | Фон распространяется до внешнего края границы (но под границей в z-порядке). |
| 1        | PaddingBoxClip | "padding-box" | Фон распространяется до внешнего края отступа. Под границей фон не рисуется. |
| 2        | ContentBoxClip | "content-box" | Фон закрашивается внутри (обрезается) поля содержимого. |

### Свойство "background"

В качестве фона View, помимо цвета, можно задать также картинки и/или градиентные заливки.
Для этого используется свойство "background". Фон может содержать несколько картинок и градиентов.
Каждый элемент фона описывается интерфейсом BackgroundElement. BackgroundElement может быть трех
типов: линейный градиент, радиальный градиент и изображение.

#### Линейный градиент

Линейный градиент создается с помощью функции

	func NewBackgroundLinearGradient(params Params) BackgroundElement

Линейный градиент имеет следующие параметры:

* Direction ("direction") - определяет направление линии градиента (линии вдоль которой меняется цвет).
Необязательный параметр. Направление по умолчанию - снизу вверх.
Может принимать или значение типа AngleUnit (угол наклона линии относительно вертикали)
или одно из следующих значений типа Int:

| Значение | Константа             | Имя               | Описание                                       |
|:--------:|-----------------------|-------------------|------------------------------------------------|
| 0        | ToTopGradient         | "to-top"          | Линия идет снизу вверх (значение по умолчанию) |
| 1        | ToRightTopGradient    | "to-right-top"    | Из левого нижнего угла в правый верхний        |
| 2        | ToRightGradient       | "to-right"        | Слева направо                                  |
| 3        | ToRightBottomGradient | "to-right-bottom" | Из левого верхнего угла в правый нижний        |
| 4        | ToBottomGradient      | "to-bottom"       | Сверху вних                                    |
| 5        | ToLeftBottomGradient  | "to-left-bottom"  | Из правого верхнего угла в левый нижний        |
| 6        | ToLeftGradient        | "to-left"         | Справа налево                                  |
| 7        | ToLeftTopGradient     | "to-left-top"     | Из правого нижнего угла в левый верхний        |

* Gradient ("gradient") - массив ключевых точек градиента (обязательный параметр). Каждая точка
описывается структурой BackgroundGradientPoint, которая имеет два поля: Pos типа SizeUnit и Color.
Pos определяет положение точки относительно начала линии градиента. Массив должен иметь не менее 2 точек.
В качестве значения градиента можно также передать массив Color. В этом случае точки равномерно
распределяются вдоль линии градиента.
Также  в качестве массива ключевых точек можно использовать массив типа []interface{}.
Элементами этого массива могут быть BackgroundGradientPoint, Color, текстовое представление BackgroundGradientPoint
или Color и имя константы

* Repeat ("repeat") - булево значение, определяющее будет ли повторяться градиент после последней
ключевой точки. Необязательный параметр. Значение по умолчанию - false (не повторять)

Текстовое представление линейного градиента имеет следующий вид:

	linear-gradient { gradient = <значение> [, direction = <значение>] [, repeat = <значение>] }

#### Радиальный градиент

Радиальный градиент создается с помощью функции

	func NewBackgroundRadialGradient(params Params) BackgroundElement

Радиальный градиент имеет следующие параметры:

* Gradient ("gradient") - массив ключевых точек градиента (обязательный параметр). Идентичен одноименному
параметру линейного градиента.

* Repeat ("repeat") - булево значение, определяющее будет ли повторяться градиент после последней
ключевой точки. Необязательный параметр. Значение по умолчанию - false (не повторять)

* RadialGradientShape ("radial-gradient-shape") или Shape ("shape") - определяет форму градиента.
Может принимать одно из двух значений типа Int:

| Значение | Константа       | Имя       | Описание                                   |
|:--------:|-----------------|-----------|--------------------------------------------|
| 0        | EllipseGradient | "ellipse" | Формой является эллипс, выровненный по оси |
| 1        | CircleGradient  | "circle"  | Формой является круг с постоянным радиусом |

Необязательный параметр. Значение по умолчанию EllipseGradient

* RadialGradientRadius ("radial-gradient-radius") или Radius ("radius") - задает радиус градиента.
Может принимать или значение типа SizeUnit или одно из следующих значений типа Int:

| Константа              | Значение | Имя               | Описание                                   |
|------------------------|:--------:|-------------------|--------------------------------------------|
| ClosestSideGradient    | 0        | "closest-side"    | Конечная форма градиента соответствует стороне прямоугольника, ближайшей к его центру (для окружностей), или обеим вертикальным и горизонтальным сторонам, ближайшим к центру (для эллипсов) |
| ClosestCornerGradient  | 1        | "closest-corner"  | Конечная форма градиента определяется таким образом, чтобы он точно соответствовал ближайшему углу окна от его центра |
| FarthestSideGradient   | 2        | "farthest-side"   | Схоже с ClosestSideGradient, кроме того что, размер формы определяется самой дальней стороной от своего центра (или вертикальных и горизонтальных сторон) |
| FarthestCornerGradient | 3        | "farthest-corner" | Конечная форма градиента определяется таким образом, чтобы он точно соответствовал самому дальнему углу прямоугольника от его центра |

Необязательный параметр. Значение по умолчанию ClosestSideGradient

* CenterX ("center-x"), CenterY ("center-y") - задает центр градиента относительно левого верхнего
угла View. Принимает значение типа SizeUnit. Необязательный параметр.
Значение по умолчанию "50%", т.е. центр градиента совпадает с центром View.

Текстовое представление линейного градиента имеет следующий вид:

	radial-gradient { gradient = <значение> [, repeat = <значение>] [, shape = <значение>]
		[, radius = <значение>][, center-x = <значение>][, center-y = <значение>]}

#### Изображение

Изображение имеет следующие параметры:

* Source ("src") - задает URL изображения

* Fit ("fit") - необязательный параметр определяющий масштабирование изображения.
Может принимать одно из следующих значений типа Int:

| Константа  | Значение | Имя       | Описание                                                         |
|------------|:--------:|-----------|------------------------------------------------------------------|
| NoneFit    | 0        | "none"    | Нет масштабирования (значение по умолчанию). Размеры изображения определяются параметрами Width и Height |
| ContainFit | 1        | "contain" | Изображение масштабирует с сохранением пропорций так, чтобы его ширина или высота равнялась ширине или высоте области фона. Изображение может обрезаться по ширине или высоте |
| CoverFit   | 2        | "cover"   | Изображение масштабирует с сохранением пропорций так, чтобы картинка целиком поместилась внутрь области фона |

* Width ("width"), Height (height) - необязательные SizeUnit параметры задающие высоту и ширину
изображения. Используется только если параметр Fit равен NoneFit. Значение по умолчанию Auto (исходный
размер). Значение в процентах задает размер относительно высоты и ширины области фона соответственно

* Attachment - ???

* Repeat (repeat) - необязательный параметр задающий повтор изображения.
Может принимать одно из следующих значений типа Int:

| Константа   | Значение | Имя         | Описание                                                      |
|-------------|:--------:|-------------|---------------------------------------------------------------|
| NoRepeat    | 0        | "no-repeat" | Изображение не повторяется (значение по умолчанию)            |
| RepeatXY    | 1        | "repeat"    | Изображение повторяется по горизонтали и вертикали            |
| RepeatX     | 2        | "repeat-x"  | Изображение повторяется только по горизонтали                 |
| RepeatY     | 3        | "repeat-y"  | Изображение повторяется только по вертикали                   |
| RepeatRound | 4        | "round"     | Изображение повторяется так, чтобы в область фона поместилось целое число рисунков; если это не удаётся сделать, то фоновые рисунки масштабируются |
| RepeatSpace | 5        | "space"     | Изображение повторяется столько раз, сколько требуется для заполнения области фона; если это не удаётся, между картинками добавляется пустое пространство |

* ImageHorizontalAlign,
* ImageVerticalAlign,

### Свойство "clip"

Свойство "clip" (константа Clip) типа ClipShape задает задает область образки.
Есть 4 типа областей обрезки

#### inset

Прямоугольная область обрезки. Создается с помощью функции:

	func InsetClip(top, right, bottom, left SizeUnit, radius RadiusProperty) ClipShape

где top, right, bottom, left это расстояние от соответственно верхней, правой, нижней и левой границы
View до одноименной границы обрезки; radius - задает радиусы скругления углов области обрезки
(описание типа RadiusProperty смотри выше). Если скругления углов не должно быть, то в качестве
radius необходимо передать nil

Текстовое описание прямоугольной области обрезки имеет следующий формат

	inset{ top = <top value>, right = <right value>, bottom = <bottom value>, left = <left value>,
		[radius = <RadiusProperty text>] }
	}

#### circle

Круглая область обрезки. Создается с помощью функции:

	func CircleClip(x, y, radius SizeUnit) ClipShape

где x, y - координаты центра окружности; radius - радиус

Текстовое описание круглой области обрезки имеет следующий формат

	circle{ x = <x value>, y = <y value>, radius = <radius value> }

#### ellipse

Эллиптическая область обрезки. Создается с помощью функции:

	func EllipseClip(x, y, rx, ry SizeUnit) ClipShape

где x, y - координаты центра эллипса; rх - радиус эллипса по оси X; ry - радиус эллипса по оси Y.

Текстовое описание эллиптической области обрезки имеет следующий формат

	ellipse{ x = <x value>, y = <y value>, radius-x = <x radius value>, radius-y = <y radius value> }

#### polygon

Многоугольная область обрезки. Создается с помощью функций:

	func PolygonClip(points []interface{}) ClipShape
	func PolygonPointsClip(points []SizeUnit) ClipShape

в качестве аргумента передается массив угловых точек многоугольника в следующем порядке: x1, y1, x2, y2, …
В качестве элементов аргумента функции PolygonClip могут быть или текстовые константы, или
текстовое представление SizeUnit, или элементы типа SizeUnit.

Текстовое описание многоугольной области обрезки имеет следующий формат

	polygon{ points = "<x1 value>, <y1 value>, <x2 value>, <y2 value>,…" }

### Свойство "оpacity"

Свойство "оpacity" (константа Opacity) типа float64 задает прозрачность View. Допустимые значения от 0 до 1.
Где 1 - View полностью непрозрачен, 0 - полностью прозрачен.

Получить значение данного свойства можно с помощью функции

	func GetOpacity(view View, subviewID string) float64

### Свойство "z-index"

Свойство "z-index" (константа ZIndex) типа int определяет положение элемента и нижестоящих элементов по оси z.
В случае перекрытия элементов, это значение определяет порядок наложения. В общем случае, элементы
с большим z-index перекрывают элементы с меньшим.

Получить значение данного свойства можно с помощью функции

	func GetZIndex(view View, subviewID string) int

### Свойство "visibility"

Свойство "visibility" (константа Visibility) типа int задает видимость View. Допустимые значения

| Значение | Константа | Имя         | Видимость                          |
|:--------:|-----------|-------------|------------------------------------|
| 0        | Visible   | "visible"   | View видим. Значение по умолчанию. |
| 1        | Invisible | "invisible" | View невидим, но занимает место.   |
| 2        | Gone      | "gone"      | View невидим и не занимает место.  |

Получить значение данного свойства можно с помощью функции

	func GetVisibility(view View, subviewID string) int

### Свойство "filter"

Свойство "filter" (константа Filter) применяет ко View такие графические эффекты, как размытие и смещение цвета.
В качестве значения свойства "filter" используется только интерфейс ViewFilter. ViewFilter создается с помощью
функции

	func NewViewFilter(params Params) ViewFilter

В аргументе перечисляются применяемые эффекты. Возможны следующие эффекты:

| Эффект        | Константа  | Тип                |  Описание                        |
|---------------|------------|--------------------|----------------------------------|
| "blur"        | Blur       | float64  0…10000px | Размытие по Гауссу               |
| "brightness"  | Brightness | float64  0…10000%  | Изменение яркости                |
| "contrast"    | Contrast   | float64  0…10000%  | Изменение контрастности          |
| "drop-shadow" | DropShadow | []ViewShadow       | Добавление тени                  |
| "grayscale"   | Grayscale  | float64  0…100%    | Преобразование к оттенкам серого |
| "hue-rotate"  | HueRotate  | AngleUnit          | Вращение оттенка                 |
| "invert"      | Invert     | float64  0…100%    | Инвертирование цветов            |
| "opacity"     | Opacity    | float64  0…100%    | Изменение прозрачности           |
| "saturate"    | Saturate   | float64  0…10000%  | Изменение насыщености            |
| "sepia"       | Sepia      | float64  0…100%    | Преобразование в серпию          |

Получить значение текущего фильтра можно с помощью функции

	func GetFilter(view View, subviewID string) ViewFilter

### Свойство "semantics"

Свойство "semantics" (константа Semantics) типа string определяет семантический смысл View.
Данное свойство может не иметь видимого эффекта, но позволяет поисковикам понимать структуру вашего приложения.
Так же оно помогает озвучивать интерфейс системам для людей с ограниченными возможностями:

| Значение | Имя              | Семантика                                            |
|:--------:|------------------|------------------------------------------------------|
| 0        | "default"        | Не определена. Значение по умолчанию.                |
| 1        | "article"        | Самостоятельная часть приложения предназначенная для независимого распространения или повторного использования. |
| 2        | "section"        | Автономный раздел который не может быть представлен более точным по семантике элементом |
| 3        | "aside"          | Часть документа, чьё содержимое только косвенно связанно с основным содержимым (сноска, метка) |
| 4        | "header"         | Заголовок приложения                                 |
| 5        | "main"           | Основной контент (содержимое) приложения             |
| 6        | "footer"         | Нижний колонтитул                                    |
| 7        | "navigation"     | Панель навигации                                     |
| 8        | "figure"         | Изображение                                          |
| 9        | "figure-caption" | Заголовок Изображения. Должно быть внутри "figure"   |
| 10       | "button"         | Кнопка                                               |
| 11       | "p"              | Параграф                                             |
| 12       | "h1"             | Заголовок текста 1-го уровня. Изменяет стиль текста  |
| 13       | "h2"             | Заголовок текста 2-го уровня. Изменяет стиль текста  |
| 14       | "h3"             | Заголовок текста 3-го уровня. Изменяет стиль текста  |
| 15       | "h4"             | Заголовок текста 4-го уровня. Изменяет стиль текста  |
| 16       | "h5"             | Заголовок текста 5-го уровня. Изменяет стиль текста  |
| 17       | "h6"             | Заголовок текста 6-го уровня. Изменяет стиль текста  |
| 18       | "blockquote"     | Цитата. Изменяет стиль текста                        |
| 19       | "code"           | Программный код. Изменяет стиль текста               |

### Свойства текста

Все перечисленные в этом разделе свойства являются наследуемыми, т.е. свойство будет применяться не только ко View
для которого оно установлено, но и ко всем View вложенным в него.

Имеются следующие свойства для настройки параметров отображения текста:

#### Свойство "font-name"

Свойство "font-name" (константа FontName) - текстовое свойство определяет имя используемого шрифта.
Может задаваться несколько шрифтов. В этом случае они разделяются пробелом.
Шрифты применяются в том порядке в котором они перечислены. Т.е. сначала
применяется первый, если он недоступен, то второй, третий и т.д.

Получить значение данного свойства можно с помощью функции

	func GetFontName(view View, subviewID string) string

#### Свойство "text-color"

Свойство "text-color" (константа TextColor) - свойство типа Color определяет цвет текста.

Получить значение данного свойства можно с помощью функции

	func GetTextColor(view View, subviewID string) Color

#### Свойство "text-size"

Свойство "text-size" (константа TextSize) - свойство типа SizeUnit определяет размер шрифта.

Получить значение данного свойства можно с помощью функции

	func GetTextSize(view View, subviewID string) SizeUnit

#### Свойство "italic"

Свойство "italic" (константа Italic) - свойство типа bool. Если значение равно true, то к тексту применяется курсивное начертание

Получить значение данного свойства можно с помощью функции

	func IsItalic(view View, subviewID string) bool
	
#### Свойство "small-caps"

Свойство "small-caps" (константа SmallCaps) - свойство типа bool. Если значение равно true, то к тексту применяется начертание капителью

Получить значение данного свойства можно с помощью функции

	func IsSmallCaps(view View, subviewID string) bool

#### Свойство "white-space"

Свойство "white-space" (константа WhiteSpace) типа int управляет тем, как обрабатываются пробельные
символы внутри View. Свойство "white-space" может принимать следующие значения:

0 (константа WhiteSpaceNormal, имя "normal") - последовательности пробелов объединяются в один пробел.
Символы новой строки в источнике обрабатываются, как отдельный пробел. Применение данного значения
при необходимости разбивает строки для того, чтобы заполнить строчные боксы.

1 (константа WhiteSpaceNowrap, имя "nowrap") - объединяет последовательности пробелов в один пробел,
как значение normal, но не переносит строки (оборачивание текста) внутри текста.

2 (константа WhiteSpacePre, имя "pre") - последовательности пробелов сохраняются так, как они указаны
в источнике. Строки переносятся только там, где в источнике указаны символы новой строки и там,
где в источнике указаны элементы "br".

3 (константа WhiteSpacePreWrap, имя "pre-wrap") - последовательности пробелов сохраняются так, как они
указаны в источнике. Строки переносятся только там, где в источнике указаны символы новой строки и там,
где в источнике указаны элементы "br", и при необходимости для заполнения строчных боксов.

4 (константа WhiteSpacePreLine, имя "pre-line") - последовательности пробелов объединяются в один пробел.
Строки разбиваются по символам новой строки, по элементам "br", и при необходимости для заполнения строчных боксов.

5 (константа WhiteSpaceBreakSpaces, имя "break-spaces") - поведение идентично pre-wrap со следующими отличиями:
* Последовательности пробелов сохраняются так, как они указаны в источнике, включая пробелы на концах строк.
* Строки переносятся по любым пробелам, в том числе в середине последовательности пробелов.
* Пробелы занимают место и не висят на концах строк, а значит влияют на внутренние размеры (min-content и max-content).

В приведённой ниже таблице указано поведение различных значений свойства "white-space"

|                       | Новые строки                | Пробелы и табуляция         | Перенос по словам | Пробелы в конце строки      |
|-----------------------|-----------------------------|-----------------------------|-------------------|-----------------------------|
| WhiteSpaceNormal      | Объединяются в одну         | Объединяются в один пробел  | Переносится       | Удаляются                   |
| WhiteSpaceNowrap      | Объединяются в одну         | Объединяются в один пробел  | Не переносится    | Удаляются                   |
| WhiteSpacePre         | Сохраняются как в источнике | Сохраняются как в источнике | Не переносится    | Сохраняются как в источнике |
| WhiteSpacePreWrap     | Сохраняются как в источнике | Сохраняются как в источнике | Переносится       | Висят                       |
| WhiteSpacePreLine     | Сохраняются как в источнике | Объединяются в один пробел  | Переносится       | Удаляются                   |
| WhiteSpaceBreakSpaces | Сохраняются как в источнике | Сохраняются как в источнике | Переносится       | Переносятся                 |

#### Свойство "word-break"

Свойство "word-break" (константа WordBreak) типа int определяет, где будет установлен перевод
на новую строку в случае превышения текстом границ блока.
Свойство "white-space" может принимать следующие значения:

0 (константа WordBreak, имя "normal) - поведение по умолчанию для расстановки перевода строк.

1 (константа WordBreakAll, имя "break-all) - при превышении границ блока, перевод строки будет
вставлен между любыми двумя символами (за исключением текста на китайском/японском/корейском языке).

2 (константа WordBreakKeepAll, имя "keep-all) - перевод строки не будет использован в тексте на
китайском/японском/корейском языке. Для текста на других языках будет применено поведение по умолчанию (normal).

3 (константа WordBreakWord, имя "break-word) - при превышении границ блока, остающиеся целыми слова
могут быть разбиты в произвольном месте, если не будет найдено более подходящее для переноса строки место.

#### Свойства "strikethrough", "overline" и "underline"

Данные свойства устанавливают декоративные линии на тексте:

| Свойство        | Константа      | Тип декоративной линии      |
|-----------------|----------------|-----------------------------|
| "strikethrough" | Strikethrough  | Линия перечеркивающая текст |
| "overline"      | Overline       | Линия над текстом           |
| "underline"     | Underline      | Линия под текстом           |

Получить значение данных свойств можно с помощью функций

	func IsStrikethrough(view View, subviewID string) bool
	func IsOverline(view View, subviewID string) bool
	func IsUnderline(view View, subviewID string) bool

#### Свойство "text-line-thickness"

Свойство "text-line-thickness" (константа TextLineThickness) - свойство типа SizeUnit.
Свойство устанавливает толщину декоративных линий на тексте заданных с помощью свойств "strikethrough", "overline" и "underline".

Получить значение данного свойства можно с помощью функции

	GetTextLineThickness(view View, subviewID string) SizeUnit

#### Свойство "text-line-style"

Свойство "text-line-style" (константа TextLineStyle) - свойство типа int.
Свойство устанавливает стиль декоративных линий на тексте заданных с помощью свойств "strikethrough", "overline" и "underline".
Возможны следующие значения:

| Значение | Константа  | Имя      | Описание                 |
|:--------:|------------|----------|--------------------------|
| 1        | SolidLine  | "solid"  | Сплошная линия           |
| 2        | DashedLine | "dashed" | Пунктирная линия         |
| 3        | DottedLine | "dotted" | Линия состоящая из точек |
| 4        | DoubleLine | "double" | Двойная сплошная линия   |
| 5        | WavyLine   | "wavy"   | Волнистая линия          |

Если свойство не определено то используется сплошная линия (SolidLine (1)).

Получить значение данного свойства можно с помощью функции

	func GetTextLineStyle(view View, subviewID string) int

#### Свойство "text-line-color"

Свойство "text-line-color" (константа TextLineColor) - свойство типа Color.
Свойство устанавливает цвет декоративных линий на тексте заданных с помощью свойств  "strikethrough", "overline" и "underline".
Если свойство не определено то для линий используется цвет текста заданный с помощью свойства "text-color".

Получить значение данного свойства можно с помощью функции

	func GetTextLineColor(view View, subviewID string) Color

#### Свойство "text-weight"

Свойство "text-weight" (константа TextWeight) - свойство типа int устанавливает начертание шрифта. Допустимые значения:

| Значение | Константа      | Общее название начертания                                       |
|:--------:|----------------|-----------------------------------------------------------------|
| 1	       | ThinFont       | Тонкий (Волосяной) Thin (Hairline)                              |
| 2	       | ExtraLightFont | Дополнительный светлый (Сверхсветлый) Extra Light (Ultra Light) |
| 3	       | LightFont      | Светлый Light                                                   |
| 4	       | NormalFont     | Нормальный Normal. Значение по умолчанию                        |
| 5	       | MediumFont     | Средний Medium                                                  |
| 6	       | SemiBoldFont   | Полужирный Semi Bold (Demi Bold)                                |
| 7	       | BoldFont       | Жирный Bold                                                     |
| 8	       | ExtraBoldFont  | Дополнительный жирный (Сверхжирный) Extra Bold (Ultra Bold)     |
| 9	       | BlackFont      | Чёрный (Густой) Black (Heavy)                                   |

Некоторые шрифты доступны только в нормальном или полужирном начертании. В этом случае значение данного свойства игнорируется

Получить значение данного свойства можно с помощью функции

	func GetTextWeight(view View, subviewID string) int

#### Свойство "text-shadow"

Свойство "text-shadow" позволяет задать тени для текста. Теней может быть несколько. Тень описывается
с помощью интерфейса ViewShadow (см. выше, раздел "Свойство 'shadow'"). Для тени текста используются только
свойства "color", "x-offset", "y-offset" и "blur". Свойства "inset" и "spread-radius" игнорируются (т.е. их
задание не является ошибкой, просто никакого влияния на тень текста они не имеют).

Для создания ViewShadow для тени текста используются функции:

	func NewTextShadow(offsetX, offsetY, blurRadius SizeUnit, color Color) ViewShadow
	func NewShadowWithParams(params Params) ViewShadow

Функция NewShadowWithParams используется когда в качестве параметров необходимо использовать
константы. Например:

	shadow := NewShadowWithParams(rui.Params{
		rui.ColorProperty : "@shadowColor",
		rui.BlurRadius : 8.0,
	})

В качестве значения свойству "text-shadow" может быть присвоено ViewShadow, массив ViewShadow,
текстовое представление ViewShadow (см. выше, раздел "Свойство 'shadow'").

Получить значение данного свойства можно с помощью функции

	func GetTextShadows(view View, subviewID string) []ViewShadow

Если тень не задана, то данная функция вернет пустой массив

#### Свойство "text-align"

Свойство "text-align" (константа TextAlign) - свойство типа int устанавливает выравнивание текста. Допустимые значения:

| Значение | Константа    | Имя       | Значение                     |
|:--------:|--------------|-----------|------------------------------|
| 0	       | LeftAlign    | "left"    | Выравнивание по левому краю  |
| 1        | RightAlign   | "right"   | Выравнивание по правому краю |
| 2        | CenterAlign  | "center"  | Выравнивание по центру       |
| 3        | JustifyAlign | "justify" | Выравнивание по ширине       |

Получить значение данного свойства можно с помощью функции

	func GetTextAlign(view View, subviewID string) int

#### Свойство "text-indent"

Свойство "text-indent" (TextIndent) - свойство типа SizeUnit определяет размер отступа (пустого места) перед первой строкой текста.

Получить значение данного свойства можно с помощью функции

	func GetTextIndent(view View, subviewID string) SizeUnit
	
#### Свойство "letter-spacing"

Свойство "letter-spacing" (LetterSpacing) - свойство типа SizeUnit определяет межбуквенное расстояние в тексте.
Значение может быть отрицательным, но при этом могут быть ограничения, зависящие от конкретной реализации.
Агент пользователя может не увеличивать или уменьшать межбуквенное расстояние для выравнивания текста.

Получить значение данного свойства можно с помощью функции

	func GetLetterSpacing(view View, subviewID string) SizeUnit

#### Свойство "word-spacing"

Свойство "word-spacing" (константа WordSpacing) - свойство типа SizeUnit определяет длину пробела между словами.
Если величина задана в процентах, то она определяет дополнительный интервал как процент от предварительной ширины символа.
В остальных случаях она определяет дополнительный интервал в дополнение к внутреннему интервалу между словами, определяемому шрифтом.

Получить значение данного свойства можно с помощью функции

	func GetWordSpacing(view View, subviewID string) SizeUnit

#### Свойство "line-height"

Свойство "line-height" (константа LineHeight) - свойство типа SizeUnit устанавливает величину пространства между строками.

Получить значение данного свойства можно с помощью функции

	func GetLineHeight(view View, subviewID string) SizeUnit

#### Свойство "text-transform"

Свойство "text-transform" (константа TextTransform) - свойство типа int определяет регистр символов. Допустимые значения:

| Значение | Константа               | Преобразование регистра                 |
|:--------:|-------------------------|-----------------------------------------|
| 0        | NoneTextTransform       | Оригинальный регистр символов           |
| 1	       | CapitalizeTextTransform | Каждое слово начинается с большой буквы |
| 2	       | LowerCaseTextTransform  | Все символы строчные                    |
| 3	       | UpperCaseTextTransform  | Все символы заглавные                   |

Получить значение данного свойства можно с помощью функции

	func GetTextTransform(view View, subviewID string) int

#### Свойство "text-direction"

Свойство "text-direction" (константа TextDirection) - свойство типа int определяет направление вывода текста. Допустимые значения:

| Значение | Константа               | Направление вывода текста                                                |
|:--------:|-------------------------|--------------------------------------------------------------------------|
| 0        | SystemTextDirection     | Системное направление. Определяется языком операционной системы.         |
| 1	       | LeftToRightDirection    | Слева направо. Используется для английского и большинства других языков. |
| 2	       | RightToLeftDirection    | Справа налево. Используется для иврит, арабский и некоторых других.      |

Получить значение данного свойства можно с помощью функции

	func GetTextDirection(view View, subviewID string) int

#### Свойство "writing-mode"

Свойство "writing-mode" (константа WritingMode) - свойство типа int определяет как располагаются строки текста
вертикально или горизонтально, а также направление в котором выводятся строки.
Возможны следующие значения:

| Значение | Константа             | Значение                                                           |
|:--------:|-----------------------|--------------------------------------------------------------------|
| 0        | HorizontalTopToBottom | Горизонтальные строки выводятся сверху сниз. Значение по умолчанию |
| 1        | HorizontalBottomToTop | Горизонтальные строки выводятся снизу вверх.                       |
| 2        | VerticalRightToLeft   | Вертикальные строки выводятся справа налево.                       |
| 3        | VerticalLeftToRight   | Вертикальные строки выводятся слева направо.                       |

Получить значение данного свойства можно с помощью функции

	func GetWritingMode(view View, subviewID string) int

#### Свойство "vertical-text-orientation"

Свойство "vertical-text-orientation" (константа VerticalTextOrientation) - свойство типа int используется, только
если "writing-mode" установлено в VerticalRightToLeft (2) или VerticalLeftToRight (3) и определяет положение
символов вертикальной строки. Возможны следующие значения:

| Значение | Константа               | Значение                                                           |
|:--------:|-------------------------|--------------------------------------------------------------------|
| 0        | MixedTextOrientation    | Символы повернуты на 90 по часовой стрелке. Значение по умолчанию. |
| 1        | UprightTextOrientation  | Символы расположены нормально (вертикально).                       |

Получить значение данного свойства можно с помощью функции

	func GetVerticalTextOrientation(view View, subviewID string) int

### Свойства трансформации

Данные свойства используются для трансформации (наклон, масштабирование и т.п.) содержимого View.

#### Свойство "perspective"

Свойство "perspective" (константа Perspective) определяет расстояние между плоскостью z = 0 и пользователем
для того чтобы придать 3D-позиционируемому элементу эффект перспективы. Каждый трансформируемый элемент с z > 0
станет больше, с z < 0 соответственно меньше.

Элементы части которые находятся за пользователем, т.е. z-координата этих элементов больше чем значение  свойства perspective, не отрисовываются.

Точка схождения по умолчанию расположена в центре элемента, но её можно переместить используя свойства
"perspective-origin-x" и "perspective-origin-y".

Получить значение данного свойства можно с помощью функции

	func GetPerspective(view View, subviewID string) SizeUnit

#### Свойства "perspective-origin-x" и "perspective-origin-y"

Свойства "perspective-origin-x" и "perspective-origin-y" (константы PerspectiveOriginX и PerspectiveOriginY)
типа SizeUnit определяют позицию, с которой смотрит зритель. Она используется свойством "perspective" как точка схода.

По умолчанию свойства "perspective-origin-x" и "perspective-origin-y" имеют значение 50%, т.е. указывают на центр View.

Получить значение данных свойств можно с помощью функции

	func GetPerspectiveOrigin(view View, subviewID string) (SizeUnit, SizeUnit)

#### Свойство "backface-visibility"

Свойство "backface-visibility" (константа BackfaceVisible) типа bool определяет, видно ли заднюю грань элемента,
когда он повёрнут к пользователю.

Задняя поверхность элемента является зеркальным отражением его передней поверхности. Однако невидимая в 2D,
задняя грань может быть видимой, когда преобразование вызывает вращение элемента в 3D пространстве.
(Это свойство не влияет на 2D-преобразования, которые не имеют перспективы.)

Получить значение данного свойства можно с помощью функции

	func GetBackfaceVisible(view View, subviewID string) bool

#### Свойства "origin-x", "origin-y" и "origin-z"

Свойства "origin-x", "origin-y" и "origin-z" (константа OriginX, OriginY и OriginZ) типа SizeUnit устанавливают
исходную точку для преобразований элемента.

Исходная точка преобразования - это точка, вокруг которой происходит преобразование. Например, вращение.

Свойство "origin-z" игнорируется если не установлено свойство "perspective".

Получить значение данных свойств можно с помощью функции

	func GetOrigin(view View, subviewID string) (SizeUnit, SizeUnit, SizeUnit)

#### Свойства "translate-x", "translate-y" и "translate-z"

Свойства "translate-x", "translate-y" и "translate-z" (константа TranslateX, TranslateY и TranslateZ) типа SizeUnit
позволяют задать смещение содержимого View.

Свойство "translate-z" игнорируется если не установлено свойство "perspective".

Получить значение данных свойств можно с помощью функции

	func GetTranslate(view View, subviewID string) (SizeUnit, SizeUnit, SizeUnit)

#### Свойства "scale-x", "scale-y" и "scale-z"

Свойства "scale-x", "scale-y" и "scale-z" (константа ScaleX, ScaleY и ScaleZ) типа float64 устанавливает
коэффициент масштабирования соответственно по осям x, y и z.
Исходный масштаб равен 1. Значение от 0 до 1 используется для уменьшения. Больше 1 - для увеличения.
Значения меньше или равное 0 являются недопустимыми (функция Set будет возвращать значение false)

Свойство "scale-z" игнорируется если не установлено свойство "perspective".

Получить значение данных свойств можно с помощью функции

	func GetScale(view View, subviewID string) (float64, float64, float64)

#### Свойства "rotate"

Свойство "rotate" (константа Rotate) типа AngleUnit задает угол поворота содержимого вокруг
вектора задаваемого свойствами "rotate-x", "rotate-y" и "rotate-z".

#### Свойства "rotate-x", "rotate-y" и "rotate-z"

Свойства "rotate-x", "rotate-y" и "rotate-z" (константа RotateX, RotateY и RotateZ) типа float64
задают вектор вокруг которого осуществляется вращение на угол заданный свойством "rotate".
Данный вектор проходит через точку заданную свойствами "origin-x", "origin-y" и "origin-z"

Свойство "rotate-z" игнорируется если не установлено свойство "perspective".

Получить значение данных свойств, а также свойства "rotate" можно с помощью функции

	func GetRotate(view View, subviewID string) (float64, float64, float64, AngleUnit)

#### Свойства "skew-x" и "skew-y"

Свойства "skew-x" и "skew-y" (константа SkewX и SkewY) типа AngleUnit задают скос (наклон) содержимого,
превращая тем самым его из прямоугольника в параллелограмм. Скос осуществляется вокруг точки,
задаваемой свойствами transform-origin-x и transform-origin-y.

Получить значение данных свойств можно с помощью функции

	func GetSkew(view View, subviewID string) (AngleUnit, AngleUnit)

### События клавиатуры

Для View получившего фокус ввода могут генерироваться два вида событий клавиатуры

| Событие          | Константа    | Описание               |
|------------------|--------------|------------------------|
| "key-down-event" | KeyDownEvent | Клавиша была нажата.   |
| "key-up-event"   | KeyUpEvent   | Клавиша была отпущена. |

Основной слушатель данных событий имеет следующий формат:
	
	func(View, KeyEvent)
	
где второй аргумент описывает параметры нажатых клавиш. Структура KeyEvent имеет следующие поля:

| Поле      | Тип    | Описание                                                                                                             |
|-----------|--------|----------------------------------------------------------------------------------------------------------------------|
| TimeStamp | uint64 | Время, когда событие было создано (в миллисекундах). Точка отсчета зависит от реализации браузера (ЭПОХА, запуск браузера и т.п.). |
| Key       | string | Значение клавиши, на которой возникло событие. Значение выдается с учетом текущего языка и регистра.                 |
| Code      | string | Код клавиши, представленного события. Значение не зависит от текущего языка и регистра.                              |
| Repeat    | bool   | Повторное нажатие: клавиша была нажата до тех пор, пока её ввод не начал автоматически повторяться.                  |
| CtrlKey   | bool   | Клавиша Ctrl была активна, когда возникло событие.                                                                   |
| ShiftKey  | bool   | Клавиша Shift была активна, когда возникло событие.                                                                  |
| AltKey    | bool   | Kлавиша Alt ( Option или ⌥ в OS X) была активна, когда возникло событие.                                             |
| MetaKey   | bool   | Kлавиша Meta (для Mac это клавиша ⌘ Command; для Windows - клавиша "Windows" ⊞) была активна, когда возникло событие.|

Можно также использовать слушателей следующих форматов:

* func(KeyEvent)
* func(View)
* func()

Получить списки слушателей событий клавиатуры можно с помощью функций:

	func GetKeyDownListeners(view View, subviewID string) []func(View, KeyEvent)
	func GetKeyUpListeners(view View, subviewID string) []func(View, KeyEvent)

### События фокуса

События фокуса возникают когда View получает или теряет фокус ввода. Соответственно возможны два событий:

| Событие            | Константа      | Описание                                        |
|--------------------|----------------|-------------------------------------------------|
| "focus-event"      | FocusEvent     | View получает фокус ввода (становится активным) |
| "lost-focus-event" | LostFocusEvent | View теряет фокус ввода (становится неактивным) |

Основной слушатель данных событий имеет следующий формат:

	func(View).

Можно также использовать слушателя следующего формата:

	func()

Получить списки слушателей событий фокуса можно с помощью функций:

	func GetFocusListeners(view View, subviewID string) []func(View)
	func GetLostFocusListeners(view View, subviewID string) []func(View)

### События мыши

Для View могут генерироваться несколько вида событий мыши

| Событие              | Константа        | Описание               |
|----------------------|------------------|------------------------|
| "mouse-down"         | MouseDown        | Клавиша мыши была нажата.                                    |
| "mouse-up"           | MouseUp          | Клавиша мыши была отпущена.                                  |
| "mouse-move"         | MouseMove        | Переместился курсор мыши                                     |
| "mouse-out"          | MouseOut         | Курсор мыши вышел за пределы View, или зашел в дочерной View |
| "mouse-over"         | MouseOver        | Курсор мыши зашел в пределы View                             |
| "click-event"        | ClickEvent       | Произошел клик мышкой                                        |
| "double-click-event" | DoubleClickEvent | Произошел двойной клик мышкой                                |
| "context-menu-event" | ContextMenuEvent | Нажета клавиша вызова контекстного меню (правая кнопка мыши) |

Основной слушатель данных событий имеет следующий формат:
	
	func(View, MouseEvent)
	
где второй аргумент описывает параметры нажатых клавиш. Структура MouseEvent имеет следующие поля:

| Поле      | Тип     | Описание                                                                                                             |
|-----------|---------|----------------------------------------------------------------------------------------------------------------------|
| TimeStamp | uint64  | Время, когда событие было создано (в миллисекундах). Точка отсчета зависит от реализации браузера (ЭПОХА, запуск браузера и т.п.). |
| Button    | int     | Номер кнопки мыши, нажатие на которую инициировало событие                                                           |
| Buttons   | int     | Битовая маска, показывающия какие кнопки мыши были нажаты в момент возникновения события                             |
| X         | float64 | Горизонтальная позиция мыши относительно начала координат View                                                       |
| Y         | float64 | Вертикальная позиция мыши относительно начала координат View                                                         |
| ClientX   | float64 | Горизонтальная позиция мыши относительно левого верхнего угда приложения                                             |
| ClientY   | float64 | Вертикальная позиция мыши относительно левого верхнего угда приложения                                               |
| ScreenX   | float64 | Горизонтальная позиция мыши относительно левого верхнего угда экрана                                                 |
| ScreenY   | float64 | Вертикальная позиция мыши относительно левого верхнего угда экрана                                                   |
| CtrlKey   | bool    | Клавиша Ctrl была активна, когда возникло событие.                                                                   |
| ShiftKey  | bool    | Клавиша Shift была активна, когда возникло событие.                                                                  |
| AltKey    | bool    | Kлавиша Alt ( Option или ⌥ в OS X) была активна, когда возникло событие.                                             |
| MetaKey   | bool    | Kлавиша Meta (для Mac это клавиша ⌘ Command; для Windows - клавиша "Windows" ⊞) была активна, когда возникло событие.|

Поле Button может принимать следующие значения

| Значение | Константа            | Описание                                                                         |
|:--------:|----------------------|----------------------------------------------------------------------------------|
| <0       |                      | Не нажата ни одна кнопка                                                         |
| 0        | PrimaryMouseButton   | Основная кнопка. Обычно левая кнопка мыши (может быть изменена в настройках ОС)  |
| 1        | AuxiliaryMouseButton | Вспомогательная кнопка. Колёсико или средняя кнопка мыши, если она есть          |
| 2        | SecondaryMouseButton | Вторичная кнопка. Обычно правая кнопка мыши (может быть изменена в настройках ОС)|
| 3        | MouseButton4         | Четвёртая кнопка мыши. Обычно кнопка браузера Назад                              |
| 4        | MouseButton5         | Пятая кнопка мыши. Обычно кнопка браузера Вперёд                                 |

Поле Button преставляет собой битовую маску объединающую (с помощью ИЛИ) следующие значения

| Значение | Константа          | Описание               |
|:--------:|--------------------|------------------------|
| 1        | PrimaryMouseMask   | Основная кнопка        |
| 2        | SecondaryMouseMask | Вторичная кнопка       |
| 4        | AuxiliaryMouseMask | Вспомогательная кнопка |
| 8        | MouseMask4         | Четвёртая кнопка       |
| 16       | MouseMask5         | Пятая кнопка           |

Можно также использовать слушателей следующих форматов:

* func(MouseEvent)
* func(View)
* func()

Получить списки слушателей событий мыши можно с помощью функций:

	func GetMouseDownListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseUpListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseMoveListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseOverListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetMouseOutListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetClickListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetDoubleClickListeners(view View, subviewID string) []func(View, MouseEvent)
	func GetContextMenuListeners(view View, subviewID string) []func(View, MouseEvent)

### События указателя

Указатель - это аппаратно-независимое представление устройств ввода (таких как мышь, перо
или точка контакта на сенсорной поверхности). Указатель может указывать на конкретную координату
(или набор координат) на контактной поверхности, например на экране.

Все указатели могут генерироваться несколько вида событий

| Событие          | Константа     | Описание                                                   |
|------------------|---------------|------------------------------------------------------------|
| "pointer-down"   | PointerDown   | Указатель был нажат.                                       |
| "pointer-up"     | PointerUp     | Указатель был отпущен.                                     |
| "pointer-move"   | PointerMove   | Указатель перемещен                                        |
| "pointer-cancel" | PointerCancel | События указателя прерваны.                                |
| "pointer-out"    | PointerOut    | Указатель вышел за пределы View, или зашел в дочерной View |
| "pointer-over"   | PointerOver   | Указатель зашел в пределы View                             |

Основной слушатель данных событий имеет следующий формат:
	
	func(View, PointerEvent)
	
где второй аргумент описывает параметры указателя. Структура PointerEvent расширяет структуру MouseEvent
и имеет следующие дополнительные поля:

| Поле               | Тип     | Описание                                                              |
|--------------------|---------|-----------------------------------------------------------------------|
| PointerID          | int     | Уникальный идентификатор указателя, вызвавшего событие.               |
| Width              | float64 | Ширина (величина по оси X) в пикселях контактной геометрии указателя. |
| Height             | float64 | Высота (величина по оси Y) в пикселях контактной геометрии указателя. |
| Pressure           | float64 | Нормализованное давление на входе указателя в диапазоне от 0 до 1, где 0 и 1 представляют минимальное и максимальное давление, которое аппаратное обеспечение способно обнаруживать, соответственно.|
| TangentialPressure | float64 | Нормализованное тангенциальное давление на входе указателя (также известное как давление в цилиндре или напряжение цилиндра) в диапазоне от -1 до 1, где 0 - нейтральное положение элемента управления.|
| TiltX              | float64 | Плоский угол (в градусах в диапазоне от -90 до 90) между плоскостью Y–Z и плоскостью, содержащей как ось указателя (например, стилуса), так и ось Y.|
| TiltY              | float64 | Плоский угол (в градусах в диапазоне от -90 до 90) между плоскостью X–Z и плоскостью, содержащей как ось указателя (например, стилуса), так и ось X.|
| Twist              | float64 | Вращение указателя (например, стилуса) по часовой стрелке вокруг своей главной оси в градусах со значением в диапазоне от 0 до 359.|
| PointerType        | string  | тип устройства, вызвавшего событие: "mouse", "pen", "touch" и т.п.     |
| IsPrimary          | bool    | указатель является первичным указателем этого типа.                    |

Можно также использовать слушателей следующих форматов:

* func(PointerEvent)
* func(View)
* func()

Получить списки слушателей событий указателя можно с помощью функций:

	func GetPointerDownListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerUpListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerMoveListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerCancelListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerOverListeners(view View, subviewID string) []func(View, PointerEvent)
	func GetPointerOutListeners(view View, subviewID string) []func(View, PointerEvent)

### Touch события

Данные события используются для отслеживания многоточечных касаний. Одиночные касания эмулируют события мыши.
Если у вас нет необходимости отслеживания многоточечных касаний, то проще использовать события мыши

| Событие        | Константа   | Описание                                      |
|----------------|-------------|-----------------------------------------------|
| "touch-start"  | TouchStart  | Произошло касание поверхности.                |
| "touch-end"    | TouchEnd    | Завершено касание поверхности.                |
| "touch-move"   | TouchMove   | Одно или несколько касаний изменили положение |
| "touch-cancel" | TouchCancel | Касание прервано.                             |

Основной слушатель данных событий имеет следующий формат:
	
	func(View, TouchEvent)

где второй аргумент описывает параметры касаний. Структура TouchEvent имеет следующие поля:

| Поле      | Тип     | Описание                                                                                                             |
|-----------|---------|----------------------------------------------------------------------------------------------------------------------|
| TimeStamp | uint64  | Время, когда событие было создано (в миллисекундах). Точка отсчета зависит от реализации браузера (ЭПОХА, запуск браузера и т.п.). |
| Touches   | []Touch | Массив структур Touch, каждая из которых описывает одно касание                                                      |
| CtrlKey   | bool    | Клавиша Ctrl была активна, когда возникло событие.                                                                   |
| ShiftKey  | bool    | Клавиша Shift была активна, когда возникло событие.                                                                  |
| AltKey    | bool    | Kлавиша Alt ( Option или ⌥ в OS X) была активна, когда возникло событие.                                             |
| MetaKey   | bool    | Kлавиша Meta (для Mac это клавиша ⌘ Command; для Windows - клавиша "Windows" ⊞) была активна, когда возникло событие.|

Структура Touch описывает одиночное касание и имеет следующие поля

| Поле          | Тип     | Описание                                                                                                                 |
|---------------|---------|--------------------------------------------------------------------------------------------------------------------------|
| Identifier    | int     | Уникальный идентификатор присваюваемый каждому касанию и не меняющийся до его завершения.                                |
| X             | float64 | Горизонтальная позиция мыши относительно начала координат View                                                           |
| Y             | float64 | Вертикальная позиция мыши относительно начала координат View                                                             |
| ClientX       | float64 | Горизонтальная позиция мыши относительно левого верхнего угда приложения                                                 |
| ClientY       | float64 | Вертикальная позиция мыши относительно левого верхнего угда приложения                                                   |
| ScreenX       | float64 | Горизонтальная позиция мыши относительно левого верхнего угда экрана                                                     |
| ScreenY       | float64 | Вертикальная позиция мыши относительно левого верхнего угда экрана                                                       |
| RadiusX       | float64 | x-радиус эллипса в пикселях, который наиболее точно ограничивает область контакта с экраном.                             |
| RadiusY       | float64 | y-радиус эллипса в пикселях, который наиболее точно ограничивает область контакта с экраном.                             |
| RotationAngle | float64 | Угол (в градусах), на который нужно повернуть по часовой стрелке эллипс, описываемый параметрами radiusX и radiusY, чтобы наиболее точно покрыть область контакта между пользователем и поверхностью. |
| Force         | float64 | Величина давления от 0,0 (без давления) до 1,0 (максимальное давление), которое пользователь прикладывает к поверхности. |

Можно также использовать слушателей следующих форматов:

* func(TouchEvent)
* func(View)
* func()

Получить списки слушателей событий касания можно с помощью функций:

	func GetTouchStartListeners(view View, subviewID string) []func(View, TouchEvent)
	func GetTouchEndListeners(view View, subviewID string) []func(View, TouchEvent)
	func GetTouchMoveListeners(view View, subviewID string) []func(View, TouchEvent)
	func GetTouchCancelListeners(view View, subviewID string) []func(View, TouchEvent)

### Событие "resize-event"

Событие "resize-event" (константа ResizeEvent) вызывается когда View меняет свои положение и/или размеры.
Основной слушатель данных событий имеет следующий формат:
	
	func(View, Frame)

где структура объявлена как

	type Frame struct {
		Left, Top, Width, Height float64
	}

Соотвественно элементы Frame содержат следующие данные
* Left - новое смещение в пикселях по горизонтали относительно родительского View (левая позиция);
* Top - новое смещение в пикселях по вертикали относительно родительского View (верхняя позиция)
* Width - новая ширина видимой части View в пикселях;
* Height - новая высота видимой части View в пикселях.

Можно также использовать слушателей следующих форматов:

* func(Frame)
* func(View)
* func()

Получить список слушателей данного события можно с помощью функции:

	func GetResizeListeners(view View, subviewID string) []func(View, Frame)

Текущие положение и размеры видимой части View можно получить с помощью функции интерфейса View:

	Frame() Frame

или глобальной функции

	func GetViewFrame(view View, subviewID string) Frame

### Событие прокрутки

Событие "scroll-event" (константа ScrollEvent) созникает при прокрутке содержимого View.
Основной слушатель данных событий имеет следующий формат:
	
	func(View, Frame)

где элементы Frame содержат следующие данные
* Left - новое смещение видимой области по горизонтали (в пикселях);
* Top - новое смещение видимой области по вертикали (в пикселях);
* Width - общая ширина View в пикселях;
* Height - общая высота View в пикселях.

Можно также использовать слушателей следующих форматов:

* func(Frame)
* func(View)
* func()

Получить список слушателей данного события можно с помощью функции:

	func GetScrollListeners(view View) []func(View, Frame)

Текущие положение видимой области и общие размеры View можно получить с помощью функции интерфейса View:

	Scroll() Frame

или глобальной функции

	func GetViewScroll(view View, subviewID string) Frame

Для программной прокрутки могут использоваться следующие глобальные функции

	func ScrollViewTo(view View, subviewID string, x, y float64)
	func ScrollViewToStart(view View, subviewID string)
	func ScrollViewToEnd(view View, subviewID string)

которые прокручивают view, соответственно, в заданную позицию, начало и конец

## ViewsContainer

Интерфейс ViewsContainer, реализующий View, описывает контейнер содержащий несколько
дочерних элементов интерфейса (View). ViewsContainer является базовым для других контейнеров
(ListLayout, GridLayout, StackLayout и т.д.) и самостоятельно не используется.

Помимо всех свойств View данный элемент имеет всего одно дополнительное свойство "content"

### "content"

Свойство "content" (константа Сontent) определяет массив дочерних View. Функция Get интерфейса
для данного свойства всегда возвращает []View.

В качестве значения свойства "content" могут быть переданы следующие 5 типов данных:

* View - преобразуется во []View, содержащий один элемент;

* []View - nil-элементы запрещены, если массив будет содержать nil, то свойство не будет
установлено, а функция Set вернет false и в лог запишется сообщение об ошибке;

* string - если строка является текстовым представление View, то создается соответствующий View,
иначе создается TextView, которому в качестве текста передается данная строка.
Далее создается []View, содержащий полученный View;

* []string - каждый элемент массива преобразуется во View как описано в предыдущем пункте;

* []interface{} - данный массив должен содержать только View и string. Каждый string-элемент
преобразуется во View, как описано выше. Если массив будет содержать недопустимае значения,
то свойство "content" не будет установлено, а функция Set вернет false  и в лог запишется сообщение об ошибке.

Поучить значение свойства "content" можно с помощи функции интерфейса ViewsContainer

	Views() []View

Для редактирования свойства "content" можно использовать следующие функции интерфейса ViewsContainer:

	Append(view View)

Данная функция добавляет аргумент в конец списка View

	Insert(view View, index uint)

Данная функция вставляет аргумент в заданную позицию списка View. Если index больше длины
списка, то View добавляется в конец списка. Если index меньше 0, то в начало списка.

	RemoveView(index uint) View

Данная функция удаляет View из заданной позиции и возвращает его. Если index указывает за
границы списка, то ничего не удаляется, а функция возвращает nil.

## ListLayout

ListLayout является контейнером, реализующим интерфейс ViewsContainer. Для его создания используется функция

	func NewListLayout(session Session, params Params) ListLayout

Элементы в данном контейнере располагаются в виде списка. Расположением дочерних элементов можно управлять. 
Для этого ListLayout имеет ряд свойств

### "orientation"

Свойство "orientation" (константа Orientation) типа int задает то как дочерние элементы будут
располагаться друг относительно друга. Свойство может принимать следующие значения:

| Значение | Константа             | Расположение                                               |
|:--------:|-----------------------|------------------------------------------------------------|
| 0        | TopDownOrientation    | Дочерние элементы располагаются в столбец сверху вниз.     |
| 1        | StartToEndOrientation | Дочерние элементы располагаются в строку с начала в конец. |
| 2        | BottomUpOrientation   | Дочерние элементы располагаются в столбец снизу вверх.     |
| 3        | EndToStartOrientation | Дочерние элементы располагаются в строку с конца в начала. |

Положение начала и конца для StartToEndOrientation и EndToStartOrientation зависит от значения
свойства "text-direction". Для языков с письмом справа налево (арабский, иврит) начало находится
справа, для остальных языков - слева.

### "wrap"

Свойство "wrap" (константа Wrap) типа int определяет расположения элементов в случае достижения
границы контейнера. Возможны три варианта:

* WrapOff (0) - колонка/строка элементов продолжается и выходит за границы видимой области.
	
* WrapOn (1) - начинается новая колонка/строка элементов. Новая колонка располагается по направлению
к концу (о положении начала и конца см. выше), новая строка - снизу.
	
* WrapReverse (2) - начинается новая колонка/строка элементов. Новая колонка располагается по направлению
к началу (о положении начала и конца см. выше), новая строка - сверху.

### "vertical-align"

Свойство "vertical-align" (константа VerticalAlign) типа int устанавливает вертикальное
выравнивание элементов в контейнере. Допустимые значения:

| Значение | Константа    | Имя       | Значение                      |
|:--------:|--------------|-----------|-------------------------------|
| 0	       | TopAlign     | "top"     | Выравнивание по верхнему краю |
| 1        | BottomAlign  | "bottom"  | Выравнивание по нижнему краю  |
| 2        | CenterAlign  | "center"  | Выравнивание по центру        |
| 3        | StretchAlign | "stretch" | Выравнивание по высоте        |

### "horizontal-align"

Свойство "horizontal-align" (константа HorizontalAlign) типа int устанавливет
горизонтальное выравнивание элементов в списке. Допустимые значения:

| Значение | Константа    | Имя       | Значение                     |
|:--------:|--------------|-----------|------------------------------|
| 0	       | LeftAlign    | "left"    | Выравнивание по левому краю  |
| 1        | RightAlign   | "right"   | Выравнивание по правому краю |
| 2        | CenterAlign  | "center"  | Выравнивание по центру       |
| 3        | StretchAlign | "stretch" | Выравнивание по ширине       |

## GridLayout

GridLayout является контейнером, реализующим интерфейс ViewsContainer. Для его создания используется функция

	func NewGridLayout(session Session, params Params) GridLayout

Пространство контейнера данного контейнера разбито на ячейки в виде таблицы.
Все дочерние элементы располагаются в ячейках таблицы. Ячейка адресуется по номеру строки и
столбца. Номера строк и столбцов начинаются с 0.

### "column" и "row"

Расположение View внутри GridLayout определяется с помощью свойств "column" и "row".
Данные свойства должны устанавливаться для каждого из дочерних View.
Дочерний View может занимать несколько ячеек внутри GridLayout. При этом они могут
перекрываться.

В качестве значения "column" и "row" можно установить:
* целое число большее или равное 0;
* текстовое представление целого числа большего или равного 0 или константу;
* структуру Range, задающую диапазон строк/столбцов:

	type Range struct {
		First, Last int
	}

где First - номер первого столбца/строки, Last - номер последнего столбца/строки;
* строка вида "<номер первого столбца/строки>:<номер последнего столбца/строки>", являющуюся
текстовым представление структуры Range

Пример

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

В данном примере view1 занимает в нулевой строке столбцы 1 и 2, а view1 занимает
в нулевом стобце строки 0, 1 и 2.

### "cell-width" и "cell-height"

По умолчанию размеры ячеек вычисляются на основе размеров помещенных в них дочерних View.
Свойства "cell-width" и "cell-height" (константы CellWidth и CellHeight) позволяют установить
фиксированную ширину и высоту ячеек независимо от размеров дочерних элементов.
Данные свойства имеют тип []SizeUnit. Каждый элемент массива определяет размер соответствующего
столбца или строки.

Данным свойствам могут быть присвоены следующие типы данных:

* SizeUnit или текстовое представление SizeUnit (или SizeUnit константа). В этом случае
соотвествующие размеры всех ячеек устанавливаются одинаковыми;

* []SizeUnit;

* string содержащая текстовые представления SizeUnit (или SizeUnit константы) разделенные запятой;

* []string. Каждый элемент должен быть текстовым представлением SizeUnit (или SizeUnit константой)

* []interface{}. Каждый элемент должен или иметь тип SizeUnit или быть текстовым
представлением SizeUnit (или SizeUnit константой)

Если количество элементов в свойствах "cell-width" и "cell-height" меньше, чем используемое число
столбцов и строк, то недостающие элементы устанавливаются в Auto.

В значениях свойств "cell-width" и "cell-height" может использоваться SizeUnit тип SizeInFraction.
Этот тип означает 1 часть. Часть вычисляется так: из размера контейнера вычитается размер всех
ячеек имеющих тип не SizeInFraction, а затем оставшийся размер делится на количество частей.
Значение SizeUnit типа SizeInFraction может быть как целым, так и дробным.

### "grid-row-gap" и "grid-column-gap"

Свойства "grid-row-gap" и "grid-column-gap" (константы GridRowGap и GridColumnGap) типа SizeUnit позволяют
установить соответственно расстояния между строками и столбцами контейнера. Значение по умолчанию 0.

### "cell-vertical-align"

Свойство "cell-vertical-align" (константа CellVerticalAlign) типа int устанавливает вертикальное
выравнивание дочерних элементов внутри занимаемой ячейки. Допустимые значения:

| Значение | Константа    | Имя       | Значение                      |
|:--------:|--------------|-----------|-------------------------------|
| 0	       | TopAlign     | "top"     | Выравнивание по верхнему краю |
| 1        | BottomAlign  | "bottom"  | Выравнивание по нижнему краю  |
| 2        | CenterAlign  | "center"  | Выравнивание по центру        |
| 3        | StretchAlign | "stretch" | Растягивание на всю высоту    |

Значение по умолчанию StretchAlign (3)

### "cell-horizontal-align"

Свойство "cell-horizontal-align" (константа CellHorizontalAlign) типа int устанавливет
горизонтальное выравнивание дочерних элементов внутри занимаемой ячейки. Допустимые значения:

| Значение | Константа    | Имя       | Значение                     |
|:--------:|--------------|-----------|------------------------------|
| 0	       | LeftAlign    | "left"    | Выравнивание по левому краю  |
| 1        | RightAlign   | "right"   | Выравнивание по правому краю |
| 2        | CenterAlign  | "center"  | Выравнивание по центру       |
| 3        | StretchAlign | "stretch" | Растягивание на всю ширину   |

Значение по умолчанию StretchAlign (3)

## ColumnLayout

ColumnLayout является контейнером, реализующим интерфейс ViewsContainer. Все дочерние View
располагаются в виде вертикального списка выровненные по левому или правому краю и разбитого
на несколько колонок. Выравнивание зависит от свойства "text-direction".

Для создания ColumnLayout используется функция

	func NewColumnLayout(session Session, params Params) ColumnLayout

### Свойство "column-count"

Свойство "column-count" (константа ColumnCount) типа int устанавливет количество колонок.

Если данное свойство равно 0 и не задано свойство "column-width", то разбитие на колонки
не выполняется, а контейнер прокручивается вниз.

Если значение данного свойства больше 0, то список разбивается на колонки. Высота колонки
равна высоте ColumnLayout, а ширина вычисляется как ширина ColumnLayout делённая на
"column-count". Каждая следующая колонка располагается в зависимости от свойства
"text-direction" справа или слева от предыдущей, а контейнер прокручивается по горизонтали.

Получить значение данного свойства можно с помощью функции

	func GetColumnCount(view View, subviewID string) int

### Свойство "column-width"

Свойство "column-width" (константа ColumnWidth) типа SizeUnit используется только если
"column-count" равно 0 и устанавливет ширину колонки.

ВАЖНО! В качестве значения "column-width" нельзя использовать проценты (т.е. если вы зададите
значение в процентах, то это проигнорируется системой)

Получить значение данного свойства можно с помощью функции

	func GetColumnWidth(view View, subviewID string) SizeUnit

### Свойство "column-gap"

Свойство "column-gap" (константа ColumnGap) типа SizeUnit устанавливает ширину разрыва между колонками.

Получить значение данного свойства можно с помощью функции

	func GetColumnGap(view View, subviewID string) SizeUnit

### Свойство "column-separator"

Свойство "column-separator" (константа ColumnSeparator) позволяет задать линию которая будет
рисоваться в разрывах колонок. Линия рамки описывается тремя атрибутами: стиль линии, толщина и цвет.

Значение свойства "column-separator" хранится в виде интерфейса ColumnSeparatorProperty,
реализующего интерфейс Properties (см. выше). ColumnSeparatorProperty может содержать следующие свойства:

| Свойство | Константа     | Тип      | Описание         |
|----------|---------------|----------|------------------|
| "style"  | Style         | int      | Стиль линии      |
| "width"  | Width         | SizeUnit | Толщина линии    |
| "color"  | ColorProperty | Color    | Цвет линии       |

Стиль линии может принимать следующие значения:

| Значение | Константа  | Имя      | Описание                 |
|:--------:|------------|----------|--------------------------|
| 0        | NoneLine   | "none"   | Нет рамки                |
| 1        | SolidLine  | "solid"  | Сплошная линия           |
| 2        | DashedLine | "dashed" | Пунктирная линия         |
| 3        | DottedLine | "dotted" | Линия состоящая из точек |
| 4        | DoubleLine | "double" | Двойная сплошная линия   |

Все другие значения стиля игнорируются.

Для создания интерфейса ColumnSeparatorProperty используется функция

	func NewColumnSeparator(params Params) ColumnSeparatorProperty

Интерфейс ColumnSeparatorProperty может быть преобразован в структуру ViewBorder с помощью
функции ViewBorder. При преобразовании все текстовые константы заменяются реальными значениями.
ViewBorder описана как

	type ViewBorder struct {
		Style int
		Color Color
		Width SizeUnit
	}

Структура ViewBorder может быть передана в качестве параметра функции Set при установке значения
свойства "column-separator". При этом ViewBorder преобразуется в ColumnSeparatorProperty.
Поэтому при чтении свойства функцией Get будет возвращен интерфейс ColumnSeparatorProperty,
а не структура ViewBorder. Получить структуру ViewBorders без дополнительных преобразований можно
с помощью глобальной функции

	func GetColumnSeparator(view View, subviewID string) ViewBorder

Вы также можете устанавливать отдельные атрибуты линии использую функцию Set интерфейса View.
Для этого используются следующие свойства

| Свойство                 | Константа            | Тип      | Описание      |
|--------------------------|----------------------|----------|---------------|
| "column-separator-style" | ColumnSeparatorStyle | int      | Стиль линии   |
| "column-separator-width" | ColumnSeparatorWidth | SizeUnit | Толщина линии |
| "column-separator-color" | ColumnSeparatorColor | Color    | Цвет линии    |

Например

	view.Set(rui.ColumnSeparatorStyle, rui.SolidBorder)
	view.Set(rui.ColumnSeparatorWidth, rui.Px(1))
	view.Set(rui.ColumnSeparatorColor, rui.Black)

эквивалентно

	view.Set(rui.ColumnSeparator, ColumnSeparatorProperty(rui.Params{
		rui.Style: rui.SolidBorder,
		rui.Width: rui.Px(1),
		rui.ColorProperty: rui.Black,
	}))

### Свойство "avoid-break"

При формировании колонок ColumnLayout может разрывать некоторые типы View, так что начало
будет в конце одной колонки, а окончание в следующей. Например, разрывается TextView,
заголовок картинки и сама картинки и т.д.

Свойство "avoid-break" (константа AvoidBreak) типа bool позволяет избежать этого эффекта.
Необходимо установить для View, который нельзя разрывать, данное свойство со значением "true".
Соответственно значение "false" данного свойства позволяет разрывать View.
Значение по умолчанию "false".

Получить значение данного свойства можно с помощью функции

	func GetAvoidBreak(view View, subviewID string) bool

## StackLayout

StackLayout является контейнером, реализующим интерфейс ViewsContainer. Все дочерние View
располагаются друг над другом и каждый занимает все пространство контейнера. В каждый момент времени
доступен только один дочерний View (текущий).

Для создания StackLayout используется функция

	func NewStackLayout(session Session, params Params) StackLayout

Помимо свойств Append, Insert, RemoveView и свойства "content" интерфейса ViewsContainer
контейнер StackLayout имеет еще две функции интерфейса для управления дочерними View: Push и Pop

	Push(view View, animation int, onPushFinished func())

Данная функция добавляет новый View в контейнер и делает его текущим. Она похожа на Append,
но в отличие от нее дабавление выполняется с использованием эффекта анимации. Вид анимации
задается вторым аргументом и может принимать следующие значения:

| Значение | Константа           | Анимация                    |
|:--------:|---------------------|-----------------------------|
| 0	       | DefaultAnimation    | Анимация по умолчанию. Для функции Push это EndToStartAnimation, для Pop - StartToEndAnimation  |
| 1        | StartToEndAnimation | Анимация из начала в конец. Начало и конец определются направлением вывода текста               |
| 2        | EndToStartAnimation | Анимация из конца в начало. |
| 3        | TopDownAnimation    | Анимация сверху вниз.       |
| 4        | BottomUpAnimation   | Анимация снизу вверх.       |

Третий аргумент onPushFinished - функция вызываемая по окончании анимации. Может быть nil.

	Pop(animation int, onPopFinished func(View)) bool

Данная функция удаляет текущий View из контейнера используя анимацию.
Второй аргумент onPopFinished - функция вызываемая по окончании анимации. Может быть nil.
Функция вернёт false если StackLayout пуст и true если текущий элемени был удален.

 Получить текущий (видимый) View можно с помощью функции интерфейса

	Peek() View

Для того чтобы сделать любой дочерний View текущим (видимым) используются функции интерфейса:

	MoveToFront(view View) bool
	MoveToFrontByID(viewID string) bool

Данная функция вернет true в случае успеха и false если дочерний View или View с таким id не существует и в
лог будет записано сообщение об ошибке.

## AbsoluteLayout

AbsoluteLayout является контейнером, реализующим интерфейс ViewsContainer. Дочерние View
могут располагатся в произвольных позициях пространства контейнера.

Для создания AbsoluteLayout используется функция

	func NewAbsoluteLayout(session Session, params Params) AbsoluteLayout

Дочерние View позиционируюся с помощью свойств типа SizeUnit: "left", "right", "top" и
"bottom" (соответственно константы Left, Right, Top и Bottom). Можно задавать любые из
этих свойств для дочернего View. Если ни "left" ни "right" не заданы, то дочерний View
будет прижат к левому краю контейнера. Если ни "top" ни "bottom" не заданы, то дочерний View
будет прижат к верхнему краю контейнера.

## DetailsView

DetailsView является контейнером, реализующим интерфейс ViewsContainer.
Для создания DetailsView используется функция

	func NewDetailsView(session Session, params Params) DetailsView

Помимо дочерних View данный контейнер имеет свойство "summary" (константа Summary).
В качестве значения свойства "summary" может быть или View или строка текста.

DetailsView может находиться в одном из двух состояний:

* отображается только содержимое свойства "summary". Дочерние View скрыты и не занимают место на экране

* отображается сначала содержимое свойства "summary", а ниже дочерние View.
Размещение дочерних View, аналогично ColumnLayout с "column-count" равным 0.

DetailsView переключается между состояниями по клику по "summary".

Для принудительного переключения состояний DetailsView используется bool свойство
"expanded" (константа Expanded). Соответственно значение "true" показывает дочерние
View, "false" - скрывает.

Получить значение свойства "expanded" можно с помощью функции

	func IsDetailsExpanded(view View, subviewID string) bool

а значение свойства "summary" можно получить с помощью функции

	func GetDetailsSummary(view View, subviewID string) View

## Resizable

Resizable является контейнером в который можно поместить только один View. Resizable позволяет
интерактивно менять размеры вложенного View.

Вокруг влоденного View создается рамка, потянув за которую можно менять размеры.

Resizable не реализует интерфейс ViewsContainer. Для управлением вложенным View используется
только свойство Content. Данному свойству может быть присвоено значение типа View или
строка текста. Во втором случае создается TextView.

Рамка вокруг вложенного View может быть как со всех сторон, так и только с отдельных.
Для задания сторон рамки используется свойство "side" (константа Side) типа int.
Оно может принимать следующие значения:

| Значение | Константа    | Имя      | Сторона рамки                      |
|:--------:|--------------|----------|------------------------------------|
| 1	       | TopSide      | "top"    | Верхняя                            |
| 2        | RightSide    | "right"  | Правая                             |
| 4        | BottomSide   | "bottom" | Нижняя                             |
| 8        | LeftSide     | "left"   | Левая                              |
| 15       | AllSides     | "all"    | Все стороны. Значение по умолчанию |

Кроме этих значений может также использоваться or-комбинация TopSide, RightSide, BottomSide и LeftSide.
AllSides определено как

	AllSides = TopSide | RightSide | BottomSide | LeftSide

Для установки ширины рамки используется SizeUnit свойство "resize-border-width" (константа ResizeBorderWidth).
Значение по умолчанию для "resize-border-width" равно 4px.

## TextView

Элемент TextView расширяющий интерфейс View предназначен для вывода текста.

Для создания TextView используется функция:

	func NewTextView(session Session, params Params) TextView

Выводимый текст задается string свойством "text" (константа Text).
Помимо метода Get значение свойства "text" может быть получено с помощью функции

	func GetText(view View, subviewID string) string

TextView наследует от View все свойства параметров текста ("font-name", "text-size", "text-color" и т.д.).
Кроме них добавляется еще один "text-overflow" (константа TextOverflow). Он определяет как обрезается
текст если он выходит за границы. Данное свойство типа int может принимать следующие значения

| Значение | Константа            | Имя        | Обрезка текста                             |
|:--------:|----------------------|------------|--------------------------------------------|
| 0	       | TextOverflowClip     | "clip"     | Текст обрезается по границе (по умолчанию) |
| 1        | TextOverflowEllipsis | "ellipsis" | В конце видимой части текста выводится '…' |

## EditView

Элемент EditView является редактором теста и расширяет интерфейс View.

Для создания EditView используется функция:

	func NewEditView(session Session, params Params) EditView

Возможно несколько вариантов редактируемого текста. Тип редактируемого текста устанвливается
с помощью int свойства "edit-view-type" (константа EditViewType).
Данное свойство может принимать следующие значения:

| Значение | Константа      | Имя         | Тип редактора                                       |
|:--------:|----------------|-------------|-----------------------------------------------------|
| 0	       | SingleLineText | "text"      | Однострочный редактор текста. Значение по умолчанию |
| 1	       | PasswordText   | "password"  | Радактор пароля. Текст скрывается звездочками       |
| 2	       | EmailText      | "email"     | Редактор для ввода одиночного e-mail                |
| 3	       | EmailsText     | "emails"    | Редактор для ввода нескольких e-mail                |
| 4	       | URLText        | "url"       | Редактор для ввода интернет адреса                  |
| 5	       | PhoneText      | "phone"     | Редактор для ввода телефонного номера               |
| 6	       | MultiLineText  | "multiline" | Многострочный редактор текста                       |

Для упрощения текста программы можно использовать свойства "type" (константа Type) вместо "edit-view-type".
Эти имена свойств синонимы. Но при описании стиля "type" использовать нельзя

Для установки/получения редактируемого текста используется string свойство "text" (константа Text)

Максимальная длина редактируемого текста устанавливается с помощью int свойства "max-length"
(константа MaxLength).

Вы можете ограничить вводимый текст с помощью регулярного выражения. Для этого используется
string свойство "edit-view-pattern" (константа EditViewPattern). Вместо "edit-view-pattern"
можно использовать синоним "pattern" (константа Pattern), за исключением описания стиля.

Для запрещения редактирования текста используется bool свойство "readonly" (константа ReadOnly).

Для включения/выключения встроенной проверки орфографии используется bool свойство "spellcheck"
(константа Spellcheck). Проверка орфографии можно включить только если тип редактора установлен
в SingleLineText или MultiLineText.

Для редактора можно установить подсказку которая будет показываться пока редактор пуст.
Для этого используется string свойство "hint" (константа Hint).

Для многострочного редактора может быть включен режим автоматического переноса. Для
этого используется bool свойство "wrap" (константа Wrap). Если "wrap" выключен (значение по умолчанию),
то используется горизонтальная прокрутка. Если включен, то по достижении границы EditView
текст переносится на новую строку.

Для получения значений свойств EditView могут использоваться следующие функции:

	func GetText(view View, subviewID string) string
	func GetHint(view View, subviewID string) string
	func GetMaxLength(view View, subviewID string) int
	func GetEditViewType(view View, subviewID string) int
	func GetEditViewPattern(view View, subviewID string) string
	func IsReadOnly(view View, subviewID string) bool
	func IsEditViewWrap(view View, subviewID string) bool
	func IsSpellcheck(view View, subviewID string) bool


Для отслеживания изменения текста используется событие "edit-text-changed" (константа
EditTextChangedEvent). Основной слушатель события имеет следующий формат:

	func(EditView, string)

где второй аргумент это новое значение текста

Получить текущий список слушателей изменения текста можно с помощью функции

	func GetTextChangedListeners(view View, subviewID string) []func(EditView, string)

## NumberPicker

Элемент NumberPicker расширяет интерфейс View и предназначен для ввода чисел.

Для создания NumberPicker используется функция:

	func NewNumberPicker(session Session, params Params) NumberPicker

NumberPicker может работать в двух режимах: редактор текста и слайдер.
Режим устанавливает int свойство "date-picker-type" (константа NumberPickerType).
Свойство "date-picker-type" может принимать следующие значения:

| Значение | Константа    | Имя      | Тип редактора                          |
|:--------:|--------------|----------|----------------------------------------|
| 0	       | NumberEditor | "editor" | Редактор текста. Значение по умолчанию |
| 1	       | NumberSlider | "slider" | Слайдер                                |

Установить/прочитать текущее значение можно с помощью свойства "date-picker-value"
(константа NumberPickerValue). В качестве значения свойству "date-picker-value" могут быть переданы:

* float64
* float32
* int
* int8…int64
* uint
* uint8…uint64
* текстовое представление любых из выше перечисленых типов

Все эти типы приводятся к float64. Соответственно функция Get всегда возвращает float64 значение.
Прочитано значение свойства "date-picker-value" может быть также с помощью функции:

	func GetNumberPickerValue(view View, subviewID string) float64

На вводимые значения могут быть наложены ограничения. Для этого используются следующие свойства:

| Свойство             | Константа        | Ограничение |
|----------------------|------------------|------------------------|
| "date-picker-min"  | NumberPickerMin  | Минимальное значение   |
| "date-picker-max"  | NumberPickerMax  | Максимальное значение  |
| "date-picker-step" | NumberPickerStep | Шаг изменения значения |

Присвоины данным свойствам могут те же типы значений, что и "date-picker-value".

По умолчанию, в случае если "date-picker-type" равно NumberSlider, минимальное значение равно 0,
максимальное - 1. Если же "date-picker-type" равно NumberEditor то вводимые числа, по умолчанию,
ограничены лишь диапазоном значений float64.

Прочитать значения данных свойств можно с помощью функций:

	func GetNumberPickerMinMax(view View, subviewID string) (float64, float64)
	func GetNumberPickerStep(view View, subviewID string) float64

Для отслеживания изменения вводимого значения используется событие "date-changed" (константа
NumberChangedEvent).  Основной слушатель события имеет следующий формат:

	func(picker NumberPicker, newValue float64)

где второй аргумент это новое значение

Получить текущий список слушателей изменения значения можно с помощью функции

	func GetNumberChangedListeners(view View, subviewID string) []func(NumberPicker, float64)

## DatePicker

Элемент DatePicker расширяет интерфейс View и предназначен для ввода дат.

Для создания DatePicker используется функция:

	func NewDatePicker(session Session, params Params) DatePicker

Установить/прочитать текущее значение можно с помощью свойства "date-picker-value"
(константа DatePickerValue). В качестве значения свойству "date-picker-value" могут быть переданы:

* time.Time
* константа
* текст, который может быть преобразован в time.Time функцией
	func time.Parse(layout string, value string) (time.Time, error)

Текст  преобразуется в time.Time. Соответственно функция Get всегда возвращает time.Time значение.
Прочитано значение свойства "date-picker-value" может быть также с помощью функции:

	func GetDatePickerValue(view View, subviewID string) time.Time

На вводимые даты могут быть наложены ограничения. Для этого используются следующие свойства:

| Свойство           | Константа      | Тип данных | Ограничение                |
|--------------------|----------------|------------|----------------------------|
| "date-picker-min"  | DatePickerMin  | time.Time  | Минимальное значение даты  |
| "date-picker-max"  | DatePickerMax  | time.Time  | Максимальное значение даты |
| "date-picker-step" | DatePickerStep | int        | Шаг изменения даты в днях  |

Прочитать значения данных свойств можно с помощью функций:

	func GetDatePickerMin(view View, subviewID string) (time.Time, bool)
	func GetDatePickerMax(view View, subviewID string) (time.Time, bool)
	func GetDatePickerStep(view View, subviewID string) int

Для отслеживания изменения вводимого значения используется событие "date-changed" (константа
DateChangedEvent).  Основной слушатель события имеет следующий формат:

	func(picker DatePicker, newDate time.Time)

где второй аргумент это новое значение даты

Получить текущий список слушателей изменения даты можно с помощью функции

	func GetDateChangedListeners(view View, subviewID string) []func(DatePicker, time.Time)

## TimePicker

Элемент TimePicker расширяет интерфейс View и предназначен для ввода времени.

Для создания TimePicker используется функция:

	func NewTimePicker(session Session, params Params) TimePicker

Установить/прочитать текущее значение можно с помощью свойства "time-picker-value"
(константа TimePickerValue). В качестве значения свойству "time-picker-value" могут быть переданы:

* time.Time
* константа
* текст, который может быть преобразован в time.Time функцией
	func time.Parse(layout string, value string) (time.Time, error)

Текст преобразуется в time.Time. Соответственно функция Get всегда возвращает time.Time значение.
Прочитано значение свойства "time-picker-value" может быть также с помощью функции:

	func GetTimePickerValue(view View, subviewID string) time.Time

На вводимое время могут быть наложены ограничения. Для этого используются следующие свойства:

| Свойство           | Константа      | Тип данных | Ограничение                      |
|--------------------|----------------|------------|----------------------------------|
| "time-picker-min"  | TimePickerMin  | time.Time  | Минимальное значение времени     |
| "time-picker-max"  | TimePickerMax  | time.Time  | Максимальное значение времени    |
| "time-picker-step" | TimePickerStep | int        | Шаг изменения времени в секундах |

Прочитать значения данных свойств можно с помощью функций:

	func GetTimePickerMin(view View, subviewID string) (time.Time, bool)
	func GetTimePickerMax(view View, subviewID string) (time.Time, bool)
	func GetTimePickerStep(view View, subviewID string) int

Для отслеживания изменения вводимого значения используется событие "time-changed" (константа
TimeChangedEvent).  Основной слушатель события имеет следующий формат:

	func(picker TimePicker, newTime time.Time)

где второй аргумент это новое значение времени

Получить текущий список слушателей изменения даты можно с помощью функции

	func GetTimeChangedListeners(view View, subviewID string) []func(TimePicker, time.Time)

## ColorPicker

Элемент ColorPicker расширяет интерфейс View и предназначен для выбора цвета в формате RGB без альфа канала.

Для создания ColorPicker используется функция:

	func NewColorPicker(session Session, params Params) ColorPicker

Установить/прочитать текущее значение можно с помощью свойства "color-picker-value"
(константа ColorPickerValue). В качестве значения свойству "color-picker-value" могут быть переданы:

* Color
* текстовое представление Color
* константа

Прочитано значение свойства "color-picker-value" может быть также с помощью функции:

	func GetColorPickerValue(view View, subviewID string) Color

Для отслеживания изменения выбранного цвета используется событие "color-changed" (константа
ColorChangedEvent).  Основной слушатель события имеет следующий формат:

	func(picker ColorPicker, newColor Color)

где второй аргумент это новое значение цвета

Получить текущий список слушателей изменения даты можно с помощью функции

	func GetColorChangedListeners(view View, subviewID string) []func(ColorPicker, Color)

## FilePicker

Элемент FilePicker расширяет интерфейс View и предназначен для выбора одного или нескольких файлов.

Для создания FilePicker используется функция:

	func NewFilePicker(session Session, params Params) FilePicker

Булевское свойство "multiple" (константа Multiple) используется для установки режима выбора нескольких файлов.
Значение "true" включает режим выбора нескольких файлов, "false" включает режим выбора одиночного файл.
Значение по умолчанию "false".

Вы можете ограничить выбор только определенными типами файлов. Для этого используется свойство "accept" (константа Accept).
Данному свойству присваивается список разрешенных расширений файлов и/или mime-типов. Значение можно задавать или в виде 
строки (элементы при этом разделяются запятыми) или в виде массива строк. Примеры

	rui.Set(view, "myFilePicker", rui.Accept, "png, jpg, jpeg")
	rui.Set(view, "myFilePicker", rui.Accept, []string{"png", "jpg", "jpeg"})
	rui.Set(view, "myFilePicker", rui.Accept, "image/*")
	
Для доступа к выбранным файлам используются две функции интерфейса FilePicker:

	Files() []FileInfo
	LoadFile(file FileInfo, result func(FileInfo, []byte))

а также соответствующие им глобальные функции

	func GetFilePickerFiles(view View, subviewID string) []FileInfo
	func LoadFilePickerFile(view View, subviewID string, file FileInfo, result func(FileInfo, []byte))

Функции Files/GetFilePickerFiles возвращают список выбранных файлов в виде среза структур FileInfo. Структура FileInfo объявлена как

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

FileInfo содержит только информацию о файле, но не сам файл. Функция LoadFile/LoadFilePickerFile позволяет загрузить 
содержимое одного из выбранных файлов. Функция LoadFile асинхронная. После загрузки содержимое выбранного файла 
передается функции-аргументу LoadFile. Пример

	if filePicker := rui.FilePickerByID(view, "myFilePicker"); filePicker != nil {
		if files := filePicker.Files(); len(files) > 0 {
			filePicker.LoadFile(files[0], func(file rui.FileInfo, data []byte) {
				if data != nil {
					// ... 
				}
			})
		}
	}

эквивалентно

	if files := rui.GetFilePickerFiles(view, "myFilePicker"); len(files) > 0 {
		rui.LoadFilePickerFile(view, "myFilePicker", files[0], func(file rui.FileInfo, data []byte) {
			if data != nil {
				// ... 
			}
		})
	}

Если во время загрузки файла произойдет ошибка, то значение data передаваемое в функцию результата будет равно nil,
а описание ошибки будет записано в лог

Для отслеживания изменения списка выбранных файлов используется событие "file-selected-event" 
(константа FileSelectedEvent). Основной слушатель события имеет следующий формат:

	func(picker FilePicker, files []FileInfo))

где второй аргумент это новое значение списка выбранных файлов.

Получить текущий список слушателей изменения списка файлов можно с помощью функции

	func GetFileSelectedListeners(view View, subviewID string) []func(FilePicker, []FileInfo)

## DropDownList

Элемент DropDownList расширяет интерфейс View и предназначен для выбора значения из выпадающего списка.

Для создания DropDownList используется функция:

	func NewDropDownList(session Session, params Params) DropDownList

Список возможных значений задается с помощью свойства "items" (константа Items).
В качестве значения свойству "items" могут быть переданы следующие типы данных

* []string
* []fmt.Stringer
* []interface{} содержащий в качестве элементов только: string, fmt.Stringer, bool, rune,
float32, float64, int, int8…int64, uint, uint8…uint64.

Все эти типы данных преопразуются в []string и присваиваются свойству "items".
Прочитать значение свойства "items" можно с помощью функции

	func GetDropDownItems(view View, subviewID string) []string

Выбранное значение определяется int свойством "current" (константа Current). Значение по умолчанию 0.
Прочитать значение данного свойства можно с помощью функции

	func GetDropDownCurrent(view View, subviewID string) int

Для отслеживания изменения свойства "current" используется событие "drop-down-event" (константа
DropDownEvent). Основной слушатель события имеет следующий формат:

	func(list DropDownList, newCurrent int)

где второй аргумент это индекс выбранного элемента

Получить текущий список слушателей изменения даты можно с помощью функции

	func GetDropDownListeners(view View, subviewID string) []func(DropDownList, int)

## ProgressBar

Элемент DropDownList расширяет интерфейс View и предназначен для отображение прогресса в виде
заполняемой полоски.

Для создания ProgressBar используется функция:

	func NewProgressBar(session Session, params Params) ProgressBar

ProgressBar имеет два свойства типа float64:
* "progress-max" (константа ProgressBarMax) - максимальное значение (по умолчанию 1);
* "progress-value" (константа ProgressBarValue) - текущее значение (по умолчанию 0).

Минимальное всегда 0.
В качестве значений этим свойствам может быть присвоино кроме float64 также float32, int,
int8…int64, uint, uint8…uint64

Прочитать значение данных свойств можно с помощью функций

	func GetProgressBarMax(view View, subviewID string) float64
	func GetProgressBarValue(view View, subviewID string) float64

## Button

Элемент Button реализует нажимаемую кнопку. Это CustomView (о нем ниже) на базе ListLayout и,
соответсвенно, обладает всеми свойствами ListLayout. Но в отличие от ListLayout может получать
фокус ввода.

Контент, по умолчанию, выровнен по центру.

Для создания Button используется функция:

	func NewButton(session Session, params Params) Button

## ListView

Элемент ListView реализует список.
Для создания ListView используется функция:

	func NewListView(session Session, params Params) ListView

### Свойство "items"

Элементы списка задаются с помощью свойства "items" (константа Items). Основным значением
свойства "items" является интерфейс ListAdapter:

	type ListAdapter interface {
		ListSize() int
		ListItem(index int, session Session) View
		IsListItemEnabled(index int) bool
	}

Соостветственно функции этого интерфейса должны возвращать количество элементов, View i-го элемента и
статус i-го элемента разрешен/запрещен.

Вы можете реализовать этот интерфейс сами или воспользоваться вспомогательными функциями:

	func NewTextListAdapter(items []string, params Params) ListAdapter
	func NewViewListAdapter(items []View) ListAdapter

NewTextListAdapter создает адаптер из массива строк, второй аргумент это параметры TextView используемого
для отобрыжения текста. NewViewListAdapter создает адаптер из массива View.

Cвойствy "items" могут быть присвоены следующие типы данных:

* ListAdapter;
* []View при присваиваниии преобразуется в ListAdapter с помощью функции NewViewListAdapter;
* []string при присваиваниии преобразуется в ListAdapter с помощью функции NewTextListAdapter;
* []interface{} который может содержать элементы типа View, string, fmt.Stringer, bool, rune,
float32, float64, int, int8…int64, uint, uint8…uint64. При присваиваниии все типы кроме
View и string преобразуется в string, далее все string в TextView и из получившегося массива View
с помощью функции NewViewListAdapter получается ListAdapter.

Если элементы списка меняются в ходе работы, то после изменения необходимо вызывать или функцию
ReloadListViewData() интерфейса ListView или глобальную функцию ReloadListViewData(view View, subviewID string).
Данные функции обновляют отображаемые элементы списка.

### Свойство "orientation"

Элементы списка могут располагаться как вертикально (колонками), так и горизонтально (строками).
Свойство "orientation" (константа Orientation) типа int задает то как элементы списка будут
располагаться друг относительно друга. Свойство может принимать следующие значения:

| Значение | Константа             | Расположение                                      |
|:--------:|-----------------------|---------------------------------------------------|
| 0        | TopDownOrientation    | Элементы располагаются в столбец сверху вниз.     |
| 1        | StartToEndOrientation | Элементы располагаются в строку с начала в конец. |
| 2        | BottomUpOrientation   | Элементы располагаются в столбец снизу вверх.     |
| 3        | EndToStartOrientation | Элементы располагаются в строку с конца в начала. |

Положение начала и конца для StartToEndOrientation и EndToStartOrientation зависит от значения
свойства "text-direction". Для языков с письмом справа налево (арабский, иврит) начало находится
справа, для остальных языков - слева.

Получить значение данного свойства можно с помощью функции

	func GetListOrientation(view View, subviewID string) int

### Свойство "wrap"

Свойство "wrap" (константа Wrap) типа int определяет расположения элементов в случае достижения
границы контейнера. Возможны три варианта:

* WrapOff (0) - колонка/строка элементов продолжается и выходит за границы видимой области.
	
* WrapOn (1) - начинается новая колонка/строка элементов. Новая колонка располагается по направлению
к концу (о положении начала и конца см. выше), новая строка - снизу.
	
* WrapReverse (2) - начинается новая колонка/строка элементов. Новая колонка располагается по направлению
к началу (о положении начала и конца см. выше), новая строка - сверху.

Получить значение данного свойства можно с помощью функции

	func GetListWrap(view View, subviewID string) int

### Свойства "item-width" и "item-height"

По умолчанию высота и ширина элементов списка вычисляется на основе их содержимого.
Это приводит к тому что элементы вертикального списка могут иметь разную высоту, а элементы
горизонтального - разную ширину.

Вы можете установить фиксированную высоту и ширину элемента списка. Для этого используются SizeUnit
свойства "item-width" и "item-height"

Получить значения данных свойств можно с помощью функций

	func GetListItemWidth(view View, subviewID string) SizeUnit
	func GetListItemHeight(view View, subviewID string) SizeUnit

### Свойство "item-vertical-align"

Свойство "item-vertical-align" (константа ItemVerticalAlign) типа int устанавливает вертикальное
выравнивание содержимого элементов списка. Допустимые значения:

| Значение | Константа    | Имя       | Значение                      |
|:--------:|--------------|-----------|-------------------------------|
| 0	       | TopAlign     | "top"     | Выравнивание по верхнему краю |
| 1        | BottomAlign  | "bottom"  | Выравнивание по нижнему краю  |
| 2        | CenterAlign  | "center"  | Выравнивание по центру        |
| 3        | StretchAlign | "stretch" | Выравнивание по высоте        |

Получить значение данного свойства можно с помощью функции

	func GetListItemVerticalAlign(view View, subviewID string) int

### Свойство "item-horizontal-align"

Свойство "item-horizontal-align" (константа ItemHorizontalAlign) типа int устанавливет
горизонтальное выравнивание содержимого элементов списка. Допустимые значения:

| Значение | Константа    | Имя       | Значение                     |
|:--------:|--------------|-----------|------------------------------|
| 0	       | LeftAlign    | "left"    | Выравнивание по левому краю  |
| 1        | RightAlign   | "right"   | Выравнивание по правому краю |
| 2        | CenterAlign  | "center"  | Выравнивание по центру       |
| 3        | StretchAlign | "stretch" | Выравнивание по ширине       |

Получить значение данного свойства можно с помощью функции

	GetListItemHorizontalAlign(view View, subviewID string) int

### Свойство "current"

ListView позволяет выбирать пункты списка имеющие статус "разрешен" (см. ListAdapter).
Элемент может быть выбран как интерактивно, так и программно. Для этого используется
int свойство "current" (константа Current). Значение "current" меньше 0 означает что
не выбран ни один пункт

Получить значение данного свойства можно с помощью функции

	func GetListViewCurrent(view View, subviewID string) int

### Свойства "list-item-style", "current-style" и "current-inactive-style"

Данные три свойства отвечают за стиль фона и свойства текста каждого элемента списка.

| Свойство                 | Константа            | Стиль                                           |
|--------------------------|----------------------|-------------------------------------------------|
| "list-item-style"        | ListItemStyle        | Стиль невыбранного элемента                     |
| "current-style"          | CurrentStyle         | Стиль выбранного элемента. ListView в фокусе    |
| "current-inactive-style" | CurrentInactiveStyle | Стиль выбранного элемента. ListView не в фокусе |

### Свойства "checkbox", "checked", "checkbox-horizontal-align" и "checkbox-vertical-align"

Свойство "current" позволяет выбрать один пункт списка.
Свойства "checkbox" позволяет добавить к каждому элементу списка чекбокс с помощью которого
можно выбрать несколько элементов списка. Свойство "checkbox" (константа ItemCheckbox) имеет тип int
и может принимать следующие значения

| Значение | Константа        | Имя        | Вид чекбокса                                     |
|:--------:|------------------|------------|--------------------------------------------------|
| 0	       | NoneCheckbox     | "none"     | Нет чекбокса. Значение по умолчанию              |
| 1	       | SingleCheckbox   | "single"   | ◉ Чекбокс позволяющий пометить только один пункт |
| 2	       | MultipleCheckbox | "multiple" | ☑ Чекбокс позволяющий пометить несколько пунктов |


Получить значение данного свойства можно с помощью функции

	func GetListViewCheckbox(view View, subviewID string) int

Получить/установить список помеченных пунктов можно с помощью свойства "checked" (константа Checked).
Данное свойство имеет тип []int и хранит индексы помеченных элементов.
Получить значение данного свойства можно с помощью функции

	func GetListViewCheckedItems(view View, subviewID string) []int

Проверить помечен ли конкретный элемент можно с помощью функции

	func IsListViewCheckedItem(view View, subviewID string, index int) bool

По умолчанию чекбокс расположен в верхнем левом углу элемента. Изменить его положение можно
с помощью int свойств "checkbox-horizontal-align" и "checkbox-vertical-align" (константы
CheckboxHorizontalAlign и CheckboxVerticalAlign)

Свойство "checkbox-horizontal-align" (константа СheckboxHorizontalAlign) может принимать следующие значения:

| Значение | Константа    | Имя      | Расположение чекбокса                           |
|:--------:|--------------|----------|-------------------------------------------------|
| 0	       | LeftAlign    | "left"   | У левого края. Контент справа                   |
| 1        | RightAlign   | "right"  | У правого края. Контент слева                   |
| 2        | CenterAlign  | "center" | По центру по горизонтали. Контент ниже или выше |

Свойство "checkbox-vertical-align" (константа CheckboxVerticalAlign) может принимать следующие значения:

| Значение | Константа    | Имя      | Значение                      |
|:--------:|--------------|----------|-------------------------------|
| 0	       | TopAlign     | "top"    | Выравнивание по верхнему краю |
| 1        | BottomAlign  | "bottom" | Выравнивание по нижнему краю  |
| 2        | CenterAlign  | "center" | Выравнивание по центру        |

Особый случай когда и "checkbox-horizontal-align" и "checkbox-vertical-align" равны CenterAlign (2).
В этом случае чекбокс расположен по центру по горизонтали, контент ниже

Получить значения свойств можно "checkbox-horizontal-align" и "checkbox-vertical-align" с помощью функций

	func GetListViewCheckboxHorizontalAlign(view View, subviewID string) int
	func GetListViewCheckboxVerticalAlign(view View, subviewID string) int

### События ListView

Для ListView есть три характерных события

* "list-item-clicked" (константа ListItemClickedEvent) - возникает когда пользователь кликает по элементу
списка. Основной слушатель данного события имеет следующий формат: func(ListView, int). Где второй
аргумент это индекс элемента.

* "list-item-selected" (константа ListItemSelectedEvent) - возникает когда пользователь элемент списка
становится выбраным. Основной слушатель данного события имеет следующий формат: func(ListView, int).
Где второй аргумент это индекс элемента.

* "list-item-checked" (константа ListItemCheckedEvent) - возникает когда пользователь ставит/снимает
пометку чекбокса элемента списка. Основной слушатель данного события имеет следующий формат: func(ListView, []int).
Где второй аргумент это массив индексов помеченных элементов.

Получить списки слушателей данных событий можно с помощью функций:

	func GetListItemClickedListeners(view View, subviewID string) []func(ListView, int)
	func GetListItemSelectedListeners(view View, subviewID string) []func(ListView, int)
	func GetListItemCheckedListeners(view View, subviewID string) []func(ListView, []int)

## TableView

Элемент TableView реализует таблицу. Для создания TableView используется функция:

	func NewTableView(session Session, params Params) TableView

### Свойство "content"

Свойство "content" определяет содержимое таблицы. Для описания содержимого необходимо реализовать
интерфейс TableAdapter объявленный как

	type TableAdapter interface {
		RowCount() int
		ColumnCount() int
		Cell(row, column int) interface{}
	}

где функции RowCount() и ColumnCount() должны возврацать количество строк и столбцов в таблице;
Cell(row, column int) возвращает содержимое ячейки таблицы. Функция Cell() может возвращать
элементы следующих типов:

* string
* rune
* float32, float64
* целочисленные значения: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
* bool
* rui.Color
* rui.View
* fmt.Stringer

Свойству "content" можно также присваивать следующие типы данных

* TableAdapter
* [][]interface{}
* [][]string

[][]interface{} и [][]string при присвоении преобразуются к TableAdapter.

### Свойство "cell-style"

Свойство "cell-style" (константа CellStyle) предназначено для настройки оформления ячейки таблицы. Данному свойству
может быть присвоено только реализация интерфейса TableCellStyle.

	type TableCellStyle interface {
		CellStyle(row, column int) Params
	}

Данный интерфейс содержит только одну функцию CellStyle, которая возвращает параметры оформления
заданной ячейки таблицы. Можно использовать любые свойства интерфейса View. Например

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		if row == 0 {
			return rui.Params {
				rui.BackgroundColor: rui.Gray,
				rui.Italic:          true,
			}
		}
		return nil
	}

Если не надо менять оформление какой-то ячейки, то для нее можно вернуть nil.

#### Свойства "row-span" и "column-span"

Помимо свойств интерфейса View, функцией CellStyle могут возвращаться еще два свойства типа int:
"row-span" (константа RowSpan) и "column-span" (константа ColumnSpan).
Данные свойства используются для объединения ячеек таблицы.

Свойство "row-span" указавает сколько ячеек надо объединить по вертикали,
а "column-span" - по горизонтали. Например

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		if row == 0 && column == 0 {
			return rui.Params { rui.RowSpan: 2 }
		}
		if row == 0 && column == 1 {
			return rui.Params { rui.ColumnSpan: 2 }
		}
		return nil
	}

В этом случае таблица будет иметь следующий вид

|------|----------------|
|      |                |
|      |-------|--------|
|      |       |        |
|------|-------|--------|

Если в качестве значения свойства "content" используется [][]interface{}, то для объединения
ячеек используются пустые структуры

	type VerticalTableJoin struct {
	}
	type HorizontalTableJoin struct {
	}

Данные структуры присоединяют ячейку, соответсвенно, к верхней/левой. Описание приведенной выше таблицы будет
иметь следующий вид

	content := [][]interface{} {
		{"", "", rui.HorizontalTableJoin{}},
		{rui.VerticalTableJoin{}, "", ""},
	}

### Свойство "row-style"

Свойство "row-style" (константа RowStyle) предназначено для настройки оформления строки таблицы.
Данному свойству может быть присвоены или реализация интерфейса TableRowStyle или []Params.
TableRowStyle объявлена как

	type TableRowStyle interface {
		RowStyle(row int) Params
	}

Функция RowStyle возвращает параметры применяемые ко всей строке таблицы.
Свойство "row-style" имеет более низкий приоритет по сравнению со сойством "cell-style",
т.е. свойства заданные в "cell-style" будут использоваться вместо заданных в "row-style"

### Свойство "column-style"

Свойство "column-style" (константа ColumnStyle) предназначено для настройки оформления столбца таблицы.
Данному свойству может быть присвоены или реализация интерфейса TableColumnStyle или []Params.
TableColumnStyle объявлена как

	type TableColumnStyle interface {
		ColumnStyle(column int) Params
	}

Функция ColumnStyle возвращает параметры применяемые ко всему столбцу таблицы.
Свойство "column-style" имеет более низкий приоритет по сравнению со сойствами "cell-style" и "row-style".

### Свойства "head-height" и "head-style"

Таблица может иметь "шапку".
Свойство "head-height" (константа HeadHeight) типа int указывает сколько первых строк таблицы образуют "шапку".
Свойство "head-style" (константа HeadStyle) задает стиль шапки. Свойству "head-style" может быть
присвоено, значение типа:

* string - имя стиля;
* []Params - перечисление свойств "шапки".

### Свойства "foot-height" и "foot-style"

Таблица может иметь в конце финализирующие строки (например строка "Итого").
Свойство "foot-height" (константа FootHeight) типа int указывает количество этих финализирующих строк.
Свойство "foot-style" (константа FootStyle) задает их стиль. Значения свойства "foot-style" аналогичны свойству "head-style".

### Свойство "cell-padding"

Свойство "cell-padding" (константа CellPadding) типа SizeUnit задает отступы от границ ячейки до
контента. Данное свойство эквивалентно

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		return rui.Params { rui.Padding: <my padding> }
	}

И введено для удобства, чтобы не надо было писать адаптер для задания отступов.
Свойство "cell-padding" имеет более низкий приоритет по сравнению со сойством "cell-style".

"cell-padding" может также использоваться при задании параметров в свойствах
"row-style", "column-style", "foot-style" и "head-style"

### Свойство "cell-border"

Свойство "cell-border" (константа CellBorder) задает памку для всех ячеек таблицы.
Данное свойство эквивалентно

	func (style *myTableCellStyle) CellStyle(row, column int) rui.Params {
		return rui.Params { rui.Border: <my padding> }
	}

И введено для удобства, чтобы не надо было писать адаптер для рамки.
Свойство "cell-border" имеет более низкий приоритет по сравнению со сойством "cell-style".

"cell-border" может также использоваться при задании параметров в свойствах
"row-style", "column-style", "foot-style" и "head-style"

### Свойство "table-vertical-align"

Свойство "table-vertical-align" (константа TableVerticalAlign) типа int задает вертикальное выравание
данных внутри ячейки таблицы. Допустимые значения:

| Значение | Константа     | Имя        | Значение                      |
|:--------:|---------------|------------|-------------------------------|
| 0	       | TopAlign      | "top"      | Выравнивание по верхнему краю |
| 1        | BottomAlign   | "bottom"   | Выравнивание по нижнему краю  |
| 2        | CenterAlign   | "center"   | Выравнивание по центру        |
| 3, 4     | BaselineAlign | "baseline" | Выравнивание по базовой линии |

Для горизонтального выравнивания используется свойство "text-align"

## Пользовательский View

Пользовательский View должен реализовывать интерфейс CustomView, который в свою очередь
расширяет интерфейсы ViewsContainer и View. Пользовательский View создается на основе другого,
который назавается Super View.

Для упрощение задачи уже имеется базовая реализация CustomView в виде структуры CustomViewData.

Создание пользовательского View рассмотрим на примере встроенного элемента Buttom:

1) объявляем интерфейс Button, как расширяющий CustomView, и структуру buttonData как расширяющую CustomViewData

	type Button interface {
		rui.CustomView
	}

	type buttonData struct {
		rui.CustomViewData
	}

2) реализуем функцию CreateSuperView

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

3) если надо, то переопределяем методы интерфейса CustomView, для Button это
функция Focusable() (так как кнопка может получать фокус, а ListLayout не получает)

	func (button *buttonData) Focusable() bool {
		return true
	}

4) пишем функцию для создания Button:

	func NewButton(session rui.Session, params rui.Params) Button {
		button := new(buttonData)
		rui.InitCustomView(button, "Button", session, params)
		return button
	}

При создании CustomView обязательным является вызов функции InitCustomView.
Данная функция инициализирует структуру CustomViewData. Первым аргументом
является указатель на инициализируемую структуру, вторым - имя присвоенное
вашему View, третьим - сессия и четвертым - параметры

5) регистрируем элемент. Регистрацию рекомендуется осуществлять в методе init пакета

	rui.RegisterViewCreator("Button", func(session rui.Session) rui.View {
		return NewButton(session, nil)
	})

Все! Новый элемент готов

## CanvasView

CanvasView это область в которой вы можете рисовать. Для создания CanvasView используется функция:

	func NewCanvasView(session Session, params Params) CanvasView

CanvasView имеет всего одно дополнительное свойство: "draw-function" (константа DrawFunction).
С помощью данного свойства задается функция рисования имеющая следующее описание

	func(Canvas)

где Canvas это контекст рисования с помощью которого осуществляется рисование

Интерфей Canvas сожержит ряд функция для настройки стилей, текста и непосредственно самого рисования.

Все координаты и размеры задаются только в пикселях, поэтому при рисовании SizeUnit не используется.
Везде используется float64

### Настройка стиля линий

Для настройки цвета линий используются следующие функции интерфейса Canvas:

* SetSolidColorStrokeStyle(color Color) - линия будет рисоваться сплошным цветом

* SetLinearGradientStrokeStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) -
линия будет рисоваться с помощью линейного градиента. Начальная точка градиента задается с помощью x0, y0 и color0,
конечная - x1, y1 и color1. Массив []GradientPoint задает промежуточные точки градиента. Если промежуточных точек нет,
то в качестве последнего параметра можно передать nil

* SetRadialGradientStrokeStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) -
линия будет рисоваться с помощью радиального градиента. x0, y0, r0, color0 - координаты центра, радиус и цвет начальной окружности.
x1, y1, r1, color1 - координаты центра, радиус и цвет конечной окружности. Массив []GradientPoint задает промежуточные точки градиента

Структура GradientPoint описана как

	type GradientPoint struct {
		Offset float64
		Color Color
	}

где Offset - значение в диапазоне от 0 до 1 задает относительное положение промежуточной точки, Color - цвет этой точки.

Толщина линии в пикселях задается функцией

	SetLineWidth(width float64)

Вид концов линии задается с помощью функции

	SetLineCap(cap int)

где cap может принимать следующие значения

| Значение | Константа | Вид                                                                              |
|:--------:|-----------|----------------------------------------------------------------------------------|
| 0        | ButtCap   | Окончании линии обрезано в конечной точке. Значение по умолчанию.                |
| 1        | RoundCap  | Окончании линии скруглено. Центр окружности находится в конечной точке.          |
| 2        | SquareCap | В конец линии добавляется прямоугольник с шириной равной половине толщины линии. |

Форма, используемая для соединения двух отрезков линии в месте их пересечения, задается функцией

	SetLineJoin(join int)

где join может принимать следующие значения

| Значение | Константа | Вид    |
|:--------:|-----------|--------|
| 0        | MiterJoin | Сегменты соединяются путем удлинения их внешних краев для соединения в одной точке с эффектом заполнения дополнительной области в форме ромба. |
| 1        | RoundJoin | Закругляет углы фигуры, заполняя дополнительный сектор диском с центром в общей конечной точке соединенных сегментов. Радиус этих закругленных углов равен ширине линии. |
| 2        | BevelJoin | Заполняет дополнительную треугольную область между общей конечной точкой соединенных сегментов и отдельными внешними прямоугольными углами каждого сегмента. |

По умолчанию рисуется сплошная линия. Если необходимо нарисовать прерывистую линию, то
необходимо сначала задать шаблон с помощью функции

	SetLineDash(dash []float64, offset float64)

где dash []float64 задает шаблон линиии в виде чередования длин отрезков и пропусков. Второй аргумент -
смещение шаблона относительно начала линии.

Пример

	canvas.SetLineDash([]float64{16, 8, 4, 8}, 0)

Линия рисуется следующим образом: отрезок длиной 16 пикселей, затеп пропуск длиной 8 пикселей,
отрезок длиной 4 пикселя, затеп пропуск длиной 8 пикселей, затем снова отрезок длиной 16 пикселей и т.д.

### Настройка стиля заливки

Для настройки стиля заливки используются следующие функции интерфейса Canvas:

* SetSolidColorFillStyle(color Color) - фигура будет заливаться сплошным цветом

* SetLinearGradientFillStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) -
фигура будет заливаться с помощью линейного градиента. Начальная точка градиента задается с помощью x0, y0 и color0,
конечная - x1, y1 и color1. Массив []GradientPoint задает промежуточные точки градиента. Если промежуточных точек нет,
то в качестве последнего параметра можно передать nil

* SetRadialGradientFillStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) -
фигура будет заливаться с помощью радиального градиента. x0, y0, r0, color0 - координаты центра, радиус и цвет начальной окружности.
x1, y1, r1, color1 - координаты центра, радиус и цвет конечной окружности. Массив []GradientPoint задает промежуточные точки градиента

### Рисование геометрических фигур

#### Прямоугольник

Для рисования прямоугольников могут использоваться три функции:

	FillRect(x, y, width, height float64)
	StrokeRect(x, y, width, height float64)
	FillAndStrokeRect(x, y, width, height float64)

FillRect рисует закрашеный прямоугольник.

StrokeRect рисует контур прямоугольника.

FillAndStrokeRect рисует контур и закрашивает внутриности.

#### Прямоугольник с закругленными углами

Аналогично прямоуголику есть три функции рисования

	FillRoundedRect(x, y, width, height, r float64)
	StrokeRoundedRect(x, y, width, height, r float64)
	FillAndStrokeRoundedRect(x, y, width, height, r float64)

где r это радиус скругления

#### Эллипс

Для рисования эллипсов также могут использоваться три функции:

	FillEllipse(x, y, radiusX, radiusY, rotation float64)
	StrokeEllipse(x, y, radiusX, radiusY, rotation float64)
	FillAndStrokeEllipse(x, y, radiusX, radiusY, rotation float64)

где x, y - центр эллипса, radiusX, radiusY - радиусы эллипса по оси X и Y,
rotation - угол поворота элилпса относительно центра в радианах

#### Path

Интерфейс Path позволяет описать сложную фигуру. Создается Path с помощью функции NewPath().

После создания вы должны описать фигуру. Для этого могут использоваться следующие функции интерфейса:

* MoveTo(x, y float64) - переместить текущую точке в заданные координаты;

* LineTo(x, y float64) - добавить линию из текущей точки в заданную;

* ArcTo(x0, y0, x1, y1, radius float64) - добавить дугу окружности, используя заданные контрольные точки и радиус.
При необходимости дуга автоматически соединяется с последней точкой пути прямой линией.
x0, y0 - координаты первой контрольной точки;
x1, y1 - координаты второй контрольной точки;
radius - радиус дуги. Должен быть неотрицательным.

* Arc(x, y, radius, startAngle, endAngle float64, clockwise bool) - добавить дугу окружности.
x, y - координаты центра дуги;
radius - радиус дуги. Должен быть неотрицательным;
startAngle - угол в радианах, под которым начинается дуга, измеряется по часовой стрелке от положительной оси X.
endAngle - угол в радианах, под которым заканчивается дуга, измеряется по часовой стрелке от положительной оси X.
clockwise - если true, дуга будет нарисована по часовой стрелке между начальным и конечным углами, иначе - против часовой стрелки

* BezierCurveTo(cp0x, cp0y, cp1x, cp1y, x, y float64) - добавить кубическую кривую Безье из текущей точки.
cp0x, cp0y - координаты первой контрольной точки;
cp1x, cp1y - координаты второй контрольной точки;
x, y - координаты конечной точки.

* QuadraticCurveTo(cpx, cpy, x, y float64) - добавить квадратичную кривую Безье из текущей точки.
cpx, cpy - координаты контрольной точки;
x, y - координаты конечной точки.

* Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, clockwise bool) - добавить эллиптическую дугу.
x, y - координаты центра эллипса;
radiusX - радиус большой оси эллипса. Должен быть неотрицательным;
radiusY - радиус малой оси эллипса. Должен быть неотрицательным;
rotation - вращение эллипса, выраженное в радианах;
startAngle - угол начала эллипса в радианах, измеренный по часовой стрелке от положительной оси x;
endAngle - угол в радианах, под которым заканчивается эллипс, измеренный по часовой стрелке от положительной оси x.
clockwise - если true, рисует эллипс по часовой стрелке, иначе - против часовой стрелки.

Функция Close() вызывается в конце и соединяет начальную и конечную точку фигуры. Используется только для
замкнутых фигур.

После того как Path сформирован его можно нарисовать использую следующие 3 функции

	FillPath(path Path)
	StrokePath(path Path)
	FillAndStrokePath(path Path)

#### Линия

Для рисования линии используется функция

	DrawLine(x0, y0, x1, y1 float64)

### Текст

Для вывода текста в заданных координатах используются две функции

	FillText(x, y float64, text string)
	StrokeText(x, y float64, text string)

Функция StrokeText рисует контур текста, FillText - рисует сам текст.

Горизонтальное выравнивание текста относительно заданных координат устанавливается с помощью функции

	SetTextAlign(align int)

где align может принимать одно из следующих значений:

| Значение | Константа   | Выравнивание                                       |
|:--------:|-------------|----------------------------------------------------|
| 0        | LeftAlign   | Заданная точка является самой левой точкой текста  |
| 1        | RightAlign  | Заданная точка является самой правой точкой текста |
| 2        | CenterAlign | Текст центрируется относительно заданной точки     |
| 3        | StartAlign  | Если текст выводиться слева направо, то вывод текста эквивалентен LeftAlign, иначе RightAlign |
| 4        | EndAlign    | Если текст выводиться слева направо, то вывод текста эквивалентен RightAlign, иначе LeftAlign |
	
Вертикальное выравнивание текста относительно заданных координат устанавливается с помощью функции

	SetTextBaseline(baseline int)

где baseline может принимать одно из следующих значений:

| Значение | Константа           | Выравнивание                                      |
|:--------:|---------------------|---------------------------------------------------|
| 0        | AlphabeticBaseline  | Относительно нормальной базовой линии текста      | 
| 1        | TopBaseline         | Относительно верхней границы текста               |
| 2        | MiddleBaseline      | Относительно середины текста                      |
| 3        | BottomBaseline      | Относительно нижней границы текста                |
| 4        | HangingBaseline     | Относительно подвешенной базовой линии текста (используется тибетскими и другими индийскими шрифтами) |
| 5        | IdeographicBaseline | Относительно идеографической базовой линии текста |

Идеографическая базовая линия это нижняя часть изображения символов, если основная часть символов выступает за базовую линию алфавита
(Используется китайскими, японскими и корейскими шрифтами).

Для установки параметров шрифта выводимого текста используются функции

	SetFont(name string, size SizeUnit)
	SetFontWithParams(name string, size SizeUnit, params FontParams)

где FontParams определена как

	type FontParams struct {
		// Italic - if true then a font is italic
		Italic bool
		// SmallCaps - if true then a font uses small-caps glyphs
		SmallCaps bool
		// Weight - a font weight. Valid values: 0...9, there
		//   0 - a weight does not specify;
		//   1 - a minimal weight;
		//   4 - a normal weight;
		//   7 - a bold weight;
		//   9 - a maximal weight.
		Weight int
		// LineHeight - the height (relative to the font size of the element itself) of a line box.
		LineHeight SizeUnit
	}

Функция TextWidth позволяет узнать ширину выводимого текста в пикселях

	TextWidth(text string, fontName string, fontSize SizeUnit) float64

### Изображение

Перед рисованием изображения его необходимо сначала загрузить. Для этого используется глобальная функция:

	func LoadImage(url string, onLoaded func(Image), session Session) Image {

Изображение загружается асинхронно. После окончания загрузки будет вызвана функция передаваемая во втором аргументе.
Если изображение было загружено успешно, то функция LoadingStatus() интерфейса Image будет возвращать значение
ImageReady (1), если при загрузке произошла ошибка, то данная функция будет возвращать ImageLoadingError (2).
Текстовое описание ошибки возвращает функция LoadingError()

В отличие от ImageView призагрузке Image не учитывается плотность пикселей. Вы должны сами определить какое изображение 
загружать. Это можно сделать так:

	var url string
	if session.PixelRatio() == 2 {
		url = "image@2x.png"
	} else {
		url = "image.png"
	}

Для рисования изображения используются следующие функции:

	DrawImage(x, y float64, image Image)
	DrawImageInRect(x, y, width, height float64, image Image)
	DrawImageFragment(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight float64, image Image)

Функция DrawImage выводит изображение как есть (без масштабирования): x, y - координаты левого верхнего угла изображения

Функция DrawImageInRect выводит изображение с масштабированием: x, y - координаты левого верхнего угла изображения,
width, height - ширина и высота результата

Функция DrawImageFragment выводит фрагмент изображения с масштабированием: srcX, srcY, srcWidth, srcHeight описывают исходную
область изображения, dstX, dstY, dstWidth, dstHeight - результирующая область.

Изображение можно также использовать в стиле заливке

	SetImageFillStyle(image Image, repeat int)

где repeat может принимать следующие значения:

| Значение | Константа | Описание                                           |
|:--------:|-----------|----------------------------------------------------|
| 0        | NoRepeat  | Изображение не повторяется                         | 
| 1        | RepeatXY  | Изображение повторяется по вертикали и горизонтали |
| 2        | RepeatX   | Изображение повторяется только по горизонтали      |
| 3        | RepeatY   | Изображение повторяется только по вертикали        |

## AudioPlayer, VideoPlayer, MediaPlayer

AudioPlayer и VideoPlayer это элементы которые предназначены для воспроизведения аудио и видео.
Оба элемента реализуют интефейс MediaPlayer. Большинство свойств и все события AudioPlayer и VideoPlayer
являются общими и реализуются через MediaPlayer.

### Свойство "src"

Свойство "src" (константа Source) задает один или несколько источников медиафайлов. Свойство "src" может принимать
значение следующих типов:

* string,
* MediaSource,
* []MediaSource.

Структура MediaSource объявлена как

	type MediaSource struct {
		Url      string
		MimeType string
	}

где Url - обязательный параметр, MimeType - необязательный mime тип файла

Так как разные браузеры поддерживают разные форматы файлов и кодеков, то рекомендуется задавать несколько
источников в разных форматах. Плеер сам выбирает из списка источников наиболее подходящий. Задание mime
типов облегчает браузеру этот процесс

### Свойство "controls"

Свойство "controls" (константа Controls) типа bool указывает, должны ли отображаться элементы пользовательского
интерфейса для управления воспроизведения медиа ресурса. Значение по умолчанию false.

Если свойство "controls" равно false для AudioPlayer, то он будет невидим и не будет занимаеть место на экране.

### Свойство "loop"

Свойство "loop" (константа Loop) типа bool. Если оно установлено в true, то медиа-файл начинаться сначала,
когда он достигает конца. Значение по умолчанию false.

### Свойство "muted"

Свойство "muted" (константа Muted) типа bool включает (true) / выключает (false) беззвучный режим. Значение по умолчанию false.

### Свойство "preload"

Свойство "preload" (константа Preload) типа int определяет какие данные должны быть предварительно загружены, если таковые имеются.
Допустимые значения:

| Значение | Константа       | Имя        | Значение                                                                               |
|:--------:|-----------------|------------|----------------------------------------------------------------------------------------|
| 0	       | PreloadNone     | "none"     | Медиа файл не должен быть предварительно загружен                                      |
| 1	       | PreloadMetadata | "metadata" | Предварительно загружаются только метаданные                                           |
| 2	       | PreloadAuto     | "auto"     | Весь медиафайл может быть загружен, даже если пользователь не должен его использовать. |

Значение по умолчанию PreloadAuto (2)

### Свойство "poster"

Свойство "poster" (константа Poster) типа string используется только для VideoPlayer.
Оно задает url картинки которая будет показываться пока видео не загрузится.
Если данное свойство не задано, то будет сначала показываться черный экран, а затем первый кадр (как только он загрузится).

### Свойства "video-width" и "video-height"

Свойства "video-width" (константа VideoWidth) и "video-height" (константа VideoHeight) типа float64 используется только для VideoPlayer.
Оно определяет ширину и высоту выводимого видео в пикселях.

Если "video-width" и "video-height" не заданы, то используются реальные размеры видео, при этом размеры контейнера в
который помещено видео игнорируются и видео может перекрывать другие элементы интерфейса. Поэтому рекоментуется задавать
эти величины, например, так

	rui.Set(view, "videoPlayerContainer", rui.ResizeEvent, func(frame rui.Frame) {
		rui.Set(view, "videoPlayer", rui.VideoWidth, frame.Width)
		rui.Set(view, "videoPlayer", rui.VideoHeight, frame.Height)
	})

Если задано только одно из свойств "video-width" или "video-height", то второе вычисляется на основе пропорций видео

### События

MediaPlayer имеет две группы событий:

1) имеет обработчик вида func(MediaPlayer) (также можно использовать func()). в эту группу входят следующие события

* "abort-event" (константа AbortEvent) - Срабатывает, когда ресурс загружен не полностью, но не в результате ошибки.

* "can-play-event" (константа CanPlayEvent) -  Запускается, когда пользовательский агент может воспроизводить
мультимедиа, но оценивает, что загружено недостаточно данных для воспроизведения мультимедиа до его конца
без необходимости остановки для дальнейшей буферизации контента.

* "can-play-through-event" (константа CanPlayThroughEvent) - Запускается, когда пользовательский агент может воспроизводить мультимедиа,
и оценивает, что было загружено достаточно данных для воспроизведения мультимедиа до его конца, без необходимости
остановки для дальнейшей буферизации контента.

* "complete-event" (константа CompleteEvent) -

* "emptied-event" (константа EmptiedEvent) - Запускается, когда носитель становится пустым; например, когда носитель
уже загружен (или частично загружен)

* "ended-event" (константа EndedEvent) - Срабатывает, когда воспроизведение останавливается, когда достигнут конец носителя
или если дальнейшие данные недоступны.

* "loaded-data-event" (константа LoadedDataEvent) - Запускается, когда первый кадр носителя завершил загрузку.

* "loaded-metadata-event" (константа LoadedMetadataEvent) - Запускается, когда метаданные были загружены.

* "loadstart-event" (константа LoadstartEvent) - Запускается, когда браузер начал загружать ресурс.

* "pause-event" (константа PauseEvent) - Вызывается, когда обрабатывается запрос на приостановку воспроизведения,
и действие переходит в состояние паузы, чаще всего это происходит, когда вызывается метод Pause().

* "play-event" (константа PlayEvent) - Срабатывает, когда начинается воспроизведение медиа файла, например,
в результате использования метода Play()

* "playing-event" (константа PlayingEvent) - Запускается, когда воспроизведение готово начать после приостановки
или задержки из-за отсутствия данных.

* "progress-event" (константа ProgressPvent) - Периодически запускается, когда браузер загружает ресурс.

* "seeked-event" (константа SeekedEvent) - Запускается, когда скорость воспроизведения изменилась.

* "seeking-event" (константа SeekingEvent) - Запускается, когда начинается операция поиска.

* "stalled-event" (константа StalledEvent) - Запускается, когда пользовательский агент пытается извлечь данные мультимедиа,
но данные неожиданно не поступают.

* "suspend-event" (константа SuspendEvent) - Запускается, когда загрузка медиа-данных была приостановлена.

* "waiting-event" (константа WaitingEvent) - Срабатывает, когда воспроизведение остановлено из-за временной нехватки данных

2) имеет обработчик вида func(MediaPlayer, float64) (также можно использовать func(float64), func(MediaPlayer) и func()).
В эту группу входят события связанные с измененим парамеров плеера. В качестве второго аргумента передается
новое значение измененного параметра.

* "duration-changed-event" (константа DurationChangedEvent) - запускается, когда атрибут продолжительности был обновлён.

* "time-updated-event" (константа TimeUpdatedEvent) - запускается, когда текущее время было обновлено.

* "volume-changed-event" (константа VolumeChangedEvent) - запускается при изменении громкости.

* "rate-changed-event" (константа RateChangedEvent) - запускается, когда скорость воспроизведения изменилась.

Отдельное событие, не относящееся к этим двум группам, "player-error-event" (константа PlayerErrorEvent) срабатывает,
когда ресурс не может быть загружен из-за ошибки (например, ошибки сети).

Обработчик данного события имеет вида func(player MediaPlayer, code int, message string) (также можно использовать
func(int, string), func(MediaPlayer) и func()). Где аргумент "message" это сообщение об ошибке, "code" - код ошибки:

| Код ошибки | Константа                     |  Значение                                                                    |
|:----------:|-------------------------------|------------------------------------------------------------------------------|
| 0	         | PlayerErrorUnknown            | Неизвестная ошибка                                                           |
| 1	         | PlayerErrorAborted            | Извлечение связанного ресурса было прервано запросом пользователя.           |
| 2	         | PlayerErrorNetwork            | Произошла какая-то сетевая ошибка, которая помешала успешному извлечению носителя, несмотря на то, что он был ранее доступен. |
| 3	         | PlayerErrorDecode             | Несмотря на то, что ранее ресурс был определён, как используемый, при попытке декодировать медиаресурс произошла ошибка. |
| 4	         | PlayerErrorSourceNotSupported | Связанный объект ресурса или поставщик мультимедиа был признан неподходящим. |

### Методы

MediaPlayer имеет ряд методов для управления параметрами плеера:

* Play() - запускает воспроизведение медиа файла;

* Pause() - ставит воспроизведение на паузу;
	
* SetCurrentTime(seconds float64) - устанавливает текущее время воспроизведения в секундах;
	
* CurrentTime() float64 - возвращает текущее время воспроизведения в секундах;

* Duration() float64 - возвращает длительность медиа файла  в секундах;

* SetPlaybackRate(rate float64) - устанавливает скорость воспроизведения. Нормальная скорость равна 1.0;

* PlaybackRate() float64 - возвращает текущую скорость воспроизведения;

* SetVolume(volume float64) - устанавливает скорост громкость в диапазоне от 0 (тишина) до 1 (максимальная громкость);
	
* Volume() float64 - возвращает текущую громкость;
	
* IsEnded() bool - возвращает true если достигнут конец медиа файла;
	
* IsPaused() bool - возвращает true если воспроизведение поставлено на паузу.

Для быстрого доступа к этим методам имеются глобальные функции:

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

где view - корневой View, playerID - id of AudioPlayer or VideoPlayer

## Анимация

Библиотека поддерживает два вида анимации:

* Анимированое изменения значения свойства (далее "анимация перехода") 
* Сценарий анимированного изменения одного или нескольких свойств (далее просто "сценарий анимации")

### Интерфейс Animation

Для задания параметров анимации мспользуется интерфейс Animation. Он расширяет интерфейс Properties.
Интерфейс создается с помощью функции:

	func NewAnimation(params Params) Animation

Часть свойств интерфейса Animation используется в обоих типах анимации, остальные используются 
только в сценариях анимации.

Общими свойствами являются

| Свойство          | Константа      | Тип     | По умолчанию | Описание                                        |
|-------------------|----------------|---------|--------------|-------------------------------------------------|
| "duration"        | Duration       | float64 | 1            | Длительность анимации в секундах                |
| "delay"           | Delay          | float64 | 0            | Длительность задержки перед анимации в секундах |
| "timing-function" | TimingFunction | string  | "ease"       | Функция изменения скорости анимации             |

Свойства используемые только в сценариях анимации будут описаны ниже

#### Свойство "timing-function"

Свойство "timing-function" описывает в текстовом виде функция изменения скорости анимации. Функции 
могут быть разделены на 2 вида: простые функции и функции с параметрами.

Простые функции

| Функция       | Константа       |  Описание                                                                     |
|---------------|-----------------|-------------------------------------------------------------------------------|
| "ease"        | EaseTiming      | скорость увеличивается к середине, а в конце замедляется.                     |
| "ease-in"     | EaseInTiming    | скорость вначале медленная, но в конце увеличивается.                         |
| "ease-out"    | EaseOutTiming   | скорость вначале быстрая, но быстро снижается. Большая часть медленная        |
| "ease-in-out" | EaseInOutTiming | скорость вначале быстрая, но быстро снижается, а в конце снова увеличивается. |
| "linear"      | LinearTiming    | постоянная скорость                                                           |

И есть две функции с параметрами:

* "steps(N)" - дискретная функция, где N - целое число задающее количество шагов. Вы можете задавать
данную функцию или в виде текста или использую функцию:

	func StepsTiming(stepCount int) string

Например

	animation := rui.NewAnimation(rui.Params{
		rui.TimingFunction: rui.StepsTiming(10),
	})

эквивалентно 

	animation := rui.NewAnimation(rui.Params{
		rui.TimingFunction: "steps(10)",
	})
	
* "cubic-bezier(x1, y1, x2, y2)" - временная функция кубической кривой Безье. x1, y1, x2, y2 имеют тип float64. 
x1 и x2 должны быть в диапазоне [0, 1]. Вы можете задавать данную функцию или в виде текста или использую функцию:

	func CubicBezierTiming(x1, y1, x2, y2 float64) string

### Анимация перехода

Анимация перехода может применяться к свойствам типа: SizeUnit, Color, AngleUnit, float64 и составным свойсвам
в составе которых имеются элементы перечисленных типов (например "shadow", "border" и т.д.).

При попытке применить анимацию к свойствам других типов (например, bool, string) ошибки не произойдет,
просто анимации не будет.

Анимация перехода бывает двух видов:
* аднократная;
* постоянная;

Однократная анимация запускается с помощью функции SetAnimated интерфейса View. Данная функция имеет следующее
описание:

	SetAnimated(tag string, value interface{}, animation Animation) bool

Она присваивает свойству новое значение, при этом изменение происходит с использованием заданной анимации.
Например,

	view.SetAnimated(rui.Width, rui.Px(400), rui.NewAnimation(rui.Params{
		rui.Duration:       0.75,
		rui.TimingFunction: rui.EaseOutTiming,
	}))

Есть также глобальная функция для анимированного однократного изменения значения свойства дочернего View

	func SetAnimated(rootView View, viewID, tag string, value interface{}, animation Animation) bool

Постоянная анимация запускается каждый раз когда изменяется значение свойства. Для задания постоянной 
анимации перехода используется свойство "transition" (константа Transition). В качества значения данному 
свойству присваивается rui.Params, где в качестве ключа должно быть имя свойства, а значение - интерфейс Animation.
Например,

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

Вызов функции SetAnimated не меняет значение свойства "transition".

Для получения текущего списка постоянных анимаций перехода используется функция

	func GetTransition(view View, subviewID string) Params

Добавлять новые анимации перехода рекомендуется с помощью функции 

	func AddTransition(view View, subviewID, tag string, animation Animation) bool

Вызов данной функции эквивалентен следующему коду

	transitions := rui.GetTransition(view, subviewID)
	transitions[tag] = animation
	rui.Set(view, subviewID, rui.Transition, transitions)

#### События анимации перехода

Анимация перехода генерирует следующие события

| Событие                   | Константа             | Описание                                           |
|---------------------------|-----------------------|----------------------------------------------------|
| "transition-run-event"    | TransitionRunEvent    | Цикл анимации перехода стартовал, т.е. до задержки |
| "transition-start-event"  | TransitionStartEvent  | Анимация перехода действительно стартовала, т.е. после задержки |
| "transition-end-event"    | TransitionEndEvent    | Анимация перехода закончена                        |
| "transition-cancel-event" | TransitionCancelEvent | Анимация перехода прервана                         |

Основной слушатель данных событий имеет следующий формат:

	func(View, string)

где второй аргумент это имя свойства.

Можно также использовать слушателя следующего формата:

	func()
	func(string)
	func(View)

Получить списки слушателей событий анимации перехода с помощью функций:

	func GetTransitionRunListeners(view View, subviewID string) []func(View, string)
	func GetTransitionStartListeners(view View, subviewID string) []func(View, string)
	func GetTransitionEndListeners(view View, subviewID string) []func(View, string)
	func GetTransitionCancelListeners(view View, subviewID string) []func(View, string)

### Cценарий анимации

Cценарий анимации описывает более сложную анимацию, по сравнению с анимацией перехода. Для этого
в Animation добавляются дополнительные свойства:

#### Свойство "property"

Свойство "property" (константа PropertyTag) описывает изменения свойств. В качестве значения ему присваивается
[]AnimatedProperty. Структура AnimatedProperty описывает изменение одного свойства. Она описана как

	type AnimatedProperty struct {
		Tag       string
		From, To  interface{}
		KeyFrames map[int]interface{}
	}

где Tag - имя свойства, From - начальное значение свойства, To - конечное значение свойства,
KeyFrames - промежуточные значения свойства (ключевые кадры).

Обязательными являются поля Tag, From, To. Поле KeyFrames опционально, может быть nil.

Поле KeyFrames описывает ключавые кадры. В качестве ключа типа int используется процент времени 
прошедший с начала анимации (именно начала самой анимации, время заданное свойством "delay" исключается).
А значание это значение свойства в данный момент времени. Например

	prop := rui.AnimatedProperty {
		Tag:       rui.Width,
		From:      rui.Px(100),
		To:        rui.Px(200),
		KeyFrames: map[int]interface{
			90: rui.Px(220),
		}
	}

В данном примере свойство "width" 90% времени будет увеличиваться со 100px до 220px. В оставшиеся 
10% времени - будет уменьшаться с 220px до 200px.

Свойству "property" присваивается []AnimatedProperty, а значит можно анимаровать сразу несколько свойств.

Вы должны задать хотя бы один эдемент "property", иначе анимация будет игнорироваться.

#### Свойство "id"

Свойство "id" (константа ID) типа string задает идентификатор анимации. Передается в качестве параметра слушателю 
события анимации. Если вы не планируете использовать слушателей событий для анимации, то данное свойство
можно не задавать.

#### Свойство "iteration-count"

Свойство "iteration-count" (константа IterationCount) типа int задает количество повторений анимации.
Значение по умолчанию 1. Значение меньше нуля заставляет повторяться анимацию бесконечно.

#### Свойство "animation-direction"

Свойство "animation-direction" (константа AnimationDirection) типа int устанавливает, должна ли анимация
воспроизводиться вперед, назад или поочередно вперед и назад между воспроизведением 
последовательности вперед и назад. Может принимать следующие значения:

| Значение | Константа                 |  Описание                                                             |
|:--------:|---------------------------|-----------------------------------------------------------------------|
| 0        | NormalAnimation           | Анимация проигрывается вперёд каждую итерацию, то есть, когда анимация заканчивается, она сразу сбрасывается в начальное положение и снова проигрывается. |
| 1        | ReverseAnimation          | Анимация проигрывается наоборот, с последнего положения до первого и потом снова сбрасывается в конечное положение и снова проигрывается. |
| 2        | AlternateAnimation        | Анимация меняет направление в каждом цикле, то есть в первом цикле она начинает с начального положения, доходит до конечного, а во втором цикле продолжает с конечного и доходит до начального и так далее |
| 3        | AlternateReverseAnimation | Анимация начинает проигрываться с конечного положения и доходит до начального, а в следующем цикле продолжая с начального переходит в конечное |

#### Запуск анимации

Для запуска сценария анимации необходимо созданный Animation интерфейс присвоить свойству "animation"
(константа AnimationTag). Если View уже отображается на экране, то анимация запускается сразу (с учетом
заданной задержки), в противоположном случае анимация запускается как только View отобразится на экране.

Свойству "animation" можно присваивать Animation и []Animation, т.е. можно запускать несколько анимаций 
одновременно для одного View

Пример,

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

#### Свойство "animation-paused"

Свойство "animation-paused" (константа AnimationPaused) типа bool позволяет приостановить анимацию.
Для того чтобы поставить анимацию на паузу необходимо данному свойству присвоить значение true, а
для возобновления - false.

Внимание. В момент присваивания значения свойству "animation" cвойство "animation-paused" сбрасывается в false.
 
#### События анимации

Сценарий анимации генерирует следующие события

| Событие                     | Константа               | Описание                         |
|-----------------------------|-------------------------|----------------------------------|
| "animation-start-event"     | AnimationStartEvent     | Анимация стартовала              |
| "animation-end-event"       | AnimationEndEvent       | Анимация закончена               |
| "animation-cancel-event"    | AnimationCancelEvent    | Анимация прервана                |
| "animation-iteration-event" | AnimationIterationEvent | Началась новая итерация анимации |

Внимание! Не все браузеры поддерживают событие "animation-cancel-event". В данное время это только
Safari и Firefox.

Основной слушатель данных событий имеет следующий формат:

	func(View, string)

где второй аргумент это id анимации.

Можно также использовать слушателя следующего формата:

	func()
	func(string)
	func(View)

Получить списки слушателей событий анимации с помощью функций:

	func GetAnimationStartListeners(view View, subviewID string) []func(View, string)
	func GetAnimationEndListeners(view View, subviewID string) []func(View, string)
	func GetAnimationCancelListeners(view View, subviewID string) []func(View, string)
	func GetAnimationIterationListeners(view View, subviewID string) []func(View, string)

## Сессия

Когда клиент создает соединение с сервером, то для этого соединения создается интерфейс Session.
Этот интерфейс используется для взаимодействия с клиентом.
Получить текущий интерфейс Session можно вызвав метод Session() интерфейса View.

При создани сессии она получает пользовательскую реализацию интерфейса SessionContent.

	type SessionContent interface {
		CreateRootView(session rui.Session) rui.View
	}

Данный интерфейс создается функцией передаваемой в качестве параметра при создании приложения
функцией NewApplication.

Кроме обязательной функции CreateRootView() SessionContent может иметь несколько опциональных
функций:

	OnStart(session rui.Session)
	OnFinish(session rui.Session)
	OnResume(session rui.Session)
	OnPause(session rui.Session)
	OnDisconnect(session rui.Session)
	OnReconnect(session rui.Session)

Сразу после создания сессии вызывается функция CreateRootView. После создания корневого View
вызывается функцию OnStart (если она реализована)

Функция OnFinish (если она реализована) вызывается когда пользователь закрывет страницу приложения в браузере

Функция OnPause вызывается когда страница приложения в браузере клиента становится неактивной.
Это происходит если пользователь переключается на другую вкладку/окно браузера, сворачивает браузер
или переключается на другое приложение.

Функция OnResume вызывается когда страница приложения в браузере клиента становится активной. Так же
эта функция вызывается сразу после OnStart

Функция OnDisconnect вызывается если сервер теряет соединение с клиентом. Это происходит либо при
обрыве связи.

Функция OnReconnect вызывается после того как сервер восстанавливает соединение с клиентом.

Интерфейс Session предоставляет следующие методы:

* DarkTheme() bool - возвращает true, если используется темная тема. Определяется настройками на стороне клиента

* TouchScreen() bool - возвращает true, если клиент поддерживает touch screen

* PixelRatio() float64 - возвращает размер логического пикселя, т.е. сколько физических пикселей образуют логический. Например, для iPhone это значение будет 2 или 3

* TextDirection() int - возвращает напрвление письма: LeftToRightDirection (1) или RightToLeftDirection (2)

* Constant(tag string) (string, bool) - возвращает значение константы

* Color(tag string) (Color, bool) - возвращает значение константы цвета

* SetCustomTheme(name string) bool - устанавливает тему заданым именем в качестве текущей. Возвращает false если тема с таким именем не найдена. Темы с именем "" это тема по умолчанию.

* Language() string - возвращает текущий язык интерфейса, например: "en", "ru", "ptBr"

* SetLanguage(lang string) - устанавливает текущий язык интерфейса (см. "Поддержка нескольких языков")

* GetString(tag string) (string, bool) - возвращает текстовое текстовое значение для текущего языка
(см. "Поддержка нескольких языков")

* Content() SessionContent - возвращает текущий экземпляр SessionContent

* RootView() View - возвращает корневой View сессии

* SetTitle(title string) - устанавливает текст заголовка окна/закладки браузера.

* SetTitleColor(color Color) устанавливает цвет панели навигации браузера. Поддерживается только в Safari и Chrome для Android.

* Get(viewID, tag string) interface{} - возвращает значение свойства View с именем tag. Эквивалентно

	rui.Get(session.RootView(), viewID, tag)

* Set(viewID, tag string, value interface{}) bool - устанавливает значение свойства View с именем tag.

	rui.Set(session.RootView(), viewID, tag, value)

* DownloadFile(path string) - загружает (сохраняет) на стороне клиента файл расположенный по заданному пути на сервере. 
Используется когда клиенту надо передать с сервера какой-либо файл.

* DownloadFileData(filename string, data []byte) - загружает (сохраняет) на стороне клиента файл с заданным именем и
заданным содержимым. Обычно используется для передачи файла сгенерированного в памяти сервера.	

## Формат описания ресурсов

Ресурсы приложения (темы, View, переводы) могут быть описаны в виде текста (utf-8). Данный текст помещается
в файл с расширением ".rui".

Корневым элементом файла ресурса должен быть объект. Он имеет следующий формат:

	<имя объекта> {
		<данные объекта>
	}

если имя объекта содержит следующие символы: '=', '{', '}', '[', ']', ',', ' ', '\t', '\n', '\'', '"',
'`', '/' и любые пробелы, то имя объекта необходимо брать в кавычки. Если  эти символы не используются,
то кавычки не обязательны.

Можно использовать три вида ковычек:

* "…" - эквивалентна такой же строке в языке go, т.е. внутри можно использовать escape последовательности:
\n, \r, \\, \", \', \0, \t, \x00, \u0000

* '…' - аналогична строке "…"

* `…` - эквивалентна такой же строке в языке go, т.е. текст внутри этой строки остается как есть. Внутри
нельзя использовать символ `.

Данные объяекта представляют собой множество пар <ключ> = <значение> разделенные запятой.

Ключ это строка текста. Правила оформления такие же как и у имени объекта.

Значения могут быть 3 видов:

* Простое значение - строка текста оформленная по тем же правилам, что и имя объекта

* Объект

* Массив значений

Массив значений заключается в квадратные скобки. Элементы массива разделяются запятыми.
Элементами могут быть простые значения или объекты.

В текста могут быть комментарии. Правила оформления такие же как в языке go: // и /* … */

Пример:

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

Для работы с текстовыми ресурсами используется интерфейс DataNode

	type DataNode interface {
		Tag() string
		Type() int
		Text() string
		Object() DataObject
		ArraySize() int
		ArrayElement(index int) DataValue
		ArrayElements() []DataValue
	}

Данный элемент описывает базовый элемент данных.

Метод Tag возвращает значение ключа.

Тип данных возвращается методом Type. Он возвращает одно из 3 значений

| Значение | Константа  | Тип данных         |
|:--------:|------------|--------------------|
| 0	       | TextNode   | Простое значение   |
| 1	       | ObjectNode | Объект             |
| 2        | ArrayNode  | Массив             |

Для получения простого значения используется метод Text.
Для получения объекта используется метод Object.
Для получения элементов массива используюся методы ArraySize, ArrayElement и ArrayElements

## Ресурсы

Ресурсы (картинки, темы, переводы и т.д.) с которыми работает приложение должны размещаться по
поддиректориям внутри одного директория ресурсов. Ресурсы должны располагаться в следующих поддерикториях:

* images - в данную поддиректорию помещаются все изображения. Здесь можно делать вложенные поддиректории.
В этом случае их надо включать в имя файла. Например, "subdir/image1.png"

* themes - в данную поддиректорию помещаются темы приложения (см. ниже)

* views - в данную поддиректорию помещаются описания View

* strings - в данную поддиректорию помещаются переводы текстовых ресурсов (см. Поддержка нескольких языков)

* raw - в данную поддиректорию помещаются все остальные ресурсы: звуки, видео, двоичные данные и т.п.

Директория с ресурсами может или включаться в исполняемый файл или располагаться отдельно.

Если ресурсы необходимо включить в исполняемый файл, то имя директории должно быть "resources" и
подключаться она должны следующим образом:

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

Если ресурсы поставляются в виде отдельной директории, то ее необходимо зарегистрировать
с помощью функции SetResourcePath до создания Application:

	func main() {
		rui.SetResourcePath(path)
		
		app := rui.NewApplication("Hello world", createHelloWorldSession)
		app.Start("localhost:8000")
	}

## Изображения для экранов с разной плотностью пикселей

Если вам необходимо добавить в ресурсы отдельные изображения для экранов с разной плотностью пикселей,
то это делается в стиле iOS. Т.е. к имени файла добавляется '@<плотность>x'. Например

	image@2x.png
	image@3x.jpg
	image@1.5x.gif

Например, у вас есть изображения для трех плотностей: image.png, image@2x.png и image@3x.png.
В этом случае полю "src" ImageView вы присваиваете только значение "image.png". Библиотека
сама найдет остальные в директории "images" и передаст клиенту изображение с
требуемой плотностью

## Темы

Тема включает в сябя три вида данных:
* константы
* константы цвета
* Стили View

Темы оформляются в виде rui файла и помещаются в папку themes.

Корневым объектом темы является объект с именем 'theme'. Данный объект может содержать следующие свойства:

* name - текстовое свойство задающее имя темы. Если данное свойство не задано или оно равно пустой строке,
то это тема по умолчанию.

* constants - свойство-объект определяющий константы. Имя объекта может быть любым. Рекомендуется использовать "_".
Объект может иметь любое количиство текстовых свойств задающих пару "имя константы" = "значение".
В данном разделе помещаются константы типа SizeUnit, AngleUnit, текстовые и числовые. Для того чтобы
присвоить константу какому либо свойству View надо свойству присвоить имя константы добавив вначале символ '@'.
Например

	theme {
		constants = _{
			defaultPadding = 4px,
			buttonPadding = @defaultPadding,
			angle = 30deg,
		}
	}

	rui.Set(view, "subView", rui.Padding, "@defaultPadding")

* constants:touch - свойство-объект определяющий константы используемые только для touch screen.
Например, как сделать отступы больше на touch screen:

	theme {
		constants = _{
			defaultPadding = 4px,
		},
		constants:touch = _{
			defaultPadding = 12px,
		},
	}

* colors - свойство-объект определяющий цветовые константы для светлой темы оформления (тема по умолчанию).
Объект может иметь любое количиство текстовых свойств задающих пару "имя цвета" = "цвет". Аналогично
константам, при присваивании необходимо вначале имени цвета добавить '@'. Например

	theme {
		colors = _{
			textColor = #FF101010,
			borderColor = @textColor,
			backgroundColor = white,
		}
	}

	rui.Set(view, "subView", rui.TextColor, "@textColor")

Имена цветов, такие как "black", "white", "red" и т.д., используются без символа '@'. При этом вы можете
задавать цветовые константы с такими же именами. Например

	theme {
		colors = _{
			red = blue,
		}
	}

	rui.Set(view, "subView", rui.TextColor, "@red") // blue text
	rui.Set(view, "subView", rui.TextColor, "red")  // red text

* colors:dark - свойство-объект определяющий цветовые константы для темной темы оформления

* styles - массив общих стилей. Каждый элемент массива должен быть объектом. Имя объекта является ивляется
именем стиля. Например,

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

Для использования стилей у View есть два текстовых свойства "style" (константа Style) и "style-disabled"
(константа StyleDisabled). Свойству "style" присваивается имя свойства которое применяется ко View при
значении свойства "disabled" равного false. Свойству "style-disabled" присваивается имя свойства
которое применяется ко View при значении свойства "disabled" равного true. Если "style-disabled"
не определен, то всойство "style" используется в обоих режимах.

Внимание! Символ '@' к имени стиля добавлять НЕ НАДО. Если вы добавите символ '@' к имени, то имя
стиля будет извлекаться из одноименной константы. Например

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

Помимо общих стилей можно затать стили для отдельных режимов работы. Для этого к имени "styles" добавляются
следующие модификаторы:

* ":portrait" или ":landscape" - соответственно стили для портретного или ладшафтного режима программы.
Внимание имеется ввиду соотношение сторон окна программы, а не экрана.

* ":width<размер>" - стили для экрана ширина которого не превышает заданный размер в логических пикселях.

* ":height<размер>" - стили для экрана высота которого не превышает заданный размер в логических пикселях.

Например

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

## Стандартные константы и стили

В библиотеке определен ряд констант и стилей. Вы их можете переопределять в своих темах.

Системные стили которые вы можете переопределять:

| Имя стиля           | Описание                                                                    |
|---------------------|-----------------------------------------------------------------------------|
| ruiApp              | Данный стиль используется для назначения стиля текста (шрифт, размер и т.д.) по умолчанию |
| ruiView             | Стиль View по умолчанию                                                     |
| ruiArticle          | Стиль используемый если свойство "semantics" установлено в "article"        |
| ruiSection          | Стиль используемый если свойство "semantics" установлено в "section"        |
| ruiAside            | Стиль используемый если свойство "semantics" установлено в "aside"          |
| ruiHeader           | Стиль используемый если свойство "semantics" установлено в "header"         |
| ruiMain             | Стиль используемый если свойство "semantics" установлено в "main"           |
| ruiFooter           | Стиль используемый если свойство "semantics" установлено в "footer"         |
| ruiNavigation       | Стиль используемый если свойство "semantics" установлено в "navigation"     |
| ruiFigure           | Стиль используемый если свойство "semantics" установлено в "figure"         |
| ruiFigureCaption    | Стиль используемый если свойство "semantics" установлено в "figure-caption" |
| ruiButton           | Стиль используемый если свойство "semantics" установлено в "button"         |
| ruiParagraph        | Стиль используемый если свойство "semantics" установлено в "paragraph"      |
| ruiH1               | Стиль используемый если свойство "semantics" установлено в "h1"             |
| ruiH2               | Стиль используемый если свойство "semantics" установлено в "h2"             |
| ruiH3               | Стиль используемый если свойство "semantics" установлено в "h3"             |
| ruiH4               | Стиль используемый если свойство "semantics" установлено в "h4"             |
| ruiH5               | Стиль используемый если свойство "semantics" установлено в "h5"             |
| ruiH6               | Стиль используемый если свойство "semantics" установлено в "h6"             |
| ruiBlockquote       | Стиль используемый если свойство "semantics" установлено в "blockquote"     |
| ruiCode             | Стиль используемый если свойство "semantics" установлено в "code"           |
| ruiTable            | Стиль TableView по умолчанию                                                |
| ruiTableHead        | Стиль заголовка TableView по умолчанию                                      |
| ruiTableFoot        | Стиль итого TableView по умолчанию                                          |
| ruiTableRow         | Стиль строки TableView по умолчанию                                         |
| ruiTableColumn      | Стиль колонки TableView по умолчанию                                        |
| ruiTableCell        | Стиль ячейки TableView по умолчанию                                         |
| ruiDisabledButton   | Стиль Button если свойство "disabled" установлено в true                    |
| ruiCheckbox         | Стиль Checkbox                                                              |
| ruiListItem         | Стиль пункта ListView                                                       |
| ruiListItemSelected | Стиль выбранного пункта ListView когда ListView не владеет фокусом          |
| ruiListItemFocused  | Стиль выбранного пункта ListView когда ListView владеет фокусом             |
| ruiPopup            | Стиль всплывающего окна                                                     |
| ruiPopupTitle       | Стиль заголовка всплывающего окна                                           |
| ruiMessageText      | Стиль текста всплывающего окна (Message, Question)                          |
| ruiPopupMenuItem    | Стиль пункта всплывающего меню                                              |

Системные цвета которые вы можете переопределять:

| Имя константы цвета        | Описание                                           |
|----------------------------|----------------------------------------------------|
| ruiBackgroundColor         | Цвет фона                                          |
| ruiTextColor               | Цвет текста                                        |
| ruiDisabledTextColor       | Цвет запрещенного текста                           |
| ruiHighlightColor          | Цвет подсветки                                     |
| ruiHighlightTextColor      | Цвет подсвеченного текста                          |
| ruiButtonColor             | Цвет кнопки                                        |
| ruiButtonActiveColor       | Цвет кнопки в фокусе                               |
| ruiButtonTextColor         | Цвет текста кнопки                                 |
| ruiButtonDisabledColor     | Цвет запрещенной кнопки                            |
| ruiButtonDisabledTextColor | Цвет текста запрещенной кнопки                     |
| ruiSelectedColor           | Цвет фона неактивного выбранного пункта ListView   |
| ruiSelectedTextColor       | Цвет текста неактивного выбранного пункта ListView |
| ruiPopupBackgroundColor    | Цвет фона всплывающего окна                        |
| ruiPopupTextColor          | Цвет текста всплывающего окна                      |
| ruiPopupTitleColor         | Цвет фона заголовка всплывающего окна              |
| ruiPopupTitleTextColor     | Цвет текста заголовка всплывающего окна            |

Константы которые вы можете переопределять:

| Имя константы                | Описание                                      |
|------------------------------|-----------------------------------------------|
| ruiButtonHorizontalPadding   | Горизонтальный отступ внутри кнопки           |
| ruiButtonVerticalPadding     | Вертикальный  отступ внутри кнопки            |
| ruiButtonMargin              | Внешний тоступ кнопки                         |
| ruiButtonRadius              | Радиус скругления углов кнопки                |
| ruiButtonHighlightDilation   | Ширина внешней рамки активной кнопки          |
| ruiButtonHighlightBlur       | Размытие рамки активной кнопки                |
| ruiCheckboxGap               | Разрыв между checkbox и содержимым            |
| ruiListItemHorizontalPadding | Горизонтальный отступ внутри пункта ListView  |
| ruiListItemVerticalPadding   | Вертикальный отступ внутри пункта ListView    |
| ruiPopupTitleHeight          | Высота заголовка всплывающего окна            |
| ruiPopupTitlePadding         | Внутренний отступ заголовка всплывающего окна |
| ruiPopupButtonGap            | Разрыв между кнопками всплывающего окна       |

## Поддержка нескольких языков

Если вы хотите добавить в программу поддержку нескольких языков, то необходимо поместить в папку "strings" ресурсов
файлы с переводом. Файлы перевода должны иметь расширение "rui" и следующий формат

	strings {
		<язык 1> = _{
			<Исходный текст 1> = <Перевод 1>,
			<Исходный текст 2> = <Перевод 2>,
			…
		},
		<язык 2> = _{
			<Исходный текст 1> = <Перевод 1>,
			<Исходный текст 2> = <Перевод 2>,
			…
		},
		…
	}

Если перевод на каждый язык помещается в отдельный файл, то можно использовать следующий формат

	strings:<язык> {
		<Исходный текст 1> = <Перевод 1>,
		<Исходный текст 2> = <Перевод 2>,
		…
	}

Например, если все переводы в одном файле strings.rui

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

Если в разных. Файл ru.rui

	strings:ru {
		"Yes" = "Да",
		"No" = "Нет",
	}

Файл de.rui

	strings:de {
		"Yes" = "Ja",
		"No" = "Nein",
	}

Перевод можно также разбивать на несколько файлов.

Переводы автоматически подставляются во всех View.

Однако если вы рисуете текст в CanvasView, то вы должны запрашивать перевод сами. Для этого в интерфейсе Session
есть метод:

	GetString(tag string) (string, bool)

Если перевода данной строки нет, то метод вернет исходную строку и false в качестве второго параметра.

Получить текущий язык можно с помощью метода Language() интерфейса Session. Текущий язык определяется настройками
браузера пользователя. Поменять язык сессии можно спомощью метода SetLanguage(lang string) интерфейса Session.

