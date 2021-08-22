# Loki with local s3 bucket

If you want to give this operator a quick test without having to setup a aws account you can use something like Minio.
This guide is just to test the loki operator, nothing else, don't use this setup in production.

## Setup Minio

Minio have a simple cli tool to make life easier, you can install the cli by following the [docs](https://github.com/minio/operator#1-install-the-minio-operator)

or you can use [krew](https://github.com/kubernetes-sigs/krew), if you have it installed.

```shell
kubectl krew install minio
```

### Install operator and tenant

This command will create a two deployments in minio-operator namespace.

```shell
kubectl minio init
```

This command will create a tenant with a few basic configs.

I use the loki namespace, since this is where I have installed the operator.

```shell
kubectl create ns loki
kubectl minio tenant create minio-tenant-1       \
  --servers                 1                    \
  --volumes                 4                    \
  --capacity                1Gi                 \
  --storage-class           standard             \
  --namespace               loki
```

The minio tenant takes up to 5 minutes to setup so be patient:

```shell
kubectl get tenants.minio.min.io -w -n loki
```

In the end you should have something like this:

```shell
kubectl get pods
NAME                                      READY   STATUS    RESTARTS   AGE
minio-tenant-1-console-859c7b6f48-9qsv4   1/1     Running   0          14m
minio-tenant-1-console-859c7b6f48-zkkq2   1/1     Running   0          14m
minio-tenant-1-ss-0-0                     1/1     Running   0          15m

kubectl get svc
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
minio                    ClusterIP   10.96.209.21    <none>        443/TCP    17m
minio-tenant-1-console   ClusterIP   10.96.192.143   <none>        9443/TCP   15m
minio-tenant-1-hl        ClusterIP   None            <none>        9000/TCP   17m

# Ignore that it says Waiting for pods to be ready, I think this is due to that i really wants for 4 replicas and not 1
kubectl get tenants.minio.min.io
NAME             STATE                          AGE
minio-tenant-1   Waiting for Pods to be ready   17m
```

### Create S3 bucket

First lets get the aws access and secret key:

```shell
# access key
kubectl get secret minio-tenant-1-creds-secret -o jsonpath={.data.accesskey} |base64 -d
# secret key
kubectl get secret minio-tenant-1-creds-secret -o jsonpath={.data.secretkey} |base64 -d

# port-forward to the s3 bucket
kubectl port-forward service/minio 9000:443
```

Now you can reach the s3 bucket from your browser on `localhost:9000`, use the access key and secret key to login.

Login to the console and in the bottom right you will find a `+` sign, push it and create a bucket.
To be able to use the script out of the box, I suggest calling the bucket loki.

## Configure loki operator

Use the config defined in hack, to create the needed secret file use `minio-deploy-example-secret.sh`.

## Connect to loki

You should now be able to use your grafana instance to connect to loki:

http://loki-query-frontend-http-lokistack-dev.loki.svc.cluster.local:3100
