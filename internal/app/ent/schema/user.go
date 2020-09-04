package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username"),
		field.String("password"),
		field.String("token"),
		field.String("email"),
		field.String("phone"),
		field.Time("login_time"),
		field.Time("logout_time"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
