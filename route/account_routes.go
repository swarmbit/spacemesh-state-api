package route

import (
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/swarmbit/spacemesh-state-api/config"
    "github.com/swarmbit/spacemesh-state-api/database"
    "github.com/swarmbit/spacemesh-state-api/network"
    "github.com/swarmbit/spacemesh-state-api/price"
    "github.com/swarmbit/spacemesh-state-api/types"
)

type AccountRoutes struct {
    db            *database.ReadDB
    networkUtils  *network.NetworkUtils
    state         *network.NetworkState
    priceResolver *price.PriceResolver
}

func NewAccountRoutes(
    readDB *database.ReadDB,
    networkUtils *network.NetworkUtils,
    state *network.NetworkState,
    priceResolver *price.PriceResolver,
) *AccountRoutes {
    return &AccountRoutes{
        db:            readDB,
        networkUtils:  networkUtils,
        state:         state,
        priceResolver: priceResolver,
    }
}

func (a *AccountRoutes) GetAccountsPost(c *gin.Context) {
    offsetStr := c.DefaultQuery("offset", "0")
    limitStr := c.DefaultQuery("limit", "20")
    sortStr := c.DefaultQuery("sort", "desc")

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
    if sortStr == "asc" {
        sort = 1
    } else {
        sort = -1
    }

    epochStr := c.Param("epoch")
    epoch, err := strconv.Atoi(epochStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "epoch must be a valid integer",
        })
        return
    }

    accounts, errAccounts := a.db.GetAccountsPostEpoch(epoch-1, int64(offset), int64(limit), sort)
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "failed to fetch accounts",
        })
        return
    }

    count, errCount := a.db.CountAccountsPostEpoch(epoch - 1)
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "failed to count accounts",
        })
        return
    }

    if errAccounts != nil || errCount != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch account",
        })
    } else if accounts != nil {

        accountsResponse := make([]*types.AccountPostResponse, len(accounts))

        for i, v := range accounts {
            accountsResponse[i] = &types.AccountPostResponse{
                Account:                v.Id.Coinbase,
                TotalEffectiveNumUnits: v.TotalEffectiveNumUnits,
                TotalAtx:               v.TotalAtx,
                TotalWeight:            v.TotalWeight,
            }
        }

        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, accountsResponse)
    } else {
        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, make([]*types.AccountAtxDoc, 0))
    }
}

func (a *AccountRoutes) GetAccounts(c *gin.Context) {

    offsetStr := c.DefaultQuery("offset", "0")
    limitStr := c.DefaultQuery("limit", "20")
    sortStr := c.DefaultQuery("sort", "desc")

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
    if sortStr == "asc" {
        sort = 1
    } else {
        sort = -1
    }

    accounts, errAccounts := a.db.GetAccounts(int64(offset), int64(limit), sort)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to get accounts",
        })
        return
    }

    count, errCount := a.db.CountAccounts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to count accounts",
        })
        return
    }

    if errAccounts != nil || errCount != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch transactions for account",
        })
    } else if accounts != nil {

        accountsResponse := make([]*types.ShortAccount, len(accounts))

        for i, v := range accounts {
            priceValue := a.priceResolver.GetPrice()
            dollarValue := int64(-1)
            if priceValue > -1 {
                dollarValue = int64(priceValue * float64(v.Balance))
            }
            accountsResponse[i] = &types.ShortAccount{
                Balance:      v.Balance,
                Address:      v.Address,
                USDValue:     dollarValue,
                TotalRewards: v.TotalRewards,
            }
        }

        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, accountsResponse)
    } else {
        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, make([]*types.Transaction, 0))
    }
}

func (a *AccountRoutes) GetAccountGroup(c *gin.Context) {
    var req types.AccounGroupRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result, err := a.db.GetAccountsGroup(req.Accounts)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch account group",
        })
        return
    }

    priceValue := a.priceResolver.GetPrice()
    dollarValue := int64(-1)
    if priceValue > -1 {
        dollarValue = int64(priceValue * float64(result.Balance))
    }

    c.JSON(200, &types.AccountGroupResponse{
        Balance:      uint64(result.Balance),
        USDValue:     dollarValue,
        TotalRewards: uint64(result.TotalRewards),
    })

}

