:: Generate stock and sales data, send stock mail and commit to git
git pull
git reset --hard origin/main
git pull

agorer.exe stock --config agora.conf --output-type json --output data/stock.json --log-dir logs --isbn-dir data
agorer.exe sales --config agora.conf --output-type json --output data/sales --log-dir logs isnb-dir data
agorer.exe stock --config mail.conf --output-type sinli --input-type json --input data/stock.json
echo %date% %time% > data/date.txt

git add data/*
git commit -m "Update data and send mail"
git push origin main

echo "Press any key to close..."
pause
