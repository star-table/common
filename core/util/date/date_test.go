package date

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/types"
	"gitea.bjx.cloud/allstar/common/core/util/tests"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestFormatTime(t *testing.T) {

	t1 := types.Time(time.Now())

	str := FormatTime(t1)
	fmt.Println(str)
	fmt.Println(ParseTime(str))

}

func TestFormat(t *testing.T) {

	ti := "2019-07-07 12:13:14"
	dt := Parse(ti)

	convey.Convey("Test Format ", t, tests.StartUp(func() {
		convey.Convey(" test format date str true ", func() {
			convey.So(Format(dt), convey.ShouldEqual, ti)
		})

		convey.Convey(" test format date str false ", func() {
			convey.So(Format(dt), convey.ShouldNotEqual, "2019-07-06 12:13:14")
		})

		convey.Convey(" test read config ", func() {
			convey.So(config.GetParameters().CodeConnector, convey.ShouldEqual, "#")
		})
	}))

}
