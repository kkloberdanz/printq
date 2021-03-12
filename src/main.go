package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello %s!\n", request.URL.Path[1:])
}

func print_file(filename string) {
	cmd := exec.Command("lp", filename)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to print file: '%s'\n", filename)
	} else {
		fmt.Printf("printing: '%s'\n", filename)
	}

	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("failed to remove file: '%s'\n", filename)
	}
}

func recieve_file(writer http.ResponseWriter, request *http.Request) {
	// you can upload file with:
	// curl -X POST -F "file=@src/main.go" localhost:8080/print
	file, header, err := request.FormFile("file")
	if err != nil {
		// TODO: how to send 4xx status code in response?
		fmt.Fprintf(writer, "error\n")
		return
	}
	defer file.Close()
	fmt.Printf("file name '%s'\n", header.Filename)
	extension := filepath.Ext(header.Filename)
	output_name := uuid.Must(uuid.NewV4())
	output_path := fmt.Sprintf("printq/%s%s", output_name, extension)
	outf, err := os.Create(output_path)
	if err != nil {
		fmt.Println("failed to create output file")
		return
	}
	defer outf.Close()
	io.Copy(outf, file)
	go print_file(output_path)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/print", recieve_file)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
