package route

import (
    "github.com/gin-gonic/gin"
    "github.com/swarmbit/spacemesh-state-api/config"
    "github.com/swarmbit/spacemesh-state-api/database"
    "github.com/swarmbit/spacemesh-state-api/network"
    "github.com/swarmbit/spacemesh-state-api/types"
    "net/http"
    "strconv"
)

type TransactionRoutes struct {
    db           *database.ReadDB
    networkUtils *network.NetworkUtils
    state        *network.NetworkState
}

func NewTransactionRoutes(db *database.ReadDB, networkUtils *network.NetworkUtils, state *network.NetworkState) *TransactionRoutes {
    routes := &TransactionRoutes{
        db:           db,
        networkUtils: networkUtils,
        state:        state,
    }
    return routes
}

func (t *TransactionRoutes) GetTransactions(c *gin.Context) {
    offsetStr := c.DefaultQuery("offset", "0")
    limitStr := c.DefaultQuery("limit", "20")
    sortStr := c.DefaultQuery("sort", "asc")
    completeStr := c.DefaultQuery("complete", "true")

    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "offset must be a valid integer",
        })
        return
    }
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "limit must be a valid integer",
        })
        return
    }

    if offset < 0 || limit < 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "offset and limit must be greater or equal to 0",
        })
        return
    }

    var sort int8
    if sortStr == "desc" {
        sort = -1
    } else {
        sort = 1
    }

    complete := completeStr == "true"

    transactions, errRewards := t.db.GetAllTransactions(int64(offset), int64(limit), sort, complete)
    count, errCount := t.db.CountAllTransactions()

    if errRewards != nil || errCount != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch transactions for layer",
        })
    } else if transactions != nil {

        transactionsResponse := make([]*types.Transaction, len(transactions))

        for i, v := range transactions {
            method := ""
            if v.Method == 0 {
                method = "Spawn"
            }
            if v.Method == 16 {
                method = "Spend"
            }
            if v.Method == 17 {
                method = "DrainVault"
            }
            transactionsResponse[i] = &types.Transaction{
                ID:               v.ID,
                Status:           v.Status,
                PrincipalAccount: v.PrincipaAccount,
                ReceiverAccount:  v.ReceiverAccount,
                VaultAccount:     v.VaultAccount,
                Fee:              v.Gas * v.GasPrice,
                Amount:           v.Amount,
                Layer:            v.Layer,
                Counter:          v.Counter,
                Method:           method,
                Timestamp:        int64(config.GenesisEpochSeconds + (v.Layer * config.LayerDuration)),
            }
        }

        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, transactionsResponse)
    } else {
        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, make([]*types.Transaction, 0))
    }

}

func (t *TransactionRoutes) GetTransaction(c *gin.Context) {
    transactionId := c.Param("transactionId")
    transaction, err := t.db.GetTransaction(transactionId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch transaction",
        })
        return
    }
    if transaction.ID == "" {
        c.JSON(http.StatusNotFound, gin.H{
            "status": "Not Found",
            "error":  "Node not found",
        })
        return
    }

    method := ""
    if transaction.Method == 0 {
        method = "Spawn"
    }
    if transaction.Method == 16 {
        method = "Spend"
    }
    if transaction.Method == 17 {
        method = "DrainVault"
    }

    c.JSON(200, &types.Transaction{
        ID:               transaction.ID,
        Status:           transaction.Status,
        PrincipalAccount: transaction.PrincipaAccount,
        ReceiverAccount:  transaction.ReceiverAccount,
        VaultAccount:     transaction.VaultAccount,
        Fee:              transaction.Gas * transaction.GasPrice,
        Amount:           transaction.Amount,
        Layer:            transaction.Layer,
        Counter:          transaction.Counter,
        Method:           method,
        Timestamp:        int64(config.GenesisEpochSeconds + (transaction.Layer * config.LayerDuration)),
    })
}
