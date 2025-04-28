# Strategy Optimization Framework for PoS Blockchain Attacks

This repository provides a reinforcement learning-based strategy optimization framework for simulating and optimizing adversarial attacks (Selfish Mining, Staircase Attack-II) on Ethereum 2.0 Proof-of-Stake systems.

We model the environment, agent, and blockchain dynamics carefully, using a recurrent actor-critic architecture (LSTM) and train the agent using Proximal Policy Optimization (PPO).

---

## 📂 Project Structure

```plaintext
env/
  - environment.py        # Ethereum PoS environment (Selfish Mining, Staircase Attack-II)
  - hlmdghost.py          # HLMD GHOST for selfish mining
  - casper.py             # Casper FFG finality gadget for staircase attack

model/
  - actor_critic.py       # Recurrent Actor-Critic model (MLP + 2-layer LSTM)

ppo/
  - ppo_agent.py          # PPO optimizer (with clipped objective and entropy bonus)

utils/
  - logger.py             # TensorBoard logger

train.py                  # Main training script
config.py                 # Hyperparameter configuration
```

---

## ⚙️ Requirements

- Python 3.10+
- PyTorch >= 2.0
- TensorBoard
- NumPy

Install dependencies:

```bash
pip install torch numpy tensorboard
```

---

## 🚀 How to Train

Train the reinforcement learning agent for Selfish Mining or Staircase Attack:

```bash
python train.py
```

You can modify `config.py` or directly edit `train.py` to choose between `selfish` and `staircase` attacks:

```python
ENV_TYPE = 'selfish'   # or 'staircase'
```

Training logs are saved in the `runs/` directory.  
To monitor training with TensorBoard:

```bash
tensorboard --logdir=runs/
```

---
