package agora

import "strings"

type IDName struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

type MethodAmount struct {
	MethodName string  `json:"MethodName"`
	Amount     float32 `json:"Amount"`
}

type Day struct {
	Invoices        []Invoice        `json:"Invoices"`
	PosCloseOuts    []PosCloseOut    `json:"PosCloseOuts"`
	SystemCloseOuts []SystemCloseOut `json:"SystemCloseOuts"`
}

type Invoice struct {
	Serie         string           `json:"Serie"`
	Number        int              `json:"Number"`
	BusinessDay   string           `json:"BusinessDay"`
	VatIncluded   bool             `json:"VatIncluded"`
	Date          string           `json:"Date"`
	Pos           IDName           `json:"Pos"`
	Workplace     IDName           `json:"Workplace"`
	User          IDName           `json:"User"`
	DocumentType  string           `json:"DocumentType"`
	InvoiceItems  []InvoiceItem    `json:"InvoiceItems"`
	Payments      []InvoicePayment `json:"Payments"`
	Totals        InvoiceTotals    `json:"Totals"`
	TicketBAIData string           `json:"TicketBAIData"`
}

type InvoiceItem struct {
	ContentType string            `json:"ContentType"`
	Pos         IDName            `json:"Pos"`
	User        IDName            `json:"User"`
	GlobalID    string            `json:"GlobalId"`
	BusinessDay string            `json:"BusinessDay"`
	PriceList   IDName            `json:"PriceList"`
	Date        string            `json:"Date"`
	Lines       []InvoiceItemLine `json:"Lines"`
	Discounts   InvoiceDiscount   `json:"Discounts"`
	Payments    []any             `json:"Payments"`
	Offers      []any             `json:"Offers"`
	VatIncluded bool              `json:"VatIncluded"`
	Totals      InvoiceTotals     `json:"Totals"`
}

type InvoiceItemLine struct {
	Index         int     `json:"Index"`
	CreationDate  string  `json:"CreationDate"`
	UserID        int     `json:"UserId"`
	ProductID     int     `json:"ProductId"`
	ProductName   string  `json:"ProductName"`
	ProductPrice  float32 `json:"ProductPrice"`
	FamilyID      int     `json:"FamilyId"`
	FamilyName    string  `json:"FamilyName"`
	VatID         int     `json:"VatId"`
	VatRate       float32 `json:"VatRate"`
	SurchargeRate float32 `json:"SurchargeRate"`
	Quantity      float32 `json:"Quantity"`
	UnitPrice     float32 `json:"UnitPrice"`
	DiscountRate  float32 `json:"DiscountRate"`
	CashDiscount  float32 `json:"CashDiscount"`
	TotalAmount   float32 `json:"TotalAmount"`
	UnitCostPrice float32 `json:"UnitCostPrice"`
	OfferID       any     `json:"OfferId"`
	OfferCode     string  `json:"OfferCode"`
	Notes         string  `json:"Notes"`
}

type InvoiceTotals struct {
	GrossAmount     float32      `json:"GrossAmount"`
	NetAmount       float32      `json:"NetAmount"`
	VatAmount       float32      `json:"VatAmount"`
	SurchargeAmount float32      `json:"SuperchargeAmount"`
	Taxes           []InvoiceTax `json:"Taxes"`
}

type InvoicePayment struct {
	MethodID         int     `json:"MethodId"`
	MethodName       string  `json:"MethodName"`
	Amount           float32 `json:"Amount"`
	PaidAmount       float32 `json:"PaidAmount"`
	ChangeAmount     float32 `json:"ChangeAmount"`
	Date             string  `json:"Date"`
	PosID            int     `json:"PosId"`
	IsPrepayment     bool    `json:"IsPrepayment"`
	ExtraInformation string  `json:"ExtraInformation"`
}

