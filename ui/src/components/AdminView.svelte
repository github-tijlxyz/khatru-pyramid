<script lang="js">
  import {
    NDKRelay,
    NDKRelaySet,
  } from "@nostr-dev-kit/ndk";
  import { ndk } from "../lib/nostr";
  import { onMount } from "svelte";
  import ReportEvent from "./ReportEvent.svelte";
  import { relayUrl } from '../lib/consts';

  let events = [];

  onMount(async () => {
    try {
      // dont know why this needs so much code
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

      let filter = { kinds: [1984], limit: 250 };
      let es = await $ndk.fetchEvents(filter, relaySet);
      events = Array.from(es);
    } catch {
      console.log("error while getting feed", error);
    }
  });
</script>

<div>
  {#each events as event}
  {#if event.relay.url == relayUrl}
    <ReportEvent {event} />
  {/if}
  {/each}
</div>
