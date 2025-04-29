package service

import "log"

type AiList struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ImgSrc      string `db:"imgSrc"`
}

// TableName 方法用于指定 GORM 使用的表名
func (AiList) TableName() string {
	return "ai_list" // 显式指定表名
}

func GetAiList() []AiList {
	var aiList []AiList
	result := dataBase.Find(&aiList)
	// 检查查询是否成功
	if result.Error != nil {
		log.Println(result.Error)
	}
	return aiList
}
