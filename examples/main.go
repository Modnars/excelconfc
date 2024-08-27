package main

import (
	"encoding/json"
	"os"

	excleconf "git.woa.com/modnarshen/excelconfc/output"
	"git.woa.com/modnarshen/excelconfc/util"
)

const (
	jsonDir = "../output/"
)

func main() {
	if err := excleconf.GetEnumTestMapInst().LoadFromJsonFile(jsonDir + "EnumTest.ec.json"); err != nil {
		util.LogError("load from json file failed|err:{%+v}", err)
		os.Exit(0)
	}
	util.LogInfo("map: %+v", excleconf.GetEnumTestMapInst())
	util.LogInfo("map: %+v", excleconf.GetEnumTestMapInst().GetVal(100))
	if outBytes, err := json.Marshal(excleconf.GetEnumTestMapInst().GetVal(100)); err != nil {
		util.LogError("marshal to json bytes failed -> %w", err)
	} else {
		util.LogInfo("json: %s", string(outBytes))
	}
	util.LogInfo("OK!")
}
