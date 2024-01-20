package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"wechat_robot/config"
	"wechat_robot/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetMySQL().Username,
		config.GetMySQL().Password,
		config.GetMySQL().IP,
		config.GetMySQL().Port,
		config.GetMySQL().DBName,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	gormDB, err = gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB}),
		&gorm.Config{
			SkipDefaultTransaction: true,
			Logger: &GormLogger{
				SlowThreshold: 2 * time.Second,
				LogLevel:      logger.Info,
			},
		})
	if err != nil {
		panic(err)
	}
}

type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level

	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logrus.GetLogger().CtxInfof(ctx, "GORM LOG %s %+v", msg, data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logrus.GetLogger().CtxWarnf(ctx, "GORM LOG %s %+v", msg, data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logrus.GetLogger().CtxErrorf(ctx, "GORM LOG %s %+v", msg, data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > logger.Silent {
		costDuration := time.Since(begin)
		cost := float64(costDuration.Nanoseconds()/1e4) / 100.0
		switch {

		// err hapends and log level is greater than 'Error'. if we shold ignore data not found err
		case err != nil && l.LogLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
			sql, _ := fc()
			logrus.GetLogger().CtxErrorf(ctx, "GORM LOG: %s, Err: %s, Cost: %.2fms", sql, err.Error(), cost)

		// slow SQL exec hapends and level is greater than 'Warn'
		case l.LogLevel >= logger.Warn && costDuration > l.SlowThreshold && l.SlowThreshold > 0:
			sql, rows := fc()
			logrus.GetLogger().CtxWarnf(ctx, "GORM LOG SLOW SQL: %s, Rows: %d, Cost: %.2fms, Limit: %s", sql, rows, cost, l.SlowThreshold)

		// normal SQL record
		case l.LogLevel >= logger.Info:
			sql, rows := fc()
			logrus.GetLogger().CtxInfof(ctx, "GORM LOG SQL: %s, Rows: %d, Cost: %.2fms", sql, rows, cost)
		}
	}
}
