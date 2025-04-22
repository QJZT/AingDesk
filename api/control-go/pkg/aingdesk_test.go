package pkg

import (
	"fmt"
	"testing"
)

func TestGetModelList(t *testing.T) {
	data, err := GetModelList()
	fmt.Println("GetModelList:ERROR: ", err)
	//遍历打印 data 中的每个元素的 Title 字段
	for _, v := range data {
		fmt.Printf("%+v\n", v.ModelName)
	}
	fmt.Println("-------------------")
	// GetRagList() 测试
	ragData, err := GetRagList()
	fmt.Println("GetRagList:ERROR: ", err)
	for _, v := range ragData {
		fmt.Printf("%+v\n", v)
	}
	fmt.Println("-------------------")
	searchRagData, err := SearchRagDocument([]string{"知识库1"}, "111")
	fmt.Println("SearchRagDocument:ERROR: ", err)
	for _, v := range searchRagData {
		fmt.Printf("%+v\n", v)
	}
}
