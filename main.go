package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	systray.Run(onReady, onExit)
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func transformEncoding(rawReader io.Reader, trans transform.Transformer) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(rawReader, trans))
	if err == nil {
		return string(ret), nil
	} else {
		return "", err
	}
}

func FromShiftJIS(str string) (string, error) {
	return transformEncoding(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
}

var usr, _ = user.Current()
var path = usr.HomeDir
var filename = "obs-currentsong.txt"
var logfilepath = filepath.Join(path, filename)

func onReady() {
	systray.SetIcon(Red)
	systray.SetTitle("OBS Current Song")
	systray.SetTooltip("OBS Current Song")
	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	for {
		cmd := exec.Command("powershell.exe", "((get-process -processname Spotify | select -unique -property MainWindowTitle) | Where-Object -Property mainwindowtitle -notlike $null | where-object -property mainwindowtitle -notlike *Spotify*).mainwindowtitle")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, _ := cmd.Output()
		currentsong := string(out)
		currentsongtrimmed := strings.TrimSpace(currentsong)
		convertsongtitletoutf8, _ := FromShiftJIS(currentsongtrimmed)
		fmt.Println(convertsongtitletoutf8)
		logfile, err := os.Create(logfilepath)
		defer logfile.Close()

		if convertsongtitletoutf8 != "" {
			systray.SetIcon(Green)
			w := bufio.NewWriter(logfile)
			_, err = fmt.Fprintf(w, "%s", convertsongtitletoutf8)
			check(err)
			w.Flush()
		} else {
			systray.SetIcon(Red)
			w := bufio.NewWriter(logfile)
			_, err = fmt.Fprintf(w, "%s", "")
			check(err)
			w.Flush()
		}
		time.Sleep(time.Second)
	}
}

func onExit() {
	// clean up here
}
