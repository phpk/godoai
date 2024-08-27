package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"

	"godoai/cmd"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cmd.Start()
}
func (a *App) shutdown(ctx context.Context) {
	cmd.Stop()
}
func (a *App) OpenDirDialog() string {
	path, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Select Folder",
	})
	if err != nil {
		wruntime.LogErrorf(a.ctx, "Error: %+v\n", err)
	}
	return path
}

func (a *App) RestartApp() error {
	name, err := os.Executable()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command(name, os.Args[1:]...)
	case "darwin": // macOS
		cmd = exec.Command("/usr/bin/open", name)
	case "linux":
		cmd = exec.Command(name, os.Args[1:]...)
		// Optionally, you could use 'xdg-open' or 'gnome-open' etc.
		// cmd = exec.Command("/usr/bin/gnome-open", name)
	default:
		return errors.New("unsupported OS")
	}

	if cmd != nil {
		cmd.Start()
		wruntime.Quit(a.ctx)
		return nil
	}

	return errors.New("failed to restart application")
}
