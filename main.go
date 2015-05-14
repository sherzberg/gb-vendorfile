package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Dependency struct {
	dep string
}

func (dep *Dependency) folder() (string, error) {

	folder := strings.TrimPrefix(dep.dep, "https://")
	folder = folder[:strings.LastIndex(folder, "/archive")]
	return folder, nil
}

func UnzipToPath(zipfile string, folder string) {

	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer reader.Close()

	prefix := ""
	for i, f := range reader.Reader.File {

		zipped, err := f.Open()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer zipped.Close()

		if i == 0 {
			prefix = f.Name
		}

		// get the individual file name and extract the current directory
		path := filepath.Join(folder, "/", strings.TrimPrefix(f.Name, prefix))

		if f.FileInfo().IsDir() {
			if i != 0 {
				// Github archives have a top level dir; skip it
				os.MkdirAll(path, f.Mode())
			}
		} else {
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			defer writer.Close()

			if _, err = io.Copy(writer, zipped); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func (d *Dependency) vendor() {
	folder, err := d.folder()
	if err != nil {
		panic(err)
	}

	vendorPath := filepath.Join("vendor/src", "/", folder)
	c := exec.Command("mkdir", "-p", vendorPath)
	err = c.Run()
	if err != nil {
		panic(err)
	}

	out, err := os.Create("/tmp/_vendorfile.zip")
	defer out.Close()

	resp, err := http.Get(d.dep)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("... vendoring %s", folder)
	UnzipToPath("/tmp/_vendorfile.zip", vendorPath)
}

func main() {

	file, err := os.Open("Vendorfile")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			d := Dependency{
				dep: line,
			}

			d.vendor()

		}
	}
}
