package main

import (
	"fmt"
	"github.com/cycloidio/raws/core"
)

func main() {
	accessKey := ""
	secretKey := ""
	region := []string{"eu-*"}

	c, _ := core.NewConnector(accessKey, secretKey, region, nil)
	elbs, _ := c.GetLoadBalancers(nil)
	for _, elb := range elbs {
		fmt.Println("=====================")
		fmt.Println(elb)
	}
}
