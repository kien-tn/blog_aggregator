package main

import (
	"github.com/kien-tn/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		panic(err)
	}
	err = c.SetUser("kien")
	if err != nil {
		panic(err)
	}
	config.ReadCfgFile()
}
