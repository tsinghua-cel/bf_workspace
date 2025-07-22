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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonValue;

/**
 * AttackerCommand defines the command types that can be returned by the attacker client API. This
 * matches the Go implementation in tsinghua-cel/attacker-client-go/client/type.go.
 */
public enum AttackerCommand {
  CMD_NULL(0),
  CMD_CONTINUE(1),
  CMD_RETURN(2),
  CMD_ABORT(3),
  CMD_SKIP(4),
  CMD_ROLE_TO_NORMAL(5), // Role changes to normal node
  CMD_ROLE_TO_ATTACKER(6), // Role changes to attacker
  CMD_EXIT(7),
  CMD_UPDATE_STATE(8),
  CMP_DELAY_BROAD_CAST(9);

  private final int value;

  AttackerCommand(final int value) {
    this.value = value;
  }

  @JsonValue
  public int getValue() {
    return value;
  }

  @JsonCreator
  public static AttackerCommand fromValue(final int value) {
    for (AttackerCommand command : AttackerCommand.values()) {
      if (command.value == value) {
        return command;
      }
    }
    throw new IllegalArgumentException("Unknown AttackerCommand value: " + value);
  }
}
