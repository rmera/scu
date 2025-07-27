package scu

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// format-panic. Panics with the given format string and its variables.
// Just because you are panicking doesn't mean you can't be methodical.
func Fpanic(s string, a ...any) {
	panic(fmt.Sprintf(s, a...))

}

// Prints the d arguments to stderr
// Meant to be a debugging function that's easy to find and delete
// from programs. Yeah, I use print statements.
func DB(d ...interface{}) {
	fmt.Fprintln(os.Stderr, d...)
}

// If cutoff  is greater or equal to verb, prints the d arguments to stderr
// In either case, returns the argument as a string.
func LogV(cutoff, verb int, d ...interface{}) string {
	if cutoff <= verb {
		fmt.Fprintln(os.Stderr, d...)
	}
	return fmt.Sprintln(d...)

}

// If cutoff  is greater or equal to verb, prints the d arguments to stdout
// Otherwise, does nothing.
func PrintV(cutoff, verb int, d ...interface{}) {
	if cutoff <= verb {
		fmt.Println(d...)
	}

}

// A not efficient function to delete from v all elements with indexes in todel
// The slice todel (i.e. NOT the slice from which elements are being deleted,
// but the one containing the indexes) is altered (sorted in descending order)
func Delete[S ~[]E, E any](v S, todel []int) S {
	slices.Sort(todel)
	slices.Reverse(todel)
	for _, w := range todel {
		v = slices.Delete(v, w, w+1)
	}
	return v
}

// if each element of the slice (except for the first) is 1+the previous element, the slice is contiguous
func IsContiguous(s []int) bool {
	for i, v := range s[1:] {
		if v != s[i-1] {
			return false
		}
	}
	return true
}

