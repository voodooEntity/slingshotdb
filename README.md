# SlingshotDB 
Welcome to the home of SlingshotDB - an in-memory entity/relation database. It's completly written in golang (vanilla, no 3rd party libraries used) and provides acces via an HTTP API. 

While the database is based on in-memory operations, it offers the option to have asynchronous persistency. This will have a rather small impact to the write/update/delete actions. If you enable persistance, the database will import all persistance datasets on startup.

In it's current state the database is shipped with minimal functionality. I will extend the functionality over time based on needs and time. 

The main target of the database is to have an easy to use high performance storage.

`Important!` Since the database is in memory and focused on high performance it uses a lot of memory for indexes and data storage. Using the database you should make sure to run it on a machine that provides a lot of memory to work fine. While there are several ways how i could reducde the memory usage i accepted this trade off for the performance.

The database itself does not ship with a user/permission management (like elastic). This decision was made because i think your security should not rely on all your softwares implementation of such, instead you should use things like 'api-gateways' or smiliar to achieve the security management you want/need.

Since this is an early alpha you should expect to encounter bugs that i didn't takle yet. Please feel free to report those as issues, and i will do my best to fix/solve them asap.

If you have suggestions on how to change the database, you can also feel free to push them as issue, but keep in mind that i coded this database for specific purposes. So even if your suggestions may seem legit, they still could be refused if they oppose the way this database is meant to work.

Finally i wanne leave a special thanks to some friends that helped through the process of creating this software by listening to hours of rage/ideas and providing suggestions that lead the way to the software you are about to use. 
* Maze (the name 'Slingshot' was his idea)
* Luxer 
* f0o

I hope you will enjoy the usage of SlingshotDB and that this will just be the start of a great project.

Sincerely yours,
voodooEntity

---
## Build
SlingshotDB needs golang to be build. Just add the directory path you cloned the database into to your gopath and fire 'go build -o slingshot' inside. No special build flags or installing of dependencies needed. After building, just start the database with "./slingshot". 


## Config
The configuration is shipped with the repo in 'config.json' file. In the current state you got three options to configure.
* `host` string  (by setting the host you limit the database to listen to a specific ip. by leaving it empty string it will listen to all IPs )
* `port` int  (the port the database will listen on)
* `persistance` bool (if set on true it will persist your datasets on your harddisk)

Example file:
```javascript
{
    "host" : "",
    "port" : 8090,
    "persistance" : false
}
```

## Usage
After building and starting the database, you are free to use it. You can create your own code to use the API based on the HTTP API v1 docs below, or if your environment is PHP you can use the 'voodooEntity/slingshotdb-php-sdk'. As time goes by i plan to add more SDK's (golang, javascript, ....)

## Future plans
I got several plans on extending the functionality of the database, but i don't have a specific order for implementing those. I plan to release a full list of incoming features in a future release of SlingshotDB. So stay tuned .)


## HTTP API v1
---
### Route: /v1

**Path:** /getEntityByTypeAndId     
**Type:** GET      
**Required Parameters:**    
* `GET` id int    
* `GET` type string    

**Optional parameters:**     
* `GET` traverse int    

