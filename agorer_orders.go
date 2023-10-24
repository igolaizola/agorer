package agorer

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/sinli"
)

func OrderDetails(ctx context.Context, s *Store, inv *agora.Invoice) ([]sinli.OrderDetail, error) {
	var details []sinli.OrderDetail

	for _, item := range inv.InvoiceItems {
		for _, l := range item.Lines {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			quantity := int(l.Quantity)
			if quantity <= 0 {
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
			price := l.ProductPrice
			if !priceList.VatIncluded {
				vat := s.Vats[p.VatID]
				rate := vat.VatRate
				// Add VAT to price
				price = price * (1 + rate)
			}
			details = append(details, sinli.OrderDetail{
				ISBN:         isbnCode,
				EAN:          fmt.Sprintf("%s00000", barCode),
				Reference:    strconv.Itoa(l.ProductID),
				Title:        p.Name,
				Quantity:     int(l.Quantity),
				PriceWithVAT: price,
				OrderSource:  sinli.OrderSourceClient,
				Code:         strconv.Itoa(inv.Number),
			})
		}
	}
	return details, nil
}
