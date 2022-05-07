package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Validator template
var templateSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"name"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "description",
			},
		},
	},
}

// List of required collections
var requiredCollectionsAndSchemas = map[string]bson.M{
	"packages":         packageSchema,
	"packagesMetadata": packageMetadataSchema,
	"users":            userSchema,
	"authentications":  authenticationSchema,
	"sources":          sourceSchema,
	"projects":         projectSchema,
	"scripts":          scriptSchema,
	"graphs":           graphSchema,
	"admins":           adminSchema,
}

// Validation for 'graphs' collections
var graphSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"name", "version", "package", "children"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "required string name of node",
			},
			"version": bson.M{
				"bsonType":    "long",
				"description": "required long version of node",
			},
			"package": bson.M{
				"bsonType":    "objectId",
				"description": "required object id for associated package",
			},
			"children": bson.M{
				"bsonType":    "array",
				"description": "required array of children of node",
				"items": bson.M{
					"bsonType": "objectId",
				},
			},
		},
	},
}

// Validation for 'packages' collection
var packageSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"name", "version", "sourceId", "graphId", "scripts"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "required string name of package",
			},
			"version": bson.M{
				"bsonType":    "long",
				"description": "required long version of package",
			},
			"sourceId": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of package source",
			},
			"graphId": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of package graph",
			},
			"scripts": bson.M{
				"bsonType":    "array",
				"description": "required array of object ids for package scripts",
				"items": bson.M{
					"bsonType": "objectId",
				},
			},
		},
	},
}

// Validation for 'packagesMetadata' collection
var packageMetadataSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"_id", "packageId", "numberDownloads", "lastDownloaded", "lastDownloadedBy"},
		"properties": bson.M{
			"_id": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of metadata",
			},
			"packageId": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of package",
			},
			"numberDownloads": bson.M{
				"bsonType":    "long",
				"description": "required long of package downloads",
			},
			"lastDownloaded": bson.M{
				"bsonType":    "date",
				"description": "required date of last package download",
			},
			"lastDownloadedBy": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of last downloading user",
			},
		},
	},
}

// Validation for 'users' collection
var userSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"name", "email", "authenticationId"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "required string name of user",
			},
			"email": bson.M{
				"bsonType":    "string",
				"description": "required string email address of user",
			},
			"authenticationId": bson.M{
				"bsonType":    primitive.ObjectID{},
				"description": "required object id for authentication attached to user",
			},
		},
	},
}

// Validation for 'authentications' collection
var authenticationSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"userId", "projectIds", "isAdmin"},
		"properties": bson.M{
			"userId": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of user",
			},
			"projectIds": bson.M{
				"bsonType":    "array",
				"description": "required array of object ids for projects",
				"items": bson.M{
					"bsonType": "objectId",
				},
			},
			"isAdmin": bson.M{
				"bsonType":    "bool",
				"description": "required bool for is user admin",
			},
		},
	},
}

// Validation for 'sources' collection
var sourceSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"path", "lastAccessedBy"},
		"properties": bson.M{
			"path": bson.M{
				"bsonType":    "string",
				"description": "required string path to source",
			},
			"lastAccessedBy": bson.M{
				"bsonType":    "date",
				"description": "required date source was last accessed",
			},
		},
	},
}

// Validation for 'projects' collection
var projectSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"_id", "name", "sourceId", "projectVersion", "packageVersion", "graphId", "packages"},
		"properties": bson.M{
			"_id": bson.M{
				"bsonType":    primitive.ObjectID{},
				"description": "required object id of project",
			},
			"name": bson.M{
				"bsonType":    "string",
				"description": "required string name of project",
			},
			"sourceId": bson.M{
				"bsonType":    "objectId",
				"description": "required object id of project source",
			},
			"projectVersion": bson.M{
				"bsonType":    "long",
				"description": "required long for project version",
			},
			"packageVersion": bson.M{
				"bsonType":    "long",
				"description": "required long for attached package version",
			},
			"graphId": bson.M{
				"bsonType":    "objectId",
				"description": "required object id for project graph",
			},
			"packages": bson.M{
				"bsonType":    "objectId",
				"description": "required object id for project packages",
				"items": bson.M{
					"bsonType": "objectId",
				},
			},
		},
	},
}

// Validation for 'scripts' collection
var scriptSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"name", "description", "body", "package"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "required string name of script",
			},
			"description": bson.M{
				"bsonType":    "string",
				"description": "required string description of script",
			},
			"body": bson.M{
				"bsonType":    "string",
				"description": "required string body of script to be executed",
			},
			"package": bson.M{
				"bsonType":    "objectId",
				"description": "required object id for associated package",
			},
			"packages": bson.M{
				"bsonType":    "array",
				"description": "optional object ids for associated packages",
				"items": bson.M{
					"bsonType": "objectId",
				},
			},
		},
	},
}

// Validation for 'admins' collection
var adminSchema = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"userId", "password"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "objectId",
				"description": "required object id for user",
			},
			"password": bson.M{
				"bsonType":    "string",
				"description": "required string password for user",
			},
		},
	},
}
