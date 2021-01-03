# keyValueDatabase

Simple key-value database engine developed in Go Golang

Avltree directory contains core layer, database is organized as a binary tree to improve efficiency. Binary tree is self-balancing as using AVL algorithm. This layer is responsible for managing CRUD operations and calling the storage layer to commit changes to permanent storage.
Caching.go file contains a caching interface that has three different strategies to manage in-memory caching of values stored in the database, keys are always loaded in memory.

Storage directory is home to a simple layer that manages persistency. One file contains keys organized in the binary tree and another one contains values.

Webserver is a directory that stores HTTP interface to the database. It is used to manage users and do CRUD operations.

Usermanager is a specific implementation of avltree key-value database that stores users data used to authorize requests to web server. Package contains UserManager class and PermissionManager interface. UserManager manages CRUD operations on users and checking credentials. PermissionManager implements different strategies to decide who is allowed to read and write to the database.

Starter creates a database instance based on config file. It constructs the required components and injects them where needed.

Inspiration http://www.aosabook.org/en/500L/dbdb-dog-bed-database.html
