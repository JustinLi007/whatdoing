package database

import (
	"service-user/internal/password"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersNoUsername(t *testing.T) {
	db, err := NewDb("")
	require.NoError(t, err)
	usersService := NewServiceUsers(db)

	// Create first user
	reqUser1 := &User{
		Email:    "nousername@test.com",
		Password: &password.Password{},
	}
	reqUser1.Password.Set("1234")
	newUser1, err := usersService.CreateUser(reqUser1)
	require.NoError(t, err)
	require.NotNil(t, newUser1)
	assert.Nil(t, newUser1.Username)
	assert.Equal(t, reqUser1.Email, newUser1.Email)
	assert.Equal(t, "regular", newUser1.Role)

	// Create second user
	reqUser2 := &User{
		Email:    "nousername2@test.com",
		Password: &password.Password{},
	}
	reqUser2.Password.Set("1234")
	newUser2, err := usersService.CreateUser(reqUser2)
	require.NoError(t, err)
	require.NotNil(t, newUser2)
	assert.Nil(t, newUser2.Username)
	assert.Equal(t, reqUser2.Email, newUser2.Email)
	assert.Equal(t, "regular", newUser2.Role)

	// Create second user again
	reqUser2Copy := &User{
		Email:    "nousername2@test.com",
		Password: &password.Password{},
	}
	reqUser2Copy.Password.Set("1234")
	newUser2Copy, err := usersService.CreateUser(reqUser2Copy)
	require.Error(t, err)
	require.Nil(t, newUser2Copy)

	// Get user1 by id
	reqGetUser1ById := &User{
		Id: newUser1.Id,
	}
	existingUser1, err := usersService.GetUserById(reqGetUser1ById)
	require.NoError(t, err)
	require.NotNil(t, existingUser1)
	assert.Equal(t, newUser1.Id, existingUser1.Id)
	assert.Equal(t, newUser1.CreatedAt, existingUser1.CreatedAt)
	assert.Equal(t, newUser1.UpdatedAt, existingUser1.UpdatedAt)
	assert.Equal(t, newUser1.Username, existingUser1.Username)
	assert.Equal(t, newUser1.Email, existingUser1.Email)
	assert.Equal(t, newUser1.Role, existingUser1.Role)

	// Get user2 by email and password
	reqGetUser2ByEmailPassword := &User{
		Email:    newUser2.Email,
		Password: &password.Password{},
	}
	reqGetUser2ByEmailPassword.Password.Set(reqUser2.Password.PlainText)
	existingUser2, err := usersService.GetUserByEmailPassword(reqGetUser2ByEmailPassword)
	require.NoError(t, err)
	require.NotNil(t, existingUser2)
	assert.Equal(t, newUser2.Id, existingUser2.Id)
	assert.Equal(t, newUser2.CreatedAt, existingUser2.CreatedAt)
	assert.Equal(t, newUser2.UpdatedAt, existingUser2.UpdatedAt)
	assert.Equal(t, newUser2.Username, existingUser2.Username)
	assert.Equal(t, newUser2.Email, existingUser2.Email)
	assert.Equal(t, newUser2.Role, existingUser2.Role)

	// Get user2 by email and wrong password
	reqGetUser2ByEmailPasswordInvalid := &User{
		Email:    newUser2.Email,
		Password: &password.Password{},
	}
	reqGetUser2ByEmailPasswordInvalid.Password.Set("wrong password")
	existingUser2Invalid, err := usersService.GetUserByEmailPassword(reqGetUser2ByEmailPasswordInvalid)
	require.Error(t, err)
	require.Nil(t, existingUser2Invalid)

	// Get user2 by wrong email and password
	reqGetUser2ByEmailPasswordInvalid2 := &User{
		Email:    "wrongemail@test.com",
		Password: &password.Password{},
	}
	reqGetUser2ByEmailPasswordInvalid2.Password.Set(reqUser2.Password.PlainText)
	existingUser2Invalid2, err := usersService.GetUserByEmailPassword(reqGetUser2ByEmailPasswordInvalid2)
	require.Error(t, err)
	require.Nil(t, existingUser2Invalid2)

	// Update user1
	reqUpdateUser1 := &User{
		Id:       newUser1.Id,
		Username: strPtr("user1Username"),
	}
	updatedUser, err := usersService.UpdateUser(reqUpdateUser1)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	assert.Equal(t, newUser1.Id, updatedUser.Id)
	assert.Equal(t, newUser1.CreatedAt, updatedUser.CreatedAt)
	assert.NotEqual(t, newUser1.UpdatedAt, updatedUser.UpdatedAt)
	assert.Equal(t, *reqUpdateUser1.Username, *updatedUser.Username)
	assert.Equal(t, newUser1.Email, updatedUser.Email)
	assert.Equal(t, newUser1.Role, updatedUser.Role)

	// Update user2
	reqUpdateUser2 := &User{
		Id:       newUser2.Id,
		Username: strPtr("user1Username"),
	}
	updatedUser2, err := usersService.UpdateUser(reqUpdateUser2)
	require.Error(t, err)
	require.Nil(t, updatedUser2)

	// Delete user1
	reqDeleteUser1 := &User{
		Id: newUser1.Id,
	}
	err = usersService.DeleteUser(reqDeleteUser1)
	require.NoError(t, err)

	// Delete user2
	reqDeleteUser2 := &User{
		Id: newUser2.Id,
	}
	err = usersService.DeleteUser(reqDeleteUser2)
	require.NoError(t, err)

	// Get deleted users
	reqGetDeletedUser1 := &User{
		Id: newUser1.Id,
	}
	deletedUser1, err := usersService.GetUserById(reqGetDeletedUser1)
	require.Error(t, err)
	require.Nil(t, deletedUser1)

	reqGetDeletedUser2 := &User{
		Id: newUser2.Id,
	}
	deletedUser2, err := usersService.GetUserById(reqGetDeletedUser2)
	require.Error(t, err)
	require.Nil(t, deletedUser2)
}

func strPtr(val string) *string {
	return &val
}
