# SlingshotDB 
< some nice text>


## HTTP API v1
**Index**
* /v1
  * [/getEntityByTypeAndId](#route-v1space)
  * [/getEntitiesByTypeAndValue](#route-v1dashboard)
  * [/getEntitiesByType](#route-v1visualisation)
  * [/getEntitiesByValue]()
  * [/getParentEntities]()
  * [/getChildEntities]()
  * [/createEntity]()
  * [/mapJson]()
  * [/updateEntity]()
  * [/deleteEntity]()
  * [/getEntityTypes]()
  * [/getRelation]()
  * [/createRelation]()
  * [/updateRelation]()
  * [/deleteRelation]()
  * [/getRelationsTo]()
  * [/getRelationsFrom]()

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
* 401 Insufficient permissions
* 404 Unknown space id given 
___

