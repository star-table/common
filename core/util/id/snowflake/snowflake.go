package snowflake

import(
	"github.com/bwmarrin/snowflake"
	"os"
	"strconv"
)

var workId, _ = strconv.ParseInt(os.Getenv("WORK_ID"), 10, 64)
var node, err = snowflake.NewNode(workId)

func Node() *snowflake.Node{
	return node
}

func Id() int64{
	return node.Generate().Int64()
}