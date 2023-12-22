package main

import (
	"account/app/dbs/sql"
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

	db, err := sql.NewDB("mySQL", "gen_sql.log", "db")
	if err != nil {
		panic(err)
	}

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary, or it will panic
	g.UseDB(db)
	defaultAuth := g.GenerateModel("default_auth",
		gen.FieldType("auth", "json.RawMessage"))
	userWorkspace := g.GenerateModel("user_workspace",
		gen.FieldType("auth", "json.RawMessage"))
	userGroup := g.GenerateModel("user_group")
	wGroup := g.GenerateModel("w_group",
		gen.FieldRelate(field.HasMany, "users", userGroup, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"group_id"}},
			RelatePointer: false}),
	)
	wUser := g.GenerateModel("w_user",
		gen.FieldRelate(field.HasMany, "groups", userGroup, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"user_id"}},
			RelatePointer: false}),
		gen.FieldRelate(field.HasMany, "workspaces", userWorkspace, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"user_id"}},
			RelatePointer: false,
		}),
	)
	wGroup = g.GenerateModel("w_group",
		gen.FieldRelate(field.Many2Many, "users", wUser, &field.RelateConfig{
			GORMTag: map[string][]string{"many2many": {"user_group"}},
		}))
	workspace := g.GenerateModel("workspace",
		gen.FieldRelate(field.HasMany, "users", userWorkspace, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"workspace_id"}},
			RelatePointer: false}),
		gen.FieldRelate(field.HasMany, "groups", wGroup, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"workspace_id"}},
			RelatePointer: false}),
	)

	g.ApplyBasic(defaultAuth, userGroup, userWorkspace, wUser, wGroup, workspace)

	// execute the action of code generation
	g.Execute()
}
