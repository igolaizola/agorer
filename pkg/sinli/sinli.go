package sinli

import "time"

// Subject to be used in the email where the sinli file is attached.
// Example: `ESFANDELXXXXXXXESFANDELXXXXXXXLIBROSNNFANDE`
type Subject struct {
	//lint:ignore U1000 Fixed value
	sourcePrefix string `sinli:"order=1,length=7,fixed=ESFANDE"`
	SourceID     string `sinli:"order=2,length=8"`
	//lint:ignore U1000 Fixed value
	destinationPrefix string      `sinli:"order=3,length=7,fixed=ESFANDE"`
	DestinationID     string      `sinli:"order=4,length=8"`
	FileType          FileType    `sinli:"order=5,length=6"`
	FileVersion       FileVersion `sinli:"order=6,length=2"`
	//lint:ignore U1000 Fixed value
	suffix string `sinli:"order=7,length=5,fixed=ESFANDE"`
}

// IdentificationHeader is the first line of a sinli file.
type IdentificationHeader struct {
	_        struct{}    `sinli:"order=1,length=1,fixed=I"`
	Format   FormatType  `sinli:"order=2,length=1"`
	Document FileType    `sinli:"order=3,length=6"`
	Version  FileVersion `sinli:"order=4,length=2"`
	// Format `Lnnnnnnn`
	SourceID string `sinli:"order=5,length=8"`
	// Format `Lnnnnnnn`
	DestinationID string `sinli:"order=6,length=8"`
	// Optionals
	Records            int    `sinli:"order=7,length=5"`
	TransmissionNumber int    `sinli:"order=8,length=7"`
	LocalSourceID      string `sinli:"order=9,length=15"`
	LocalDestinationID string `sinli:"order=10,length=15"`
	FreeText           string `sinli:"order=11,length=7"`
	//lint:ignore U1000 Fixed value
	suffix string `sinli:"order=12,length=5,fixed=FANDE"`
}

type FormatType string

const (
	FormatTypeNormalized FormatType = "N"
	FormatTypeFree       FormatType = "L"
)

// Identification is the second line of a sinli file.
type Identification struct {
	_                struct{}    `sinli:"order=1,length=1,fixed=I"`
	SourceEmail      string      `sinli:"order=2,length=50"`
	DestinationEmail string      `sinli:"order=3,length=50"`
	FileType         FileType    `sinli:"order=4,length=6"`
	FileVersion      FileVersion `sinli:"order=5,length=2"`
	// Optionals
	TransmissionNumber int `sinli:"order=6,length=8"`
}

type FileVersion int

const (
	FileVersionStock  FileVersion = 2
	FileVersionOrder  FileVersion = 7
	FileVersionReturn FileVersion = 2
	FileVersionSale   FileVersion = 3
)

type FileType string

const (
	FileTypeStock  FileType = "CEGALD"
	FileTypeOrder  FileType = "PEDIDO"
	FileTypeReturn FileType = "DEVOLU"
	FileTypeSale   FileType = "CEGALV"
)

// Order is a sinli order. Code: `PEDIDO`
type Order struct {
	IdentificationHeader IdentificationHeader `sinli:"order=1"`
	Identification       Identification       `sinli:"order=2"`
	Header               OrderHeader          `sinli:"order=3"`
	Details              []OrderDetail        `sinli:"order=4"`
}

type OrderHeader struct {
	_            struct{}   `sinli:"order=1,length=1,fixed=C"`
	ClientName   string     `sinli:"order=2,length=40"`
	ProviderName string     `sinli:"order=3,length=40"`
	OrderDate    time.Time  `sinli:"order=4,length=8"`
	OrderCode    string     `sinli:"order=5,length=10"`
	OrderType    OrderType  `sinli:"order=6,length=1"`
	Coin         LegacyCoin `sinli:"order=7,length=1"`

	// Optionals
	PrintOnDemand           bool       `sinli:"order=8,length=1"`
	RequestedDeliveryDate   *time.Time `sinli:"order=9,length=8"`
	LastDeliveryAllowedDate *time.Time `sinli:"order=10,length=8"`
	LastDeliveryExpiration  bool       `sinli:"order=11,length=1"`

	Batch string `sinli:"order=12,length=15"`
}

type OrderType string

const (
	OrderTypeNormal  OrderType = "N"
	OrderTypeFair    OrderType = "F"
	OrderTypeDeposit OrderType = "D"
	OrderTypeOther   OrderType = "O"
)

type LegacyCoin string

const (
	LegacyCoinEuro   LegacyCoin = "E"
	LegacyCoinPeseta LegacyCoin = "P"
)

type OrderDeliveryPoint struct {
	_          struct{} `sinli:"order=1,length=1,fixed=E"`
	Name       string   `sinli:"order=2,length=50"`
	Address    string   `sinli:"order=3,length=80"`
	PostalCode string   `sinli:"order=4,length=5"`
	City       string   `sinli:"order=5,length=50"`
	Province   string   `sinli:"order=6,length=40"`
}

type OrderDetail struct {
	_            struct{}    `sinli:"order=1,length=1,fixed=D"`
	ISBN         string      `sinli:"order=2,length=17"`
	EAN          string      `sinli:"order=3,length=18"`
	Reference    string      `sinli:"order=4,length=15"`
	Title        string      `sinli:"order=5,length=50"`
	Quantity     int         `sinli:"order=6,length=6"`
	PriceWithVAT float32     `sinli:"order=7,length=10"`
	WantPending  bool        `sinli:"order=8,length=1"`
	OrderSource  OrderSource `sinli:"order=9,length=1"`
	FastDelivery bool        `sinli:"order=10,length=1"`
	Code         string      `sinli:"order=11,length=10"`
}

