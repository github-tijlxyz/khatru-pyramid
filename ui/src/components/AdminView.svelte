<script>
  import { NDKRelay, NDKRelaySet } from "@nostr-dev-kit/ndk";
  import { ndk } from "../lib/nostr";
  import { onMount } from "svelte";
  import ReportEvent from "./ReportEvent.svelte";
  import { relayUrl } from "../lib/consts";

  let events = [];

  onMount(async () => {
    try {
      // dont know why this needs so much code
      let specificRelay = [new NDKRelay(relayUrl)];
      const relaySet = new NDKRelaySet(specificRelay, $ndk);
      relaySet.relays.forEach(async (relay) => {
        await relay.connect().catch((err) => {
          console.log("error while connecting to relay", err);
        });
        relay.on("connect", () => {
          console.log("connected");
        });
      });

      let filter = { kinds: [1984], limit: 250 };
      let options = { closeOnEose: true };
      let es = await $ndk.fetchEvents(filter, options, relaySet);
      events = Array.from(es);
    } catch (error) {
      console.log("error while getting feed", error);
    }
  });
</script>

<div>
  {#each events as event}
    <ReportEvent {event} />
  {/each}
  {#if events.length == 0}
    <span>Didn't found any events</span>
  {/if}
</div>
