import numpy as np
from hlmdghost import LMDGhost
from env.casper import Casper

class EthereumAttackEnv:
    def __init__(self, num_validators=128, attack_type='selfish'):
        self.attack_type = attack_type
        self.num_validators = num_validators
        self.honest_validators = list(range(0, num_validators-1))
        self.byzantine_validator = num_validators-1
        self.chain = LMDGhost(num_validators)
        self.slot = 0
        self.epoch = 0
        self.block_counter = 0
        if attack_type == 'staircase':
            self.casper = Casper(threshold=0.67)
            self.casper.total_validators = num_validators
        self.reset()

    def reset(self):
        self.slot = 0
        self.epoch = 0
        self.block_counter = 1
        self.chain = LMDGhost(self.num_validators)
        if self.attack_type == 'staircase':
            self.casper.reset()
        return self._get_state()

    def step(self, action):
        if self.attack_type == 'selfish':
            return self._step_selfish(action)
        elif self.attack_type == 'staircase':
            return self._step_staircase(action)

    def _step_selfish(self, action):
        proposer_action, attester_action = action
        self.slot += 1

        proposer_id = self._get_proposer(self.slot)
        parent_id = self.chain.fork_choice()

        if proposer_id == self.byzantine_validator:
            if proposer_action == 0:
                self.chain.add_block(str(self.block_counter), parent_id, proposer_id, withheld=0, self.slot)
            elif proposer_action == 1:
                self.chain.add_block(str(self.block_counter), parent_id, proposer_id, withheld=1, self.slot)
            elif proposer_action == 2:
                for block in self.chain.blocks.values():
                    if block.proposer == proposer_id and block.withheld == 1:
                        block.withheld = 0
        else:
            self.chain.add_block(str(self.block_counter), parent_id, proposer_id, withheld=0, self.slot)

        self.chain.set_proposer_boost(str(self.block_counter))
        self.block_counter += 1

        for v_id in self.honest_validators:
            self.chain.update_vote(v_id, self.chain.fork_choice(), self.slot)

        if attester_action == 0:
            self.chain.update_vote(self.byzantine_validator, self.chain.fork_choice(), self.slot)

        reward = self._compute_reward_selfish()
        next_state = self._get_state()
        done = False
        info = {}

        reward = self.chain.compute_selfish_reward(self.byzantine_validator)

        return next_state, reward, done, info

    def _step_staircase(self, action):
        p1, p2, p3, a1 = action
        self.epoch += 1

        parent_id = self.chain.fork_choice()
        proposer_id = self._get_proposer(self.epoch)

        if proposer_id == self.byzantine_validator:
            withheld = 1 if p1 == 1 else 0
            self.chain.add_block(str(self.block_counter), parent_id, proposer_id, withheld=withheld, epoch=self.epoch)
        else:
            self.chain.add_block(str(self.block_counter), parent_id, proposer_id, withheld=0, epoch=self.epoch)

        self.chain.set_proposer_boost(str(self.block_counter))
        self.block_counter += 1

        # At the beginning of each epoch, add new checkpoint to Casper
        self.casper.add_checkpoint(str(self.block_counter))

        # Honest validators always vote for fork choice
        head = self.chain.fork_choice()
        for v_id in self.honest_validators:
            self.casper.vote(v_id, head)

        # Byzantine behavior based on action
        if a1 == 0:
            self.casper.vote(self.byzantine_validator, head)

        # Update justification and finalization state
        self.casper.update_justification()

        reward = self._compute_reward_staircase()
        next_state = self._get_state()
        done = False
        info = {}

        reward = self.casper.compute_staircase_reward()

        return next_state, reward, done, info

    def _get_proposer(self, index):
        return index % self.num_validators

    def _get_state(self):
        if self.attack_type == 'selfish':
            proposer_duty = np.zeros(8)
            for i in range(8):
                proposer_duty[i] = 1 if self._get_proposer(self.slot + i) == self.byzantine_validator else 0

            head = self.chain.fork_choice()
            pub_weight = self.chain.blocks[head].weight
            priv_weight = sum(b.weight for b in self.chain.blocks.values() if b.proposer == self.byzantine_validator)

            return np.concatenate([proposer_duty, [pub_weight, priv_weight]]).astype(np.float32)

        elif self.attack_type == 'staircase':
            proposer_duty = np.zeros(64)
            for i in range(64):
                proposer_duty[i] = 1 if self._get_proposer(self.epoch + i) == self.byzantine_validator else 0

            j_pub = int(self.casper.get_last_justified() or 0)
            j_priv = max((int(b.epoch) for b in self.chain.blocks.values() if b.proposer == self.byzantine_validator), default=0)
            j_global = min(j_pub, j_priv)

            return np.concatenate([proposer_duty, [j_pub, j_priv, j_global]]).astype(np.float32)

