# Changelog    
## Release v0.1.2-alpha `04.03.2019`
* Update of README.md, [About Slingshot](https://github.com/voodooEntity/slingshotdb/blob/master/docs/ABOUT_SLINGSHOT.md) and [HTTP API](https://github.com/voodooEntity/slingshotdb/blob/master/docs/HTTP_API_V1.md)
* Adding 'context' filter filter param to multiple methods
    * `getChildEntities` (applies on the child entity not the relation) 
    * `getParentEntities` (applies on the parental entity not the relation)
    * `getRelationsTo` (applies on the relation)
    * `getRelationsFrom` (applies on the relation)
    * `getEntitiesByValue` 
* Fixing a possible concurrency bug in `storage.GetParentRelationsByTargetTypeAndTargetId`
* Little comment cleaning in storage

## Release v0.1.2-alpha `21.12.2019`    
* Update of README.md, [About Slingshot](https://github.com/voodooEntity/slingshotdb/blob/master/docs/ABOUT_SLINGSHOT.md) and [HTTP API](https://github.com/voodooEntity/slingshotdb/blob/master/docs/HTTP_API_V1.md)
* Adding optional param `mode` string to `getEntitiesByValue` and `getEntitiesByTypeAndValue` to decide wich compare mode should be used. Available now are
    * `match` (exact match, default mode)
    * `prefix` (value must begin with)
    * `suffix` (value must end with)
    * `contain` (value must contain)
    * `regex` (value must match regex pattern)
* Adding optional param `context` string to `getEntitiesByType` and `getEntitiesByTypeAndValue` api methods as filter.
* Some code comment cleaning in storage.go

## Release v0.1.1-alpha `19.12.2019`    
* Fixed a bug that could lead to concurrent map read&write actions
* Added Version to all datasets (Relations & Entities) that have to match on update to fix clientside race condition problems.
* Updated README.md
* Adding ABOUT_SLINGSHOT.md doc for better explaination of the database

## Release v0.1-alpha    
* Initial release of first alpha version