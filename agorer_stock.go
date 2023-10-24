package agorer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/isbn"
	"github.com/igolaizola/agorer/pkg/mail"
	"github.com/igolaizola/agorer/pkg/sinli"
)

func Stock(ctx context.Context, c *Config) error {
	// Validate config
	if c.Input == "" {
		return errors.New("input must be provided")
	}
	if c.LogDir == "" {
		return errors.New("log dir must be provided")
	}
	// Validate input type
	var agoraHost string
	switch c.InputType {
	case "agora":
		if c.AgoraToken == "" {
			return errors.New("agora token must be provided")
		}
		agoraHost = c.Input
	case "agora-json":
		port, err := agora.MockServe(ctx, ":0", c.Input)
		if err != nil {
			return fmt.Errorf("couldn't mock serve agora: %w", err)
		}
		agoraHost = fmt.Sprintf("http://localhost:%d", port)
	case "json":
	default:
		return fmt.Errorf("invalid input type %s", c.InputType)
	}
	if agoraHost != "" {
		if c.ISBNDir == "" {
			return errors.New("isbn dir must be provided")
		}
	}

	// Validate output type
	output := c.Output
	switch c.OutputType {
	case "json":
		if output == "" {
			output = filepath.Join(c.LogDir, fmt.Sprintf("stock_%s.json", time.Now().Format("20060102_150405")))
		}
	case "sinli":
		if output == "" {
			output = filepath.Join(c.LogDir, fmt.Sprintf("sinli_N_%s_%s.snl", time.Now().Format("20060102_150405"), c.SINLISourceID))
		}
		if c.SINLISourceEmail == "" {
			return errors.New("sinli source email must be provided")
		}
		if c.SINLISourceID == "" {
			return errors.New("sinli source id must be provided")
		}
		if c.SINLIDestinationEmail == "" {
			return errors.New("sinli destination email must be provided")
		}
		if c.SINLIDestinationID == "" {
			return errors.New("sinli destination id must be provided")
		}
		if c.SINLIClientName == "" {
			return errors.New("sinli client name must be provided")
		}
		if !c.Mail.Dry {
			if c.Mail.Host == "" {
				return errors.New("mail host must be provided")
			}
			if c.Mail.Port == 0 {
				return errors.New("mail port must be provided")
			}
			if c.Mail.Username == "" {
				return errors.New("mail username must be provided")
			}
			if c.Mail.Password == "" {
				return errors.New("mail password must be provided")
			}
		}
	default:
		return fmt.Errorf("invalid output type %s", c.OutputType)
	}

	// Create log dir if it doesn't exist
	if c.LogDir != "" {
		if err := os.MkdirAll(c.LogDir, 0755); err != nil {
			return fmt.Errorf("couldn't create log dir %s: %w", c.LogDir, err)
		}
	}

	var stockItems []StockItem
	if agoraHost != "" {
		// Export master data from Agora
		client := agora.New(agoraHost, c.AgoraToken, c.LogDir)
		master, err := client.ExportMaster(ctx)
		if err != nil {
			return fmt.Errorf("couldn't get master: %w", err)
		}

		// Create isbn client
		isbnClient, err := isbn.New(filepath.Join(c.ISBNDir, "isbn.json"), filepath.Join(c.ISBNDir, "isbn_err.json"))
		if err != nil {
			return fmt.Errorf("couldn't create isbn client: %w", err)
		}

		// Create store using master data and isbn client
		s := NewStore(ctx, master, isbnClient)

		var conflicts []Conflict
		stockItems, conflicts, err = StockItems(ctx, s)
		if err != nil {
			return fmt.Errorf("couldn't generate stock: %w", err)
		}

		// Write conflicts to file
		if len(conflicts) > 0 {
			b, err := json.MarshalIndent(conflicts, "", "  ")
			if err != nil {
				return fmt.Errorf("couldn't marshal conflicts: %w", err)
			}
			f := filepath.Join(c.LogDir, fmt.Sprintf("conflicts_%s.json", time.Now().Format("20060102_150405")))
			if err := os.WriteFile(f, b, 0644); err != nil {
				return fmt.Errorf("couldn't write file %s: %w", f, err)
			}
		}
	} else {
		// Read stock from json file
		b, err := os.ReadFile(c.Input)
		if err != nil {
			return fmt.Errorf("couldn't read file %s: %w", c.Input, err)
		}
		if err := json.Unmarshal(b, &stockItems); err != nil {
			return fmt.Errorf("couldn't unmarshal json: %w", err)
		}
	}

	if c.OutputType == "json" {
		// Write stock to json file
		b, err := json.MarshalIndent(stockItems, "", "  ")
		if err != nil {
			return fmt.Errorf("couldn't marshal json: %w", err)
		}
		if err := os.WriteFile(output, b, 0644); err != nil {
			return fmt.Errorf("couldn't write file %s: %w", output, err)
		}
		return nil
	}

	// Generate sinli stock
	stockDetails, err := StockDetails(ctx, stockItems)
	if err != nil {
		return fmt.Errorf("couldn't generate stock: %w", err)
	}

	// Create sinli stock
	stock := sinli.Stock{
		IdentificationHeader: sinli.IdentificationHeader{
			Format:        sinli.FormatTypeNormalized,
			Document:      sinli.FileTypeStock,
			Version:       sinli.FileVersionStock,
			SourceID:      c.SINLISourceID,
			DestinationID: c.SINLIDestinationID,
		},
		Identification: sinli.Identification{
			SourceEmail:      c.SINLISourceEmail,
			DestinationEmail: c.SINLIDestinationEmail,
			FileType:         sinli.FileTypeStock,
			FileVersion:      sinli.FileVersionStock,
		},
		Header: sinli.StockHeader{
			ClientName: c.SINLIClientName,
			StockDate:  time.Now(),
			StockCoin:  sinli.CoinEuro,
		},
		Details: stockDetails,
	}

	// Write sinli stock to output
	b, err := sinli.Marshal(stock)
	if err != nil {
		return fmt.Errorf("couldn't marshal sinli stock: %w", err)
	}
	if err := os.WriteFile(output, b, 0644); err != nil {
		return fmt.Errorf("couldn't write file %s: %w", output, err)
	}

	// Marshal subject
	sinliSubject := sinli.Subject{
		SourceID:      c.SINLISourceID,
		DestinationID: c.SINLIDestinationID,
		FileType:      sinli.FileTypeStock,
		FileVersion:   sinli.FileVersionStock,
	}
	b, err = sinli.Marshal(sinliSubject)
	if err != nil {
		return fmt.Errorf("couldn't marshal sinli subject: %w", err)
	}
	subject := strings.TrimSpace(string(b))

	// Send email
	if err := mail.Send(ctx, &c.Mail, c.SINLISourceEmail, c.SINLIDestinationEmail, subject, "", output); err != nil {
		return fmt.Errorf("couldn't send email: %w", err)
	}
	return nil
}

