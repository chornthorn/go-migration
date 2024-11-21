package templates

func getSQLiteTemplates(templateType string) (string, string) {
	switch templateType {
	case "create_table":
		return `
BEGIN TRANSACTION;

CREATE TABLE {{.TableName}} (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
`, `
BEGIN TRANSACTION;

DROP TABLE IF EXISTS {{.TableName}};

COMMIT;
`
	// Add more SQLite-specific templates
	default:
		return DefaultUpTemplate, DefaultDownTemplate
	}
}
