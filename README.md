### Build Project:
```
go build -o wc main.go
```

### Create CA Cert:
```
./wc initCA
```

###  Generate Server Cert:
```
./wc generateServerCertificate
```

### Generate Client Cert:
```
./wc generateClientCertificate
```

### Run Server:
```
./wc runServer <ca.crt path> <server.crt path> <server.key path>
```

### Run Client:
```
./wc runClient <ca.crt path> <client.crt path> <client.key path> <FirstOpearand> <SecondOperand>
```