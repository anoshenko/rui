package rui

import "strings"

// ViewByID return a View with id equal to the argument of the function or nil if there is no such View
func ViewByID(rootView View, id string) View {
	if rootView == nil {
		ErrorLog(`ViewByID(nil, "` + id + `"): rootView is nil`)
		return nil
	}
	if rootView.ID() == id {
		return rootView
	}

	if container, ok := rootView.(ParentView); ok {
		if view := viewByID(container, id); view != nil {
			return view
		}
	}

	if index := strings.IndexRune(id, '/'); index > 0 {
		if view2 := ViewByID(rootView, id[:index]); view2 != nil {
			if view := ViewByID(view2, id[index+1:]); view != nil {
				return view
			}
			return nil
		}
		return nil
	}

	ErrorLog(`ViewByID(_, "` + id + `"): View not found`)
	return nil
}

func viewByID(rootView ParentView, id string) View {
	for _, view := range rootView.Views() {
		if view != nil {
			if view.ID() == id {
				return view
			}
			if container, ok := view.(ParentView); ok {
				if v := viewByID(container, id); v != nil {
					return v
				}
			}
		}
	}

	return nil
}

// ViewsContainerByID return a ViewsContainer with id equal to the argument of the function or
// nil if there is no such View or View is not ViewsContainer
func ViewsContainerByID(rootView View, id string) ViewsContainer {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(ViewsContainer); ok {
			return list
		}
		ErrorLog(`ViewsContainerByID(_, "` + id + `"): The found View is not ViewsContainer`)
	}
	return nil
}

// ListLayoutByID return a ListLayout with id equal to the argument of the function or
// nil if there is no such View or View is not ListLayout
func ListLayoutByID(rootView View, id string) ListLayout {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(ListLayout); ok {
			return list
		}
		ErrorLog(`ListLayoutByID(_, "` + id + `"): The found View is not ListLayout`)
	}
	return nil
}

// StackLayoutByID return a StackLayout with id equal to the argument of the function or
// nil if there is no such View or View is not StackLayout
func StackLayoutByID(rootView View, id string) StackLayout {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(StackLayout); ok {
			return list
		}
		ErrorLog(`StackLayoutByID(_, "` + id + `"): The found View is not StackLayout`)
	}
	return nil
}

// GridLayoutByID return a GridLayout with id equal to the argument of the function or
// nil if there is no such View or View is not GridLayout
func GridLayoutByID(rootView View, id string) GridLayout {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(GridLayout); ok {
			return list
		}
		ErrorLog(`GridLayoutByID(_, "` + id + `"): The found View is not GridLayout`)
	}
	return nil
}

// ColumnLayoutByID return a ColumnLayout with id equal to the argument of the function or
// nil if there is no such View or View is not ColumnLayout
func ColumnLayoutByID(rootView View, id string) ColumnLayout {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(ColumnLayout); ok {
			return list
		}
		ErrorLog(`ColumnLayoutByID(_, "` + id + `"): The found View is not ColumnLayout`)
	}
	return nil
}

// DetailsViewByID return a ColumnLayout with id equal to the argument of the function or
// nil if there is no such View or View is not DetailsView
func DetailsViewByID(rootView View, id string) DetailsView {
	if view := ViewByID(rootView, id); view != nil {
		if details, ok := view.(DetailsView); ok {
			return details
		}
		ErrorLog(`DetailsViewByID(_, "` + id + `"): The found View is not DetailsView`)
	}
	return nil
}

// DropDownListByID return a DropDownListView with id equal to the argument of the function or
// nil if there is no such View or View is not DropDownListView
func DropDownListByID(rootView View, id string) DropDownList {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(DropDownList); ok {
			return list
		}
		ErrorLog(`DropDownListByID(_, "` + id + `"): The found View is not DropDownList`)
	}
	return nil
}

// TabsLayoutByID return a TabsLayout with id equal to the argument of the function or
// nil if there is no such View or View is not TabsLayout
func TabsLayoutByID(rootView View, id string) TabsLayout {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(TabsLayout); ok {
			return list
		}
		ErrorLog(`TabsLayoutByID(_, "` + id + `"): The found View is not TabsLayout`)
	}
	return nil
}

// ListViewByID return a ListView with id equal to the argument of the function or
// nil if there is no such View or View is not ListView
func ListViewByID(rootView View, id string) ListView {
	if view := ViewByID(rootView, id); view != nil {
		if list, ok := view.(ListView); ok {
			return list
		}
		ErrorLog(`ListViewByID(_, "` + id + `"): The found View is not ListView`)
	}
	return nil
}

// TextViewByID return a TextView with id equal to the argument of the function or
// nil if there is no such View or View is not TextView
func TextViewByID(rootView View, id string) TextView {
	if view := ViewByID(rootView, id); view != nil {
		if text, ok := view.(TextView); ok {
			return text
		}
		ErrorLog(`TextViewByID(_, "` + id + `"): The found View is not TextView`)
	}
	return nil
}

// ButtonByID return a Button with id equal to the argument of the function or
// nil if there is no such View or View is not Button
func ButtonByID(rootView View, id string) Button {
	if view := ViewByID(rootView, id); view != nil {
		if button, ok := view.(Button); ok {
			return button
		}
		ErrorLog(`ButtonByID(_, "` + id + `"): The found View is not Button`)
	}
	return nil
}

// CheckboxByID return a Checkbox with id equal to the argument of the function or
// nil if there is no such View or View is not Checkbox
func CheckboxByID(rootView View, id string) Checkbox {
	if view := ViewByID(rootView, id); view != nil {
		if checkbox, ok := view.(Checkbox); ok {
			return checkbox
		}
		ErrorLog(`CheckboxByID(_, "` + id + `"): The found View is not Checkbox`)
	}
	return nil
}

