package banner

import (
	"api-test/src/config"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v4/host"
)

type APIInfo struct {
	Name        string
	Author      string
	Description string
	StartedAt   time.Time
	GitRepo     string
	Environment string
	Port        int
}

func Banner(config *config.Config) {
	apiInfo := APIInfo{
		Name:        "KOSVI",
		Author:      "KOSVI Team",
		Description: "API REST KOSVI",
		StartedAt:   time.Now(),
		GitRepo:     "github.com/kosvi/kosvi",
		Environment: config.Environment.Name,
		Port:        config.Port,
	}
	clearScreen()
	printLogoWithDetails(apiInfo)
	fmt.Println()
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func printLogoWithDetails(apiInfo APIInfo) {
	logoLines := strings.Split(getLogo(), "\n")

	detailsLines := getDetailsLines(apiInfo)

	maxLines := len(logoLines)
	if len(detailsLines) > maxLines {
		maxLines = len(detailsLines)
	}

	fmt.Println()

	spacing := "     "

	for i := 0; i < maxLines; i++ {
		if i < len(logoLines) {
			color.RGB(106, 101, 255).Add(color.Bold).Print(logoLines[i])
		} else {
			fmt.Print(strings.Repeat(" ", len(logoLines[0])))
		}

		fmt.Print(spacing)

		if i < len(detailsLines) {
			fmt.Println(detailsLines[i])
		} else {
			fmt.Println()
		}
	}
}

func getLogo() string {
	return `⢀⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣦⡀
⣾⣿⣿⠉⠉⠉⠉⣿⣿⣿⣿⡿⠋⠉⠉⠉⠙⣿⣿⣧
⣿⣿⣿⠀⠀⠀⣠⣿⣿⡿⠋⠀⠀⠀⠀⣠⣾⣿⣿⣿
⣿⣿⣿⠀⣠⣾⣿⣿⠟⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿
⣿⣿⣿⣿⣿⣿⠟⠁⠀⠀⠀⠠⣾⣿⣿⣿⣿⣿⣿⣿
⣿⣿⣿⣿⠟⠁⠀⠀⠀⠀⠀⠀⠘⢿⣿⣿⣿⣿⣿⣿
⣿⣿⣿⠁⠀⠀⠀⣠⣶⡄⠀⠀⠀⠀⠹⣿⣿⣿⣿⣿
⣿⣿⣿⠀⠀⠀⠀⣿⣿⣿⣦⡀⠀⠀⠀⠈⢿⣿⣿⣿
⢻⣿⣿⣀⣀⣀⣀⣿⣿⣿⣿⣷⣄⣀⣀⣀⣀⣿⣿⡟
⠈⠻⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠟⠁`
}

func getDetailsLines(apiInfo APIInfo) []string {
	white := color.New(color.FgHiWhite).Add(color.Bold)
	green := color.New(color.FgGreen).Add(color.Bold)
	blue := color.New(color.FgHiBlue).Add(color.Bold)
	kosviColor := color.RGB(106, 101, 255).Add(color.Bold)
	hostInfo, _ := host.Info()

	var lines []string

	lines = append(lines, "")

	var kosviLine strings.Builder
	kosviColor.Fprintf(&kosviLine, "KOSVI")
	lines = append(lines, kosviLine.String())

	lines = append(lines, "")

	var osLine strings.Builder
	white.Fprint(&osLine, "OS: ")
	osLine.WriteString(fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion))
	lines = append(lines, osLine.String())

	var apiLine strings.Builder
	white.Fprint(&apiLine, "API: ")
	apiLine.WriteString(apiInfo.Name)
	lines = append(lines, apiLine.String())

	var descLine strings.Builder
	white.Fprint(&descLine, "Descripción: ")
	descLine.WriteString(apiInfo.Description)
	lines = append(lines, descLine.String())

	var authorLine strings.Builder
	white.Fprint(&authorLine, "Autor: ")
	authorLine.WriteString(apiInfo.Author)
	lines = append(lines, authorLine.String())

	var repoLine strings.Builder
	white.Fprint(&repoLine, "Repositorio: ")
	repoLine.WriteString(apiInfo.GitRepo)
	lines = append(lines, repoLine.String())

	var envLine strings.Builder
	white.Fprint(&envLine, "Entorno: ")
	green.Fprintf(&envLine, "%s", apiInfo.Environment)
	lines = append(lines, envLine.String())

	var hostLine strings.Builder
	white.Fprint(&hostLine, "Host: ")
	blue.Add(color.Underline).Fprintf(&hostLine, "http://0.0.0.0:%d", apiInfo.Port)
	lines = append(lines, hostLine.String())

	return lines
}
