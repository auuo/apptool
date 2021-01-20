package main

import (
	"apptool/biz"
	"flag"
	"fmt"
)

func main() {
	cmdType := flag.String("type", "new", "new or update")
	modName := flag.String("mod", "", "mod name")
	idlPath := flag.String("idl", "", "idl file path")
	dir := flag.String("dir", "", "project path")
	flag.Parse()
	if *modName == "" {
		panic("mod name must not be empty")
	}
	if *cmdType != "new" && *cmdType != "update" {
		panic("type must be 'new' or 'update'")
	}
	if *dir == "" {
		*dir = "./" + *modName
	}
	if err := biz.Generator(*cmdType, *dir, *modName, *idlPath); err != nil {
		panic(err)
	}
	fmt.Println("success, enjoy!")
}
