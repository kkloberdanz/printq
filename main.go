/*
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
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

	time.Sleep(10 * time.Minute)
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("failed to remove file: '%s'\n", filename)
	}
}

func receive_file(writer http.ResponseWriter, request *http.Request) {
	// you can upload file with:
	// curl -X POST -F "file=@src/main.go" localhost:8080/print
	file, header, err := request.FormFile("file")
	if err != nil {
		http.Error(
			writer,
			"file upload requires key of 'file'",
			http.StatusBadRequest,
		)
		fmt.Fprintf(writer, "error\n")
		return
	}
	fmt.Printf("file name '%s'\n", header.Filename)
	extension := filepath.Ext(header.Filename)
	output_name := uuid.Must(uuid.NewV4(), nil)
	output_path := fmt.Sprintf(".printq/%s%s", output_name, extension)
	outf, err := os.Create(output_path)
	if err != nil {
		fmt.Println("failed to create output file")
		return
	}
	io.Copy(outf, file)
	outf.Close()
	file.Close()
	go print_file(output_path)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/print", receive_file)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
