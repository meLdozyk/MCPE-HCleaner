package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var choice int

func CMD(title string) {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		panic(err)
	}
	defer syscall.FreeLibrary(kernel32)

	proc, err := syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
	if err != nil {
		panic(err)
	}

	syscall.Syscall(uintptr(proc), 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
}

func main() {
	for {
		titles := []string{"Panic cleaner by meLdozy (meldozy.devv)"}
		i := 0
		go func() {
			for {
				title := titles[i%len(titles)]
				CMD(title)
				i++
				time.Sleep(500 * time.Millisecond)
			}
		}()
		for {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
			fmt.Println("\033[32m")
			fmt.Println(`
██████╗  █████╗ ███╗   ██╗██╗ ██████╗
██╔══██╗██╔══██╗████╗  ██║██║██╔════╝
██████╔╝███████║██╔██╗ ██║██║██║     
██╔═══╝ ██╔══██║██║╚██╗██║██║██║     
██║     ██║  ██║██║ ╚████║██║╚██████╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝ ╚═════╝
                                     
`)
			fmt.Println("1. Фулл очиста")
			fmt.Println("2. Автозакрытие програм")
			fmt.Println("3. Вернуть папки")
			fmt.Println("P.S: если папки не убрались, то просто скройте папку в проводнике.")
			fmt.Print("Выберите действие: ")
			fmt.Scanln(&choice)
			processesToClose := []string{"lastactivityview.exe", "everything.exe", "shellbag_analyzer_cleaner", "SystemSettings.exe"}
			if choice == 2 {
				for {
					for _, processName := range processesToClose {
						if IsRunning(processName) {
							kill(processName)
						} else {

						}
					}
					time.Sleep(1 * time.Second)
				}
			}
			var (
				sourceDir      = filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages", "Microsoft.MinecraftUWP_8wekyb3d8bbwe", "RoamingState")
				destinationDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages", "Microsoft.MinecraftUWP_8wekyb3d8bbwe", "LocalCM")
			)

			if choice == 1 {
				fmt.Println("Идет поиск...")
				rootDir := "C:\\"
				targetFiles := map[string]string{
					"Horion.exe":         "Horion.exe",
					"Vector.dll":         "Vector.dll",
					"Prax.dll":           "Prax.dll",
					"Panic.exe":          "Panic.exe",
					"Horion.dll":         "Horion.dll",
					"HorionInjector.exe": "HorionInjector.exe",
				}

				findAndRenameFiles(rootDir, targetFiles)
				err := Create(destinationDir)
				if err != nil {
					fmt.Println("Error creating directory:", err)
					return
				}
				err = moveFiles(sourceDir, destinationDir)
				if err != nil {
					fmt.Println("Error moving files:", err)
					return
				}
				err = rename(destinationDir)
				if err != nil {
					fmt.Println("Error renaming files:", err)
					return
				}
				err = move(destinationDir, sourceDir)
				if err != nil {
					fmt.Println("Error moving files back:", err)
					return
				}
				cmd1 := exec.Command("cmd", "/c", "start", "Tool1.bat")
				cmd1.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
				err1 := cmd1.Run()
				if err1 != nil {
					fmt.Println("Error:", err1)
					return
				}
			}
			if choice == 3 {
				err := test(sourceDir, destinationDir)
				if err != nil {
					fmt.Println("Error", err)
				}
			}
		}
	}

}
func Create(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func move(sourceDir, destinationDir string) error {
	fileInfos, err := ioutil.ReadDir(destinationDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		oldPath := filepath.Join(destinationDir, fileInfo.Name())
		newPath := filepath.Join(sourceDir, fileInfo.Name())
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		fmt.Printf("Перемещен файл: %s -> %s\n", oldPath, newPath)
	}
	return nil
}
func test(sourceDir, destinationDir string) error {
	fileList, err := filepath.Glob(destinationDir + "/*")
	if err != nil {
		return err
	}
	for _, file := range fileList {
		newPath := filepath.Join(sourceDir, remove(filepath.Base(file)))
		err = os.Rename(file, newPath)
		if err != nil {
			return err
		}
		fmt.Printf("Перемещен файл: %s -> %s\n", file, newPath)
	}
	return nil
}

func remove(s string) string {
	var b strings.Builder
	for _, c := range s {
		if c != '0' {
			b.WriteRune(c)
		}
	}
	return b.String()
}
func rename(dirPath string) error {
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		oldPath := filepath.Join(dirPath, fileInfo.Name())
		newName := addZerosToString(fileInfo.Name())
		newPath := filepath.Join(filepath.Dir(oldPath), newName)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		fmt.Printf("Переименован файл: %s -> %s\n", oldPath, newPath)
	}
	return nil
}

func addZerosToString(s string) string {
	var b strings.Builder
	for i, c := range s {
		b.WriteRune(c)
		if i < len(s)-1 {
			b.WriteString("0")
		}
	}
	return b.String()
}
func moveFiles(sourceDir, destinationDir string) error {
	fileList, err := filepath.Glob(sourceDir + "/*")
	if err != nil {
		return err
	}
	for _, file := range fileList {
		newPath := filepath.Join(destinationDir, filepath.Base(file))
		err = os.Rename(file, newPath)
		if err != nil {
			return err
		}
		fmt.Printf("Перемещен файл: %s -> %s\n", file, newPath)
		fmt.Scanln()
	}
	return nil
}
func IsRunning(processName string) bool {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	for _, proc := range processes {
		name, _ := proc.Name()
		if strings.EqualFold(name, processName) {
			return true
		}
	}
	return false
}

func kill(processName string) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for _, proc := range processes {
		name, _ := proc.Name()
		if strings.EqualFold(name, processName) {
			pid := proc.Pid
			err := proc.Terminate()
			if err != nil {
				fmt.Printf("Error terminating %s (PID %d): %s\n", processName, pid, err)
			} else {
				fmt.Printf("%s (PID %d) terminated\n", processName, pid)
			}
		}
	}
}
func findAndRenameFiles(rootDir string, targetFiles map[string]string) {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			fileName := filepath.Base(path)
			if newName, ok := targetFiles[fileName]; ok {
				log.Printf("Found file: %s\n", path)

				randomName := generateRandomName(6)
				dir := filepath.Dir(path)
				newFilePath := filepath.Join(dir, randomName+filepath.Ext(newName))

				err := os.Rename(path, newFilePath)
				if err != nil {
					log.Printf("Error renaming file: %v\n", err)
				} else {
					log.Printf("File renamed to: %s\n", newFilePath)
				}
			}
		}
		return nil
	})
}

func generateRandomName(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomName := make([]byte, length)
	for i := range randomName {
		randomName[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomName)
}
