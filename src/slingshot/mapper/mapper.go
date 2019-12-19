package mapper

// handle all the imports
import (
	"encoding/json"
	"errors"
	"slingshot/storage"
	"slingshot/types"
)

// - - - - - - - - - - - - - - - - - - - - - - - - - -
//                  service interfaces
// - - - - - - - - - - - - - - - - - - - - - - - - - -
func GetEntityByTypeAndId(entityType string, id int) (types.MapperTransport, error) {
	// check if the entity type exists in first place
	ok := storage.TypeExists(entityType)
	if !ok {
		return types.MapperTransport{}, errors.New("Unknown entity type given")
	}

	// get the entity
	entityTypeID, _ := storage.GetTypeIdByString(entityType)
	entity, err := storage.GetEntityByPath(entityTypeID, id)
	if nil != err {
		return types.MapperTransport{}, errors.New("Unknown entity id given")
	}

	// prepare the return entry
	entityTypeString, _ := storage.GetTypeStringById(entity.Type)
	var tmpArray = []types.MapperEntity{
		types.MapperEntity{
			Type:       *entityTypeString,
			ID:         entity.ID,
			Properties: entity.Properties,
			Context:    entity.Context,
			Value:      entity.Value,
			Version:    entity.Version,
			Children:   []types.MapperEntity{},
		},
	}

	// build up the return
	transport := types.MapperTransport{
		Entities: tmpArray,
	}
	return transport, nil
}

func GetEntitiesRecursive(entityType int, entityID int, depth int, relations *[]types.MapperRelation) (types.MapperEntity, error) {
	// lets  get the entity
	entity, err := storage.GetEntityByPath(entityType, entityID)

	// if it doesnt exist we stop
	if err != nil {
		return types.MapperEntity{}, errors.New("Entity with give path does not exist")
	}

	// now we define the data into
	// mapper entity struct
	entityTypeString, _ := storage.GetTypeStringById(entity.Type)
	returnEntity := types.MapperEntity{
		ID:         entity.ID,
		Type:       *entityTypeString,
		Value:      entity.Value,
		Context:    entity.Context,
		Properties: entity.Properties,
		Version:    entity.Version,
	}

	// if we didnt reach the destinated depth yet
	// we retrieve all relations from this entity
	if 0 < depth {
		depth--
		tmpRelations, _ := storage.GetChildRelationsBySourceTypeAndSourceId(entity.Type, entity.ID)

		// lets check if we got anyrelations
		if 0 != len(tmpRelations) {
			returnEntity.Children = []types.MapperEntity{}
			// there are children lets iterate
			// through the map
			for _, tmpRelation := range tmpRelations {
				// prepare and add the reltion to the relations array
				sourceType, _ := storage.GetTypeStringById(tmpRelation.SourceType)
				targetType, _ := storage.GetTypeStringById(tmpRelation.TargetType)
				transRelation := types.MapperRelation{
					SourceID:   tmpRelation.SourceID,
					SourceType: *sourceType,
					TargetID:   tmpRelation.TargetID,
					TargetType: *targetType,
					Properties: tmpRelation.Properties,
					Context:    tmpRelation.Context,
					Version:    tmpRelation.Version,
				}
				*relations = append(*relations, transRelation)

				// call the function recursive and add the object
				var tmpEntity, _ = GetEntitiesRecursive(tmpRelation.TargetType,
					tmpRelation.TargetID,
					depth,
					relations)

				// store the subentity in child field
				returnEntity.Children = append(returnEntity.Children, tmpEntity)
			}
		}
	}
	// fill children with empty array
	if nil == returnEntity.Children {
		returnEntity.Children = []types.MapperEntity{}
	}

	// return the entity
	return returnEntity, nil
}

func DeleteEntity(Type string, id int) error {
	typeID, err := storage.GetTypeIdByString(Type)
	if nil != err {
		return err
	}
	storage.DeleteEntity(typeID, id)
	return nil
}

