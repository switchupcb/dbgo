package jet

import (
	"database/sql"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/switchupcb/dbgo/cmd/dbgo_query/jet_internal/3rdparty/snaker"
)

// ToGoIdentifier converts database to Go identifier.
func ToGoIdentifier(databaseIdentifier string) string {
	return snaker.SnakeToCamel(replaceInvalidChars(databaseIdentifier))
}

// ToGoEnumValueIdentifier converts enum value name to Go identifier name.
func ToGoEnumValueIdentifier(enumName, enumValue string) string {
	enumValueIdentifier := ToGoIdentifier(enumValue)
	if !unicode.IsLetter([]rune(enumValueIdentifier)[0]) {
		return ToGoIdentifier(enumName) + enumValueIdentifier
	}

	return enumValueIdentifier
}

// ToGoFileName converts database identifier to Go file name.
func ToGoFileName(databaseIdentifier string) string {
	return strings.ToLower(replaceInvalidChars(databaseIdentifier))
}

// SaveGoFile saves go file at folder dir, with name fileName and contents text.
func SaveGoFile(dirPath, fileName string, text []byte) error {
	newGoFilePath := filepath.Join(dirPath, fileName) + ".go"

	file, err := os.Create(newGoFilePath)

	if err != nil {
		return err
	}

	defer file.Close()

	p, err := format.Source(text)
	if err != nil {
		return err
	}

	_, err = file.Write(p)

	if err != nil {
		return err
	}

	return nil
}

// EnsureDirPath ensures dir path exists. If path does not exist, creates new path.
func EnsureDirPath(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)

		if err != nil {
			return err
		}
	}

	return nil
}

// CleanUpGeneratedFiles deletes everything at folder dir.
func CleanUpGeneratedFiles(dir string) error {
	exist, err := DirExists(dir)

	if err != nil {
		return err
	}

	if exist {
		err := os.RemoveAll(dir)

		if err != nil {
			return err
		}
	}

	return nil
}

// DBClose closes non nil db connection
func DBClose(db *sql.DB) {
	if db == nil {
		return
	}

	db.Close()
}

// DirExists checks if folder at path exist.
func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func replaceInvalidChars(str string) string {
	str = strings.Replace(str, " ", "_", -1)
	str = strings.Replace(str, "-", "_", -1)
	str = strings.Replace(str, ".", "_", -1)

	return str
}

// IsNil check if v is nil
func IsNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

// MustBeTrue panics when condition is false
func MustBeTrue(condition bool, errorStr string) {
	if !condition {
		panic(errorStr)
	}
}

// MustBe panics with errorStr error, if v interface is not of reflect kind
func MustBe(v interface{}, kind reflect.Kind, errorStr string) {
	if reflect.TypeOf(v).Kind() != kind {
		panic(errorStr)
	}
}

// ValueMustBe panics with errorStr error, if v value is not of reflect kind
func ValueMustBe(v reflect.Value, kind reflect.Kind, errorStr string) {
	if v.Kind() != kind {
		panic(errorStr)
	}
}

// TypeMustBe panics with errorStr error, if v type is not of reflect kind
func TypeMustBe(v reflect.Type, kind reflect.Kind, errorStr string) {
	if v.Kind() != kind {
		panic(errorStr)
	}
}

// MustBeInitializedPtr panics with errorStr if val interface is nil
func MustBeInitializedPtr(val interface{}, errorStr string) {
	if IsNil(val) {
		panic(errorStr)
	}
}

// PanicOnError panics if err is not nil
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// ErrorCatch is used in defer to recover from panics and to set err
func ErrorCatch(err *error) {
	recovered := recover()

	if recovered == nil {
		return
	}

	recoveredErr, isError := recovered.(error)

	if isError {
		*err = recoveredErr
	} else {
		*err = fmt.Errorf("%v", recovered)
	}
}

// StringSliceContains checks if slice of strings contains a string
func StringSliceContains(strings []string, contains string) bool {
	for _, str := range strings {
		if str == contains {
			return true
		}
	}

	return false
}

// ExtractDateTimeComponents extracts number of days, hours, minutes, seconds, microseconds from duration
func ExtractDateTimeComponents(duration time.Duration) (days, hours, minutes, seconds, microseconds int64) {
	days = int64(duration / (24 * time.Hour))
	reminder := duration % (24 * time.Hour)

	hours = int64(reminder / time.Hour)
	reminder = reminder % time.Hour

	minutes = int64(reminder / time.Minute)
	reminder = reminder % time.Minute

	seconds = int64(reminder / time.Second)
	reminder = reminder % time.Second

	microseconds = int64(reminder / time.Microsecond)

	return
}

// SerializeClauseList func
func SerializeClauseList(statement StatementType, clauses []Serializer, out *SQLBuilder) {

	for i, c := range clauses {
		if i > 0 {
			out.WriteString(", ")
		}

		if c == nil {
			panic("jet: nil clause")
		}

		c.serialize(statement, out)
	}
}

