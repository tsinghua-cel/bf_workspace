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

import com.googlecode.jsonrpc4j.JsonRpcHttpClient;
import tech.pegasys.teku.infrastructure.async.SafeFuture;

public class AttackService {

  private static final String BLOCK_MODULE = "block";
  private static final String ATTEST_MODULE = "attest";

  private final JsonRpcHttpClient jsonRpcClient;

  public AttackService() {
    this.jsonRpcClient = JsonRpcClientUtil.getInstance().getJsonRpcClient();
  }

  public AttackService(final String url) {
    this.jsonRpcClient = JsonRpcClientUtil.getInstance().getJsonRpcClientByUrl(url);
  }

  public boolean enabled() {
    return this.jsonRpcClient != null;
  }

  private SafeFuture<AttackerResponse> invokeRpc(final String method, final Object[] params) {
    try {
      AttackerResponse response = jsonRpcClient.invoke(method, params, AttackerResponse.class);
      return SafeFuture.completedFuture(response);
    } catch (Exception e) {
      return SafeFuture.failedFuture(e);
    } catch (Throwable e) {
      throw new RuntimeException(e);
    }
  }

  public SafeFuture<AttackerResponse> blockGetNewParentRoot(
      final long slot, final String pubkey, final String parentRoot) {
    return invokeRpc(BLOCK_MODULE + "_getNewParentRoot", new Object[] {slot, pubkey, parentRoot});
  }

  public SafeFuture<AttackerResponse> delayForReceiveBlock(final long slot) {
    return invokeRpc(BLOCK_MODULE + "_delayForReceiveBlock", new Object[] {slot});
  }

  public SafeFuture<AttackerResponse> blockBeforeBroadcast(final long slot) {
    return invokeRpc(BLOCK_MODULE + "_beforeBroadCast", new Object[] {slot});
  }

  public SafeFuture<AttackerResponse> blockAfterBroadcast(final long slot) {
    return invokeRpc(BLOCK_MODULE + "_afterBroadCast", new Object[] {slot});
  }

  public SafeFuture<AttackerResponse> blockBeforeSign(
      final long slot, final String pubkey, final String blockDataBase64) {
    return invokeRpc(BLOCK_MODULE + "_beforeSign", new Object[] {slot, pubkey, blockDataBase64});
  }

  public SafeFuture<AttackerResponse> blockAfterSign(
      final long slot, final String pubkey, final String signedBlockDataBase64) {
    return invokeRpc(
        BLOCK_MODULE + "_afterSign", new Object[] {slot, pubkey, signedBlockDataBase64});
  }

  public SafeFuture<AttackerResponse> blockBeforePropose(
      final long slot, final String pubkey, final String signedBlockDataBase64) {
    return invokeRpc(
        BLOCK_MODULE + "_beforePropose", new Object[] {slot, pubkey, signedBlockDataBase64});
  }

  public SafeFuture<AttackerResponse> blockAfterPropose(
      final long slot, final String pubkey, final String signedBlockDataBase64) {
    return invokeRpc(
        BLOCK_MODULE + "_afterPropose", new Object[] {slot, pubkey, signedBlockDataBase64});
  }

  public SafeFuture<AttackerResponse> attestBeforeBroadcast(final long slot) {
    return invokeRpc(ATTEST_MODULE + "_beforeBroadCast", new Object[] {slot});
  }

  public SafeFuture<AttackerResponse> attestAfterBroadcast(final long slot) {
    return invokeRpc(ATTEST_MODULE + "_afterBroadCast", new Object[] {slot});
  }

  public SafeFuture<AttackerResponse> attestBeforeSign(
      final long slot, final String pubkey, final String attestDataBase64) {
    return invokeRpc(ATTEST_MODULE + "_beforeSign", new Object[] {slot, pubkey, attestDataBase64});
  }

  public SafeFuture<AttackerResponse> attestAfterSign(
      final long slot, final String pubkey, final String signedAttestDataBase64) {
    return invokeRpc(
        ATTEST_MODULE + "_afterSign", new Object[] {slot, pubkey, signedAttestDataBase64});
  }

  public SafeFuture<AttackerResponse> attestBeforePropose(
      final long slot, final String pubkey, final String signedAttestDataBase64) {
    return invokeRpc(
        ATTEST_MODULE + "_beforePropose", new Object[] {slot, pubkey, signedAttestDataBase64});
  }

  public SafeFuture<AttackerResponse> attestAfterPropose(
      final long slot, final String pubkey, final String signedAttestDataBase64) {
    return invokeRpc(
        ATTEST_MODULE + "_afterPropose", new Object[] {slot, pubkey, signedAttestDataBase64});
  }
}
