from collections import defaultdict

WEIGHT_PER_SLOT = 2

class Block:
    def __init__(self, block_id, parent_id=None, proposer_id=None, withheld=0, slot=0):
        self.id = block_id
        self.parent = parent_id
        self.proposer = proposer_id
        self.weight = 0
        self.withheld = withheld
        self.slot = slot
        self.votes = set()

class Validator:
    def __init__(self, vid, stake=1, slot = 0):
        self.id = vid
        self.stake = stake
        self.latest_vote = None
        self.slot = slot

class LMDGhost:
    def __init__(self, num_validators, genesis_id="genesis"):
        self.blocks = {genesis_id: Block(genesis_id)}
        self.children = defaultdict(list)
        self.validators = {}
        for i in range(num_validators):
            v = Validator(i)
            v.latest_vote = genesis_id
            self.validators[i] = v
        self.boost_block_id = None
        self.proposer_boost_value = 0.4 * WEIGHT_PER_SLOT

    def add_block(self, block_id, parent_id, proposer_id, withheld=0, slot=0):
        b = Block(block_id, parent_id, proposer_id, withheld, slot)
        self.blocks[block_id] = b
        self.children[parent_id].append(block_id)
        ancestry = set()
        current = parent_id
        while current:
            ancestry.add(current)
            current = self.blocks[current].parent
        for block in self.blocks.values():
            for v in block.votes:
                if v not in self.blocks[block_id].votes and block.id not in ancestry:
                    self.blocks[block_id].votes.add(v)

    def update_vote(self, vid, target):
        validator = self.validators.get(vid)
        if validator:
            validator.latest_vote = target
            self.blocks[target].votes.add(vid)
            current = target
            while current:
                self.blocks[current].weight += validator.stake
                current = self.blocks[current].parent

    def fork_choice(self):
        current = "genesis"
        while True:
            children = [self.blocks[child_id] for child_id in self.children[current] if self.blocks[child_id].withheld == 0]
            if not children:
                return current
            current = max(
                children, key=lambda block: (block.weight, block.id)
            ).id

    def set_proposer_boost(self, block_id):
        if block_id in self.blocks:
            current = block_id
            while current:
                self.blocks[current].weight += self.proposer_boost_value
                current = self.blocks[current].parent
            self.boost_block_id = block_id

    def clear_proposer_boost(self):
        if self.boost_block_id:
            current = self.boost_block_id
            while current:
                self.blocks[current].weight -= self.proposer_boost_value
                current = self.blocks[current].parent
            self.boost_block_id = None

    def get_votes(self, block_id):
        return self.blocks[block_id].votes

    def compute_selfish_reward(self, byzantine_id):
        rewards = defaultdict(float)
        # Build chain from genesis to head
        chain = []
        current = self.fork_choice()
        while current != "genesis":
            chain.append(current)
            current = self.blocks[current].parent
        chain = chain[::-1]  # from genesis to head

        # Only consider up to the last 8 blocks
        chain = chain[:-8] if len(chain) > 8 else []

        for block_id in chain:
            block = self.blocks[block_id]
            block_slot = int(block.id)
            for v in block.votes:
                # get the slot of the latest vote of v
                vote_block_id = self.validators[v].latest_vote
                vote_slot = int(self.blocks[vote_block_id].slot)
                if vote_slot + 1 == block_slot:
                    rewards[v] += 1.0
                    rewards[block.proposer] += 1/8
                else:
                    rewards[v] += 40/54
                    rewards[block.proposer] += (1/8)*(40/54)
        return dict(rewards)
