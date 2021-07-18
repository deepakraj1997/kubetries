package main

import (
	"fmt"

	"github.com/deepakraj1997/kubetries/nginx"
)

func main() {
	fmt.Printf("Deploying Nginx Stateless")
	// nginx.DeployNginxStateless()
	// nginx.BackupNginxStateless()
	nginx.RestoreNginxStateless()
}
