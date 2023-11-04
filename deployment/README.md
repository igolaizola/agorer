# agorer deployment template

This template creates a deployment for agorer.

Use the provided scripts to generate the latest stock and sales data and commit them to the repository from the computer where Agora is installed.

GitHub Actions will send nightly emails with the updated stock and sales data.

## Set up

Create a new repository using this template.

Place the latest `agorer.exe` binary in the root of the repository.

Create `agora.conf` and `mail.conf` files using provided examples.

Add the following github action secrets:

```
MAIL_PASS <mail-password>
```

Add the following github action variables:

```
MAIL_HOST <mail-host>
MAIL_PORT <mail-port-587>
MAIL_USER <mail-user-or-email>
SINLI_CLIENT_NAME <YOUR-AWESOME-LIBRARY-NAME>
SINLI_DESTINATION_EMAIL sinli@cegal.es
SINLI_DESTINATION_ID LIB00022
SINLI_SOURCE_EMAIL <your-sinli-email>
SINLI_SOURCE_ID <L000XXXX>
```

## Scripts

`run-sync-exit.sh`

This script generates the latest stock and sales data and commits it to the repository.
Afterward, it shuts down the computer.

`run-sync-mail.bat`

This script generates the latest stock and sales data and commits it to the repository.
Then, it sends an email with the latest stock data.
Use this script if you need to ensure your latest stock data is synchronized.

`run-sync.bat`

This script generates the latest stock and sales data and commits it to the repository.

`run-sync.sh`

This script generates the latest stock and sales data but does not commit it to the repository.
