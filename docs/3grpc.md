# gRPC

gRPC (acronimo di gRPC Remote Procedure Calls) è un framework multipiattaforma ad alte prestazioni per chiamate di procedura remota (RPC). gRPC è stato inizialmente creato da Google, ma è open source e viene utilizzato in molte organizzazioni. I casi d'uso spaziano dai microservizi all'"ultimo miglio" dell'informatica (mobile, web e Internet of Things). gRPC utilizza HTTP/2 per il trasporto, Protocol Buffers come linguaggio di descrizione dell'interfaccia e fornisce funzionalità quali autenticazione, streaming bidirezionale e controllo del flusso, binding bloccanti o non bloccanti. Genera binding client e server multipiattaforma per molti linguaggi. Gli scenari di utilizzo più comuni includono la connessione di servizi in un'architettura in stile microservizi o la connessione di client di dispositivi mobili a servizi backend.

A partire dal 2019, l'utilizzo di HTTP/2 da parte di gRPC rende impossibile implementare un client gRPC in un browser, richiedendo invece un proxy.

gRPC supporta l'utilizzo di Transport Layer Security (TLS) e l'autenticazione basata su token. Per l'autorizzazione basata su token, gRPC fornisce Server Interceptor e un Client Interceptor.

Più che altro ne conosco i vantaggi, per esempio sul milione di richieste è molto più performante rispetto all'HTTP ed è molto comodo per lo streaming di dati. Ti puoi sicuramente allontanare dall'implementare delle API con il modello REST visto che sono delle procedure remote che chiami direttamente

poi utilizza protocol buffer che per la serializzazione/deserializzazione è molto più performante rispetto ad un JSON visto che sono byte scambiati. Lo svantaggio è che il debugging è più ostico perchè devi avere un qualcosa che ti trasforma quei dati in formato leggibili (altrimenti rimangono byte) (modificato).


# gRPC senza Istio opzione con headless service

Il motivo principale per cui è difficile bilanciare il traffico gRPC è che le persone vedono gRPC come HTTP ed è qui che inizia il problema, in base alla progettazione sono diversi, mentre HTTP crea e chiude le connessioni per richiesta, gRPC opera su un protocollo HTTP2 che funziona su una connessione TCP di lunga durata che rende più difficile il bilanciamento poiché più richieste passano attraverso la stessa connessione grazie alla funzione multiplex.

Senza la service mesh è necessario implementare soluzioni alternative per gestire efficacemente il bilanciamento del carico del traffico gRPC in Kubernetes. Queste soluzioni possono includere la configurazione diretta dei client per gestire più connessioni server o l'utilizzo di proxy come meccanismo di bilanciamento, poiché averlo nell'applicazione principale aumenterà la complessità e il consumo di risorse nel tentativo di mantenere molte connessioni aperte e ribilanciarle. Immagina questo, con il bilanciamento finale nell'app avresti 1 istanza connessa a N server (1-N), ma con un proxy avresti 1 istanza connessa a M proxy connessi a N server (1-M-N) dove sicuramente M < N poiché ogni istanza proxy può gestire molte connessioni ai diversi server.