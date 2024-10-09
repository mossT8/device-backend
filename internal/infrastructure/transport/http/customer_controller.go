package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"mossT8.github.com/device-backend/internal/domain/customer"
	"mossT8.github.com/device-backend/internal/domain/customer/model/entity"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/request"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/response"
)

type CustomerController struct {
	customerDomain customer.CustomerDomain
}

func NewCustomerController(conn *datastore.MySqlDataStore, server *iris.Application, custDomain customer.CustomerDomain) CustomerController {
	ac := CustomerController{
		customerDomain: custDomain,
	}

	server.Post(constants.ApiPrefix+"/account", ac.HandlePostAccount)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/update", ac.HandlePutAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/fetch", ac.HandleGetAccount)
	server.Get(constants.ApiPrefix+"/account/list", ac.HandleGetAccounts)

	server.Post(constants.ApiPrefix+"/account/{accountID:int64}/address", ac.HandlePostAddressForAccount)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/address/{addressID:int64}/update", ac.HandlePutAddressForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/address/{addressID:int64}/fetch", ac.HandleGetAddressForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/address/list", ac.HandleGetAddressesForAccount)

	server.Post(constants.ApiPrefix+"/account/{accountID:int64}/user", ac.HandlePostUserForAccount)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/user/{userID:int64}/update", ac.HandlePutUserForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/user/{userID:int64}/fetch", ac.HandleGetUserForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/user/list", ac.HandleGetUsersForAccount)

	return ac
}

func (ac *CustomerController) HandlePostAccount(ctx iris.Context) {
	var req request.Account
	requestId := GetRequestID(ctx)

	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account := entity.NewAccount(req.Email, req.Name, time.Now())
	if pErr := account.SetPassword(req.Password, uuid.New().String()); pErr != nil {
		RespondWithMappingError(ctx.ResponseWriter(), pErr.Error(), requestId)
		return
	}
	account.SetReceivesUpdates(req.ReceivesUpdates)

	if err := ac.customerDomain.AddAccount(requestId, &account); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Account{
		ID:              account.GetID(),
		Email:           account.GetEmail(),
		Name:            account.GetName(),
		ReceivesUpdates: account.GetReceivesUpdates(),
		CreatedAt:       account.GetCreatedAt(),
		ModifiedAt:      account.GetModifiedAt(),
	}, http.StatusCreated, requestId)
}

func (ac *CustomerController) HandlePutAccount(ctx iris.Context) {
	var req request.Account
	requestId := GetRequestID(ctx)

	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, err := ac.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	account.SetName(req.Name)
	account.SetReceivesUpdates(req.ReceivesUpdates)

	if err := ac.customerDomain.UpdateAccount(requestId, account); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Account{
		ID:              account.GetID(),
		Email:           account.GetEmail(),
		Name:            account.GetName(),
		ReceivesUpdates: account.GetReceivesUpdates(),
		CreatedAt:       account.GetCreatedAt(),
		ModifiedAt:      account.GetModifiedAt(),
	}, http.StatusOK, requestId)
}

func (ac *CustomerController) HandleGetAccount(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, err := ac.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Account{
		ID:              account.GetID(),
		Email:           account.GetEmail(),
		Name:            account.GetName(),
		ReceivesUpdates: account.GetReceivesUpdates(),
		CreatedAt:       account.GetCreatedAt(),
		ModifiedAt:      account.GetModifiedAt(),
	}, http.StatusOK, requestId)
}