// Returns a slice witht the elements v1 and v2 have in common.
func Intersection[S ~[]E, E comparable](v1 S, v2 S) S {
	ret := make([]E, 0, 1)
	for _, v := range v1 {
		if slices.Contains(v2, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Is subset a subset of set?
// 2 identical sets are subsets of each other.
func IsSubset[S ~[]E, E comparable](subset S, set S) bool {
	if len(subset) > len(set) {
		return false //saves some time
	}
	for _, v := range subset {
		if !slices.Contains(set, v) {
			return false
		}
	}
	return true
}

// MolAtom is a simple structure to keep track of goChem atoms, keeping only the molID and atname
type MolAtom struct {
	molid  int
	atname string
}

// Molid returns the molid of the atom
func (M *MolAtom) Molid() int {
	return M.molid
}

// AtName returns the name of the atom
func (M *MolAtom) AtName() string {
	return M.atname
}

// ReplaceInFile replaces the occurences of regex in the file inpfile, into a new file
// outfile. If inpfile and outfile are the same, it creates a temporal file with the
// replacement, which then replaces inpfile by renaming.
func ReplaceInFile(inpfile, outfile, regex, replacement string) error {
	out := outfile
	if outfile == inpfile {
		out = "repl.tmp"
	}
	fin, err := NewMustReadFile(inpfile)
	if err != nil {
		return err
	}
	fout, err := os.Create(out)
	if err != nil {
		return err
	}
	defer func() {
		if fin != nil {
			fin.Close()
		}
		if fout != nil {
			fout.Close()
		}
	}()
	re, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	var i string
	var ferr error
	for i, ferr = fin.ErrNext(); ferr != nil; i, ferr = fin.ErrNext() {
		j := re.ReplaceAllString(i, replacement)
		fout.WriteString(j)
	}
	if ferr.Error() != "EOF" {
		return ferr
	}

	if inpfile == outfile {
		fin.Close()
		fin = nil
		fout.Close()
		fout = nil
		err = os.Rename(out, inpfile)
		if err != nil {
			return err
		}
	}
	return nil
}

// IndexFileParse will read a file which contains one line with integer numbers separated by spaces. It returns those numbers
// as a slice of ints, and an error or nil.
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

// IndexesFileParse will read a file which contains several lines with integer numbers separated by spaces. It returns those numbers
// as a slice of slices of ints, and an error or nil.
func IndexesFileParse(fname string) ([][]int, error) {
	ret := make([][]int, 0, 3)
	f, err := NewMustReadFile(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	for l := f.Next(); l != "EOF"; l = f.Next() {
		s, err := IndexStringParse(strings.Replace(l, "\n", "", -1))
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

// IndexStringParse will read a string that contains integer numbers separated by spaces. It returns those numbers
// as a slice of ints, and an error or nil.
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

// MolAtomFileParse parses a file with a line that contains one or more "atom" info (pairs of molid, atomname)
// only the first line of the file is read. A slice of *MolAtom is returned
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

// MolAtomStringParse Parses a string that contains one or more "atom" info (pairs of molid, atomname)
// only the first line of the file is read. A slice of *MolAtom is returned
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

// BackwardSearch (DEPRECATED) search a file backwards, i.e., starting from the end,
// for a string. Returns the line that contains the string, or an empty
// string. This will fail if the seeked line is the first one. Instead
// of fixing it, I wrote a more general way of reading a file backwards
// the BWFile structure and it's methods. Use that instead of this.
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

// BWFile is a structure for
// reading a file line by line, but starting from
// the file's end and backwards towards its
// begining.
type BWFile struct {
	fname    string
	f        *os.File
	readable bool
	EOF      bool
	i        int64
}

// Creates a BWFile structure from the name of a file, returns a pointer
// to it and an error.
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

// Closes the file.
// when you finished reading it
func (B *BWFile) Close() {
	B.f.Close()
	B.readable = false
	B.EOF = false
}

// Reads the previous line of the file. Returns the read string and an error
// An attempt to read a line after the first one has ben read will return
// an EOF error and close the file. Further attempts, as any attempt to
// read an unopened file will return a "not readable" error.
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

// MustReadFile is a structure to
// somewhat simplify reading text from a file
// Including the option for reading line by line with
// panics on failure, instead of error
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

// ErrNext does the same as just reading the file with
// as a bufio.Reader, except that in case of EOF it marks
// the file as unreadable
func (F *MustReadFile) ErrNext() (string, error) {
	line, err := F.buf.ReadString('\n')
	if err != nil && err.Error() == "EOF" {
		F.readable = false
	}
	return line, err
}

// reads the next line of a file and returns it. Panics on error, except
// EOF, in which case, it returns the string "EOF". If a further read is
// attempted after an EOF, Next will panic.
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

// Closes the file.
// when you finished reading it
func (F *MustReadFile) Close() {
	F.f.Close()
	F.readable = false
}

// Opens a the file 'name' to append. Creates a new file if
// it doesn't exist
func OpenToAppend(name string) (*os.File, error) {
	return os.OpenFile(name,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// Parses the int present in the string number, after removing
// leading and trailing spaces (as defined by Unicode, so \n and \t get removed also.
// Panics if it can't parse the float.
func MustAtoi(number string) int {
	num, err := strconv.Atoi(strings.TrimSpace(number))
	if err != nil {
		panic(err.Error())
	}
	return num

}

// Parses the float present in the string number, after removing
// leading and trailing spaces (as defined by Unicode, so \n and \t get removed also.
// Panics if it can't parse the float.
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

// QErr takes an error and panics unless the error is nil.
// It is a silly little function to quickly deal with errors
// in small, throw-away programs, not to be used on longer
// programs that are meant to be mantained.
func QErr(err error, additionalmsg ...string) {
	if err != nil {
		str := ""
		if len(additionalmsg) > 0 {
			str = strings.Join(additionalmsg, " - ")
		}
		panic(str + " " + err.Error())
	}
}

//Legacy things from before generics

// appends test to containter only if it's not already present
func AppendNRString(test string, container []string) []string {
	if !IsInString(test, container) {
		return append(container, test)
	}
	return container
}

// appends test to containter only if it's not already present
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

// Same as the previous, but with strings.
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

// IsIn returns the position of test in the slice set, or
// -1 if test is not present in set. Panics if set is not a slice.
// This function was mostly written as a toy. At least for my use cases,
// The two copy/pasted versions above are enough, and reflection doesn't seem justified.
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
