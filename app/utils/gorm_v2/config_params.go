package gorm_v2

// 数据库参数配置，结构体
// 用于解决复杂的业务场景连接到多台服务器部署的 mysql、sqlserver、postgresql 数据库
// 具体用法参见单元测试(test/gormv2_test.go)文件，TestCustomeParamsConnMysql 函数代码段

type ConfigParams struct {
	Write ConfigParamsDetail
	Read  ConfigParamsDetail
}
type ConfigParamsDetail struct {
	Host     string
	DataBase string
	Port     int
	Prefix   string
	User     string
	Pass     string
	Charset  string
}
