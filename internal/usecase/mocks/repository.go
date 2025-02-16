// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "merch-shop/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// BuyMerch provides a mock function with given fields: ctx, userID, itemName, itemPrice
func (_m *Repository) BuyMerch(ctx context.Context, userID uint64, itemName string, itemPrice uint64) error {
	ret := _m.Called(ctx, userID, itemName, itemPrice)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string, uint64) error); ok {
		r0 = rf(ctx, userID, itemName, itemPrice)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: ctx, creds
func (_m *Repository) CreateUser(ctx context.Context, creds domain.Credentials) (uint64, error) {
	ret := _m.Called(ctx, creds)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, domain.Credentials) uint64); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Credentials) error); ok {
		r1 = rf(ctx, creds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMerchPrice provides a mock function with given fields: ctx, itemName
func (_m *Repository) GetMerchPrice(ctx context.Context, itemName string) (uint64, error) {
	ret := _m.Called(ctx, itemName)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, string) uint64); ok {
		r0 = rf(ctx, itemName)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, itemName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ctx, userID
func (_m *Repository) GetUserByID(ctx context.Context, userID uint64) (domain.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 domain.User
	if rf, ok := ret.Get(0).(func(context.Context, uint64) domain.User); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: ctx, username
func (_m *Repository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	ret := _m.Called(ctx, username)

	var r0 domain.User
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInventory provides a mock function with given fields: ctx, userID
func (_m *Repository) GetUserInventory(ctx context.Context, userID uint64) ([]domain.Inventory, error) {
	ret := _m.Called(ctx, userID)

	var r0 []domain.Inventory
	if rf, ok := ret.Get(0).(func(context.Context, uint64) []domain.Inventory); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Inventory)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserTransactions provides a mock function with given fields: ctx, userID
func (_m *Repository) GetUserTransactions(ctx context.Context, userID uint64) (domain.CoinHistory, error) {
	ret := _m.Called(ctx, userID)

	var r0 domain.CoinHistory
	if rf, ok := ret.Get(0).(func(context.Context, uint64) domain.CoinHistory); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(domain.CoinHistory)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferCoins provides a mock function with given fields: ctx, fromUserID, toUserID, amount
func (_m *Repository) TransferCoins(ctx context.Context, fromUserID uint64, toUserID uint64, amount uint64) error {
	ret := _m.Called(ctx, fromUserID, toUserID, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, uint64) error); ok {
		r0 = rf(ctx, fromUserID, toUserID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
