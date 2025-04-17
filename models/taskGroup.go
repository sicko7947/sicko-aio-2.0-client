package models

import "github.com/gogf/gf/os/gmutex"

type TaskGroupID string
type TaskID string

type TaskGroupSetting struct {
	Name       string   `json:"name"`
	Category   CATEGORY `json:"category"`
	Country    string   `json:"country"`
	Domain     string   `json:"domain,omitempty"`
	Currency   string   `json:"currency,omitempty"`
	Locale     string   `json:"locale,omitempty"`
	MerchGroup string   `json:"merchGroup,omitempty"`
	Language   string   `json:"language,omitempty"`
	TaskType   string   `json:"taskType"`
}

type WorkerTaskProduct struct {
	StyleColor         string   `json:"styleColor,omitempty"`
	ProductName        string   `json:"productName,omitempty"`
	ProductDescription string   `json:"productDescription,omitempty"`
	QuantityLimit      int      `json:"quantityLimit,omitempty"`
	ProductID          string   `json:"productId,omitempty"`
	LaunchID           string   `json:"launchId,omitempty"`
	SkuID              string   `json:"skuId,omitempty"`
	Gtin               string   `json:"gtin,omitempty"`
	Size               string   `json:"size,omitempty"`
	Price              string   `json:"price,omitempty"`
	ImageURL           string   `json:"imageUrl,omitempty"`
	Other              []string `json:"other,omitempty"`
}

type WorkerTaskInfo struct {
	Email string `json:"email,omitemtpy"` // email used for placing the order if user choose to use domain

	XB3TracedID       string `json:"xb3tracedid,omitempty"`
	XNikeVisitorID    string `json:"xnikevisitorid,omitempty"`
	Authorization     string `json:"authorization,omitempty"`
	CreditCardToken   string `json:"creditCardToken,omitempty"`
	PaymentToken      string `json:"paymentToken,omitempty"`
	CheckoutID        string `json:"checkoutId,omitempty"`
	RedirectURL       string `json:"redirectUrl,omitempty"`
	OrderID           string `json:"orderId,omitempty"`
	PaymentApprovalID string `json:"paymentApprovalId,omitempty"`
	FulfillmentID     string `json:"fulfillmentId,omitempty"`
	MinDate           string `json:"minDate,omitempty"`
	MaxDate           string `json:"maxDate,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	Cookies           string `json:"cookies,omitempty"`
	Phone             string `json:"phone,omitempty"`
	Token             string `json:"token,omitempty"`

	Discount      string      `json:"discount,omitempty"`
	Profile       *Profile    `json:"profile,omitempty"`
	Account       *Account    `json:"account,omitempty"`
	GiftCardGroup []*GiftCard `json:"giftCardGroup,omitempty"`
}

type SizeSkuMap struct {
	SkuId string `json:"skuId,omitempty"`
	Gtin  string `json:"gtin,omitempty"`
}
type ScraperTaskProduct struct {
	StyleColor string `json:"styleColor,omitempty"`
	Keyword    string `json:"keyword,omitempty"`

	ProductName        string                 `json:"productName,omitempty"`
	ProductURL         string                 `json:"productUrl,omitempty"`
	LaunchID           string                 `json:"launchId,omitempty"`
	QuantityLimit      int                    `json:"quantityLimit,omitempty"`
	ProductDescription string                 `json:"productDescription,omitempty"`
	ProductID          string                 `json:"productId,omitempty"`
	PublishType        string                 `json:"publishType,omitempty"`
	Price              string                 `json:"price,omitempty"`
	SizeSkuMap         map[string]*SizeSkuMap `json:"sizeSkuMap,omitempty"`
	CountryExclustion  []string               `json:"countryExclusion,omitempty"`
	CommerceStartDate  string                 `json:"commerceStartDate,omitempty"`
	ImageURL           string                 `json:"imageUrl,omitempty"`
	IsLive             bool                   `json:"isLive,omitempty"`
	Other              []string               `json:"other,omitempty"`

	Status  string `json:"status,omitempty"`
	Scrapes int64  `json:"scrapes,omitempty"`

	Size  string `json:"size,omitempty"`  // for normal task invoke
	SkuID string `json:"skuId,omitempty"` // for normal task invoke
}

type TaskWorker struct {
	GroupID TaskGroupID   `json:"groupId,omitempty"`
	TaskID  TaskID        `json:"taskId,omitempty"`
	Mutex   *gmutex.Mutex `json:"muxtex,omitempty"`

	StartTime string `json:"startTime,omitempty"`
	Quantity  int    `json:"quantity"`
	MaxCart   bool   `json:"maxCart"`
	Restart   bool   `json:"restart"`

	Product                *WorkerTaskProduct `json:"product,omitempty"`
	TaskInfo               *WorkerTaskInfo    `json:"taskInfo,omitempty"`
	ProfileGroupName       ProfileGroupName   `json:"profileGroupName"`
	CheckoutProxyGroupName ProxyGroupName     `json:"checkoutProxyGroupName"`
	AccountGroupName       AccountGroupName   `json:"accountGroupName,omitempty"`
	DiscountGroupName      DiscountGroupName  `json:"discountGroupName,omitempty"`
	GiftCardGroupName      GiftCardGroupName  `json:"giftCardGroupName,omitempty"`
	DesireSizes            []string           `json:"desireSizes,omitempty"`
}

type TaskScraper struct {
	GroupID TaskGroupID `json:"groupId,omitempty"`
	TaskID  TaskID      `json:"taskId"`
	Mutex   *gmutex.Mutex

	DesireSizes  []string `json:"desireSizes"`
	MonitorMode  string   `json:"monitorMode,omitempty"`
	ScraperNum   int      `json:"scraperNum"`
	ScraperDelay int      `json:"scraperDelay"`
	TriggerNum   int      `json:"triggerNum"`
	TriggerDelay int      `json:"triggerDelay"`
	CheckLive    bool     `json:"checkLive"`

	Product               *ScraperTaskProduct `json:"product"`
	ScraperProxyGroupName ProxyGroupName      `json:"scraperProxyGroupName"`
}

type TaskGroup struct {
	GroupID          TaskGroupID             `json:"groupId"`
	TaskGroupSetting *TaskGroupSetting       `json:"taskGroupSetting"`
	TaskWorkers      map[TaskID]*TaskWorker  `json:"taskWorkers,omitempty"`
	TaskScrapers     map[TaskID]*TaskScraper `json:"taskScrapers,omitempty"`
}
