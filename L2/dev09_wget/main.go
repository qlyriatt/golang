package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func urlToFilename(url string, contentType string) string {
	cType, _, _ := strings.Cut(contentType, ";")
	_, fileExt, _ := strings.Cut(cType, "/")

	filename := url
	filename = strings.TrimPrefix(filename, "https://")
	filename = strings.TrimPrefix(filename, "http://")
	filename = strings.TrimPrefix(filename, "ftp://")
	filename = strings.ReplaceAll(filename, "/", "")

	return filename + "." + fileExt
}

func wget(args args, urls []string) {

	errors := os.Stderr

	if args.quiet {
		errors, _ = os.Create(os.DevNull)
		goto fetch
	}

	if args.logTo != "" {
		file, err := os.Create(args.logTo)
		if err != nil {
			fmt.Fprintln(errors, err)
		} else {
			errors = file
		}
		defer file.Close()
		goto fetch
	}

	if args.isBG {
		file, err := os.Create("wgetlog.txt")
		if err != nil {
			fmt.Fprintln(errors, err)
		} else {
			errors = file
		}
		defer file.Close()
		goto fetch
	}

fetch:

	if args.urlsFile != "" && len(urls) != 0 {
		fmt.Fprintln(errors, "Встречены лишние аргументы")
		return
	}

	if args.urlsFile != "" {
		file, err := os.Open(args.urlsFile)
		if err != nil {
			fmt.Fprintln(errors, "Файл с url не существует")
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if scanner.Err() != nil {
				fmt.Fprintln(errors, err)
				return
			}
			urls = append(urls, scanner.Text())
		}
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		url := url
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintln(errors, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintln(errors, "bad status:"+resp.Status)
				return
			}

			filename := urlToFilename(url, resp.Header.Get("content-type"))

			file, err := os.Create(filename)
			if err != nil {
				fmt.Fprintln(errors, err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				return
			}
		}()
	}

	wg.Wait()
}

type args struct {
	isBG     bool
	BG       bool
	quiet    bool
	logTo    string
	urlsFile string
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Недостаточно аргументов")
		return
	}

	var args args
	flag.BoolVar(&args.isBG, "restart", false, "")
	flag.BoolVar(&args.BG, "b", false, "run wget in background")
	flag.BoolVar(&args.quiet, "q", false, "supress output")
	flag.StringVar(&args.logTo, "o", "", "custom file to log to")
	flag.StringVar(&args.urlsFile, "i", "", "file to read urls from")
	flag.Parse()

	if args.BG && !args.isBG {
		name, err := os.Executable()
		if err != nil {
			fmt.Println(err)
			return
		}

		passArgs := flag.Args()
		passArgs = append(passArgs, "-restart")

		cmd := exec.Command(name, passArgs...)
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("PID:", cmd.Process.Pid)
		return
	}

	wget(args, flag.Args())
}
