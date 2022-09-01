/*
Copyright 2022 k0s authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"strconv"
	"errors"
	"os/exec"
	"sort"
	"reflect"

)

// Exists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func DirExist(DirName string) bool {
	_, err := os.Stat(DirName)
	 return  !os.IsNotExist(err)  
}


// Copy copies file from src to dst
func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error reading source file (%v): %v", src, err)
	}

	err = os.WriteFile(dst, input, sourceFileStat.Mode())
	if err != nil {
		return fmt.Errorf("error writing destination file (%v): %v", dst, err)
	}
	return nil
}

func WriteTmpFile(data string, prefix string) (path string, err error) {
	tmpFile, err := os.CreateTemp("", prefix)
	if err != nil {
		return "", fmt.Errorf("cannot create temporary file: %v", err)
	}

	text := []byte(data)
	if _, err = tmpFile.Write(text); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}

	return tmpFile.Name(), nil
}


// GetUID returns uid of given username and logs a warning if its missing
func GetUID(name string) (int, error) {
	entry, err := user.Lookup(name)
	if err == nil {
		return strconv.Atoi(entry.Uid)
	}
	if errors.Is(err, user.UnknownUserError(name)) {
		// fallback to call external `id` in case NSS is used
		out, err := exec.Command("/usr/bin/id", "-u", name).CombinedOutput()
		if err == nil {
			return strconv.Atoi(strings.TrimSuffix(string(out), "\n"))
		}
	}
	return 0, err
}


func Contains(strSlice []string, str string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}

	return false
}

// IsEqual returns true if an array of strings is equal, regardless of order
func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	if len(a1) == len(a2) {
		return reflect.DeepEqual(a1, a2)
	}
	return false
}

// Unique returns only the unique items from given input slice
func Unique(input []string) []string {
	m := make(map[string]bool)
	result := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := m[s]; !ok {
			m[s] = true
			result = append(result, s)
		}
	}
	return result
}
