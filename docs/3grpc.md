# gRPC

gRPC (acronimo di gRPC Remote Procedure Calls) è un framework multipiattaforma ad alte prestazioni per chiamate di procedura remota (RPC). gRPC è stato inizialmente creato da Google, ma è open source e viene utilizzato in molte organizzazioni. I casi d'uso spaziano dai microservizi all'"ultimo miglio" dell'informatica (mobile, web e Internet of Things). gRPC utilizza HTTP/2 per il trasporto, Protocol Buffers come linguaggio di descrizione dell'interfaccia e fornisce funzionalità quali autenticazione, streaming bidirezionale e controllo del flusso, binding bloccanti o non bloccanti. Genera binding client e server multipiattaforma per molti linguaggi. Gli scenari di utilizzo più comuni includono la connessione di servizi in un'architettura in stile microservizi o la connessione di client di dispositivi mobili a servizi backend.

A partire dal 2019, l'utilizzo di HTTP/2 da parte di gRPC rende impossibile implementare un client gRPC in un browser, richiedendo invece un proxy.

gRPC supporta l'utilizzo di Transport Layer Security (TLS) e l'autenticazione basata su token. Per l'autorizzazione basata su token, gRPC fornisce Server Interceptor e un Client Interceptor.

Per quanto riguarda i vantaggi, ad esempio, su un milione di richieste è molto più performante rispetto all'HTTP ed è molto comodo per lo streaming di dati. Ci si può sicuramente allontanare dall'implementare delle API con il modello REST visto che sono delle procedure remote che si chiamano direttamente

Utilizza protocol buffer che per la serializzazione/deserializzazione è molto più performante rispetto ad un JSON visto che sono byte scambiati. Tra gli svantaggi si può notare che il debugging, ad esempio, è più ostico perchè si necessita di avere un qualcosa che trasforma quei dati in formato leggibili (altrimenti rimangono byte).



# gRPC senza Istio opzione con headless service

Il motivo principale per cui è difficile bilanciare il traffico gRPC è che le persone vedono gRPC come HTTP ed è qui che inizia il problema, in base alla progettazione sono diversi, mentre HTTP crea e chiude le connessioni per richiesta, gRPC opera su un protocollo HTTP2 che funziona su una connessione TCP di lunga durata che rende più difficile il bilanciamento poiché più richieste passano attraverso la stessa connessione grazie alla funzione multiplex.

Senza la service mesh è necessario implementare soluzioni alternative per gestire efficacemente il bilanciamento del carico del traffico gRPC in Kubernetes. Queste soluzioni possono includere la configurazione diretta dei client per gestire più connessioni server o l'utilizzo di proxy come meccanismo di bilanciamento

### 1) proxy load balancing

Nel bilanciamento del carico proxy, il client invia gli RPC a un proxy Load Balancer (LB). Il LB distribuisce la chiamata RPC a uno dei server backend disponibili che implementano la logica effettiva per servire la chiamata. Il LB tiene traccia del carico su ogni backend e implementa algoritmi per distribuire equamente il carico. I client stessi non conoscono i server backend e possono non essere attendibili. Questa architettura è in genere utilizzata per servizi rivolti all'utente in cui i client da Internet aperta possono connettersi ai server

### 2) client-side load balancing

Nel bilanciamento del carico lato client, il client è a conoscenza di molti server backend e ne sceglie uno da utilizzare per ogni RPC. Se il client lo desidera, può implementare gli algoritmi di bilanciamento del carico in base al report di carico dal server. Per una distribuzione semplice, i client possono effettuare il round-robin delle richieste tra i server disponibili.

### Considerazioni finali

Kubernetes consente ai client di scoprire gli IP dei pod tramite ricerche DNS. Di solito, quando si esegue una ricerca DNS per un servizio, il server DNS restituisce un singolo IP, il cluser IP del servizio. Ma se si comunica a Kubernetes che non è necessario un cluster IP per il servizio (lo si fa impostando il campo clusterIP su None nella specifica del servizio), il server DNS restituirà gli IP dei pod anziché il singolo IP del servizio. Invece di restituire un singolo record DNS A, il server DNS restituirà più record A per il servizio, ognuno dei quali punta all'IP di un singolo pod che supporta il servizio in quel momento. I client possono quindi effettuare una semplice ricerca di record DNS A e ottenere gli IP di tutti i pod che fanno parte del servizio. Il client può quindi utilizzare tali informazioni per connettersi a uno, molti o tutti.

Impostando il campo clusterIP in una specifica di servizio su None, il servizio diventa headless, poiché Kubernetes non gli assegnerà un IP cluster tramite il quale i client potrebbero connettersi ai pod che lo supportano.