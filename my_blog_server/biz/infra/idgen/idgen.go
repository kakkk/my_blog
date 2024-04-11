package idgen

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func MustInit() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func GenID() int64 {
	return node.Generate().Int64()
}