type InvoiceTax struct {
	VatRate         float32 `json:"VatRate"`
	SurchargeRate   float32 `json:"SurchargeRate"`
	NetAmount       float32 `json:"NetAmount"`
	VatAmount       float32 `json:"VatAmount"`
	SurchargeAmount float32 `json:"SurchargeAmount"`
}

type InvoiceDiscount struct {
	DiscountRate float32 `json:"DiscountRate"`
	CashDiscount float32 `json:"CashDiscount"`
}

type PosCloseOut struct {
	ID                int                  `json:"Id"`
	PosID             int                  `json:"PosId"`
	WorkplaceID       int                  `json:"WorkplaceId"`
	BusinessDay       string               `json:"BusinessDay"`
	InitialAmount     float32              `json:"InitialAmount"`
	ExpectedEndAmount float32              `json:"ExpectedEndAmount"`
	ActualEndAmount   float32              `json:"ActualEndAmount"`
	Incident          string               `json:"Incident"`
	OpenDate          string               `json:"OpenDate"`
	OpenerUserId      int                  `json:"OpenerUserId"`
	CloseDate         string               `json:"CloseDate"`
	CloserUserId      int                  `json:"CloserUserId"`
	VerificationCode  string               `json:"VerificationCode"`
	Balances          []PosCloseOutBalance `json:"Balances"`
}

type PosCloseOutBalance struct {
	PaymentMethodID   int     `json:"PaymentMethodId"`
	InitialAmount     float32 `json:"InitialAmount"`
	ExpectedEndAmount float32 `json:"ExpectedEndAmount"`
	ActualEndAmount   float32 `json:"ActualEndAmount"`
}

type SystemCloseOut struct {
	Number               int                      `json:"Number"`
	BusinessDay          string                   `json:"BusinessDay"`
	OpenDate             string                   `json:"OpenDate"`
	CloseDate            string                   `json:"CloseDate"`
	OpenerUserId         int                      `json:"OpenerUserId"`
	CloserUserId         int                      `json:"CloserUserId"`
	WorkplaceID          int                      `json:"WorkplaceId"`
	Documents            []SystemCloseOutDocument `json:"Documents"`
	Amounts              SystemCloseOutAmount     `json:"Amounts"`
	InvoicePayments      []MethodAmount           `json:"InvoicePayments"`
	SalesOrderPayments   []MethodAmount           `json:"SalesOrderPayments"`
	DeliveryNotePayments []MethodAmount           `json:"DeliveryNotePayments"`
	TicketPayment        []MethodAmount           `json:"TicketPayment"`
}

type SystemCloseOutDocument struct {
	Serie       string  `json:"Serie"`
	Amount      float32 `json:"Amount"`
	FirstNumber int     `json:"FirstNumber"`
	LastNumber  int     `json:"LastNumber"`
	Count       int     `json:"Count"`
}

type SystemCloseOutAmount struct {
	NetAmount       float32 `json:"NetAmount"`
	GrossAmount     float32 `json:"GrossAmount"`
	SurchargeAmount float32 `json:"SurchargeAmount"`
	VatAmount       float32 `json:"VatAmount"`
}

type Master struct {
	WorkplacesSummary []WorkplacesSummary `json:"WorkplacesSummary"`
	Vats              []Vat               `json:"Vats"`
	Series            []Serie             `json:"Series"`
	PriceLists        []PriceList         `json:"PriceLists"`
	Users             []User              `json:"Users"`
	PaymentMethods    []PaymentMethod     `json:"PaymentMethods"`
	Warehouses        []Warehouse         `json:"Warehouses"`
	Customers         []Customer          `json:"Customers"`
	Families          []Family            `json:"Families"`
	Products          []Product           `json:"Products"`
	Stocks            []Stock             `json:"Stocks"`
	Suppliers         []Supplier          `json:"Suppliers"`
}

type WorkplacesSummary struct {
}

