package model

type FileData struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	FileType     string `json:"file_type"`     // 文件类型
	FileName     string `json:"file_name"`     // 文件名称
	DataType     string `json:"data_type"`     // 数据类型
	AssociatedID int    `json:"associated_id"` // 关联ID
	FilePath     string `json:"file_path"`     // 文件绝对路径
	Description  string `json:"description"`   // 文件说明
	FileSize     int64  `json:"file_size"`     // 文件大小（字节）
	Metadata     string `json:"metadata"`      // 文件其他信息（JSON格式）
}
