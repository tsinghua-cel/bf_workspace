/*
 * Copyright Consensys Software Inc., 2025
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package tech.pegasys.teku.attacker.client;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

import org.junit.jupiter.api.Test;

public class AttackServiceTest {

  @Test
  public void attackIsNotEnabled() {
    AttackService s = new AttackService();
    assertThat(s.enabled()).isEqualTo(false);
  }

  //  @Test
  //  public void attackIsEnabled() {
  //    AttackService s = new AttackService("http://127.0.0.1:12000");
  //    assertThat(s.enabled()).isEqualTo(true);
  //  }

  //  @Test
  //  public void blockGetNewParentRoot() throws ExecutionException, InterruptedException {
  //    AttackService s = new AttackService("http://127.0.0.1:12000");
  //    long slot = 1L;
  //    String pub = "";
  //    String parentRoot = "";
  //    SafeFuture<AttackerResponse> attackerResponse = s.blockGetNewParentRoot(slot, pub,
  // parentRoot);
  //    assertThat(attackerResponse.get().getCmd().getValue())
  //        .isEqualTo(AttackerCommand.CMD_NULL.getValue());
  //  }
  //    public void blockBeforeBroadcast() {
  //        AttackService s = new AttackService("http://127.0.0.1:12000");
  //        long slot = 1L;
  //        boolean skipBroadCast = false;
  //
  //        try {
  //            AttackerResponse attackerResponse = s.blockBeforeBroadcast(slot).get(); // Blocking
  // call to get the result
  //            switch (attackerResponse.getCmd()) {
  //                case CMD_EXIT:
  //                case CMD_ABORT:
  //                    System.exit(-1); // Terminate the process
  //                    break;
  //                case CMD_SKIP:
  //                    skipBroadCast = true; // Skip broadcast
  //                    break;
  //                case CMD_RETURN:
  //                    // Simulate returning a response (adjust as per actual method requirements)
  //                    return; // Exit the method
  //                case CMD_NULL:
  //                case CMD_CONTINUE:
  //                    // Do nothing
  //                    break;
  //                default:
  //                    throw new IllegalStateException("Unexpected command received: " +
  // attackerResponse.getCmd());
  //            }
  //        } catch (Exception e) {
  //            e.printStackTrace();
  //            throw new RuntimeException("Failed to process attacker response", e);
  //        }
  //
  //        System.out.println("Block before broadcast completed for slot: " + slot + ", skip:" +
  // skipBroadCast);
  //        // then sleep 10s
  //    }
}