type Vat struct {
	ID            int     `json:"Id"`
	Name          string  `json:"Name"`
	VatRate       float32 `json:"VatRate"`
	SurchargeRate float32 `json:"SurchargeRate"`
	Enabled       bool    `json:"Enabled"`
}

type Serie struct {
	Name         SerieName         `json:"Name"`
	LastNumber   int               `json:"LastNumber"`
	DocumentType SerieDocumentType `json:"DocumentType"`
}

type SerieName string

const (
	SerieNameBasicInvoice    SerieName = "T"
	SerieNameStandardInvoice SerieName = "F"
	SerieNameBasicRefund     SerieName = "TD"
	SerieNameStandardRefund  SerieName = "FD"
	SerieNameDeliveryNote    SerieName = "A"
	SerieNameSalesOrder      SerieName = "P"
)

type SerieDocumentType string

const (
	SerieDocumentTypeBasicInvoice    SerieDocumentType = "BasicInvoice"
	SerieDocumentTypeStandardInvoice SerieDocumentType = "StandardInvoice"
	SerieDocumentTypeBasicRefund     SerieDocumentType = "BasicRefund"
	SerieDocumentTypeStandardRefund  SerieDocumentType = "StandardRefund"
	SerieDocumentTypeDeliveryNote    SerieDocumentType = "DeliveryNote"
	SerieDocumentTypeSalesOrder      SerieDocumentType = "SalesOrder"
)

type PriceList struct {
	ID          int    `json:"Id"`
	Name        string `json:"Name"`
	VatIncluded bool   `json:"VatIncluded"`
}

type User struct {
}

type PaymentMethod struct {
}

type Warehouse struct {
}

type Customer struct {
}

type Family struct {
}

type Product struct {
	ID                   int                `json:"Id"`
	Name                 string             `json:"Name"`
	ButtonText           string             `json:"ButtonText"`
	Color                string             `json:"Color"`
	PLU                  string             `json:"PLU"`
	FamilyID             int                `json:"FamilyId"`
	VatID                int                `json:"VatId"`
	UseAsDirectSale      bool               `json:"UseAsDirectSale"`
	Saleable             bool               `json:"Saleable"`
	IsSoldByWeight       bool               `json:"IsSoldByWeight"`
	PrintWhenPriceIsZero bool               `json:"PrintWhenPriceIsZero"`
	SizeGroupID          *int               `json:"SizeGroupId"`
	ColorGroupID         *int               `json:"ColorGroupId"`
	CostPrice            float32            `json:"CostPrice"`
	Barcodes             []ProductBarcode   `json:"Barcodes"`
	StorageOptions       []ProductStorage   `json:"StorageOptions"`
	Prices               []ProductPrice     `json:"Prices"`
	CostPrices           []ProductCostPrice `json:"CostPrices"`
	DeletionDate         string             `json:"DeletionDate"`
}

func (p Product) Barcode() string {
	if len(p.Barcodes) == 0 {
		return ""
	}
	barcode := p.Barcodes[0].Value
	for _, c := range []string{" ", "-", "_", ".", ",", ";"} {
		barcode = strings.ReplaceAll(barcode, c, "")
	}
	return barcode
}

type ProductBarcode struct {
	Value string `json:"Value"`
}

type ProductStorage struct {
	WarehouseID int     `json:"WarehouseId"`
	Location    string  `json:"Location"`
	MinStock    float32 `json:"MinStock"`
	MaxStock    float32 `json:"MaxStock"`
}

type ProductPrice struct {
	PriceListID int     `json:"PriceListId"`
	Price       float32 `json:"Price"`
}

type ProductCostPrice struct {
	WarehouseID int     `json:"WarehouseId"`
	CostPrice   float32 `json:"CostPrice"`
}

type Stock struct {
	WarehouseID int     `json:"WarehouseId"`
	ProductID   int     `json:"ProductId"`
	Quantity    float32 `json:"Quantity"`
}

type Supplier struct {
}
