package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println("file2tcp")

	var (
		file     string
		host     string
		port     uint
		interval bool
		crlf     bool
	)
	flag.StringVar(&file, "f", "", "File")
	flag.StringVar(&host, "h", "127.0.0.1", "Hostname")
	flag.UintVar(&port, "p", 0, "Port")
	flag.BoolVar(&interval, "i", false, "Add interval of 1 sec after every line.")
	flag.BoolVar(&crlf, "crlf", false, "Pass CRLF as EOL character. Default is LF.")
	flag.Parse()

	// Validations:
	if file == "" {
		fmt.Println("No file given.")
		os.Exit(1)
	}
	if port == 0 {
		fmt.Println("Port is not set.")
		os.Exit(1)
	}

	// Read and send to socket:
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		fmt.Println("Error opening TCP")
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Conn est")

	fp, err := os.Open(file)
	if err != nil {
		fmt.Println("File open error")
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	eol := "\n"
	if crlf {
		eol = "\r\n"
	}
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Println(txt)
		fmt.Fprint(conn, txt+eol)
		if interval {
			time.Sleep(1 * time.Second)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error")
		panic(err)
	}
}
