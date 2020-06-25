package main

import (
	"archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//Unzip the file
func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	fmt.Printf("filenames: %v \n", filenames)
	return filenames, nil
}

func gzipFile(uncompressedName string, compressedString string) {
	// Open file on disk.
	name := uncompressedName
	f, _ := os.Open(name)

	// Create a Reader and use ReadAll to get all the bytes from the file.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Replace txt extension with gz extension.
	name = strings.Replace(compressedString, ".zip", ".gz", -1)

	// Open file for writing.
	f, _ = os.Create(name)

	// Write compressed data.
	w := gzip.NewWriter(f)
	w.Name = uncompressedName
	w.Write(content)
	w.Close()

	// Done.
	fmt.Printf("Done for %v \n", name)
}

func main() {
	files, err := ioutil.ReadDir("./") // consider making this windows compatible
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "zip") {
			fmt.Printf("Decompressing: %v", f.Name())
			newFilename, _ := unzip(f.Name(), "")
			if len(newFilename) > 0 {
				gzipFile(newFilename[0], f.Name())
				os.Remove(newFilename[0])
			} else {
				fmt.Printf("No filename found array: %v", newFilename)
			}

		}

	}
}