func serializeExpressionList(statement StatementType, expressions []Expression, separator string, out *SQLBuilder) {

	for i, value := range expressions {
		if i > 0 {
			out.WriteString(separator)
		}

		value.serialize(statement, out)
	}
}

// SerializeProjectionList func
func SerializeProjectionList(statement StatementType, projections []Projection, out *SQLBuilder) {
	for i, col := range projections {
		if i > 0 {
			out.WriteString(",")
			out.NewLine()
		}

		if col == nil {
			panic("jet: Projection is nil")
		}

		col.serializeForProjection(statement, out)
	}
}

// SerializeColumnNames func
func SerializeColumnNames(columns []Column, out *SQLBuilder) {
	for i, col := range columns {
		if i > 0 {
			out.WriteString(", ")
		}

		if col == nil {
			panic("jet: nil column in columns list")
		}

		out.WriteIdentifier(col.Name())
	}
}

// SerializeColumnExpressionNames func
func SerializeColumnExpressionNames(columns []ColumnExpression, statementType StatementType,
	out *SQLBuilder, options ...SerializeOption) {
	for i, col := range columns {
		if i > 0 {
			out.WriteString(", ")
		}

		if col == nil {
			panic("jet: nil column in columns list")
		}

		col.serialize(statementType, out, options...)
	}
}

// ExpressionListToSerializerList converts list of expressions to list of serializers
func ExpressionListToSerializerList(expressions []Expression) []Serializer {
	var ret []Serializer

	for _, expr := range expressions {
		ret = append(ret, expr)
	}

	return ret
}

// ColumnListToProjectionList func
func ColumnListToProjectionList(columns []ColumnExpression) []Projection {
	var ret []Projection

	for _, column := range columns {
		ret = append(ret, column)
	}

	return ret
}

// ToSerializerValue creates Serializer type from the value
func ToSerializerValue(value interface{}) Serializer {
	if clause, ok := value.(Serializer); ok {
		return clause
	}

	return literal(value)
}

// UnwindRowFromModel func
func UnwindRowFromModel(columns []Column, data interface{}) []Serializer {
	structValue := reflect.Indirect(reflect.ValueOf(data))

	row := []Serializer{}

	ValueMustBe(structValue, reflect.Struct, "jet: data has to be a struct")

	for _, column := range columns {
		columnName := column.Name()
		structFieldName := ToGoIdentifier(columnName)

		structField := structValue.FieldByName(structFieldName)

		if !structField.IsValid() {
			panic("missing struct field for column : " + columnName)
		}

		var field interface{}

		if structField.Kind() == reflect.Ptr && structField.IsNil() {
			field = nil
		} else {
			field = reflect.Indirect(structField).Interface()
		}

		row = append(row, literal(field))
	}

	return row
}

// UnwindRowsFromModels func
func UnwindRowsFromModels(columns []Column, data interface{}) [][]Serializer {
	sliceValue := reflect.Indirect(reflect.ValueOf(data))
	ValueMustBe(sliceValue, reflect.Slice, "jet: data has to be a slice.")

	rows := [][]Serializer{}

	for i := 0; i < sliceValue.Len(); i++ {
		structValue := sliceValue.Index(i)

		rows = append(rows, UnwindRowFromModel(columns, structValue.Interface()))
	}

	return rows
}

// UnwindRowFromValues func
func UnwindRowFromValues(value interface{}, values []interface{}) []Serializer {
	row := []Serializer{}

	allValues := append([]interface{}{value}, values...)

	for _, val := range allValues {
		row = append(row, ToSerializerValue(val))
	}

	return row
}

// UnwindColumns func
func UnwindColumns(column1 Column, columns ...Column) []Column {
	columnList := []Column{}

	if val, ok := column1.(ColumnList); ok {
		for _, col := range val {
			columnList = append(columnList, col)
		}
		columnList = append(columnList, columns...)
	} else {
		columnList = append(columnList, column1)
		columnList = append(columnList, columns...)
	}

	return columnList
}

// UnwidColumnList func
func UnwidColumnList(columns []Column) []Column {
	ret := []Column{}

	for _, col := range columns {
		if columnList, ok := col.(ColumnList); ok {
			for _, c := range columnList {
				ret = append(ret, c)
			}
		} else {
			ret = append(ret, col)
		}
	}

	return ret
}

// OptionalOrDefaultString will return first value from variable argument list str or
// defaultStr if variable argument list is empty
func OptionalOrDefaultString(defaultStr string, str ...string) string {
	if len(str) > 0 {
		return str[0]
	}

	return defaultStr
}

// OptionalOrDefaultExpression will return first value from variable argument list expression or
// defaultExpression if variable argument list is empty
func OptionalOrDefaultExpression(defaultExpression Expression, expression ...Expression) Expression {
	if len(expression) > 0 {
		return expression[0]
	}

	return defaultExpression
}
