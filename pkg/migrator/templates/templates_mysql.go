package templates

func getMySQLTemplates(templateType string) (string, string) {
	switch templateType {
	case "create_table":
		return `
START TRANSACTION;

CREATE TABLE {{.TableName}} (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
`, `
START TRANSACTION;

DROP TABLE IF EXISTS {{.TableName}};

COMMIT;
`
	// Add more MySQL-specific templates
	default:
		return DefaultUpTemplate, DefaultDownTemplate
	}
}
