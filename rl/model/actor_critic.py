import torch
import torch.nn as nn

class RecurrentActorCritic(nn.Module):
    def __init__(self, state_dim, action_dim):
        super().__init__()
        self.encoder = nn.Sequential(
            nn.Linear(state_dim, 512),
            nn.ReLU()
        )
        self.lstm = nn.LSTM(512, 256, num_layers=2, batch_first=True)
        self.actor_head = nn.Linear(256, action_dim)
        self.critic_head = nn.Linear(256, 1)

    def forward(self, state, hidden):
        x = self.encoder(state)
        x, hidden = self.lstm(x.unsqueeze(1), hidden)
        x = x.squeeze(1)
        action_logits = self.actor_head(x)
        state_value = self.critic_head(x)
        return action_logits, state_value, hidden

    def init_hidden(self, batch_size=1):
        return (torch.zeros(2, batch_size, 256), torch.zeros(2, batch_size, 256))