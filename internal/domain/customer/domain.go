package customer

import (
	"mossT8.github.com/device-backend/internal/domain/customer/model/entity"
	"mossT8.github.com/device-backend/internal/infrastructure/logger"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

type CustomerDomain interface {
	AddAccount(requestId string, account *entity.Account) error
	FetchAccount(requestId string, accountId int64) (*entity.Account, error)
	ListAccounts(requestId string, page, pageSize int64) ([]entity.Account, *int64, error)
	UpdateAccount(requestId string, account *entity.Account) error
	DeleteAccount(requestId string, accountId int64) error

	AddAddressForAccount(requestId string, account entity.Account, address *entity.Address) error
	FetchAddressForAccount(requestId string, account entity.Account, addressId int64) (*entity.Address, error)
	ListAddressesForAccount(requestId string, account entity.Account, page, pageSize int64) ([]entity.Address, *int64, error)
	UpdateAddressForAccount(requestId string, account entity.Account, address *entity.Address) error
	DeleteAddressForAccount(requestId string, account entity.Account, addressId int64) error

	AddUserForAccount(requestId string, account entity.Account, user *entity.User) error
	FetchUserForAccount(requestId string, account entity.Account, userId int64) (*entity.User, error)
	ListUsersForAccount(requestId string, account entity.Account, page, pageSize int64) ([]entity.User, *int64, error)
	UpdateUserForAccount(requestId string, account entity.Account, user *entity.User) error
	DeleteUserForAccount(requestId string, account entity.Account, userId int64) error
}

type CustomerDomainImpl struct {
	dbConn *datastore.MySqlDataStore
}

func NewCustomerDomain(conn *datastore.MySqlDataStore) CustomerDomain {
	return &CustomerDomainImpl{
		dbConn: conn,
	}
}

// Account operations
func (u *CustomerDomainImpl) AddAccount(requestId string, account *entity.Account) error {
	if aErr := account.AddAccount(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to create account %+v", account)
		return aErr
	}
	return nil
}

func (u *CustomerDomainImpl) FetchAccount(requestId string, accountId int64) (*entity.Account, error) {
	account := &entity.Account{}
	account.SetID(accountId)
	if aErr := account.GetAccountByID(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to get account by ID %d", accountId)
		return nil, aErr
	}
	return account, nil
}

func (u *CustomerDomainImpl) UpdateAccount(requestId string, account *entity.Account) error {
	if aErr := account.UpdateAccount(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to update account %+v", account)
		return aErr
	}
	return nil
}

func (u *CustomerDomainImpl) DeleteAccount(requestId string, accountId int64) error {
	account := &entity.Account{}
	account.SetID(accountId)
	if aErr := account.DeleteAccount(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to delete account by ID %d", accountId)
		return aErr
	}
	return nil
}

func (u *CustomerDomainImpl) ListAccounts(requestId string, page, pageSize int64) ([]entity.Account, *int64, error) {
	queryAccount := entity.Account{}
	accounts, err := queryAccount.ListAccounts(*u.dbConn, page, pageSize)
	if err != nil {
		logger.Errorf(requestId, "unable to list accounts")
		return nil, nil, err
	}

	total, err := queryAccount.CountAccounts(*u.dbConn)
	if err != nil {
		logger.Errorf(requestId, "unable to count all accounts")
		return nil, nil, err
	}

	return accounts, total, nil
}

// Address operations
func (u *CustomerDomainImpl) AddAddressForAccount(requestId string, account entity.Account, address *entity.Address) error {
	if aErr := address.AddAddress(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to create address %+v", address)
		return aErr
	}
	return nil
}

func (u *CustomerDomainImpl) FetchAddressForAccount(requestId string, account entity.Account, addressId int64) (*entity.Address, error) {
	address := &entity.Address{}
	address.SetID(addressId)
	if aErr := address.GetAddressByID(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to get address by ID %d", addressId)
		return nil, aErr
	}
	return address, nil
}

func (u *CustomerDomainImpl) UpdateAddressForAccount(requestId string, account entity.Account, address *entity.Address) error {
	if aErr := address.UpdateAddress(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to update address %+v", address)
		return aErr
	}
	return nil
}

func (u *CustomerDomainImpl) DeleteAddressForAccount(requestId string, account entity.Account, addressId int64) error {
	address := &entity.Address{}
	address.SetID(addressId)
	if aErr := address.DeleteAddress(*u.dbConn); aErr != nil {
		logger.Errorf(requestId, "unable to delete address by ID %d", addressId)
		return aErr
	}
	return nil
}

func (u *CustomerDomainImpl) ListAddressesForAccount(requestId string, account entity.Account, page, pageSize int64) ([]entity.Address, *int64, error) {
	queryAddress := entity.Address{}
	queryAddress.SetAccountId(account.GetID())

	addresses, err := queryAddress.ListAddresses(*u.dbConn, page, pageSize)
	if err != nil {
		logger.Errorf(requestId, "unable to list addresses for account ID %d", account.GetID())
		return nil, nil, err
	}

	total, err := queryAddress.CountAddresses(*u.dbConn)
	if err != nil {
		logger.Errorf(requestId, "unable to count all addresses for account ID %d", account.GetID())
		return nil, nil, err
	}

	return addresses, total, nil
}

// User operations
func (u *CustomerDomainImpl) AddUserForAccount(requestId string, account entity.Account, user *entity.User) error {
	if uErr := user.AddUser(*u.dbConn); uErr != nil {
		logger.Errorf(requestId, "unable to create user %+v", user)
		return uErr
	}
	return nil
}

func (u *CustomerDomainImpl) FetchUserForAccount(requestId string, account entity.Account, userId int64) (*entity.User, error) {
	user := &entity.User{}
	user.SetID(userId)
	if uErr := user.GetUserByID(*u.dbConn); uErr != nil {
		logger.Errorf(requestId, "unable to get user by ID %d", userId)
		return nil, uErr
	}
	return user, nil
}

func (u *CustomerDomainImpl) UpdateUserForAccount(requestId string, account entity.Account, user *entity.User) error {
	if uErr := user.UpdateUser(*u.dbConn); uErr != nil {
		logger.Errorf(requestId, "unable to update user %+v", user)
		return uErr
	}
	return nil
}

func (u *CustomerDomainImpl) DeleteUserForAccount(requestId string, account entity.Account, userId int64) error {
	user := &entity.User{}
	user.SetID(userId)
	if uErr := user.DeleteUser(*u.dbConn); uErr != nil {
		logger.Errorf(requestId, "unable to delete user by ID %d", userId)
		return uErr
	}
	return nil
}

func (u *CustomerDomainImpl) ListUsersForAccount(requestId string, account entity.Account, page, size int64) ([]entity.User, *int64, error) {
	queryUser := entity.User{}
	queryUser.SetAccountId(account.GetID())

	users, err := queryUser.ListUsers(*u.dbConn, page, size)
	if err != nil {
		logger.Errorf(requestId, "unable to list users for account ID %d", account.GetID())
		return nil, nil, err
	}

	total, err := queryUser.CountUsers(*u.dbConn)
	if err != nil {
		logger.Errorf(requestId, "unable to count all users for account ID %d", account.GetID())
		return nil, nil, err
	}

	return users, total, nil
}
