package rui

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

type formElement struct {
	data     any
	index    int
	id       string
	listener func(string)
}

type formDropDown struct {
	formElement
}

type formNumber struct {
	formElement
}

type formFlag struct {
	formElement
}

type formText struct {
	formElement
}

func (element *formElement) fieldValue() reflect.Value {
	structValue := reflect.ValueOf(element.data)
	if structValue.Kind() == reflect.Pointer {
		structValue = structValue.Elem()
	}

	return structValue.Field(element.index)
}

func (formFlag *formFlag) changed(checkbox Checkbox, checked bool) {
	val := formFlag.fieldValue()
	if val.Kind() == reflect.Bool && val.Bool() != checked {
		val.SetBool(checked)
		formFlag.listener(formFlag.id)
	}
}

func (text *formText) changed(checkbox EditView, newText, oldText string) {
	val := text.fieldValue()
	if val.Kind() == reflect.String && val.String() != newText {
		val.SetString(newText)
		text.listener(text.id)
	}
}

func (dropDown *formDropDown) changed(list DropDownList, index int) {
	switch val := dropDown.fieldValue(); val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		newInt := int64(index)
		if newInt == val.Int() {
			return
		}
		val.SetInt(newInt)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if index < 0 {
			return
		}
		newUint := uint64(index)
		if newUint == val.Uint() {
			return
		}
		val.SetUint(newUint)

	default:
		return
	}

	if dropDown.listener != nil {
		dropDown.listener(dropDown.id)
	}
}

func (number *formNumber) changed(picker NumberPicker, newValue float64, oldValue float64) {
	switch val := number.fieldValue(); val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		newInt := int64(newValue)
		if newInt == val.Int() {
			return
		}
		val.SetInt(newInt)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		newUint := uint64(newValue)
		if newUint == val.Uint() {
			return
		}
		val.SetUint(newUint)

	case reflect.Float32, reflect.Float64:
		if newValue == val.Float() {
			return
		}
		val.SetFloat(newValue)

	default:
		return
	}

	if number.listener != nil {
		number.listener(number.id)
	}
}

