import NDK from "@nostr-dev-kit/ndk";
import { writable, type Writable } from "svelte/store";
import { relayUrl, someRelays } from "./consts";

const relays = [relayUrl, ...someRelays];

const Ndk: NDK = new NDK({ explicitRelayUrls: relays });

Ndk.connect().then(() => console.log("ndk connected"));

export const ndk: Writable<NDK> = writable(Ndk);
export const userPublickey: Writable<string | undefined> = writable(undefined);
