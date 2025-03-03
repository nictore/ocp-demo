# Service Mesh Demo

## Requisiti

### Setup Operators

1. OpenShift Elasticsearch Operator
2. Red Hat OpenShift distributed tracing platform
3. Kiali Operator
4. Red Hat OpenShift Service Mesh

![image](images/1.png)

### Definizione ServiceMeshControlPlane

```yaml smcp.yaml
apiVersion: maistra.io/v2
kind: ServiceMeshControlPlane
metadata:
  name: basic
  namespace: istio-system
spec:
  addons:
    grafana:
      enabled: true
    jaeger:
      install:
        storage:
          type: Memory
    kiali:
      enabled: true
    prometheus:
      enabled: true
  gateways:
    openshiftRoute:
      enabled: true
  mode: MultiTenant
  policy:
    type: Istiod
  profiles:
    - default
  telemetry:
    type: Istiod
  tracing:
    sampling: 10000
    type: Jaeger
  version: v2.6
```

## Aggiunta di servizi in Service Mesh

### Definizione ServiceMeshMemberRoll

Questo oggetto fornisce agli amministratori di OpenShift Service Mesh un modo per delegare le autorizzazioni e per aggiungere progetti a una mesh. 

```yaml smmr.yaml
apiVersion: maistra.io/v1
kind: ServiceMeshMemberRoll
metadata:
  name: default
  namespace: istio-system
spec:
  members:
    - bookinfo
```

La Service Mesh in questo modo definisce anche network policy nella control plane della Service Mesh e nei namespace che partecipano alla mesh per consentire il traffico tra di essi.

```bash
oc get netpol -n bookinfo

istio-expose-route-basic
istio-mesh-basic
```

### Deploy Bookinfo

L'applicazione Bookinfo visualizza informazioni simili ad un negozio di libri online.
L'applicazione mostra una pagina che descrive il libro, i dettagli del libro (ISBN, numero di pagine e altre informazioni) e le recensioni del libro.

L'applicazione Bookinfo è composta da questi microservizi:

- Il microservizio productpage chiama i microservizi details e reviews per popolare la pagina.
- Il microservizio details contiene informazioni sui libri.
- Il microservizio reviews contiene le recensioni dei libri. Chiama anche il microservizio delle ratings.
- Il microservizio delle ratings contiene le informazioni sulle classifiche dei libri che accompagnano le recensioni.

Esistono tre versioni del microservizio reviews:

- La versione v1 non chiama il servizio di ratings.
- La versione v2 chiama il Servizio reviews e visualizza ogni valutazione con stelle nere.
- La versione v3 chiama il Servizio reviews e visualizza ogni valutazione con stelle rosse.

### Sidecar Injection

Annotations nei deployment per abilitare injection del proxy istio

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    [...]
spec:
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: 'true'
```

E' possibile sfruttare l'injection automatica dei sidecar configurando una label direttamente sul namespace:
`$ oc label namespace <nome_namespace> istio-injection=enabled`

### Versioning del deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    [...]
spec:
  template:
    metadata:
      labels:
        app: reviews
        version: v1
```

### Definizione ingressGateway

Una risorsa gateway descrive un bilanciatore di carico che opera ai margini della mesh ricevendo connessioni HTTP/TCP in entrata o in uscita. La specifica descrive:

- un set di porte che devono essere esposte
- il tipo di protocollo da utilizzare
- la configurazione SNI per il bilanciatore di carico e così via.

```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: bookinfo-gateway
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 8080
      name: http
      protocol: HTTP
    hosts: 
    - "*"
```

A differenza di una Ingress o Rotta standard, non include alcuna configurazione di routing del traffico. Il routing del traffico è invece configurato utilizzando l'oggetto virtualServices.

### Definizione virtualServices

Per specificare il routing e per far funzionare il gateway come previsto, devi anche associare il gateway a un virtualServices:

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: bookinfo
spec:
  hosts:
  - "*"
  gateways:
  - bookinfo-gateway
  http:
  - match:
    - uri:
        exact: /productpage
    - uri:
        prefix: /static
    - uri:
        exact: /login
    - uri:
        exact: /logout
    - uri:
        prefix: /api/v1/products
    route:
    - destination:
        host: productpage
        port:
          number: 9080
```

## Gestione del traffico

Per il microservizio reviews definiamo un oggetto DestinationRule per identificare i subset in base alla versione del deployment, configura quindi tre diversi sottoinsiemi:

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: reviews
  namespace: bookinfo
spec:
  host: reviews
  subsets:
    - labels:
        version: v1
      name: v1
    - labels:
        version: v2
      name: v2
    - labels:
        version: v3
      name: v3
  trafficPolicy:
    loadBalancer:
      simple: RANDOM
```

