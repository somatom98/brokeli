package import_transactions

import (
	"net/http"
)

func (f *Feature) handleImportTransactions(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file_path")
	if filePath == "" {
		filePath = "internal/features/import_transactions/transactions.csv"
	}

	if err := f.ImportTransactions(r.Context(), filePath); err != nil {
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
