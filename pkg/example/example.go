package example

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/igolaizola/agorer"
	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/isbn"
	"github.com/igolaizola/agorer/pkg/sinli"
)

type Config struct {
	SourceEmail      string
	DestinationEmail string
	SourceID         string
	DestinationID    string
	ProviderName     string
	ClientName       string

	MasterFile string
	DayFile    string
	OutputDir  string
}

func Run(ctx context.Context, c *Config) error {
	b, err := os.ReadFile(c.MasterFile)
	if err != nil {
		return fmt.Errorf("couldn't read %s: %w", c.MasterFile, err)
	}
	var master agora.Master
	if err := json.Unmarshal(b, &master); err != nil {
		return fmt.Errorf("couldn't unmarshal master: %w", err)
	}

	b, err = os.ReadFile(c.DayFile)
	if err != nil {
		return fmt.Errorf("couldn't read %s: %w", c.DayFile, err)
	}
	var day agora.Day
	if err := json.Unmarshal(b, &day); err != nil {
		return fmt.Errorf("couldn't unmarshal %s: %w", c.DayFile, err)
	}

	isbnClient, err := isbn.New("data/isbn.json", "data/isbn_err.json")
	if err != nil {
		return fmt.Errorf("couldn't create isbn client: %w", err)
	}
	s := agorer.NewStore(ctx, &master, isbnClient)

	stockDate := time.Now()
	stockItems, _, err := agorer.StockItems(ctx, s)
	if err != nil {
		return fmt.Errorf("couldn't generate stock items: %w", err)
	}
	stockDetails, err := agorer.StockDetails(ctx, stockItems)
	if err != nil {
		return fmt.Errorf("couldn't generate stock: %w", err)
	}
	stock := sinli.Stock{
		IdentificationHeader: sinli.IdentificationHeader{
			Format:        sinli.FormatTypeNormalized,
			Document:      sinli.FileTypeStock,
			Version:       sinli.FileVersionStock,
			SourceID:      c.SourceID,
			DestinationID: c.DestinationID,
		},
		Identification: sinli.Identification{
			SourceEmail:      c.SourceEmail,
			DestinationEmail: c.DestinationEmail,
			FileType:         sinli.FileTypeStock,
			FileVersion:      sinli.FileVersionStock,
		},
		Header: sinli.StockHeader{
			ClientName: c.ProviderName,
			StockDate:  stockDate,
			StockCoin:  sinli.CoinEuro,
		},
		Details: stockDetails,
	}
	subject := sinli.Subject{
		SourceID:      c.SourceID,
		DestinationID: c.DestinationID,
		FileType:      sinli.FileTypeStock,
		FileVersion:   sinli.FileVersionStock,
	}

	// Marshal sinli stock
	b, err = sinli.Marshal(stock)
	if err != nil {
		return fmt.Errorf("couldn't marshal stock: %w", err)
	}
	file := filepath.Join(c.OutputDir, fmt.Sprintf("%s.txt", sinli.FileTypeStock))
	if err := os.WriteFile(file, b, 0644); err != nil {
		return fmt.Errorf("couldn't write file %s: %w", file, err)
	}
	// Marshal sinli subject
	b, err = sinli.Marshal(subject)
	if err != nil {
		return fmt.Errorf("couldn't marshal subject: %w", err)
	}
	file = fmt.Sprintf("%s.subject", file)
	if err := os.WriteFile(file, b, 0644); err != nil {
		return fmt.Errorf("couldn't write file %s: %w", file, err)
	}

	for _, inv := range day.Invoices {
		details, err := agorer.OrderDetails(ctx, s, &inv)
		if err != nil {
			return fmt.Errorf("couldn't generate orders: %w", err)
		}
		if len(details) == 0 {
			continue
		}
		date, err := time.Parse("2006-01-02T15:04:05", inv.Date)
		if err != nil {
			return fmt.Errorf("couldn't parse date: %w", err)
		}
		order := sinli.Order{
			IdentificationHeader: sinli.IdentificationHeader{
				Format:        sinli.FormatTypeNormalized,
				Document:      sinli.FileTypeOrder,
				Version:       sinli.FileVersionOrder,
				SourceID:      c.SourceID,
				DestinationID: c.DestinationID,
			},
			Identification: sinli.Identification{
				SourceEmail:      c.SourceEmail,
				DestinationEmail: c.DestinationEmail,
				FileType:         sinli.FileTypeOrder,
				FileVersion:      sinli.FileVersionOrder,
			},
			Header: sinli.OrderHeader{
				ClientName:   c.ClientName,
				ProviderName: c.ProviderName,
				OrderDate:    date,
				OrderCode:    strconv.Itoa(inv.Number),
				OrderType:    sinli.OrderTypeNormal,
				Coin:         sinli.LegacyCoinEuro,
				Batch:        strconv.Itoa(inv.Number),
			},
			Details: details,
		}
		subject := sinli.Subject{
			SourceID:      c.SourceID,
			DestinationID: c.DestinationID,
			FileType:      sinli.FileTypeOrder,
			FileVersion:   sinli.FileVersionOrder,
		}

		// Marshal sinli order
		b, err = sinli.Marshal(order)
		if err != nil {
			return fmt.Errorf("couldn't marshal order: %w", err)
		}
		file := filepath.Join(c.OutputDir, fmt.Sprintf("%s_%d.txt", sinli.FileTypeOrder, inv.Number))
		if err := os.WriteFile(file, b, 0644); err != nil {
			return fmt.Errorf("couldn't write file %s: %w", file, err)
		}
		// Marshal sinli subject
		b, err = sinli.Marshal(subject)
		if err != nil {
			return fmt.Errorf("couldn't marshal subject: %w", err)
		}
		file = fmt.Sprintf("%s.subject", file)
		if err := os.WriteFile(file, b, 0644); err != nil {
			return fmt.Errorf("couldn't write file %s: %w", file, err)
		}
	}

	for _, inv := range day.Invoices {
		details, err := agorer.ReturnDetails(ctx, s, &inv)
		if err != nil {
			return fmt.Errorf("couldn't generate returns: %w", err)
		}
		if len(details) == 0 {
			continue
		}
		date, err := time.Parse("2006-01-02T15:04:05", inv.Date)
		if err != nil {
			return fmt.Errorf("couldn't parse date: %w", err)
		}
		ret := sinli.Return{
			IdentificationHeader: sinli.IdentificationHeader{
				Format:        sinli.FormatTypeNormalized,
				Document:      sinli.FileTypeReturn,
				Version:       sinli.FileVersionReturn,
				SourceID:      c.SourceID,
				DestinationID: c.DestinationID,
			},
			Identification: sinli.Identification{
				SourceEmail:      c.SourceEmail,
				DestinationEmail: c.DestinationEmail,
				FileType:         sinli.FileTypeReturn,
				FileVersion:      sinli.FileVersionReturn,
			},
			Header: sinli.ReturnHeader{
				ClientName:   c.ClientName,
				ProviderName: c.ProviderName,
				OrderCode:    strconv.Itoa(inv.Number),
				DocumentDate: date,
				DocumentType: sinli.ReturnDocumentTypeDefinitive,
				ReturnType:   sinli.ReturnTypeDefinitive,
				Coin:         sinli.LegacyCoinEuro,
			},
			Details: details,
		}
		subject := sinli.Subject{
			SourceID:      c.SourceID,
			DestinationID: c.DestinationID,
			FileType:      sinli.FileTypeReturn,
			FileVersion:   sinli.FileVersionReturn,
		}

		// Marshal sinli return
		b, err = sinli.Marshal(ret)
		if err != nil {
			return fmt.Errorf("couldn't marshal return: %w", err)
		}
		file := filepath.Join(c.OutputDir, fmt.Sprintf("%s_%d.txt", sinli.FileTypeReturn, inv.Number))
		if err := os.WriteFile(file, b, 0644); err != nil {
			return fmt.Errorf("couldn't write file %s: %w", file, err)
		}
		// Marshal sinli subject
		b, err = sinli.Marshal(subject)
		if err != nil {
			return fmt.Errorf("couldn't marshal subject: %w", err)
		}
		file = fmt.Sprintf("%s.subject", file)
		if err := os.WriteFile(file, b, 0644); err != nil {
			return fmt.Errorf("couldn't write file %s: %w", file, err)
		}
	}

	return nil
}
