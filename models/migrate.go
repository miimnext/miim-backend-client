package models

import "go_core/config"

// Migrate 执行数据库迁移
func Migrate() {
	// 执行所有模型的迁移
	err := config.DB.AutoMigrate(
		&User{},
		&Post{},         // 需要包括 Post 模型
		&PostReaction{}, // 需要包括 Post 模型
		&Tag{},          // 需要包括 Tag 模型
		&Category{},     // 需要包括 Category 模型
		&Comment{},      // 需要包括 Comment 模型
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}
