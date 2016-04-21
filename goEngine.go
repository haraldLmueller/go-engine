package main

import (
	"os"
	"runtime"

	"github.com/walesey/go-engine/editor"
	"github.com/walesey/go-engine/examples"

	"github.com/codegangsta/cli"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "goEngine"
	app.Usage = "This is a basic cli for goEngine"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{

		{
			Name:   "editor",
			Usage:  "Starts up the asset editor",
			Action: startEditor,
		},

		//DEMOS
		{
			Name:   "demo",
			Usage:  "very basic demo",
			Action: examples.Demo,
		},
		{
			Name:   "particles",
			Usage:  "run a particle effect example",
			Action: examples.Particles,
		},
		{
			Name:   "gun",
			Usage:  "run a demo of a gun model",
			Action: examples.GunDemo,
		},
		{
			Name:   "bullet",
			Usage:  "run the bullet physics demo",
			Action: examples.BulletDemo,
		},
	}

	app.Run(os.Args)
}

//CLI start the editor
func startEditor(c *cli.Context) {
	editor.New().Start()
}
