// Code generated by excelconfc. DO NOT EDIT.

package excelconf

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

type DateTime struct {
	time.Time
}

func (s *DateTime) UnmarshalJSON(b []byte) error {
	// 去掉引号
	str := string(b)
	str = str[1 : len(str)-1]

	// 解析时间字符串
	t, err := time.Parse(time.DateTime, str)
	if err != nil {
		return err
	}

	s.Time = t
	return nil
}

func (s *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Time.Format(time.DateTime))
}

func (s *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var dateStr string

	// 解析 XML 元素的字符数据到 dateStr
	if err := d.DecodeElement(&dateStr, &start); err != nil {
		return err
	}

	// 解析时间字符串，使用正确的格式
	parsedTime, err := time.Parse(time.DateTime, dateStr)
	if err != nil {
		return err
	}

	// 将解析后的时间赋值给 DateTime
	s.Time = parsedTime
	return nil
}
