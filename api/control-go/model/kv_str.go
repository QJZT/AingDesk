package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口，将 JSONSlice 序列化为 JSON 字符串
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口，从数据库字段反序列化为 JSONSlice
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = JSONMap{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("unsupported type for JSONSlice")
	}
	return json.Unmarshal(bytes, j)
}

type KvStr struct {
	Key string  `gorm:"primarykey"` // 内容
	V   JSONMap `json:"v"`          // 值
}
