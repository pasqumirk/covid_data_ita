# covid_data_ita
Questo tool recupera i dati ufficiali della protezione civile per quanto riguarda il COVID 19 e plotta su un grafico l'andamento di malati, ricoveri, terapie intensive e decessi.

Rifiutando qualisasi etichettatura (come novax e provax) il tool deve essere soltanto di supporto per un monitoraggio o controllo.

I dati vengono scaricati automaticamente dal repository della protezione civile (https://github.com/pcm-dpc/COVID-19/) a garanzia di trasparenza.

Il tool è scritto in go (https://golang.org/). Può essere eseguito con "go run main.go" oppure compilato con "go build": è stato testato su windows ma dovrebbe funzionare correttamente anche su linux e mac

TODO: Integrare i dati con quelli di https://github.com/italia/covid19-opendata-vaccini
