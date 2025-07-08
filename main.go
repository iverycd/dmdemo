package main

import (
	"database/sql"
	_ "dm"
	"dmdemo/config"
	"dmdemo/loggerzap"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var Host, Pwd arrayFlags

func init() {
	// 从命令行获取flag参数
	flag.Var(&Host, "host", "ip")
	flag.Var(&Pwd, "pwd", "password")
	// 通过Parse解析获取到配置文件
	flag.Parse()
	if len(Host) == 0 && len(Pwd) == 0 {
		panic(errors.New("config file not specified."))
	}

}

func killSlow(sqlStr string) {
	// 带日志
	//dsn := "dm://SYSDBA:SYSDBA@192.168.212.203:5236?appName=MacPro&connectTimeout=3000&logLevel=all&clientEncoding=GB18030"
	dsn := fmt.Sprintf("dm://SYSDBA:%s@%s:5236?appName=MacPro&connectTimeout=3000&clientEncoding=GB18030", Pwd[0], Host[0])
	db, err := sql.Open("dm", dsn)
	if err != nil {
		loggerzap.L().Fatal("", zap.Error(err))
		return
	}
	loggerzap.L().Info("连接数据库成功")
	var sqlRet, sql_text, user_name, clnt_ip, clnt_host, appname, osname, clnt_type, last_send_time string
	// 创建一个每隔1秒触发一次的Ticker
	ticker := time.NewTicker(100 * time.Millisecond)

	for range ticker.C {
		rows, err := db.Query(sqlStr)
		if err != nil {
			loggerzap.L().Fatal("", zap.Error(err))
			return
		}
		for rows.Next() {
			err := rows.Scan(&sqlRet, &sql_text, &user_name, &clnt_ip, &clnt_host, &appname, &osname, &clnt_type, &last_send_time)
			if err != nil {
				loggerzap.L().Fatal("", zap.Error(err))
				return
			}
			sqlRet = strings.ReplaceAll(sqlRet, "\n", "")
			sql_text = strings.ReplaceAll(sql_text, "\n", "")
			user_name = strings.ReplaceAll(user_name, "\n", "")
			clnt_ip = strings.ReplaceAll(clnt_ip, "\n", "")
			clnt_host = strings.ReplaceAll(clnt_host, "\n", "")
			appname = strings.ReplaceAll(appname, "\n", "")
			osname = strings.ReplaceAll(osname, "\n", "")
			clnt_type = strings.ReplaceAll(clnt_type, "\n", "")
			last_send_time = strings.ReplaceAll(last_send_time, "\n", "")

			//fmt.Println(sqlRet,
			//	sql_text,
			//	user_name,
			//	clnt_ip,
			//	clnt_host,
			//	appname,
			//	osname,
			//	clnt_type,
			//	last_send_time)
			_, err = db.Exec(sqlRet)
			if err != nil {
				loggerzap.L().Fatal("", zap.Error(err))
				return
			}
			loggerzap.L().Info("sql",
				zap.String("sqlRet", sqlRet),
				zap.String("sql_text", sql_text),
				zap.String("user_name", user_name),
				zap.String("clnt_ip", clnt_ip),
				zap.String("clnt_host", clnt_host),
				zap.String("appname", appname),
				zap.String("osname", osname),
				zap.String("clnt_type", clnt_type),
				zap.String("last_send_time", last_send_time),
			)
		}
	}

}

func basicConn() {
	// 带日志
	//dsn := "dm://SYSDBA:SYSDBA@192.168.212.2:5236?appName=MacPro&connectTimeout=3000&logLevel=all&clientEncoding=GB18030"
	// 不带日志
	dsn := "dm://SYSDBA:SYSDBA@192.168.212.2:5236?appName=MacPro&connectTimeout=3000&clientEncoding=GB18030"
	db, err := sql.Open("dm", dsn)
	if err != nil {
		return
	}
	sqlText := "select id,name from test_go"
	var id, col2 string

	err = db.QueryRow(sqlText).Scan(&id, &col2)
	if err != nil {
		return
	}
	fmt.Println(id, col2)
}
func main() {
	//全局panic捕获
	defer func() {
		if r := recover(); r != nil {
			loggerzap.L().Panic("panic", zap.Any("error", r))
		}
	}()

	// 初始化zap配置
	logCfg := &config.LogConfig{
		Level:      "info",
		FilePath:   "logs/run.log",
		MaxSize:    100, // 100MB
		MaxBackups: 20,  // 控制保留的历史文件数量上限
		//MaxAge:   5, // 日志文件的最大保留天数
		Compress: false,
	}

	// 初始化日志
	if err := loggerzap.Init("production", logCfg); err != nil {
		panic(err)
	}
	defer loggerzap.Sync()
	//sqlStr := "select 'call sp_close_session('||a.sess_id||');',a.sql_text,a.user_name,a.clnt_ip from v$sessions a   where   a.sql_text like '%DBMS_METADATA%' and a.user_name !='SYSDBA';"
	//sqlStr := "select 'call sp_close_session('||a.sess_id||');',a.sql_text,a.user_name,a.clnt_ip,a.clnt_host,a.appname,a.osname,a.clnt_type,a.last_send_time from v$sessions a   where   (a.sql_text like '%as table_cat%' or a.sql_text like '%DBMS_METADATA%' or a.sql_text like '%gateway metadata JDBC_getColumns(%') and a.user_name !='SYSDBA';"
	sqlStr := "select 'call sp_close_session('||a.sess_id||');',a.sql_text,a.user_name,a.clnt_ip,a.clnt_host,a.appname,a.osname,a.clnt_type,a.last_send_time from v$sessions a   where   (a.sql_text like '%DBMS_METADATA%' or a.sql_text like '%gateway metadata JDBC_getColumns(%') and a.user_name !='SYSDBA';"
	killSlow(sqlStr)
}