func (a *AccountRoutes) GetAccount(c *gin.Context) {
    accountAddress := c.Param("accountAddress")
    account, err := a.db.GetAccount(accountAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch account",
        })
        return
    }
    if account.Address == "" {
        c.JSON(http.StatusNotFound, gin.H{
            "status": "Not Found",
            "error":  "Account not found",
        })
        return
    }
    numberOfTransactions, err := a.db.CountTransactions(accountAddress)
    if err != nil {
        log.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch account",
        })
        return
    }
    numberOfRewards, err := a.db.CountRewards(accountAddress, -1, -1)
    if err != nil {
        log.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch account",
        })
        return
    }

    priceValue := a.priceResolver.GetPrice()
    dollarValue := int64(-1)
    if priceValue > -1 {
        dollarValue = int64(priceValue * float64(account.Balance))
    }

    c.JSON(200, &types.Account{
        Balance:  account.Balance,
        USDValue: dollarValue,
        // legacy
        BalanceDisplay:       "",
        Address:              accountAddress,
        TotalRewards:         account.TotalRewards,
        NumberOfTransactions: numberOfTransactions,
        Counter:              numberOfTransactions,
        NumberOfRewards:      numberOfRewards,
    })
}

func (a *AccountRoutes) GetAccountRewards(c *gin.Context) {
    offsetStr := c.DefaultQuery("offset", "0")
    limitStr := c.DefaultQuery("limit", "20")
    sortStr := c.DefaultQuery("sort", "asc")

    firstLayerStr := c.DefaultQuery("firstLayer", "-1")
    lastLayerStr := c.DefaultQuery("lastLayer", "-1")

    firstLayer, err := strconv.Atoi(firstLayerStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "firstLayer must be a valid integer",
        })
        return
    }

    lastLayer, err := strconv.Atoi(lastLayerStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "lastLayer must be a valid integer",
        })
        return
    }

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

    accountAddress := c.Param("accountAddress")
    rewards, errRewards := a.db.GetRewards(accountAddress, int64(offset), int64(limit), sort, firstLayer, lastLayer)
    count, errCount := a.db.CountRewards(accountAddress, firstLayer, lastLayer)

    if errRewards != nil || errCount != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch rewards for account",
        })
    } else if rewards != nil {

        rewardsResponse := make([]*types.Reward, len(rewards))

        for i, v := range rewards {
            rewardsResponse[i] = &types.Reward{
                Rewards: int64(v.TotalReward),
                // legacy
                RewardsDisplay: "",
                Layer:          v.Layer,
                SmesherId:      v.NodeId,
                // legacy
                Time:      "2023-09-05T00:00:00Z",
                Timestamp: config.GenesisEpochSeconds + (v.Layer * config.LayerDuration),
            }
        }

        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, rewardsResponse)
    } else {
        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, make([]*types.Reward, 0))
    }
}

func (a *AccountRoutes) GetAccountTransactions(c *gin.Context) {
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

    accountAddress := c.Param("accountAddress")
    transactions, errRewards := a.db.GetTransactions(accountAddress, int64(offset), int64(limit), sort, complete)
    count, errCount := a.db.CountTransactions(accountAddress)

    if errRewards != nil || errCount != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch transactions for account",
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
                Type:             v.Type,
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

func (a *AccountRoutes) GetAccountRewardsDetails(c *gin.Context) {
    accountAddress := c.Param("accountAddress")

    networkInfo := a.state.GetInfo()
    epoch := networkInfo.Epoch

    a.getAccountRewardDetailsForEpoch(c, accountAddress, int(epoch))

}

func (a *AccountRoutes) FilterEpochActiveNodes(c *gin.Context) {

    accountAddress := c.Param("accountAddress")

    epochStr := c.Param("epoch")
    epoch, err := strconv.Atoi(epochStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "epoch must be a valid integer",
        })
        return
    }

    var req types.NodeFilterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    nodes := req.Nodes

    if epoch == 8 {
        c.JSON(200, &types.ActiveNodesEpoch{
            Nodes: nodes,
        })
    } else {
        activeNodes, err := a.db.FilterAccountAtxNodesForEpoch(accountAddress, uint64(epoch-1), nodes)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to filter nodes",
            })
            return
        }

        c.JSON(200, &types.ActiveNodesEpoch{
            Nodes: activeNodes,
        })
    }

}

