package excel_multi_process

import (
	"errors"
	"fmt"
	"testing"
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
		if len(er.GetLastErrorMap()) != 1033 {
			t.Error("error is not 1033")
		}
		fmt.Printf("total %d success %d\n", total, success)
	}
}

func TestExcelReaderBuilder_AddPageRange(t *testing.T) {
	erb := NewExcelReaderBuilder().
		SetCacheSize(1).
		SetChanSize(40).
		SetProcessId("id").
		AddPageRange(0, 1, 1).
		SetFilePath("test_data/test.xlsx")
	if er, err := erb.Build(); err != nil {
		t.Error(err)
	} else {
		total, success := er.Process(processAddPageRangeFunc)
		if total != 1032 || success != 1032 {
			t.Error("number not equal 1032 or success is not 1032")
		}
		if len(er.GetLastErrorMap()) != 0 {
			t.Error("error is not 0")
		}
		fmt.Printf("total %d success %d\n", total, success)
	}
}

func TestExcelReaderBuilder_Chinese(t *testing.T) {
	erb := NewExcelReaderBuilder().
		SetCacheSize(1).
		SetChanSize(40).
		SetProcessId("id").
		AddPageRange(0, 1, 1).
		SetFilePath("test_data/test1.xlsx")
	if er, err := erb.Build(); err != nil {
		t.Error(err)
	} else {
		total, success := er.Process(processChineseFunc)
		if total != 1 || success != 1 {
			t.Error("number not equal 1 or success is not 1")
		}
		if len(er.GetLastErrorMap()) != 0 {
			t.Error("error is not 0")
		}
		fmt.Printf("total %d success %d\n", total, success)
	}
}

func processChineseFunc(pageIndex int, rowIndex int, data []string, id string) error {
	//fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
	if data[0] != "中文测试" {
		return errors.New("中文 error")
	} else {
		return nil
	}
}

func processAddPageRangeFunc(pageIndex int, rowIndex int, data []string, id string) error {
	//fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
	if len(data) != 11 {
		return errors.New("error count")
	}
	return nil
}

func processFunc(pageIndex int, rowIndex int, data []string, id string) error {
	//fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
	return nil
}

func processErrFunc(pageIndex int, rowIndex int, data []string, id string) error {
	//fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
	return errors.New("test error")
}
