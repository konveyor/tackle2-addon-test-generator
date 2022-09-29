package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	pathlib "path"

	"github.com/konveyor/tackle2-addon/command"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/pelletier/go-toml/v2"
)

//
// tkltest application analyzer.
type Tkltest struct {
	application *api.Application
	*Data
}

// Run tkltest add on
func (r *Tkltest) Run() (err error) {

	_ = r.BuildConfig()

	cmd := command.Command{Path: "tkltest-unit"}
	cmd.Dir = r.application.Bucket
	cmd.Options, err = r.options()
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		r.reportLog()
	} else {
		return err
	}

	return
}

func (r *Tkltest) BuildConfig() (err error) {

	if r.Data.General.JavaJDKHome == "" {
		r.Data.General.JavaJDKHome = JavaHome
	}
	for i, path := range r.Data.General.MonolithAppPath {
		r.Data.General.MonolithAppPath[i] = pathlib.Join(AppDir, path)
	}
	if r.Data.General.AppClasspathFile != "" {
		r.Data.General.AppClasspathFile =  pathlib.Join(AppDir, r.Data.General.AppClasspathFile)
	}
	for i, path := range r.Data.Generate.AppBuildFile {
		r.Data.Generate.AppBuildFile[i] = pathlib.Join(AppDir, path)
	}

	data, err := toml.Marshal(&r.Data)

	if err != nil {

		log.Fatal(err)
		return err
	}

	err = ioutil.WriteFile(r.output(), data, 0)

	if err != nil {
		return err
	}

	fmt.Println("config file created")
	return
}

//
// output returns output directory.
func (r *Tkltest) output() string {
	return pathlib.Join(
		r.application.Bucket,
		"config.toml")
}

func (r *Tkltest) options() (options command.Options, err error) {
	options = command.Options{"--verbose"}
	options.Add("--config-file", r.output())
	options.Add("generate")
	options.Add("ctd-amplified")
	return
}


//reportLog reports the log content.
func (r *Tkltest) reportLog() {
	path := pathlib.Join(
		HomeDir,
		".mta",
		"log",
		"mta.log")
	f, err := os.Open(path)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addon.Activity(">> %s\n", scanner.Text())
	}
	_ = f.Close()
}
