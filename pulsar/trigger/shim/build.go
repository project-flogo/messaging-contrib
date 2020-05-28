//+build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	fmt.Println("Running build script for the Pulsar Flogo trigger")

	//Get the dir where build.go is present
	appDir, err := os.Getwd()
	if err != nil {
		fmt.Printf(err.Error())
	}

	var cmd *exec.Cmd

	// Clean up
	fmt.Println("Cleaning up previous executables")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "del", "/q", "pflogoFunc")
	} else {
		cmd = exec.Command("rm", "pflogoFunc")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Build an executable for Linux
	fmt.Println("Building a new handler Bin")
	cmd = exec.Command("go", "build", "-o", "pflogoFunc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(appDir)

	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

	pFlogoFuncPath := filepath.Join(cmd.Dir, "pflogoFunc")
	err = CopyFile(pFlogoFuncPath, filepath.Join(cmd.Dir, "..", "bin", "pflogoFunc"))
	if err != nil {
		fmt.Println("Failed to copy zip file to bin: %v", err)
	}

}

func CopyFile(srcFile, destFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
