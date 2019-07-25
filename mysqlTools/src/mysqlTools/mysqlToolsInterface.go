package mysqlTools

import "database/sql"

// 									目标目录					包名			包路径	接口名
//go:generate mockgen -destination=./mysqlMock.go -package=mysqlTools mysqlTools MysqlTools

type MysqlTools interface {
	SetConfig(host, user, passwd, db string, port int, charset string, autocommit bool)		bool
	Connect()	bool
	Close() 	error
	SetLogHandle(funcName func(msg interface{}))
	SetErrLog(msg string, value ...interface{})
	SetMaxConn(num int)
	SetSqlStatement(key, value string) bool
	QueryOneStmt(key string, value ...interface{}) (sql.Row, bool)
	ChangeStmt(key string, value ...interface{}) (sql.Result, bool)
	Query(sqlCmd string, args ...interface{}) func(isClose bool) ([]sql.RawBytes, bool)
	ExecCmd(sqlCmd string, values ...interface{}) sql.Result
}
