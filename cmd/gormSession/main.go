package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gormSession/configs"
	"gormSession/internal/db"
	"gormSession/internal/models"
	"gormSession/internal/query"
	"gormSession/internal/repository"
	"gormSession/pkg/logging"
	"log"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	//日志
	configs.InitConfig(configPath)
	logging.InitLogger(logging.Config{
		Debug:     false,
		InfoFile:  configs.GetApp().Log.Info.Filename,
		ErrorFile: configs.GetApp().Log.Error.Filename,
	})
	db.InitDb(configs.GetApp().Mysql)

	logging.Logger.Info("社保中心")

	appConfig := configs.GetApp()
	if appConfig.Web.RunMode == "Prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	engine.GET("/user", func(ctx *gin.Context) {
		/*dbLogger := gormLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gLogger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})*/

		traceId := uuid.New().String()

		sessionCtx := context.WithValue(context.Background(), "uuid", traceId)

		//sessionDb := db.GetDb().Session(&gorm.Session{Context: sessionCtx, Logger: dbLogger})
		//sessionDb := db.GetDb().Session(&gorm.Session{Context: sessionCtx})
		sessionDb := db.GetDb()

		repUser := repository.NewUser(sessionCtx, sessionDb)

		err := repUser.Create(nil, &models.User{
			Username: "hejinxue",
			Password: "123456",
			Status:   0,
		})
		if err != nil {
			ctx.JSON(200, gin.H{"err": err})
			return
		}
		user, err := repUser.FindByName(nil, "hejinxue")
		if err != nil {
			ctx.JSON(200, gin.H{"err": err})
			return
		}

		var user2 *models.User
		err = repository.GetQuery(sessionDb).Transaction(func(tx *query.Query) error {
			user2, err = repUser.FindSubQuery(tx, "hejinxue")

			return err
		})

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(200, gin.H{"err": err})
			return
		}

		ctx.JSON(200, gin.H{"user": user, "user2": user2})
	})

	log.Fatal(engine.Run(fmt.Sprintf("%s:%d", appConfig.Web.Ip, appConfig.Web.Port)))
}
