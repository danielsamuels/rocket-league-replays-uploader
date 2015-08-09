package main

import (
	"fmt"
	"github.com/andlabs/ui"
	"github.com/cratonica/trayhost"
	"github.com/go-fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"runtime"
)

var window ui.Window

func generate_ui() {

}

func main() {
	// Before we do anything, check if the user wants to log in.
	go ui.Do(func() {
		name := ui.NewTextField()
		button := ui.NewButton("Greet")
		greeting := ui.NewLabel("")
		stack := ui.NewVerticalStack(
			ui.NewLabel("Enter your name:"),
			name,
			button,
			greeting)
		window = ui.NewWindow("Hello", 200, 100, stack)
		button.OnClicked(func() {
			greeting.SetText("Hello, " + name.Text() + "!")
		})
		window.OnClosing(func() bool {
			ui.Stop()
			return true
		})
		window.Show()
	})

	fmt.Println("Starting UI")
	err := ui.Go()
	if err != nil {
		panic(err)
	}

	// EnterLoop must be called on the OS's main thread
	fmt.Println("Locking OS Thread")
	runtime.LockOSThread()

	// Enter the host system's event loop
	fmt.Println("Starting trayhost")
	iconData, _ := ioutil.ReadFile("img/logo.png")
	trayhost.EnterLoop("Rocket League Replays Uploader", iconData)

	go func() {
		// Run your application/server code in here. Most likely you will
		// want to start an HTTP server that the user can hit with a browser
		// by clicking the tray icon.
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op == 1 {
						log.Println("New file created:", event.Name)
					}
				case err := <-watcher.Errors:
					log.Println("error:", err)
				}
			}
		}()

		err = watcher.Add("./testdir")
		if err != nil {
			log.Fatal(err)
		}
		<-done

		// Be sure to call this to link the tray icon to the target url
		trayhost.SetUrl("http://www.rocketleaguereplays.com")
	}()

	// This is only reached once the user chooses the Exit menu item
	fmt.Println("Exiting")
}
