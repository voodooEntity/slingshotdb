
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
      "Context": "Example",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
      "Children": [
        {
          "ID": 1,
          "Type": "Port",
          "Context": "Example",
          "Value": "80",
          "Properties": {
            "created_at": "sometimestamp"
          },
          "Version" : 1,
          "Children": []
        },
        {
          "ID": 2,
          "Type": "Port",
          "Context": "Example",
          "Value": "443",
          "Properties": {
            "created_at": "sometimestamp"
          },
          "Version" : 1,
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
      "Version" : 1,
      "Properties": null
    },
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 2,
      "Context": "",
      "Version" : 1,
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

**Optional parameters:**       
*  `GET` mode string   (compare mode used with `GET` value string)
    * `match` (exact match, default mode)
    * `prefix` (value must begin with)
    * `suffix` (value must end with)
    * `contain` (value must contain)
    * `regex` (value must match regex pattern)
* `GET` context string (can be used as filter, requires exact match)    


**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Example",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
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

**Optional parameters:**      
*  `GET` mode string   (compare mode used with `GET` value string)
    * `match` (exact match, default mode)
    * `prefix` (value must begin with)
    * `suffix` (value must end with)
    * `contain` (value must contain)
    * `regex` (value must match regex pattern)
* `GET` context string (can be used as filter, requires exact match)    

**Body:**  none        
**Return:**    
```javascript
{
  "Entities": [
    {
      "ID": 1,
      "Type": "IP",
      "Context": "Example",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
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
      "Context": "Example",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
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
      "Context": "Example",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
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
      "Context": "Example",
      "Value": "80",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
      "Children": []
    },
    {
      "ID": 2,
      "Type": "Port",
      "Context": "Example",
      "Value": "443",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
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
		"the" : "mentor"
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
      "Version" : 1,
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
    "Context" : "Example",
    "Children" : [
        {
            "Type" : "Port",
            "Value" : "80",
            "Context" : "Example",
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
            "Context" : "Example",
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
      "Context": "Example",
      "Value": "127.0.0.1",
      "Properties": {
        "created_at": "sometimestamp"
      },
      "Version" : 1,
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
      "Version" : 1,
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
      "Version" : 1,
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
  "Version" : 1,
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
      "Version" : 1,
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
      "Version" : 1,
      "Properties": null
    },
    {
      "SourceType": "IP",
      "SourceID": 1,
      "TargetType": "Port",
      "TargetID": 2,
      "Context": "",
      "Version" : 1,
      "Properties": null
    }
  ]
}
```
**Errors:**  
* TODO
___
