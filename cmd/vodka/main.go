package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

const banner = `
__/\\\________/\\\_______/\\\\\_______/\\\\\\\\\\\\_____/\\\________/\\\_____/\\\\\\\\\____        
 _\/\\\_______\/\\\_____/\\\///\\\____\/\\\////////\\\__\/\\\_____/\\\//____/\\\\\\\\\\\\\__       
  _\//\\\______/\\\____/\\\/__\///\\\__\/\\\______\//\\\_\/\\\__/\\\//______/\\\/////////\\\_      
   __\//\\\____/\\\____/\\\______\//\\\_\/\\\_______\/\\\_\/\\\\\\//\\\_____\/\\\_______\/\\\_     
    ___\//\\\__/\\\____\/\\\_______\/\\\_\/\\\_______\/\\\_\/\\\//_\//\\\____\/\\\\\\\\\\\\\\\_    
     ____\//\\\/\\\_____\//\\\______/\\\__\/\\\_______\/\\\_\/\\\____\//\\\___\/\\\/////////\\\_   
      _____\//\\\\\_______\///\\\__/\\\____\/\\\_______/\\\__\/\\\_____\//\\\__\/\\\_______\/\\\_  
       ______\//\\\__________\///\\\\\/_____\/\\\\\\\\\\\\/___\/\\\______\//\\\_\/\\\_______\/\\\_ 
        _______\///_____________\/////_______\////////////_____\///________\///__\///________\///__
`

var (
	appCmd *exec.Cmd
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println(Cyan + banner + Reset)
		log.Println(Blue + "Vodka Watching for Empty Glasses..." + Reset)
		watchBackend()
		return
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println(Red + "Usage: vodka create <project-name>" + Reset)
			return
		}
		createProject(os.Args[2])

	case "run":
		if len(os.Args) >= 3 && os.Args[2] == "dev" {
			fmt.Println(Cyan + banner + Reset)
			log.Println(Blue + "Starting Full-Stack Dev Environment..." + Reset)
			runDev()
		} else {
			fmt.Println(Red + "Usage: vodka run dev" + Reset)
		}

	default:
		fmt.Printf(Red+"Unknown command '%s'.\n"+Reset, os.Args[1])
		fmt.Println(Cyan + "Available commands:" + Reset + "\n  " + Green + "vodka" + Reset + "\n  " + Green + "vodka create <name>" + Reset + "\n  " + Green + "vodka run dev" + Reset)
	}
}

func runDev() {
	go func() {
		frontendDir := filepath.Join(".", "frontend")

		if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
			log.Println(Red + "Error: 'frontend' directory not found. Are you in a vodka project?" + Reset)
			return
		}

		log.Println(Cyan + "Starting Vite Frontend..." + Reset)

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", "npm run dev")
		} else {
			cmd = exec.Command("npm", "run", "dev")
		}

		cmd.Dir = frontendDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Println(Red+"Frontend server crashed:"+Reset, err)
		}
	}()

	time.Sleep(2 * time.Second)
	watchBackend()
}

func createProject(name string) {
	fmt.Printf(Cyan+"Distilling your project: %s...\n"+Reset, name)

	os.Mkdir(name, 0755)

	fmt.Println(Gray + "Initializing Go backend..." + Reset)
	runCmd(name, "go", "mod", "init", name)
	runCmd(name, "go", "get", "github.com/DevanshuTripathi/vodka@latest")

	mainGoContent := `package main

import (
	"github.com/DevanshuTripathi/vodka"
	"` + name + `/routes"
)

func main() {
	app := vodka.DefaultRouter() // Creates a Default Router with Logger and Recovery Middleware

	allowedOrigins := []string{"http://localhost:5173"}

	app.Use(vodka.AllowCORS(allowedOrigins))

	routes.Setup(app)

	app.Run(":8080")
}
`
	routesContent := `package routes

import (
	"github.com/DevanshuTripathi/vodka"
	"` + name + `/controllers"
)

func Setup(app *vodka.Engine) {
	app.GET("/ping", controllers.Pong)

	app.GET("/hello/:name", controllers.Hello)
}
`
	controllersContent := `package controllers

import (
	"github.com/DevanshuTripathi/vodka"
)

func Pong(c *vodka.Context) {
	c.String(200, "Pong!")
}

func Hello(c *vodka.Context) {
	name := c.Param("name")

	c.String(200, "Hello "+ name +"!")
}
`
	os.MkdirAll(filepath.Join(name, "controllers"), 0755)
	os.MkdirAll(filepath.Join(name, "routes"), 0755)

	os.WriteFile(filepath.Join(name, "main.go"), []byte(mainGoContent), 0644)
	os.WriteFile(filepath.Join(name, "controllers", "ping.go"), []byte(controllersContent), 0644)
	os.WriteFile(filepath.Join(name, "routes", "routes.go"), []byte(routesContent), 0644)

	fmt.Println(Gray + "Spinning up React frontend with Vite..." + Reset)
	if runtime.GOOS == "windows" {
		runCmd(name, "cmd", "/C", "npm create vite@latest frontend -- --template react")
	} else {
		runCmd(name, "npm", "create", "vite@latest", "frontend", "--", "--template", "react")
	}

	fmt.Printf(Green+"\nProject %s is ready!\n"+Reset, name)
	fmt.Printf(Cyan+"Next steps:\n"+Reset+"  "+Green+"cd %s\n"+Reset+"  "+Green+"cd frontend && npm install\n"+Reset+"  "+Green+"cd ..\n"+Reset+"  "+Green+"vodka run dev\n"+Reset, name)
}

func runCmd(dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf(Red+"Error running %s: %v\n"+Reset, name, err)
	}
}

func watchBackend() {
	buildAndRun()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if strings.HasSuffix(event.Name, ".go") {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(500*time.Millisecond, func() {
					log.Printf(Yellow+"File changed: %s. Rebuilding...\n"+Reset, event.Name)
					buildAndRun()
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println(Red+"Watcher Error: "+Reset, err)
		}
	}
}

func buildAndRun() {
	if appCmd != nil && appCmd.Process != nil {
		log.Println(Gray + "Stopping old process..." + Reset)
		appCmd.Process.Kill()
		appCmd.Wait()
	}

	os.Mkdir("tmp", 0755)

	binaryPath := filepath.Join("tmp", "vodka-build")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		log.Println(Red + "Build failed. Waiting for changes..." + Reset)
		return
	}

	runPath := "." + string(filepath.Separator) + binaryPath
	appCmd = exec.Command(runPath)
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	if err := appCmd.Start(); err != nil {
		log.Println(Red+"Failed to start app:"+Reset, err)
		return
	}

	log.Println(Green + "App is running!" + Reset)
}

