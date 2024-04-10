package main

import (
	"flag"
	"fmt"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gormSession/configs"
	"gormSession/internal/db"
	"gormSession/internal/models"
	"gormSession/pkg/logging"
	"os"
)

var (
	configPath string
)

// generate code
func main() {
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	if currentPath, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		fmt.Println("当前工作目录:" + currentPath)
	}
	configs.InitConfig(configPath)
	logging.InitLogger(logging.Config{
		Debug:     false,
		InfoFile:  configs.GetApp().Log.Info.Filename,
		ErrorFile: configs.GetApp().Log.Error.Filename,
	})
	db.InitDb(configs.GetApp().Mysql)

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/query/",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery, // generate mode
		ModelPkgPath: "models",
	})

	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"decimal": func(detailType gorm.ColumnType) (dataType string) { return "decimal.Decimal" },
		"datetime": func(detailType gorm.ColumnType) (dataType string) {
			if detailType.Name() == "deleted_at" {
				return "gorm.DeletedAt"
			}
			return "*time.Time"
		},
	}

	// 要先于`ApplyBasic`执行
	g.WithDataTypeMap(dataMap)

	tags := map[string]string{
		"created_at":      "sql_datetime",
		"updated_at":      "sql_datetime",
		"birthday":        "sql_date",
		"work_start_date": "sql_date",
		"work_end_date":   "sql_date",
	}

	g.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
		if v, ok := tags[columnName]; ok {
			return columnName + "\" time_format:\"" + v
		} else {
			return columnName
		}
	})

	g.UseDB(db.GetDb())

	_ = g.GenerateAllTable()

	g.ApplyInterface(
		func() {},
		models.User{},
	)
	g.Execute()
}
