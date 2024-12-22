package main

import (
	"fmt"
	"log"
	"reflect"

	. "github.com/tigerbeetle/tigerbeetle-go"
	. "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Since we only require Go 1.17 we can't do this as a generic function
// even though that would be fine. So do the dynamic approach for now.
func assert(a, b any, field string) {
	if !reflect.DeepEqual(a, b) {
		log.Fatalf("Expected %s to be [%+v (%T)], got: [%+v (%T)]", field, b, b, a, a)
	}
}

func main() {
	fmt.Println("Import ok!")
	client, err := NewClient(ToUint128(0), []string{"127.0.0.1:3000", "127.0.0.1:3001", "127.0.0.1:3002"})
	if err != nil {
		log.Printf("Error creating client: %s", err)
		return
	}
	// https://docs.tigerbeetle.com/coding/recipes/balance-bounds

	destAcID := ToUint128(100)
	sourceAcID := ToUint128(101)

	operatorAcId := ToUint128(900)
	controlAcId := ToUint128(901)

	limitAmount := ToUint128(5000)
	controlLedger := 1
	dummyTransferCode := 420

	account0 := Account{
		ID:     destAcID,
		Ledger: 1,
		Code:   718,
		Flags: AccountFlags{
			DebitsMustNotExceedCredits: true,
		}.ToUint16(),
	}
	account1 := Account{
		ID:     sourceAcID,
		Ledger: 1,
		Code:   718,
		Flags: AccountFlags{
			CreditsMustNotExceedDebits: true,
		}.ToUint16(),
	}

	accountErrors, err := client.CreateAccounts([]Account{
		account0,
		account1,
		{
			ID:     operatorAcId,
			Ledger: uint32(controlLedger),
			Code:   uint16(dummyTransferCode),
			Flags: AccountFlags{
				DebitsMustNotExceedCredits: true,
			}.ToUint16(),
		},
		{
			ID:     controlAcId,
			Ledger: uint32(controlLedger),
			Code:   uint16(dummyTransferCode),
			Flags: AccountFlags{
				CreditsMustNotExceedDebits: true,
			}.ToUint16(),
		},
	})
	log.Println(accountErrors)

	curTid := ID()
	pendingId := ID()
	log.Println(curTid, pendingId)
	// Start a pending transfer
	transferRes, err := client.CreateTransfers([]Transfer{
		{
			ID:              curTid,
			DebitAccountID:  sourceAcID,
			CreditAccountID: destAcID,
			Amount:          ToUint128(500),
			Ledger:          1,
			Code:            1,
			Flags:           TransferFlags{Linked: true}.ToUint16(),
		},
		{
			ID:              ID(),
			DebitAccountID:  controlAcId,
			CreditAccountID: operatorAcId,
			Amount:          limitAmount, // LIMIT
			Ledger:          uint32(controlLedger),
			Code:            1,
			Flags:           TransferFlags{Linked: true}.ToUint16(),
		},
		{
			ID:              pendingId,
			DebitAccountID:  destAcID,
			CreditAccountID: controlAcId,
			Amount:          AmountMax,
			Ledger:          uint32(controlLedger),
			Code:            1,
			Flags: TransferFlags{
				Linked:         true,
				BalancingDebit: true,
				Pending:        true,
			}.ToUint16(),
		},
		{
			ID:              ID(),
			DebitAccountID:  ToUint128(0),
			CreditAccountID: ToUint128(0),
			Amount:          ToUint128(0),
			PendingID:       pendingId,
			Ledger:          uint32(controlLedger),
			Code:            1,
			Flags: TransferFlags{
				Linked:              true,
				VoidPendingTransfer: true,
			}.ToUint16(),
		},
		{
			ID:              ID(),
			DebitAccountID:  operatorAcId,
			CreditAccountID: controlAcId,
			Amount:          limitAmount,
			Ledger:          uint32(controlLedger),
			Code:            1,
			Flags:           TransferFlags{Linked: false}.ToUint16(),
		},
	})
	log.Println(transferRes, "<-transferRes")
	if err != nil {
		log.Fatalf("Error creating transfer: %s", err)
	}

	for _, err := range transferRes {
		log.Fatalf("Error creating transfer: %s", err.Result)
	}

	// Validate accounts pending and posted debits/credits before finishing the two-phase transfer
	accounts, err := client.LookupAccounts([]Uint128{destAcID, sourceAcID})
	log.Println(accounts, "accounts")
	if err != nil {
		log.Fatalf("Could not fetch accounts: %s", err)
	}

	defer client.Close()
}
