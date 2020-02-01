package mail

import (
	"errors"
	"strconv"
	"sync"

	"github.com/galaxy-book/common/core/config"
	"gopkg.in/gomail.v2"
)

var mailConn map[string]string

var mu sync.Mutex

func initConfig() {
	if mailConn == nil {
		mu.Lock()
		defer mu.Unlock()
		if mailConn == nil {
			if config.GetMailConfig() == nil {
				panic(errors.New("Mysql Datasource Configuration is missing!"))
			}

			conf := config.GetMailConfig()
			mailConn = map[string]string{
				"alias": conf.Alias,
				"user": conf.Usr,
				"pass": conf.Pwd,
				"host": conf.Host,
				"port": strconv.Itoa(conf.Port),
			}
		}
	}
}

func SendMail(mailTo []string, subject string, body string) error {

	initConfig()

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(mailConn["user"], mailConn["alias"])) //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                     //发送给多个用户
	m.SetHeader("Subject", subject)                  //设置邮件主题
	m.SetBody("text/html", body)                     //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}
