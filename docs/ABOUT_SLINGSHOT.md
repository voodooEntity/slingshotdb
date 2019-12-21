
# About SlingshotDB
This document should provide a small introduction about how SlingshotDB works.

SlingshotDB is tool for storing graphs and traversing in them. In SlingshotDB Entities are taking the role of nodes and Relations the role of edges.

#### Entity Example    
```go
type Entity struct {
	ID         int
	Type       string
	Value      string
	Context    string
	Properties map[string]string
	Version    int
}
```
* `Type` Entity type 
* `ID` auto increment per Type
* `Value` value for this Entity
* `Context` can be used for search in different methods
* `Properties` a dynamic list of key/value strings that can be used to store any kind of information (from timestamp to description) 
* `Version` used to tackle race condition problems on multithread clients

#### Relation Example
```go
type Relation struct {
	SourceType string
	SourceID   int
	TargetType string
	TargetID   int
	Context    string
	Properties map[string]string
	Version    int
}
```
* `SourceType` Entity type of source entity
* `SourceID` auto increment identifier of source entity
* `TargetType` Entity type of target entity
* `TargetID` auto increment identifier of target entity
* `Context` can be used for search in different methods (future)
* `Properties` a dynamic list of key/value strings that can be used to store any kind of information (from timestamp to description) 
* `Version` used to tackle race condition problems on multithread clients    

While SlingshotDB allows you to create Network structures, due to the nature of JSON as transport format those will be flattened in the output if using the traverse option. To retrieve a Network structure you need to read it step by step (in future there will be a network retrievel method).    

#### A simple example of data would be the Simpsons family tree:    
![Simpsons family tree](http://scriptjungle.de/slingshotdb/simpsons.png)     

------

An other of how your input could look like when mapping a tree structure:
```javascript
{
    "Type": "House",
    "Context": "Rental",
    "Value": "Holmes Home",
    "Properties": {
        "City" : "London"
        "Streetname": "Baker Street",
        "Streetnumber" : "221b"
    },
    "Children": [
        {
            "Type": "Floor",
            "Context": "Rental",
            "Value": "Ground",
            "Properties": {
                "Rooms" : "2"
            },
            "Children": [
                {
                    "Type": "Room",
                    "Context": "Rental",
                    "Value": "Living room",
                    "Properties": {
                        "width" : "10m",
                        "length": "12m",
                        "windows" : "4"
                    },
                    "Children": []
                },
                {
                    "Type": "Room",
                    "Context": "Rental",
                    "Value": "Bathroom",
                    "Properties": {
                        "width" : "14m",
                        "length": "8m",
                        "windows" : "2"
                    },
                    "Children": []
                }
            ]
        },
        {
            "Type": "Floor",
            "Context": "Rental",
            "Value": "First",
            "Properties": {
                "Rooms" : "3"
            },
            "Children": [
                {
                    "Type": "Room",
                    "Context": "Rental",
                    "Value": "Bedroom",
                    "Properties": {
                        "width" : "8m",
                        "length": "3m",
                        "windows" : "3"
                    },
                    "Children": []
                },
                {
                    "Type": "Room",
                    "Context": "Rental",
                    "Value": "Bedroom",
                    "Properties": {
                        "width" : "4m",
                        "length": "4m",
                        "windows" : "3"
                    },
                    "Children": []
                },
                {
                    "Type": "Room",
                    "Context": "Rental",
                    "Value": "Bathroom",
                    "Properties": {
                        "width" : "6m",
                        "length": "3m",
                        "windows" : "1"
                    },
                    "Children": []
                }
            ]
        }
    ]
}

```
In this case House is the most parental entity. It has two child entities, representing the houses floors. Those themself got multiple children, representing the rooms in that floor. The properties are used to hold some information about the datasets. Deciding if an information should be a property or a new child node itself should depend on your usage. When mapped in this structure the system will create relations from entities containing a children array to the entities inside the children array. 

In the current state SlingshotDB can be used via the [HTTP API](https://github.com/voodooEntity/slingshotdb/blob/master/docs/HTTP_API_V1.md). 

While the database is based on in-memory operations, it offers the option to have asynchronous persistency. This will have a rather small impact to the write/update/delete actions. If you enable persistance in config.php , the database will import all persistance datasets on every startup.

The database itself does not ship with a user/permission management (like elastic). This decision was made because i think your security should not rely on all your softwares implementation of such, instead you should use things like 'api-gateways' or smiliar to achieve the security management you want/need.

To tackle clientside race-condition problems the database implements a version number for each dataset. When updating a dataset you need to retrieve the dataset before to know the current version. In case your client is working from multple parallel instances, there is always the chance that two instances try to edit the same dataset. To prevent this the database checks if the version in your update dataset matches the current one in the database. On mismatch the update action will be canceld and the client will be informed with a proper error message.
