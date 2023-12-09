package migrators

import (
	"fmt"
	"gorm.io/gorm"
)

type MSSQLMigrator struct{}

func (m *MSSQLMigrator) Migrate(db *gorm.DB, model interface{}) error {
	if err := db.AutoMigrate(model); err != nil {
		return fmt.Errorf("gorm auto migrate: %w", err)
	}
	return nil
	//if err := adjustModelForMSSQL(db, model); err != nil {
	//	return fmt.Errorf("adjusting model for MSSQL: %w", err)
	//}
	//
	//if err := db.AutoMigrate(model); err != nil {
	//	return fmt.Errorf("gorm auto migrate: %w", err)
	//}
	//
	//return m.addCheckConstraints(db, model)
}

func (m *MSSQLMigrator) addCheckConstraints(db *gorm.DB, model interface{}) error {
	enumFields, err := getEnumFieldsInfo(db, model)
	if err != nil {
		return err
	}

	tableName := getTableName(model)
	if err != nil {
		return err
	}

	for columnName, enumValues := range enumFields {
		constraintSQL := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT ck_%s CHECK (%s IN (%s))", tableName, columnName, columnName, enumValues)
		if err := db.Exec(constraintSQL).Error; err != nil {
			return err
		}
	}

	return nil
}
