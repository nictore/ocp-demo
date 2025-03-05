# Helm Chart con Istio

## Demo app rilasciata come Helm Chart con Istio

Rilasciamo stessa applicazione Quarkus client-server con istio per verificare gestione delle connessioni multiple grazie agli envoy

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