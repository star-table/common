package uuid

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUuid(t *testing.T) {

	fmt.Println(NewUuid())
}

func TestNewUuid(t *testing.T) {
	Convey("TestNewUuid", t, func() {
		Convey(" test new uuid ", func() {
			So(NewUuid(), ShouldNotBeEmpty)
		})
	})
}
