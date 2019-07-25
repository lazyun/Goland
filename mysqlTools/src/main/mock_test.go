package main

import (
	"mysqlTools"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSomeOne(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mysqlMock := mysqlTools.NewMockMysqlTools(ctl)

	var mysqlCfg = mysqlTools.MysqlConfig{}
	gomock.InOrder(
		mysqlMock.EXPECT().SetConfig(mysqlCfg.Host, mysqlCfg.User, mysqlCfg.Passwd, mysqlCfg.Db, mysqlCfg.Port, mysqlCfg.Charset, mysqlCfg.Autocommit).Return(true),
	)

	if mysqlMock.SetConfig(mysqlCfg.Host, mysqlCfg.User, mysqlCfg.Passwd, mysqlCfg.Db, mysqlCfg.Port, mysqlCfg.Charset, mysqlCfg.Autocommit) {
		t.Log("Connect Mysql success!")
	} else {
		t.Error("Connect Mysql fail!")
	}
}
