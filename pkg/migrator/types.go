package migrator

import "time"

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	AppliedAt time.Time
}

// MigrationData represents template data for generating migrations
type MigrationData struct {
	Name       string
	Timestamp  string
	TableName  string
	ColumnName string
	ColumnType string
	IndexName  string

	// Foreign key related
	ConstraintName  string
	ReferenceTable  string
	ReferenceColumn string
	OnDelete        string
	OnUpdate        string

	// Enum related
	EnumName   string
	EnumValues string

	// View related
	ViewName  string
	ViewQuery string

	// Check constraint
	CheckExpression string

	// Index related
	ColumnNames string
	Where       string

	// Column properties
	DefaultValue string
	NotNull      bool

	// Trigger related
	FunctionName  string
	TriggerName   string
	TriggerLogic  string
	TriggerTiming string
	TriggerEvent  string
	TriggerLevel  string

	// Partition related
	Columns              string
	PartitionType        string
	PartitionKey         string
	PartitionDefinitions string
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   string
	Name      string
	Status    string
	AppliedAt time.Time
}
