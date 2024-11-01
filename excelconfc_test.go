package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	excelconf "git.woa.com/modnarshen/excelconfc/testdata/excelconf"
)

func TestExcelconfcCmd(t *testing.T) {
	cmd := exec.Command("go", "build")
	output, err := cmd.Output()
	require.Equal(t, nil, err, "build excelconfc should be succeed")
	require.Equal(t, 0, len(output), "len(output) should be 0")

	output, err = exec.Command("rm", "-rf", "./output").Output()
	require.Equal(t, nil, err)
	require.Equal(t, 0, len(output))

	output, err = exec.Command("mkdir", "-p", "./output").Output()
	require.Equal(t, nil, err)
	require.Equal(t, 0, len(output))

	testSheetNames := []string{`ActConf`, `ActTaskConf`, `ArrayAndBDTVecConf`, `GroupFlagTestConf`, `NestFieldsTestConf`}
	appendArguments := [][]string{{`-add_enum`}, {}, {}, {}, {}}
	for i, testSheetName := range testSheetNames {
		baseArgs := []string{
			`-excel=./testdata/ExcelConfTest.xlsx`,
			`-sheet=` + testSheetName,
			`-go_package=git.woa.com/modnarshen/uasvr-go/configs/excelconf;excelconf`,
			`-outdir=./output`,
			`--group=server`,
		}

		args := append(baseArgs, appendArguments[i]...)
		cmd1 := exec.Command(`./excelconfc`, args...)
		output, err = cmd1.Output()
		require.Equal(t, nil, err, "execute result should be success, real: %v", err)
		require.Equal(t, 0, len(output), "len(output) should be 0")

		cmd2 := exec.Command(`diff`, `./output/`+testSheetName+`.ec.proto`, `./testdata/excelconf/`+testSheetName+`.ec.proto`)
		output, err = cmd2.Output()
		require.Equal(t, nil, err, "execute result should be success, real: %v", err)
		require.Equal(t, 0, len(output), "len(output) should be 0")

		cmd3 := exec.Command(`diff`, `./output/`+testSheetName+`.ec.go`, `./testdata/excelconf/`+testSheetName+`.ec.go`)
		output, err = cmd3.Output()
		require.Equal(t, nil, err, "execute result should be success, real: %v", err)
		require.Equal(t, 0, len(output), "len(output) should be 0")

		outFileBytes, err := os.ReadFile(`./output/` + testSheetName + `.ec.json`)
		require.Equal(t, nil, err)
		obj1 := map[string]any{}
		err = json.Unmarshal(outFileBytes, &obj1)
		require.Equal(t, nil, err)

		stdFileBytes, err := os.ReadFile(`./testdata/excelconf/` + testSheetName + `.ec.json`)
		require.Equal(t, nil, err)
		obj2 := map[string]any{}
		err = json.Unmarshal(stdFileBytes, &obj2)
		require.Equal(t, nil, err)

		require.True(t, reflect.DeepEqual(obj1, obj2))

		t.Logf("sheet: %s passed", testSheetName)
	}
}

func TestLoadFromJson(t *testing.T) {
	jsonFileDir := `./testdata/excelconf/`
	var err error

	err = excelconf.GetActConfMapInst().LoadFromJsonFile(jsonFileDir + `ActConf.ec.json`)
	require.Equal(t, nil, err)
	require.Greater(t, len(excelconf.GetActConfMapInst()), 0)

	err = excelconf.GetActTaskConfMapInst().LoadFromJsonFile(jsonFileDir + `ActTaskConf.ec.json`)
	require.Equal(t, nil, err)
	require.Greater(t, len(excelconf.GetActTaskConfMapInst()), 0)

	err = excelconf.GetArrayAndBDTVecConfMapInst().LoadFromJsonFile(jsonFileDir + `ArrayAndBDTVecConf.ec.json`)
	require.Equal(t, nil, err)
	require.Greater(t, len(excelconf.GetArrayAndBDTVecConfMapInst()), 0)

	err = excelconf.GetGroupFlagTestConfMapInst().LoadFromJsonFile(jsonFileDir + `GroupFlagTestConf.ec.json`)
	require.Equal(t, nil, err)
	require.Greater(t, len(excelconf.GetGroupFlagTestConfMapInst()), 0)

	err = excelconf.GetNestFieldsTestConfMapInst().LoadFromJsonFile(jsonFileDir + `NestFieldsTestConf.ec.json`)
	require.Equal(t, nil, err)
	require.Greater(t, len(excelconf.GetNestFieldsTestConfMapInst()), 0)
	require.Equal(t, "20002", excelconf.GetNestFieldsTestConfMapInst().GetVal(2).A[0].AA[0].Aa2)
}
