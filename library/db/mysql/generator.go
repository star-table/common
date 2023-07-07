package mysql

import (
	"errors"
	"github.com/polaris-team/converter"
	"github.com/star-table/common/core/config"
	"runtime"
	"strconv"
	"strings"
)

func Generate(savePath string, tables []string) error{
	if config.GetMysqlConfig() == nil {
		panic(errors.New("Mysql Datasource Configuration is missing!"))
	}
	t2t := converter.NewTable2Struct()

	var pkg string
	split := "/"
	index := strings.LastIndex(savePath, split)
	len := len(savePath)
	if index == -1{
		pkg = savePath
	}else{
		pkg = savePath[index + 1:len]
	}

	dbConf := config.GetMysqlConfig()

	t2t.Config(&converter.T2tConfig{
		TagToLower: false,
		JsonTagToHump: true,
		StructNameToHump: true,
	})

	t2t.EnableJsonTag(true).
		PackageName(pkg).
		TagKey("db").
		RealNameMethod("TableName").
		Dsn(dbConf.Usr + ":" + dbConf.Pwd + "@tcp(" + dbConf.Host + ":" + strconv.Itoa(dbConf.Port) + ")/" + dbConf.Database + "?charset=utf8")

	for i := range tables{
		err := t2t.SavePath(savePath + split + tables[i] + ".go").Table(tables[i]).Run()
		if err != nil{
			return err
		}
	}
	return nil
}

func GetFileSplit() string{
	if runtime.GOOS == "windows"{
		return "\\"
	}
	return "/"
}