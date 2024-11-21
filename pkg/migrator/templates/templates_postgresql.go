package templates

func getPostgresTemplates(templateType string) (string, string) {
	switch templateType {
	case "create_table":
		return `
CREATE TABLE {{.TableName}} (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`, `
DROP TABLE IF EXISTS {{.TableName}};
`
	// Add more PostgreSQL-specific templates
	default:
		return DefaultUpTemplate, DefaultDownTemplate
	}
}
