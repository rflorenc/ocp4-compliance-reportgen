package main

import (
	"path/filepath"
	"time"

	rgexec "github.com/rflorenc/ocp4-compliance-reportgen/exec"
	"github.com/rflorenc/ocp4-compliance-reportgen/report"
)

func pollReports() {
	for {
		time.Sleep(10 * time.Second)
		report.GenerateReportHTML()
	}
}

func main() {
	rgexec.SigHandler()
	rgexec.ExecuteCmd("rm", "-f", filepath.Join(report.RootHTMLDir, "css"))
	rgexec.ExecuteCmd("ln", "-s", "/opt/nginx/css", filepath.Join(report.RootHTMLDir, "css"))

	report.GenerateReportHTML()

	rgexec.StartWebServer()
	pollReports()
}
