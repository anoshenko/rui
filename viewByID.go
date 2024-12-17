package rui

import "strings"

// ViewByID returns the child View path to which is specified using the arguments id, ids. Example
//
//	view := ViewByID(rootView, "id1", "id2", "id3")
//	view := ViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found, the function will return nil
func ViewByID(rootView View, id string, ids ...string) View {
	if rootView == nil {
		ErrorLog(`ViewByID(nil, "` + id + `"): rootView is nil`)
		return nil
	}

	path := []string{id}
	if len(ids) > 0 {
		path = append(path, ids...)
	}

	result := rootView
	for _, id := range path {
		if result.ID() != id {
			if container, ok := result.(ParentView); ok {
				if view := viewByID(container, id); view != nil {
					result = view
				} else if index := strings.IndexRune(id, '/'); index > 0 {
					if view := ViewByID(result, id[:index], id[index+1:]); view != nil {
						result = view
					} else {
						ErrorLog(`ViewByID(_, "` + id + `"): View not found`)
						return nil
					}
				} else {
					ErrorLog(`ViewByID(_, "` + id + `"): View not found`)
					return nil
				}
			}
		}
	}
	return result
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

// ViewsContainerByID return the ViewsContainer path to which is specified using the arguments id, ids. Example
//
//	view := ViewsContainerByID(rootView, "id1", "id2", "id3")
//	view := ViewsContainerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ViewsContainer, the function will return nil
func ViewsContainerByID(rootView View, id string, ids ...string) ViewsContainer {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(ViewsContainer); ok {
			return list
		}
		ErrorLog(`ViewsContainerByID(_, "` + id + `"): The found View is not ViewsContainer`)
	}
	return nil
}

// ListLayoutByID return the ListLayout  path to which is specified using the arguments id, ids. Example
//
//	view := ListLayoutByID(rootView, "id1", "id2", "id3")
//	view := ListLayoutByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ListLayout, the function will return nil
func ListLayoutByID(rootView View, id string, ids ...string) ListLayout {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(ListLayout); ok {
			return list
		}
		ErrorLog(`ListLayoutByID(_, "` + id + `"): The found View is not ListLayout`)
	}
	return nil
}

// StackLayoutByID return the StackLayout path to which is specified using the arguments id, ids. Example
//
//	view := StackLayoutByID(rootView, "id1", "id2", "id3")
//	view := StackLayoutByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not StackLayout, the function will return nil
func StackLayoutByID(rootView View, id string, ids ...string) StackLayout {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(StackLayout); ok {
			return list
		}
		ErrorLog(`StackLayoutByID(_, "` + id + `"): The found View is not StackLayout`)
	}
	return nil
}

// GridLayoutByID return the GridLayout path to which is specified using the arguments id, ids. Example
//
//	view := GridLayoutByID(rootView, "id1", "id2", "id3")
//	view := GridLayoutByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not GridLayout, the function will return nil
func GridLayoutByID(rootView View, id string, ids ...string) GridLayout {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(GridLayout); ok {
			return list
		}
		ErrorLog(`GridLayoutByID(_, "` + id + `"): The found View is not GridLayout`)
	}
	return nil
}

// ColumnLayoutByID return the ColumnLayout path to which is specified using the arguments id, ids. Example
//
//	view := ColumnLayoutByID(rootView, "id1", "id2", "id3")
//	view := ColumnLayoutByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ColumnLayout, the function will return nil
func ColumnLayoutByID(rootView View, id string, ids ...string) ColumnLayout {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(ColumnLayout); ok {
			return list
		}
		ErrorLog(`ColumnLayoutByID(_, "` + id + `"): The found View is not ColumnLayout`)
	}
	return nil
}

// DetailsViewByID return the ColumnLayout path to which is specified using the arguments id, ids. Example
//
//	view := DetailsViewByID(rootView, "id1", "id2", "id3")
//	view := DetailsViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not DetailsView, the function will return nil
func DetailsViewByID(rootView View, id string, ids ...string) DetailsView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if details, ok := view.(DetailsView); ok {
			return details
		}
		ErrorLog(`DetailsViewByID(_, "` + id + `"): The found View is not DetailsView`)
	}
	return nil
}

