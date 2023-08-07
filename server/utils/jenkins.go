package utils

import (
	"encoding/xml"
	"strconv"
	"unicode"
)

func ChineseToHTMLEntity(entity rune) string {
	if unicode.Is(unicode.Han, entity) {
		return "&#" + strconv.FormatInt(int64(entity), 10) + ";"
	} else {
		return string(entity)
	}
}


func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}