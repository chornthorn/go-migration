package templates

const (
	DefaultUpTemplate = `-- Migration: {{.Name}}
-- Created at: {{.Timestamp}}

-- Write your UP migration SQL here
`

	DefaultDownTemplate = `-- Migration: {{.Name}}
-- Created at: {{.Timestamp}}

-- Write your DOWN migration SQL here
`

	CreateTableUpTemplate = `-- Migration: {{.Name}}
-- Created at: {{.Timestamp}}

CREATE TABLE {{.TableName}} (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`

	CreateTableDownTemplate = `-- Migration: {{.Name}}
-- Created at: {{.Timestamp}}

DROP TABLE IF EXISTS {{.TableName}};
`

	AddColumnUpTemplate = `-- Migration: {{.Name}}
-- Created at: {{.Timestamp}}

ALTER TABLE {{.TableName}}
    ADD COLUMN {{.ColumnName}} {{.ColumnType}};
`

	AddColumnDownTemplate = `-- Migration: {{.Name}}
-- Created at: {{.Timestamp}}

ALTER TABLE {{.TableName}}
    DROP COLUMN IF EXISTS {{.ColumnName}};
`
)

func GetTemplatesByDriver(driver, templateType string) (string, string) {
	switch driver {
	case "postgres", "postgresql":
		return getPostgresTemplates(templateType)
	case "mysql":
		return getMySQLTemplates(templateType)
	case "sqlite", "sqlite3":
		return getSQLiteTemplates(templateType)
	default:
		return DefaultUpTemplate, DefaultDownTemplate
	}
}
