package main

import (
	"github.com/littlebluewhite/Account/app/dbs/sql"
	"github.com/littlebluewhite/Account/util/config"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: "./dal/query",
		/* Mode: gen.WithoutContext,*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable:  true,
		FieldCoverable: true,
	})

	db, err := sql.NewDB("mySQL", "gen_sql.log", config.SQLConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "123456",
		DB:       "account",
	})
	if err != nil {
		panic(err)
	}

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary, or it will panic
	g.UseDB(db)
	defaultAuth := g.GenerateModel("default_auth",
		gen.FieldType("auth", "json.RawMessage"))
	wUserGroup := g.GenerateModel("w_user_group")
	wGroup := g.GenerateModel("w_group")
	wUser := g.GenerateModel("w_user",
		gen.FieldType("auth", "json.RawMessage"),
		//gen.FieldRelate(field.Many2Many, "WGroups", wGroup, &field.RelateConfig{
		//	GORMTag: map[string][]string{"many2many": {"w_user_group"}},
		//}),
		gen.FieldRelate(field.HasMany, "WUserGroups", wUserGroup, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"w_user_id"}},
			RelatePointer: false,
		}),
	)
	user := g.GenerateModel("user",
		gen.FieldRelate(field.HasMany, "WUsers", wUser, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"user_id"}},
			RelatePointer: false,
		}),
	)
	wGroup = g.GenerateModel("w_group",
		//gen.FieldRelate(field.Many2Many, "WUsers", wUser, &field.RelateConfig{
		//	GORMTag: map[string][]string{"many2many": {"w_user_group"}},
		//}),
		gen.FieldType("default_auth", "json.RawMessage"),
		gen.FieldRelate(field.HasMany, "WUserGroups", wUserGroup, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"w_group_id"}},
			RelatePointer: false,
		}),
	)

	workspaceBase := g.GenerateModel("workspace")
	workspace := g.GenerateModel("workspace",
		gen.FieldType("auth", "json.RawMessage"),
		gen.FieldType("user_auth_const", "json.RawMessage"),
		gen.FieldType("user_auth_pass_down", "json.RawMessage"),
		gen.FieldType("user_auth_custom", "json.RawMessage"),
		gen.FieldRelate(field.HasMany, "WUsers", wUser, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"workspace_id"}},
			RelatePointer: false}),
		gen.FieldRelate(field.HasMany, "WGroups", wGroup, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"workspace_id"}},
			RelatePointer: false}),
		gen.FieldRelate(field.HasMany, "NextWorkspaces", workspaceBase, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"pre_workspace_id"}},
			RelatePointer: false}),
	)

	g.ApplyBasic(defaultAuth, wUserGroup, wUser, user, wGroup, workspace)

	// execute the action of code generation
	g.Execute()
}
