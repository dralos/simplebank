package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/dralos/simplebank/db/mock"
	db "github.com/dralos/simplebank/db/sqlc"
	"github.com/dralos/simplebank/token"
	"github.com/dralos/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransferAPI(t *testing.T) {
	fromAccount := db.Account{ID: 1, Owner: "alice", Balance: 1000, Currency: util.USD}
	toAccount := db.Account{ID: 2, Owner: "bob", Balance: 500, Currency: util.USD}
	result := db.TransferTxResult{
		Transfer:    db.Transfer{ID: 10, FromAccountID: 1, ToAccountID: 2, Amount: 100},
		FromEntry:   db.Entry{ID: 11, AccountID: 1, Amount: -100},
		ToEntry:     db.Entry{ID: 12, AccountID: 2, Amount: 100},
		FromAccount: fromAccount,
		ToAccount:   toAccount,
	}

	testCases := []struct {
		name          string
		body          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(2)).Return(toAccount, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, arg db.TransferTxParams) (db.TransferTxResult, error) {
						require.EqualValues(t, 1, arg.FromAccountID)
						require.EqualValues(t, 2, arg.ToAccountID)
						require.EqualValues(t, 100, arg.Amount)
						return result, nil
					},
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var got db.TransferTxResult
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
				require.EqualValues(t, result.Transfer.ID, got.Transfer.ID)
				require.EqualValues(t, result.FromEntry.ID, got.FromEntry.ID)
				require.EqualValues(t, result.ToEntry.ID, got.ToEntry.ID)
			},
		},
		{
			name: "InvalidJSON",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(2)).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "CurrencyMismatch",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(db.Account{ID: 1, Currency: util.EUR}, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(db.Account{}, errors.New("db down"))
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "alice", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(2)).Return(toAccount, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Return(db.TransferTxResult{}, errors.New("tx failed"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// No authorization header
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: `{"from_account_id":1,"to_account_id":2,"amount":100,"currency":"USD"}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// No authorization header
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewBufferString(tc.body))
			request.Header.Set("Content-Type", "application/json")
			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
