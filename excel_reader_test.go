package excel_multi_process

import ("testing"
    "fmt"
    "time"
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
    if er, err := erb.build();err!=nil{
        println(err.Error())
    }else{
        total, success := er.Process(processFunc)
        if(total!=1033||success!=1033){
            t.Error("number not equal 1033")
        }
        fmt.Printf("total %d success %d\n", total, success)
    }

}

func processFunc(pageIndex int, rowIndex int,data []string, id string) error{
    //fmt.Printf("    %d page %d page data %s\n", pageIndex, rowIndex, strings.Join(data, ","))
    time.Sleep(2*time.Millisecond)
    return nil
}