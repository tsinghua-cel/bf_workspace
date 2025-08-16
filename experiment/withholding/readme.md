# Withholding Test Data Description
`data.tar.gz` contains the data and logs of beacon nodes, validators, and bunnyfinder during the test. After extraction, the directory structure is as follows:
```txt
└── withholding
    ├── data
    │   ├── attacker1
    │   ├── beacon1
    │   ├── beacon2
    │   ├── beacon3
    │   ├── beacon4
    │   ├── beacon5
    │   ├── validator1
    │   ├── validator2
    │   ├── validator3
    │   ├── validator4
    │   └── validator5
    └── database
        └── data
```
beacon1 ~ beacon5 are the data directories of 5 beacon nodes, validator1 ~ validator5 are the keystore and password files of 5 validators, and attacker1 is the data directory of bunnyfinder, which only retains runtime logs.
The data collected by bunnyfinder is stored in `database/data`.

beacon1 and validator1 are malicious nodes. By searching the reorg logs in `beacon2/`, the occurrence of the withholding attack can be observed.
```shell
grep -nr "reorg" ./withholding/beacon2/
```

# Retrieve data in MySQL
Boot up mysql with given data.
```shell
docker compose -f mysql.yml up -d 
```

Connect to mysql with password `12345678`.
```shell
docker exec -it withholding-ethmysql-1 mysql -u root -p eth
```

