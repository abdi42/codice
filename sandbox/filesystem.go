package sandbox

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetFiles(name string) []string {
	var files []string
	root := os.Getenv("TEMP_PATH") + "/" + name

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fmt.Println(filepath.Dir(path))
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func CreatePayload(path string, lang Lang) {
	tempDir := path
	err := os.Mkdir(tempDir, os.ModePerm)
	if err != nil {
		log.Fatal("Cannot create directory", err)
	}

	file, err := os.Create(tempDir + "/" + lang.FileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	text := []byte("")
	if _, err = file.Write(text); err != nil {
		log.Fatal("Failed to write to file", err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	input, err := ioutil.ReadFile("./program_runner.sh")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(tempDir+"/program_runner.sh", input, 0777)
	if err != nil {
		fmt.Println("Error creating", "program_runner")
		fmt.Println(err)
	}
}

func Update(id string, filename string, source string) {
	path := "./temp/" + id + "/" + filename
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if isError(err) {
		return
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		fmt.Println(err.Error())
	}

	_, err = file.WriteString(source)
	if isError(err) {
		return
	}

	err = file.Sync()
	if isError(err) {
		return
	}
}

func GetOutputs(id string) (string, string) {
	stdout, err := ioutil.ReadFile("./temp/" + id + "/log.out")
	if err != nil {
		panic(err)
	}
	stderr, er := ioutil.ReadFile("./temp/" + id + "/error.out")
	if er != nil {
		fmt.Println(er.Error())
	}

	go ClearFiles(id)

	return string(stdout), string(stderr)
}

func ClearFiles(id string) {
	dir := "./temp/" + id

	stdout, err := os.OpenFile(dir+"/log.out", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	defer stdout.Close()

	stderr, err := os.OpenFile(dir+"/error.out", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	defer stderr.Close()

	mainfile, err := os.OpenFile(dir+"/main.rb", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	defer mainfile.Close()

	if err := stdout.Truncate(0); err != nil {
		fmt.Println(err.Error())
	}

	if err := stderr.Truncate(0); err != nil {
		fmt.Println(err.Error())
	}

	if err := mainfile.Truncate(0); err != nil {
		fmt.Println(err.Error())
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
