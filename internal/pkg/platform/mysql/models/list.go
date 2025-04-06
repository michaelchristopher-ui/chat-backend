package models

/*
	ModelsList represent the gorm tagged structs that define the tables.
	Insert pointers to gorm tagged structs here.
*/
var ModelsList = []interface{}{
	&UserFriends{},
	&Messages{},
	&Account{},
}
