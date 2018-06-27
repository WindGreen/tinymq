package main

import (
	"flag"
	"tinymq"
)

func main(){
	var(
		num int
	)
	flag.IntVar(&num,"c",10,"num of chanel")

	mq:=tinymq.MQ{
		Size:num,
	}
	mq.Run()
}
