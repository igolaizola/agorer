:: Generate stock and sales data, commit to git and shutdown
git pull
git reset --hard origin/main
git pull

agorer.exe stock --config agora.conf --output-type json --output data/stock.json --log-dir logs --isbn-dir data
agorer.exe sales --config agora.conf --output-type json --output data/sales --log-dir logs isnb-dir data
echo %date% %time% > data/date.txt

git add data/*
git commit -m "Update data"
git push origin main

shutdown.exe /s /t 30
echo "Press enter to abort shutdown."
pause
shutdown.exe /a
pause
