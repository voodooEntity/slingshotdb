# About SlingshotDB
This document should provide a small introduction about how SlingshotDB works.

Slingshot stores its data in two kinds of datasets, entities and relations. Entities represent data nodes wich have a certrain type (entity type) and a value. Relations are links between those data nodes. Relations have source and target entitys wich defines if the relation is a parent or child relation. 

Both types of data (entities and relations) can hold a dynamic list of properties (key/value pairs) and a context. Those can be used in multiple ways. From storing a timestamp over storing an order id to holding a list of specific properties or many more different options. 

There is no `correct` way to utilize this - you should always model your storage fitting to your needs. 

An example of how your datastructure could look like:
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
In this case House is the most parental entity. It has two child entities, representing the houses floors. Those themself got multiple children, representing the rooms in that floor. The properties are used to hold some information about the datasets. Deciding if an information should be a property or a new child node itself should depend on your usage.

While the data inside SlingshotDB can be mapped in a network like way, the input and output format are flattened onto a tree like format. This allows you to minimize the amount of data stored in the database, while still beeing able to retrieve an easy to parse and use format.

In the current state SlingshotDB can be used via the [HTTP API](https://github.com/voodooEntity/slingshotdb/blob/master/docs/HTTP_API_V1.md). 

While the database is based on in-memory operations, it offers the option to have asynchronous persistency. This will have a rather small impact to the write/update/delete actions. If you enable persistance, the database will import all persistance datasets on startup.

The database itself does not ship with a user/permission management (like elastic). This decision was made because i think your security should not rely on all your softwares implementation of such, instead you should use things like 'api-gateways' or smiliar to achieve the security management you want/need.

To tackle clientside race-condition problems the database implements a version number for each dataset. When updating a dataset you need to retrieve the dataset before to know the current version. In case your client is working from multple parallel instances, there is always the chance that two instances try to edit the same dataset. To prevent this the database checks if the version in your update dataset matches the current one in the database. On mismatch the update action will be canceld and the client will be informed with a proper error message.