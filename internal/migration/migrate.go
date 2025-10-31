package migration

import (
	"gorm.io/gorm"
)

// DBMigrate will create & migrate the tables, then make the some relationships if necessary
func DbMigrate(db *gorm.DB) {
	if err := initTables(db); err != nil {
		panic(err)
	}
	if err := initData(db); err != nil {
		panic(err)
	}
}
