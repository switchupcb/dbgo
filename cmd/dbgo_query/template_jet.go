package query

import (
	"github.com/switchupcb/jet/v2/generator/metadata"
	"github.com/switchupcb/jet/v2/generator/template"
	"github.com/switchupcb/jet/v2/postgres"
)

var (
	modelPkg = "model"
	tablePkg = "table"
	viewPkg  = "view"
	enumPkg  = "enum"
)

// Copyright 2019 Goran Bjelanovic
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

func genTemplate() template.Template {
	return template.Default(postgres.Dialect).
		UseSchema(func(schemaMetaData metadata.Schema) template.Schema {
			return template.DefaultSchema(schemaMetaData).
				UseModel(template.DefaultModel().UsePath(modelPkg).
					UseTable(func(table metadata.Table) template.TableModel {

						return template.DefaultTableModel(table)
					}).
					UseView(func(view metadata.Table) template.ViewModel {
						return template.DefaultViewModel(view)
					}).
					UseEnum(func(enum metadata.Enum) template.EnumModel {
						return template.DefaultEnumModel(enum)
					}),
				).
				UseSQLBuilder(template.DefaultSQLBuilder().
					UseTable(func(table metadata.Table) template.TableSQLBuilder {
						return template.DefaultTableSQLBuilder(table).UsePath(tablePkg)
					}).
					UseView(func(table metadata.Table) template.ViewSQLBuilder {
						return template.DefaultViewSQLBuilder(table).UsePath(viewPkg)
					}).
					UseEnum(func(enum metadata.Enum) template.EnumSQLBuilder {
						return template.DefaultEnumSQLBuilder(enum).UsePath(enumPkg)
					}),
				)
		})
}
