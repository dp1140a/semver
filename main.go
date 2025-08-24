/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/dp1140a/semver/cmd"
	_ "github.com/dp1140a/semver/cmd/bump"
	_ "github.com/dp1140a/semver/cmd/set"
	_ "github.com/dp1140a/semver/cmd/version"
)

func main() {
	cmd.Execute()
}
