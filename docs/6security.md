# Security

Service Mesh consente di crittografare tutto il traffico senza richiedere modifiche al codice, senza complicati aggiornamenti di rete e senza installare/utilizzare strumenti esterni.

Per impostazione predefinita, mTLS in OpenShift Service Mesh viene abilitato e impostato Permissive Mode, i sidecar in Service Mesh accettano sia il traffico in plain text sia le connessioni crittografate tramite mTLS.

Abilitando mTLS nella mesh a livello di Control Plane (ServiceMeshControlPlane) è possibile proteggere i namespace dichiarati nella mesh. Per personalizzare le connessioni di crittografia del traffico i namespace devono essere configurati con le risorse PeerAuthentication e DestinationRule.

![image info](images/4.png)

La CA di Istio genera automaticamente certificati per supportare le connessioni mTLS e li inietta nei pod dell'applicazione. In questo caso, l'utilizzo di mTLS comporta un ulteriore vantaggio poiché consente agli amministratori di creare regole di controllo degli accessi basate sui ruoli (RBAC) nel cluster OpenShift per specificare quale client può connettersi a quali servizi.

# Abilitare mTLS

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