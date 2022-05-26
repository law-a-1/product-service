// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/law-a-1/product-service/ent/predicate"
	"github.com/law-a-1/product-service/ent/product"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeProduct = "Product"
)

// ProductMutation represents an operation that mutates the Product nodes in the graph.
type ProductMutation struct {
	config
	op            Op
	typ           string
	id            *int
	name          *string
	description   *string
	price         *int
	addprice      *int
	stock         *int
	addstock      *int
	image         *string
	video         *string
	created_at    *time.Time
	updated_at    *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Product, error)
	predicates    []predicate.Product
}

var _ ent.Mutation = (*ProductMutation)(nil)

// productOption allows management of the mutation configuration using functional options.
type productOption func(*ProductMutation)

// newProductMutation creates new mutation for the Product entity.
func newProductMutation(c config, op Op, opts ...productOption) *ProductMutation {
	m := &ProductMutation{
		config:        c,
		op:            op,
		typ:           TypeProduct,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withProductID sets the ID field of the mutation.
func withProductID(id int) productOption {
	return func(m *ProductMutation) {
		var (
			err   error
			once  sync.Once
			value *Product
		)
		m.oldValue = func(ctx context.Context) (*Product, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Product.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withProduct sets the old Product of the mutation.
func withProduct(node *Product) productOption {
	return func(m *ProductMutation) {
		m.oldValue = func(context.Context) (*Product, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ProductMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ProductMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ProductMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ProductMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Product.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *ProductMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ProductMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *ProductMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *ProductMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *ProductMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldDescription(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDescription is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDescription requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDescription: %w", err)
	}
	return oldValue.Description, nil
}

// ResetDescription resets all changes to the "description" field.
func (m *ProductMutation) ResetDescription() {
	m.description = nil
}

// SetPrice sets the "price" field.
func (m *ProductMutation) SetPrice(i int) {
	m.price = &i
	m.addprice = nil
}

// Price returns the value of the "price" field in the mutation.
func (m *ProductMutation) Price() (r int, exists bool) {
	v := m.price
	if v == nil {
		return
	}
	return *v, true
}

// OldPrice returns the old "price" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldPrice(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPrice is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPrice requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPrice: %w", err)
	}
	return oldValue.Price, nil
}

// AddPrice adds i to the "price" field.
func (m *ProductMutation) AddPrice(i int) {
	if m.addprice != nil {
		*m.addprice += i
	} else {
		m.addprice = &i
	}
}

// AddedPrice returns the value that was added to the "price" field in this mutation.
func (m *ProductMutation) AddedPrice() (r int, exists bool) {
	v := m.addprice
	if v == nil {
		return
	}
	return *v, true
}

// ResetPrice resets all changes to the "price" field.
func (m *ProductMutation) ResetPrice() {
	m.price = nil
	m.addprice = nil
}

// SetStock sets the "stock" field.
func (m *ProductMutation) SetStock(i int) {
	m.stock = &i
	m.addstock = nil
}

// Stock returns the value of the "stock" field in the mutation.
func (m *ProductMutation) Stock() (r int, exists bool) {
	v := m.stock
	if v == nil {
		return
	}
	return *v, true
}

// OldStock returns the old "stock" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldStock(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStock is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStock requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStock: %w", err)
	}
	return oldValue.Stock, nil
}

// AddStock adds i to the "stock" field.
func (m *ProductMutation) AddStock(i int) {
	if m.addstock != nil {
		*m.addstock += i
	} else {
		m.addstock = &i
	}
}

// AddedStock returns the value that was added to the "stock" field in this mutation.
func (m *ProductMutation) AddedStock() (r int, exists bool) {
	v := m.addstock
	if v == nil {
		return
	}
	return *v, true
}

// ResetStock resets all changes to the "stock" field.
func (m *ProductMutation) ResetStock() {
	m.stock = nil
	m.addstock = nil
}

// SetImage sets the "image" field.
func (m *ProductMutation) SetImage(s string) {
	m.image = &s
}

// Image returns the value of the "image" field in the mutation.
func (m *ProductMutation) Image() (r string, exists bool) {
	v := m.image
	if v == nil {
		return
	}
	return *v, true
}

// OldImage returns the old "image" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldImage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldImage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldImage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldImage: %w", err)
	}
	return oldValue.Image, nil
}

// ResetImage resets all changes to the "image" field.
func (m *ProductMutation) ResetImage() {
	m.image = nil
}

// SetVideo sets the "video" field.
func (m *ProductMutation) SetVideo(s string) {
	m.video = &s
}

// Video returns the value of the "video" field in the mutation.
func (m *ProductMutation) Video() (r string, exists bool) {
	v := m.video
	if v == nil {
		return
	}
	return *v, true
}

// OldVideo returns the old "video" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldVideo(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldVideo is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldVideo requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldVideo: %w", err)
	}
	return oldValue.Video, nil
}

