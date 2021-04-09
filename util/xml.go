package util

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

/**
 * 根据xml字符串解析成map
 */
func EncodeXmlToMap(xmlStr string) map[string]string {

	params := make(map[string]string)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	var (
		key   string
		value string
	)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params[key] = value
			}
		}
	}

	return params
}

func EncodeXMLFromMap(w io.Writer, m map[string]string, rootname string) (err error) {
	switch v := w.(type) {
	case *bytes.Buffer:
		bufw := v
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return nil
	case *strings.Builder:
		bufw := v
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return nil

	case *bufio.Writer:
		bufw := v
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return bufw.Flush()

	default:
		bufw := bufio.NewWriterSize(w, 256)
		if err = bufw.WriteByte('<'); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}

		for k, v := range m {
			if err = bufw.WriteByte('<'); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}

			if err = xml.EscapeText(bufw, []byte(v)); err != nil {
				return
			}

			if _, err = bufw.WriteString("</"); err != nil {
				return
			}
			if _, err = bufw.WriteString(k); err != nil {
				return
			}
			if err = bufw.WriteByte('>'); err != nil {
				return
			}
		}

		if _, err = bufw.WriteString("</"); err != nil {
			return
		}
		if _, err = bufw.WriteString(rootname); err != nil {
			return
		}
		if err = bufw.WriteByte('>'); err != nil {
			return
		}
		return bufw.Flush()
	}
}
