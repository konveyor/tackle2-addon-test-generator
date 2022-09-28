package main

import (
	"os"
	"path"
	"strings"

	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-addon/ssh"
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
	JavaHome = "/usr/lib/jvm/java-11-openjdk-11.0.16.1.1-1.el8_6.x86_64"
)

func init() {
	Dir, _ = os.Getwd()
	HomeDir, _ = os.UserHomeDir()
	SourceDir = path.Join(Dir, "source")
	BinDir = path.Join(Dir, "dependencies")
}

type SoftError = hub.SoftError

// addon data passed in secret
//type Data struct {
//	Repository repository.Repository
//}

type General struct {
	AppName                string   `json:"app_name" toml:"app_name"`
	MonolithAppPath        []string `json:"monolith_app_path" toml:"monolith_app_path"`
	AppClasspathFile       string   `json:"app_classpath_file,omitempty" toml:"app_classpath_file,omitempty"`
	JavaJDKHome            string   `json:"java_jdk_home,omitempty" toml:"java_jdk_home"`
	OfflineInstrumentation bool     `json:"offline_instrumentation" toml:"offline_instrumentation"`
	BuildType              string   `json:"build_type" toml:"build_type"`
}

type Generate struct {
	TimeLimit       int          `json:"time_limit" toml:"time_limit"`
	AddAssertions   bool         `json:"add_assertions" toml:"add_assertions"`
	AppBuildFile    []string     `json:"app_build_files" toml:"app_build_files"`
	TargetClassList []string     `json:"target_class_list" toml:"target_class_list"`
	CtdAmplified    CtdAmplified `json:"ctd_amplified" toml:"ctd_amplified"`
}

type CtdAmplified struct {
	BaseTestGenerator string `json:"base_test_generator" toml:"base_test_generator"`
	InteractionLevel  int    `json:"interaction_level" toml:"interaction_level"`
	NoCtdCoverage     bool   `json:"no_ctd_coverage" toml:"no_ctd_coverage"`
	NumSeqExecutions  int    `json:"num_seq_executions" toml:"num_seq_executions"`
}

type Data struct {
	Name     string   `json:"name" toml:"name"`
	General  General  `json:"general" toml:"general"`
	Generate Generate `json:"generate" toml:"generate"`
}

//
// main
func main() {

	addon.Run(func() (err error) {
		//
		// Get the addon data associated with the task.
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			err = &SoftError{Reason: err.Error()}
			return err
		}

		// tkltest
		tkltest := Tkltest{}
		tkltest.Data = d
		//
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err == nil {
			tkltest.application = application
		} else {
			return
		}

		// SSH
		agent := ssh.Agent{}
		err = agent.Start()
		if err != nil {
			return
		}

		addon.Total(2)
		if application.Repository == nil {
			err = &SoftError{Reason: "Application repository not defined."}
			return err
		}
		SourceDir = path.Join(
			Dir,
			strings.Split(
				path.Base(
					application.Repository.URL),
				".")[0])
		AppDir = path.Join(SourceDir, application.Repository.Path)
		var r repository.Repository
		r, err = repository.New(SourceDir, application)
		if err != nil {
			return
		}
		err = r.Fetch()
		if err == nil {
			addon.Increment()
		} else {
			return
		}

		//
		// Run windup.
		err = tkltest.Run()
		if err == nil {
			addon.Increment()
		} else {
			err = &SoftError{Reason: err.Error()}
			return
		}

		return
	})
}

// sample config
//config := Config{
//	Name: "TKLTEST_CONFIG_FILE",
//	General: General{
//		AppName:                "irs",
//		MonolithAppPath:        []string{"test/data/irs/monolith/target/classes"},
//		AppClasspathFile:       "test/data/irs/irsMonoClasspath.txt",
//		JavaJDKHome:            "",
//		OfflineInstrumentation: true,
//		BuildType:              "ant",
//	},
//	Generate: Generate{
//		TimeLimit:       2,
//		AddAssertions:   false,
//		AppBuildFile:    []string{"test/data/irs/monolith/for_tests_build.xml"},
//		TargetClassList: []string{},
//		CtdAmplified: CtdAmplified{
//			BaseTestGenerator: "combined",
//			InteractionLevel:  1,
//			NoCtdCoverage:     false,
//			NumSeqExecutions:  2,
//		},
//	},
//}
