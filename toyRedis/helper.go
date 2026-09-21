package main

import (
	"fmt"
	"os"
	"strings"
)

func CustomSplit(s, sep string) ([]string, error) {
	var result []string

	if sep == "" {
		for _, runeValue := range s {
			result = append(result, string(runeValue))
			return result, nil
		}
	}
	for {
		parenindex := strings.Index(s, "(")
		if parenindex != 0 {
			idx := strings.Index(s, sep)

			if idx == -1 {
				result = append(result, s)
				break
			}

			result = append(result, s[:idx])
			s = s[idx+len(sep):]

		} else {
			idx := strings.Index(s, ")")
			if idx == -1 {
				return make([]string, 0), fmt.Errorf("No Closing Parenthesis")
			}
			result = append(result, s[parenindex+1:idx])
			s = s[idx+len(")"):]
		}

	}
	return result, nil
}

func Append(append []string) error {
	line := strings.Join(append, " ")
	f, err := os.OpenFile(Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
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
