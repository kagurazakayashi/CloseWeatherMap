package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func getHostsFilePath() string {
	switch runtime.GOOS {
	case "windows":
		systemRoot := os.Getenv("SystemRoot")
		if systemRoot == "" {
			systemRoot = "C:\\Windows"
		}
		return systemRoot + "\\System32\\drivers\\etc\\hosts"
	default:
		systemRoot := os.Getenv("SystemRoot")
		if systemRoot == "" {
			systemRoot = "/"
		}
		return systemRoot + "etc/hosts"
	}
}

func entryExists(entry, filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), entry) {
			return true
		}
	}
	return false
}

func ensureTrailingNewline(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(content) > 0 && content[len(content)-1] != '\n' {
		return os.WriteFile(filePath, append(content, '\n'), 0644)
	}
	return nil
}

func addEntryToHosts(entry, filePath string) {
	if entryExists(entry, filePath) {
		log.Println("hosts 条目已经存在:", entry)
		return
	}

	if err := ensureTrailingNewline(filePath); err != nil {
		fmt.Println("错误：不能写入 hosts 文件:", err)
		return
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("错误：不能打开 hosts 文件:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(entry + "\n"); err != nil {
		log.Println("错误：不能写入 hosts 文件:", err)
		return
	}

	log.Println("已添加到 hosts 文件:", entry)
}

func removeEntryFromHosts(entry, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("错误：不能打开 hosts 文件:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, entry) {
			continue
		}
		lines = append(lines, line)
	}

	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("错误：不能打开 hosts 文件:", err)
		return
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			log.Println("错误：不能写入 hosts 文件:", err)
			return
		}
	}

	log.Println("已从 hosts 文件中删除:", entry)
}

func hostsAdd() {
	if len(hostEntry) == 0 {
		return
	}
	hostsFilePath := getHostsFilePath()
	addEntryToHosts(hostEntry, hostsFilePath)
}

func hostsRm() {
	if len(hostEntry) == 0 {
		return
	}
	hostsFilePath := getHostsFilePath()
	removeEntryFromHosts(hostEntry, hostsFilePath)
}
