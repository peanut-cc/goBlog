package schema

import (
	"time"

	"github.com/facebook/ent/dialect"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Unique(),
		field.String("body").SchemaType(map[string]string{
			dialect.MySQL: "longtext",
		}),
		field.Time("created_time").Default(time.Now),
		field.Time("modified_time").Default(time.Now),
		field.String("excerpt").Optional(),
		field.String("author"),
		// 是否是草稿
		field.Bool("is_Draft").Default(false),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).
			Ref("posts").
			Unique(),
		edge.From("tags", Tag.Type).Ref("posts"),
	}
}
