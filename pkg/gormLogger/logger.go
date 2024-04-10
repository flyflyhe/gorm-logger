package gormLogger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogLevel log level
type LogLevel int

var (
	// Discard logger will print any log to io.Discard
	Discard = New(log.New(io.Discard, "", log.LstdFlags), gLogger.Config{})
	// Default Default logger
	Default = New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gLogger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	// Recorder logger records running SQL into a recorder instance
	Recorder = traceRecorder{Interface: Default, BeginAt: time.Now()}
)

// New initialize logger
func New(writer gLogger.Writer, config gLogger.Config) gLogger.Interface {
	var (
		infoStr      = "[%s]%s\n[info] "
		warnStr      = "[%s]%s\n[warn] "
		errStr       = "[%s]%s\n[error] "
		traceStr     = "[%s]%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "[%s]%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "[%s]%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = "[%s]" + Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = "[%s]" + BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = "[%s]" + Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = "[%s]" + Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = "[%s]" + Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = "[%s]" + RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &logger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	gLogger.Writer
	gLogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	uuid := ctx.Value("uuid")
	uuidStr, ok := uuid.(string)
	if !ok {
		uuidStr = "-----"
	}

	if l.LogLevel >= gLogger.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{uuidStr, utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	uuid := ctx.Value("uuid")
	uuidStr, ok := uuid.(string)
	if !ok {
		uuidStr = "----"
	}
	if l.LogLevel >= gLogger.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{uuidStr, utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	uuid := ctx.Value("uuid")
	uuidStr, ok := uuid.(string)
	if !ok {
		uuidStr = "----"
	}
	if l.LogLevel >= gLogger.Error {
		l.Printf(l.errStr+msg, append([]interface{}{uuidStr, utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	uuid := ctx.Value("uuid")
	uuidStr, ok := uuid.(string)
	if !ok {
		uuidStr = "----"
	}

	if l.LogLevel <= gLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gLogger.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, uuidStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, uuidStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, uuidStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, uuidStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, uuidStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, uuidStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// ParamsFilter filter params
func (l *logger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

type traceRecorder struct {
	gLogger.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

// New trace recorder
func (l *traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

// Trace implement logger interface
func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
