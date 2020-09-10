// Code generated by entc, DO NOT EDIT.

package post

import (
	"time"
)

const (
	// Label holds the string label denoting the post type in the database.
	Label = "post"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldBody holds the string denoting the body field in the database.
	FieldBody = "body"
	// FieldCreatedTime holds the string denoting the created_time field in the database.
	FieldCreatedTime = "created_time"
	// FieldModifiedTime holds the string denoting the modified_time field in the database.
	FieldModifiedTime = "modified_time"
	// FieldExcerpt holds the string denoting the excerpt field in the database.
	FieldExcerpt = "excerpt"
	// FieldAuthor holds the string denoting the author field in the database.
	FieldAuthor = "author"
	// FieldIsDraft holds the string denoting the is_draft field in the database.
	FieldIsDraft = "is_draft"

	// EdgeCategory holds the string denoting the category edge name in mutations.
	EdgeCategory = "category"
	// EdgeTags holds the string denoting the tags edge name in mutations.
	EdgeTags = "tags"

	// Table holds the table name of the post in the database.
	Table = "posts"
	// CategoryTable is the table the holds the category relation/edge.
	CategoryTable = "posts"
	// CategoryInverseTable is the table name for the Category entity.
	// It exists in this package in order to avoid circular dependency with the "category" package.
	CategoryInverseTable = "categories"
	// CategoryColumn is the table column denoting the category relation/edge.
	CategoryColumn = "category_posts"
	// TagsTable is the table the holds the tags relation/edge. The primary key declared below.
	TagsTable = "tag_posts"
	// TagsInverseTable is the table name for the Tag entity.
	// It exists in this package in order to avoid circular dependency with the "tag" package.
	TagsInverseTable = "tags"
)

// Columns holds all SQL columns for post fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldBody,
	FieldCreatedTime,
	FieldModifiedTime,
	FieldExcerpt,
	FieldAuthor,
	FieldIsDraft,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Post type.
var ForeignKeys = []string{
	"category_posts",
}

var (
	// TagsPrimaryKey and TagsColumn2 are the table columns denoting the
	// primary key for the tags relation (M2M).
	TagsPrimaryKey = []string{"tag_id", "post_id"}
)

var (
	// DefaultCreatedTime holds the default value on creation for the created_time field.
	DefaultCreatedTime func() time.Time
	// DefaultModifiedTime holds the default value on creation for the modified_time field.
	DefaultModifiedTime func() time.Time
)
