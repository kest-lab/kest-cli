package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

type postgresForeignKeyConstraint struct {
	TableName  string
	Name       string
	Definition string
}

type postgresForeignKeyConstraintRow struct {
	ConstraintName string `gorm:"column:constraint_name"`
	TableName      string `gorm:"column:table_name"`
	Definition     string `gorm:"column:definition"`
	ColumnName     string `gorm:"column:column_name"`
	RefTableName   string `gorm:"column:ref_table_name"`
	RefColumnName  string `gorm:"column:ref_column_name"`
}

func dropForeignKeysForTargetColumns(
	db *gorm.DB,
	targetColumns map[string]struct{},
) ([]postgresForeignKeyConstraint, error) {
	if len(targetColumns) == 0 {
		return nil, nil
	}

	const query = `
SELECT
	con.conname AS constraint_name,
	child.relname AS table_name,
	pg_get_constraintdef(con.oid) AS definition,
	child_att.attname AS column_name,
	parent.relname AS ref_table_name,
	parent_att.attname AS ref_column_name
FROM pg_constraint con
JOIN pg_class child ON child.oid = con.conrelid
JOIN pg_namespace child_ns ON child_ns.oid = child.relnamespace
JOIN pg_class parent ON parent.oid = con.confrelid
JOIN pg_namespace parent_ns ON parent_ns.oid = parent.relnamespace
JOIN LATERAL generate_subscripts(con.conkey, 1) AS ord(pos) ON true
JOIN pg_attribute child_att
  ON child_att.attrelid = con.conrelid
 AND child_att.attnum = con.conkey[ord.pos]
JOIN pg_attribute parent_att
  ON parent_att.attrelid = con.confrelid
 AND parent_att.attnum = con.confkey[ord.pos]
WHERE con.contype = 'f'
  AND child_ns.nspname = current_schema()
  AND parent_ns.nspname = current_schema()
ORDER BY child.relname, con.conname, ord.pos
`

	var rows []postgresForeignKeyConstraintRow
	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	constraints := make([]postgresForeignKeyConstraint, 0)

	for _, row := range rows {
		childKey := targetColumnKey(row.TableName, row.ColumnName)
		parentKey := targetColumnKey(row.RefTableName, row.RefColumnName)
		if !targetColumnSetContains(targetColumns, childKey) && !targetColumnSetContains(targetColumns, parentKey) {
			continue
		}

		seenKey := row.TableName + "\x00" + row.ConstraintName
		if _, ok := seen[seenKey]; ok {
			continue
		}
		seen[seenKey] = struct{}{}

		constraints = append(constraints, postgresForeignKeyConstraint{
			TableName:  row.TableName,
			Name:       row.ConstraintName,
			Definition: row.Definition,
		})
	}

	for _, constraint := range constraints {
		if err := db.Exec(
			fmt.Sprintf(
				"ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s",
				quoteIdent(constraint.TableName),
				quoteIdent(constraint.Name),
			),
		).Error; err != nil {
			return nil, err
		}
	}

	return constraints, nil
}

func restoreForeignKeys(
	db *gorm.DB,
	constraints []postgresForeignKeyConstraint,
) error {
	for _, constraint := range constraints {
		if err := db.Exec(
			fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT %s %s",
				quoteIdent(constraint.TableName),
				quoteIdent(constraint.Name),
				constraint.Definition,
			),
		).Error; err != nil {
			return err
		}
	}

	return nil
}

func targetColumnKey(tableName, columnName string) string {
	return tableName + "." + columnName
}

func targetColumnSetContains(targetColumns map[string]struct{}, key string) bool {
	_, ok := targetColumns[key]
	return ok
}
