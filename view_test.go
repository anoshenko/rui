package rui

/*
import (
	"testing"
)

func testViewCreate(t *testing.T, session Session, viewText string) View {
	if obj, err := ParseDataText(viewText); err == nil {
		if view := CreateView(session, obj); view != nil {
			writer := newCompactDataWriter()
			WriteViewData(writer, view)
			if str := writer.String(); str != viewText {
				t.Errorf("\n  result: \"%s\"\nexpected: \"%s\"", str, viewText)
			}
			return view
		}
		t.Errorf("CreateView(`%s`) == nil", viewText)

	} else {
		t.Error(err)
	}
	return nil
}

func TestViewCreate(t *testing.T) {

	testView1 := `View{id=View1,width=100%,height=20cm,margin="0px,0.8in,0px,16mm",padding="10px,8px,12px,16px",visibility=invisible,p=x}`
	session := newSession(nil, 0, "", false, false)

	if obj, err := ParseDataText(testView1); err == nil {
		if view := CreateView(session, obj); view != nil {
			//view.ParseProperties(obj)
			if view.ID() != "View1" {
				t.Errorf(`view.ID() != "%s"`, view.ID())
			}
			if view.Tag() != "View" {
				t.Errorf(`view.Tag() != "%s"`, view.Tag())
			}
			if view.Width() != Percent(100) {
				t.Errorf(`view.Width() == "%s"`, view.Width().String())
			}
			if view.Height() != Cm(20) {
				t.Errorf(`view.Height() == "%s"`, view.Height().String())
			}
			if view.Visibility() != Invisible {
				t.Error(`view.Visibility() != Invisible`)
			}
		}
	} else {
		t.Error(err)
	}

	testViewCreate(t, session, testView1)
}
*/
