
# How to run: 

You need to have go installed in your machine, and then you can simply run: 

`go run server/main.go` from a terminal and `go run client/main.go` from a different terminal.

if you want to more options, you can run `go run client/main.go --help` 

# Instructions:

There is a simple implementation of file server based on http.FileServer handle ( https://pkg.go.dev/net/http#example-FileServer ).

The server instance is running on top of simple file folder which doesn’t have nested subfolders.


Please implement client which downloads files using this server.

You should download a file containing char 'A' on earlier position than other files.

In case several files have the 'A' char on the same the earliest position you should download all of them.


Each TCP connection is limited by speed. The total bandwidth is unlimited.

You can use any disk space for temporary files.


The goal is to minimize execution time and data size to be transferred.

 

======================================================

Example


If the folder contains the following files on server:


'file1' with contents: "---A---"

'file2' with contents: "--A------"  

'file3' with contents: "------------"

'file4' with contents: "==A=========="

 

then 'file2' and 'file4' should be downloaded


