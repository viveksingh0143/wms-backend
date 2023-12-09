package migrators

import "gorm.io/gorm"

type MySQLMigrator struct{}

func (m *MySQLMigrator) Migrate(db *gorm.DB, model interface{}) error {
	// MySQL might not need specific logic if enum is supported natively
	return db.AutoMigrate(model)
}
