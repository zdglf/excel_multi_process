package excel_multi_process

import (
    "github.com/tealeg/xlsx"
    "sync"
)

const CONTROL_FLAG = 0

type ExcelReader struct {
    //Excel Offset 读取
    pageStartIndex map[int]ExcelOffset
    //携程缓存
    chainBuffer chan *ExcelStatistics
    //控制同时运行携程数量
    chainCtrl chan int

    totalChainSize int
    chainSize int
    chainSizeLeft int
    cacheSize int
    lock sync.Mutex
    filePath string
    xlFile *xlsx.File

}

func newExcelReader(pageStartIndex map[int]ExcelOffset, chanSize int, cacheSize int, filePath string)(er *ExcelReader, err error){
    er = new(ExcelReader)
    er.pageStartIndex = pageStartIndex
    er.chainBuffer = make(chan *ExcelStatistics, chanSize)
    er.chainCtrl = make(chan int ,chanSize)
    er.chainSize = chanSize
    er.cacheSize = cacheSize
    er.filePath = filePath

    if er.xlFile, err = xlsx.OpenFile(er.filePath);err!=nil{
        return
    }
    return

}

func (er *ExcelReader)Process(task func(pageIndex int, rowIndex int,data []string)error)(totalCount, dealCount int){

    er.lock.Lock()
    defer er.lock.Unlock()
    er.totalChainSize = 0

    totalCount = 0;
    dealCount = 0

    sheetSize := len(er.xlFile.Sheets)
    for k,v := range er.pageStartIndex{
        if(k<sheetSize){
            sheet := er.xlFile.Sheets[k]
            rowSize := len(sheet.Rows)
            for i:=v.TopOffset;i<rowSize;i+=er.cacheSize{
                var endIndex = min(rowSize,i+er.cacheSize)
                loop:
                for{
                    if(len(er.chainCtrl)<er.chainSize){
                        er.chainCtrl<-CONTROL_FLAG
                        er.totalChainSize+=1
                        go er.innerProcess(k,i,endIndex,sheet.Rows,task)
                        break loop
                    }else{
                        <-er.chainCtrl
                        s := <- er.chainBuffer
                        er.totalChainSize-=1
                        totalCount+=s.TotalCount
                        dealCount+=s.SuccessCount
                    }
                }
            }
        }
    }

    for i:=0;i<er.totalChainSize;i++{
        s := <- er.chainBuffer
        totalCount+=s.TotalCount
        dealCount+=s.SuccessCount

    }
    return
}


func (er *ExcelReader)innerProcess(pageIndex,startIndex,endIndex int, rows []*xlsx.Row, task func(pageIndex int, rowIndex int,data []string)error){
    statistics := ExcelStatistics{SuccessCount:0,TotalCount:0}
    defer func(){
        er.chainBuffer <- (&statistics)
    }()

    leftOffset := er.pageStartIndex[pageIndex].LeftOffset

    for i:=startIndex;i<endIndex;i++{
        statistics.TotalCount+=1

        dataBuffer := make([]string, 0, len(rows[i].Cells))
        for j:=leftOffset;j<len(rows[i].Cells);j++{
            dataBuffer = append(dataBuffer, rows[i].Cells[j].Value)
        }

        if err :=task(pageIndex, i, dataBuffer);err==nil{
            statistics.SuccessCount+=1
        }
    }

}


