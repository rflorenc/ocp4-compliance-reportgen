package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	rgexec "github.com/rflorenc/ocp4-compliance-reportgen/exec"
	"github.com/rflorenc/ocp4-compliance-reportgen/report"
)

func pollReports() {
	for {
		time.Sleep(time.Duration(10) * time.Second)
		report.GenerateReportHTML()
	}
}

func main() {
	termChan := make(chan os.Signal)
	done := make(chan bool, 1)
	signal.Notify(termChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-termChan
		rgexec.Shutdown()
		fmt.Printf("Caught %v", sig)
		done <- true
	}()

	rgexec.ExecuteCmd("rm", "-f", filepath.Join(report.RootHTMLDir, "css"))
	rgexec.ExecuteCmd("ln", "-s", "/opt/nginx/css", filepath.Join(report.RootHTMLDir, "css"))

	report.GenerateReportHTML()

	rgexec.StartWebServer()
	pollReports()

	fmt.Println("Awaiting signal...")
	<-done
	fmt.Println("Quit.")
}
