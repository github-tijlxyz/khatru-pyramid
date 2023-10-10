import NDK, { NDKNip07Signer } from '@nostr-dev-kit/ndk';
import '../style.css';
import { loadUsernames } from './user.js';
import { loadActions } from './actions.js';
// import NDKCacheAdapterDexie from "@nostr-dev-kit/ndk-cache-dexie";

export let ndk: NDK;
export let relayUrl: string;

async function load() {
  // get relay websocket url
  var loc = window.location;
  if (loc.protocol === "https:") {
    relayUrl = "wss:";
  } else {
    relayUrl = "ws:";
  }
  relayUrl += "//" + loc.host;

  // init NDK
  // const dexieAdapter = new NDKCacheAdapterDexie({ dbName: 'khatru-invite-ndk-cache' })
  const nip07signer = new NDKNip07Signer()
  ndk = new NDK({
    explicitRelayUrls: [
      relayUrl,
      "wss://nos.lol",
      "wss://nostr-pub.wellorder.net",
    ], signer: nip07signer
  });
  await ndk.connect();

  loadActions();
  loadUsernames();
}

// load
load();




