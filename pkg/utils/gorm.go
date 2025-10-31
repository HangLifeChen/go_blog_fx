package utils

import (
	"gorm.io/gorm"
)

// Paginate returns a function that can be used to paginate a query.
func Paginate(curPage, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if curPage <= 0 {
			curPage = 1
		}
		if perPage <= 0 || perPage > 200 {
			perPage = 20
		}
		offset := (curPage - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}

// 根据对象解析数据库中对应的表信息
func GetDbTableName(db *gorm.DB, model interface{}) string {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(model)
	return stmt.Schema.Table
}

//使用说明
// type User struct {
//     ID   uint
//     Name string
// }

// name := GetDbTableName(db, &User{})
// fmt.Println(name) // 输出：users
