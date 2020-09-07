package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// Blog holds the schema definition for the Blog entity.
type Blog struct {
	ent.Schema
}

// Fields of the Blog.
func (Blog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("default_page_num"),
		field.String("blog_name"),
		field.String("btitle"),
		field.String("subtitle"),
		field.String("beian"),
		field.String("copy_right"),
	}
}

// Edges of the Blog.
func (Blog) Edges() []ent.Edge {
	return nil
}
