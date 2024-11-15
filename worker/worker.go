package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	line string
	lineNum int
	path string
}

type Results struct {
	Inner []Result
}

// creates a new Result struct
func New(line string,lineNum int,path string) Result{
	return Result{line,lineNum,path}
}

// function to search a file with a key string
func SearchFile(path string, key string) *Results {
	file,err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println("Error : ", err)
		return nil
	}

	results := Results{make([]Result, 0)}

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan(){
		if strings.Contains(scanner.Text(),key){
			r := New(scanner.Text(),lineNum,path)
			results.Inner = append(results.Inner, r)
		}
		lineNum++
	}
	if len(results.Inner) == 0{
		return nil
	}else{
		return &results
	}

}

func ResultTempl(result *Result) string {
	return fmt.Sprintf("%v[%v]:%v",result.path,result.lineNum,result.line)
}