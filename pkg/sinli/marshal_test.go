package sinli

import (
	"bytes"
	"testing"
	"time"
)

// Example usage
type Bar struct {
	_       struct{} `sinli:"order=1,length=1,fixed=B"`
	Text    string   `sinli:"order=2,length=10"`
	Number  int      `sinli:"order=3,length=5"`
	Decimal float32  `sinli:"order=4,length=5"`
	Boolean bool     `sinli:"order=5,length=1"`
}

type Foo struct {
	Header Bar   `sinli:"order=1"`
	Body   []Bar `sinli:"order=2"`
}

func TestSinliMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{
			name: "struct",
			input: Bar{
				Text:    "Hello",
				Number:  42,
				Decimal: 12.34,
				Boolean: true,
			},
			want: []byte("BHello     0004201234S\r\n"),
		},
		{
			name: "pointer",
			input: &Bar{
				Text:    "Hello",
				Number:  -42,
				Decimal: -12.34,
				Boolean: true,
			},
			want: []byte("BHello     -0042-1234S\r\n"),
		},
		{
			name: "slice",
			input: []Bar{
				{
					Text:    "Hello",
					Number:  42,
					Decimal: 12.34,
					Boolean: true,
				},
				{
					Text:    "World",
					Number:  -43,
					Decimal: -12.35,
					Boolean: false,
				},
			},
			want: []byte("BHello     0004201234S\r\nBWorld     -0043-1235N\r\n"),
		},
		{
			name: "array",
			input: [2]Bar{
				{
					Text:    "Hello",
					Number:  42,
					Decimal: 12.34,
					Boolean: true,
				},
				{
					Text:    "World",
					Number:  -43,
					Decimal: -12.35,
					Boolean: false,
				},
			},
			want: []byte("BHello     0004201234S\r\nBWorld     -0043-1235N\r\n"),
		},
		{
			name: "nested",
			input: Foo{
				Header: Bar{
					Text:    "Header",
					Number:  1,
					Decimal: 0.01,
					Boolean: true,
				},
				Body: []Bar{
					{
						Text:    "Hello",
						Number:  42,
						Decimal: 12.34,
						Boolean: true,
					},
					{
						Text:    "World",
						Number:  -43,
						Decimal: -12.35,
						Boolean: false,
					},
				},
			},
			want: []byte("BHeader    0000100001S\r\nBHello     0004201234S\r\nBWorld     -0043-1235N\r\n"),
		},
		{
			name: "stock",
			input: Stock{
				IdentificationHeader: IdentificationHeader{
					Format:        FormatTypeNormalized,
					Document:      FileTypeStock,
					Version:       FileVersionStock,
					SourceID:      "12345678",
					DestinationID: "12345678",
				},
				Identification: Identification{
					SourceEmail:      "source@fakemail.com",
					DestinationEmail: "destination@fakemail.com",
					FileType:         FileTypeStock,
					FileVersion:      FileVersionStock,
				},
				Header: StockHeader{
					ClientName: "Client Name",
					StockDate:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					StockCoin:  CoinEuro,
				},
				Details: []StockDetail{
					{
						ISBN:            "9781234567890",
						Quantity:        1,
						PriceWithoutVAT: 1.01,
					},
					{
						ISBN:            "9781234567891",
						Quantity:        1,
						PriceWithoutVAT: 0.99,
					},
				},
			},
			want: []byte("INCEGALD021234567812345678000000000000                                     FANDE\r\n" +
				"Isource@fakemail.com                               destination@fakemail.com                          CEGALD0200000000\r\n" +
				"CClient Name                             20210101EUR\r\n" +
				"D9781234567890    0000010000000101\r\n" +
				"D9781234567891    0000010000000099\r\n",
			),
		},
		{
			name: "sale",
			input: Sale{
				IdentificationHeader: IdentificationHeader{
					Format:        FormatTypeNormalized,
					Document:      FileTypeSale,
					Version:       FileVersionSale,
					SourceID:      "12345678",
					DestinationID: "12345678",
				},
				Identification: Identification{
					SourceEmail:      "source@fakemail.com",
					DestinationEmail: "destination@fakemail.com",
					FileType:         FileTypeSale,
					FileVersion:      FileVersionSale,
				},
				Header: SaleHeader{
					ClientName:   "Client Name",
					DispatchDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Coin:         CoinEuro,
				},
				Tickets: []SaleTicket{
					{
						ClientNumber: 777,
						SaleDate:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						SaleNumber:   "1234",
						NetAmount:    2.0,
						Details: []SaleDetail{
							{
								ISBN:            "9781234567890",
								Quantity:        1,
								PriceWithoutVAT: 1.01,
							},
							{
								ISBN:            "9781234567891",
								Quantity:        1,
								PriceWithoutVAT: 0.99,
							},
						},
					},
					{
						ClientNumber: 777,
						SaleDate:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						SaleNumber:   "1234",
						NetAmount:    2.0,
						Details: []SaleDetail{
							{
								ISBN:            "9781234567890",
								Quantity:        1,
								PriceWithoutVAT: 1.01,
							},
							{
								ISBN:            "9781234567891",
								Quantity:        1,
								PriceWithoutVAT: 0.99,
							},
						},
					},
				},
			},
			want: []byte("INCEGALV031234567812345678000000000000                                     FANDE\r\n" +
				"Isource@fakemail.com                               destination@fakemail.com                          CEGALV0300000000\r\n" +
				"CClient Name                             20210101EUR\r\n" +
				"T0000000777202101011234      0000000200\r\n" +
				"D9781234567890    0000010000000101\r\n" +
				"D9781234567891    0000010000000099\r\n" +
				"T0000000777202101011234      0000000200\r\n" +
				"D9781234567890    0000010000000101\r\n" +
				"D9781234567891    0000010000000099\r\n",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			got, err := Marshal(tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil, got error: %s", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Fatalf("want:\n'%s'\ngot:\n'%s'\n", tt.want, got)
			}
		})
	}
}