func UpdateEntity(entityType string, entityID int, value string, properties map[string]string, context string, Version int) error {
	// first we create the entity type or retrieve the id (smart function is smart)
	entityTypeID, _ := storage.CreateEntityType(entityType)

	// now we create the new entity for faster updating
	newEntity := types.StorageEntity{
		ID:         entityID,
		Type:       entityTypeID,
		Value:      value,
		Properties: properties,
		Context:    context,
		Version:    Version,
	}
	// now we pass the new entity object to
	// the storage
	err := storage.UpdateEntity(newEntity)
	// pass er through no matter if its error or nil
	return err
}

func CreateEntity(entityType string, value string, properties map[string]string, context string) (types.MapperTransport, error) {
	// check if the entity type exists in first place
	// if not we auto create it
	entityTypeID, _ := storage.CreateEntityType(entityType)

	// build the entity object
	entity := types.StorageEntity{
		Type:       entityTypeID,
		Value:      value,
		Properties: properties,
		Context:    context,
	}

	// finally create the entity
	id, err := storage.CreateEntity(entity)
	if err != nil {
		return types.MapperTransport{}, err
	}

	// prepare the return entry
	var tmpArray = []types.MapperEntity{
		types.MapperEntity{
			Type:       entityType,
			ID:         id,
			Properties: properties,
			Context:    context,
			Value:      value,
			Version:    1,
			Children:   []types.MapperEntity{},
		},
	}

	// build up the return
	transport := types.MapperTransport{
		Entities:  tmpArray,
		Relations: []types.MapperRelation{},
	}

	return transport, nil
}

func GetEntitiesByType(entityType string) (types.MapperTransport, error) {
	// check if the entity type exists in first place
	ok := storage.TypeExists(entityType)
	if !ok {
		return types.MapperTransport{}, errors.New("Unknown entity type given")
	}

	// ok type exists lets retrieve all fitting entities for
	// this type
	entities, _ := storage.GetEntitiesByType(entityType)
	var tmpArray []types.MapperEntity
	// if we got more than 0 entries we transform it
	if 0 < len(entities) {
		for id := range entities {
			tmpArray = append(tmpArray, types.MapperEntity{
				ID:         entities[id].ID,
				Value:      entities[id].Value,
				Type:       entityType,
				Context:    entities[id].Context,
				Properties: entities[id].Properties,
				Children:   []types.MapperEntity{},
				Version:    entities[id].Version,
			})
		}
	}
	// build up the return
	transport := types.MapperTransport{
		Entities:  tmpArray,
		Relations: []types.MapperRelation{},
	}

	return transport, nil
}

func GetEntitiesByTypeAndValue(entityType string, entityValue string) (types.MapperTransport, error) {
	// check if the entity type exists in first place
	ok := storage.TypeExists(entityType)
	if !ok {
		return types.MapperTransport{}, errors.New("Unknown entity type given")
	}

	// ok type exists lets retrieve all fitting entities for
	// this type
	entities, _ := storage.GetEntitiesByType(entityType)
	var tmpArray []types.MapperEntity
	if 0 < len(entities) {
		for id := range entities {
			if entities[id].Value == entityValue {
				tmpArray = append(tmpArray, types.MapperEntity{
					ID:         entities[id].ID,
					Value:      entities[id].Value,
					Type:       entityType,
					Properties: entities[id].Properties,
					Context:    entities[id].Context,
					Version:    entities[id].Version,
					Children:   []types.MapperEntity{},
				})
			}
		}
	}
	// build up the return
	transport := types.MapperTransport{
		Entities:  tmpArray,
		Relations: []types.MapperRelation{},
	}

	return transport, nil
}

func MapJson(data []byte) (types.MapperTransport, error) {
	// define an entity struct as template
	// for json unmarshal and parse the
	// json byte array into the var
	var entity types.MapperEntity
	if err := json.Unmarshal(data, &entity); err != nil {
		return types.MapperTransport{}, errors.New("Invalid Json")
	}

	// lets start recursive mapping of the
	// unmarshalled data
	newID, err := MapEntitiesRecursive(entity, -1, -1)
	if err != nil {
		return types.MapperTransport{}, errors.New("Couldnt map entities .... why tho?")
	}

	// we got it done lets wrap our data
	// in an transport object
	entities := []types.MapperEntity{}
	entities = append(entities, types.MapperEntity{
		ID:         newID,
		Type:       entity.Type,
		Value:      entity.Value,
		Properties: entity.Properties,
		Context:    entity.Context,
		Version:    1,
		Children:   []types.MapperEntity{},
	})

	transport := types.MapperTransport{
		Entities:  entities,
		Relations: []types.MapperRelation{},
	}

	return transport, nil
}

