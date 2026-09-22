package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

// Runs append.txt backup storage
func runBackup(store *Store) error {
	//Opens File
	f, err := os.Open(Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	//Scans File
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		//Splits Command
		command, err := CustomSplit(scanner.Text(), " ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(command)
		//Runs Each Method. It does not append methods to append.txt. This would create an infinite loop
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
	//Scans commands
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		//Splits Command using Helper.CustomerSplit
		command, err := CustomSplit(scanner.Text(), " ")
		if err != nil {
			conn.Write([]byte(err.Error()))
		}
		//Gets the Value, Append String and Error from RunMethod
		value, append, err := runmethod(command, store)
		if err != nil {
			conn.Write([]byte(err.Error()))
		}
		//Appends the string to append.txt
		Append(append)
		//Writes response to Client
		conn.Write([]byte(value + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid Input: %s", err)
	}

}

func runmethod(command []string, store *Store) (string, []string, error) {
	//Takes command and determines which method to run
	append := []string{}
	//If the command has no data returns error
	if len(command) == 0 {
		return "", append, fmt.Errorf("Command Misformatted")
	}
	//Puts the Method at zero index
	method := command[0]
	command = command[1:]
	key := command[0]
	//Determines which method to use
	switch method {
	case "SET":
		//runs Set
		v, err := store.Set(command)
		if err != nil {
			return "", append, err
		}
		//Creates Append Stribg
		if v.t != "" {
			append = []string{"SET", key, command[1], "EX", "(" + v.t + ")"}
		} else {
			append = []string{"SET", key, command[1]}
		}
		if err != nil {
			fmt.Println(err)
			return "", append, err
		}
		//Returns Complete with Append String.
		return "SET Command Complete", append, nil
	case "GET":
		//Runs Get
		value, err := store.Get(key)
		if err != nil {
			return "", append, err
		}
		//Creates Append String
		append = []string{"GET", key}
		return key + ":" + value, append, nil
	case "DELETE":
		//Runs Delete
		store.Delete(key)
		//Creates Append String
		append = []string{"DELETE", key}
		return "", append, nil
	default:
		return "", append, fmt.Errorf("Word Format Invalid! Expected SET, GET, or DELETE")
	}
}

// Main Function
func main() {
	//Starts Boot Time
	boottime := time.Now()
	store := Create()
	//Runs every command in append.txt
	runBackup(store)
	//Cleans up expired keys from append.txt
	store.CleanUp()
	uptime := time.Since(boottime)
	//Determines how long server boot up took
	fmt.Println("Server Boot Time:", uptime)
	//Starts Server
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Listener Error:", err)
		return
	}
	defer listener.Close()
	//Handles Concurrent Connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection", err)
		}
		go HandleConnection(conn, store)
	}
}
