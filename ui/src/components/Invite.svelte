<script>
  import { NDKEvent, NDKRelay, NDKRelaySet } from "@nostr-dev-kit/ndk";
  import { ndk } from "../lib/nostr";
  import { nip19 } from "nostr-tools";
  import { relayUrl } from "../lib/consts";

  let pubKeyToInvite = "";
  export let reload;

  async function invite() {
    if (!pubKeyToInvite) return;
    try {
      // only publish to the relay in question, dont know why this needs so much code
      let specificRelay = [new NDKRelay(relayUrl)];
      const relaySet = new NDKRelaySet(specificRelay, $ndk);
      relaySet.relays.forEach(async (relay) => {
        await relay.connect().catch((err) => {
          console.log("RELAY CONNECT ERROR");
          console.error(err);
        });
        relay.on("connect", () => {
          console.log(relay.url, "connected");
        });
      });

      const pk = pubKeyToInvite.startsWith("npub")
        ? nip19.decode(pubKeyToInvite).data
        : pubKeyToInvite;
      const event = new NDKEvent($ndk);
      event.kind = 20201;
      event.tags.push(["p", pk.toString()]);
      await event.publish(relaySet).then(() => reload());
      pubKeyToInvite = "";
    } catch (error) {
      console.log("error while publishing", error);
    }
  }
</script>

<div class="relative flex items-stretch flex-grow focus-within:z-10">
  <input
    bind:value={pubKeyToInvite}
    class="focus:ring-indigo-500 focus:border-indigo-500 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300"
    placeholder="hex or npub"
  />
  <button
    on:click={invite}
    type="submit"
    class="-ml-px relative inline-flex items-center space-x-2 px-3 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 text-white"
    >Go</button
  >
</div>