func MapEntitiesRecursive(entity types.MapperEntity, parentType int, parentID int) (int, error) {
	// first we get the right TypeID
	var TypeID = HandleType(entity.Type)
	// now we create the fitting entity
	tmpEntity := types.StorageEntity{
		ID:         -1,
		Type:       TypeID,
		Value:      entity.Value,
		Context:    entity.Context,
		Version:    1,
		Properties: entity.Properties,
	}
	// now we create the entity
	var newID, _ = storage.CreateEntity(tmpEntity)
	// lets check if there are child elements
	if len(entity.Children) != 0 {
		// there are children lets iteater over
		// the map
		for _, childEntity := range entity.Children {
			// pas the child entity and the parent coords to
			// create the relation after inserting the entity
			MapEntitiesRecursive(childEntity, TypeID, newID)
		}
	}
	// now lets check if ourparent Type and id
	// are not -1 , if so we need to create
	// a relation
	if parentType != -1 && parentID != -1 {
		// lets create the relation to our parent
		tmpRelation := types.StorageRelation{
			SourceType: parentType,
			SourceID:   parentID,
			TargetType: tmpEntity.Type,
			TargetID:   newID,
			Version:    1,
		}
		storage.CreateRelation(parentType, parentID, tmpEntity.Type, newID, tmpRelation)
	}
	// only the first return is interesting since it
	// returns the most parent id
	return newID, nil
}

