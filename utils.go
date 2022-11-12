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

//these things will go away when generics reach the standard library.

//appends test to containter only if it's not already present
func AppendNRString(test string, container []string) []string {
	if !IsInString(test, container) {
		return append(container, test)
	}
	return container
}

//appends test to containter only if it's not already present
func AppendNRInt(test int, container []int) []int {
	if !IsInInt(test, container) {
		return append(container, test)
	}
	return container
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
//This will fail if the seeked line is the first one. Instead of fixing it, I wrote a more general way of reading a file backwards
//the BWFile structure and it's methods. Use that instead of this.
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

//BWFile is a structure for
//reading a file line by line, but starting from
//the file's end and backwards towards its
//begining.
type BWFile struct {
	fname    string
	f        *os.File
	readable bool
	EOF      bool
	i        int64
}

//Creates a BWFile structure from the name of a file, returns a pointer
//to it and an error.
func NewBWFile(fname string) (*BWFile, error) {
	var err error
	r := new(BWFile)
	r.f, err = os.Open(fname)
	if err != nil {
		return nil, err
	}
	r.readable = true
	r.i = 1
	r.fname = fname
	return r, err

}

//Closes the file.
//when you finished reading it
func (B *BWFile) Close() {
	B.f.Close()
	B.readable = false
	B.EOF = false
}

//Reads the previous line of the file. Returns the read string and an error
//An attempt to read a line after the first one has ben read will return
//an EOF error and close the file. Further attempts, as any attempt to
//read an unopened file will return a "not readable" error.
func (B *BWFile) PrevLine() (string, error) {
	if B.EOF {
		B.Close()
		B.EOF = false
		return "EOF", fmt.Errorf("EOF")
	}
	if !B.readable {
		return "", fmt.Errorf("scu/BWFile: %s is not readable", B.fname)
	}
	var ini int64 = 0
	var end int64 = 0
	buf := make([]byte, 1)
	var last bool

	for ; ; B.i++ {
		//	if last {
		///		return "", fmt.Errorf("EOF")
		//	}
		if _, err := B.f.Seek(-1*B.i, 2); err != nil {
			if !strings.Contains(err.Error(), "invalid argument") {
				return "", err
			}
			last = true
			B.i--
			B.f.Seek(-1*B.i, 2)
			B.i++

		}
		if _, err := B.f.Read(buf); err != nil {
			B.Close()
			return "", err
		}
		if last {
			buf[0] = byte('\n')
		}
		if buf[0] == byte('\n') && end == 0 {
			end = B.i
		} else if buf[0] == byte('\n') && ini == 0 {
			B.i--
			ini = B.i
			B.f.Seek(-1*(ini), 2)
			bufF := make([]byte, ini-end)
			B.f.Read(bufF)
			B.i++
			if last {
				B.EOF = true
			}
			return string(bufF), nil
		}

	}

}

//MustReadFile is a structure to
//somewhat simplify reading text from a file
//Including the option for reading line by line with
//panics on failure, instead of error
type MustReadFile struct {
	fname    string
	f        *os.File
	buf      *bufio.Reader
	readable bool
}

func NewMustReadFile(name string) (*MustReadFile, error) {
	var err error
	r := new(MustReadFile)
	r.f, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	r.buf = bufio.NewReader(r.f)
	r.readable = true
	return r, err

}

//ErrNext does the same as just reading the file with
//as a bufio.Reader, except that in case of EOF it marks
//the file as unreadable
func (F *MustReadFile) ErrNext() (string, error) {
	line, err := F.buf.ReadString('\n')
	if err != nil && err.Error() == "EOF" {
		F.readable = false
	}
	return line, err
}

//reads the next line of a file and returns it. Panics on error, except
//EOF, in which case, it returns the string "EOF". If a further read is
//attempted after an EOF, Next will panic.
func (F *MustReadFile) Next() string {
	if !F.readable {
		panic(fmt.Sprintf("MustReadFile: %s not readable", F.fname))
	}
	line, err := F.buf.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			F.readable = false
			return err.Error()
		} else {
			panic(err.Error())
		}
	}
	return line
}

//Closes the file.
//when you finished reading it
func (F *MustReadFile) Close() {
	F.f.Close()
	F.readable = false
}

//Opens a the file 'name' to append. Creates a new file if
//it doesn't exist
func OpenToAppend(name string) (*os.File, error) {
	return os.OpenFile(name,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

//Parses the int present in the string number, after removing
//leading and trailing spaces (as defined by Unicode, so \n and \t get removed also.
//Panics if it can't parse the float.
func MustAtoi(number string) int {
	num, err := strconv.Atoi(strings.TrimSpace(number))
	if err != nil {
		panic(err.Error())
	}
	return num

}

//Parses the float present in the string number, after removing
//leading and trailing spaces (as defined by Unicode, so \n and \t get removed also.
//Panics if it can't parse the float.
func MustParseFloat(number string, size ...int) float64 {
	s := 64
	if len(size) < 0 && size[0] == 32 {
		s = size[0]
	}
	num, err := strconv.ParseFloat(strings.TrimSpace(number), s)
	if err != nil {
		panic(err.Error())
	}
	return num

}

//QErr takes an error and panics unless the error is nil.
//It is a silly little function to quickly deal with errors
//in small, throw-away programs, not to be used on longer
//programs that are meant to be mantained.
func QErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
