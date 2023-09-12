package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)

	account2 := createRandomAccount(t)

	// sending 10 in 5 go routines in order to check concurrency

	n := 5

	amount := int64(10)

	errs := make(chan error)

	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err

			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)

		result := <-results

		require.NotEmpty(t, result)

		// checking transfer

		transfer := result.Transfer

		require.NotEmpty(t, transfer)

		require.Equal(t, account1.ID, transfer.FromAccountID)

		require.Equal(t, account2.ID, transfer.ToAccountID)

		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)

		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err)

		// checking from entry

		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)

		require.Equal(t, account1.ID, fromEntry.AccountID)

		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		require.NoError(t, err)

		// checking to entry

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)

		require.Equal(t, account2.ID, toEntry.AccountID)

		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)

		require.NoError(t, err)

		// check accounts

		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)

		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount

		require.NotEmpty(t, toAccount)

		require.Equal(t, account2.ID, toAccount.ID)

		// check the balance of the accounts

		diff1 := account1.Balance - fromAccount.Balance

		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)

		require.True(t, diff1 > 0)

		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)

		existed[k] = true
	}

	// check updated balances

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)

	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)

	account2 := createRandomAccount(t)

	// sending 10 in 5 go routines in order to check concurrency

	n := 10

	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {

		fromAccountID := account1.ID

		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check updated balances

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)

	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}