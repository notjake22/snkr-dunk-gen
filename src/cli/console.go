package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var cmd exec.Cmd

func ClearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "mac":
		exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Start()
	default:
		fmt.Println(runtime.GOOS)
	}
}
