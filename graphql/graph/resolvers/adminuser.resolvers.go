package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"fmt"
	"graphql/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, name string, email string) (*model.AdminUser, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, name *string, email *string) (*model.AdminUser, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.AdminUser, error) {
	panic(fmt.Errorf("not implemented: DeleteUser - deleteUser"))
}

// AdminUser is the resolver for the adminUser field.
func (r *queryResolver) AdminUser(ctx context.Context, id *string, name *string) (*model.AdminUser, error) {
	panic(fmt.Errorf("not implemented: AdminUser - adminUser"))
}

// AdminUsers is the resolver for the adminUsers field.
func (r *queryResolver) AdminUsers(ctx context.Context) ([]*model.AdminUser, error) {
	panic(fmt.Errorf("not implemented: AdminUsers - adminUsers"))
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) DeleteUser11(ctx context.Context, id string) (*model.AdminUser, error) {
	panic(fmt.Errorf("not implemented: DeleteUser11 - deleteUser11"))
}
*/