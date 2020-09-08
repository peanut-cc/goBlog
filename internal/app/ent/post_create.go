// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/peanut-cc/goBlog/internal/app/ent/category"
	"github.com/peanut-cc/goBlog/internal/app/ent/post"
	"github.com/peanut-cc/goBlog/internal/app/ent/tag"
)

// PostCreate is the builder for creating a Post entity.
type PostCreate struct {
	config
	mutation *PostMutation
	hooks    []Hook
}

// SetTitle sets the title field.
func (pc *PostCreate) SetTitle(s string) *PostCreate {
	pc.mutation.SetTitle(s)
	return pc
}

// SetBody sets the body field.
func (pc *PostCreate) SetBody(s string) *PostCreate {
	pc.mutation.SetBody(s)
	return pc
}

// SetCreatedTime sets the created_time field.
func (pc *PostCreate) SetCreatedTime(t time.Time) *PostCreate {
	pc.mutation.SetCreatedTime(t)
	return pc
}

// SetModifiedTime sets the modified_time field.
func (pc *PostCreate) SetModifiedTime(t time.Time) *PostCreate {
	pc.mutation.SetModifiedTime(t)
	return pc
}

// SetExcerpt sets the excerpt field.
func (pc *PostCreate) SetExcerpt(s string) *PostCreate {
	pc.mutation.SetExcerpt(s)
	return pc
}

// SetNillableExcerpt sets the excerpt field if the given value is not nil.
func (pc *PostCreate) SetNillableExcerpt(s *string) *PostCreate {
	if s != nil {
		pc.SetExcerpt(*s)
	}
	return pc
}

// SetAuthor sets the author field.
func (pc *PostCreate) SetAuthor(s string) *PostCreate {
	pc.mutation.SetAuthor(s)
	return pc
}

// SetIsDraft sets the is_Draft field.
func (pc *PostCreate) SetIsDraft(b bool) *PostCreate {
	pc.mutation.SetIsDraft(b)
	return pc
}

// SetCategoryID sets the category edge to Category by id.
func (pc *PostCreate) SetCategoryID(id int) *PostCreate {
	pc.mutation.SetCategoryID(id)
	return pc
}

// SetNillableCategoryID sets the category edge to Category by id if the given value is not nil.
func (pc *PostCreate) SetNillableCategoryID(id *int) *PostCreate {
	if id != nil {
		pc = pc.SetCategoryID(*id)
	}
	return pc
}

// SetCategory sets the category edge to Category.
func (pc *PostCreate) SetCategory(c *Category) *PostCreate {
	return pc.SetCategoryID(c.ID)
}

// AddTagIDs adds the tags edge to Tag by ids.
func (pc *PostCreate) AddTagIDs(ids ...int) *PostCreate {
	pc.mutation.AddTagIDs(ids...)
	return pc
}

// AddTags adds the tags edges to Tag.
func (pc *PostCreate) AddTags(t ...*Tag) *PostCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pc.AddTagIDs(ids...)
}

// Mutation returns the PostMutation object of the builder.
func (pc *PostCreate) Mutation() *PostMutation {
	return pc.mutation
}

// Save creates the Post in the database.
func (pc *PostCreate) Save(ctx context.Context) (*Post, error) {
	if err := pc.preSave(); err != nil {
		return nil, err
	}
	var (
		err  error
		node *Post
	)
	if len(pc.hooks) == 0 {
		node, err = pc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PostMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pc.mutation = mutation
			node, err = pc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pc.hooks) - 1; i >= 0; i-- {
			mut = pc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PostCreate) SaveX(ctx context.Context) *Post {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pc *PostCreate) preSave() error {
	if _, ok := pc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New("ent: missing required field \"title\"")}
	}
	if _, ok := pc.mutation.Body(); !ok {
		return &ValidationError{Name: "body", err: errors.New("ent: missing required field \"body\"")}
	}
	if _, ok := pc.mutation.CreatedTime(); !ok {
		return &ValidationError{Name: "created_time", err: errors.New("ent: missing required field \"created_time\"")}
	}
	if _, ok := pc.mutation.ModifiedTime(); !ok {
		return &ValidationError{Name: "modified_time", err: errors.New("ent: missing required field \"modified_time\"")}
	}
	if _, ok := pc.mutation.Author(); !ok {
		return &ValidationError{Name: "author", err: errors.New("ent: missing required field \"author\"")}
	}
	if _, ok := pc.mutation.IsDraft(); !ok {
		return &ValidationError{Name: "is_Draft", err: errors.New("ent: missing required field \"is_Draft\"")}
	}
	return nil
}

func (pc *PostCreate) sqlSave(ctx context.Context) (*Post, error) {
	po, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	po.ID = int(id)
	return po, nil
}

func (pc *PostCreate) createSpec() (*Post, *sqlgraph.CreateSpec) {
	var (
		po    = &Post{config: pc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: post.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: post.FieldID,
			},
		}
	)
	if value, ok := pc.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: post.FieldTitle,
		})
		po.Title = value
	}
	if value, ok := pc.mutation.Body(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: post.FieldBody,
		})
		po.Body = value
	}
	if value, ok := pc.mutation.CreatedTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: post.FieldCreatedTime,
		})
		po.CreatedTime = value
	}
	if value, ok := pc.mutation.ModifiedTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: post.FieldModifiedTime,
		})
		po.ModifiedTime = value
	}
	if value, ok := pc.mutation.Excerpt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: post.FieldExcerpt,
		})
		po.Excerpt = value
	}
	if value, ok := pc.mutation.Author(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: post.FieldAuthor,
		})
		po.Author = value
	}
	if value, ok := pc.mutation.IsDraft(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: post.FieldIsDraft,
		})
		po.IsDraft = value
	}
	if nodes := pc.mutation.CategoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   post.CategoryTable,
			Columns: []string{post.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: category.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   post.TagsTable,
			Columns: post.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return po, _spec
}

// PostCreateBulk is the builder for creating a bulk of Post entities.
type PostCreateBulk struct {
	config
	builders []*PostCreate
}

// Save creates the Post entities in the database.
func (pcb *PostCreateBulk) Save(ctx context.Context) ([]*Post, error) {
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Post, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				if err := builder.preSave(); err != nil {
					return nil, err
				}
				mutation, ok := m.(*PostMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				id := specs[i].ID.Value.(int64)
				nodes[i].ID = int(id)
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX calls Save and panics if Save returns an error.
func (pcb *PostCreateBulk) SaveX(ctx context.Context) []*Post {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