// DropDownListByID return the DropDownListView path to which is specified using the arguments id, ids. Example
//
//	view := DropDownListByID(rootView, "id1", "id2", "id3")
//	view := DropDownListByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not DropDownList, the function will return nil
func DropDownListByID(rootView View, id string, ids ...string) DropDownList {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(DropDownList); ok {
			return list
		}
		ErrorLog(`DropDownListByID(_, "` + id + `"): The found View is not DropDownList`)
	}
	return nil
}

// TabsLayoutByID return the TabsLayout path to which is specified using the arguments id, ids. Example
//
//	view := TabsLayoutByID(rootView, "id1", "id2", "id3")
//	view := TabsLayoutByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not TabsLayout, the function will return nil
func TabsLayoutByID(rootView View, id string, ids ...string) TabsLayout {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(TabsLayout); ok {
			return list
		}
		ErrorLog(`TabsLayoutByID(_, "` + id + `"): The found View is not TabsLayout`)
	}
	return nil
}

// ListViewByID return the ListView path to which is specified using the arguments id, ids. Example
//
//	view := ListViewByID(rootView, "id1", "id2", "id3")
//	view := ListViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ListView, the function will return nil
func ListViewByID(rootView View, id string, ids ...string) ListView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if list, ok := view.(ListView); ok {
			return list
		}
		ErrorLog(`ListViewByID(_, "` + id + `"): The found View is not ListView`)
	}
	return nil
}

// TextViewByID return the TextView path to which is specified using the arguments id, ids. Example
//
//	view := TextViewByID(rootView, "id1", "id2", "id3")
//	view := TextViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not TextView, the function will return nil
func TextViewByID(rootView View, id string, ids ...string) TextView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if text, ok := view.(TextView); ok {
			return text
		}
		ErrorLog(`TextViewByID(_, "` + id + `"): The found View is not TextView`)
	}
	return nil
}

// ButtonByID return the Button path to which is specified using the arguments id, ids. Example
//
//	view := ButtonByID(rootView, "id1", "id2", "id3")
//	view := ButtonByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not Button, the function will return nil
func ButtonByID(rootView View, id string, ids ...string) Button {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if button, ok := view.(Button); ok {
			return button
		}
		ErrorLog(`ButtonByID(_, "` + id + `"): The found View is not Button`)
	}
	return nil
}

// CheckboxByID return the Checkbox path to which is specified using the arguments id, ids. Example
//
//	view := CheckboxByID(rootView, "id1", "id2", "id3")
//	view := CheckboxByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not Checkbox, the function will return nil
func CheckboxByID(rootView View, id string, ids ...string) Checkbox {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if checkbox, ok := view.(Checkbox); ok {
			return checkbox
		}
		ErrorLog(`CheckboxByID(_, "` + id + `"): The found View is not Checkbox`)
	}
	return nil
}

// EditViewByID return the EditView path to which is specified using the arguments id, ids. Example
//
//	view := EditViewByID(rootView, "id1", "id2", "id3")
//	view := EditViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not EditView, the function will return nil
func EditViewByID(rootView View, id string, ids ...string) EditView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if buttons, ok := view.(EditView); ok {
			return buttons
		}
		ErrorLog(`EditViewByID(_, "` + id + `"): The found View is not EditView`)
	}
	return nil
}

// ProgressBarByID return the ProgressBar path to which is specified using the arguments id, ids. Example
//
//	view := ProgressBarByID(rootView, "id1", "id2", "id3")
//	view := ProgressBarByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ProgressBar, the function will return nil
func ProgressBarByID(rootView View, id string, ids ...string) ProgressBar {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if buttons, ok := view.(ProgressBar); ok {
			return buttons
		}
		ErrorLog(`ProgressBarByID(_, "` + id + `"): The found View is not ProgressBar`)
	}
	return nil
}

// ColorPickerByID return the ColorPicker path to which is specified using the arguments id, ids. Example
//
//	view := ColorPickerByID(rootView, "id1", "id2", "id3")
//	view := ColorPickerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ColorPicker, the function will return nil
func ColorPickerByID(rootView View, id string, ids ...string) ColorPicker {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if input, ok := view.(ColorPicker); ok {
			return input
		}
		ErrorLog(`ColorPickerByID(_, "` + id + `"): The found View is not ColorPicker`)
	}
	return nil
}

// NumberPickerByID return the NumberPicker path to which is specified using the arguments id, ids. Example
//
//	view := NumberPickerByID(rootView, "id1", "id2", "id3")
//	view := NumberPickerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not NumberPicker, the function will return nil
func NumberPickerByID(rootView View, id string, ids ...string) NumberPicker {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if input, ok := view.(NumberPicker); ok {
			return input
		}
		ErrorLog(`NumberPickerByID(_, "` + id + `"): The found View is not NumberPicker`)
	}
	return nil
}

