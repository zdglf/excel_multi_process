package excel_multi_process

import (
	"fmt"
	"testing"
	"time"
    "errors"
)

func TestNewExcelReaderBuilder(t *testing.T) {
	erb := NewExcelReaderBuilder().
		SetCacheSize(1).
		SetChanSize(40).
		SetProcessId("id").
		AddPageRange(0, 0, 0).
		AddPageRange(1, 0, 0).
		AddPageRange(2, 0, 0).
		SetFilePath("test_data/test.xlsx")
	if er, err := erb.Build(); err != nil {
		t.Error(err)
	} else {
		total, success := er.Process(processFunc)
		if total != 1033 || success != 1033 {
			t.Error("number not equal 1033")
		}
		fmt.Printf("total %d success %d\n", total, success)
	}

}

func TestExcelReader_GetLastErrorMap(t *testing.T) {

    erb := NewExcelReaderBuilder().
        SetCacheSize(1).
        SetChanSize(40).
        SetProcessId("id").
        AddPageRange(0, 0, 0).
        AddPageRange(1, 0, 0).
        AddPageRange(2, 0, 0).
        SetFilePath("test_data/test.xlsx")
    if er, err := erb.Build(); err != nil {
        t.Error(err)
    } else {
        total, success := er.Process(processErrFunc)
        if total != 1033 || success != 0 {
            t.Error("number not equal 1033 or success is not 0")
        }
        if(len(er.GetLastErrorMap())!=1033){
            t.Error("error is not 1033")
        }
        fmt.Printf("total %d success %d\n", total, success)
    }
}

func processFunc(pageIndex int, rowIndex int, data []string, id string) error {
	//fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
	time.Sleep(2 * time.Millisecond)
	return nil
}


func processErrFunc(pageIndex int, rowIndex int, data []string, id string) error {
    //fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
    return errors.New("test error")
}