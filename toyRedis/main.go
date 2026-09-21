package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func runBackup(store *Store) error {
	f, err := os.Open(Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		command, err := CustomSplit(scanner.Text(), " ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(command)
		_, _, err = runmethod(command, store)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func HandleConnection(conn net.Conn, store *Store) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		command, err := CustomSplit(scanner.Text(), " ")
		if err != nil {
			conn.Write([]byte(err.Error()))
		}
		value, append, err := runmethod(command, store)
		if err != nil {
			conn.Write([]byte(err.Error()))
		}
		Append(append)
		conn.Write([]byte(value + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid Input: %s", err)
	}

}

func runmethod(command []string, store *Store) (string, []string, error) {
	append := []string{}
	if len(command) == 0 {
		return "", append, fmt.Errorf("Command Misformatted")
	}
	method := command[0]
	command = command[1:]
	key := command[0]
	switch method {
	case "SET":
		v, err := store.Set(command)
		if err != nil {
			return "", append, err
		}
		if v.t != "" {
			append = []string{"SET", key, command[1], "EX", "(" + v.t + ")"}
		} else {
			append = []string{"SET", key, command[1]}
		}
		if err != nil {
			fmt.Println(err)
			return "", append, err
		}
		return "SET Command Complete", append, nil
	case "GET":
		value, err := store.Get(key)
		if err != nil {
			return "", append, err
		}
		append = []string{"GET", key}
		return key + ":" + value, append, nil
	case "DELETE":
		store.Delete(key)
		append = []string{"DELETE", key}
		return "", append, nil
	default:
		return "", append, fmt.Errorf("Word Format Invalid! Expected SET, GET, or DELETE")
	}
}

func main() {
	boottime := time.Now()
	store := Create()
	runBackup(store)

	store.CleanUp()
	uptime := time.Since(boottime)
	fmt.Println("Server Boot Time:", uptime)
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Listener Error:", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection", err)
		}
		go HandleConnection(conn, store)
	}
}