- Scenario 1: veicoliamo tutto il traffico solo per la versione v1 di review e poi solo per v2

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v1   #v2
```

- Scenario 2: veicoliamo il percentuale il traffico sulle 2 istanze v1 e v2
  
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
    - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v1
      weight: 80
    - destination:
        host: reviews
        subset: v2
      weight: 20
```

- Scenario 3: set header http, veicola traffico solo su v2 solo se corrisponde un determinato utente
  
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - match:
    - headers:
        end-user:
          exact: jason
    route:
    - destination:
        host: reviews
        subset: v2
  - route:
    - destination:
        host: reviews
        subset: v3
```

- Scenario 4: fault injection microservizio details // Aprire jaeger

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: details
spec:
  host: details
  subsets:
  - name: v1
    labels:
      version: v1
---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: details
spec:
  hosts:
  - details
  http:
  - fault:
      abort:
        httpStatus: 555
        percentage:
          value: 100
    route:
    - destination:
        host: details
        subset: v1
```

- Scenario 5: delay

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: details
spec:
  hosts:
  - details
  http:
  - fault:
      delay:
        fixedDelay: 7s
        percentage:
          value: 100
    route:
    - destination:
        host: details
        subset: v1
```

- Scenerio 6: mirroring del traffico // Aprire grafana

Il mirroring invia una copia del traffico live a un servizio mirrorato.

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
    - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v1
      weight: 100
    mirror:
        host: reviews
        subset: v2
    mirrorPercentage:
    value: 100.0
```

- Scenario 7: circuit breaking

Il circuit breaking è un pattern importante per la creazione di applicazioni microservice resilienti. Il circuit breaking consente di scrivere applicazioni che limitano l'impatto di guasti, picchi di latenza e altri effetti indesiderati delle peculiarità della rete.

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: details
spec:
  host: details
  subsets:
  - name: v1
    labels:
      version: v1
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 1
      http:
        http1MaxPendingRequests: 1
        maxRequestsPerConnection: 1
    outlierDetection:
      consecutive5xxErrors: 1
      interval: 1s
      baseEjectionTime: 3m
      maxEjectionPercent: 100
```

`maxConnections:` 1 e `http1MaxPendingRequests: 1`: Queste regole indicano che se si supera più di una connessione e contemporanea, dovrebbero verificarsi alcuni errori quando istio-proxy tenta di aprire ulteriori richieste e connessioni.


## gRPC


Più che altro ne conosco i vantaggi, per esempio sul milione di richieste è molto più performante rispetto all'HTTP ed è molto comodo per lo streaming di dati. Ti puoi sicuramente allontanare dall'implementare delle API con il modello REST visto che sono delle procedure remote che chiami direttamente

poi utilizza protocol buffer che per la serializzazione/deserializzazione è molto più performante rispetto ad un JSON visto che sono byte scambiati. Lo svantaggio è che il debugging è più ostico perchè devi avere un qualcosa che ti trasforma quei dati in formato leggibili (altrimenti rimangono byte) (modificato).


### Demo app rilasciata come Helm Chart con ArgoCD

Rilasciamo una applicazione Quarkus che stabilisce comunicazione client server con protocollo grpc per verificare la mancanza di bilanciamento (multiplexing)

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

Rilascio stessa applicazione Quarkus client-server con istio per verificare gestione delle connessioni multiple grazie agli envoy

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
    repoURL: https://github.com/nictore/grpc-demo.git
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

### Aggiornare ServiceMeshMemberRoll

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

## Security

Service Mesh consente di crittografare tutto il traffico senza richiedere modifiche al codice, senza complicati aggiornamenti di rete e senza installare/utilizzare strumenti esterni.

Per impostazione predefinita, mTLS in OpenShift Service Mesh viene abilitato e impostato **Permissive Mode**, i sidecar in Service Mesh accettano sia il traffico in plain text sia le connessioni crittografate tramite mTLS.

Abilitando mTLS nella mesh a livello di Control Plane (**ServiceMeshControlPlane**) è possibile proteggere i namespace dichiarati nella mesh. Per personalizzare le connessioni di crittografia del traffico i namespace devono essere configurati con le risorse **PeerAuthentication** e **DestinationRule**.

![image](images/4.png)

La CA di Istio genera automaticamente certificati per supportare le connessioni mTLS e li inietta nei pod dell'applicazione. In questo caso, l'utilizzo di mTLS comporta un ulteriore vantaggio poiché consente agli amministratori di creare regole di controllo degli accessi basate sui ruoli (RBAC) nel cluster OpenShift per specificare quale client può connettersi a quali servizi.

### Abilitare mTLS

```yaml
apiVersion: maistra.io/v2
kind: ServiceMeshControlPlane
spec:
  version: v2.2
  security:
    dataPlane:
      mtls: true
```

WARNING: Il passaggio alla modalità Enforce è facilmente attuabile se i workload non necessitano di comunicare con risorse esterne, perchè il traffico in egress dagli envoy viaggia cifrato


