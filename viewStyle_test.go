package rui

/*
import (
	"strings"
	"testing"
)

func TestViewStyleCreate(t *testing.T) {

	app := new(application)
	app.init("")
	session := newSession(app, 1, "", false, false)

	var style viewStyle
	style.init()

	data := []struct{ property, value string }{
		{Width, "100%"},
		{Height, "400px"},
		{Margin, "4px"},
		{Margin + "-bottom", "auto"},
		{Padding, "1em"},
		{Font, "Arial"},
		{BackgroundColor, "#FF008000"},
		{TextColor, "#FF000000"},
		{TextSize, "1.25em"},
		{TextWeight, "bold"},
		{TextAlign, "center"},
		{TextTransform, "uppercase"},
		{TextIndent, "0.25em"},
		{LetterSpacing, "1.5em"},
		{WordSpacing, "8px"},
		{LineHeight, "2em"},
		{Italic, "on"},
		{TextDecoration, "strikethrough | overline | underline"},
		{SmallCaps, "on"},
	}

	for _, prop := range data {
		style.Set(prop.property, prop.value)
	}

	style.AddShadow(NewViewShadow(SizeUnit{Auto, 0}, SizeUnit{Auto, 0}, Px(4), Px(6), 0xFF808080))

	expected := `width: 100%; height: 400px; font-size: 1.25rem; text-indent: 0.25rem; letter-spacing: 1.5rem; word-spacing: 8px; ` +
		`line-height: 2rem; padding: 1rem; margin-left: 4px; margin-top: 4px; margin-right: 4px; box-shadow: 0 0 4px 6px rgb(128,128,128); ` +
		`background-color: rgb(0,128,0); color: rgb(0,0,0); font-family: Arial; font-weight: bold; font-style: italic; font-variant: small-caps; ` +
		`text-align: center; text-decoration: line-through overline underline; text-transform: uppercase;`

	buffer := strings.Builder{}
	style.cssViewStyle(&buffer, session)
	if text := strings.Trim(buffer.String(), " "); text != expected {
		t.Error("\nresult  : " + text + "\nexpected: " + expected)
	}

		w := newCompactDataWriter()
		w.StartObject("_")
		style.writeStyle(w)
		w.FinishObject()
		expected2 := `_{width=100%,height=400px,margin="4px,4px,auto,4px",padding=1em,background-color=#FF008000,shadow=_{color=#FF808080,blur=4px,spread-radius=6px},font=Arial,text-color=#FF000000,text-size=1.25em,text-weight=bold,italic=on,small-caps=on,text-decoration=strikethrough|overline|underline,text-align=center,text-indent=0.25em,letter-spacing=1.5em,word-spacing=8px,line-height=2em,text-transform=uppercase}`

		if text := w.String(); text != expected2 {
			t.Error("\n result: " + text + "\nexpected: " + expected2)
		}

		var style1 viewStyle
		style1.init()
		if obj, err := ParseDataText(expected2); err == nil {
			style1.parseStyle(obj, new(sessionData))
			buffer.Reset()
			style.cssStyle(&buffer)
			if text := buffer.String(); text != expected {
				t.Error("\n result: " + text + "\nexpected: " + expected)
			}
		} else {
			t.Error(err)
		}

		var style2 viewStyle
		style2.init()

		style2.textWeight = 4
		style2.textAlign = RightAlign
		style2.textTransform = LowerCaseTextTransform
		style2.textDecoration = NoneDecoration
		style2.italic = Off
		style2.smallCaps = Off

		expected = `font-weight: normal; font-style: normal; font-variant: normal; text-align: right; text-decoration: none; text-transform: lowercase; `
		buffer.Reset()
		style2.cssStyle(&buffer)
		if text := buffer.String(); text != expected {
			t.Error("\n result: " + text + "\nexpected: " + expected)
		}

		w.Reset()
		w.StartObject("_")
		style2.writeStyle(w)
		w.FinishObject()
		expected = `_{text-weight=normal,italic=off,small-caps=off,text-decoration=none,text-align=right,text-transform=lowercase}`

		if text := w.String(); text != expected {
			t.Error("\n result: " + text + "\nexpected: " + expected)
		}

		style2.textWeight = 5
		style2.textAlign = JustifyTextAlign
		style2.textTransform = CapitalizeTextTransform
		style2.textDecoration = Inherit
		style2.italic = Inherit
		style2.smallCaps = Inherit

		expected = `font-weight: 500; text-align: justify; text-transform: capitalize; `
		buffer.Reset()
		style2.cssStyle(&buffer)
		if text := buffer.String(); text != expected {
			t.Error("\n  result: " + text + "\nexpected: " + expected)
		}

		w.Reset()
		w.StartObject("_")
		style2.writeStyle(w)
		w.FinishObject()
		expected = `_{text-weight=5,text-align=justify,text-transform=capitalize}`

		if text := w.String(); text != expected {
			t.Error("\n  result: " + text + "\nexpected: " + expected)
		}
}
*/
