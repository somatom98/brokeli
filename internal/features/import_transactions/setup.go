package import_transactions

func (f *Feature) Setup() {
	f.httpHandler.HandleFunc("POST /api/import-transactions", f.handleImportTransactions)
}
