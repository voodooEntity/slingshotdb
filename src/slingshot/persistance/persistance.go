package persistance

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"slingshot/config"
	"slingshot/types"
	"strconv"
	"time"
)

// persistance flag
var PersistanceChan chan types.PersistancePayload
var PersistanceFlag bool

func Boot() chan types.PersistancePayload {
	// first we make sure we have a storage directories and they are writable
	err := handleDirectory("storage/entities/")
	if nil != err {
		config.Logger.Print(err.Error())
		os.Exit(1)
	}
	err = handleDirectory("storage/relations/")
	if nil != err {
		config.Logger.Print(err.Error())
		os.Exit(1)
	}

	// if we got persistance
	if true == config.Conf.Persistance {
		// lets created the persistance channel & flag
		PersistanceFlag = true
		PersistanceChan = make(chan types.PersistancePayload, 100000)
		// and add the temporary channel for import
		importChan := make(chan types.PersistancePayload, 1000000)
		go startWorker(importChan)
		return importChan
	}
	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
//  worker
func startWorker(importChan chan types.PersistancePayload) {
	config.Logger.Print("> Persistance worker started")
	// first we import existing data
	importData(importChan)
	// now we handle further persistance
	var err error

	for elem := range PersistanceChan {
		switch elem.Type {
		case "Entity":
			err = handleEntity(elem)
		case "Relation":
			err = handleRelation(elem)
		case "EntityType":
			err = handleEntityType(elem)
		}
		if nil != err {
			config.Logger.Print(err.Error())
			os.Exit(1)
		}
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// Handle storage directory
func handleDirectory(direcory string) error {
	if _, err := os.Stat(direcory); os.IsNotExist(err) {
		// directory doesnt exist, lets create it
		dirErr := os.MkdirAll(direcory, os.ModePerm)
		if nil != dirErr {
			return dirErr
		}
	}
	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// create a  entity on hdd
func createEntity(entity types.StorageEntity) error {
	// first we handle the path
	path := "storage/entities/" + strconv.Itoa(entity.Type) + "/"
	err := handleDirectory(path)
	if nil != err {
		return err
	}

	// path seems to exist , lets create
	// the entity json
	data, err := json.Marshal(entity)
	if nil != err {
		return err
	}

	// we got the json, lets extend the path
	// variable and write the file
	path = path + "/" + strconv.Itoa(entity.ID)
	err = writeFile(data, path)
	if nil != err {
		return err
	}

	return nil
}

func deleteEntity(Type int, id int) error {
	// first we handle the path
	path := "storage/entities/" + strconv.Itoa(Type) + "/" + strconv.Itoa(id)
	var err = os.Remove(path)
	if nil != err {
		return err
	}

	return nil
}

func updateEntity(entity types.StorageEntity) error {
	// first we handle the path
	path := "storage/entities/" + strconv.Itoa(entity.Type) + "/"
	err := handleDirectory(path)
	if nil != err {
		return err
	}

	// path seems to exist , lets create
	// the entity json
	data, err := json.Marshal(entity)
	if nil != err {
		return err
	}

	// we got the json, lets extend the path
	// variable and write the file
	path = path + "/" + strconv.Itoa(entity.ID)
	err = writeFile(data, path)
	if nil != err {
		return err
	}

	return nil
}

func createRelation(relation types.StorageRelation) error {
	// first we handle the path
	path := "storage/relations/" + strconv.Itoa(relation.SourceType) + "/" + strconv.Itoa(relation.SourceID) + "/" + strconv.Itoa(relation.TargetType)
	err := handleDirectory(path)
	if nil != err {
		return err
	}

	// path seems to exist , lets create
	// the entity json
	data, err := json.Marshal(relation)
	if nil != err {
		return err
	}

	// we got the json, lets extend the path
	// variable and write the file
	path = path + "/" + strconv.Itoa(relation.TargetID)
	err = writeFile(data, path)
	if nil != err {
		return err
	}

	return nil
}

func updateRelation(relation types.StorageRelation) error {
	// first we handle the path
	path := "storage/relations/" + strconv.Itoa(relation.SourceType) + "/" + strconv.Itoa(relation.SourceID) + "/" + strconv.Itoa(relation.TargetType)
	err := handleDirectory(path)
	if nil != err {
		return err
	}

	// path seems to exist , lets create
	// the entity json
	data, err := json.Marshal(relation)
	if nil != err {
		return err
	}

	// we got the json, lets extend the path
	// variable and write the file
	path = path + "/" + strconv.Itoa(relation.TargetID)
	err = writeFile(data, path)
	if nil != err {
		return err
	}

	return nil
}

func deleteRelation(srcType int, srcID int, targetType int, targetID int) error {
	// first we handle the path
	path := "storage/relations/" + strconv.Itoa(srcType) + "/" + strconv.Itoa(srcID) + "/" + strconv.Itoa(targetType) + "/" + strconv.Itoa(targetID)
	var err = os.Remove(path)
	if nil != err {
		return err
	}

	return nil
}

func writeFile(content []byte, path string) error {
	// write content to file
	err := ioutil.WriteFile(path, content, 0777)
	if nil != err {
		return err
	}
	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// Handle entity
func handleEntity(payload types.PersistancePayload) error {
	var err error
	// dispatch the action
	switch payload.Method {
	case "Create":
		err = createEntity(payload.Entity)
	case "Update":
		err = updateEntity(payload.Entity)
	case "Delete":
		err = deleteEntity(payload.Entity.Type, payload.Entity.ID)
	}

	// an error occured?
	if nil != err {
		return err
	}

	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// Handle relation
func handleRelation(payload types.PersistancePayload) error {
	var err error
	// dispatch the action
	switch payload.Method {
	case "Create":
		err = createRelation(payload.Relation)
	case "Update":
		err = updateRelation(payload.Relation)
	case "Delete":
		err = deleteRelation(payload.Relation.SourceType, payload.Relation.SourceID, payload.Relation.TargetType, payload.Relation.TargetID)
	}

	// an error occured?
	if nil != err {
		return err
	}

	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// Handle entity type
func handleEntityType(payload types.PersistancePayload) error {
	//ltes marshall the entity type map
	data, err := json.Marshal(payload.EntityTypes)
	if nil != err {
		return err
	}

	// we got the json, lets build the path and write the file
	path := "storage/entityTypes"
	err = writeFile(data, path)
	if nil != err {
		return err
	}

	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// import the persistance data
func importData(importChan chan types.PersistancePayload) {
	// first we import the entity types
	config.Logger.Print("> - Importing entity types")
	importEntityTypes(importChan)

	// than the entities
	config.Logger.Print("> - Importing entities")
	importEntities(importChan)

	// and finally the relations
	config.Logger.Print("> - Importing relations")
	importRelations(importChan)

	// finally we check on the channel until its empty to close it
	for 0 < len(importChan) {
		time.Sleep(1000000)
	}
	config.Logger.Print("> Closing import channel")
	close(importChan)

}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// import the persistance types
func importEntityTypes(importChan chan types.PersistancePayload) {
	// first we read the entityTypes file
	entityTypesJsonBytes, err := readFile("storage/entityTypes")

	// if it is an error
	if nil != err {
		config.Logger.Print(err.Error())
		os.Exit(1)
	}

	// seems fine lets unmarshall it
	var entityTypes map[int]string
	err = json.Unmarshal(entityTypesJsonBytes, &entityTypes)
	if nil != err {
		config.Logger.Print(err.Error())
		os.Exit(1)
	}

	// ok we got the entity types, lets pack a payload and send it to storage
	payload := types.PersistancePayload{
		EntityTypes: entityTypes,
		Type:        "EntityTypes",
	}
	importChan <- payload
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// import the entities
func importEntities(importChan chan types.PersistancePayload) {
	// predefinitions
	path := "storage/entities/"

	// first we get the entity type dirs
	entityTypeDirs, _ := readDir(path)

	// if there are any entity type dirs
	if 0 < len(entityTypeDirs) {
		// now we walk through the entity Type dirs
		for _, entityType := range entityTypeDirs {
			// now we check the directory
			entityTypePath := path + "/" + entityType
			entityIDs, _ := readDir(entityTypePath)

			// if there are any
			if 0 < len(entityIDs) {
				// now we talk through the entity IDs
				for _, entityID := range entityIDs {
					// read the file
					file := entityTypePath + "/" + entityID
					entityFile, _ := readFile(file)

					// seems fine lets unmarshall it
					var entity types.StorageEntity
					err := json.Unmarshal(entityFile, &entity)
					if nil != err {
						config.Logger.Print(err.Error())
						os.Exit(1)
					}

					// ok this worked, so we pack it
					// into a persistance payload and send it
					payload := types.PersistancePayload{
						Type:   "Entity",
						Entity: entity,
					}
					importChan <- payload
				}
			}
		}
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// import the relation
func importRelations(importChan chan types.PersistancePayload) {
	// predefinitions
	path := "storage/relations/"

	// first we get the source Type id's
	srcTypeIDs, _ := readDir(path)

	// if there are any source type dir
	if 0 < len(srcTypeIDs) {
		// now we walk through the source id directories
		for _, srcTypeID := range srcTypeIDs {
			// get then fitting srcIDs
			srcTypeIDpath := path + srcTypeID + "/"
			srcIDs, _ := readDir(srcTypeIDpath)

			// if there are any
			if 0 < len(srcIDs) {
				// now we walk through the source id's
				for _, srcID := range srcIDs {
					// get then fitting target type id's
					srcIDpath := srcTypeIDpath + srcID + "/"
					targetTypeIDs, _ := readDir(srcIDpath)

					// if there are any
					if 0 < len(targetTypeIDs) {
						// walk through the target type ids
						for _, targetTypeID := range targetTypeIDs {
							// get then fitting target ids
							targetTypePath := srcIDpath + targetTypeID + "/"
							targetIDs, _ := readDir(targetTypePath)

							// if there are any
							if 0 < len(targetIDs) {
								// walk through all the target ID's
								for _, targetID := range targetIDs {
									// finally we made it !!!! epic...
									// ltes build the path and retrieve the relation
									fullpath := targetTypePath + targetID
									relationBytes, _ := readFile(fullpath)

									// seems fine lets unmarshall it
									var relation types.StorageRelation
									err := json.Unmarshal(relationBytes, &relation)
									if nil != err {
										config.Logger.Print(err.Error())
										os.Exit(1)
									}

									// ok this worked, so we pack it
									// into a persistance payload and send it
									payload := types.PersistancePayload{
										Type:     "Relation",
										Relation: relation,
									}
									importChan <- payload
								}
							}

						}
					}
				}

			}

			// now we check the directory
			//entityTypePath := path + "/" + entityType
			//entityIDs, _ := readDir(entityTypePath)
		}
	}
}

func readDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func readFile(filePath string) ([]byte, error) {
	// first we read the json data
	data, err := ioutil.ReadFile(filePath)
	if nil != err {
		config.Logger.Print("> Error reading persistant storage file. Check your permissions")
		os.Exit(1)
	}
	return data, nil
}
