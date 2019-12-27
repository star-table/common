package temp

import (
	"bytes"
	"text/template"
)

//模板解析工具
func Render(str string, params interface{}) (string, error) {
	tmpl, err := template.New("test").Parse(str) //建立一个模板
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")

	err = tmpl.Execute(buf, params) //将struct与模板合成，合成结果放到os.Stdout里
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderIgnoreError(str string, params interface{}) string {
	afterRender, err := Render(str, params)
	if err != nil{
		return str
	}
	return afterRender
}