func StockDetails(ctx context.Context, items []StockItem) ([]sinli.StockDetail, error) {
	var details []sinli.StockDetail
	for _, item := range items {
		details = append(details, sinli.StockDetail{
			ISBN:            item.ISBN,
			Quantity:        item.Quantity,
			PriceWithoutVAT: item.PriceWithoutVAT,
		})
	}
	return details, nil
}

type StockItem struct {
	Name            string  `json:"name"`
	ISBN            string  `json:"isbn"`
	Quantity        int     `json:"quantity"`
	PriceWithVAT    float32 `json:"price_with_vat"`
	PriceWithoutVAT float32 `json:"price_without_vat"`
}

type Conflict struct {
	ISBN  string
	Names []string
}

func StockItems(ctx context.Context, s *Store) ([]StockItem, []Conflict, error) {
	var items []StockItem
	for id, qty := range s.Quantity {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}
		p := s.Books[id]
		if qty == 0 {
			continue
		}
		if qty < 0 {
			continue
		}
		if len(p.Prices) == 0 {
			log.Println("❌ no price for", p.ID, p.Name)
			continue
		}
		if len(p.Prices) > 1 {
			log.Println("❌ more than one price for", p.ID, p.Name)
			continue
		}
		priceData := p.Prices[0]
		if priceData.Price == 0 {
			log.Println("❌ price is 0 for", p.ID, p.Name)
			continue
		}
		priceList, ok := s.PriceLists[priceData.PriceListID]
		if !ok {
			log.Println("❌ price list not found for", p.ID, p.Name)
			continue
		}
		vat, ok := s.Vats[p.VatID]
		if !ok {
			log.Println("❌ vat not found for", p.ID, p.Name)
			continue
		}
		isbnCode := s.ISBNs[p.ID]
		priceWithVAT := priceData.Price
		priceWithoutVAT := priceData.Price
		if priceList.VatIncluded {
			rate := vat.VatRate
			// Remove VAT from price
			priceWithoutVAT = priceData.Price / (1 + rate)
		} else {
			rate := vat.VatRate
			// Add VAT to price
			priceWithVAT = priceData.Price * (1 + rate)
		}

		items = append(items, StockItem{
			Name:            p.Name,
			ISBN:            isbnCode,
			Quantity:        qty,
			PriceWithoutVAT: priceWithoutVAT,
			PriceWithVAT:    priceWithVAT,
		})
	}

	// Group items by ISBN
	groups := make(map[string][]StockItem)
	for _, item := range items {
		groups[item.ISBN] = append(groups[item.ISBN], item)
	}

	// Find conflicts
	var conflicts []Conflict
	var filtered []StockItem
	for key, group := range groups {
		if len(group) == 1 {
			filtered = append(filtered, group[0])
			continue
		}
		var names []string
		for _, item := range group {
			names = append(names, item.Name)
		}
		sort.Strings(names)
		conflicts = append(conflicts, Conflict{
			ISBN:  key,
			Names: names,
		})
	}

	// Sort items by ISBN to be deterministic
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ISBN < filtered[j].ISBN
	})
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].ISBN < conflicts[j].ISBN
	})

	return filtered, conflicts, nil
}