**Body:** none    
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Pentest",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": [
        {
          "ID": 1,
          "Type": "Port",
          "Context": "Pentest",
          "Value": "80",
          "Properties": {
            "created_at": "sometimestamp"
          },
          "Children": []
        },
        {
          "ID": 2,
          "Type": "Port",
          "Context": "Pentest",
          "Value": "443",
          "Properties": {
            "created_at": "sometimestamp"
          },
          "Children": []
        }
      ]
    }
  ],
  "Relations": [
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 1,
      "Context": "",
      "Properties": null
    },
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 2,
      "Context": "",
      "Properties": null
    }
  ]
}
```
**Errors:**     
* TODO    
___ 
**Path:** /getEntitiesByTypeAndValue     
**Type:** GET      
**Required Parameters:**      
* `GET` value string    
* `GET` type string    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Pentest",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**  
* TODO
___
**Path:** /getEntitiesByType     
**Type:** GET      
**Required Parameters:**      
* `GET` type string    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Pentest",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**      
* TODO    
___
**Path:** /getEntitiesByValue    
**Type:** GET      
**Required Parameters:**      
* `GET` value string    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Pentest",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**      
* TODO    
___
**Path:** /getParentEntities     
**Type:** GET      
**Required Parameters:**      
* `GET` type string    
* `GET` id int    

**Optional parameters:** none      
**Body:**  none    
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Pentest",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**  
* TODO
___
**Path:** /getChildEntities     
**Type:** GET      
**Required Parameters:**      
* `GET` type string    
* `GET` id int    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "Port",
      "Context": "Pentest",
      "Value": "80",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    },
    {
      "ID": 2,
      "Type": "Port",
      "Context": "Pentest",
      "Value": "443",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**  
* TODO
___
**Path:** /createEntity     
**Type:** POST     
**Required Parameters:** none     
**Optional parameters:** none     
**Body:**      
```javascript
{
	"Type" : "Person",
	"Properties" : {
		"the" : "menthor"
	},
	"Value" : "This is it… this is where I belong…",
	"Context" : "manifesto"
}
```
**Return:**     
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "Person",
      "Context": "manifesto",
      "Value": "This is it… this is where I belong…",
      "Properties": {
        "the": "menthor"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**      
* TODO    
___
**Path:** /mapJson     
**Type:** POST     
**Required Parameters:** none     
**Optional parameters:** none     
**Body:**      
```javascript
{
    "Type" : "IP",
    "Value" : "127.0.0.1",
    "Properties" : {
        "created_at" : "sometimestamp"
    },
    "Context" : "Pentest",
    "Children" : [
        {
            "Type" : "Port",
            "Value" : "80",
            "Context" : "Pentest",
            "Properties" : {
                "created_at" : "sometimestamp"
            },
            "Children" : [
                {
                    "Type" : "State",
                    "Value" : "Open",
                    "Properties" : {
                        "created_at" : "sometimestamp"
                    }
                }
            ]
        },
        {
            "Type" : "Port",
            "Value" : "443",
            "Context" : "Pentest",
            "Properties" : {
                "created_at" : "sometimestamp"
            },
            "Children" : [
                {
                    "Type" : "State",
                    "Value" : "Closed",
                    "Properties" : {
                        "created_at" : "sometimestamp"
                    }
                }
            ]
        }
    ]
}
```
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Pentest",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Children": []
    }
  ],
  "Relations": []
}
```
**Errors:**      
* TODO    
___
**Path:** /updateEntity     
**Type:** PUT     
**Required Parameters:** none     
**Optional parameters:** none     
**Body:**     
```javascript
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Other Context",
      "Value": "127.0.0.3",
      "Properties": {
        "created_not": "no timestamp"
      },
      "Children": null
    }
```
**Return:** none      
**Errors:**      
* TODO    
___
**Path:** /deleteEntity     
**Type:** DELETE     
**Required Parameters:**    
* `GET` type string      
* `GET` id int      

**Optional parameters:** none     
**Body:**  none      
**Return:** none      
**Errors:**      
* TODO
___

**Path:** /getEntityTypes     
**Type:** GET      
**Required Parameters:** none      
**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
[
  "IP",
  "Port",
  "State"
]
```
**Errors:**  
* TODO
___ 
**Path:** /getRelation     
**Type:** GET      
**Required Parameters:**      
* `GET` srcType string    
* `GET` srcID int    
* `GET` targetType string    
* `GET` targetID int    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [],
  "Relations": [
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 1,
      "Context": "",
      "Properties": null
    }
  ]
}
```
**Errors:**  
* TODO
___
**Path:** /createRelation     
**Type:** POST     
**Required Parameters:** none     
**Optional parameters:** none     
**Body:**      
```javascript
{
	"SourceType" : "IP",
	"SourceID" : 1,
	"TargetType" : "State",
	"TargetID" : 1,
	"Context" : "created",
	"Properties" : {
		"created" : "property"
	}
}
```
**Return:** none     
**Errors:**      
* TODO    
___
**Path:** /updateRelation     
**Type:** PUT     
**Required Parameters:** none     
**Optional parameters:** none     
**Body:**     
```javascript
{
	"SourceType" : "IP",
	"SourceID" : 1,
	"TargetType" : "Port",
	"TargetID" : 1,
	"Context" : "updated",
	"Properties" : {
		"new" : "property"
	}
}
```
**Return:** none      
**Errors:**      
* TODO    
___
**Path:** /deleteRelation     
**Type:** DELETE      
**Required Parameters:**      
* `GET` srcType string    
* `GET` srcID int    
* `GET` targetType string    
* `GET` targetID int    

**Optional parameters:** none      
**Body:**  none        
**Return:** none    
**Errors:**  
* TODO
___
**Path:** /getRelationsTo     
**Type:** GET      
**Required Parameters:**      
* `GET` type string    
* `GET` id int    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [],
  "Relations": [
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 2,
      "Context": "",
      "Properties": null
    }
  ]
}
```
**Errors:**  
* TODO
___
**Path:** /getRelationsFrom     
**Type:** GET      
**Required Parameters:**      
* `GET` type string    
* `GET` id int    

**Optional parameters:** none      
**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [],
  "Relations": [
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 1,
      "Context": "",
      "Properties": null
    },
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 2,
      "Context": "",
      "Properties": null
    }
  ]
}
```
**Errors:**  
* TODO
___