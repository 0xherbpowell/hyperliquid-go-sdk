package types

// LimitOrderType represents a limit order
type LimitOrderType struct {
	Tif Tif `json:"tif"`
}

// TriggerOrderType represents a trigger order
type TriggerOrderType struct {
	TriggerPx float64 `json:"triggerPx"`
	IsMarket  bool    `json:"isMarket"`
	Tpsl      Tpsl    `json:"tpsl"`
}

// TriggerOrderTypeWire represents a trigger order for wire format
type TriggerOrderTypeWire struct {
	TriggerPx FloatString `json:"triggerPx"`
	IsMarket  bool        `json:"isMarket"`
	Tpsl      Tpsl        `json:"tpsl"`
}

// OrderType represents the type of order (limit or trigger)
type OrderType struct {
	Limit   *LimitOrderType   `json:"limit,omitempty"`
	Trigger *TriggerOrderType `json:"trigger,omitempty"`
}

// OrderTypeWire represents the order type in wire format
type OrderTypeWire struct {
	Limit   *LimitOrderType       `json:"limit,omitempty"`
	Trigger *TriggerOrderTypeWire `json:"trigger,omitempty"`
}

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Coin       string    `json:"coin"`
	IsBuy      bool      `json:"is_buy"`
	Sz         float64   `json:"sz"`
	LimitPx    float64   `json:"limit_px"`
	OrderType  OrderType `json:"order_type"`
	ReduceOnly bool      `json:"reduce_only"`
	Cloid      *Cloid    `json:"cloid,omitempty"`
}

// OidOrCloid represents either an order ID (int) or client order ID (Cloid)
type OidOrCloid interface {
	IsOid() bool
	GetOid() int
	GetCloid() *Cloid
}

// Oid represents an order ID
type Oid int

func (o Oid) IsOid() bool      { return true }
func (o Oid) GetOid() int      { return int(o) }
func (o Oid) GetCloid() *Cloid { return nil }

func (c *Cloid) IsOid() bool      { return false }
func (c *Cloid) GetOid() int      { return 0 }
func (c *Cloid) GetCloid() *Cloid { return c }

// ModifyRequest represents a request to modify an order
type ModifyRequest struct {
	Oid   OidOrCloid   `json:"oid"`
	Order OrderRequest `json:"order"`
}

// CancelRequest represents a request to cancel an order
type CancelRequest struct {
	Coin string `json:"coin"`
	Oid  int    `json:"oid"`
}

// CancelByCloidRequest represents a request to cancel an order by client order ID
type CancelByCloidRequest struct {
	Coin  string `json:"coin"`
	Cloid *Cloid `json:"cloid"`
}

// OrderWire represents an order in wire format for API communication
type OrderWire struct {
	A int           `json:"a"`           // Asset
	B bool          `json:"b"`           // Is buy
	P FloatString   `json:"p"`           // Price
	S FloatString   `json:"s"`           // Size
	R bool          `json:"r"`           // Reduce only
	T OrderTypeWire `json:"t"`           // Order type
	C *string       `json:"c,omitempty"` // Client order ID
}

// ModifyWire represents a modify request in wire format
type ModifyWire struct {
	Oid   interface{} `json:"oid"` // Can be int or string
	Order OrderWire   `json:"order"`
}

// Grouping represents order grouping type
type Grouping string

const (
	GroupingNa           Grouping = "na"
	GroupingNormalTpsl   Grouping = "normalTpsl"
	GroupingPositionTpsl Grouping = "positionTpsl"
)

// OrderAction represents the action to place orders
type OrderAction struct {
	Type     string       `json:"type"` // "order"
	Orders   []OrderWire  `json:"orders"`
	Grouping Grouping     `json:"grouping"`
	Builder  *BuilderInfo `json:"builder,omitempty"`
}

// CancelAction represents the action to cancel orders
type CancelAction struct {
	Type    string `json:"type"` // "cancel"
	Cancels []struct {
		A int `json:"a"` // Asset
		O int `json:"o"` // Order ID
	} `json:"cancels"`
}

// CancelByCloidAction represents the action to cancel orders by client order ID
type CancelByCloidAction struct {
	Type    string `json:"type"` // "cancelByCloid"
	Cancels []struct {
		Asset int    `json:"asset"`
		Cloid string `json:"cloid"`
	} `json:"cancels"`
}

// BatchModifyAction represents the action to modify multiple orders
type BatchModifyAction struct {
	Type     string       `json:"type"` // "batchModify"
	Modifies []ModifyWire `json:"modifies"`
}

// ScheduleCancelAction represents the action to schedule order cancellation
type ScheduleCancelAction struct {
	Type string `json:"type"`           // "scheduleCancel"
	Time *int64 `json:"time,omitempty"` // Optional cancel time in milliseconds
}

// UpdateLeverageAction represents the action to update leverage
type UpdateLeverageAction struct {
	Type     string `json:"type"` // "updateLeverage"
	Asset    int    `json:"asset"`
	IsCross  bool   `json:"isCross"`
	Leverage int    `json:"leverage"`
}

// UpdateIsolatedMarginAction represents the action to update isolated margin
type UpdateIsolatedMarginAction struct {
	Type  string `json:"type"` // "updateIsolatedMargin"
	Asset int    `json:"asset"`
	IsBuy bool   `json:"isBuy"`
	Ntli  int64  `json:"ntli"` // Amount in integer format
}

// SetReferrerAction represents the action to set a referrer
type SetReferrerAction struct {
	Type string `json:"type"` // "setReferrer"
	Code string `json:"code"`
}

// CreateSubAccountAction represents the action to create a sub-account
type CreateSubAccountAction struct {
	Type string `json:"type"` // "createSubAccount"
	Name string `json:"name"`
}

// OrderStatus represents the status of an order
type OrderStatus struct {
	Order           OrderRequest `json:"order"`
	Status          string       `json:"status"`
	StatusTimestamp int64        `json:"statusTimestamp"`
}

// OrderResponse represents the response from placing an order
type OrderResponse struct {
	Status   string `json:"status"`
	Response struct {
		Type string `json:"type"`
		Data struct {
			Statuses []struct {
				Resting *struct {
					Oid int `json:"oid"`
				} `json:"resting,omitempty"`
				Filled *struct {
					Oid     int         `json:"oid"`
					TotalSz FloatString `json:"totalSz"`
					AvgPx   FloatString `json:"avgPx"`
				} `json:"filled,omitempty"`
				Error *string `json:"error,omitempty"`
			} `json:"statuses"`
		} `json:"data"`
	} `json:"response"`
}
