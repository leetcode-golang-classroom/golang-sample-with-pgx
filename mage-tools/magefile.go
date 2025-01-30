//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// clean the build binary
func Clean() error {
	return sh.Rm("bin")
}

// update the dependency
func Update() error {
	return sh.Run("go", "mod", "download")
}

// setupdata for database
func SetupData() error {
	return sh.Run("docker", "compose", "up", "-d", "postgres")
}

// build Creates the binary in the current directory.
func Build() error {
	mg.Deps(Clean)
	mg.Deps(Update)
	err := sh.Run("go", "build", "-o", "./bin/connection-pool-sample", "./cmd/connection-pool-sample/main.go")
	if err != nil {
		return err
	}
	return nil
}

// LaunchConnectionPoolSample start the connection pool sample
func LaunchConnectionPoolSample() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/connection-pool-sample")
	if err != nil {
		return err
	}
	return nil
}

// run the test
func Test() error {
	err := sh.RunV("go", "test", "-v", "./...")
	if err != nil {
		return err
	}
	return nil
}
