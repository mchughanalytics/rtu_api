package ui

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// OpenBrowser launches the UI
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		c := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		err = c.Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
