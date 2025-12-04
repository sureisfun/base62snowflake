package main

import (
	"fmt"
	"github.com/sureisfun/base62snowflake/snowflake"
)

func main() {
	// Simple usage - just get IDs
	fmt.Println("Simple API:")
	for i := 0; i < 5; i++ {
		fmt.Println(base62snowflake.GetSnowflakeID())
	}

	// Advanced usage - custom node
	fmt.Println("\nCustom node:")
	node, _ := base62snowflake.NewNode(42)
	for i := 0; i < 5; i++ {
		fmt.Println(node.Generate())
	}

	// Set custom default node (before first use)
	// base62snowflake.SetDefaultNodeID(123)


	fmt.Println("----")
	fmt.Println(base62snowflake.GetSnowflakeID())
}
