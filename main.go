package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func printErr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func readFromServer(conn net.Conn) {
	rdr := bufio.NewReader(conn)
	for true {
		str, err := rdr.ReadString('\n')
		if err != nil {
			printErr("Error reading")
			panic(err)
		}
		fmt.Println(str)
	}
}

func main() {
	var (
		file     string
		host     string
		port     uint
		interval bool
		crlf     bool
		verbose  bool
	)
	flag.StringVar(&file, "f", "", "File.")
	flag.StringVar(&host, "h", "127.0.0.1", "Hostname.")
	flag.UintVar(&port, "p", 0, "Port")
	flag.BoolVar(&interval, "i", false, "Add interval of 1 sec after every line.")
	flag.BoolVar(&crlf, "crlf", false, "Pass CRLF as EOL character. Default is LF.")
	flag.BoolVar(&verbose, "v", false, "Verbose.")
	flag.Parse()

	// Validations:
	if file == "" {
		printErr("No file given.")
		os.Exit(1)
	}
	if port == 0 {
		printErr("Port is not set.")
		os.Exit(1)
	}

	// Read and send to socket:
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		printErr("Error opening TCP")
		panic(err)
	}
	defer conn.Close()

	go readFromServer(conn)

	fp, err := os.Open(file)
	if err != nil {
		printErr("File open error")
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
		if verbose {
			fmt.Println(txt)
		}
		fmt.Fprint(conn, txt+eol)
		if interval {
			time.Sleep(1 * time.Second)
		}
	}
	if err := scanner.Err(); err != nil {
		printErr("Scanner error")
		panic(err)
	}
}
