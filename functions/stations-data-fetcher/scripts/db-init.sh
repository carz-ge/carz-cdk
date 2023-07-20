mkdir $HOME/.postgresql
wget -O $HOME/.postgresql/root.crt 'https://cockroachlabs.cloud/clusters/b548e392-7d3c-45b4-8d87-2e7e69c9aa03/cert'
#curl --create-dirs -o $HOME/.postgresql/root.crt 'https://cockroachlabs.cloud/clusters/b548e392-7d3c-45b4-8d87-2e7e69c9aa03/cert'