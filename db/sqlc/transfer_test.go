package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zsoltardai/simple_bank/util"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)

	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)

	require.Equal(t, account1.ID, transfer.FromAccountID)

	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)

	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)

	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)

	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.Error(t, err)

	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)

	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)

	require.Equal(t, arg.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}
