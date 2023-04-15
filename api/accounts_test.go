package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Hypersus/simplebank/db/mock"
	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/token"
	"github.com/Hypersus/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user := randomUsername()
	account := randomAccountByUsername(user)
	var testCases = []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authType, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMacthAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authType, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authType, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authType, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	ctrl := gomock.NewController(t)
	for i := range testCases {
		tc := testCases[i]
		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)
		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/accounts/%d", tc.accountID)
		request := httptest.NewRequest(http.MethodGet, url, nil)
		tc.setupAuth(t, request, server.tokenMaker)
		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func randomUsername() string {
	return util.RandomOwner()
}

func randomAccountByUsername(username string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMacthAccount(t *testing.T, respBody *bytes.Buffer, account db.Account) {
	var gotAccount db.Account
	body := respBody.Bytes()
	require.NoError(t, json.Unmarshal(body, &gotAccount))
	require.Equal(t, account, gotAccount)
}
