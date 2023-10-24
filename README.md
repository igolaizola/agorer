# agorer

**agorer** is a tool to sync [Agora Retail software](https://www.agorapos.com/) with [Todostuslibros.com](https://www.todostuslibros.com/) using [SINLI](http://www.fande.es/sinli_indicedocumentos.html) format.

## 📦 Installation

You can use the Golang binary to install **agorer**:

```bash
go install github.com/igolaizola/agorer/cmd/agorer@latest
```

Or you can download the binary from the [releases](https://github.com/igolaizola/agorer/releases)

## 🕹️ Usage

Create a configuration file.
You can use `example.conf` as a template.

Use the help command to see the available options:

```bash
agorer --help
agorer stock --help
agorer sales --help
```

### stock

Run this command to obtain the stock data from Agora Retail and send it by email in SINLI format:

```bash
agorer stock --config stock.conf
```

### sales

Run this command to obtain sales data of a given day from Agora Retail and send it by email in SINLI format:

```bash
agorer stock --config sales.conf --day 2023-02-28
```

If `day` is not specified, the current day is used.

## 📚 Resources

 - [SINLI documentation](http://www.fande.es/sinli_indicedocumentos.html)
 - [Agora Retail documentation](https://www.agorapos.com/manual/agora-retail/guia-integracion-agora-retail.pdf)
