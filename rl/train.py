import torch
from env.environment import EthereumAttackEnv
from model.actor_critic import RecurrentActorCritic
from ppo.ppo_agent import PPOAgent
from utils.logger import Logger


ENV_TYPE = 'selfish'
TOTAL_TIMESTEPS = 100_000
SAVE_PATH = './model_ckpt.pth'
LOG_DIR = './runs/'


env = EthereumAttackEnv(attack_type=ENV_TYPE)
state_dim = len(env.reset())
action_dim = 3 if ENV_TYPE == 'selfish' else 4
model = RecurrentActorCritic(state_dim, action_dim)
agent = PPOAgent(model)
logger = Logger(LOG_DIR)


obs = torch.tensor(env.reset()).float().unsqueeze(0)
hidden = model.init_hidden()
timestep = 0

while timestep < TOTAL_TIMESTEPS:
    with torch.no_grad():
        action_logits, _, hidden = model(obs, hidden)
        action_dist = torch.distributions.Categorical(logits=action_logits)
        action = action_dist.sample()
        log_prob = action_dist.log_prob(action)

    next_obs, reward, done, _ = env.step(action.item())
    next_obs = torch.tensor(next_obs).float().unsqueeze(0)

    advantage = torch.tensor([reward])
    returns = advantage

    agent.update(obs, action, advantage, returns, log_prob)

    logger.log('reward', reward, timestep)

    obs = next_obs
    timestep += 1

    if timestep % 5000 == 0:
        torch.save(model.state_dict(), SAVE_PATH)
        print(f"Saved model at {SAVE_PATH}")

logger.close()
print("Training Done!")