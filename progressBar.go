package rui

import (
	"strconv"
	"strings"
)

const (
	ProgressBarMax   = "progress-max"
	ProgressBarValue = "progress-value"
)

// ProgressBar - ProgressBar view
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
	return NewProgressBar(session, nil)
}

func (progress *progressBarData) init(session Session) {
	progress.viewData.init(session)
	progress.tag = "ProgressBar"
}

func (progress *progressBarData) String() string {
	return getViewString(progress)
}

func (progress *progressBarData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Max, "progress-bar-max", "progressbar-max":
		return ProgressBarMax

	case Value, "progress-bar-value", "progressbar-value":
		return ProgressBarValue
	}
	return tag
}

func (progress *progressBarData) Remove(tag string) {
	progress.remove(progress.normalizeTag(tag))
}

func (progress *progressBarData) remove(tag string) {
	progress.viewData.remove(tag)
	progress.propertyChanged(tag)
}

func (progress *progressBarData) propertyChanged(tag string) {
	if progress.created {
		switch tag {
		case ProgressBarMax:
			updateProperty(progress.htmlID(), Max, strconv.FormatFloat(GetProgressBarMax(progress), 'f', -1, 32), progress.session)

		case ProgressBarValue:
			updateProperty(progress.htmlID(), Value, strconv.FormatFloat(GetProgressBarValue(progress), 'f', -1, 32), progress.session)
		}
	}
}

func (progress *progressBarData) Set(tag string, value any) bool {
	return progress.set(progress.normalizeTag(tag), value)
}

func (progress *progressBarData) set(tag string, value any) bool {
	if progress.viewData.set(tag, value) {
		progress.propertyChanged(tag)
		return true
	}
	return false
}

func (progress *progressBarData) Get(tag string) any {
	return progress.get(progress.normalizeTag(tag))
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
