package sql

import (
	"errors"
	"fmt"

	"github.com/eunnseo/AirPost/application/domain/model"
	"github.com/eunnseo/AirPost/application/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Setup() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", setting.Databasesetting.User, setting.Databasesetting.Pass, setting.Databasesetting.Server, setting.Databasesetting.Database)
	dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(errors.New("DB connection fail"))
	}

	// AutoMigrate creates the foreign-key constraints declared via the
	// `constraint:...` gorm tags on each model's associations (gorm v2 API).
	// The old gorm v1 AddForeignKey() calls no longer exist, so the FKs are
	// now expressed on the models themselves and enforced here.
	if err = dbConn.AutoMigrate(
		&model.Topic{}, &model.LogicService{},
		&model.Sink{}, &model.Node{},
		&model.SensorValue{}, &model.Logic{},
		&model.Delivery{}, &model.Path{}, &model.StationDrone{},
		&model.DroneBusy{},
	); err != nil {
		panic(errors.New("DB migration failed: " + err.Error()))
	}
}