// TimePickerByID return the TimePicker path to which is specified using the arguments id, ids. Example
//
//	view := TimePickerByID(rootView, "id1", "id2", "id3")
//	view := TimePickerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not TimePicker, the function will return nil
func TimePickerByID(rootView View, id string, ids ...string) TimePicker {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if input, ok := view.(TimePicker); ok {
			return input
		}
		ErrorLog(`TimePickerByID(_, "` + id + `"): The found View is not TimePicker`)
	}
	return nil
}

// DatePickerByID return the DatePicker path to which is specified using the arguments id, ids. Example
//
//	view := DatePickerByID(rootView, "id1", "id2", "id3")
//	view := DatePickerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not DatePicker, the function will return nil
func DatePickerByID(rootView View, id string, ids ...string) DatePicker {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if input, ok := view.(DatePicker); ok {
			return input
		}
		ErrorLog(`DatePickerByID(_, "` + id + `"): The found View is not DatePicker`)
	}
	return nil
}

// FilePickerByID return the FilePicker path to which is specified using the arguments id, ids. Example
//
//	view := FilePickerByID(rootView, "id1", "id2", "id3")
//	view := FilePickerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not FilePicker, the function will return nil
func FilePickerByID(rootView View, id string, ids ...string) FilePicker {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if input, ok := view.(FilePicker); ok {
			return input
		}
		ErrorLog(`FilePickerByID(_, "` + id + `"): The found View is not FilePicker`)
	}
	return nil
}

// CanvasViewByID return the CanvasView path to which is specified using the arguments id, ids. Example
//
//	view := CanvasViewByID(rootView, "id1", "id2", "id3")
//	view := CanvasViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not CanvasView, the function will return nil
func CanvasViewByID(rootView View, id string, ids ...string) CanvasView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if canvas, ok := view.(CanvasView); ok {
			return canvas
		}
		ErrorLog(`CanvasViewByID(_, "` + id + `"): The found View is not CanvasView`)
	}
	return nil
}

// AudioPlayerByID return the AudioPlayer path to which is specified using the arguments id, ids. Example
//
//	view := AudioPlayerByID(rootView, "id1", "id2", "id3")
//	view := AudioPlayerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not AudioPlayer, the function will return nil
func AudioPlayerByID(rootView View, id string, ids ...string) AudioPlayer {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if canvas, ok := view.(AudioPlayer); ok {
			return canvas
		}
		ErrorLog(`AudioPlayerByID(_, "` + id + `"): The found View is not AudioPlayer`)
	}
	return nil
}

// VideoPlayerByID return the VideoPlayer path to which is specified using the arguments id, ids. Example
//
//	view := VideoPlayerByID(rootView, "id1", "id2", "id3")
//	view := VideoPlayerByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not VideoPlayer, the function will return nil
func VideoPlayerByID(rootView View, id string, ids ...string) VideoPlayer {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if canvas, ok := view.(VideoPlayer); ok {
			return canvas
		}
		ErrorLog(`VideoPlayerByID(_, "` + id + `"): The found View is not VideoPlayer`)
	}
	return nil
}

// ImageViewByID return the ImageView path to which is specified using the arguments id, ids. Example
//
//	view := ImageViewByID(rootView, "id1", "id2", "id3")
//	view := ImageViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not ImageView, the function will return nil
func ImageViewByID(rootView View, id string, ids ...string) ImageView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if canvas, ok := view.(ImageView); ok {
			return canvas
		}
		ErrorLog(`ImageViewByID(_, "` + id + `"): The found View is not ImageView`)
	}
	return nil
}

// TableViewByID return the TableView path to which is specified using the arguments id, ids. Example
//
//	view := TableViewByID(rootView, "id1", "id2", "id3")
//	view := TableViewByID(rootView, "id1/id2/id3")
//
// These two function calls are equivalent.
// If a View with this path is not found or View is not TableView, the function will return nil
func TableViewByID(rootView View, id string, ids ...string) TableView {
	if view := ViewByID(rootView, id, ids...); view != nil {
		if canvas, ok := view.(TableView); ok {
			return canvas
		}
		ErrorLog(`TableViewByID(_, "` + id + `"): The found View is not TableView`)
	}
	return nil
}
