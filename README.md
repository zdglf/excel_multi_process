## Excel 多协程处理

> go get github.com/zdglf/excel_multi_process

### Example

```
    import "github.com/zdglf/excel_multi_process"
    
    erb := excel_multi_process.NewExcelReaderBuilder().
		SetCacheSize(1).
		SetChanSize(40).
		SetProcessId("id"). //id 可不设置
		AddPageRange(0, 0, 0).
		AddPageRange(1, 0, 0).
		AddPageRange(2, 0, 0).
		SetFilePath("test_data/test.xlsx")
	if er, err := erb.Build();err!=nil{
		fmt.Println(err.Error())
	}else{
		total, success := er.Process(processFunc)

		fmt.Printf("total %d success %d\n", total, success)
	}

}

func processFunc(pageIndex int, rowIndex int,data []string, processId string) error{

}

```

### 设置同时运行协程的数量

> ExcelReaderBuilder.SetChanSize

```
erb := NewExcelReaderBuilder().
        SetChanSize(40)
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

### 设置处理Id

> ExcelReaderBuilder.SetProcessId


```
//设置id可用于记录同一批导入数据
erb := NewExcelReaderBuilder().
        SetProcessId("id")

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
func processFunc(pageIndex int, rowIndex int,data []string, processId string) error{
    return nil
}

```