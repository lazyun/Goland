# mysqlTools 使用说明

详细使用请查看 `main.go`

## 创建 mysqlTools 对象

`sqlT := new(SqlT.SqlTools)`

## 连接数据库

`ok := sqlT.SetConfig(cfg ...)`

输入：IP、用户、密码、DB、端口、字符集、是否自动提交。
输出：true or false。true 连接成功、false 连接失败，异常信息在 `sqlT.SqlErr` 中。

## 使用模版操作数据库

使用场景：使用同样的 SQL 命令操作数据库。

### 设置模版信息

使用方法：`ok := sqlT.SetSqlStatement(key, sql)`
例子：`ok := sqlT.SetSqlStatement("findId", "select id from goTest where created > ? limit 1")`

说明：key 是此 sql 语句的别名，调用使用。
输出：ok：true or false。true 设置成功、false 设置失败，异常信息在 `sqlT.SqlErr` 中。

### 调用模版

调用模块分两个接口：查询接口、修改接口。

#### 查询接口（仅限查询一条数据）

`ok, QueryOneStmtRet := sqlT.QueryOneStmt("findId", 1531800702)`

输出：ok：true or false。true 成功、false 失败，异常信息在 `sqlT.SqlErr` 中。 QueryOneStmtRet 返回结果只有当 ok 为 true 时有意义。

获取查询结果例子：

```
ok, QueryOneStmtRet := sqlT.QueryOneStmt("findId", 1531800702)
var id int
err := QueryOneStmtRet.Scan(&id)
```

#### 修改接口

```
ok = sqlT.SetSqlStatement("insertGoTest", "insert into goTest(name, created) values(?, ?)")
ok, ChangeStmtRet := sqlT.ChangeStmt("insertGoTest", "la~la~la~", 1531800702)
idChangeStmt, err := ChangeStmtRet.LastInsertId()
rowChangeStmt, err := ChangeStmtRet.RowsAffected()
```

## 查询（支持多条数据）

使用此接口要关闭，否则导致连接丢失，造成 Mysql 大量连接

```
queryRet := sqlT.Query("select id, name from goTest limit ?", 3)
for {
		ok, queryRetValue := queryRet(false)
		if !ok {
			if nil == sqlT.SqlErr {
				fmt.Printf("Get query result all!\n")
				break
			}

			fmt.Printf("Get query result fail error %s !\n", sqlT.SqlErr.Error())
			break
		}

		queryId := string( queryRetValue[0] )
		queryName := string( queryRetValue[1] )
		fmt.Printf("Query result id is %s, name is %s\n", queryId, queryName)
}
```

## 执行命令

```
execRet := sqlT.ExecCmd("insert into goTest(name, created) values(?, ?)", "la~la~la~", 1531800702)
if nil == execRet {
    fmt.Printf("ExecCmd fail error is %s\n", sqlT.SqlErr.Error() )
} else {
    execRetValue1, _ := execRet.LastInsertId()
    execRetValue2, _ := execRet.RowsAffected()
    fmt.Printf("ExecCmd insert row %d, last id is %d\n", execRetValue2, execRetValue1)
}
```