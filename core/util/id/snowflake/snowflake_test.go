package snowflake

import (
	"os"
	"strconv"
	"testing"
)

func TestId(t *testing.T) {
	os.Setenv("WORK_ID", "1")
	t.Log(strconv.FormatInt(Id(), 10))
	t.Log(strconv.FormatInt(Id(), 10))
	t.Log(strconv.FormatInt(Id(), 10))
	t.Log(strconv.FormatInt(Id(), 10))
	t.Log(strconv.FormatInt(Id(), 10))

}