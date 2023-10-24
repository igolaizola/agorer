package agorer

import (
	"context"
	"errors"
	"log"

	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/isbn"
	"github.com/igolaizola/agorer/pkg/mail"
)

type Config struct {
	Debug  bool
	LogDir string

	Input      string
	InputType  string
	Output     string
	OutputType string

	AgoraToken string

	ISBNDir string

	SINLISourceEmail      string
	SINLISourceID         string
	SINLIDestinationEmail string
	SINLIDestinationID    string
	SINLIClientName       string

	Mail mail.Config
}

type Store struct {
	Products   map[int]agora.Product
	Books      map[int]agora.Product
	Vats       map[int]agora.Vat
	PriceLists map[int]agora.PriceList
	Quantity   map[int]int
	ISBNs      map[int]string
}

func NewStore(ctx context.Context, master *agora.Master, isbnCli *isbn.Client) *Store {
	vats := map[int]agora.Vat{}
	for _, vat := range master.Vats {
		vats[vat.ID] = vat
	}

	books := map[int]agora.Product{}
	isbns := map[int]string{}
	for _, pr := range master.Products {
		if pr.DeletionDate != "" {
			continue
		}
		barcode := pr.Barcode()
		if !isbn.Valid(barcode) {
			continue
		}
		vat, ok := vats[pr.VatID]
		if !ok {
			continue
		}

		// Books have VatRate of 0.04 or less
		if vat.VatRate > 0.04 {
			continue
		}

		isbnCode, err := isbnCli.Hyphenate(ctx, barcode, pr.Name)
		if err != nil {
			if !errors.Is(err, isbn.ErrNotFound) {
				log.Println("❌ couldn't get isbn for", pr.Name, barcode, err)
			}
			continue
		}
		isbns[pr.ID] = isbnCode
		books[pr.ID] = pr
	}

	priceLists := map[int]agora.PriceList{}
	for _, pl := range master.PriceLists {
		priceLists[pl.ID] = pl
	}

	quantity := map[int]int{}
	for _, st := range master.Stocks {
		if _, ok := books[st.ProductID]; !ok {
			continue
		}
		quantity[st.ProductID] += int(st.Quantity)
	}
	return &Store{
		Books:      books,
		Vats:       vats,
		PriceLists: priceLists,
		Quantity:   quantity,
		ISBNs:      isbns,
	}
}
