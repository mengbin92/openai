package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var (
	gid int
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("file:%s has changed\n", event.Name)
				stopWork()
				startWork()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("error: %s", err)
			}
		}
	}()

	go startWork()

	err = watcher.Add("./conf")
	if err != nil {
		log.Fatal("Add failed:", err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown watcher")
}

func startWork() {
	cmd := exec.Command("./openai")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		log.Printf("start apiserver error: %s", err.Error())
		os.Exit(-1)
	}
	gid = cmd.Process.Pid
}

func stopWork() {
	cmd := exec.Command("kill", "-9", strconv.Itoa(gid))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		log.Printf("start apiserver error: %s", err.Error())
		os.Exit(-1)
	}
}
