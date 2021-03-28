GOOS=linux go build -o ./app .
docker build .
kubectl apply -f main.yaml
kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default
