package main 

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func main() {
	//Create a TCP listener on localhost with porth 27001
	server, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	//Spawn a new goroutine whenever a client connects
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go sendFileToClient(connection)
	}
}

func sendFileToClient(connection net.Conn) {
	fmt.Println("A client has connected!")
	defer connection.Close()
	//Create buffer to read in the name and size of the file
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	bufferFilePath := make([]byte, 60)
	//Get the FilePath
	connection.Read(bufferFilePath)
	filePath := strings.Trim(string(bufferFilePath), ":")
	fmt.Println(filePath)
	//Get the filesize
	connection.Read(bufferFileSize)
	//Strip the ':' from the received size, convert it to a int64
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	//Get the filename
	connection.Read(bufferFileName)
	//Strip the ':' once again but from the received file name now
	fileName := strings.Trim(string(bufferFileName), ":")
	//Create a new file to write in
	if fileName == "" {
		return
	}
	newFile, err := os.Create(filePath + "/" + fileName)
	if err != nil {
		panic(err)
		return
	}
	defer newFile.Close()
	//Create a variable to store in the total amount of data that we received already
	var receivedBytes int64
	//Start writing in the file
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			//Empty the remaining bytes that we don't need from the network buffer
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			//We are done writing the file, break out of the loop
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		//Increment the counter
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
	return
}