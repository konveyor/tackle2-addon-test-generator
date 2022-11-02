package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/konveyor/tackle2-addon/command"
	"github.com/pelletier/go-toml/v2"
)

// tkltest application analyzer.
type Tkltest struct {
	appName string
	*TackleTestConfig
}

type General struct {
	AppName       string `json:"-" toml:"app_name"`
	TestDirectory string `json:"test_directory" toml:"test_directory"`
}

type Generate struct {
	AppBuildFiles     []string `json:"app_build_files,omitempty" toml:"app_build_files,omitempty"`
	TargetClassList   []string `json:"target_class_list,omitempty" toml:"target_class_list,omitempty"`
	ExcludedClassList []string `json:"excluded_class_list,omitempty" toml:"excluded_class_list,omitempty"`
	TimeLimit         int      `json:"time_limit,omitempty" toml:"time_limit,omitempty"`
}

type TackleTestConfig struct {
	General  General  `json:"general" toml:"general"`
	Generate Generate `json:"generate" toml:"generate"`
}

// Run tkltest add on
func (r *Tkltest) Run() error {

	// get the app name from the application
	r.TackleTestConfig.General.AppName = r.appName

	configTOML, err := toml.Marshal(&r.TackleTestConfig)
	if err != nil {
		return err
	}

	tkltestConfig := path.Join(AppDir, TKLTEST_CONFIG_FILE)
	err = ioutil.WriteFile(tkltestConfig, configTOML, 0)
	if err != nil {
		return err
	}
	addon.Activity("[TklTest] created config: %s.", tkltestConfig)
	addon.Increment()
	fmt.Println(string(configTOML))

	// Mvn compile
	mvnCommand := command.Command{Path: "mvn"}
	mvnCommand.Dir = AppDir
	mvnCommand.Options.Add("compile")
	err = mvnCommand.Run()
	if err != nil {
		return err
	}
	addon.Increment()

	// Run tkltest-unit
	tklTestCommand := command.Command{Path: "tkltest-unit"}
	tklTestCommand.Dir = AppDir
	tklTestCommand.Options = command.Options{"--verbose"}
	tklTestCommand.Options.Add("generate")
	tklTestCommand.Options.Add("ctd-amplified")
	err = tklTestCommand.Run()
	if err != nil {
		r.reportLog()
		return err
	}

	addon.Activity("[TklTest] testcases generated.")
	addon.Increment()
	return nil
}

// reportLog reports the log content.
func (r *Tkltest) reportLog() {
	logPath := path.Join(
		HomeDir,
		".mta",
		"log",
		"mta.log")
	f, err := os.Open(logPath)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addon.Activity(">> %s\n", scanner.Text())
	}
	_ = f.Close()
}
