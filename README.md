# BunnyFinder: Finding Incentive Flaws for Ethereum Consensus

## Overview

This repository contains the implementation for our paper "BunnyFinder: Finding Incentive Flaws for Ethereum Consensus." BunnyFinder is a framework for finding incentive flaws in Ethereum with little manual effort, inspired by failure injection techniques commonly used in software testing. The repository includes:

- Implementation code for BunnyFinder

- Two versions of the [Prysm](https://github.com/OffchainLabs/prysm) client (v4.0.2 and v5.2.0) with injection points
- Experimental datasets and results from our paper

## Ethical Considerations

Our research adheres to responsible disclosure principles:

- All experiments are conducted exclusively on isolated local testnets
- No testing occurs on the live Ethereum network  
- The attacks we analyzed do not disclose any new vulnerability information or additional exploit techniques

## System Requirements

### Hardware Dependencies

The experiments do not require any specialized hardware. Our reference system configuration:

| Component | Specification       |
| --------- | ------------------- |
| CPU       | 8-core processor    |
| Memory    | 16 GB RAM           |
| Storage   | 100 GB              |
| Network   | 100 Mbps connection |

### Software Dependencies

Our experiments require:
- Ubuntu 22.04 or later
- Docker Engine version 24.0.6 or higher 
- docker-compose plugin
- Kurtosis framework (optional)

Installation instructions are available in the [official Docker documentation](https://docs.docker.com/engine/install/)

## Installation & Configuration

After installing Docker, follow these steps:

1. Git clone the repository:
   ```bash
   git clone https://github.com/tsinghua-cel/bf_workspace.git
   ```

2. Enter the repository directory:
   ```bash
   cd bf_workspace
   ```

3. Build the required Docker image in the repository root directory:
   ```bash
   ./build.sh
   ```

## Quick Start

After building the Docker image, run the basic test with:

```bash
./attack.sh none
```

Expected output:
```console
casetype is none
[+] Running 2/2
 ✔ Network case_default       Created     0.2s 
 ✔ Container case-ethmysql-1  Started     0.6s 
run strategy none
INFO[0000] Specified a chain config file: /root/config/config.yml  prefix=genesis
INFO[0000] No genesis time specified, defaulting to now()  prefix=genesis
INFO[0000] Delaying genesis 1752219566 by 15 seconds prefix=genesis
INFO[0000] Genesis is now 1752219581                 prefix=genesis
INFO[0000] Setting fork geth times cancun=1752219581 prague=1761819581 prefix=genesis shanghai=1752219581
INFO[0000] Done writing genesis state to /root/config/genesis.ssz  prefix=genesis
INFO[0000] Command completed                         prefix=genesis
[+] Running 17/17
 ✔ Network none_meta            Created   0.2s 
 ✔ Container none-execute3-1    Started   1.3s 
 ✔ Container none-execute4-1    Started   1.1s 
 ✔ Container none-execute2-1    Started   0.6s 
 ✔ Container none-attacker1-1   Started   1.5s 
 ✔ Container none-execute5-1    Started   1.5s 
 ✔ Container none-execute1-1    Started   1.5s 
 ✔ Container none-beacon1-1     Started   3.3s 
 ✔ Container none-beacon3-1     Started   3.0s 
 ✔ Container none-beacon2-1     Started   1.7s 
 ✔ Container none-beacon4-1     Started   2.0s 
 ✔ Container none-beacon5-1     Started   3.3s 
 ✔ Container none-validator3-1  Started   4.5s 
 ✔ Container none-validator5-1  Started   4.8s 
 ✔ Container none-validator4-1  Started   3.8s 
 ✔ Container none-validator1-1  Started   4.8s 
 ✔ Container none-validator2-1  Started   3.5s 
wait 360 seconds
```

After running for six minutes, the experiment stops with output like:
```console
[+] Running 17/17
 ✔ Container none-attacker1-1   Removed   0.7s 
 ✔ Container none-validator4-1  Removed  12.0s 
 ✔ Container none-validator3-1  Removed  11.5s 
 ✔ Container none-validator2-1  Removed  11.8s 
 ✔ Container none-validator5-1  Removed  12.3s 
 ✔ Container none-validator1-1  Removed  12.1s 
 ✔ Container none-beacon3-1     Removed  10.5s 
 ✔ Container none-beacon2-1     Removed  10.9s 
 ✔ Container none-beacon4-1     Removed  10.8s 
 ✔ Container none-beacon1-1     Removed  11.3s 
 ✔ Container none-beacon5-1     Removed  11.0s 
 ✔ Container none-execute3-1    Removed  10.6s 
 ✔ Container none-execute2-1    Removed  11.2s 
 ✔ Container none-execute4-1    Removed  10.8s 
 ✔ Container none-execute5-1    Removed  10.8s 
 ✔ Container none-execute1-1    Removed  10.9s 
 ✔ Network none_meta            Removed   0.8s 
result collect
[+] Running 2/2
 ✔ Container case-ethmysql-1  Removed     2.2s 
 ✔ Network case_default       Removed     0.9s
```

## Experiments

### Experiment 1: Known Incentive Attacks (E1)
**Duration**: 30 human-minutes + 5 compute-hours

Reproduce five known incentive attacks:

#### Ex-ante Reorg Attack (Prysm 5.2.0)
```bash
./attack.sh exante
```

#### Sandwich Reorg Attack (Prysm 5.2.0)  
```bash
./attack.sh sandwich
```

#### Unrealized Justification Attack (Prysm 4.0.5)
```bash
./attack.sh unrealized
```

#### Justification Withholding Attack (Prysm 4.0.5)
```bash
./attack.sh withholding
```

#### Staircase Attack (Prysm 4.0.5)
```bash
./attack.sh staircase
```

**Expected Results**: After completion, you should see output indicating reorganized blocks, confirming successful attack reproduction:

```console
[+] Running 17/17
 ✔ Container exante-validator5-1  Removed 12.7s
 ✔ Container exante-validator3-1  Removed 12.5s
 ✔ Container exante-validator4-1  Removed 12.3s
 ✔ Container exante-validator1-1  Removed 11.9s
 ✔ Container exante-attacker1-1   Removed  0.8s
 ✔ Container exante-validator2-1  Removed 12.2s
 ✔ Container exante-beacon1-1     Removed 10.8s
 ✔ Container exante-beacon2-1     Removed 11.1s
 ✔ Container exante-beacon4-1     Removed 11.6s
 ✔ Container exante-beacon3-1     Removed 11.1s
 ✔ Container exante-beacon5-1     Removed 11.4s
 ✔ Container exante-execute1-1    Removed 10.8s
 ✔ Container exante-execute2-1    Removed 10.8s
 ✔ Container exante-execute3-1    Removed 10.8s
 ✔ Container exante-execute4-1    Removed 10.8s
 ✔ Container exante-execute5-1    Removed 10.8s
 ✔ Network exante_meta            Removed  0.8s
result collect
exante attack occurs reorganize blocks in slot 8.
exante attack occurs reorganize blocks in slot 23.
...
test finished and all nodes data in $HOME/results/exante
```

### Experiment 2: New Incentive Attacks (E2)
**Duration**: 30 human-minutes + 3 compute-hours

Identify three previously unknown incentive attacks:

#### Selfish Mining Attack (Prysm 5.2.0)
```bash
./attack.sh selfish
```

#### Staircase Attack-II (Prysm 5.2.0)
```bash
./attack.sh staircase-ii
```

#### Pyrrhic Victory Attack (Prysm 5.2.0)
```bash
./attack.sh pyrrhic-victory
```

**Expected Results**: Similar output showing reorganized blocks during the attack execution.

```console
[+] Running 17/17
 - Container staircaseii-validator5-1  Removed 12.7s
 - Container staircaseii-validator3-1  Removed 12.5s
 - Container staircaseii-validator4-1  Removed 12.3s
 - Container staircaseii-validator1-1  Removed 11.9s
 - Container staircaseii-attacker1-1   Removed  0.8s
 - Container staircaseii-validator2-1  Removed 12.2s
 - Container staircaseii-beacon1-1     Removed 10.8s
 - Container staircaseii-beacon2-1     Removed 11.1s
 - Container staircaseii-beacon4-1     Removed 11.6s
 - Container staircaseii-beacon3-1     Removed 11.1s
 - Container staircaseii-beacon5-1     Removed 11.4s
 - Container staircaseii-execute1-1    Removed 10.8s
 - Container staircaseii-execute2-1    Removed 10.8s
 - Container staircaseii-execute3-1    Removed 10.8s
 - Container staircaseii-execute4-1    Removed 10.8s
 - Container staircaseii-execute5-1    Removed 10.8s
 - Network staircaseii_meta            Removed  0.8s
result collect
staircaseii attack occurs reorganize blocks in slot 
152-216.
staircaseii attack occurs reorganize blocks in slot 
542-595.
test finished and all nodes data in /home/ec2-user/
bf_workspace/results/staircaseii
[+] Running 2/2
 - Container case-ethmysql-1  Removed           2.0s
 - Network case_default       Removed
```



### Experiment 3: Attack Database Analysis (E3)
**Duration**: 30 human-minutes

Query and analyze attack instances from our comprehensive attack database containing 7,991 completed attack strategy records.

#### Database Connection
```bash
export MYSQL_PASSWORD=<Please contact us to get the password> 
./tool/connect_ndss.sh
```

#### Example Queries

Count total attack strategies:
```sql
SELECT COUNT(1) FROM t_strategy WHERE is_end=1;
```

Top 10 most effective attacks by honest validator loss rate:
```sql
SELECT uuid,category,honest_lose_rate_avg,attacker_lose_rate_avg 
FROM t_strategy 
ORDER BY honest_lose_rate_avg DESC 
LIMIT 10;
```

Query specific strategy details:
```sql
SELECT content FROM t_strategy WHERE uuid = 'your_uuid_here';
```

**Expected Results**:
```console
mysql> SELECT COUNT(1) FROM t_strategy where is_end=1;
+----------+
| COUNT(1) |
+----------+
|     7991 |
+----------+
1 row in set (0.24 sec)

mysql> SELECT uuid,category,honest_lose_rate_avg,attacker_lose_rate_avg FROM t_strategy ORDER BY honest_lose_rate_avg DESC LIMIT 10;
+--------------------------------------+-----------------+----------------------+------------------------+
| uuid                                 | category        | honest_lose_rate_avg | attacker_lose_rate_avg |
+--------------------------------------+-----------------+----------------------+------------------------+
| be0392e4-2af5-4328-ae73-75cb940183fb | ext_unrealized  |                    2 |                      2 |
| 8165a654-fed8-4c7c-8735-b1cb2665fa39 | ext_withholding |   1.4351343907591927 |     0.4257843676895559 |
| 0908f1a1-7106-4e89-bbe5-b5b004842708 | ext_withholding |   1.0404429370482193 |    0.36005933160736164 |
| 36917a9b-7589-422f-a4d5-12c13966a3fe | ext_unrealized  |   0.8446009559734197 |    0.23046993257190004 |
| 1d80817a-3b87-4961-ac24-8c501e1b99f4 | ext_exante      |   0.8082500263074806 |     0.3611931605441617 |
+--------------------------------------+-----------------+----------------------+------------------------+
```

## Support

For questions about the artifact or research, please contact the authors through the paper submission system.

## License

This project is licensed under the same terms as the original Prysm client. See LICENSE files in respective directories.

