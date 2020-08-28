# pizza-apiserver

# update the generated files 

bash ./hack/update-codegen.sh

# to build the code 
docker build -t akankshakumari393/custom-pizza-api-server:latest .

# Pushing the code to dockerhub
docker push akankshakumari393/custom-pizza-api-server:latest



