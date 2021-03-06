# SlingshotDB 
Welcome to the home of SlingshotDB - an in-memory graph database. It's completly written in golang (vanilla, no 3rd party libraries used) and provides acces via an [HTTP API](https://github.com/voodooEntity/slingshotdb/blob/master/docs/HTTP_API_V1.md). For further information about the structure/usage of this software please check [About Slingshot](https://github.com/voodooEntity/slingshotdb/blob/master/docs/ABOUT_SLINGSHOT.md)

The main target of the database is to have a easy to use high performance storage. In it's current state the database is shipped with minimal functionality. I will extend the functionality over time based on needs and time. 

`Important!` Since the database is in memory and focused on high performance it uses a lot of memory for indexes and data storage. Using the database you should make sure to run it on a machine that provides a lot of memory to work fine. While there are several ways how i could reducde the memory usage i accepted this trade off for the performance.

Finally i wanne leave a special thanks to some friends that helped through the process of creating this software by listening to hours of rage/ideas and providing suggestions that lead the way to the software you are about to use. 
* Maze (the name 'Slingshot' was his idea)
* Luxer 
* f0o

I hope you will enjoy the usage of SlingshotDB and that this will just be the start of a great project.

Sincerely yours,
voodooEntity

---
## Contributing 
Please read the  [Code of Condut](https://github.com/voodooEntity/slingshotdb/blob/master/CODE_OF_CONDUCT.md) and [Contributing](https://github.com/voodooEntity/slingshotdb/blob/master/CONTRIBUTING.md) docs in advance.     

Since this is an early alpha you should expect to encounter bugs that i didn't takle yet. Please feel free to report those as issues, and i will do my best to fix/solve them asap.     

If you have suggestions on how to change the database, you can also feel free to message me, but keep in mind that i coded this database for specific purposes. So even if your suggestions may seem legit, they still could be refused if they oppose the way this database is meant to work.      


## Build from source
SlingshotDB needs golang to be build. Just add the directory path you cloned the database into to your gopath and fire 
```bash
go build -o slingshot
```
inside. No special build flags or installing of dependencies needed. After building, just start the database with 
```bash
./slingshot
```

## Run as docker
SlingshotDB is shipped with a Dockerfile. This will compile the current cloned version and can be run. Using this you can easily update the git repo und just rebuild the docker iamge to keep your software up to date An example of how you could do it:    

```bash
# first we build the docker image tagged as slingshot db
$ docker build . -t slingshotdb

# than we gonne run the docker image as container exposing port 8090
$ docker run -p 8090:8090 slingshotdb

# you can modify the port the database will use in your host by changing
# the first number after the -p param - if you want to run it on port 1234
# it would look like following 
$ docker run -p 1234:8090 slingshotdb

# to enable persistance in the host directory use the following
# example wich mounts the current directory into the docker to
# enable the database write into the host storage directory
$ docker run -p 8090:8090 --mount type=bind,source="$(pwd)"/,target=/go/src/app slingshotdb

# for further options check the docker doku
# https://docs.docker.com/engine/reference/commandline/run/
```


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
After building and starting the database, you are free to use it. You can create your own code to use the API based on the [HTTP API v1 docs](https://github.com/voodooEntity/slingshotdb/blob/master/docs/HTTP_API_V1.md), or if your environment is PHP you can use the 'voodooEntity/slingshotdb-php-sdk'. As time goes by i plan to add more SDK's (golang, javascript, ....)

## Future plans
* Extend binary to create storage directories intially
* Extend mapJson method to include existing entities
* Adding a network retrievel function that will ship the data in single dataset lists. This enables the traversing through the data in a non flattened structure.
* Adding offest and length parameters to retrieval HTTP methods that potentially can return giant amounts of data
* Adding optional filter params for properties
* Extended Docs with more examples
* `System` methods such as `getMemoryStats` , `shutdown` and more
* A simple tool to visualize the data (probably browser based)
* ..... and much more =)
