package main

import (
	"fmt"
	"go/build"
	"os"
	"runtime"

	"github.com/haraldLmueller/go-engine/controller"
	"github.com/haraldLmueller/go-engine/editor"
	"github.com/haraldLmueller/go-engine/glfwController"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
	// set working dir to access assets
	p, _ := build.Import("github.com/haraldLmueller/go-engine", "", build.FindOnly)
	os.Chdir(p.Dir)
}

func main() {
	assetDir := "."
	if len(os.Args) >= 2 {
		assetDir = os.Args[1]
	}

	fmt.Println("Using assetDir: ", assetDir)
	editor.New(assetDir).Start()
}
