package main

import (
	"fmt"
	"go/build"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

const maxN = 10000000
const maxBufferSize = 1024 * 1024

func main() {
	if err := runTest(); err != nil {
		log.Fatalf("Test failed: %v", err)
	}
	log.Printf("Test passed!")
}

func runTest() error {
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}

	log.Printf("Compiling kafka-go")
	out, err := exec.Command("go", "install", "-v", "github.com/hou-siyuan-coding/kafka-go").
		CombinedOutput()
	if err != nil {
		log.Printf("Faild to build: %v", err)
		return fmt.Errorf("compilation failed: %v (out: %s)", err, string(out))
	}

	port := 7357
	dbPath := "/tmp/kafka.db"
	os.Remove(dbPath)

	log.Printf("Running kafka-go on port %d", port)

	cmd := exec.Command(goPath+"/bin/kafka-go", "-filename="+dbPath, fmt.Sprintf("-port=%d", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Start()
	defer cmd.Process.Kill()

	log.Printf("Waiting for the port localhost:%d to open", port)
	for i := 0; i <= 100; i++ {
		timeout := time.Millisecond * 50
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprint(port)), timeout)
		if err != nil {
			time.Sleep(timeout)
			continue
		}
		conn.Close()
		break
	}

	return nil
}