func GetChildEntities(strType string, id int) (types.MapperTransport, error) {
	// first we check if the entity type exists
	Type, err := storage.GetTypeIdByString(strType)
	if err != nil {
		return types.MapperTransport{}, err
	}

	// now we check if the entity exists
	exists := storage.EntityExists(Type, id)
	if false == exists {
		return types.MapperTransport{}, nil
	}

	// ok the entity seems to exist lets retrieve all the child relations
	relations, err := storage.GetChildRelationsBySourceTypeAndSourceId(Type, id)
	if 0 == len(relations) {
		// if we dont find any we return an empty Transport
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok seems like we actually found child entities, lets retrieve
	// them
	transport := types.MapperTransport{
		Entities:  []types.MapperEntity{},
		Relations: []types.MapperRelation{},
	}

	for id := range relations {
		// lets retrieve the entity
		entity, err := storage.GetEntityByPath(relations[id].TargetType, relations[id].TargetID)
		if nil != err {
			return types.MapperTransport{}, err
		}

		// make the type readable and add the entity to our transport struct
		strType, _ := storage.GetTypeStringById(entity.Type)
		tmpEntity := types.MapperEntity{
			ID:         entity.ID,
			Type:       *strType,
			Value:      entity.Value,
			Properties: entity.Properties,
			Context:    entity.Context,
			Version:    entity.Version,
			Children:   []types.MapperEntity{},
		}
		transport.Entities = append(transport.Entities, tmpEntity)
	}

	return transport, nil
}

func GetParentEntities(strType string, id int) (types.MapperTransport, error) {
	// first we check if the entity type exists
	Type, err := storage.GetTypeIdByString(strType)
	if err != nil {
		return types.MapperTransport{}, err
	}

	// now we check if the entity exists
	exists := storage.EntityExists(Type, id)
	if false == exists {
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok the entity seems to exist lets retrieve all the child relations
	relations, err := storage.GetParentRelationsByTargetTypeAndTargetId(Type, id)
	if 0 == len(relations) {
		// if we dont find any we return an empty Transport
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok seems like we actually found child entities, lets retrieve
	// them
	transport := types.MapperTransport{
		Entities:  []types.MapperEntity{},
		Relations: []types.MapperRelation{},
	}

	for id := range relations {
		// lets retrieve the entity
		entity, err := storage.GetEntityByPath(relations[id].SourceType, relations[id].SourceType)
		if nil != err {
			return types.MapperTransport{}, err
		}

		// make the type readable and add the entity to our transport struct
		strType, _ := storage.GetTypeStringById(entity.Type)
		tmpEntity := types.MapperEntity{
			ID:         entity.ID,
			Type:       *strType,
			Value:      entity.Value,
			Properties: entity.Properties,
			Context:    entity.Context,
			Version:    entity.Version,
			Children:   []types.MapperEntity{},
		}
		transport.Entities = append(transport.Entities, tmpEntity)
	}

	return transport, nil
}

func GetRelation(srcType string, srcID int, targetType string, targetID int) (types.MapperTransport, error) {
	// prepare the data for storage
	srcTypeID, _ := storage.GetTypeIdByString(srcType)
	targetTypeID, _ := storage.GetTypeIdByString(targetType)

	// get the relation
	relation, err := storage.GetRelation(srcTypeID, srcID, targetTypeID, targetID)
	if nil != err {
		return types.MapperTransport{}, err
	}

	// prepare return transport
	transport := types.MapperTransport{
		Entities:  []types.MapperEntity{},
		Relations: []types.MapperRelation{},
	}

	// transform the relation dataset
	sourceTypeString, _ := storage.GetTypeStringById(relation.SourceType)
	targetTypeString, _ := storage.GetTypeStringById(relation.TargetType)
	tmpRelation := types.MapperRelation{
		SourceType: *sourceTypeString,
		SourceID:   relation.SourceID,
		TargetType: *targetTypeString,
		TargetID:   relation.TargetID,
		Context:    relation.Context,
		Properties: relation.Properties,
		Version:    relation.Version,
	}

	// add relation to transport and return
	transport.Relations = append(transport.Relations, tmpRelation)
	return transport, nil
}

func UpdateRelation(srcType string, srcID int, targetType string, targetID int, relation types.MapperRelation) error {
	// prepare the data for storage
	srcTypeID, _ := storage.GetTypeIdByString(srcType)
	targetTypeID, _ := storage.GetTypeIdByString(targetType)
	tmpRelation := types.StorageRelation{
		SourceType: srcTypeID,
		SourceID:   srcID,
		TargetType: targetTypeID,
		TargetID:   targetID,
		Properties: relation.Properties,
		Context:    relation.Context,
		Version:    relation.Version,
	}

	// update the relation
	_, err := storage.UpdateRelation(srcTypeID, srcID, targetTypeID, targetID, tmpRelation)
	if nil != err {
		return err
	}

	return nil
}

func CreateRelation(srcType string, srcID int, targetType string, targetID int, relation types.MapperRelation) error {
	// get the IDs for storage handling
	srcTypeID, _ := storage.GetTypeIdByString(srcType)
	targetTypeID, _ := storage.GetTypeIdByString(targetType)

	// check if the entities in first place todo make sure noone can explaint the small timing window between checking
	// the existence and actually updating the relation. ^^
	source := storage.EntityExists(srcTypeID, srcID)
	target := storage.EntityExists(targetTypeID, targetID)
	if !source || !target {
		return errors.New("Either source or target entity are not existing")
	}

	// ok we seem to be fine, lets update the relation
	tmpRelation := types.StorageRelation{
		SourceType: srcTypeID,
		SourceID:   srcID,
		TargetType: targetTypeID,
		TargetID:   targetID,
		Properties: relation.Properties,
		Context:    relation.Context,
	}

	// update the relation
	_, err := storage.CreateRelation(srcTypeID, srcID, targetTypeID, targetID, tmpRelation)
	if nil != err {
		return err
	}

	return nil
}

func DeleteRelation(srcType string, srcID int, targetType string, targetID int) error {
	// get the IDs for storage handling
	srcTypeID, _ := storage.GetTypeIdByString(srcType)
	targetTypeID, _ := storage.GetTypeIdByString(targetType)

	// check if that relation exists, todo make sure the timing doesnt fuck up
	// between existence check and actual deletion
	link := storage.RelationExists(srcTypeID, srcID, targetTypeID, targetID)
	if !link {
		return errors.New("Cant delete non existing relation")
	}

	// delete the relation todo fix safety between check and deletion.....
	storage.DeleteRelation(srcTypeID, srcID, targetTypeID, targetID)

	return nil
}

func GetRelationsTo(strType string, id int) (types.MapperTransport, error) {
	// first we check if the entity type exists
	Type, err := storage.GetTypeIdByString(strType)
	if err != nil {
		return types.MapperTransport{}, err
	}

	// now we check if the entity exists
	exists := storage.EntityExists(Type, id)
	if false == exists {
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok the entity seems to exist lets retrieve all the child relations
	relations, err := storage.GetParentRelationsByTargetTypeAndTargetId(Type, id)
	if 0 == len(relations) {
		// if we dont find any we return an empty Transport
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok lets build the final return
	transport := types.MapperTransport{
		Entities:  []types.MapperEntity{},
		Relations: []types.MapperRelation{},
	}
	for _, relation := range relations {
		tmpSrcType, _ := storage.GetTypeStringById(relation.SourceType)
		tmpTargetType, _ := storage.GetTypeStringById(relation.TargetType)
		tmpRelation := types.MapperRelation{
			SourceType: *tmpSrcType,
			SourceID:   relation.SourceID,
			TargetType: *tmpTargetType,
			TargetID:   relation.TargetID,
			Context:    relation.Context,
			Properties: relation.Properties,
			Version:    relation.Version,
		}
		transport.Relations = append(transport.Relations, tmpRelation)
	}

	return transport, nil
}

func GetRelationsFrom(strType string, id int) (types.MapperTransport, error) {
	// first we check if the entity type exists
	Type, err := storage.GetTypeIdByString(strType)
	if err != nil {
		return types.MapperTransport{}, err
	}

	// now we check if the entity exists
	exists := storage.EntityExists(Type, id)
	if false == exists {
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok the entity seems to exist lets retrieve all the child relations
	relations, err := storage.GetChildRelationsBySourceTypeAndSourceId(Type, id)
	if 0 == len(relations) {
		// if we dont find any we return an empty Transport
		return types.MapperTransport{
			Entities:  []types.MapperEntity{},
			Relations: []types.MapperRelation{},
		}, nil
	}

	// ok lets build the final return
	transport := types.MapperTransport{
		Entities:  []types.MapperEntity{},
		Relations: []types.MapperRelation{},
	}
	for _, relation := range relations {
		tmpSrcType, _ := storage.GetTypeStringById(relation.SourceType)
		tmpTargetType, _ := storage.GetTypeStringById(relation.TargetType)
		tmpRelation := types.MapperRelation{
			SourceType: *tmpSrcType,
			SourceID:   relation.SourceID,
			TargetType: *tmpTargetType,
			TargetID:   relation.TargetID,
			Context:    relation.Context,
			Properties: relation.Properties,
			Version:    relation.Version,
		}
		transport.Relations = append(transport.Relations, tmpRelation)
	}

	return transport, nil
}

func GetEntitiesByValue(value string) (types.MapperTransport, error) {
	// first we make sure we dont search for empty
	if "" == value {
		return types.MapperTransport{}, errors.New("Dont search for empty value ~.~")
	}

	// ok seems fine lets search for entties by value
	entities := storage.GetEntitiesByValue(value)

	// prepare return
	transport := types.MapperTransport{
		Entities:  []types.MapperEntity{},
		Relations: []types.MapperRelation{},
	}

	// now if we have any hits, we transform them into
	// our transport return
	if 0 < len(entities) {
		for _, entity := range entities {
			tmpType, _ := storage.GetTypeStringById(entity.Type)
			tmpEntity := types.MapperEntity{
				ID:         entity.ID,
				Type:       *tmpType,
				Value:      entity.Value,
				Context:    entity.Context,
				Properties: entity.Properties,
				Version:    entity.Version,
				Children:   []types.MapperEntity{},
			}
			transport.Entities = append(transport.Entities, tmpEntity)
		}
	}

	return transport, nil
}

func HandleType(strType string) int {
	// get the Type id by string
	id, _ := storage.GetTypeIdByString(strType)
	if -1 != id {
		// it does lets return it
		return id
	}
	// it didnt exist so we create it and return
	// the new id
	newID, _ := storage.CreateEntityType(strType)
	return newID
}

func GetEntityTypes() []string {
	return storage.GetEntityTypes()
}
