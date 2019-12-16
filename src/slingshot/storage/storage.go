package storage

// handle all the imports
import (
	"errors"
	"slingshot/config"
	"slingshot/persistance"
	"slingshot/types"
	"sync"
)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage map            [Type] [ID]
var EntityStorage = make(map[int]map[int]types.StorageEntity)

// entity storage master mutex
var EntityStorageMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity storage id max         [Type]
var EntityIDMax = make(map[int]int)

// master mutexd for EntityIdMax
var EntityIDMaxMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// maps to translate Types to their INT and reverse
var EntityTypes = make(map[int]string)
var EntityRTypes = make(map[string]int)

// and a fitting max ID
var EntityTypeIDMax int = 0

// entity Type mutex (for adding and deleting Type types)
var EntityTypeMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// s prefix = source
// t prefix = target
// relation storage map             [sType][sId]   [tType][tId]
var RelationStorage = make(map[int]map[int]map[int]map[int]types.StorageRelation)

// and relation reverse storage map
// (for faster queries)              [tType][Tid]   [sType][sId]
var RelationRStorage = make(map[int]map[int]map[int]map[int]bool)

// relation storage master mutex
var RelationStorageMutex = &sync.RWMutex{}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// + + + + + + + + + +  PUBLIC  + + + + + + + + + + +
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// init/construct function for storage package
func Boot() {
	// check for persistance
	persist := config.Conf.Persistance
	if true == persist {
		// lets boot the persistance worker and check
		// if we get eny import data
		importChan := persistance.Boot()
		if nil != importChan {
			handleImport(importChan)
		}
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// Create an entity Type
func CreateEntityType(name string) (int, error) {
	// first of allw e lock
	//printMutexActions("CreateEntityType.EntityTypeMutex.Lock");
	EntityTypeMutex.Lock()
	// lets check if the Type allready exists
	// if it does we just return the ID
	if id, ok := EntityRTypes[name]; ok {
		// dont forget to unlock
		//printMutexActions("CreateEntityType.EntityTypeMutex.Unlock");
		EntityTypeMutex.Unlock()
		return id, nil
	}
	// ok entity doesnt exist yet, lets
	// upcount our ID Max and copy it
	// into another variable so we can be sure
	// between unlock of the ressource and return
	// it doesnt get upcounted
	EntityTypeIDMax++
	var newID = EntityTypeIDMax
	// finally create the new Type in our
	// EntityTypes index and reverse index
	EntityTypes[newID] = name
	EntityRTypes[name] = newID
	// and create mutex for EntityStorage Type+type
	//printMutexActions("CreateEntityType.EntityStorageMutex.Lock");
	EntityStorageMutex.Lock()
	// now we prepare the submaps in the entity
	// storage itseöf....
	EntityStorage[newID] = make(map[int]types.StorageEntity)
	// set the maxID for the new
	// Type type
	EntityIDMax[newID] = 0
	//printMutexActions("CreateEntityType.EntityStorageMutex.Unlock");
	EntityStorageMutex.Unlock()
	// create the base maps in relation storage
	//printMutexActions("CreateEntityType.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	RelationStorage[newID] = make(map[int]map[int]map[int]types.StorageRelation)
	RelationRStorage[newID] = make(map[int]map[int]map[int]bool)
	//printMutexActions("CreateEntityType.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	// and create the basic submaps for
	// the relation storage
	// now we unlock the mutex
	// and return the new id
	//printMutexActions("CreateEntityType.EntityTypeMutex.Unlock");
	// - - - - - - - - - - - - - - - - -
	// persistance handling
	if true == persistance.PersistanceFlag {
		persistance.PersistanceChan <- types.PersistancePayload{
			Type:        "EntityType",
			EntityTypes: EntityTypes,
		}
	}
	// - - - - - - - - - - - - - - - - -
	EntityTypeMutex.Unlock()
	return newID, nil
}

func CreateEntity(entity types.StorageEntity) (int, error) {
	//types.PrintMemUsage()
	// first we lock the entity Type mutex
	// to make sure while we check for the
	// existence it doesnt get deletet, this
	// may sound like a very rare upcoming case,
	//but better be safe than sorry
	//printMutexActions("CreateEntity.EntityTypeMutex.RLock");
	EntityTypeMutex.RLock()
	// now
	if _, ok := EntityTypes[entity.Type]; !ok {
		// the Type doest exist, lets unlock
		// the Type mutex and return -1 for fail0r
		//printMutexActions("CreateEntity.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return -1, errors.New("CreateEntity.Entity Type not existing")
	}
	// the Type seems to exist, now lets lock the
	// storage mutex before Unlocking the Entity
	// Type mutex to prevent the Type beeing
	// deleted before we start locking (small
	// timing still possible )
	//printMutexActions("CreateEntity.EntityTypeMutex.RUnlock");
	EntityTypeMutex.RUnlock()
	// upcount our ID Max and copy it
	// into another variable so we can be sure
	// between unlock of the ressource and return
	// it doesnt get upcounted
	// and set the IDMaxMutex on write Lock
	// lets upcount the entity id max fitting to
	//         [Type]
	//printMutexActions("CreateEntity.EntityStorageMutex.Lock");
	EntityStorageMutex.Lock()
	//fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Unlock");
	//EntityIDMaxMasterMutex.Unlock()
	EntityIDMax[entity.Type]++
	var newID = EntityIDMax[entity.Type]
	//fmt.Println("CreateEntity.EntityIDMaxMasterMutex.Lock");
	//EntityIDMaxMasterMutex.Lock()
	// and tell the entity its own id
	entity.ID = newID
	// - - - - - - - - - - - - - - - - -
	// persistance handling
	if true == persistance.PersistanceFlag {
		persistance.PersistanceChan <- types.PersistancePayload{
			Type:   "Entity",
			Method: "Create",
			Entity: entity,
		}
	}
	// - - - - - - - - - - - - - - - - -
	// now we store the entity element
	// in the EntityStorage
	EntityStorage[entity.Type][newID] = entity
	//printMutexActions("CreateEntity.EntityStorageMutex.Unlock");
	EntityStorageMutex.Unlock()
	// create the mutex for our ressource on
	// relation. we have to create the sub maps too
	// golang things....
	//printMutexActions("CreateEntity.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	RelationStorage[entity.Type][newID] = make(map[int]map[int]types.StorageRelation)
	RelationRStorage[entity.Type][newID] = make(map[int]map[int]bool)
	//printMutexActions("CreateEntity.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	// since we now stored the entity and created all
	// needed ressources we can unlock
	// the storage ressource and return the ID (or err)
	return newID, nil
}

func GetEntityByPath(Type int, id int) (types.StorageEntity, error) {
	// lets check if entity witrh the given path exists
	//printMutexActions("GetEntityByPath.EntityStorageMutex.Lock");
	EntityStorageMutex.Lock()
	if entity, ok := EntityStorage[Type][id]; ok {
		// if yes we return the entity
		// and nil for error
		//printMutexActions("GetEntityByPath.EntityStorageMutex.Unlock");
		EntityStorageMutex.Unlock()
		return entity, nil
	}
	//printMutexActions("GetEntityByPath.EntityStorageMutex.Unlock");
	EntityStorageMutex.Unlock()
	// the path seems to result empty , so
	// we throw an error
	return types.StorageEntity{}, errors.New("Entity on given path does not exist.")
}

func GetEntitiesByType(Type string) (map[int]types.StorageEntity, error) {
	// retrieve the fitting id
	entityTypeID, _ := GetTypeIdByString(Type)
	// lock retrieve und unlock the storage
	EntityStorageMutex.RLock()
	entities := EntityStorage[entityTypeID]
	EntityStorageMutex.RUnlock()
	// return the entity map
	return entities, nil
}

func GetEntitiesByValue(value string) map[int]types.StorageEntity {
	// lets prepare the return map
	entities := make(map[int]types.StorageEntity)
	i := 0

	// first we lock the storage
	EntityStorageMutex.RLock()

	// than we iterate through all entity storage to find a fitting value
	if 0 < len(EntityStorage) {
		for typeID, _ := range EntityStorage {
			if 0 < len(EntityStorage[typeID]) {
				for _, entity := range EntityStorage[typeID] {
					if entity.Value == value {
						entities[i] = entity
						i++
					}
				}
			}
		}
	}

	// unlock storage again and return
	EntityStorageMutex.RUnlock()
	return entities
}

func UpdateEntity(entity types.StorageEntity) error {
	EntityStorageMutex.Lock()
	if _, ok := EntityStorage[entity.Type][entity.ID]; ok {
		// - - - - - - - - - - - - - - - - -
		// persistance handling
		if true == persistance.PersistanceFlag {
			persistance.PersistanceChan <- types.PersistancePayload{
				Type:   "Entity",
				Method: "Update",
				Entity: entity,
			}
		}
		// - - - - - - - - - - - - - - - - -
		EntityStorage[entity.Type][entity.ID] = entity
		EntityStorageMutex.Unlock()
		return nil
	}
	EntityStorageMutex.Unlock()
	return errors.New("Cant update non existing entity")
}

func DeleteEntity(Type int, id int) {
	// we gonne lock the mutex and
	// delete the element
	EntityStorageMutex.Lock()
	// - - - - - - - - - - - - - - - - -
	// persistance handling
	if true == persistance.PersistanceFlag {
		persistance.PersistanceChan <- types.PersistancePayload{
			Type:   "Entity",
			Method: "Delete",
			Entity: types.StorageEntity{
				ID:   id,
				Type: Type,
			},
		}
	}
	// - - - - - - - - - - - - - - - - -
	delete(EntityStorage[Type], id)
	EntityStorageMutex.Unlock()
	// now we delete the relations from and to this entity
	// first child
	DeleteChildRelations(Type, id)
	// than parent
	DeleteParentRelations(Type, id)
}

func GetRelation(srcType int, srcID int, targetType int, targetID int) (types.StorageRelation, error) {
	// first we lock the relation storage
	RelationStorageMutex.RLock()
	if _, firstOk := RelationStorage[srcType]; firstOk {
		if _, secondOk := RelationStorage[srcType][srcID]; secondOk {
			if _, thirdOk := RelationStorage[srcType][srcID][targetType]; thirdOk {
				if relation, fourthOk := RelationStorage[srcType][srcID][targetType][targetID]; fourthOk {
					RelationStorageMutex.RUnlock()
					return relation, nil
				}
			}
		}
	}
	RelationStorageMutex.RUnlock()
	return types.StorageRelation{}, errors.New("Non existing relation requested")
}

// maybe deprecated, check later
func RelationExists(srcType int, srcID int, targetType int, targetID int) bool {
	// first we lock the relation storage
	RelationStorageMutex.RLock()
	if srcTypeMap, firstOk := RelationStorage[srcType]; firstOk {
		if srcIDMap, secondOk := srcTypeMap[srcID]; secondOk {
			if targetTypeMap, thirdOk := srcIDMap[targetType]; thirdOk {
				if _, fourthOk := targetTypeMap[targetID]; fourthOk {
					RelationStorageMutex.RUnlock()
					return true
				}
			}
		}
	}
	RelationStorageMutex.RUnlock()
	return false
}

func DeleteRelationList(relationList map[int]types.StorageRelation) {
	// lets walk through the iterations and delete all
	// corrosponding Relation & RRelation index entries
	if 0 < len(relationList) {
		for _, relation := range relationList {
			DeleteRelation(relation.SourceType, relation.SourceID, relation.TargetType, relation.TargetID)
		}
	}
	return
}

func DeleteRelation(sourceType int, sourceID int, targetType int, targetID int) {
	RelationStorageMutex.Lock()
	// - - - - - - - - - - - - - - - - -
	// persistance handling
	if true == persistance.PersistanceFlag {
		persistance.PersistanceChan <- types.PersistancePayload{
			Type:   "Relation",
			Method: "Delete",
			Relation: types.StorageRelation{
				SourceID:   sourceID,
				SourceType: sourceType,
				TargetID:   targetID,
				TargetType: targetType,
			},
		}
	}
	// - - - - - - - - - - - - - - - - -
	delete(RelationStorage[sourceType][sourceID][targetType], targetID)
	delete(RelationRStorage[targetType][targetID][sourceType], sourceID)
	RelationStorageMutex.Unlock()
}

func DeleteChildRelations(Type int, id int) error {
	childRelations, err := GetChildRelationsBySourceTypeAndSourceId(Type, id)
	if nil != err {
		return err
	}
	DeleteRelationList(childRelations)
	return nil
}

func DeleteParentRelations(Type int, id int) error {
	parentRelations, err := GetParentRelationsByTargetTypeAndTargetId(Type, id)
	if nil != err {
		return err
	}
	DeleteRelationList(parentRelations)
	return nil
}

func CreateRelation(srcType int, srcID int, targetType int, targetID int, relation types.StorageRelation) (bool, error) {
	// first we Readlock the EntityTypeMutex
	//printMutexActions("CreateRelation.EntityTypeMutex.RLock");
	EntityTypeMutex.RLock()
	// lets make sure the source Type exist
	if _, ok := EntityTypes[srcType]; !ok {
		//printMutexActions("CreateRelation.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return false, errors.New("Source Type not existing")
	}
	// and the target Type exists too
	if _, ok := EntityTypes[targetType]; !ok {
		//printMutexActions("CreateRelation.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return false, errors.New("Target Type not existing")
	}
	// finally unlock the TypeMutex again if both checks were successfull
	EntityTypeMutex.RUnlock()
	//// - - - - - - - - - - - - - - - - -
	// now we lock the relation mutex
	//printMutexActions("CreateRelation.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	// lets check if their exists a map for our
	// source entity to the target Type if not
	// create it.... golang things...
	if _, ok := RelationStorage[srcType][srcID][targetType]; !ok {
		RelationStorage[srcType][srcID][targetType] = make(map[int]types.StorageRelation)
		// if the map doesnt exist in this direction
		// it wont exist in the other as in reverse
		// map either so we should create it too
		// but we will store a pointer to the other
		// maps Relation instead of the complete
		// relation twice - defunct, refactor later (may create more problems then help)
		//RelationStorage[targetType][targetID][srcType] = make(map[int]Relation)
	}
	// now we prepare the reverse storage if necessary
	if _, ok := RelationRStorage[targetType][targetID][srcType]; !ok {
		RelationRStorage[targetType][targetID][srcType] = make(map[int]bool)
	}
	// now we store the relation
	RelationStorage[srcType][srcID][targetType][targetID] = relation
	// - - - - - - - - - - - - - - - - -
	// persistance handling
	if true == persistance.PersistanceFlag {
		persistance.PersistanceChan <- types.PersistancePayload{
			Type:     "Relation",
			Method:   "Create",
			Relation: relation,
		}
	}
	// - - - - - - - - - - - - - - - - -
	// and an entry into the reverse index, its existence
	// allows us to use the coords in the normal index to revtrieve
	// the Relation. We dont create a pointer because golang doesnt
	// allow pointer on submaps in nested maps
	RelationRStorage[targetType][targetID][srcType][srcID] = true
	// we are done now we can unlock the entity Types
	//// - - - - - - - - - - - - - - - -
	//and finally unlock the relation Type and return
	//printMutexActions("CreateRelation.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	return true, nil
}

func UpdateRelation(srcType int, srcID int, targetType int, targetID int, relation types.StorageRelation) (types.StorageRelation, error) {
	// first we lock the relation storage
	RelationStorageMutex.RLock()
	if _, firstOk := RelationStorage[srcType]; firstOk {
		if _, secondOk := RelationStorage[srcType][srcID]; secondOk {
			if _, thirdOk := RelationStorage[srcType][srcID][targetType]; thirdOk {
				if rel, fourthOk := RelationStorage[srcType][srcID][targetType][targetID]; fourthOk {
					RelationStorageMutex.RUnlock()
					RelationStorageMutex.Lock()
					rel.Context = relation.Context
					rel.Properties = relation.Properties
					RelationStorage[srcType][srcID][targetType][targetID] = rel
					// - - - - - - - - - - - - - - - - -
					// persistance handling
					if true == persistance.PersistanceFlag {
						persistance.PersistanceChan <- types.PersistancePayload{
							Type:     "Relation",
							Method:   "Create",
							Relation: rel,
						}
					}
					// - - - - - - - - - - - - - - - - -
					RelationStorageMutex.Unlock()
					return relation, nil
				}
			}
		}
	}
	RelationStorageMutex.RUnlock()
	return types.StorageRelation{}, errors.New("Cant update non existing relation")
}

func GetChildRelationsBySourceTypeAndSourceId(Type int, id int) (map[int]types.StorageRelation, error) {
	// initialice the return map
	var mapRet = make(map[int]types.StorageRelation)
	// set counter for the loop
	var cnt = 0
	// copy the pool we have to search in
	// to prevent crashes on RW concurrency
	// we lock the RelationStorage mutex with
	// fitting Type. this allows us to proceed
	// faster since we just block to copy instead
	// of blocking for the whole process
	//printMutexActions("GetChildRelationsBySourceTypeAndSourceId.RelationStorageMutex.Lock");
	RelationStorageMutex.Lock()
	var pool = RelationStorage[Type][id]
	//printMutexActions("GetChildRelationsBySourceTypeAndSourceId.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
	// for each possible targtType
	for _, targetTypeMap := range pool {
		// for each possible targetId per targetType
		for _, relation := range targetTypeMap {
			// copy the relation into the return map
			// and upcount the int
			mapRet[cnt] = relation
			cnt++
		}
	}
	// + + + + + + +
	//fmt.Println("Relations: ",cnt," - ",len(mapRet))
	// + + + + + + +
	return mapRet, nil
}

func GetParentRelationsByTargetTypeAndTargetId(targetType int, targetID int) (map[int]types.StorageRelation, error) {
	// initialice the return map
	var mapRet = make(map[int]types.StorageRelation)
	// set counter for the loop
	var cnt = 0
	// copy the pool we have to search in
	// to prevent crashes on RW concurrency
	// we lock the RelationStorage mutex with
	// fitting Type. this allows us to proceed
	// faster since we just block to copy instead
	// of blocking for the whole process
	//printMutexActions("GetChildRelationsBySourceTypeAndSourceId.RelationStorageMutex.Lock");
	RelationStorageMutex.RLock()
	var pool = RelationRStorage[targetType][targetID]
	//printMutexActions("GetChildRelationsBySourceTypeAndSourceId.RelationStorageMutex.Unlock");
	RelationStorageMutex.RUnlock()
	// for each possible targtType
	for sourceTypeID, targetTypeMap := range pool {
		// for each possible targetId per targetType
		for sourceRelationID, _ := range targetTypeMap {
			// copy the relation into the return map
			// and upcount the int
			RelationStorageMutex.RLock()
			mapRet[cnt] = RelationStorage[sourceTypeID][sourceRelationID][targetType][targetID]
			RelationStorageMutex.RUnlock()
			cnt++
		}
	}
	// + + + + + + +
	//fmt.Println("Relations: ",cnt," - ",len(mapRet))
	// + + + + + + +
	return mapRet, nil
}

func GetEntityTypes() []string {
	// prepare the return array
	types := []string{}

	// now we lock the storage
	EntityTypeMutex.RLock()
	for _, Type := range EntityTypes {
		types = append(types, Type)
	}

	// unlock the mutex and return
	EntityTypeMutex.RUnlock()
	return types
}

func TypeExists(strType string) bool {
	//printMutexActions("TypeExists.EntityTypeMutex.RLock");
	EntityTypeMutex.RLock()
	// lets check if this Type exists
	if _, ok := EntityRTypes[strType]; ok {
		// it does lets return it
		//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return true
	}
	//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
	EntityTypeMutex.RUnlock()
	return false
}

func EntityExists(Type int, id int) bool {
	//printMutexActions("TypeExists.EntityTypeMutex.RLock");
	EntityStorageMutex.RLock()
	// lets check if this Type exists
	if _, ok := EntityStorage[Type][id]; ok {
		// it does lets return it
		//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
		EntityStorageMutex.RUnlock()
		return true
	}
	//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
	EntityStorageMutex.RUnlock()
	return false
}

func TypeIdExists(id int) bool {
	//printMutexActions("TypeExists.EntityTypeMutex.RLock");
	EntityTypeMutex.RLock()
	// lets check if this Type exists
	if _, ok := EntityTypes[id]; ok {
		// it does lets return it
		//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return true
	}
	//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
	EntityTypeMutex.RUnlock()
	return false
}

func GetTypeIdByString(strType string) (int, error) {
	//printMutexActions("TypeExists.EntityTypeMutex.RLock");
	EntityTypeMutex.RLock()
	// lets check if this Type exists
	if id, ok := EntityRTypes[strType]; ok {
		// it does lets return it
		//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return id, nil
	}
	//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
	EntityTypeMutex.RUnlock()
	return -1, errors.New("Entity Type string does not exist")
}

func GetTypeStringById(intType int) (*string, error) {
	//printMutexActions("TypeExists.EntityTypeMutex.RLock");
	EntityTypeMutex.RLock()
	// lets check if this Type exists
	if strType, ok := EntityTypes[intType]; ok {
		// it does lets return it
		//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
		EntityTypeMutex.RUnlock()
		return &strType, nil
	}
	//printMutexActions("TypeExists.EntityTypeMutex.RUnlock");
	EntityTypeMutex.RUnlock()
	return nil, errors.New("Entity Type string does not exist")
}

func printMutexActions(param string) {
	//fmt.Println(param)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// + + + + + + + + + +  PRIVATE  + + + + + + + + + + +
// - - - - - - - - - - - - - - - - - - - - - - - - - -

func handleImport(importChan chan types.PersistancePayload) {
	for elem := range importChan {
		switch elem.Type {
		case "Entity":
			importEntity(elem)
		case "Relation":
			importRelation(elem)
		case "EntityTypes":
			importEntityTypes(elem)
		}
	}
}

func importEntityTypes(payload types.PersistancePayload) {
	// than we lock the entity type mutex and relationstorage mutex
	EntityTypeMutex.Lock()
	EntityStorageMutex.Lock()
	RelationStorageMutex.Lock()
	//presets
	maxID := 0

	// first we build the corrosponding
	// reverse index and determine the maxID
	reverseMap := make(map[string]int)
	for key, value := range payload.EntityTypes {
		// store the reverse index
		reverseMap[value] = key

		// if bigger replace the max ID
		if maxID < key {
			maxID = key
		}

		// and prepare the relation storage
		RelationStorage[key] = make(map[int]map[int]map[int]types.StorageRelation)
		RelationRStorage[key] = make(map[int]map[int]map[int]bool)

		// same as the entity storage
		EntityStorage[key] = make(map[int]types.StorageEntity)
	}

	// store typemap, rmap and max id
	EntityTypeIDMax = maxID
	EntityTypes = payload.EntityTypes
	EntityRTypes = reverseMap

	//  unlock the mutex's again
	RelationStorageMutex.Unlock()
	EntityStorageMutex.Unlock()
	EntityTypeMutex.Unlock()
}

func importEntity(payload types.PersistancePayload) {
	// first we handle the ID max
	EntityIDMaxMutex.Lock()
	if EntityIDMax[payload.Entity.Type] < payload.Entity.ID {
		EntityIDMax[payload.Entity.Type] = payload.Entity.ID
	}
	EntityIDMaxMutex.Unlock()

	// now we create the entity themself
	// first we lock the storage
	EntityStorageMutex.Lock()

	// and put the entity in the EntityStorage
	EntityStorage[payload.Entity.Type][payload.Entity.ID] = payload.Entity

	// than unlock the entity storage again
	EntityStorageMutex.Unlock()

	// now we handle the relations , prepare the maps
	RelationStorageMutex.Lock()

	// create all the maps
	RelationStorage[payload.Entity.Type][payload.Entity.ID] = make(map[int]map[int]types.StorageRelation)
	RelationRStorage[payload.Entity.Type][payload.Entity.ID] = make(map[int]map[int]bool)

	// and unlock the relation storage again
	RelationStorageMutex.Unlock()
}

func importRelation(payload types.PersistancePayload) {
	//// - - - - - - - - - - - - - - - - -
	// now we lock the relation mutex
	RelationStorageMutex.Lock()

	// lets check if their exists a map for our
	// source entity to the target Type if not
	// create it.... golang things...
	if _, ok := RelationStorage[payload.Relation.SourceType][payload.Relation.SourceID][payload.Relation.TargetType]; !ok {
		RelationStorage[payload.Relation.SourceType][payload.Relation.SourceID][payload.Relation.TargetType] = make(map[int]types.StorageRelation)
	}

	// now we prepare the reverse storage if necessary
	if _, ok := RelationRStorage[payload.Relation.TargetType][payload.Relation.TargetID][payload.Relation.SourceType]; !ok {
		RelationRStorage[payload.Relation.TargetType][payload.Relation.TargetID][payload.Relation.SourceType] = make(map[int]bool)
	}

	// now we store the relation
	RelationStorage[payload.Relation.SourceType][payload.Relation.SourceID][payload.Relation.TargetType][payload.Relation.TargetID] = payload.Relation

	// - - - - - - - - - - - - - - - - -
	// and an entry into the reverse index, its existence
	// allows us to use the coords in the normal index to revtrieve
	// the Relation. We dont create a pointer because golang doesnt
	// allow pointer on submaps in nested maps
	RelationRStorage[payload.Relation.TargetType][payload.Relation.TargetID][payload.Relation.SourceType][payload.Relation.SourceID] = true

	// we are done now we can unlock the entity Types
	//// - - - - - - - - - - - - - - - -
	//and finally unlock the relation Type and return
	//printMutexActions("CreateRelation.RelationStorageMutex.Unlock");
	RelationStorageMutex.Unlock()
}