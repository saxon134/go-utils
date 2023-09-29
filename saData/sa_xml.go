package saData

import (
	"encoding/xml"
	"io"
)

type StringMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

/******** 外部接口 ********/

func MapToXml(m *map[string]string) string {
	buf, _ := xml.Marshal(StringMap(*m))
	return string(buf)
}

func XmlToMap(s string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := xml.Unmarshal([]byte(s), &m)
	return m, err
}

/******** 重写解析接口 ********/

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		if err = e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v}); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = StringMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}
