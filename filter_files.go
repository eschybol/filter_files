package main
import (
	"os"
	"fmt"
	"regexp"
	"log"
	"errors"
	"flag"
)

// custom datatypes
type File struct {
	Title string
	Content []byte
}

// flag vars -> here are your commandline values my friend ;)
var searchStringFlag string
var searchRegexPatternFlag  string
var fileNameFlag string

func initFlagVars() {
	defaultValue := ""  
	// search string command loine flags. (-s <string>, --search_string <string>)
	flag.StringVar(&searchStringFlag, "s", defaultValue, "Enter a Search Phrase (Shorthand)")
	flag.StringVar(&searchStringFlag, "search_string", defaultValue, "Enter a Search Phrase")

	flag.StringVar(&searchRegexPatternFlag, "r", defaultValue,"Enter a Regular Expression Pattern (shorthand)")
	flag.StringVar(&searchRegexPatternFlag, "regular_expression", defaultValue, "Enter a Regular Expression Pattern")
	
	flag.StringVar(&fileNameFlag, "f", defaultValue, "Filename (shorthand)")
	flag.StringVar(&fileNameFlag, "file_name", defaultValue, "Filename")


	flag.Parse()
}

func xor(a, b string) bool {
	if ((a == "" || b == "") && !(a == "" && b == "")) {
		return true
	}
	return false

}

func checkRequiredArgs() bool {
	if (fileNameFlag != "") && xor(searchStringFlag, searchRegexPatternFlag) {
		return true
	}
	return false
}


// get file
func openFile(filename string) (*File, error) {
	// defer ?
	content, err := os.ReadFile(filename) 
	if err != nil {
		return nil, err
	}
	return &File{Title: filename, Content:[]byte(content)}, nil
}


/* this should return somethings depending on the filter method ie. 
if regex is used, some parsed regex string, 
or if simple string pattern matching the string*/
func (f *File) filter() ([][]byte, error)  {
	var checkPhrase string
	// check which filter method is used by referencing the command line arguments
	if searchStringFlag != "" && searchRegexPatternFlag == "" {
		checkPhrase = searchStringFlag
	} else if searchStringFlag == "" && searchRegexPatternFlag != "" {
		checkPhrase = searchRegexPatternFlag
	} else {
		log.Fatal(errors.New("Input: Provide Command Line Arguments"))
	}
	// sanitzation?`!
	re := regexp.MustCompile(checkPhrase)
	content := re.FindAll(f.Content, -1)
	if content == nil {
		return [][]byte(nil), fmt.Errorf("Filter: Expression (%s) not Found", checkPhrase)
	}
	return content, nil
}



func main() {
	// runtime vars 
	initFlagVars()
	if !checkRequiredArgs() {
		log.Fatal(errors.New("Input: Provide a Filename and eather a search string or a regex string"))
	} 

	file, err := openFile(fileNameFlag)
	if err != nil {
		log.Fatal(err)
	}

	content, err := file.filter()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q\n" ,content) 
}