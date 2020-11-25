package main

import (
	"fmt"
	"young/tEtcd/tEtcd"
)

func main(){
	fmt.Println("main")
	//tEtcd.TPut()
	//tEtcd.TGet()
	//tEtcd.TGetDir()
	//tEtcd.TDelete()
	//tEtcd.TSetLease()
	//tEtcd.TWatch()
	//tEtcd.TOp()
	tEtcd.DistributedOptimisticLock()
}
