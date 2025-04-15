# Service Dependency API
_This is a work in progress_


This api is designed allow mapping dependencies between services, as well as provide 
basic information about the services such as

- name
- database dependencies
- owner
- github repo

Alternatively, you will be able to associate a GitHub Release with a service


## Neo4j Schenanigans
To make everything unit testable, `neo4j_interfaces` and `neo4j_wrappers` were implemented.
The interfaces mimic the basic required methods from the neo4j library. To get around a types error,
I had to add wrapper classes around
- Driver
- Session
- Transaction
- Result
- Record

Each of these structs take in its counterpart neo4j struct and store it as a private property.
Then it mimics the required methods to implement the appropriate interface. 

> [!WARNING]
> There has to be an easier way to do all this, it seems excessive