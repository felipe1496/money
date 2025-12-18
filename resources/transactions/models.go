package transactions

// ==============================================================================
// 1. HTTP MODELS
//    Models that represents request or response objects
// ==============================================================================

// Request body to create a simple expense
type CreateSimpleExpenseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,lt=0,gte=-999999"`
	Period      string  `json:"period" binding:"required,len=6"`
	Description string  `json:"description" biding:"min=0,max=400"`
}

// ==============================================================================
// 2. DTO MODELS
//    Models that represents data transfer objects between api layers
// ==============================================================================

// Payload to create a transaction in the database
type CreateTransactionDTO struct {
	UserID      string
	Type        TransactionType
	Name        string
	Description *string
}

// Payload to create an entry in the database
type CreateEntryDTO struct {
	TransactionID string
	Amount        float64
	Period        string
}

// Payload to the use case create a simple expense
type CreateSimpleExpenseDTO struct {
	Name        string
	Amount      float64
	Period      string
	Description string
	UserID      string
}

// ==============================================================================
// 3. DATABASE
//    Models that represents database objects
// ==============================================================================

// View that mixes the entries with the transaction information, riched with some valuable information about the totality of this relationship
type ViewEntry struct {
	ID                string          `json:"id"`
	TransactionID     string          `json:"transaction_id"`
	Name              string          `json:"name"`
	Description       *string         `json:"description"`
	Amount            float64         `json:"amount"`
	Period            string          `json:"period"`
	UserID            string          `json:"user_id"`
	Type              TransactionType `json:"type"`
	TotalAmount       float64         `json:"total_amount"`
	Installment       int             `json:"installment"`
	TotalInstallments int             `json:"total_installments"`
	CreatedAt         string          `json:"created_at"`
}

// Entries table record
type Entry struct {
	ID            string
	TransactionID string
	Amount        float64
	Period        string
	CreatedAt     string
}

// Transactions table record
type Transaction struct {
	ID          string
	UserID      string
	Type        TransactionType
	Name        string
	Description string
	CreatedAt   string
}
