package scu

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type MolAtom struct {
	molid  int
	atname string
}

func (M *MolAtom) Molid() int {
	return M.molid
}

func (M *MolAtom) AtName() string {
	return M.atname
}

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
	ret, err := IndexStringParse(line)
	return ret, err
}

func IndexStringParse(str string) ([]int, error) {
	var err error
	fields := strings.Fields(str)
	ret := make([]int, len(fields))
	for key, val := range fields {
		ret[key], err = strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func MolAtomFileParse(filename string) ([]*MolAtom, error) {
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
	ret, err := MolAtomStringParse(line)
	return ret, err
}

func MolAtomStringParse(str string) ([]*MolAtom, error) {
	var err error
	fields := strings.Fields(str)
	l := len(fields)
	if l%2 != 0 {
		return nil, fmt.Errorf("The string to process must have an even number of fields but it has: %d", l)
	}
	ret := make([]*MolAtom, 0, len(fields)/2)
	var m *MolAtom
	for key, val := range fields {
		if (key+2)%2 == 0 {
			m = new(MolAtom)
			m.molid, err = strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
		} else {
			m.atname = val
			ret = append(ret, m)
		}
	}
	return ret, nil
}

//returns true if test is in container, false otherwise.

func IsInInt(test int, container []int) bool {
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
func IsInString(test string, container []string) bool {
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
// -1 if test is not present in set. Panics if set is not a slice.
//This function was mostly written as a toy. At least for my use cases,
//The two copy/pasted versions above are enough, and reflection doesn't seem justified.
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

//search a file backwards, i.e., starting from the end, for a string. Returns the line that contains the string, or an empty string
func BackwardsSearch(filename, str string) string {
	var ini int64 = 0
	var end int64 = 0
	var first bool
	first = true
	buf := make([]byte, 1)
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()
	var i int64 = 1
	for ; ; i++ {
		if _, err := f.Seek(-1*i, 2); err != nil {
			return ""
		}
		if _, err := f.Read(buf); err != nil {
			return ""
		}
		if buf[0] == byte('\n') && first == false {
			first = true
		} else if buf[0] == byte('\n') && end == 0 {
			end = i
		} else if buf[0] == byte('\n') && ini == 0 {
			i--
			ini = i
			f.Seek(-1*(ini), 2)
			bufF := make([]byte, ini-end)
			f.Read(bufF)
			if strings.Contains(string(bufF), str) {
				return string(bufF)
			}
			//	first=false
			end = 0
			ini = 0
		}

	}
}
