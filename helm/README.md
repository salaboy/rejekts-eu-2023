# Rejekts EU 2023 :: Demo using Knative and Dapr

This app requires to have installed Knative Serving in the cluster and Dapr and a Redis Instance called `redis` installed in the namespace where the chart is deployed. 

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis --set architecture=standalone
```

For more info check: [https://github.com/salaboy/rejekts-eu-2023](https://github.com/salaboy/rejekts-eu-2023)