package excel_multi_process

import (
	"github.com/tealeg/xlsx"
	"sync"
    "fmt"
)

type ExcelReader struct {
	//Excel Offset 读取
	pageStartIndex map[int]excelOffset
	//携程缓存
	chanBuffer chan *excelStatistics
	//控制同时运行携程数量
	chanCtrl  chan int
	processId string

	totalChainSize int
	chanSize       int
	chainSizeLeft  int
	cacheSize      int
	lock           sync.Mutex
	filePath       string
    xlFile *xlsx.File
    errMap map[string]error
    errMapLock sync.Mutex
}

func newExcelReader(pageStartIndex map[int]excelOffset, chanSize int, cacheSize int, filePath string, id string) (er *ExcelReader, err error) {
	er = new(ExcelReader)
	er.pageStartIndex = pageStartIndex
	er.processId = id
	er.chanSize = chanSize
	er.cacheSize = cacheSize
	er.filePath = filePath

    if er.xlFile, err = xlsx.OpenFile(er.filePath); err != nil {
        return
    }
	return

}

func (er *ExcelReader) Process(task func(pageIndex int, rowIndex int, data []string, processId string) error) (totalCount, dealCount int) {

	er.chanBuffer = make(chan *excelStatistics, er.chanSize)
	er.chanCtrl = make(chan int, er.chanSize)
    er.errMap = make(map[string]error)

	er.lock.Lock()
	defer er.lock.Unlock()
	er.totalChainSize = 0

	totalCount = 0
	dealCount = 0

	sheetSize := len(er.xlFile.Sheets)
	for k, v := range er.pageStartIndex {
		if k < sheetSize {
			sheet := er.xlFile.Sheets[k]
			rowSize := len(sheet.Rows)
			for i := v.TopOffset; i < rowSize; i += er.cacheSize {
				var endIndex = min(rowSize, i+er.cacheSize)
			loop:
				for {
					if len(er.chanCtrl) < er.chanSize {
                        //任意数字都可以，并不会处理这个数字。
						er.chanCtrl <- 0
						er.totalChainSize += 1
						go er.innerProcess(k, i, endIndex, sheet.Rows, task)
						break loop
					} else {
                        //携程已满，直接输出
						<-er.chanCtrl
						s := <-er.chanBuffer
						er.totalChainSize -= 1
						totalCount += s.TotalCount
						dealCount += s.SuccessCount
					}
				}
			}
		}
	}

	for i := 0; i < er.totalChainSize; i++ {
		s := <-er.chanBuffer
		totalCount += s.TotalCount
		dealCount += s.SuccessCount

	}
	close(er.chanCtrl)
	close(er.chanBuffer)
	return
}

func (er *ExcelReader) innerProcess(pageIndex, startIndex, endIndex int, rows []*xlsx.Row, task func(pageIndex int, rowIndex int, data []string, processId string) error) {
	statistics := excelStatistics{SuccessCount: 0, TotalCount: 0}
	defer func() {
		er.chanBuffer <- (&statistics)
	}()

	leftOffset := er.pageStartIndex[pageIndex].LeftOffset

	for i := startIndex; i < endIndex; i++ {
		statistics.TotalCount += 1

		dataBuffer := make([]string, 0, len(rows[i].Cells))
		for j := leftOffset; j < len(rows[i].Cells); j++ {
			dataBuffer = append(dataBuffer, rows[i].Cells[j].Value)
		}

		if err := task(pageIndex, i, dataBuffer, er.processId); err != nil {
            er.errMapLock.Lock()
            er.errMap[fmt.Sprintf("%d_%d",pageIndex, i)] = err
            er.errMapLock.Unlock()
		}else{
            statistics.SuccessCount += 1
        }
	}

}

func (er *ExcelReader) GetLastErrorMap()map[string]error{
    return er.errMap
}
