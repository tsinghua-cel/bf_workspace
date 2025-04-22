# BunnyFinder Repository

This repository contains the implementation for the paper "BunnyFinder: Finding Incentive Flaws for
Ethereum Consensus". It includes multiple modified [Prysm](https://github.com/OffchainLabs/prysm) repos in different versions and the attacks analyzed in the paper.

---

## 1. Ethical Considerations

All experiments are conducted on a local testnet (running on a single local machine). No experiments are conducted on the live Ethereum network. This repository does not uncover new vulnerabilities but instead analyzes known malicious reorganization attacks. We will not provide any additional exploit information.

---

## 2. Requirements

### 2.1. Hardware Dependencies

These experiments do not require any specialized hardware. The computer used in our setup is configured as follows:

- **CPU:** 16-core
- **RAM:** 32 GB
- **Storage:** 512 GB
- **Network Bandwidth:** 100 Mbps

---

### 2.2. Software Dependencies

#### 2.2.1. Docker

We use Docker to run our experiments. You can install Docker by following the instructions provided in the [official Docker documentation](https://docs.docker.com/engine/install/). Ensure that your Docker Engine version is at least Docker 24.

---

## 3. Run the Experiments Step by Step

We have designed two experiments that can be executed on different versions of Prysm. This repository only retains the operational scripts for Prysm v4.0.5 and Prysm v5.2.0, located in `v4/runtest.sh` and `v5/runtest.sh`, respectively.

- Simple experiments
- Strategy experiments

---

### 3.1. Build the Docker Image

After entering the repository directory (referred to as `$HOME`), run the following script to build the Docker image:

```shell
./build.sh
```

---

### 3.2. Simple experiments

To run the simple experiments, execute the following commands for each version of Prysm:

For Prysm v4.0.5:

```shell
./v4/runtest.sh normal
```

For Prysm v5.2.0:

```shell
./v5/runtest.sh normal
```

This experiment will take approximately 25 hours. If you need to stop the experiment early, you can execute `./v4/stop.sh` or `./v5/stop.sh`.

---

### 3.3. Strategy experiment

To run the strategy experiments, execute the following commands for each version of Prysm:

For Prysm v4.0.5:

```shell
./v4/runtest.sh strategy
```

For Prysm v5.2.0:

```shell
./v5/runtest.sh strategy
```

This experiment will take approximately 100 hours. If you need to stop the experiment early, you can execute `./v4/stop.sh` or `./v5/stop.sh`.

This experiment will take approximately 80 minutes. After completion, the results can be found in the `$HOME` directory.

---

### 3.5. Expected Output

The output of each experiment should look like this:

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
