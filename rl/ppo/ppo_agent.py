import torch
import torch.optim as optim
import torch.nn.functional as F

class PPOAgent:
    def __init__(self, model, lr=3e-4):
        self.model = model
        self.optimizer = optim.Adam(model.parameters(), lr=lr)

    def compute_loss(self, states, actions, advantages, returns, old_log_probs, clip_eps=0.2, entropy_coef=0.01):
        action_logits, state_values, _ = self.model(states, self.model.init_hidden(states.size(0)))
        dist = torch.distributions.Categorical(logits=action_logits)

        new_log_probs = dist.log_prob(actions)
        ratio = (new_log_probs - old_log_probs).exp()

        policy_loss = -torch.min(
            ratio * advantages,
            torch.clamp(ratio, 1-clip_eps, 1+clip_eps) * advantages
        ).mean()

        value_loss = F.mse_loss(state_values.squeeze(-1), returns)
        entropy_loss = dist.entropy().mean()

        return policy_loss + 0.5 * value_loss - entropy_coef * entropy_loss

    def update(self, states, actions, advantages, returns, old_log_probs):
        loss = self.compute_loss(states, actions, advantages, returns, old_log_probs)
        self.optimizer.zero_grad()
        loss.backward()
        torch.nn.utils.clip_grad_norm_(self.model.parameters(), 0.5)
        self.optimizer.step()