// Code generated by excelconfc. DO NOT EDIT.
// source:
//   file: ExcelConfTest.xlsx
//   sheet: GroupFlagTestConf

package excelconf

import (
	"encoding/json"
	"encoding/xml"
	"os"
)

type TaskRef struct {
	TaskId int32 `json:"task_id,omitempty" xml:"task_id"`
}

type ClientTaskResource struct {
	TaskId    int32  `json:"task_id,omitempty" xml:"task_id"`
	NewTaskId int32  `json:"new_task_id,omitempty" xml:"new_task_id"`
	Desc      string `json:"desc,omitempty" xml:"desc"`
}

type GroupFlagTestConf struct {
	Id            int32              `json:"id,omitempty" xml:"id"`
	TaskRefVec    []TaskRef          `json:"task_ref_vec,omitempty" xml:"task_ref_vec>item"`
	HotUpdateItem ClientTaskResource `json:"hot_update_item,omitempty" xml:"hot_update_item"`
	ActId         int64              `json:"act_id,omitempty" xml:"act_id"`
}

type GroupFlagTestConfMap map[int32]*GroupFlagTestConf

func (s GroupFlagTestConfMap) LoadFromJsonFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	jsonData := struct {
		Data []*GroupFlagTestConf `json:"data"`
	}{}
	if err := json.Unmarshal(fileBytes, &jsonData); err != nil {
		return err
	}
	for _, conf := range jsonData.Data {
		s[conf.Id] = conf
	}
	return nil
}

func (s GroupFlagTestConfMap) LoadFromXmlFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	xmlData := struct {
		Data []*GroupFlagTestConf `xml:"all_infos>item"`
	}{}
	if err := xml.Unmarshal(fileBytes, &xmlData); err != nil {
		return err
	}
	for _, conf := range xmlData.Data {
		s[conf.Id] = conf
	}
	return nil
}

func (s GroupFlagTestConfMap) GetVal(key int32) *GroupFlagTestConf {
	return s[key]
}

var instanceGroupFlagTestConfMap = make(GroupFlagTestConfMap)

func GetGroupFlagTestConfMapInst() GroupFlagTestConfMap {
	return instanceGroupFlagTestConfMap
}
