package constants

const (
	MainnetAPIURL = "https://api.hyperliquid.xyz"
	TestnetAPIURL = "https://api.hyperliquid-testnet.xyz"
	LocalAPIURL   = "http://localhost:3001"
)

const (
	DefaultSlippage = 0.05 // 5%
	MaxDecimals     = 6    // For perps
	SpotMaxDecimals = 8    // For spot
)

const (
	SignatureChainID = "0x66eee"
	MainnetChain     = "Mainnet"
	TestnetChain     = "Testnet"
)

const (
	ExchangeDomain    = "Exchange"
	HyperliquidDomain = "HyperliquidSignTransaction"
	DomainVersion     = "1"
	VerifyingContract = "0x0000000000000000000000000000000000000000"
	EIP712ChainID     = 1337
)

// Order sides
const (
	SideAsk = "A"
	SideBid = "B"
)

// Time in force
const (
	TifAlo = "Alo" // Add liquidity only
	TifIoc = "Ioc" // Immediate or cancel
	TifGtc = "Gtc" // Good till cancelled
)

// TP/SL types
const (
	TpslTp = "tp" // Take profit
	TpslSl = "sl" // Stop loss
)

// Leverage types
const (
	LeverageCross    = "cross"
	LeverageIsolated = "isolated"
)

// Action types
const (
	ActionOrder                  = "order"
	ActionCancel                 = "cancel"
	ActionCancelByCloid          = "cancelByCloid"
	ActionBatchModify            = "batchModify"
	ActionScheduleCancel         = "scheduleCancel"
	ActionUpdateLeverage         = "updateLeverage"
	ActionUpdateIsolatedMargin   = "updateIsolatedMargin"
	ActionSetReferrer            = "setReferrer"
	ActionCreateSubAccount       = "createSubAccount"
	ActionUsdClassTransfer       = "usdClassTransfer"
	ActionSendAsset              = "sendAsset"
	ActionSubAccountTransfer     = "subAccountTransfer"
	ActionSubAccountSpotTransfer = "subAccountSpotTransfer"
	ActionVaultTransfer          = "vaultTransfer"
	ActionUsdSend                = "usdSend"
	ActionSpotSend               = "spotSend"
	ActionTokenDelegate          = "tokenDelegate"
	ActionWithdraw3              = "withdraw3"
	ActionApproveAgent           = "approveAgent"
	ActionApproveBuilderFee      = "approveBuilderFee"
	ActionConvertToMultiSigUser  = "convertToMultiSigUser"
	ActionSpotDeploy             = "spotDeploy"
	ActionPerpDeploy             = "perpDeploy"
	ActionMultiSig               = "multiSig"
	ActionEvmUserModify          = "evmUserModify"
	ActionNoop                   = "noop"
)

// Info request types
const (
	InfoClearinghouseState          = "clearinghouseState"
	InfoSpotClearinghouseState      = "spotClearinghouseState"
	InfoOpenOrders                  = "openOrders"
	InfoFrontendOpenOrders          = "frontendOpenOrders"
	InfoAllMids                     = "allMids"
	InfoUserFills                   = "userFills"
	InfoUserFillsByTime             = "userFillsByTime"
	InfoMeta                        = "meta"
	InfoMetaAndAssetCtxs            = "metaAndAssetCtxs"
	InfoPerpDexs                    = "perpDexs"
	InfoSpotMeta                    = "spotMeta"
	InfoSpotMetaAndAssetCtxs        = "spotMetaAndAssetCtxs"
	InfoFundingHistory              = "fundingHistory"
	InfoUserFunding                 = "userFunding"
	InfoL2Book                      = "l2Book"
	InfoCandleSnapshot              = "candleSnapshot"
	InfoUserFees                    = "userFees"
	InfoDelegatorSummary            = "delegatorSummary"
	InfoDelegations                 = "delegations"
	InfoDelegatorRewards            = "delegatorRewards"
	InfoDelegatorHistory            = "delegatorHistory"
	InfoOrderStatus                 = "orderStatus"
	InfoReferral                    = "referral"
	InfoSubAccounts                 = "subAccounts"
	InfoUserToMultiSigSigners       = "userToMultiSigSigners"
	InfoPerpDeployAuctionStatus     = "perpDeployAuctionStatus"
	InfoHistoricalOrders            = "historicalOrders"
	InfoUserNonFundingLedgerUpdates = "userNonFundingLedgerUpdates"
	InfoPortfolio                   = "portfolio"
	InfoUserTwapSliceFills          = "userTwapSliceFills"
	InfoUserVaultEquities           = "userVaultEquities"
	InfoUserRole                    = "userRole"
	InfoUserRateLimit               = "userRateLimit"
	InfoSpotDeployState             = "spotDeployState"
	InfoExtraAgents                 = "extraAgents"
)

// WebSocket subscription types
const (
	WsAllMids                     = "allMids"
	WsL2Book                      = "l2Book"
	WsTrades                      = "trades"
	WsUserEvents                  = "userEvents"
	WsUserFills                   = "userFills"
	WsCandle                      = "candle"
	WsOrderUpdates                = "orderUpdates"
	WsUserFundings                = "userFundings"
	WsUserNonFundingLedgerUpdates = "userNonFundingLedgerUpdates"
	WsWebData2                    = "webData2"
	WsBbo                         = "bbo"
	WsActiveAssetCtx              = "activeAssetCtx"
	WsActiveAssetData             = "activeAssetData"
)
