package rui

import (
	"testing"
)

func TestParseDataText(t *testing.T) {

	SetErrorLog(func(text string) {
		t.Error(text)
	})

	text := `obj1 {
	key1 = val1,
	key2=obj2{
		key2.1=[val2.1,obj2.2{}, obj2.3{}],
		"key 2.2"='val 2.2'
		// Comment
		key2.3/* comment */ = {
		}
		/*
		Multiline comment
		*/
		'key2.4' = obj2.3{ text = " "},
		key2.5= [],
	},
	key3 = "\n \t \\ \r \" ' \X4F\x4e \U01Ea",` +
		"key4=`" + `\n \t \\ \r \" ' \x8F \UF80a` + "`\r}"

	obj := ParseDataText(text)
	if obj != nil {
		if obj.Tag() != "obj1" {
			t.Error(`obj.Tag() != "obj1"`)
		}
		if !obj.IsObject() {
			t.Error(`!obj.IsObject()`)
		}
		if obj.PropertyCount() != 4 {
			t.Error(`obj.PropertyCount() != 4`)
		}

		if obj.Property(-1) != nil {
			t.Error(`obj.Property(-1) != nil`)
		}

		if val, ok := obj.PropertyValue("key1"); !ok || val != "val1" {
			t.Errorf(`obj.PropertyValue("key1") result: ("%s",%v)`, val, ok)
		}

		if val, ok := obj.PropertyValue("key3"); !ok || val != "\n \t \\ \r \" ' \x4F\x4e \u01Ea" {
			t.Errorf(`obj.PropertyValue("key3") result: ("%s",%v)`, val, ok)
		}

		if val, ok := obj.PropertyValue("key4"); !ok || val != `\n \t \\ \r \" ' \x8F \UF80a` {
			t.Errorf(`obj.PropertyValue("key4") result: ("%s",%v)`, val, ok)
		}

		if o := obj.PropertyObject("key2"); o == nil {
			t.Error(`obj.PropertyObject("key2") == nil`)
		}

		if o := obj.PropertyObject("key1"); o != nil {
			t.Error(`obj.PropertyObject("key1") != nil`)
		}

		if o := obj.PropertyObject("key5"); o != nil {
			t.Error(`obj.PropertyObject("key5") != nil`)
		}

		if val, ok := obj.PropertyValue("key2"); ok {
			t.Errorf(`obj.PropertyValue("key2") result: ("%s",%v)`, val, ok)
		}

		if val, ok := obj.PropertyValue("key5"); ok {
			t.Errorf(`obj.PropertyValue("key5") result: ("%s",%v)`, val, ok)
		}

		testKey := func(obj DataObject, index int, tag string, nodeType DataNodeType) DataNode {
			key := obj.Property(index)
			if key == nil {
				t.Errorf(`%s.Property(%d) == nil`, obj.Tag(), index)
			} else {
				if key.Tag() != tag {
					t.Errorf(`%s.Property(%d).Tag() != "%s"`, obj.Tag(), index, tag)
				}

				if key.Type() != nodeType {
					switch nodeType {
					case TextNode:
						t.Errorf(`%s.Property(%d) is not text`, obj.Tag(), index)

					case ObjectNode:
						t.Errorf(`%s.Property(%d) is not object`, obj.Tag(), index)

					case ArrayNode:
						t.Errorf(`%s.Property(%d) is not array`, obj.Tag(), index)
					}
				}
			}

			return key
		}

		if key := testKey(obj, 0, "key1", TextNode); key != nil {
			if key.Text() != "val1" {
				t.Error(`key1.Value() != "val1"`)
			}
		}

		if key := testKey(obj, 1, "key2", ObjectNode); key != nil {
			o := key.Object()
			if o == nil {
				t.Error(`key2.Value().Object() == nil`)
			} else {
				if o.PropertyCount() != 5 {
					t.Error(`key2.Value().Object().PropertyCount() != 4`)
				}

				type testKeyData struct {
					tag      string
					nodeType DataNodeType
				}

				data := []testKeyData{
					{tag: "key2.1", nodeType: ArrayNode},
					{tag: "key 2.2", nodeType: TextNode},
					{tag: "key2.3", nodeType: ObjectNode},
					{tag: "key2.4", nodeType: ObjectNode},
					{tag: "key2.5", nodeType: ArrayNode},
				}

				for i, d := range data {
					testKey(o, i, d.tag, d.nodeType)
				}
			}
		}

		node1 := obj.Property(1)
		if node1 == nil {
			t.Error("obj.Property(1) != nil")
		} else if node1.Type() != ObjectNode {
			t.Error("obj.Property(1).Type() != ObjectNode")
		} else if obj := node1.Object(); obj != nil {
			if key := obj.Property(0); key != nil {
				if key.Type() != ArrayNode {
					t.Error("obj.Property(1).Object().Property(0)..Type() != ArrayNode")
				} else {
					if key.ArraySize() != 3 {
						t.Error("obj.Property(1).Object().Property(0).ArraySize() != 3")
					}

					if e := key.ArrayElement(0); e == nil {
						t.Error("obj.Property(1).Object().Property(0).ArrayElement(0) == nil")
					} else if e.IsObject() {
						t.Error("obj.Property(1).Object().Property(0).ArrayElement(0).IsObject() == true")
					}

					if e := key.ArrayElement(2); e == nil {
						t.Error("obj.Property(1).Object().Property(0).ArrayElement(2) == nil")
					} else if !e.IsObject() {
						t.Error("obj.Property(1).Object().Property(0).ArrayElement(2).IsObject() == false")
					} else if e.Value() != "" {
						t.Error(`obj.Property(1).Object().Property(0).ArrayElement(2).Value() != ""`)
					}

					if e := key.ArrayElement(3); e != nil {
						t.Error("obj.Property(1).Object().Property(0).ArrayElement(3) != nil")
					}
				}
			}
		} else {
			t.Error("obj.Property(1).Object() == nil")
		}
	}

	SetErrorLog(func(text string) {
	})

	failText := []string{
		" ",
		"obj[]",
		"obj={}",
		"obj{key}",
		"obj{key=}",
		"obj{key=val",
		"obj{key=obj2{}",
		"obj{key=obj2{key2}}",
		"obj{key=\"val}",
		"obj{key=val\"}",
		"obj{key=\"val`}",
		"obj{key=[}}",
		"obj{key=[val",
		"obj{key=[val,",
		"obj{key=[obj2{]",
		`obj{key="""}`,
		`obj{key="\z"}`,
		`obj{key="\xG6"}`,
		`obj{key="\uG678"}`,
		`obj{key="\x6"}`,
		`obj{key="\u678"}`,
		`obj{key1=val1 key2=val2}`,
		`obj{key=//"\u678"}`,
		`obj{key="\u678" /*}`,
	}

	for _, txt := range failText {
		if obj := ParseDataText(txt); obj != nil {
			t.Errorf("result ParseDataText(\"%s\") must be fail", txt)
		}
	}
}
