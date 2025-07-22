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

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * AttackerResponse defines the response from the attacker client API. This matches the Go
 * implementation in tsinghua-cel/attacker-client-go/client/type.go.
 */
public class AttackerResponse {
  @JsonProperty("cmd")
  private AttackerCommand cmd;

  @JsonProperty("result")
  private String result;

  public AttackerResponse() {
    // Default constructor for Jackson deserialization
  }

  public AttackerResponse(final AttackerCommand cmd, final String result) {
    this.cmd = cmd;
    this.result = result;
  }

  public AttackerCommand getCmd() {
    return cmd;
  }

  public String getResult() {
    return result;
  }

  /**
   * Convenience method to check if the command is CMD_CONTINUE.
   *
   * @return true if the command is CMD_CONTINUE, false otherwise
   */
  public boolean shouldContinue() {
    return cmd == AttackerCommand.CMD_CONTINUE;
  }

  /**
   * Convenience method to check if the command is CMD_ABORT.
   *
   * @return true if the command is CMD_ABORT, false otherwise
   */
  public boolean shouldAbort() {
    return cmd == AttackerCommand.CMD_ABORT;
  }

  @Override
  public String toString() {
    return "AttackerResponse{" + "cmd=" + cmd + ", result='" + result + '\'' + '}';
  }
}
