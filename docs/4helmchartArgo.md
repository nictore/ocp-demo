# Helm Chart con Argo CD

## Demo app rilasciata come Helm Chart con ArgoCD

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

