package agorer

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/sinli"
)

func ReturnDetails(ctx context.Context, s *Store, inv *agora.Invoice) ([]sinli.ReturnDetail, error) {
	var details []sinli.ReturnDetail
	for _, item := range inv.InvoiceItems {
		for _, l := range item.Lines {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			quantity := int(l.Quantity)
			if quantity >= 0 {
				continue
			}
			p, ok := s.Books[l.ProductID]
			if !ok {
				log.Println("❌ product not found for", l.ProductID)
				continue
			}
			barCode := p.Barcode()
			isbnCode := s.ISBNs[p.ID]
			priceList := s.PriceLists[item.PriceList.ID]
			priceWithVAT := l.ProductPrice
			priceWithoutVAT := l.ProductPrice
			if !priceList.VatIncluded {
				vat := s.Vats[p.VatID]
				rate := vat.VatRate
				// Add VAT to price
				priceWithVAT = priceWithVAT * (1 + rate)
			} else {
				vat := s.Vats[p.VatID]
				rate := vat.VatRate
				// Remove VAT from price
				priceWithoutVAT = priceWithoutVAT / (1 + rate)
			}
			details = append(details, sinli.ReturnDetail{
				ISBN:            isbnCode,
				EAN:             fmt.Sprintf("%s00000", barCode),
				Reference:       strconv.Itoa(l.ProductID),
				Title:           p.Name,
				Quantity:        int(math.Abs(float64(l.Quantity))),
				PriceWithoutVAT: priceWithoutVAT,
				PriceWithVAT:    priceWithVAT,
				Discount:        100.0,
				PriceType:       sinli.PriceTypeFixed,
			})
		}
	}
	return details, nil
}
