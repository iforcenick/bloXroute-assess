You need to implement a **Client-Server application** with the following requirements:
* multiple-threaded server;
* clients;
* External queue between the clients and server;

Clients:
* [x] Should be configured from a command line or from a file (you decide);
* [x] Can read data from a file or from a command line (you decide);
* [x] Can request server to AddItem(), RemoveItem(), GetItem(), GetAllItems()
* [x] Data is in the form of strings;

* Clients can be added / removed while not intefring to the server or other clients ;

Server:
* [x] Has data structure(s) that holds the data in the memory while keeping the order of items as they added (Ordered Map for C++);
  - The data structure must keep the order of items as they added. 
    For example: If client added the following keys in the following order A, B, D, E, C. 
    The GetAllItems returns A, B, D, E, C
	If item D was removed, the GetAllItems return A, B, E, C
* [x] Server should be able to add an item, remove an item, get a single or all item from the data structure;

External queue:
* [x] Can be Amazon Simple Queue Service (SQS) or RabbitMQ (you decide);


Clients send requests to the external queue - while the server reads those and execute them on its data structure. You define the structure of the messages (AddItem, RemoveItem, GetItem, GetAllItems)


The flow of the project:
1. [x] Multiple clients are sending requests to the queue (and not waiting for the response).
2. [x] Server is reading requests from the queue and processing them, the output of the server is written to a log file
3. [x] Server should be able to process items in parallel
4. [x] log messages (debug, error) are written to stdout

   
Definition of success:
* [x] Working project that can be executed on your computer (preferred OS = linux);
* [x] Being able to explain how the project works and how to deploy the project (for the first time) on another computer;
* [x] If you take something from the Internet or consult anyone for anything, you should be able to understand it perfectly;
* [x] Code has no bugs, no dangling references / assets / resources, no resource leaks;
* [x] Code is clean and readable;
* [x] Code is reasonably efficient (server idle time will be measured).
* [x] Working with channels when needed
* [x] You implement the data structue(s) by yourself


You should develop the project using GOLang.

Good luck!