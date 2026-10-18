package rui

import (
	"testing"
)

func TestParseDataText(t *testing.T) {

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

	obj, err := ParseDataText(text)
	if err != nil {
		t.Error(err)
	} else {
		if obj.Tag() != "obj1" {
			t.Error(`obj.Tag() != "obj1"`)
		}
		if obj.Type() != ObjectDataValue {
			t.Error(`obj.Type() != ObjectDataValue`)
		}
		if obj.PropertyCount() != 4 {
			t.Error(`obj.PropertyCount() != 4`)
		}

		if val, ok := obj.PropertyText("key1"); !ok || val != "val1" {
			t.Errorf(`obj.PropertyText("key1") result: ("%s",%v)`, val, ok)
		}

		if val, ok := obj.PropertyText("key3"); !ok || val != "\n \t \\ \r \" ' \x4F\x4e \u01Ea" {
			t.Errorf(`obj.PropertyText("key3") result: ("%s",%v)`, val, ok)
		}

		if val, ok := obj.PropertyText("key4"); !ok || val != `\n \t \\ \r \" ' \x8F \UF80a` {
			t.Errorf(`obj.PropertyText("key4") result: ("%s",%v)`, val, ok)
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

		if val, ok := obj.PropertyText("key2"); ok {
			t.Errorf(`obj.PropertyText("key2") result: ("%s",%v)`, val, ok)
		}

		if val, ok := obj.PropertyText("key5"); ok {
			t.Errorf(`obj.PropertyText("key5") result: ("%s",%v)`, val, ok)
		}

		testKey := func(obj DataObject, tag string, nodeType DataValueType) DataValue {
			key := obj.Property(tag)
			if key == nil {
				t.Errorf(`%s.Property(%s) == nil`, obj.Tag(), tag)
			} else if key.Type() != nodeType {
				switch nodeType {
				case TextDataValue:
					t.Errorf(`%s.Property(%s) is not text`, obj.Tag(), tag)

				case ObjectDataValue:
					t.Errorf(`%s.Property(%s) is not object`, obj.Tag(), tag)

				case ArrayDataValue:
					t.Errorf(`%s.Property(%s) is not array`, obj.Tag(), tag)
				}
			}

			return key
		}

		if key := testKey(obj, "key1", TextDataValue); key != nil {
			if key.String() != "val1" {
				t.Error(`key1.Value() != "val1"`)
			}
		}

		if key := testKey(obj, "key2", ObjectDataValue); key != nil {
			o := key.Object()
			if o == nil {
				t.Error(`key2.Value().Object() == nil`)
			} else {
				if o.PropertyCount() != 5 {
					t.Error(`key2.Value().Object().PropertyCount() != 4`)
				}

				type testKeyData struct {
					tag      string
					nodeType DataValueType
				}

				data := []testKeyData{
					{tag: "key2.1", nodeType: ArrayDataValue},
					{tag: "key 2.2", nodeType: TextDataValue},
					{tag: "key2.3", nodeType: ObjectDataValue},
					{tag: "key2.4", nodeType: ObjectDataValue},
					{tag: "key2.5", nodeType: ArrayDataValue},
				}

				for _, d := range data {
					testKey(o, d.tag, d.nodeType)
				}
			}
		}

		node1 := obj.Property("key2")
		if node1 == nil {
			t.Error("obj.Property(\"key2\") == nil")
		} else if node1.Type() != ObjectDataValue {
			t.Error("obj.Property(\"key2\").Type() != ObjectNode")
		} else if obj := node1.Object(); obj != nil {
			if key := obj.Property("key2.1"); key != nil {
				if key.Type() != ArrayDataValue {
					t.Error("obj.Property(\"key2\").Object().Property(\"key2.1\").Type() != ArrayNode")
				} else {
					array := key.Array()
					if array == nil {
						t.Error("obj.Property(\"key2\").Object().Property(\"key2.1\").Array() == nil")
					} else {

						if array.Size() != 3 {
							t.Error("obj.Property(\"key2\").Object().Property(\"key2.1\").Array().Size() != 3")
						}

						for i, tp := range []DataValueType{TextDataValue, ObjectDataValue, ObjectDataValue} {
							if e := array.Element(i); e == nil {
								t.Error("obj.Property(\"key2\").Object().Property(\"key2.1\").Array().Element(0) == nil")
							} else if e.Type() != tp {
								t.Errorf("obj.Property(\"key2\").Object().Property(\"key2.1\").Array().Element(0).Type() != %d", tp)
							}
						}

						if e := array.Element(3); e != nil {
							t.Error("obj.Property(\"key2\").Object().Property(\"key2.1\").Array().Element(3) != nil")
						}
					}
				}
			}
		} else {
			t.Error("obj.Property(\"key 2.2\").Object() == nil")
		}
	}

	expected := `obj1{
	key1 = val1,
	key2 = obj2{
		key2.1 = [
			val2.1,
			obj2.2{},
			obj2.3{},
		],
		"key 2.2" = "val 2.2",
		key2.3 = _{},
		key2.4 = obj2.3{
			text = " ",
		},
		key2.5 = [],
	},
	key3 = "\n \t \\ \r \" ' ON Ǫ",
	key4 = "\\n \\t \\\\ \\r \\\" ' \\x8F \\UF80a",
}`
	if expected != obj.String() {
		t.Error("obj.String() result fail")
	}

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
		if _, err := ParseDataText(txt); err == nil {
			t.Errorf("result ParseDataText(\"%s\") must be fail", txt)
		}
	}
}
