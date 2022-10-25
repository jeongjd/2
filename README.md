# 2

## Instructions

## (NEED TO CHANGE name of REPO!!!) 
In your terminal, 

    git clone git@github.com:jeongjd/2.git 
   
### To Create The Server:

    cd 2/server
    go run server.go 

You will be prompted to enter a port number. This will create a TCP server that listens to that port for connections. 

    $ Enter a port number: 

### To Create The Client:

Open up a new terminal. 

    cd 2/client
    go run client.go 
    
You will be prompted to enter a host address, port number, and username. 
To create a client with a connection to a server, Enter the SAME port number used to create the server. 

    $ Enter a host address:
    $ Enter a port number: 
    $ Enter your username in the format /name: 

For every new client, open up a new terminal to connect to the server. 

### To Send A Private Message
In the client window, Enter the command in the following format: {To} {From} {Message}

For example, 

    $ Bob Alex hello 



