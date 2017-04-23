package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
//	"strings"
)

//Define that the binairy data of the file will be sent 1024 bytes at a time
const BUFFERSIZE = 1024

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("./senfFile.go <fileSourceLocation> <fileDestLocation>")
		os.Exit(1)
	}

	fileSourceLocation := os.Args[1]
	fileDestLocation := os.Args[2]

	//Open the file that needs to be send to the client
	file, err := os.Open(fileSourceLocation)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//Make the connection
	connection, err := net.Dial("tcp", "localhost:27001")
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to server, start receiving the file name and file size")
	fmt.Println("A client has connected!")	

	//Get the filename and filesize
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	filePath := fillString(fileDestLocation,60)
	//Send the file header first so the client knows the filename and how long it has to read the incomming file
	fmt.Println("Sending filename and filesize!")
	fmt.Println(fileSize)
	fmt.Println(fileName)
	fmt.Println(filePath)
	//Write first 10 bytes to client telling them the filesize
	connection.Write([]byte(filePath))
	connection.Write([]byte(fileSize))
	//Write 64 bytes to client containing the filename
	connection.Write([]byte(fileName))
		//Initialize a buffer for reading parts of the file in
	sendBuffer := make([]byte, BUFFERSIZE)
	//Start sending the file to the client
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			//End of file reached, break out of for loop
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	defer connection.Close()
}

//This function is to 'fill'
func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

