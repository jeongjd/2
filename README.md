# 2

This Go Program simulates a TCP chat. Clients can send private encoded messages to each other via TCP chat server. 

## Instructions

In your terminal, git clone the repository 

    git clone git@github.com:jeongjd/2.git 
   
### To Create The Server:

    cd 2/server
    go run server.go 

You will be prompted to enter a port number. This will create a TCP server that listens to that port for connections. 

    $ Enter a port number: 

### To Create The Client:

Open up a new terminal. cd into the directory that contains the cloned repo 2. 

    cd 2/client
    go run client.go 
    
You will be prompted to enter a host address, port number, and username. 
To create a client with a connection to a server, enter the SAME port number used to create the server. 

    $ Enter a host address:
    $ Enter a port number: 
    $ Enter your username in the format /name: 

For every new client, open up a new terminal to connect to the server. 

### To Send A Private Message
In the client window, enter the command in the following format: {To} {From} {Message}

For example, 

    $ Bob Alex hello 

### To Exit The Client Connection

Type "EXIT" (upper case only!) in the client side terminal 

    $ EXIT
    $ Exiting the client... 

### To Shut Down The Server 

Type "EXIT" (upper case only) in the server side terminal 

    $ EXIT
    $ Server is shutting down... 
    

## Design Choices

Gob Serialization: 

Instead of using the write method from the net package, the program uses gob serialization. The client encodes the input command and sends it to the server which decodes it, parses it, encodes it and then sends it to the recipient client. The recipient client decodes the message from the server and prints it.. Gob serialization was chosen because it is extremely efficient and fast compared to other serialization methods such as JSON. Since this program and its communication is done entirely in Go, it is safe to use gob serialization.

Mutex: 

Instead of using a channel to communicate between the go threads, a global variable was used. A hashtable keeps track of the client username and connection. The hashtable is used by the go threads to cache the usernames and connections. That way each client can send a message to any client that joins since the hash table is a global variable. However since the hashtable is being mutated by multiple threads, a read and write mutex was used to synchronize access, that way whenever a go thread updates the hashtable by either adding or removing a connection another go threads cannot access it until it is done updating. This makes the server safe and prevents race problems. 

Goto statement: 

It is used for error handling. If there is an error, it will skip over to the LAST: to invoke broadcastErrorMessage function. 


## The Code Flow 

The program creates a TCP server, which constantly listens for a new client connection and accepts new client connections with the same port number. When a client connection is established, the program executes a go routine handleConnections. 

The server constantly decodes and parses the messages from the client, stores the variables and encodes the message to send it to a correct recipient. 

When the client writes the username and sends it to the server, the server decodes it and stores the client username as a key in the clientConnetions map.  If the username is already taken, it will print an error message to the client and ask for another input. Otherwise, it will restart the loop and wait for clients to send a private message to the server. 

When the client writes a message, it is encoded and sent to the server, which decodes and parses the message. It stores the message in a struct with three string variables: receiverID, senderID, and messageContent. If the message is not in the correct format of {To} {From} {Message}, it will print an error statement to the client. The checkClients function checks if both receiverID and senderID exist in the clientConnections map and returns true. If they both exist, the server sends the encoded message to the correct recipient. Otherwise, it will determine which error it is to set an error option number, then return false. 

The broadcastErrorMessage function prints an error statement according to the error option number. 

If at any point a client closes its connection, the handleConnections function will delete the client username from the map clientConnections and print that the user disconnected from the server. 




