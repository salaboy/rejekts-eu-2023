# Cloud-Native Rejekts EU 2023 :: Step-by-step tutorial

[Lessons learnt from creating platforms on Kubernetes](https://cfp.cloud-native.rejekts.io/cloud-native-rejekts-eu-amsterdam-2023/talk/PTCMVR/)

On this short tutorial we will be looking at three main aspects of creating platforms on top of Kubernetes. 

To build successful Platforms on top of Kubernetes you need to: 

- Glue things together: reduce the cognitive load, be ready to pivot. Understand and join the Cloud Native and CNCF ecosystem and projects to understand where the industry is going and what other companies are doing
- Understand your teams:  and then provide self-service APIs for them to do their work (no more Jira OPS!)
- A powerful End User Experience: will boost your teams productivity. Make sure that you have tailored experiences for example: Developer Experiences targeting specific tech stacks or Data Scientist workflows.

Before jumping into the sections make sure you follow the [prerequisites and installation section here](prerequisites.md).


For the purpose of this tutorial are creating a platform to help development teams and data scientist to work together, by exposing clear interfaces that they can use to provision the resources that they need and then have the tools to do the work. 

## Glue things together

Keeping an eye on the CNCF ecosystem is a full time job, but if you are serious about adopting Kubernetes you want to stay up to date to make sure that you levarage what these projects are doing, so you don't need to build your in-house solution. 

In this section will we look at creating our Platform using a set of tools that accomodate different teams with different expectations. 

For this we will install the following tools into our Kubernetes Cluster that we will call the Platform Cluster: 

- Crossplane
- ArgoCD
- Dapr

These three very popular tools provide a set of key features that enable us to build more complex platforms on top of Kubernetes. 

```
kind create cluster
```

Let's install [Crossplane](https://crossplane.io) into its own namespace using Helm: 

```

helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --wait
```

Install the `kubectl crossplane` plugin: 

```
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
sudo mv kubectl-crossplane /usr/local/bin
```

Then install the Crossplane Helm provider: 
```
kubectl crossplane install provider crossplane/provider-helm:v0.10.0
```

We need to get the correct ServiceAccount to create a new ClusterRoleBinding so the Helm Provider can install Charts on our behalf. 

```
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

```
kubectl apply -f crossplane/config/helm-provider-config.yaml
```

Let's install ArgoCD into our platform cluster with: 

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```


You can access the ArgoCD dashboard by using `kubectl port-forward` (in a separate terminal):

```
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Then you can point your browser to [http://localhost:8080](http://localhost:8080)

And you can get the `admin` user's password by running the following command: 

```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

You can create the production namespace by running: 

```
kubectl create ns production
```

In our production environment we will install a Redis instance using helm. 

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-redis bitnami/redis --set architecture=standalone -n production
```


## Understand your teams

Once we have the main tools to build platforms we need to combine them in a way that make sense for our teams. If we have Developers and Data Scientists we cannot give the same tools to them, as their work is completely different and the tools that they use have different requirements. 

In this section we will be creating two different Crossplane Compositions. One for our Development Teams to create Development Environments, and the other one for our Data Scientists and their exotic tools. 

The [crossplane](crossplane/) directory contains one `CompositeResourceDefinition` and two `Composition`s that enable both our Developers and our Data Scientists to create environment for them to work. 

Let's install these resources:

```
kubectl apply -f crossplane/env-resource-definition.yaml
kubectl apply -f crossplane/composition-devenv.yaml
kubectl apply -f crossplane/composition-mlenv.yaml
```

Now we can request new ML and Dev Environments by just creating Environment Resources and using labels to define what kind of Environment we want: 

```
kubectl apply -f team-b-ml-env.yaml
```


## A powerful end user experience



