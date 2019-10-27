package excel_multi_process

import (
    "errors"

)

type ExcelReaderBuilder struct {
    cacheSize      int
    chanSize       int
    filePath       string
    pageStartIndex map[int]ExcelOffset

}

func NewExcelReaderBuilder()(erb *ExcelReaderBuilder){
    erb = new(ExcelReaderBuilder)
    erb.pageStartIndex = make(map[int]ExcelOffset)
    erb.cacheSize = 4
    erb.chanSize = 4
    return erb;
}

func (erb *ExcelReaderBuilder) AddPageRange(page, topOffset, leftOffset int)(eb* ExcelReaderBuilder){
    eb = erb
    if page>=0&&topOffset>=0&&leftOffset>=0 {
        erb.pageStartIndex[page] = ExcelOffset{TopOffset:topOffset, LeftOffset:leftOffset};

    }
    return
}

func (erb *ExcelReaderBuilder) SetCacheSize(size int)(eb* ExcelReaderBuilder){
    erb.cacheSize = size
    eb = erb
    return
}

func (erb *ExcelReaderBuilder) SetChanSize(size int)(eb* ExcelReaderBuilder){
    erb.chanSize = size
    eb = erb
    return
}

func (erb *ExcelReaderBuilder) SetFilePath(filePath string)(eb* ExcelReaderBuilder){
    erb.filePath = filePath
    eb = erb
    return
}

func (erb *ExcelReaderBuilder)build()(eb *ExcelReader, err error){
    if erb.filePath ==""{
        err = errors.New("filePath is not set")
        return
    }
    if erb.pageStartIndex ==nil {
        err = errors.New("page range is not set")
        return
    }
    if erb.chanSize <1{
        err = errors.New("chansize small then 1")
        return
    }

    if erb.cacheSize <1{
        err = errors.New("chansize small then 1")
        return
    }

    eb, err = newExcelReader(erb.pageStartIndex,erb.chanSize,erb.cacheSize, erb.filePath)
    return ;
}


