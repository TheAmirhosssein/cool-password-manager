// opaqueClient.js
import { OpaqueClient, getOpaqueConfig, RegistrationResponse, OPAQUE_P256 } from "@cloudflare/opaque-ts";

const encoder = new TextEncoder();

export class OpaqueClientWrapper {
    constructor(serverIdentity) {
        this.cfg = getOpaqueConfig(3);
        this.client = new OpaqueClient(this.cfg);
        this.state = null;
        this.serverIdentity = serverIdentity
    }

    async registerInit(password) {
        const registrationRequest = await this.client.registerInit(password)
        return registrationRequest.serialize()
    }

    async registerFinish(serverResponseBytes, clientIdentity) {
        const deserRes = RegistrationResponse.deserialize(this.cfg, Array.from(serverResponseBytes))

        const rec = await this.client.registerFinish(deserRes, this.serverIdentity, clientIdentity)

        const { record, _ } = rec
        return record.serialize()
    }

    async loginInit(password) {
        const { credentialRequest, state } =
            await this.client.loginInit(encoder.encode(password));
        this.state = state;

        return credentialRequest.serialize();
    }

    async loginFinish(serverResponseBytes) {
        const credentialResponse =
            this.client.deserializeCredentialResponse(serverResponseBytes);

        const result = await this.client.loginFinish(this.state, credentialResponse);

        return {
            sessionKey: result.sessionKey,
            exportKey: result.exportKey,
        };
    }
}
