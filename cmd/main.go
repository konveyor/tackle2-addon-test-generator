package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/FabianWe/goslugify"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
)

const (
	TKLTEST_CONFIG_FILE = "tkltest_config.toml"
	TKLTEST_LOG_FILE    = "tkltest_unit.log"
	TKLTEST_TESTDIR_FMT = "tkltest-output-unit-%s/"
	COMMIT_MSG          = "tackl-test-generator: commit generated tests"
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

type SoftError = hub.SoftError

// addon data passed in secret
type Data struct {
	Branch string            `json:"branch_name"`
	Config *TackleTestConfig `json:"tkltest_config"`
}

// main
func main() {

	addon.Run(func() error {

		Dir, _ = os.Getwd()
		HomeDir, _ = os.UserHomeDir()
		SourceDir = path.Join(Dir, "source")
		BinDir = path.Join(Dir, "dependencies")

		// Get the addon data associated with the task.
		d := &Data{}
		if err := addon.DataWith(d); err != nil {
			return &SoftError{Reason: err.Error()}
		}

		if d.Branch == "" {
			return fmt.Errorf("Must specify branch_name, where changes will be written.")
		}
		if d.Config == nil {
			return fmt.Errorf("Must specify tkltest_config.")
		}

		// Setup tkltest
		tkltest := Tkltest{}
		tkltest.TackleTestConfig = d.Config

		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err != nil {
			return err
		}
		// NOTE: We slugify to handle application names with spaces intelligently
		tkltest.appName = goslugify.GenerateSlug(application.Name)

		// SSH
		agent := ssh.Agent{}
		if err = agent.Start(); err != nil {
			return err
		}

		addon.Total(5)
		if application.Repository == nil {
			return &SoftError{Reason: "Application repository not defined."}
		}
		SourceDir = path.Join(
			Dir,
			strings.Split(
				path.Base(
					application.Repository.URL),
				".")[0])
		AppDir = path.Join(SourceDir, application.Repository.Path)
		repo, err := repository.New(SourceDir, application)
		if err != nil {
			return err
		}
		err = repo.Fetch()
		if err != nil {
			return err
		}
		addon.Increment()

		// Run tkltest.
		if err = tkltest.Run(); err != nil {
			return &SoftError{Reason: err.Error()}
		}
		addon.Increment()

		// Commit the results
		addon.Activity("Branching %v", d.Branch)
		if err = repo.Branch(d.Branch); err != nil {
			return err
		}
		files := []string{
			TKLTEST_CONFIG_FILE,
			TKLTEST_LOG_FILE,
			fmt.Sprintf(TKLTEST_TESTDIR_FMT, tkltest.appName),
		}

		addon.Activity("Committing")
		if err = repo.Commit(files, COMMIT_MSG); err != nil {
			return err
		}

		return nil
	})
}
