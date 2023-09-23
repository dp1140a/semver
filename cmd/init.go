/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dp1140a/semver/util"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a semver project",
	Long: `Will launch an interactive console to launch a semver project.  This must be done in an existing git repo. 
It will create a file called VERSION that will be used to track version information.  If an Existing VERSION file is found it will ask if you want to overwrite.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runInit() {
	cwd, _ := os.Getwd() // Get Current Directory
	//Is this is a git project
	util.InGitDir(cwd) // If not will exit
	reader := bufio.NewReader(os.Stdin)

	/**
	Check if VERSION file exists
	*/
	if util.VersionFileExists(cwd) { // If file exists ask to overwrite

		fmt.Printf("VERSION fle was found do you want to overwrite it [Y/n]?")
		overwrite, _ := reader.ReadString('\n')
		overwrite = strings.TrimSuffix(overwrite, "\n")
		if overwrite == "" {
			overwrite = "y"
		}
		fmt.Println(overwrite)
		if strings.ToLower(overwrite) == "y" {
			fmt.Println("Overwriting VERSION file.  Continuing . . . ")
		} else if strings.ToLower(overwrite) == "n" {
			fmt.Println("Please delete the VERSION file and restart")
			os.Exit(0)
		} else {
			fmt.Println("I don't understand your response. Exiting.")
			os.Exit(0)
		}
	}

	fmt.Printf("Creating Semver for project: %v\n", filepath.Base(cwd))
	isValid := false
	var startingVersion = "0.1.0"

	for isValid != true {
		fmt.Print("Starting Version [0.1.0]: ")
		//reader := bufio.NewReader(os.Stdin)
		var err error
		// ReadString will block until the delimiter is entered
		startingVersion, err = reader.ReadString('\n')
		if startingVersion == "\n" {
			startingVersion = "0.1.0"
		}
		startingVersion = strings.TrimSuffix(startingVersion, "\n")

		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again", err)
			return
		}

		//validate Starting Version
		isValid = util.ValidVersionString(startingVersion)
		if !isValid {
			fmt.Printf("\"%v\" is an invalid semver format.  Must be in the form of X.Y.Z where each is a number.\n", startingVersion)
		} else {
			break
		}
	}
	fmt.Println(startingVersion)
	fmt.Printf("Creating VERSION file with starting version %v\n", startingVersion)
	err := util.WriteVersionFile(startingVersion)
	if err != nil {
		fmt.Printf("Error writing VERSION file: %v. Exiting", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