func (ac *CustomerController) HandleGetAccounts(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	paginatedList, total, err := ac.customerDomain.ListAccounts(requestId, *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	accounts := make([]response.Account, 0)
	for _, account := range paginatedList {
		accounts = append(accounts, response.Account{
			ID:              account.GetID(),
			Email:           account.GetEmail(),
			Name:            account.GetName(),
			ReceivesUpdates: account.GetReceivesUpdates(),
			CreatedAt:       account.GetCreatedAt(),
			ModifiedAt:      account.GetModifiedAt(),
		})
	}

	RespondWithList(ctx.ResponseWriter(), accounts, *page, *pageSize, *total, http.StatusOK, requestId)
}

func (ac *CustomerController) HandlePostAddressForAccount(ctx iris.Context) {
	var req request.Address
	requestId := GetRequestID(ctx)

	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, aErr := ac.customerDomain.FetchAccount(requestId, accountID)
	if aErr != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, aErr)
		return
	}

	address := entity.NewAddress(account.GetID(), time.Now())
	address.SetName(req.Name)
	address.SetAddressLine1(req.AddressLine1)
	address.SetAddressLine2(req.AddressLine2)
	address.SetCity(req.City)
	address.SetState(req.State)
	address.SetPostalCode(req.PostalCode)
	address.SetCountry(req.Country)

	if err := ac.customerDomain.AddAddressForAccount(requestId, *account, &address); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Address{
		ID:           address.GetID(),
		Name:         address.GetName(),
		AddressLine1: address.GetAddressLine1(),
		AddressLine2: address.GetAddressLine2(),
		City:         address.GetCity(),
		State:        address.GetState(),
		PostalCode:   address.GetPostalCode(),
		Country:      address.GetCountry(),
		CreatedAt:    address.GetCreatedAt(),
		ModifiedAt:   address.GetModifiedAt(),
	}, http.StatusCreated, requestId)
}

func (ac *CustomerController) HandlePutAddressForAccount(ctx iris.Context) {
	var req request.Address
	requestId := GetRequestID(ctx)

	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, aErr := ac.customerDomain.FetchAccount(requestId, accountID)
	if aErr != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, aErr)
		return
	}

	addressID, err := ctx.Params().GetInt64("addressID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	address, adErr := ac.customerDomain.FetchAddressForAccount(requestId, *account, addressID)
	if adErr != nil {
		RespondWithMappingError(ctx.ResponseWriter(), adErr.Error(), requestId)
		return
	}

	address.SetName(req.Name)
	address.SetAddressLine1(req.AddressLine1)
	address.SetAddressLine2(req.AddressLine2)
	address.SetCity(req.City)
	address.SetState(req.State)
	address.SetPostalCode(req.PostalCode)
	address.SetCountry(req.Country)

	if err := ac.customerDomain.UpdateAddressForAccount(requestId, *account, address); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Address{
		ID:           address.GetID(),
		Name:         address.GetName(),
		AddressLine1: address.GetAddressLine1(),
		AddressLine2: address.GetAddressLine2(),
		City:         address.GetCity(),
		State:        address.GetState(),
		PostalCode:   address.GetPostalCode(),
		Country:      address.GetCountry(),
		CreatedAt:    address.GetCreatedAt(),
		ModifiedAt:   address.GetModifiedAt(),
	}, http.StatusCreated, requestId)
}

func (ac *CustomerController) HandleGetAddressForAccount(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, err := ac.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	addressID, err := ctx.Params().GetInt64("addressID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	address, err := ac.customerDomain.FetchAddressForAccount(requestId, *account, addressID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Address{
		ID:           address.GetID(),
		Name:         address.GetName(),
		AddressLine1: address.GetAddressLine1(),
		AddressLine2: address.GetAddressLine2(),
		City:         address.GetCity(),
		State:        address.GetState(),
		PostalCode:   address.GetPostalCode(),
		Country:      address.GetCountry(),
	}, http.StatusOK, requestId)
}

