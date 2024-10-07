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

type AccountController struct {
	customerDomain customer.CustomerDomain
}

func NewAccountController(conn *datastore.MySqlDataStore, server *iris.Application, custDomain customer.CustomerDomain) AccountController {
	ac := AccountController{
		customerDomain: custDomain,
	}

	server.Post(constants.ApiPrefix+"/account}", ac.HandlePostAccount)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/update", ac.HandlePutAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/fetch", ac.HandleGetAccount)
	server.Get(constants.ApiPrefix+"/account/list", ac.HandleGetAccounts)

	server.Post(constants.ApiPrefix+"/account/{accountID:int64}/address", ac.HandlePostAddressForAccount)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/adress/{addressID:int64}/update", ac.HandlePutAddressForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/address/{addressID:int64}/fetch", ac.HandleGetAddressForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/address/list", ac.HandleGetAddressesForAccount)

	server.Post(constants.ApiPrefix+"/account/{accountID:int64}/user", ac.HandlePostUserForAccount)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/user/{userId:int64}/update", ac.HandlePutUserForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/user/{userId:int64}/fetch", ac.HandleGetUserForAccount)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/user/list", ac.HandleGetUsersForAccount)

	return ac
}

func (ac *AccountController) HandlePostAccount(ctx iris.Context) {
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
		Email:           account.GetEmail(),
		Name:            account.GetName(),
		ReceivesUpdates: account.GetReceivesUpdates(),
	}, http.StatusCreated, requestId)
}

func (ac *AccountController) HandlePutAccount(ctx iris.Context) {
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
	account.SetEmail(req.Email)
	account.SetReceivesUpdates(req.ReceivesUpdates)

	if req.Password != "" {
		if pErr := account.SetPassword(req.Password, uuid.New().String()); pErr != nil {
			RespondWithMappingError(ctx.ResponseWriter(), pErr.Error(), requestId)
			return
		}
	}

	if err := ac.customerDomain.UpdateAccount(requestId, account); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), response.Account{
		Email:           account.GetEmail(),
		Name:            account.GetName(),
		ReceivesUpdates: account.GetReceivesUpdates(),
	}, http.StatusOK, requestId)
}

func (ac *AccountController) HandleGetAccount(ctx iris.Context) {
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
		Email:           account.GetEmail(),
		Name:            account.GetName(),
		ReceivesUpdates: account.GetReceivesUpdates(),
	}, http.StatusOK, requestId)
}

func (ac *AccountController) HandleGetAccounts(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	page, pageSize, err := GetPageAndPageSize(ctx)
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
			Email:           account.GetEmail(),
			Name:            account.GetName(),
			ReceivesUpdates: account.GetReceivesUpdates(),
		})
	}

	RespondWithList(ctx.ResponseWriter(), accounts, *page, *pageSize, *total, http.StatusOK, requestId)
}

func (ac *AccountController) HandlePostAddressForAccount(ctx iris.Context) {
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
		Name:         address.GetName(),
		AddressLine1: address.GetAddressLine1(),
		AddressLine2: address.GetAddressLine2(),
		City:         address.GetCity(),
		State:        address.GetState(),
		PostalCode:   address.GetPostalCode(),
		Country:      address.GetCountry(),
	}, http.StatusCreated, requestId)
}

func (ac *AccountController) HandlePutAddressForAccount(ctx iris.Context) {
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
		Name:         address.GetName(),
		AddressLine1: address.GetAddressLine1(),
		AddressLine2: address.GetAddressLine2(),
		City:         address.GetCity(),
		State:        address.GetState(),
		PostalCode:   address.GetPostalCode(),
		Country:      address.GetCountry(),
	}, http.StatusCreated, requestId)
}

func (ac *AccountController) HandleGetAddressForAccount(ctx iris.Context) {
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
		Name:         address.GetName(),
		AddressLine1: address.GetAddressLine1(),
		AddressLine2: address.GetAddressLine2(),
		City:         address.GetCity(),
		State:        address.GetState(),
		PostalCode:   address.GetPostalCode(),
		Country:      address.GetCountry(),
	}, http.StatusOK, requestId)
}

func (ac *AccountController) HandleGetAddressesForAccount(ctx iris.Context) {
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

	page, pageSize, err := GetPageAndPageSize(ctx)
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
			Name:         address.GetName(),
			AddressLine1: address.GetAddressLine1(),
			AddressLine2: address.GetAddressLine2(),
			City:         address.GetCity(),
			State:        address.GetState(),
			PostalCode:   address.GetPostalCode(),
			Country:      address.GetCountry(),
		})
	}

	RespondWithList(ctx.ResponseWriter(), addressList, *page, *pageSize, *total, http.StatusOK, requestId)
}

func (ac *AccountController) HandlePostUserForAccount(ctx iris.Context) {
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
		Email:           user.GetEmail(),
		Cell:            user.GetCell(),
		FirstName:       user.GetFirstName(),
		LastName:        user.GetLastName(),
		Verified:        user.GetVerified(),
		ReceivesUpdates: user.GetReceivesUpdates(),
	}, http.StatusCreated, requestId)
}

func (ac *AccountController) HandlePutUserForAccount(ctx iris.Context) {
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
		Email:           user.GetEmail(),
		Cell:            user.GetCell(),
		FirstName:       user.GetFirstName(),
		LastName:        user.GetLastName(),
		Verified:        user.GetVerified(),
		ReceivesUpdates: user.GetReceivesUpdates(),
	}, http.StatusOK, requestId)
}

func (ac *AccountController) HandleGetUserForAccount(ctx iris.Context) {
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
		Email:           user.GetEmail(),
		Cell:            user.GetCell(),
		FirstName:       user.GetFirstName(),
		LastName:        user.GetLastName(),
		Verified:        user.GetVerified(),
		ReceivesUpdates: user.GetReceivesUpdates(),
	}, http.StatusOK, requestId)
}

func (ac *AccountController) HandleGetUsersForAccount(ctx iris.Context) {
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

	page, pageSize, err := GetPageAndPageSize(ctx)
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
			Email:           user.GetEmail(),
			Cell:            user.GetCell(),
			FirstName:       user.GetFirstName(),
			LastName:        user.GetLastName(),
			Verified:        user.GetVerified(),
			ReceivesUpdates: user.GetReceivesUpdates(),
		})
	}

	RespondWithList(ctx.ResponseWriter(), userList, *page, *pageSize, *total, http.StatusOK, requestId)
}
