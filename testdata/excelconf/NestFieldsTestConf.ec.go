// Code generated by excelconfc. DO NOT EDIT.
// source:
//   file: ExcelConfTest.xlsx
//   sheet: NestFieldsTestConf

package excelconf

import (
	"encoding/json"
	"encoding/xml"
	"os"
)

type TypeAA struct {
	Aa1 int32  `json:"aa1,omitempty" xml:"aa1"`
	Aa2 string `json:"aa2,omitempty" xml:"aa2"`
}

type TypeA struct {
	A1 int32    `json:"a1,omitempty" xml:"a1"`
	A2 int32    `json:"a2,omitempty" xml:"a2"`
	AA []TypeAA `json:"AA,omitempty" xml:"AA>item"`
	A3 string   `json:"a3,omitempty" xml:"a3"`
}

type TypeB struct {
	B1 int32  `json:"b1,omitempty" xml:"b1"`
	B2 int32  `json:"b2,omitempty" xml:"b2"`
	B3 string `json:"b3,omitempty" xml:"b3"`
}

type NestFieldsTestConf struct {
	Id int32   `json:"id,omitempty" xml:"id"`
	A  []TypeA `json:"A,omitempty" xml:"A>item"`
	B  TypeB   `json:"B,omitempty" xml:"B"`
	C  int64   `json:"C,omitempty" xml:"C"`
}

type NestFieldsTestConfMap map[int32]*NestFieldsTestConf

func (s NestFieldsTestConfMap) LoadFromJsonFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	jsonData := struct {
		Data []*NestFieldsTestConf `json:"data"`
	}{}
	if err := json.Unmarshal(fileBytes, &jsonData); err != nil {
		return err
	}
	for _, conf := range jsonData.Data {
		s[conf.Id] = conf
	}
	return nil
}

func (s NestFieldsTestConfMap) LoadFromXmlFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	xmlData := struct {
		Data []*NestFieldsTestConf `xml:"all_infos>item"`
	}{}
	if err := xml.Unmarshal(fileBytes, &xmlData); err != nil {
		return err
	}
	for _, conf := range xmlData.Data {
		s[conf.Id] = conf
	}
	return nil
}

func (s NestFieldsTestConfMap) GetVal(key int32) *NestFieldsTestConf {
	return s[key]
}

var instanceNestFieldsTestConfMap = make(NestFieldsTestConfMap)

func GetNestFieldsTestConfMapInst() NestFieldsTestConfMap {
	return instanceNestFieldsTestConfMap
}
