package agorer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/isbn"
	"github.com/igolaizola/agorer/pkg/mail"
	"github.com/igolaizola/agorer/pkg/sinli"
)

func Sales(ctx context.Context, c *Config, day time.Time) error {
	// Validate config
	input := c.Input
	if input == "" {
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
		// Check if input is a directory
		if fi, err := os.Stat(c.Input); err == nil && fi.IsDir() {
			input = filepath.Join(input, fmt.Sprintf("%s.json", day.Format("2006-01-02")))
		}
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
			output = filepath.Join(c.LogDir, fmt.Sprintf("sale_%s.json", time.Now().Format("20060102_150405")))
		}
		// Check if output is a directory
		if fi, err := os.Stat(output); err == nil && fi.IsDir() {
			output = filepath.Join(output, fmt.Sprintf("%s.json", day.Format("2006-01-02")))
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

	tickets := []SaleTicket{}
	if agoraHost != "" {
		// Export day data data from Agora
		client := agora.New(agoraHost, c.AgoraToken, c.LogDir)
		d, err := client.ExportDay(ctx, day)
		if err != nil {
			return fmt.Errorf("couldn't get day: %w", err)
		}

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

		for _, inv := range d.Invoices {
			date, err := time.Parse("2006-01-02T15:04:05", inv.Date)
			if err != nil {
				return fmt.Errorf("couldn't parse date %s: %w", inv.Date, err)
			}
			var netAmount float32
			ticket := SaleTicket{
				SaleDate:   date,
				SaleNumber: strconv.Itoa(inv.Number),
			}
			for _, item := range inv.InvoiceItems {
				for _, line := range item.Lines {
					isbnCode, ok := s.ISBNs[line.ProductID]
					if !ok {
						continue
					}
					priceWithVAT := float32(line.ProductPrice)
					priceWithoutVAT := priceWithVAT / (1 + line.VatRate)
					item := SaleItem{
						Name:            line.ProductName,
						ISBN:            isbnCode,
						Quantity:        int(line.Quantity),
						PriceWithoutVAT: priceWithoutVAT,
					}
					ticket.Items = append(ticket.Items, item)
					netAmount += priceWithoutVAT * float32(line.Quantity)
				}
			}
			if len(ticket.Items) == 0 {
				continue
			}
			ticket.NetAmount = netAmount
			tickets = append(tickets, ticket)
		}
	} else {
		// Read stock from json file
		b, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("couldn't read file %s: %w", input, err)
		}
		if err := json.Unmarshal(b, &tickets); err != nil {
			return fmt.Errorf("couldn't unmarshal json: %w", err)
		}
	}

	if c.OutputType == "json" {
		// Write stock to json file
		b, err := json.MarshalIndent(tickets, "", "  ")
		if err != nil {
			return fmt.Errorf("couldn't marshal json: %w", err)
		}
		if err := os.WriteFile(output, b, 0644); err != nil {
			return fmt.Errorf("couldn't write file %s: %w", output, err)
		}
		return nil
	}

	// Generate sinli stock
	sinliTickets, err := SaleTickets(ctx, tickets)
	if err != nil {
		return fmt.Errorf("couldn't generate stock: %w", err)
	}

	// Create sinli stock
	stock := sinli.Sale{
		IdentificationHeader: sinli.IdentificationHeader{
			Format:        sinli.FormatTypeNormalized,
			Document:      sinli.FileTypeSale,
			Version:       sinli.FileVersionSale,
			SourceID:      c.SINLISourceID,
			DestinationID: c.SINLIDestinationID,
		},
		Identification: sinli.Identification{
			SourceEmail:      c.SINLISourceEmail,
			DestinationEmail: c.SINLIDestinationEmail,
			FileType:         sinli.FileTypeStock,
			FileVersion:      sinli.FileVersionStock,
		},
		Header: sinli.SaleHeader{
			ClientName:   c.SINLIClientName,
			DispatchDate: day,
			Coin:         sinli.CoinEuro,
		},
		Tickets: sinliTickets,
	}

	// Write sinli stock to output
	b, err := sinli.Marshal(stock)
	if err != nil {
		return fmt.Errorf("couldn't marshal sinli sale: %w", err)
	}
	if err := os.WriteFile(output, b, 0644); err != nil {
		return fmt.Errorf("couldn't write file %s: %w", output, err)
	}

	// Marshal subject
	sinliSubject := sinli.Subject{
		SourceID:      c.SINLISourceID,
		DestinationID: c.SINLIDestinationID,
		FileType:      sinli.FileTypeSale,
		FileVersion:   sinli.FileVersionSale,
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

type SaleTicket struct {
	SaleDate   time.Time  `json:"sale_date"`
	SaleNumber string     `json:"sale_number"`
	NetAmount  float32    `json:"net_amount"`
	Items      []SaleItem `json:"items"`
}

type SaleItem struct {
	Name            string  `json:"name"`
	ISBN            string  `json:"isbn"`
	Quantity        int     `json:"quantity"`
	PriceWithoutVAT float32 `json:"price_without_vat"`
}

func SaleTickets(ctx context.Context, ts []SaleTicket) ([]sinli.SaleTicket, error) {
	var tickets []sinli.SaleTicket
	for _, t := range ts {
		sinliTicket := sinli.SaleTicket{
			SaleDate:   t.SaleDate,
			SaleNumber: t.SaleNumber,
			NetAmount:  t.NetAmount,
		}
		for _, item := range t.Items {
			sinliTicket.Details = append(sinliTicket.Details, sinli.SaleDetail{
				ISBN:            item.ISBN,
				Quantity:        item.Quantity,
				PriceWithoutVAT: item.PriceWithoutVAT,
			})
		}
		tickets = append(tickets, sinliTicket)
	}
	return tickets, nil
}
