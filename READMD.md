## Excel 多协程处理

### 设置同时运行协程的数量

> ExcelReaderBuilder.SetChanSize

```
erb := NewExcelReaderBuilder().
        SetChanSize(40).
```

### 设置一个协程同时处理的Excel行数 

> ExcelReaderBuilder.SetCacheSize


```
erb := NewExcelReaderBuilder().
        SetCacheSize(10)
```

### 设置文件路径

> ExcelReaderBuilder.SetFilePath

```
erb := NewExcelReaderBuilder().
        SetFilePath("test_data/test.xlsx")

```

### 添加页码和偏移

> ExcelReaderBuilder.AddPageRange

```
//pageIndex 设置需要读取的页码
//设置从上往下开始的索引。
//设置从左往右开始的索引。
erb := NewExcelReaderBuilder().
        AddPageRange(0, 0, 0).
```

### 开始处理数据

> ExcelReader.Process

```
if excelReader, err := erb.build();err!=nil{
    println(err.Error())
}else{
    total, success := excelReader.Process(processFunc)
    
}
```

### processFunc 函数

> 用户自定义传入,需注意多协程并发问题

```
//pageIndex 处理的页面索引 从0开始
//rowIndex 处理的行数索引 一般从0开始。可以设置从头开始忽略多少行
//data 处理的数据。如列数offset从1开始，data[0] 对应列数 excel column[1]。
//不定长
func processFunc(pageIndex int, rowIndex int,data []string) error{
    return nil
}

```