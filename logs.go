package logs

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
	"os"
	"io/ioutil"
	"strconv"
)

//you can configurat logs package by the following values
const(
	// is the path saving the logs files
	logs_root= `./logs/logsfile/`
	// you can set mew_log = 0 when you are testing, then it will cover the oldest flode to save the logsfiles 
	//note that new_log should be 0 or 1 !
	create_floder = true
)

//name of logs files
var(
	error_log_name = "error.log"
	warn_log_name = "warning.log"
	info_log_name = "info.log"
)
//pointer of log files
var(
 	errlog  *log.Logger
	warnlog  *log.Logger
	infolog  *log.Logger
)


func init() {
	log_floder := createfloder()
	errorp, err := os.Create(log_floder + error_log_name)
	warnp, _ := os.Create(log_floder + warn_log_name)
	infop, _ := os.Create(log_floder + info_log_name)
	if err != nil  {
		panic("logs: can not create logsfile!")
	}
	errlog = log.New(errorp, "", 3)
	warnlog = log.New(warnp, "", 3)
	infolog = log.New(infop, "", 3)
}

//=======================================================================================
//===================== the following is  public value or function ======================


//the following constant are the paramater used in Log() style
const (
	Err = 400
	Warn = 200
	Q_err =-400
	Q_warn =-200	
)

//record message into two dirfferent log file with different ways
func Log(style int, msg ...interface{}){
	path := getpath()
	path = formatPath(path)
	s := formatInterface(path, msg...)
	if style > 0 {
		fmt.Println(s)
	}else{
		style = 0-style
	} 
	switch style {
	case Err :
		errlog.Println(s)
	case Warn:
		warnlog.Println(s)
	}
}

//print the message both on std output and logs file
func Println(any ...interface{}) {
	path := getpath()
	path = formatPath(path)
	s := formatInterface(path , any...)
	fmt.Println(s)
	infolog.Println(s)
}

//When serious error happen, record the error and save in error.log
func Fatal(msg ...interface{}){
	path := getpath()
	path = formatPath(path)
	s := formatInterface(path, msg...)
	errlog.Println("Fatal !!!")
	errlog.Println(s)
	fmt.Println(s)
	os.Exit(1)
}

//=======================================================================================
//===================== the following is tools function ================================= 

//return the user file location, like "main.go"
func getpath()string{
	//Log()/Println() --> getpath() -->  Caller()
	pc,file,_,ok := runtime.Caller(2)
	if ok==false{
		return "???"
	}
	//get caller code file name
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	short2 := ""
	//get caller function name
	fun := runtime.FuncForPC(pc)
	pcname := fun.Name()
	for i:=len(pcname)-1;i>0;i--{
		if pcname[i] == '.' {
			short2 = pcname[i+1:]
			break;
		}
	}
	return short + " -> "+short2 +"()"
}

//set the format of pathstring
func formatPath(path string) string{
	return  "[" + path + "] : "
}

//the raw ouput is like :[a,b,c]
//and after it function handle, it will be : a  b  c
func formatInterface(prefix string, any ...interface{})string{
	formatstr := strings.Repeat("%v  ", len(any))
	prefix += fmt.Sprintf(formatstr, any...)
	return prefix
}


//create a new floder to save the logs file, return the path of new floder
//note that you can not use Println and Log function in createfloder()
func createfloder()string{
	//1-> check logs_root if exist, create new directory if not exist
	rd, err := ioutil.ReadDir(logs_root)
	if err != nil {
		fmt.Println("logs: Can not read logs_root !", err)
		fmt.Println("logs: Trying to create make directory ：", logs_root)
		err = os.Mkdir(logs_root, os.ModeDir)
		if err!=nil {	//read and create floder fall 
			pwd , _ := os.Getwd()
			fmt.Println("logs: Can not create logs_root!",err)
			fmt.Println("提示：程序当前运行的路径为 :",pwd)
			fmt.Println("提示：如果发现以上目录与代码目录不一致，请用go run main.go 来启动程序")
			os.Exit(1)
		}		
		//can not read and already make an directory
		fmt.Println("logs: Create logs_roots scuess! :",logs_root)
		rd , err = ioutil.ReadDir(logs_root)
		if err!=nil {
			fmt.Println("logs: An Error happen, create logs_roots but can't read it ,",err)
			os.Exit(1)
		}
	}
	tempPath := logs_root //the path to save the logs files
	ps := getPathSeperator()
	if create_floder == false {	//2->if it is testing, all logs_files save in a same floder
		tempPath = logs_root + ps + "2019-05-12#1" +ps
	}else{ 	//3-> in classic model, create new floder for each running
		filenum := 1
		datestr := time.Now().Format("2006-01-02")
		for _, fi := range rd {
			if fi.IsDir() && len(fi.Name())>10 {
				fd := fi.Name()[0:10]
				if fd==datestr{
					filenum++
				}
			}
		}
		tempPath = logs_root + ps +datestr + `#` + strconv.Itoa(filenum) + ps
	}
	err = os.Mkdir(tempPath, os.ModeDir)
	if err!=nil && os. IsExist(err)==false {	//if the reason of mkdir err is because the directory already exist
		fmt.Println("logs : Can not make directory in logs_root!",err)
		os.Exit(1)
	}
	return tempPath
}

//get system pathseperator character
func getPathSeperator() string {
	ps := `\`
	if os.PathSeparator==47 {	//PathSeparator of system
		ps = `/`
	}
	return ps
}