func (a *AccountRoutes) GetEpochAtx(c *gin.Context) {
    accountAddress := c.Param("accountAddress")

    epochStr := c.Param("epoch")
    epoch, err := strconv.Atoi(epochStr)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "epoch must be a valid integer",
        })
        return
    }

    offsetStr := c.DefaultQuery("offset", "0")
    limitStr := c.DefaultQuery("limit", "20")
    sortStr := c.DefaultQuery("sort", "asc")

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

    atxs, errAtx := a.db.GetAccountAtxEpoch(accountAddress, uint64(epoch-1), int64(offset), int64(limit), sort)
    count, errCount := a.db.CountAccountAtxEpoch(accountAddress, uint64(epoch-1))

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "failed to filter nodes",
        })
        return
    }

    if errAtx != nil || errCount != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to fetch atx for account",
        })
    } else if atxs != nil {

        atxResponse := make([]*types.Atx, len(atxs))

        for i, a := range atxs {
            atxResponse[i] = &types.Atx{
                NodeId:            a.NodeID,
                AtxId:             a.AtxID,
                EffectiveNumUnits: a.EffectiveNumUnits,
                Received:          a.Received,
            }
        }

        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, atxResponse)
    } else {
        c.Header("total", strconv.FormatInt(count, 10))
        c.JSON(200, make([]*types.Atx, 0))
    }

}

func (a *AccountRoutes) GetAccountRewardsDetailsEpoch(c *gin.Context) {
    accountAddress := c.Param("accountAddress")

    epochStr := c.Param("epoch")
    epoch, err := strconv.Atoi(epochStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "epoch must be a valid integer",
        })
        return
    }
    if epoch < 2 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "epoch should be equal or greater than 2",
        })
        return
    }

    a.getAccountRewardDetailsForEpoch(c, accountAddress, epoch)
}

func (a *AccountRoutes) getAccountRewardDetailsForEpoch(c *gin.Context, accountAddress string, epoch int) {
    epochAtx, err := a.db.GetAtxEpoch(uint64(epoch - 1))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to get atx epoch",
        })
        return
    }

    if epochAtx.TotalWeight == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "status": "Not found",
            "error":  "No details for epoch",
        })
        return
    }

    firstLayer := uint32(epoch * config.LayersPerEpoch)
    lastLayer := firstLayer + config.LayersPerEpoch

    countEpochResult, err := a.db.CountRewards(accountAddress, int(firstLayer), int(lastLayer))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to get epoch rewards count",
        })
        return
    }

    sumEpochResult, err := a.db.SumRewardsLayers(accountAddress, firstLayer, lastLayer)
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to get epoch rewards sum",
        })
        return
    }

    accountAtxs, err := a.db.GetAccountAtxList(accountAddress, uint64(epoch-1))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "Internal Error",
            "error":  "Failed to get account weight",
        })
        return
    }

    eligibilityCount := int32(0)
    totalWeight := uint64(0)
    totalEffectiveNumUnits := uint32(0)
    for _, atx := range accountAtxs {
        eligibilityCountTemp, err := a.networkUtils.GetNumberOfSlots(uint64(atx.Weight), epochAtx.TotalWeight, uint32(epoch))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "Internal Error",
                "error":  "Failed to get eligibility",
            })
            return
        }
        eligibilityCount += eligibilityCountTemp
        totalWeight += atx.Weight
        totalEffectiveNumUnits += atx.EffectiveNumUnits
    }

    if totalWeight == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "status": "Not found",
            "error":  "Account not active for epoch",
        })
        return
    }

    unitReward := a.state.GetEpochSubsidy(uint32(epoch)) / epochAtx.TotalWeight
    predictedRewards := unitReward * uint64(totalWeight)

    c.JSON(200, &types.RewardDetailsEpoch{
        Epoch:        int64(epoch),
        RewardsSum:   sumEpochResult,
        RewardsCount: countEpochResult,
        Eligibility: &types.Eligibility{
            Count:             eligibilityCount,
            EffectiveNumUnits: int64(totalEffectiveNumUnits),
            PredictedRewards:  uint64(predictedRewards),
        },
    })
}
