# KServe

`kind create cluster --image kindest/node:v1.25.3`
`kind create cluster`

In order to install KServe, Istio and Cert Manager need to be installed first. Instructions are taken from here: `https://kserve.github.io/website/master/admin/kubernetes_deployment/`

`https://github.com/theofpa/kserve-tutorial/blob/main/quick_install_kind.sh`

## Install Istio

1. `helm repo add istio https://istio-release.storage.googleapis.com/charts` 
1. `helm repo update` 
1. Create Istio CRDs: `helm install istio-base istio/base -n istio-system --create-namespace` 
1. Install Istiod: `helm install istiod istio/istiod -n istio-system --wait` 
1. Check if Istio is installed: `kubectl get pods -n istio-system` 
1. You should see istiod pod running: 

```bash
NAME                      READY   STATUS    RESTARTS   AGE
istiod-56f5f79c79-295mq   1/1     Running   0          6m15s
```

1. Install istio-ingressgateway: `helm install istio-ingress istio/gateway -n istio-ingress --create-namespace`

## Install Cert Manager

1. `helm repo add jetstack https://charts.jetstack.io` 
1. `helm repo update`  
1. `helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace  --version v1.11.0 --set  installCRDs=true` 
1. Check if cert-manager is installed: `kubectl get pods -n cert-manager` 
1. You should see cert-manager pod running: 

```bash
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-589f57598d-s9zwh              1/1     Running   0          80s
cert-manager-cainjector-6b7bf5fc86-9gsh6   1/1     Running   0          80s
cert-manager-webhook-fbc478968-7q86q       1/1     Running   0          80s
``` 

## Install KServe

`helm repo add community-charts https://community-charts.github.io/helm-charts`

`helm repo update`

`helm install kserve community-charts/kserve`

`kubectl apply -f kserve.yaml`