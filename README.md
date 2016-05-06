# Fake CAS

Download the binary from [here](https://github.com/CenterForOpenScience/fakecas/releases/download/0.7.0/fakecas)

```bash
cd ~/Downloads # cd to where you downloaded the file to
chmod +x fakecas # Make the server executable
./fakecas # Run the server

./fakecas -h  # Print possible configuration options
# Usage of ./fakecas:
#   -dbaddress="localhost:27017": The address of your mongodb. ie: localhost:27017
#   -dbname="osf20130903": The name of your OSF database
#   -host="localhost:8080": The host to bind to
```