// ResetVideo resets all changes to the "video" field.
func (m *ProductMutation) ResetVideo() {
	m.video = nil
}

// SetCreatedAt sets the "created_at" field.
func (m *ProductMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *ProductMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *ProductMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *ProductMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *ProductMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the Product entity.
// If the Product object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProductMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *ProductMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// Where appends a list predicates to the ProductMutation builder.
func (m *ProductMutation) Where(ps ...predicate.Product) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *ProductMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Product).
func (m *ProductMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ProductMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.name != nil {
		fields = append(fields, product.FieldName)
	}
	if m.description != nil {
		fields = append(fields, product.FieldDescription)
	}
	if m.price != nil {
		fields = append(fields, product.FieldPrice)
	}
	if m.stock != nil {
		fields = append(fields, product.FieldStock)
	}
	if m.image != nil {
		fields = append(fields, product.FieldImage)
	}
	if m.video != nil {
		fields = append(fields, product.FieldVideo)
	}
	if m.created_at != nil {
		fields = append(fields, product.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, product.FieldUpdatedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ProductMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case product.FieldName:
		return m.Name()
	case product.FieldDescription:
		return m.Description()
	case product.FieldPrice:
		return m.Price()
	case product.FieldStock:
		return m.Stock()
	case product.FieldImage:
		return m.Image()
	case product.FieldVideo:
		return m.Video()
	case product.FieldCreatedAt:
		return m.CreatedAt()
	case product.FieldUpdatedAt:
		return m.UpdatedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ProductMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case product.FieldName:
		return m.OldName(ctx)
	case product.FieldDescription:
		return m.OldDescription(ctx)
	case product.FieldPrice:
		return m.OldPrice(ctx)
	case product.FieldStock:
		return m.OldStock(ctx)
	case product.FieldImage:
		return m.OldImage(ctx)
	case product.FieldVideo:
		return m.OldVideo(ctx)
	case product.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case product.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Product field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ProductMutation) SetField(name string, value ent.Value) error {
	switch name {
	case product.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case product.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case product.FieldPrice:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPrice(v)
		return nil
	case product.FieldStock:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStock(v)
		return nil
	case product.FieldImage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetImage(v)
		return nil
	case product.FieldVideo:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVideo(v)
		return nil
	case product.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case product.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Product field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ProductMutation) AddedFields() []string {
	var fields []string
	if m.addprice != nil {
		fields = append(fields, product.FieldPrice)
	}
	if m.addstock != nil {
		fields = append(fields, product.FieldStock)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ProductMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case product.FieldPrice:
		return m.AddedPrice()
	case product.FieldStock:
		return m.AddedStock()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ProductMutation) AddField(name string, value ent.Value) error {
	switch name {
	case product.FieldPrice:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddPrice(v)
		return nil
	case product.FieldStock:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddStock(v)
		return nil
	}
	return fmt.Errorf("unknown Product numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ProductMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ProductMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ProductMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Product nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ProductMutation) ResetField(name string) error {
	switch name {
	case product.FieldName:
		m.ResetName()
		return nil
	case product.FieldDescription:
		m.ResetDescription()
		return nil
	case product.FieldPrice:
		m.ResetPrice()
		return nil
	case product.FieldStock:
		m.ResetStock()
		return nil
	case product.FieldImage:
		m.ResetImage()
		return nil
	case product.FieldVideo:
		m.ResetVideo()
		return nil
	case product.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case product.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	}
	return fmt.Errorf("unknown Product field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ProductMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ProductMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ProductMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ProductMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ProductMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ProductMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ProductMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Product unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ProductMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Product edge %s", name)
}