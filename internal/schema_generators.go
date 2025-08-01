package internal

import (
	"fmt"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/codecrafters-io/tester-utils/random"
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
      create table %v (id integer primary key, %v);
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

func (r Record) ValueFor(columnName string) string {
	return r.ValuesFor([]string{columnName})[0]
}

func (r Record) ValuesFor(columnNames []string) []string {
	values := []string{}

	for _, columnName := range columnNames {
		values = append(values, r.ColumnNamesToValuesMap[columnName])
	}

	return values
}

func generateRandomTable() Table {
	return Table{
		Name:        random.RandomWord(),
		ColumnNames: random.RandomWords(NUMBER_OF_COLUMNS),
	}
}

func generateRandomRecord(table Table) Record {
	record := Record{ColumnNamesToValuesMap: map[string]string{}}

	for _, columnName := range table.ColumnNames {
		record.ColumnNamesToValuesMap[columnName] = faker.FirstNameFemale()
	}

	return record
}