func (ac *CustomerController) HandleGetAddressesForAccount(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, err := ac.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	addresses, total, err := ac.customerDomain.ListAddressesForAccount(requestId, *account, *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	addressList := make([]response.Address, 0)
	for _, address := range addresses {
		addressList = append(addressList, response.Address{
			ID:           address.GetID(),
			Name:         address.GetName(),
			AddressLine1: address.GetAddressLine1(),
			AddressLine2: address.GetAddressLine2(),
			City:         address.GetCity(),
			State:        address.GetState(),
			PostalCode:   address.GetPostalCode(),
			Country:      address.GetCountry(),
			CreatedAt:    address.GetCreatedAt(),
			ModifiedAt:   address.GetModifiedAt(),
		})
	}

	RespondWithList(ctx.ResponseWriter(), addressList, *page, *pageSize, *total, http.StatusOK, requestId)
}

func (ac *CustomerController) HandlePostUserForAccount(ctx iris.Context) {
	var req request.User
	requestId := GetRequestID(ctx)

	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, aErr := ac.customerDomain.FetchAccount(requestId, accountID)
	if aErr != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, aErr)
		return
	}

	user := entity.NewUser(account.GetID(), req.Email, time.Now())
	user.SetCell(req.Cell)
	user.SetFirstName(req.FirstName)
	user.SetLastName(req.LastName)
	user.SetVerified(req.Verified)
	user.SetReceivesUpdates(req.ReceivesUpdates)

	if err := ac.customerDomain.AddUserForAccount(requestId, *account, &user); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.User{
		ID:              user.GetID(),
		Email:           user.GetEmail(),
		Cell:            user.GetCell(),
		FirstName:       user.GetFirstName(),
		LastName:        user.GetLastName(),
		Verified:        user.GetVerified(),
		ReceivesUpdates: user.GetReceivesUpdates(),
		CreatedAt:       user.GetCreatedAt(),
		ModifiedAt:      user.GetModifiedAt(),
	}, http.StatusCreated, requestId)
}

func (ac *CustomerController) HandlePutUserForAccount(ctx iris.Context) {
	var req request.User
	requestId := GetRequestID(ctx)

	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, aErr := ac.customerDomain.FetchAccount(requestId, accountID)
	if aErr != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, aErr)
		return
	}

	userID, err := ctx.Params().GetInt64("userID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	user, uErr := ac.customerDomain.FetchUserForAccount(requestId, *account, userID)
	if uErr != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, uErr)
		return
	}

	user.SetCell(req.Cell)
	user.SetFirstName(req.FirstName)
	user.SetLastName(req.LastName)
	user.SetVerified(req.Verified)
	user.SetReceivesUpdates(req.ReceivesUpdates)

	if err := ac.customerDomain.UpdateUserForAccount(requestId, *account, user); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.User{
		ID:              user.GetID(),
		Email:           user.GetEmail(),
		Cell:            user.GetCell(),
		FirstName:       user.GetFirstName(),
		LastName:        user.GetLastName(),
		Verified:        user.GetVerified(),
		ReceivesUpdates: user.GetReceivesUpdates(),
		CreatedAt:       user.GetCreatedAt(),
		ModifiedAt:      user.GetModifiedAt(),
	}, http.StatusOK, requestId)
}

func (ac *CustomerController) HandleGetUserForAccount(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, err := ac.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	userID, err := ctx.Params().GetInt64("userID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	user, err := ac.customerDomain.FetchUserForAccount(requestId, *account, userID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.User{
		ID:              user.GetID(),
		Email:           user.GetEmail(),
		Cell:            user.GetCell(),
		FirstName:       user.GetFirstName(),
		LastName:        user.GetLastName(),
		Verified:        user.GetVerified(),
		ReceivesUpdates: user.GetReceivesUpdates(),
		CreatedAt:       user.GetCreatedAt(),
		ModifiedAt:      user.GetModifiedAt(),
	}, http.StatusOK, requestId)
}

func (ac *CustomerController) HandleGetUsersForAccount(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	account, err := ac.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	users, total, err := ac.customerDomain.ListUsersForAccount(requestId, *account, *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	userList := make([]response.User, 0)
	for _, user := range users {
		userList = append(userList, response.User{
			ID:              user.GetID(),
			Email:           user.GetEmail(),
			Cell:            user.GetCell(),
			FirstName:       user.GetFirstName(),
			LastName:        user.GetLastName(),
			Verified:        user.GetVerified(),
			ReceivesUpdates: user.GetReceivesUpdates(),
			CreatedAt:       user.GetCreatedAt(),
			ModifiedAt:      user.GetModifiedAt(),
		})
	}

	RespondWithList(ctx.ResponseWriter(), userList, *page, *pageSize, *total, http.StatusOK, requestId)
}
