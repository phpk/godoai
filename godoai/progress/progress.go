package progress

import (
	"fmt"
	"godoai/libs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

type Process struct {
	Name     string
	Running  bool
	ExitCode int
	Cmd      *exec.Cmd
}

var (
	processesMu sync.RWMutex
	processes   = make(map[string]*Process)
)

func RegisterProcess(name string, cmdstr *exec.Cmd) {
	processesMu.Lock()
	defer processesMu.Unlock()
	processes[name] = &Process{
		Name:    name,
		Running: true,
		Cmd:     cmdstr,
	}
}
func GetCmd(name string) *Process {
	processesMu.Lock()
	defer processesMu.Unlock()
	return processes[name]
}

// StartCmd 执行指定名称的脚本。
// 参数：
// name - 脚本的名称。
// 返回值：
// 返回可能遇到的错误，如果执行成功，则返回nil。
func StartCmd(name string) error {
	info, ok := processes[name]
	if ok && info.Running {
		return fmt.Errorf("process information for '%s' is runing", name)
	}
	scriptPath, err := libs.GetCmdPath(name, name)
	if err != nil {
		return err
	}
	//log.Printf("ollama path %s", libs.GetOllamaModelDir())
	// 设置并启动脚本执行命令
	var cmd *exec.Cmd
	switch name {
	case "ollama":
		os.Setenv("OLLAMA_HOST", "localhost:56715")
		os.Setenv("OLLAMA_MODELS", libs.GetOllamaModelDir())
		os.Setenv("OLLAMA_ORIGINS", "wails://*")
		if runtime.GOOS == "windows" {
			runDir, err := libs.GetRunDir()
			if err != nil {
				return err
			}
			os.Setenv("OLLAMA_RUNNERS_DIR", filepath.Join(runDir, name, "ollama_runners"))
		}
		params := []string{
			"serve",
		}
		cmd = exec.Command(scriptPath, params...)
	default:
		cmd = exec.Command(scriptPath)
	}

	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = SetHideConsoleCursor(cmd)
	}
	go func() {
		// 启动脚本命令并返回可能的错误
		if err := cmd.Start(); err != nil {
			log.Printf("failed to start process %s: %v", name, err)
			return
		}
		RegisterProcess(name, cmd)
		// 等待命令完成
		if err := cmd.Wait(); err != nil {
			log.Printf("command failed for %s: %v", name, err)
		} else {
			log.Printf("%s command completed successfully", name)
		}

		// 命令完成后，更新进程信息
		processesMu.Lock()
		defer processesMu.Unlock()
		if p, ok := processes[name]; ok {
			p.Running = false
			p.ExitCode = cmd.ProcessState.ExitCode()
		}
	}()
	return nil
}
func StopCmd(name string) error {
	cmd, ok := processes[name]
	if !ok {
		return fmt.Errorf("process information for '%s' not found", name)
	}

	// 停止进程并更新status
	if err := cmd.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process %s: %v", name, err)
	}
	//delete(processes, name) // 更新status，表示进程已停止
	cmd.Running = false
	return nil
}
func RestartCmd(name string) error {
	if err := StopCmd(name); err != nil {
		return err
	}
	return StartCmd(name)
}
func StopAllCmd() error {
	processesMu.Lock()
	defer processesMu.Unlock()

	for name, cmd := range processes {
		if err := cmd.Cmd.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("failed to stop process %s: %v", name, err)
		}
	}
	return nil
}
