package chemu

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

//IndexFileParse will read a file which contains one line with integer numbers separated by spaces. It returns those numbers
//as a slice of ints, and an error or nil.
func IndexFileParse(filename string) ([]int, error) {
	parfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer parfile.Close()
	indexes := bufio.NewReader(parfile)
	line, err := indexes.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(line)
	ret := make([]int, len(fields))
	for key, val := range fields {
		ret[key], err = strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

//returns true if test is in container, false otherwise.

func IsInInt(container []int, test int) bool {
	if container == nil {
		return false
	}
	for _, i := range container {
		if test == i {
			return true
		}
	}
	return false
}

//Same as the previous, but with strings.
func IsInString(container []string, test string) bool {
	if container == nil {
		return false
	}
	for _, i := range container {
		if test == i {
			return true
		}
	}
	return false
}

//IsIn returns the position of test in the slice set, or
// -1 if test is not present in set. Panics if set is not a slice
func IsIn(test interface{}, set interface{}) int {
	vset := reflect.ValueOf(set)
	if reflect.TypeOf(set).Kind().String() != "slice" {
		panic("IsIn function needs a slice as second argument!")
	}
	if vset.Len() < 0 {
		return 1
	}
	for i := 0; i < vset.Len(); i++ {
		vcomp := vset.Index(i)
		comp := vcomp.Interface()
		if reflect.DeepEqual(test, comp) {
			return i
		}
	}
	return -1
}
