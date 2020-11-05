package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	rgerr "github.com/rflorenc/ocp4-compliance-reportgen/error"
)

// ExecuteCmd call wraps os/exec.Command
func ExecuteCmd(cmdName string, cmdArgs ...string) {
	var err error

	cmd := exec.Command(cmdName, cmdArgs...)
	err = cmd.Run()
	rgerr.CheckError(fmt.Sprintf("%s.Run():", cmdName), err)
}

// SigHandler handles signals and calls Shutdown()
func SigHandler() {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	Shutdown()
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
