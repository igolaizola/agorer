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
You can use `stock.conf.example` as a template.
Use the help command to see the available options:

```bash
agorer stock --help
```

Then run this command to obtain the stock data from Agora Retail and send it by email in SINLI format:

```bash
agorer stock --config stock.conf
```

## 📚 Resources

 - [SINLI documentation](http://www.fande.es/sinli_indicedocumentos.html)
 - [Agora Retail documentation](https://www.agorapos.com/manual/agora-retail/guia-integracion-agora-retail.pdf)