type OrderSource string

const (
	OrderSourceNormal OrderSource = "N"
	OrderSourceClient OrderSource = "C"
)

// Return is a sinli return. Code: `DEVOLU`
type Return struct {
	IdentificationHeader IdentificationHeader `sinli:"order=1"`
	Identification       Identification       `sinli:"order=2"`
	Header               ReturnHeader         `sinli:"order=3"`
	Details              []ReturnDetail       `sinli:"order=4"`
}

type ReturnHeader struct {
	_            struct{}           `sinli:"order=1,length=1,fixed=C"`
	ClientName   string             `sinli:"order=2,length=40"`
	ProviderName string             `sinli:"order=3,length=40"`
	OrderCode    string             `sinli:"order=4,length=10"`
	DocumentDate time.Time          `sinli:"order=5,length=8"`
	DocumentType ReturnDocumentType `sinli:"order=6,length=1"`
	ReturnType   ReturnType         `sinli:"order=7,length=1"`
	BookFair     bool               `sinli:"order=8,length=1"`
	Coin         LegacyCoin         `sinli:"order=9,length=1"`
}

type ReturnDocumentType string

const (
	ReturnDocumentTypeDefinitive ReturnDocumentType = "D"
	ReturnDocumentTypeRequested  ReturnDocumentType = "P"
)

type ReturnType string

const (
	ReturnTypeDefinitive ReturnType = "F"
	ReturnTypeDeposit    ReturnType = "D"
)

type ReturnDetail struct {
	_               struct{}  `sinli:"order=1,length=1,fixed=D"`
	ISBN            string    `sinli:"order=2,length=17"`
	EAN             string    `sinli:"order=3,length=18"`
	Reference       string    `sinli:"order=4,length=15"`
	Title           string    `sinli:"order=5,length=50"`
	Quantity        int       `sinli:"order=6,length=6"`
	PriceWithoutVAT float32   `sinli:"order=7,length=10"`
	PriceWithVAT    float32   `sinli:"order=8,length=10"`
	Discount        float32   `sinli:"order=9,length=10"`
	PriceType       PriceType `sinli:"order=10,length=1"`

	// Optionals
	Novelty          bool         `sinli:"order=11,length=1"`
	PurchaseDocument string       `sinli:"order=12,length=10"`
	PurchaseDate     *time.Time   `sinli:"order=13,length=8"`
	ReturnCause      *ReturnCause `sinli:"order=14,length=1"`
}

type PriceType string

const (
	PriceTypeFixed PriceType = "F"
	PriceTypeFree  PriceType = "L"
)

type ReturnCause string

const (
	ReturnCauseDamaged  ReturnCause = "0"
	ReturnCauseOutdated ReturnCause = "1"
	ReturnCauseIncident ReturnCause = "2"
)

// Stock is a sinli stock. Code: `CEGALD`
type Stock struct {
	IdentificationHeader IdentificationHeader `sinli:"order=1"`
	Identification       Identification       `sinli:"order=2"`
	Header               StockHeader          `sinli:"order=3"`
	Details              []StockDetail        `sinli:"order=4"`
}

type StockHeader struct {
	_          struct{}  `sinli:"order=1,length=1,fixed=C"`
	ClientName string    `sinli:"order=2,length=40"`
	StockDate  time.Time `sinli:"order=3,length=8"`
	StockCoin  Coin      `sinli:"order=4,length=3"`
}

// Coin as in ISO 4217
type Coin string

const (
	CoinEuro Coin = "EUR"
)

type StockDetail struct {
	_               struct{} `sinli:"order=1,length=1,fixed=D"`
	ISBN            string   `sinli:"order=2,length=17"`
	Quantity        int      `sinli:"order=3,length=6"`
	PriceWithoutVAT float32  `sinli:"order=4,length=10"`
}

// Sale is a sinli sale. Code: `CEGALV`
type Sale struct {
	IdentificationHeader IdentificationHeader `sinli:"order=1"`
	Identification       Identification       `sinli:"order=2"`
	Header               SaleHeader           `sinli:"order=3"`
	Tickets              []SaleTicket         `sinli:"order=4"`
}

type SaleHeader struct {
	_            struct{}  `sinli:"order=1,length=1,fixed=C"`
	ClientName   string    `sinli:"order=2,length=40"`
	DispatchDate time.Time `sinli:"order=3,length=8"`
	Coin         Coin      `sinli:"order=4,length=3"`
}

type SaleTicket struct {
	_ struct{} `sinli:"order=1,length=1,fixed=T"`
	// Leave unset for generic client
	ClientNumber int          `sinli:"order=2,length=10"`
	SaleDate     time.Time    `sinli:"order=3,length=8"`
	SaleNumber   string       `sinli:"order=4,length=10"`
	NetAmount    float32      `sinli:"order=5,length=10"`
	Details      []SaleDetail `sinli:"order=6"`
}

type SaleDetail struct {
	_               struct{} `sinli:"order=1,length=1,fixed=D"`
	ISBN            string   `sinli:"order=2,length=17"`
	Quantity        int      `sinli:"order=3,length=6"`
	PriceWithoutVAT float32  `sinli:"order=4,length=10"`
}
