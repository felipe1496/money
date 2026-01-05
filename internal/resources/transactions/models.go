package transactions

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
)

// ==============================================================================
// 1. HTTP MODELS
//    Models that represents request or response objects
// ==============================================================================

// Request body to create a simple expense
type CreateSimpleExpenseRequest struct {
	Name          string  `json:"name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gte=0,lte=999999"`
	ReferenceDate string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
	Description   string  `json:"description" biding:"min=0,max=400"`
	CategoryID    *string `json:"category_id"`
}

// Request body to create an income
type CreateIncomeRequest struct {
	Name          string  `json:"name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gte=0,lte=999999"`
	ReferenceDate string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
	Description   string  `json:"description" biding:"min=0,max=400"`
	CategoryID    *string `json:"category_id"`
}

type CreateIncomeResponse struct {
	Data CreateIncomeResponseData `json:"data"`
}

type CreateIncomeResponseData struct {
	Entry ViewEntry `json:"entry"`
}

// Request body to create an installment
type CreateInstallmentRequest struct {
	Name              string  `json:"name" binding:"required"`
	TotalAmount       float64 `json:"total_amount" binding:"required,gt=0,lte=999999"`
	TotalInstallments int     `json:"total_installments" binding:"required,gt=1,lte=100"`
	ReferenceDate     string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
	Description       string  `json:"description" biding:"min=0,max=400"`
	CategoryID        *string `json:"category_id"`
}

type CreateInstallmentResponse struct {
	Data CreateInstallmentResponseData `json:"data"`
}

type CreateInstallmentResponseData struct {
	Entries []ViewEntry `json:"entries"`
}

type CreateSimpleExpenseResponse struct {
	Data CreateSimpleExpenseResponseData `json:"data"`
}

type CreateSimpleExpenseResponseData struct {
	Entry ViewEntry `json:"entry"`
}

type ListEntriesResponse struct {
	Data  ListEntriesResponseData `json:"data"`
	Query utils.QueryMeta         `json:"query"`
}

type ListEntriesResponseData struct {
	Entries []ViewEntry `json:"entries"`
}

type UpdateSimpleExpenseResponse struct {
	Data UpdateSimpleExpenseResponseData `json:"data"`
}

type UpdateSimpleExpenseResponseData struct {
	Entry ViewEntry `json:"entry"`
}

type UpdateIncomeResponse struct {
	Data UpdateIncomeResponseData `json:"data"`
}

type UpdateIncomeResponseData struct {
	Entry ViewEntry `json:"entry"`
}

type UpdateSimpleExpenseRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Amount        *float64 `json:"amount"`
	ReferenceDate *string  `json:"reference_date" binding:"datetime=2006-01-02"`
	CategoryID    *string  `json:"category_id"`
}

type UpdateIncomeRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Amount        *float64 `json:"amount"`
	ReferenceDate *string  `json:"reference_date" binding:"datetime=2006-01-02"`
	CategoryID    *string  `json:"category_id"`
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
	CategoryID  *string
}

// Payload to create an entry in the database
type CreateEntryDTO struct {
	TransactionID string
	Amount        float64
	ReferenceDate string
}

// Payload to the use case create a simple expense
type CreateSimpleExpenseDTO struct {
	Name          string
	Amount        float64
	Description   string
	ReferenceDate string
	UserID        string
	CategoryID    *string
}

type CreateIncomeDTO struct {
	Name          string
	Amount        float64
	ReferenceDate string
	Description   string
	UserID        string
	CategoryID    *string
}

type CreateInstallmentDTO struct {
	Name              string
	TotalAmount       float64
	TotalInstallments int
	ReferenceDate     string
	Description       string
	UserID            string
	CategoryID        *string
}

type UpdateSimpleExpenseDTO struct {
	Name          *string
	Description   *string
	Amount        *float64
	ReferenceDate *string
	CategoryID    *string
}

type UpdateIncomeDTO struct {
	Name          *string
	Description   *string
	Amount        *float64
	ReferenceDate *string
	CategoryID    *string
}

type UpdateTransactionDTO struct {
	Name        *string
	Description *string
	CategoryID  *string
}

type UpdateEntryDTO struct {
	Amount        *float64
	ReferenceDate *string
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
	CreatedAt         time.Time       `json:"created_at"`
	ReferenceDate     string          `json:"reference_date"`
	CategoryID        *string         `json:"category_id,omitempty"`
	CategoryName      *string         `json:"category_name,omitempty"`
	CategoryColor     *string         `json:"category_color,omitempty"`
}

// Entries table record
type Entry struct {
	ID            string
	TransactionID string
	Amount        float64
	ReferenceDate string
	CreatedAt     time.Time
}

// Transactions table record
type Transaction struct {
	ID          string
	UserID      string
	Type        TransactionType
	Name        string
	Description *string
	CreatedAt   time.Time
	CategoryID  *string
}
