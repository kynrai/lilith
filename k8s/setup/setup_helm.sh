#!/usr/bin/env bash

# Must have kubectl configured
# Must have helm installed
kubectl delete svc tiller-deploy -n kube-system
kubectl -n kube-system delete deploy tiller-deploy
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller

# Install cert manager
helm install \
    --name cert-manager \
    --namespace kube-system \
    stable/cert-manager
