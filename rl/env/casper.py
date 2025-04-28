class Casper:
    def __init__(self, threshold=0.67):
        """
        Args:
            threshold: fraction of validators needed to justify (e.g., 0.67)
        """
        self.threshold = threshold
        self.checkpoints = {}  # checkpoint_id -> vote_count
        self.justified = set()  # justified checkpoint ids
        self.finalized = set()  # finalized checkpoint ids
        self.latest_justified = None
        self.latest_finalized = None
        self.total_validators = 0  # must be set externally

    def add_checkpoint(self, checkpoint_id):
        """Register a new checkpoint block"""
        if checkpoint_id not in self.checkpoints:
            self.checkpoints[checkpoint_id] = 0

    def vote(self, validator_id, checkpoint_id):
        """Simulate a validator voting for a checkpoint"""
        if checkpoint_id in self.checkpoints:
            self.checkpoints[checkpoint_id] += 1

    def update_justification(self):
        """After all votes are cast, update justification and finalization"""
        # Check justification
        for cp_id, vote_count in self.checkpoints.items():
            if cp_id not in self.justified:
                if self.total_validators > 0 and (vote_count / self.total_validators) >= self.threshold:
                    self.justified.add(cp_id)
                    self.latest_justified = cp_id

        # Check finalization: if two consecutive justified checkpoints
        justified_sorted = sorted(self.justified, key=lambda x: int(x))  # assume id can be compared
        if len(justified_sorted) >= 2:
            prev_cp, last_cp = justified_sorted[-2], justified_sorted[-1]
            if prev_cp not in self.finalized:
                self.finalized.add(prev_cp)
                self.latest_finalized = prev_cp

    def get_last_justified(self):
        """Return the latest justified checkpoint id"""
        return self.latest_justified

    def get_last_finalized(self):
        """Return the latest finalized checkpoint id"""
        return self.latest_finalized

    def reset(self):
        """Reset all state (useful between episodes)"""
        self.checkpoints.clear()
        self.justified.clear()
        self.finalized.clear()
        self.latest_justified = None
        self.latest_finalized = None

    def compute_staircase_reward(self):
        finalized = self.casper.get_last_finalized()
        if finalized is None:
            return 1.0
        else:
            return 0.0