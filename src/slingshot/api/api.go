package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"slingshot/config"
	"slingshot/mapper"
	"slingshot/storage"
	"slingshot/types"
	"strconv"
)

func Start() {
	fmt.Println("> Bootin HTTP API")
	h := http.NewServeMux()

	// Route: /v1/ping
	h.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		respond("pong", 200, w)
	})

	// Route: /v1/getEntityByTypeAndId
	h.HandleFunc("/v1/getEntityByTypeAndId", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["id"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// now we get optional params
		optionalUrlParams := make(map[string]string)
		optionalUrlParams["traverse"] = ""
		urlParams = getOptionalUrlParams(optionalUrlParams, urlParams, r)

		// int conv id
		id, err := strconv.Atoi(urlParams["id"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// lets check if we got an traverse params,
		// and if so intval it for further use
		traverse := 0
		if _, ok := urlParams["traverse"]; urlParams["traverse"] != "" && ok {
			parsedTraverse, err := strconv.Atoi(urlParams["traverse"])
			if nil != err {
				// handle error
				http.Error(w, "Invalid param traverseid given", 422)
				return
			}
			traverse = parsedTraverse
		}

		// ok we seem to be fine on types, lets call the actual getter method
		// based on if we need to get recursive or just an single entry
		returnData := types.MapperTransport{}
		if 0 == traverse {
			// read the data
			data, err := mapper.GetEntityByTypeAndId(urlParams["type"], id)
			returnData = data

			// if error respond
			if nil != err {
				http.Error(w, string(err.Error()), 404)
				return
			}
		} else {
			// get the data
			relations := []types.MapperRelation{}
			entityTypeID, _ := storage.CreateEntityType(urlParams["type"])
			data, err := mapper.GetEntitiesRecursive(entityTypeID, id, traverse, &relations)
			returnData.Entities = append(returnData.Entities, data)
			returnData.Relations = relations

			// if error respond error
			if nil != err {
				http.Error(w, string(err.Error()), 404)
				return
			}
		}

		// all seems fine lets return the data
		respondOk(returnData, w)
	})

	// Route: /v1/createEntity
	h.HandleFunc("/v1/createEntity", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "POST" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// retrieve data from request
		body, err := getRequestBody(r)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed or no body. ", 422)
			return
		}

		// lets see if the body json is valid tho
		var newEntity types.MapperEntity
		err = json.Unmarshal(body, &newEntity)
		if nil != err {
			http.Error(w, "Malformed json body.", 422)
			return
		}

		// finally we create the entity
		responseData, err := mapper.CreateEntity(newEntity.Type, newEntity.Value, newEntity.Properties, newEntity.Context)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(responseData, w)
	})

	// Route: /v1/getEntitiesByType
	h.HandleFunc("/v1/getEntitiesByType", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 403)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// now we get optional params
		optionalUrlParams := make(map[string]string)
		optionalUrlParams["context"] = ""
		urlParams = getOptionalUrlParams(optionalUrlParams, urlParams, r)

		// lets make a default for mode and
		// overwrite if given
		context := ""
		if _, ok := urlParams["context"]; ok {
			context = urlParams["context"]
		}

		// ok we seem to be fine on types, lets call the actual getter method
		responseData, err := mapper.GetEntitiesByType(urlParams["type"], context)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		// all seems fine lets return the data
		respondOk(responseData, w)
	})

	// Route: /v1/getEntitiesByTypeAndValue
	h.HandleFunc("/v1/getEntitiesByTypeAndValue", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["value"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// now we get optional params
		optionalUrlParams := make(map[string]string)
		optionalUrlParams["mode"] = ""
		optionalUrlParams["context"] = ""
		urlParams = getOptionalUrlParams(optionalUrlParams, urlParams, r)

		// lets make a default for mode and
		// overwrite if given
		mode := "match"
		if _, ok := urlParams["mode"]; urlParams["mode"] != "" && ok {
			mode = urlParams["mode"]
		}

		// lets make a default for mode and
		// overwrite if given
		context := ""
		if _, ok := urlParams["context"]; ok {
			context = urlParams["context"]
		}

		// ok we seem to be fine on types, lets call the actual getter method
		responseData, err := mapper.GetEntitiesByTypeAndValue(urlParams["type"], urlParams["value"], mode, context)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		// all seems fine lets return the data
		respondOk(responseData, w)
	})

	// Route: /v1/mapJson
	h.HandleFunc("/v1/mapJson", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "POST" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// retrieve data from request
		body, err := getRequestBody(r)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed or no body. ", 422)
			return
		}

		// lets pass the body to our mapper
		// that will recursive map the entities
		responseData, err := mapper.MapJson(body)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(responseData, w)
	})

	// Route: /v1/deleteEntity
	h.HandleFunc("/v1/deleteEntity", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "DELETE" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["id"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		// now we get optional params ###todo add unsafe param for delete without relations
		//optionalUrlParams := make(map[string]string)
		//optionalUrlParams["traverse"] = ""
		//urlParams = getOptionalUrlParams(optionalUrlParams, urlParams, r)

		// int conv id
		id, err := strconv.Atoi(urlParams["id"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// finally we delete the entity
		mapper.DeleteEntity(urlParams["type"], id)

		respond("", 200, w)
	})

	// Route: /v1/updateEntity
	h.HandleFunc("/v1/updateEntity", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "PUT" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// retrieve data from request
		body, err := getRequestBody(r)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed or no body. ", 422)
			return
		}

		// lets see if the body json is valid tho
		var newEntity types.MapperEntity
		err = json.Unmarshal(body, &newEntity)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed json body.", 422)
			return
		}

		// finally we update the entity
		err = mapper.UpdateEntity(newEntity.Type, newEntity.ID, newEntity.Value, newEntity.Properties, newEntity.Context, newEntity.Version)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respond("", 200, w)
	})

	// Route: /v1/getChildEntities
	h.HandleFunc("/v1/getChildEntities", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["id"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// int conv id
		id, err := strconv.Atoi(urlParams["id"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// retrieve the child entities if given
		transport, err := mapper.GetChildEntities(urlParams["type"], id)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(transport, w)
	})

	// Route: /v1/getParentEntities
	h.HandleFunc("/v1/getParentEntities", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["id"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// int conv id
		id, err := strconv.Atoi(urlParams["id"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// retrieve the child entities if given
		transport, err := mapper.GetParentEntities(urlParams["type"], id)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(transport, w)
	})

	// Route: /v1/getRelationsTo
	h.HandleFunc("/v1/getRelationsTo", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["id"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// int conv id
		id, err := strconv.Atoi(urlParams["id"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// retrieve the child entities if given
		transport, err := mapper.GetRelationsTo(urlParams["type"], id)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(transport, w)
	})

	// Route: /v1/getRelationsTo
	h.HandleFunc("/v1/getRelation", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["srcType"] = ""
		requiredUrlParams["srcID"] = ""
		requiredUrlParams["targetType"] = ""
		requiredUrlParams["targetID"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// int conv id
		srcID, err := strconv.Atoi(urlParams["srcID"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// int conv id
		targetID, err := strconv.Atoi(urlParams["targetID"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// retrieve the child entities if given
		transport, err := mapper.GetRelation(urlParams["srcType"], srcID, urlParams["targetType"], targetID)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(transport, w)
	})

	// Route: /v1/getRelationsFrom
	h.HandleFunc("/v1/getRelationsFrom", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["type"] = ""
		requiredUrlParams["id"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// int conv id
		id, err := strconv.Atoi(urlParams["id"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// retrieve the child entities if given
		transport, err := mapper.GetRelationsFrom(urlParams["type"], id)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(transport, w)
	})

	// Route: /v1/getEntitiesByValue
	h.HandleFunc("/v1/getEntitiesByValue", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["value"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r) // todo

		// required params check
		if nil != err {
			// handle error
			http.Error(w, err.Error(), 422)
			return
		}

		// now we get optional params
		optionalUrlParams := make(map[string]string)
		optionalUrlParams["mode"] = ""
		optionalUrlParams["context"] = ""
		urlParams = getOptionalUrlParams(optionalUrlParams, urlParams, r)

		// lets make a default for mode and
		// overwrite if given
		mode := "match"
		if _, ok := urlParams["mode"]; urlParams["mode"] != "" && ok {
			mode = urlParams["mode"]
		}

		// lets make a default context and
		// overwrite if given
		context := ""
		if _, ok := urlParams["context"]; ok {
			context = urlParams["context"]
		}

		// retrieve the entities
		transport, err := mapper.GetEntitiesByValue(urlParams["value"], mode, context)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respondOk(transport, w)
	})

	// Route: /v1/getEntityTypes
	h.HandleFunc("/v1/getEntityTypes", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "GET" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// retrieve all entity types
		types := mapper.GetEntityTypes()

		// than we gonne json encode it
		// build the json
		responseData, err := json.Marshal(types)
		if nil != err {
			http.Error(w, "Error building response data json", 500)
			return
		}
		// finally we gonne send our response
		respond(string(responseData), 200, w)
	})

	// Route: /v1/updateRelation
	h.HandleFunc("/v1/updateRelation", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "PUT" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// retrieve data from request
		body, err := getRequestBody(r)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed or no body. ", 422)
			return
		}

		// lets see if the body json is valid tho
		var newRelation types.MapperRelation
		err = json.Unmarshal(body, &newRelation)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed json body.", 422)
			return
		}

		// finally we update the entity
		err = mapper.UpdateRelation(newRelation.SourceType, newRelation.SourceID, newRelation.TargetType, newRelation.TargetID, newRelation)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respond("", 200, w)
	})

	// Route: /v1/createRelation
	h.HandleFunc("/v1/createRelation", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "POST" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// retrieve data from request
		body, err := getRequestBody(r)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed or no body. ", 422)
			return
		}

		// lets see if the body json is valid tho
		var newRelation types.MapperRelation
		err = json.Unmarshal(body, &newRelation)
		if nil != err {
			fmt.Print(err.Error())
			http.Error(w, "Malformed json body.", 422)
			return
		}

		// finally we update the entity
		err = mapper.CreateRelation(newRelation.SourceType, newRelation.SourceID, newRelation.TargetType, newRelation.TargetID, newRelation)
		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		respond("", 200, w)
	})

	// Route: /v1/createRelation
	h.HandleFunc("/v1/deleteRelation", func(w http.ResponseWriter, r *http.Request) {
		// check http method
		if "DELETE" != r.Method {
			http.Error(w, "Invalid http method for this path", 422)
			return
		}

		// first we get the params
		requiredUrlParams := make(map[string]string)
		requiredUrlParams["srcType"] = ""
		requiredUrlParams["srcID"] = ""
		requiredUrlParams["targetType"] = ""
		requiredUrlParams["targetID"] = ""
		urlParams, err := getRequiredUrlParams(requiredUrlParams, r)

		if nil != err {
			http.Error(w, err.Error(), 422)
			return
		}

		// int conv id's
		srcID, err := strconv.Atoi(urlParams["srcID"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// int conv id's
		targetID, err := strconv.Atoi(urlParams["targetID"])
		if nil != err {
			// handle error
			http.Error(w, "Invalid param id given", 422)
			return
		}

		// finally we delete the entity
		mapper.DeleteRelation(urlParams["srcType"], srcID, urlParams["targetType"], targetID)

		respond("", 200, w)
	})

	// -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
	// NOT IMPLEMENTED YET (seperator)
	// -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -

	// Route: /v1/template
	//h.HandleFunc("/v1/template", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintln(w, "Hello, you hit foo!")
	//})

	// building server listen string by
	// config values and print it - than listen
	connectString := buildHttpListenConfigString()
	fmt.Println("> Server listening settings by config (", connectString, ")")
	http.ListenAndServe(connectString, h)
}

func getOptionalUrlParams(optionalUrlParams map[string]string, urlParams map[string]string, r *http.Request) map[string]string {
	tmpParams := r.URL.Query()
	for paramName := range optionalUrlParams {
		val, ok := tmpParams[paramName]
		if ok {
			urlParams[paramName] = val[0]
		}
	}
	return urlParams
}

func getRequiredUrlParams(requiredUrlParams map[string]string, r *http.Request) (map[string]string, error) {
	urlParams := r.URL.Query()
	for paramName := range requiredUrlParams {
		val, ok := urlParams[paramName]
		if !ok {
			return nil, errors.New("Missing required url param")
		}
		requiredUrlParams[paramName] = val[0]
	}
	return requiredUrlParams, nil
}

func respond(message string, responseCode int, w http.ResponseWriter) {
	w.WriteHeader(responseCode)
	messageBytes := []byte(message)
	_, err := w.Write(messageBytes)
	if nil != err {
		fmt.Print(err)
	}
}

func respondOk(data types.MapperTransport, w http.ResponseWriter) {
	// than we gonne json encode it
	// build the json
	responseData, err := json.Marshal(data)
	if nil != err {
		http.Error(w, "Error building response data json", 500)
		return
	}

	// finally we gonne send our response
	w.WriteHeader(200)
	_, err = w.Write(responseData)
	if nil != err {
		fmt.Print(err)
	}
}

func getRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func buildHttpListenConfigString() string {
	var connectString string
	connectString += config.Conf.Host
	connectString += ":"
	connectString += strconv.Itoa(config.Conf.Port)
	return connectString
}
