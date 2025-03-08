package example

import "github.com/switchupcb/jet/v2/postgres"

var (
	Accounts	= newAccountsTable("public", "accounts", "")
	Users		= newUsersTable("public", "users", "")
)

type (
	accountsTable	struct {
		postgres.Table
		ID		postgres.ColumnInteger
		FirstName	postgres.ColumnString
		LastName	postgres.ColumnString
		Email		postgres.ColumnString
		CreatedAt	postgres.ColumnTimestamp
		UpdatedAt	postgres.ColumnTimestamp
		AllColumns	postgres.ColumnList
		MutableColumns	postgres.ColumnList
		DefaultColumns	postgres.ColumnList
	}
	AccountsTable	struct {
		accountsTable
		EXCLUDED	accountsTable
	}
	usersTable	struct {
		postgres.Table
		ID		postgres.ColumnInteger
		Name		postgres.ColumnString
		Password	postgres.ColumnString
		Email		postgres.ColumnString
		CreatedAt	postgres.ColumnTimestamp
		UpdatedAt	postgres.ColumnTimestamp
		AllColumns	postgres.ColumnList
		MutableColumns	postgres.ColumnList
		DefaultColumns	postgres.ColumnList
	}
	UsersTable	struct {
		usersTable
		EXCLUDED	usersTable
	}
)

func UseSchema(schema string) {
	Accounts = Accounts.FromSchema(schema)
	Users = Users.FromSchema(schema)
}
func (a UsersTable) AS(alias string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName(), alias)
}
func (a UsersTable) FromSchema(schemaName string) *UsersTable {
	return newUsersTable(schemaName, a.TableName(), a.Alias())
}
func (a UsersTable) WithPrefix(prefix string) *UsersTable {
	return newUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}
func (a UsersTable) WithSuffix(suffix string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}
func newUsersTable(schemaName, tableName, alias string) *UsersTable {
	return &UsersTable{usersTable: newUsersTableImpl(schemaName, tableName, alias), EXCLUDED: newUsersTableImpl("", "excluded", "")}
}
func newUsersTableImpl(schemaName, tableName, alias string) usersTable {
	var (
		IDColumn	= postgres.IntegerColumn("id")
		NameColumn	= postgres.StringColumn("name")
		PasswordColumn	= postgres.StringColumn("password")
		EmailColumn	= postgres.StringColumn("email")
		CreatedAtColumn	= postgres.TimestampColumn("created_at")
		UpdatedAtColumn	= postgres.TimestampColumn("updated_at")
		allColumns	= postgres.ColumnList{IDColumn, NameColumn, PasswordColumn, EmailColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns	= postgres.ColumnList{NameColumn, PasswordColumn, EmailColumn, CreatedAtColumn, UpdatedAtColumn}
		defaultColumns	= postgres.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn}
	)
	return usersTable{Table: postgres.NewTable(schemaName, tableName, alias, allColumns...), ID: IDColumn, Name: NameColumn, Password: PasswordColumn, Email: EmailColumn, CreatedAt: CreatedAtColumn, UpdatedAt: UpdatedAtColumn, AllColumns: allColumns, MutableColumns: mutableColumns, DefaultColumns: defaultColumns}
}
func (a AccountsTable) AS(alias string) *AccountsTable {
	return newAccountsTable(a.SchemaName(), a.TableName(), alias)
}
func (a AccountsTable) FromSchema(schemaName string) *AccountsTable {
	return newAccountsTable(schemaName, a.TableName(), a.Alias())
}
func (a AccountsTable) WithPrefix(prefix string) *AccountsTable {
	return newAccountsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}
func (a AccountsTable) WithSuffix(suffix string) *AccountsTable {
	return newAccountsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}
func newAccountsTable(schemaName, tableName, alias string) *AccountsTable {
	return &AccountsTable{accountsTable: newAccountsTableImpl(schemaName, tableName, alias), EXCLUDED: newAccountsTableImpl("", "excluded", "")}
}
func newAccountsTableImpl(schemaName, tableName, alias string) accountsTable {
	var (
		IDColumn	= postgres.IntegerColumn("id")
		FirstNameColumn	= postgres.StringColumn("first_name")
		LastNameColumn	= postgres.StringColumn("last_name")
		EmailColumn	= postgres.StringColumn("email")
		CreatedAtColumn	= postgres.TimestampColumn("created_at")
		UpdatedAtColumn	= postgres.TimestampColumn("updated_at")
		allColumns	= postgres.ColumnList{IDColumn, FirstNameColumn, LastNameColumn, EmailColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns	= postgres.ColumnList{FirstNameColumn, LastNameColumn, EmailColumn, CreatedAtColumn, UpdatedAtColumn}
		defaultColumns	= postgres.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn}
	)
	return accountsTable{Table: postgres.NewTable(schemaName, tableName, alias, allColumns...), ID: IDColumn, FirstName: FirstNameColumn, LastName: LastNameColumn, Email: EmailColumn, CreatedAt: CreatedAtColumn, UpdatedAt: UpdatedAtColumn, AllColumns: allColumns, MutableColumns: mutableColumns, DefaultColumns: defaultColumns}
}
