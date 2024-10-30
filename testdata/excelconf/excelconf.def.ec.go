// Code generated by excelconfc. DO NOT EDIT.

package excelconf

import (
	"encoding/json"
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

func (s DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Time.Format(time.DateTime))
}