package exec

import (
	"fmt"
	"log"
	"os/exec"

	rgerr "github.com/rflorenc/ocp4-compliance-reportgen/error"
)

// ExecuteCmd call wraps os/exec.Command
func ExecuteCmd(cmdName string, cmdArgs ...string) {
	var err error

	cmd := exec.Command(cmdName, cmdArgs...)
	err = cmd.Run()
	rgerr.CheckError(fmt.Sprintf("%s.Run()", cmdName), err)
}

// Shutdown kills all nginx processes
func Shutdown() {
	ExecuteCmd("killall", "nginx")
	log.Printf("Stopped NGINX")
}

// StartWebServer starts nginx
func StartWebServer() {
	log.Printf("Starting NGINX")
	ExecuteCmd("nginx", "-c", "/opt/nginx/nginx.conf")
}
