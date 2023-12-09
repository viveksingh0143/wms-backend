package migrators

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"sync"
)

type DatabaseMigrator interface {
	Migrate(db *gorm.DB, model interface{}) error
}

func NewDatabaseMigrator(dbDialect string) DatabaseMigrator {
	switch dbDialect {
	case "mysql":
		return &MySQLMigrator{}
	case "sqlserver":
		return &MSSQLMigrator{}
	default:
		return nil
	}
}

func parseModelSchema(db *gorm.DB, model interface{}) (*schema.Schema, error) {
	parsedSchema, err := schema.Parse(model, &sync.Map{}, db.NamingStrategy)
	if err != nil {
		return nil, err
	}
	return parsedSchema, nil
}

func adjustModelForMSSQL(db *gorm.DB, model interface{}) error {
	parsedSchema, err := parseModelSchema(db, model)
	if err != nil {
		return err
	}

	for _, field := range parsedSchema.Fields {
		if isEnumField(field) {
			// Adjust the field type to varchar for MSSQL
			field.DataType = schema.String
			field.Size = 255
			// Update tag settings as necessary
			field.TagSettings["SIZE"] = "255"
			field.TagSettings["TYPE"] = "varchar(255)"
		}
	}

	return nil
}

func getEnumFieldsInfo(db *gorm.DB, model interface{}) (map[string]string, error) {
	parsedSchema, err := parseModelSchema(db, model)
	if err != nil {
		return nil, err
	}

	enumFields := make(map[string]string)
	for _, field := range parsedSchema.Fields {
		if isEnumField(field) {
			enumValues := extractEnumValues(string(field.TagSettings["GORM"]))
			enumFields[field.DBName] = enumValues
		}
	}
	return enumFields, nil
}

func isEnumField(field *schema.Field) bool {
	fieldDataType := string(field.DataType)
	return strings.Contains(fieldDataType, "enum(")
}

func extractEnumValues(gormTag string) string {
	// Find the "enum(" substring in the GORM tag
	enumStartIndex := strings.Index(gormTag, "enum(")
	if enumStartIndex == -1 {
		return "" // Not found, return empty string
	}

	// Extract the enum values substring
	enumValuesSubstring := gormTag[enumStartIndex+len("enum("):]

	// Find the closing parenthesis
	closingParenthesisIndex := strings.Index(enumValuesSubstring, ")")
	if closingParenthesisIndex == -1 {
		return "" // Missing closing parenthesis, return empty string
	}

	// Extract the enum values (remove leading and trailing spaces)
	enumValues := strings.TrimSpace(enumValuesSubstring[:closingParenthesisIndex])

	return enumValues
}

func getTableName(model interface{}) string {
	// Get the schema for the model
	modelSchema, _ := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})

	// Get the table name from the model schema
	tableName := modelSchema.Table
	return tableName
}
