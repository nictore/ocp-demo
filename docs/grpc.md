# gRPC

Utilizza protocol buffer che per la serializzazione/deserializzazione è molto più performante rispetto ad un JSON visto che sono byte scambiati. Tra gli svantaggi si può notare che il debugging, ad esempio, è più ostico perchè si necessita di avere un qualcosa che trasforma quei dati in formato leggibili (altrimenti rimangono byte).

## gRPC senza Istio opzione con headless service

Il motivo principale per cui è difficile bilanciare il traffico gRPC è che le persone vedono gRPC come HTTP ed è qui che inizia il problema, in base alla progettazione sono diversi, mentre HTTP crea e chiude le connessioni per richiesta, gRPC opera su un protocollo HTTP2 che funziona su una connessione TCP di lunga durata che rende più difficile il bilanciamento poiché più richieste passano attraverso la stessa connessione grazie alla funzione multiplex.

Senza la service mesh è necessario implementare soluzioni alternative per gestire efficacemente il bilanciamento del carico del traffico gRPC in Kubernetes. Queste soluzioni possono includere la configurazione diretta dei client per gestire più connessioni server o l'utilizzo di proxy come meccanismo di bilanciamento

### proxy load balancing

Nel bilanciamento del carico proxy, il client invia gli RPC a un proxy Load Balancer (LB).
Il LB distribuisce la chiamata RPC a uno dei server backend disponibili che implementano la logica effettiva per servire la chiamata.
Il LB tiene traccia del carico su ogni backend e implementa algoritmi per distribuire equamente il carico.
I client stessi non conoscono i server backend e possono non essere attendibili.

### client-side load balancing

Nel bilanciamento del carico lato client, il client è a conoscenza di molti server backend e ne sceglie uno da utilizzare per ogni RPC.
Se il client lo desidera, può implementare gli algoritmi di bilanciamento del carico in base al report di carico dal server.
Per una distribuzione semplice, i client possono effettuare il round-robin delle richieste tra i server disponibili.

## App GRPC rilasciata come Helm Chart in ArgoCD

### 1. Creazione namespace grpc-demo

Aggiungere annotazione:

```yaml
labels: argocd.argoproj.io/managed-by: openshift-gitops
```

### 2. Rilascio applicazione Quarkus

Rilasciamo una applicazione Quarkus che stabilisce una comunicazione client-server con protocollo grpc per verificare la mancanza di bilanciamento (multiplexing)

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: grpc-demo
  namespace: openshift-gitops
spec:
  destination:
    namespace: grpc-demo
    server: https://kubernetes.default.svc
  source:
    path: grpc/helm-charts/grpc-demo-services
    repoURL: https://github.com/nictore/ocp-demo.git
    targetRevision: HEAD
    helm:
      releaseName: grpc-demo-services
      valueFiles:
      - values.yaml
  sources: []
  project: default
  syncPolicy:
    syncOptions:
      - CreateNamespace=true
```

### 3. Test di comunicazione

Accedere al client grpc ed eseguire:

```bash
oc -n test exec deploy/grpc-client -- curl http://localhost:8080/hello/grpc
```

Verificare i log lato server, solo una istanza ha ricevuto i messaggi

### 4. Osservazioni

**Perché il traffico gRPC non è bilanciato correttamente in Kubernetes?**

Il motivo principale per cui è difficile bilanciare il traffico gRPC è che spesso... si pensa che gRPC sia come HTTP ed è qui che inizia il problema
Mentre HTTP crea e chiude connessioni per richiesta, gRPC opera su protocollo HTTP2 che funziona su una connessione TCP di lunga durata rendendo più difficile il bilanciamento poiché più richieste passano attraverso la stessa connessione grazie alla funzionalità di multiplexing. Tuttavia, ci sono altri errori comuni:

- Configurazione errata del client gRPC
- Configurazione errata del servizio Kubernetes

#### Configurazione errata del client gRPC

Il client gRPC come configurazione predefinita prevede un tipo di connessione 1–1, in ambiente produttivo non funziona come vorremmo.
Il client gRPC predefinito offre la possibilità di connettersi con un semplice record IP/DNS che crea una sola connessione con il servizio di destinazione.
Ecco perché è necessario effettuare una configurazione diversa.

Impostazione di default:

```go
func main(){ 
  conn, err := grpc.Dial("my-domain:50051", grpc.WithInsecure()) 
  if err != nil { 
    log.Fatalf("errore di connessione con il server gRPC: %v", err) 
  }
  ...
```

Impostazione consigliata:

```go
func main(){ 
  addr := fmt.Sprintf("%s:///%s", "dns", " my-domain :50051") 
  conn, err := grpc.Dial(addr, grpc.WithInsecure(),grpc.WithBalancerName(roundrobin.Name)) 
  if err != nil { 
    log.Fatalf("connessione non riuscita: %v", err) 
  } 
  ...
```

In questo modo nel caso in cui il nostro client si connetta a più server, ora il nostro client gRPC è in grado di bilanciare le richieste in base all'algoritmo di bilanciamento scelto.

## App GRPC con Istio

Rilasciamo la stessa applicazione Quarkus client-server con istio per verificare gestione delle connessioni multiple grazie agli envoy

### 1. Creazione namespace grpc-demo-istio

Aggiungere annotazione:

```yaml
labels: argocd.argoproj.io/managed-by: openshift-gitops
```

### 2. Aggiornare ServiceMeshMemberRoll

```yaml
apiVersion: maistra.io/v1
kind: ServiceMeshMemberRoll
metadata:
  name: default
  namespace: istio-system
spec:
  members:
    - bookinfo
    - grpc-demo-istio
```

### 3. Rilascio app gRPC con Istio

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: grpc-demo-istio
  namespace: openshift-gitops
spec:
  destination:
    namespace: grpc-demo-istio
    server: https://kubernetes.default.svc
  source:
    path: grpc/helm-charts/grpc-demo-services-istio
    repoURL: https://github.com/nictore/ocp-demo.git
    targetRevision: HEAD
    helm:
      releaseName: grpc-demo-services-istio
      valueFiles:
      - values.yaml
  sources: []
  project: default
  syncPolicy:
    syncOptions:
      - CreateNamespace=true
```

### 4. Test di comunicazione

Accedere al client grpc ed eseguire:

```bash
oc -n test exec deploy/grpc-client -- curl http://localhost:8080/hello/grpc
```

Verificare i log lato server, traffico bilanciato correttamente

### Considerazioni finali

Kubernetes consente ai client di scoprire gli IP dei pod tramite ricerche DNS. Di solito, quando si esegue una ricerca DNS per un servizio, il server DNS restituisce un singolo IP, il cluser IP del servizio. Ma se si comunica a Kubernetes che non è necessario un cluster IP per il servizio (lo si fa impostando il campo clusterIP su None nella specifica del servizio), il server DNS restituirà gli IP dei pod anziché il singolo IP del servizio. Invece di restituire un singolo record DNS A, il server DNS restituirà più record A per il servizio, ognuno dei quali punta all'IP di un singolo pod che supporta il servizio in quel momento. I client possono quindi effettuare una semplice ricerca di record DNS A e ottenere gli IP di tutti i pod che fanno parte del servizio. Il client può quindi utilizzare tali informazioni per connettersi a uno, molti o tutti.

Impostando il campo clusterIP in una specifica di servizio su None, il servizio diventa headless, poiché Kubernetes non gli assegnerà un IP cluster tramite il quale i client potrebbero connettersi ai pod che lo supportano.