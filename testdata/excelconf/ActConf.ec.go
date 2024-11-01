// Code generated by excelconfc. DO NOT EDIT.
// source:
//   file: ExcelConfTest.xlsx
//   sheet: ActConf

package excelconf

import (
	"encoding/json"
	"os"
)

type ActType int32
type OpenCond int32
type ActTaskTargetType int32

const (
	ACT_TYPE_CHECK_IN ActType = 1
	ACT_TYPE_B        ActType = 2
	ACT_TYPE_C        ActType = 3
	ACT_TYPE_D        ActType = 4

	OPEN_COND_A OpenCond = 1
	OPEN_COND_B OpenCond = 2
	OPEN_COND_C OpenCond = 3

	ACT_TASK_TAR_TYPE_LOGIN ActTaskTargetType = 1
)

var (
	ActType_name = map[int32]string{
		1: "ACT_TYPE_CHECK_IN",
		2: "ACT_TYPE_B",
		3: "ACT_TYPE_C",
		4: "ACT_TYPE_D",
	}
	ActType_value = map[string]int32{
		"ACT_TYPE_CHECK_IN": 1,
		"ACT_TYPE_B":        2,
		"ACT_TYPE_C":        3,
		"ACT_TYPE_D":        4,
	}

	OpenCond_name = map[int32]string{
		1: "OPEN_COND_A",
		2: "OPEN_COND_B",
		3: "OPEN_COND_C",
	}
	OpenCond_value = map[string]int32{
		"OPEN_COND_A": 1,
		"OPEN_COND_B": 2,
		"OPEN_COND_C": 3,
	}

	ActTaskTargetType_name = map[int32]string{
		1: "ACT_TASK_TAR_TYPE_LOGIN",
	}
	ActTaskTargetType_value = map[string]int32{
		"ACT_TASK_TAR_TYPE_LOGIN": 1,
	}
)

func (x ActType) String() string {
	return ActType_name[int32(x)]
}

func (x OpenCond) String() string {
	return OpenCond_name[int32(x)]
}

func (x ActTaskTargetType) String() string {
	return ActTaskTargetType_name[int32(x)]
}

type ActConf struct {
	ActId     uint32  `json:"act_id,omitempty"`
	Title     string  `json:"title,omitempty"`
	ActType   ActType `json:"act_type,omitempty"`
	BeginTime string  `json:"begin_time,omitempty"`
	EndTime   string  `json:"end_time,omitempty"`
}

type ActConfMap map[uint32]*ActConf

func (s ActConfMap) LoadFromJsonFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	jsonData := struct {
		Data []ActConf `json:"data"`
	}{}
	json.Unmarshal(fileBytes, &jsonData)
	for _, conf := range jsonData.Data {
		s[conf.ActId] = &conf
	}
	return nil
}

func (s ActConfMap) GetVal(key uint32) *ActConf {
	return s[key]
}

var instanceActConfMap = make(ActConfMap)

func GetActConfMapInst() ActConfMap {
	return instanceActConfMap
}
