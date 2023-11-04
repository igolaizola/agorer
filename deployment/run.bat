:: Generate stock and sales data
agorer.exe stock --config agora.conf --output-type json --output data/stock.json --log-dir logs --isbn-dir data
agorer.exe sales --config agora.conf --output-type json --output data/sales --log-dir logs isnb-dir data
echo %date% %time% > data/date.txt
