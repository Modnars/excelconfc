// Code generated by excelconfc. DO NOT EDIT.
// source:
//   file: ExcelConfTest.xlsx
//   sheet: ActTaskConf

package excelconf

import (
	"encoding/json"
	"os"
)

type AwardItem struct {
	ItemId uint32 `json:"item_id,omitempty"`
	Num    uint32 `json:"num,omitempty"`
}

type ActTaskConf struct {
	TaskId     uint32      `json:"task_id,omitempty"`
	Title      string      `json:"title,omitempty"`
	TargetId   uint32      `json:"target_id,omitempty"`
	TargetType uint32      `json:"target_type,omitempty"`
	Progress   uint32      `json:"progress,omitempty"`
	Awards     []AwardItem `json:"awards,omitempty"`
}

type ActTaskConfMap map[uint32]*ActTaskConf

func (s ActTaskConfMap) LoadFromJsonFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	jsonData := struct {
		Data []ActTaskConf `json:"data"`
	}{}
	json.Unmarshal(fileBytes, &jsonData)
	for _, conf := range jsonData.Data {
		s[conf.TaskId] = &conf
	}
	return nil
}

func (s ActTaskConfMap) GetVal(key uint32) *ActTaskConf {
	return s[key]
}

var instanceActTaskConfMap = make(ActTaskConfMap)

func GetActTaskConfMapInst() ActTaskConfMap {
	return instanceActTaskConfMap
}