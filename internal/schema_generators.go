package internal

import (
	"fmt"
	"strings"
)

type Table struct {
	ColumnNames []string
	Name        string
}

func (t Table) CreateTableSQL() string {
	columnWithTypeList := []string{}

	for _, columnName := range t.ColumnNames {
		columnWithTypeList = append(columnWithTypeList, fmt.Sprintf("%v text", columnName))
	}

	return fmt.Sprintf(`
      create table %v (id primary key, %v);
    `, t.Name, strings.Join(columnWithTypeList, ","))
}

func (t Table) InsertRecordsSQL(records []Record) string {
	valuesSqlList := []string{}

	for _, record := range records {
		valuesSqlList = append(valuesSqlList, record.ValuesSQL(t.ColumnNames))
	}

	return fmt.Sprintf(
		`insert into %v (%v) VALUES %v`,
		t.Name,
		strings.Join(t.ColumnNames, ", "),
		strings.Join(valuesSqlList, ", "),
	)
}

type Record struct {
	ColumnNamesToValuesMap map[string]string
}

func (r Record) ValuesSQL(columnNames []string) string {
	quotedValues := []string{}

	for _, columnName := range columnNames {
		value := r.ColumnNamesToValuesMap[columnName]
		quotedValues = append(quotedValues, fmt.Sprintf("'%v'", value))
	}

	return fmt.Sprintf("(%v)", strings.Join(quotedValues, ", "))
}

func (r Record) ValuesFor(columnNames []string) []string {
	values := []string{}

	for _, columnName := range columnNames {
		values = append(values, r.ColumnNamesToValuesMap[columnName])
	}

	return values
}
