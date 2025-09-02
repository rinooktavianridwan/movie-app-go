package jobs

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "movie-app-go/internal/enums"
    "movie-app-go/internal/models"

    "github.com/hibiken/asynq"
    "gorm.io/gorm"
)

type PaymentTimeoutJob struct {
    TransactionID uint `json:"transaction_id"`
}

type PaymentJobHandler struct {
    DB *gorm.DB
}

func NewPaymentJobHandler(db *gorm.DB) *PaymentJobHandler {
    return &PaymentJobHandler{DB: db}
}

func (h *PaymentJobHandler) HandlePaymentTimeout(ctx context.Context, t *asynq.Task) error {
    var payload PaymentTimeoutJob
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v", err)
    }

    log.Printf("Processing payment timeout for transaction ID: %d", payload.TransactionID)

    err := h.DB.Transaction(func(tx *gorm.DB) error {
        var transaction models.Transaction
        if err := tx.First(&transaction, payload.TransactionID).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                log.Printf("Transaction %d not found, might be already processed", payload.TransactionID)
                return nil
            }
            return err
        }

        if transaction.PaymentStatus != enums.PaymentStatusPending {
            log.Printf("Transaction %d status is %s, skipping", payload.TransactionID, transaction.PaymentStatus)
            return nil
        }

        transaction.PaymentStatus = enums.PaymentStatusFailed
        if err := tx.Save(&transaction).Error; err != nil {
            return err
        }
        
        if err := tx.Model(&models.Ticket{}).
            Where("transaction_id = ?", payload.TransactionID).
            Update("status", enums.TicketStatusCancelled).Error; err != nil {
            return err
        }

        log.Printf("Transaction %d marked as failed and tickets cancelled", payload.TransactionID)
        return nil
    })

    return err
}

func CreatePaymentTimeoutPayload(transactionID uint) ([]byte, error) {
    payload := PaymentTimeoutJob{
        TransactionID: transactionID,
    }
    return json.Marshal(payload)
}