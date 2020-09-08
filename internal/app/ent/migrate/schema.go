// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebook/ent/dialect/sql/schema"
	"github.com/facebook/ent/schema/field"
)

var (
	// BlogsColumns holds the columns for the "blogs" table.
	BlogsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "default_page_num", Type: field.TypeInt},
		{Name: "blog_name", Type: field.TypeString},
		{Name: "btitle", Type: field.TypeString},
		{Name: "subtitle", Type: field.TypeString},
		{Name: "beian", Type: field.TypeString},
		{Name: "copy_right", Type: field.TypeString},
	}
	// BlogsTable holds the schema information for the "blogs" table.
	BlogsTable = &schema.Table{
		Name:        "blogs",
		Columns:     BlogsColumns,
		PrimaryKey:  []*schema.Column{BlogsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// CategoriesColumns holds the columns for the "categories" table.
	CategoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// CategoriesTable holds the schema information for the "categories" table.
	CategoriesTable = &schema.Table{
		Name:        "categories",
		Columns:     CategoriesColumns,
		PrimaryKey:  []*schema.Column{CategoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// PostsColumns holds the columns for the "posts" table.
	PostsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "body", Type: field.TypeString},
		{Name: "created_time", Type: field.TypeTime},
		{Name: "modified_time", Type: field.TypeTime},
		{Name: "excerpt", Type: field.TypeString, Nullable: true},
		{Name: "author", Type: field.TypeString},
		{Name: "is_draft", Type: field.TypeBool},
		{Name: "category_posts", Type: field.TypeInt, Nullable: true},
	}
	// PostsTable holds the schema information for the "posts" table.
	PostsTable = &schema.Table{
		Name:       "posts",
		Columns:    PostsColumns,
		PrimaryKey: []*schema.Column{PostsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "posts_categories_posts",
				Columns: []*schema.Column{PostsColumns[8]},

				RefColumns: []*schema.Column{CategoriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// TagsColumns holds the columns for the "tags" table.
	TagsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// TagsTable holds the schema information for the "tags" table.
	TagsTable = &schema.Table{
		Name:        "tags",
		Columns:     TagsColumns,
		PrimaryKey:  []*schema.Column{TagsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "username", Type: field.TypeString, Unique: true},
		{Name: "password", Type: field.TypeString},
		{Name: "token", Type: field.TypeString, Nullable: true},
		{Name: "email", Type: field.TypeString},
		{Name: "phone", Type: field.TypeString},
		{Name: "login_time", Type: field.TypeTime, Nullable: true},
		{Name: "logout_time", Type: field.TypeTime, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:        "users",
		Columns:     UsersColumns,
		PrimaryKey:  []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// TagPostsColumns holds the columns for the "tag_posts" table.
	TagPostsColumns = []*schema.Column{
		{Name: "tag_id", Type: field.TypeInt},
		{Name: "post_id", Type: field.TypeInt},
	}
	// TagPostsTable holds the schema information for the "tag_posts" table.
	TagPostsTable = &schema.Table{
		Name:       "tag_posts",
		Columns:    TagPostsColumns,
		PrimaryKey: []*schema.Column{TagPostsColumns[0], TagPostsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "tag_posts_tag_id",
				Columns: []*schema.Column{TagPostsColumns[0]},

				RefColumns: []*schema.Column{TagsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:  "tag_posts_post_id",
				Columns: []*schema.Column{TagPostsColumns[1]},

				RefColumns: []*schema.Column{PostsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		BlogsTable,
		CategoriesTable,
		PostsTable,
		TagsTable,
		UsersTable,
		TagPostsTable,
	}
)

func init() {
	PostsTable.ForeignKeys[0].RefTable = CategoriesTable
	TagPostsTable.ForeignKeys[0].RefTable = TagsTable
	TagPostsTable.ForeignKeys[1].RefTable = PostsTable
}
