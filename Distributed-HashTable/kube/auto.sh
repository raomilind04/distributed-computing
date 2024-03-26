#! /bin/bash

kubectl delete -f ConfigMap.yml
kubectl delete -f server-config.yml
kubectl delete -f pod-config1.yml
kubectl delete -f pod-config4.yml
kubectl delete -f pod-config9.yml
kubectl delete -f pod-config11.yml
kubectl delete -f pod-config14.yml
kubectl delete -f pod-config18.yml
kubectl delete -f pod-config20.yml
kubectl delete -f pod-config21.yml
kubectl delete -f pod-config28.yml

sleep 5

kubectl apply -f ConfigMap.yml
kubectl apply -f server-config.yml
kubectl apply -f pod-config1.yml
kubectl apply -f pod-config4.yml
kubectl apply -f pod-config9.yml
kubectl apply -f pod-config11.yml
kubectl apply -f pod-config14.yml
kubectl apply -f pod-config18.yml
kubectl apply -f pod-config20.yml
kubectl apply -f pod-config21.yml
kubectl apply -f pod-config28.yml
