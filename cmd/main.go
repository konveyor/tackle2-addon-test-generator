package main

import (
	"fmt"

	hub "github.com/konveyor/tackle2-hub/addon"
)

var (
	// hub integration.
	addon = hub.Addon
	// HomeDir directory.
	HomeDir   = ""
	BinDir    = ""
	SourceDir = ""
	AppDir    = ""
	Dir       = ""
)

//
// main
func main() {
	addon.Run(func() (err error) {
		//
		// Get the addon data associated with the task.
		// https://github.com/konveyor/tackle2-addon-windup might be useful as an example
		fmt.Println("Addon executed - implement me!")
		app, _ := addon.Task.Application()
		fmt.Printf("Triggered for application: %d %s", app.ID, app.Name)
		fmt.Printf("With data: %v", addon.Task.Data())
		return
	})
}
