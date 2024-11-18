package rui

import (
	"strconv"
	"strings"
)

// Constants for [ProgressBar] specific properties and events
const (
	// ProgressBarMax is the constant for "progress-max" property tag.
	//
	// Used by `ProgressBar`.
	// Maximum value, default is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	ProgressBarMax PropertyName = "progress-max"

	// ProgressBarValue is the constant for "progress-value" property tag.
	//
	// Used by `ProgressBar`.
	// Current value, default is 0.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	ProgressBarValue PropertyName = "progress-value"
)

// ProgressBar represents a ProgressBar view
type ProgressBar interface {
	View
}

type progressBarData struct {
	viewData
}

// NewProgressBar create new ProgressBar object and return it
func NewProgressBar(session Session, params Params) ProgressBar {
	view := new(progressBarData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newProgressBar(session Session) View {
	return new(progressBarData)
}

func (progress *progressBarData) init(session Session) {
	progress.viewData.init(session)
	progress.tag = "ProgressBar"
	progress.normalize = normalizeProgressBarTag
	progress.changed = progress.propertyChanged
}

func normalizeProgressBarTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Max, "progress-bar-max", "progressbar-max":
		return ProgressBarMax

	case Value, "progress-bar-value", "progressbar-value":
		return ProgressBarValue
	}
	return tag
}

func (progress *progressBarData) propertyChanged(tag PropertyName) {

	switch tag {
	case ProgressBarMax:
		progress.Session().updateProperty(progress.htmlID(), "max",
			strconv.FormatFloat(GetProgressBarMax(progress), 'f', -1, 32))

	case ProgressBarValue:
		progress.Session().updateProperty(progress.htmlID(), "value",
			strconv.FormatFloat(GetProgressBarValue(progress), 'f', -1, 32))

	default:
		progress.viewData.propertyChanged(tag)
	}
}

func (progress *progressBarData) htmlTag() string {
	return "progress"
}

func (progress *progressBarData) htmlProperties(self View, buffer *strings.Builder) {
	progress.viewData.htmlProperties(self, buffer)

	buffer.WriteString(` max="`)
	buffer.WriteString(strconv.FormatFloat(GetProgressBarMax(progress), 'f', -1, 64))
	buffer.WriteByte('"')

	buffer.WriteString(` value="`)
	buffer.WriteString(strconv.FormatFloat(GetProgressBarValue(progress), 'f', -1, 64))
	buffer.WriteByte('"')
}

// GetProgressBarMax returns the max value of ProgressBar subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetProgressBarMax(view View, subviewID ...string) float64 {
	return floatStyledProperty(view, subviewID, ProgressBarMax, 1)
}

// GetProgressBarValue returns the value of ProgressBar subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetProgressBarValue(view View, subviewID ...string) float64 {
	return floatStyledProperty(view, subviewID, ProgressBarValue, 0)
}