func createForm(session Session, data any, structValue reflect.Value, listener func(string)) View {
	result := NewGridLayout(session, Params{
		ColumnWidth: []SizeUnit{AutoSize(), AutoSize()},
		ColumnGap:   Px(16),
	})

	lineNumber := 0
	params := map[string]string{}

	getFloatParam := func(tag string) (float64, bool) {
		if value, ok := params[tag]; ok {
			n, err := strconv.ParseFloat(value, 64)
			if err == nil {
				return n, true
			}
			ErrorLog(err.Error())
		}
		return 0, false
	}

	getItems := func() []string {
		if value, ok := params["items"]; ok {
			if length := len(value); length > 0 && value[0] == '[' && value[length-1] == ']' {
				return strings.Split(value[1:length-1], ",")
			}
			return []string{}
		}

		return nil
	}

	createDropDown := func(index int, id string, current int) View {
		items := getItems()
		if items != nil {
			if len(items) > 0 {
				dropDown := new(formDropDown)
				dropDown.data = data
				dropDown.index = index
				dropDown.id = id
				dropDown.listener = listener

				return NewDropDownList(session, Params{
					ID:            id,
					Items:         items,
					Current:       current,
					DropDownEvent: dropDown.changed,
				})
			} else {
				ErrorLog("invalid items list")
			}
		}
		return nil
	}

	createNumberPicker := func(index int, id string, current float64, step int, minValue int, maxValue int) View {

		numberType := NumberEditor
		if value, ok := params["type"]; ok && strings.ToLower(value) == "slider" {
			numberType = NumberSlider
		}

		numberData := new(formNumber)
		numberData.data = data
		numberData.index = index
		numberData.id = id
		numberData.listener = listener

		picker := NewNumberPicker(session, Params{
			ID:                 id,
			NumberPickerType:   numberType,
			NumberPickerValue:  current,
			NumberChangedEvent: numberData.changed,
		})

		if n, ok := getFloatParam("step"); ok {
			picker.Set(NumberPickerStep, n)
		} else if step != 0 {
			picker.Set(NumberPickerStep, step)
		}

		if n, ok := getFloatParam("min"); ok {
			picker.Set(NumberPickerMin, n)
		} else if minValue != -1 {
			picker.Set(NumberPickerMin, minValue)
		}

		if n, ok := getFloatParam("max"); ok {
			picker.Set(NumberPickerMax, n)
		} else if maxValue != 0 {
			picker.Set(NumberPickerMax, maxValue)
		}

		return picker
	}

	for i := range structValue.NumField() {
		val := structValue.Field(i)
		field := structValue.Type().Field(i)

		if tag := field.Tag.Get("rui"); tag != "" {
			params = parseFieldTag(tag)
		} else {
			params = map[string]string{}
		}

		var id string
		if name, ok := params["id"]; ok && name != "" {
			id = name
		} else {
			id = field.Name
		}

		var controlView View = nil

		switch val.Kind() {

		case reflect.Int, reflect.Int64:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, -1, 0)
			}

		case reflect.Int8:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, math.MinInt8, math.MaxInt8)
			}

		case reflect.Int16:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, math.MinInt16, math.MaxInt16)
			}

		case reflect.Int32:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, math.MinInt32, math.MaxInt32)
			}

		case reflect.Uint, reflect.Uint64:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, 0, 0)
			}

		case reflect.Uint8:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, 0, math.MaxUint8)
			}

		case reflect.Uint16:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, 0, math.MaxUint16)
			}

		case reflect.Uint32:
			controlView = createDropDown(i, id, int(val.Int()))
			if controlView == nil {
				controlView = createNumberPicker(i, id, float64(val.Int()), 1, 0, math.MaxUint32)
			}

		case reflect.Float32, reflect.Float64:
			controlView = createNumberPicker(i, id, val.Float(), 0, -1, 0)

		case reflect.Bool:
			label := params["label"]
			flag := new(formFlag)
			flag.data = data
			flag.index = i
			flag.id = id
			flag.listener = listener

			result.Append(NewCheckbox(session, Params{
				ID:                   id,
				Row:                  lineNumber,
				Column:               Range{First: 0, Last: 1},
				Content:              label,
				Checked:              val.Bool(),
				CheckboxChangedEvent: flag.changed,
			}))
			lineNumber++

		case reflect.String:
			textType := ""
			if value, ok := params["type"]; ok {
				textType = strings.ToLower(value)
			}

			text := new(formText)
			text.data = data
			text.index = i
			text.id = id
			text.listener = listener

			//valueColumn := addLabel(id, textType == "multiline")
			editor := NewEditView(session, Params{
				ID:                   id,
				Text:                 val.String(),
				EditTextChangedEvent: text.changed,
			})

			if textType != "" {
				editor.Set(EditViewType, textType)
			}

			maxlen := ""
			if text, ok := params["max-length"]; ok {
				maxlen = text
			} else if text, ok := params["maxlength"]; ok {
				maxlen = text
			} else if text, ok := params["maxlen"]; ok {
				maxlen = text
			}

			if maxlen != "" {
				n, err := strconv.Atoi(maxlen)
				if err != nil {
					ErrorLog(err.Error())
				} else if n > 0 {
					editor.Set(MaxLength, n)
				} else {
					ErrorLog(`max-length <= 0`)
				}
			}

			controlView = editor

		//case reflect.Struct:
		//case reflect.Slice, reflect.Array:
		//case reflect.Map:

		default:
			DebugLog("unsupported element type") // TODO
		}

		if controlView != nil {
			label, hasLabel := params["label"]
			if !hasLabel {
				controlView.Set(Row, lineNumber)
				controlView.Set(Column, Range{First: 0, Last: 1})
				result.Append(controlView)
			} else {
				_, twoLine := params["two-line"]

				if twoLine { //} || alwaysTwoLine {
					result.Append(NewListLayout(session, Params{
						Row:         lineNumber,
						Column:      Range{First: 0, Last: 1},
						Orientation: VerticalOrientation,
						Content: []View{
							NewTextView(session, Params{
								ID:   id + "-label",
								Text: label,
							}),
							controlView,
						},
					}))
				} else {
					result.Append(NewTextView(session, Params{
						ID:     id + "-label",
						Row:    lineNumber,
						Column: 0,
						Text:   label,
					}))

					controlView.Set(Row, lineNumber)
					controlView.Set(Column, 1)
					result.Append(controlView)
				}
			}
			lineNumber++
		}
	}

	return result
}

func NewForm(session Session, data any, listener func(string)) View {

	if session == nil {
		ErrorLog(`nil session argument of NewForm function`)
		return nil
	}

	if data == nil {
		ErrorLog(`nil data argument of NewForm function`)
		return nil
	}

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		ErrorLog(`An argument to NewForm function can only be a structure or a pointer to a structure`)
		return nil
	}

	return createForm(session, data, value, listener)
}
