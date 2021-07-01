package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.Lshortfile)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file", err)
	}

	host := os.Getenv("COMPRESS_FILE_HOST")
	token := os.Getenv("COMPRESS_FILE_TOKEN")
	path := "/version/all/global"
	version := getVersion(host + path + "?token=" + token)

	tempFile := "./compress_file_temp.tar.gz"
	fileURL := host + path + "/" + version + "?token=" + token
	err = DownloadFile(tempFile, fileURL)
	if err != nil {
		log.Fatal("download file error, ", err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("get Home dir error: ", err)
	}
	basePath := homeDir + "/Downloads/compressfile_" + version + "/"
	log.Println("download to ", basePath)

	err = decompress(tempFile, basePath)
	if err != nil {
		log.Fatal("decompress base file error, ", err)
	}

	toDecompress := []struct {
		filename string
		subPath  string
	}{
		{"dns", "dns/"},
		{"leadns", "proxy/cn/"},
		{"leadns", "proxy/global/"},
		{"leadns", "web/cn/"},
		{"leadns", "web/global/"},
	}

	for _, t := range toDecompress {
		file := basePath + t.subPath + t.filename + ".tar.gz"
		err := decompress(file, basePath+t.subPath)
		if err != nil {
			log.Fatal("decompress file error, ", err)
		}
		err = os.Remove(file)
		if err != nil {
			log.Fatal("remove file error, ", err)
		}
	}

	err = os.Remove(tempFile)
	if err != nil {
		log.Fatal("remove temp file error, ", err)
	}

	log.Printf("donwload compress file completed, version: %s\n", version)
}

func getVersion(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("get version error, ", err)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read response bode error, ", err)
	}

	return strings.Trim(strings.Split(string(bytes), "\n")[0], "version:")
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// copy from https://studygolang.com/articles/7481
//解压 tar.gz
func decompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return fmt.Errorf("os.Open failed %v", err)
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("gzip.NewReader failed %v", err)
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + header.Name
		file, err := createFile(filename)
		if err != nil {
			//ignore
			//fmt.Printf("createFile failed %v\n", err)
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
