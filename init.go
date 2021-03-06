package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

var cmdInit = &Command{
	Usage: "init [package]",
	Short: "init goautoenv",
	Long: `

Initialize goautoenv. It will generate some scripts and directory will be set as "GOPATH". also it will create a symbolic link. On the windows, This program needs admin permisson for creating symbolic link. This maybe prompt privilege elevation.

After initialization, you can activate by running ".goenv/bin/activate" or ".goenv/bin/activate.ps1" or ".goenv/bin/activate.bat". You can deactivate this by running "deactivate"`,
	Run: commandInit,
}

func commandInit(cmd *Command, args []string) bool {
	root, e := getRoot()
	if e != nil {
		panic(fmt.Sprintf("Error occured while getting root of this source tree : %q", e))
	}

	var package_name string
	if len(args) < 1 {
		package_name, _ = getPackage()
		log.Printf(package_name)
		if len(package_name) == 0 {
			panic(fmt.Sprintf("Error. package named is needed."))
		}
	} else {
		package_name = args[0]
	}

	package_name_splits := strings.Split(package_name, "/")
	package_name_prefix := package_name_splits[:len(package_name_splits)-1]
	package_name_base := package_name_splits[len(package_name_splits)-1]

	goenv_root := filepath.Join(root, ".goenv")
	goenv_bin := filepath.Join(goenv_root, "bin")
	goenv_workspace := filepath.Join(goenv_root, "src", filepath.Join(package_name_prefix...))
	env := Environment{package_name, root, goenv_root}
	mkdir(goenv_bin)
	mkdir(goenv_workspace)
	MakeSymbolicLink(filepath.Join(goenv_workspace, package_name_base), root)
	writeEnvScripts(&env, goenv_bin)

	return true
}
