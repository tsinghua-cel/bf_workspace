# BunnyFinder Repository

## Overview

This repository contains the implementation for our research paper "BunnyFinder: Finding Incentive Flaws for Ethereum Consensus." It provides comprehensive resources including:

- Complete implementation code for BunnyFinder
- Integration frameworks for multiple versions of the [Prysm](https://github.com/OffchainLabs/prysm) client
- Experimental datasets and results documented in our paper

## 1. Ethical Considerations

Our research adheres to responsible disclosure principles:

- All experiments are conducted exclusively on isolated local testnets
- No testing occurs on the live Ethereum network
- This repository analyzes only previously documented malicious reorganization attacks
- We do not disclose any new vulnerability information or additional exploit techniques

## 2. System Requirements

### 2.1. Hardware Specifications

Our experiments are designed to run on standard computing equipment. Our reference system configuration:

| Component | Specification       |
| --------- | ------------------- |
| CPU       | 16-core processor   |
| Memory    | 32 GB RAM           |
| Storage   | 512 GB              |
| Network   | 100 Mbps connection |

### 2.2. Software Dependencies

#### 2.2.1. Docker Environment

We utilize Docker to ensure consistent experimental environments. Requirements:

- Docker Engine version 24 or newer
- Installation instructions available in the [official Docker documentation](https://docs.docker.com/engine/install/)

## 3. Experimental Data and Execution Timeline

Our experimental data is stored in `./experiment/ndss_strategy.csv`. We utilized three c5.4xlarge AWS EC2 instances to conduct experiments over a total duration of 128 hours.

## 4. Experimental Methodology

Our research implements two distinct experimental approaches across multiple Prysm client versions. This repository includes operational scripts for Prysm v4.0.5 and v5.2.0, located in their respective directories.

### Experiment Types:

1. **Simple Experiments**: Generate and evaluate randomized attack strategies
2. **Strategy Experiments**: Extend and refine known attack vectors to develop advanced exploitation techniques

### 4.1. Environment Setup

Begin by building the required docker image in repository root directory:

```bash
# Build docker image with all dependencies
./build.sh
```

### 4.2. Simple Experiments

Execute randomized attack strategy tests:

```bash
# For Prysm v4.0.5
./v4/runtest.sh normal

# For Prysm v5.2.0
./v5/runtest.sh normal
```

**Note**: This experiment suite requires approximately 25 hours to complete. To terminate early, use `./v4/stop.sh` or `./v5/stop.sh`.

### 4.3. Strategy Experiments

Execute advanced attack vector tests:

```bash
# For Prysm v4.0.5
./v4/runtest.sh strategy

# For Prysm v5.2.0
./v5/runtest.sh strategy
```

**Note**: This comprehensive experiment suite requires approximately 100 hours to complete. To terminate early, use `./v4/stop.sh` or `./v5/stop.sh`.

### 4.4. Expected Output

A successful experiment launch will produce output similar to:

```
[+] Running 22/22
 ✔ Network basic_meta              Created                  0.0s
 ✔ Container basic-execute5-1      Started                  0.1s
 ✔ Container basic-execute3-1      Started                  0.1s
 ✔ Container basic-execute4-1      Started                  0.1s
 ✔ Container basic-execute1-1      Started                  0.1s
 ✔ Container basic-execute2-1      Started                  0.1s
 ✔ Container basic-attacker1-1     Started                  0.1s
 ✔ Container basic-beacon3-1       Started                  0.1s
 ✔ Container basic-beacon1-1       Started                  0.1s
 ✔ Container basic-beacon2-1       Started                  0.1s
 ✔ Container basic-beacon5-1       Started                  0.1s
 ✔ Container basic-beacon4-1       Started                  0.1s
 ✔ Container basic-validator4-1    Started                  0.1s
 ✔ Container basic-validator5-1    Started                  0.1s
 ✔ Container basic-validator2-1    Started                  0.1s
 ✔ Container basic-validator1-1    Started                  0.1s
 ✔ Container basic-validator3-1    Started                  0.1s
```

### 4.5. Result Querier

At the end of the experiment, the test results will be summarized and printed to the screen. An example output is shown below:

```txt
Executing metrics queries...

mysql: [Warning] Using a password on the command line interface can be insecure.
+----------+-------+
| Metric   | Value |
+----------+-------+
| Metric 1 |  2586 |
| Metric 2 |     1 |
| Metric 3 |     2 |
| Metric 4 |    43 |
+----------+-------+

Execution completed - container has been automatically removed
```

## 5. Run Individual Attacks

If you want to run a specific attack individually, you can use the following commands. Each attack will last for 30 minutes or an hour.

Run the modified exante reorg attack:
```bash
./attack.sh exante
```

Run the sandwich reorg attack:
```bash
./attack.sh sandwich
```

Run the unrealized justification reorg attack:
```bash
./attack.sh unrealized
```

Run the justification withholding reorg attack:
```bash
./attack.sh withholding
```

Run the staircase attack:
```bash
./attack.sh staircase
```

Run the selfish mining attack:
```bash
./attack.sh selfish
```

Run the staircase attack-II:
```bash
./attack.sh staircase-ii
```

Run the pyrrhic victory attack:
```bash
./attack.sh pyrrhic-victory
```
