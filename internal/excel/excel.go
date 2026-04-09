package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
	"sync"
)

// ExcelUtil Excel 操作工具类
type ExcelUtil struct{}

var (
	instance *ExcelUtil
	once     sync.Once
)

// GetInstance 获取 Excel 工具单例
func GetInstance() *ExcelUtil {
	once.Do(func() {
		instance = &ExcelUtil{}
	})
	return instance
}

// ================== 新增：通用二维数组读写 ==================

// ReadToArray 读取整个工作表，返回二维字符串数组（行列结构）
// filePath: 文件路径
// sheetName: 工作表名称（空字符串则使用第一个工作表）
func (e *ExcelUtil) ReadToArray(filePath, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开 Excel 文件失败: %w", err)
	}
	defer f.Close()

	if sheetName == "" {
		sheetName = f.GetSheetName(0)
		if sheetName == "" {
			return nil, fmt.Errorf("工作簿中没有工作表")
		}
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表行失败: %w", err)
	}
	return rows, nil
}

// WriteFromArray 将二维数组写入 Excel 文件
// filePath: 输出文件路径
// sheetName: 工作表名称（默认 Sheet1）
// data: 二维数组，支持 [][]interface{} 或 [][]string
func (e *ExcelUtil) WriteFromArray(filePath, sheetName string, data interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	if sheetName == "" {
		sheetName = "Sheet1"
	}
	idx, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("创建工作表失败: %w", err)
	}
	f.SetActiveSheet(idx)

	// 反射获取二维数组的值
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data 必须是二维切片")
	}

	for i := 0; i < val.Len(); i++ {
		rowVal := val.Index(i)
		if rowVal.Kind() != reflect.Slice {
			return fmt.Errorf("第 %d 行不是切片类型", i+1)
		}
		for j := 0; j < rowVal.Len(); j++ {
			cellVal := rowVal.Index(j).Interface()
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return err
			}
			f.SetCellValue(sheetName, cellName, cellVal)
		}
	}
	return f.SaveAs(filePath)
}

// ================== 原有：结构体映射读写 ==================

// ReadToStructs 读取 Excel 文件，将工作表数据映射到结构体切片
// filePath: 文件路径
// sheetName: 工作表名称（空字符串则使用第一个工作表）
// headerRow: 表头所在行号（从1开始）
// target: 目标切片指针，例如 &[]User{}
func (e *ExcelUtil) ReadToStructs(filePath, sheetName string, headerRow int, target interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开 Excel 文件失败: %w", err)
	}
	defer f.Close()

	if sheetName == "" {
		sheetName = f.GetSheetName(0)
		if sheetName == "" {
			return fmt.Errorf("工作簿中没有工作表")
		}
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取工作表行失败: %w", err)
	}
	if len(rows) < headerRow {
		return fmt.Errorf("表头行 %d 超出总行数 %d", headerRow, len(rows))
	}

	headers := rows[headerRow-1]
	slicePtr := reflect.ValueOf(target)
	if slicePtr.Kind() != reflect.Ptr || slicePtr.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("target 必须是切片的指针")
	}
	sliceVal := slicePtr.Elem()
	elemType := sliceVal.Type().Elem()

	// 构建列索引到结构体字段的映射
	colFieldMap := make(map[int]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Name
		}
		for colIdx, header := range headers {
			if header == tag {
				colFieldMap[colIdx] = i
				break
			}
		}
	}

	// 遍历数据行
	for rowIdx := headerRow; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()
		for colIdx, cellVal := range row {
			fieldIdx, ok := colFieldMap[colIdx]
			if !ok {
				continue
			}
			field := elem.Field(fieldIdx)
			if !field.CanSet() {
				continue
			}
			if err := setFieldValue(field, cellVal); err != nil {
				return fmt.Errorf("第 %d 行第 %d 列值 '%s' 转换失败: %w", rowIdx+1, colIdx+1, cellVal, err)
			}
		}
		sliceVal.Set(reflect.Append(sliceVal, elem))
	}
	return nil
}

// WriteFromStructs 将结构体切片写入 Excel 文件
// filePath: 输出文件路径
// sheetName: 工作表名称（默认 Sheet1）
// data: 结构体切片
// includeHeader: 是否写入表头
func (e *ExcelUtil) WriteFromStructs(filePath, sheetName string, data interface{}, includeHeader bool) error {
	f := excelize.NewFile()
	defer f.Close()

	if sheetName == "" {
		sheetName = "Sheet1"
	}
	idx, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("创建工作表失败: %w", err)
	}
	f.SetActiveSheet(idx)

	sliceVal := reflect.ValueOf(data)
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("data 必须是切片")
	}
	if sliceVal.Len() == 0 {
		return f.SaveAs(filePath)
	}

	elemType := sliceVal.Type().Elem()
	var headers []string
	var fieldIndexes []int
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Name
		}
		headers = append(headers, tag)
		fieldIndexes = append(fieldIndexes, i)
	}

	// 写入表头
	if includeHeader {
		for colIdx, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}
	}

	startRow := 1
	if includeHeader {
		startRow = 2
	}
	for rowIdx := 0; rowIdx < sliceVal.Len(); rowIdx++ {
		item := sliceVal.Index(rowIdx)
		for colIdx, fieldIdx := range fieldIndexes {
			field := item.Field(fieldIdx)
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+startRow)
			f.SetCellValue(sheetName, cell, field.Interface())
		}
	}
	return f.SaveAs(filePath)
}

// setFieldValue 辅助函数：将字符串值设置到反射字段
func setFieldValue(field reflect.Value, str string) error {
	if str == "" {
		return nil
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		field.SetBool(b)
	default:
		return fmt.Errorf("不支持的类型: %v", field.Kind())
	}
	return nil
}
