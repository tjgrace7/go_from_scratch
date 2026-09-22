package main

import (
	"fmt"
	"os"
	"strings"
)

var Path = "append.txt"

// Split is Designed to Seperate Strings at a sep & Keep all strings within () together even with a space.
func CustomSplit(s, sep string) ([]string, error) {
	var result []string
	//If the sep is an empty string. Append each char as a seperate string
	if sep == "" {
		for _, runeValue := range s {
			result = append(result, string(runeValue))
			return result, nil
		}
	}
	//Iterates through the rest of s
	for {
		//finds an open (. If it is not at the zero index, runs standard seperation logic
		parenindex := strings.Index(s, "(")
		if parenindex != 0 {
			//finds next seperator
			idx := strings.Index(s, sep)
			//If it doesn't exist, appends the rest of s
			if idx == -1 {
				result = append(result, s)
				break
			}
			//Appends everything in s until the seperator
			result = append(result, s[:idx])
			//adjusts the slice to remove everything until the seperator (Including the seperator
			s = s[idx+len(sep):]

		} else {
			//IF ( is at the zero index. Searches until it finds ) then appends everything
			idx := strings.Index(s, ")")
			//If ) is not found returns error
			if idx == -1 {
				return make([]string, 0), fmt.Errorf("No Closing Parenthesis")
			}
			result = append(result, s[parenindex+1:idx])
			s = s[idx+len(")"):]
		}

	}
	return result, nil
}

// Appends the command to the append.txt file
func Append(append []string) error {
	//creates a single string from []string
	line := strings.Join(append, " ")
	//Opens the file
	f, err := os.OpenFile(Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	//Writes line to the file
	if _, err := f.Write([]byte(line + "\n")); err != nil {
		f.Close()
		fmt.Println(err)
		return err
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