// EditViewByID return a EditView with id equal to the argument of the function or
// nil if there is no such View or View is not EditView
func EditViewByID(rootView View, id string) EditView {
	if view := ViewByID(rootView, id); view != nil {
		if buttons, ok := view.(EditView); ok {
			return buttons
		}
		ErrorLog(`EditViewByID(_, "` + id + `"): The found View is not EditView`)
	}
	return nil
}

// ProgressBarByID return a ProgressBar with id equal to the argument of the function or
// nil if there is no such View or View is not ProgressBar
func ProgressBarByID(rootView View, id string) ProgressBar {
	if view := ViewByID(rootView, id); view != nil {
		if buttons, ok := view.(ProgressBar); ok {
			return buttons
		}
		ErrorLog(`ProgressBarByID(_, "` + id + `"): The found View is not ProgressBar`)
	}
	return nil
}

// ColorPickerByID return a ColorPicker with id equal to the argument of the function or
// nil if there is no such View or View is not ColorPicker
func ColorPickerByID(rootView View, id string) ColorPicker {
	if view := ViewByID(rootView, id); view != nil {
		if input, ok := view.(ColorPicker); ok {
			return input
		}
		ErrorLog(`ColorPickerByID(_, "` + id + `"): The found View is not ColorPicker`)
	}
	return nil
}

// NumberPickerByID return a NumberPicker with id equal to the argument of the function or
// nil if there is no such View or View is not NumberPicker
func NumberPickerByID(rootView View, id string) NumberPicker {
	if view := ViewByID(rootView, id); view != nil {
		if input, ok := view.(NumberPicker); ok {
			return input
		}
		ErrorLog(`NumberPickerByID(_, "` + id + `"): The found View is not NumberPicker`)
	}
	return nil
}

// TimePickerByID return a TimePicker with id equal to the argument of the function or
// nil if there is no such View or View is not TimePicker
func TimePickerByID(rootView View, id string) TimePicker {
	if view := ViewByID(rootView, id); view != nil {
		if input, ok := view.(TimePicker); ok {
			return input
		}
		ErrorLog(`TimePickerByID(_, "` + id + `"): The found View is not TimePicker`)
	}
	return nil
}

// DatePickerByID return a DatePicker with id equal to the argument of the function or
// nil if there is no such View or View is not DatePicker
func DatePickerByID(rootView View, id string) DatePicker {
	if view := ViewByID(rootView, id); view != nil {
		if input, ok := view.(DatePicker); ok {
			return input
		}
		ErrorLog(`DatePickerByID(_, "` + id + `"): The found View is not DatePicker`)
	}
	return nil
}

// FilePickerByID return a FilePicker with id equal to the argument of the function or
// nil if there is no such View or View is not FilePicker
func FilePickerByID(rootView View, id string) FilePicker {
	if view := ViewByID(rootView, id); view != nil {
		if input, ok := view.(FilePicker); ok {
			return input
		}
		ErrorLog(`FilePickerByID(_, "` + id + `"): The found View is not FilePicker`)
	}
	return nil
}

// CanvasViewByID return a CanvasView with id equal to the argument of the function or
// nil if there is no such View or View is not CanvasView
func CanvasViewByID(rootView View, id string) CanvasView {
	if view := ViewByID(rootView, id); view != nil {
		if canvas, ok := view.(CanvasView); ok {
			return canvas
		}
		ErrorLog(`CanvasViewByID(_, "` + id + `"): The found View is not CanvasView`)
	}
	return nil
}

/*
// TableViewByID return a TableView with id equal to the argument of the function or
// nil if there is no such View or View is not TableView
func TableViewByID(rootView View, id string) TableView {
	if view := ViewByID(rootView, id); view != nil {
		if canvas, ok := view.(TableView); ok {
			return canvas
		}
		ErrorLog(`TableViewByID(_, "` + id + `"): The found View is not TableView`)
	}
	return nil
}
*/

// AudioPlayerByID return a AudioPlayer with id equal to the argument of the function or
// nil if there is no such View or View is not AudioPlayer
func AudioPlayerByID(rootView View, id string) AudioPlayer {
	if view := ViewByID(rootView, id); view != nil {
		if canvas, ok := view.(AudioPlayer); ok {
			return canvas
		}
		ErrorLog(`AudioPlayerByID(_, "` + id + `"): The found View is not AudioPlayer`)
	}
	return nil
}

// VideoPlayerByID return a VideoPlayer with id equal to the argument of the function or
// nil if there is no such View or View is not VideoPlayer
func VideoPlayerByID(rootView View, id string) VideoPlayer {
	if view := ViewByID(rootView, id); view != nil {
		if canvas, ok := view.(VideoPlayer); ok {
			return canvas
		}
		ErrorLog(`VideoPlayerByID(_, "` + id + `"): The found View is not VideoPlayer`)
	}
	return nil
}

// ImageViewByID return a ImageView with id equal to the argument of the function or
// nil if there is no such View or View is not ImageView
func ImageViewByID(rootView View, id string) ImageView {
	if view := ViewByID(rootView, id); view != nil {
		if canvas, ok := view.(ImageView); ok {
			return canvas
		}
		ErrorLog(`ImageViewByID(_, "` + id + `"): The found View is not ImageView`)
	}
	return nil
}

// TableViewByID return a TableView with id equal to the argument of the function or
// nil if there is no such View or View is not TableView
func TableViewByID(rootView View, id string) TableView {
	if view := ViewByID(rootView, id); view != nil {
		if canvas, ok := view.(TableView); ok {
			return canvas
		}
		ErrorLog(`TableViewByID(_, "` + id + `"): The found View is not TableView`)
	}
	return nil
}
