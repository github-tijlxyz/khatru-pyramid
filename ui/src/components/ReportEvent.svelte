<script>
  import { NDKEvent, NDKRelay, NDKRelaySet } from "@nostr-dev-kit/ndk";
  import { relayUrl } from "../lib/consts";
  import { ndk } from "../lib/nostr";
  export let event;

  let show = true;

  async function dismissReportEvent() {
    let confirmation = confirm(
      "Are you sure you want to delete this event? (You can only do this if you are a relay master)",
    );
    if (confirmation) {
      try {
        // remove report event
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

        const newevent = new NDKEvent($ndk);
        newevent.kind = 20203;
        newevent.tags.push(["e", event.id]);
        await newevent.publish(relaySet);
        show = false;
      } catch (error) {
        console.log("error publishing event", error);
      }
    }
  }

  async function removeUser(username, pk) {
    let confirmation = confirm(
      `Are you sure you want to remove ${
        username ? username : pk
      }? All people they invited will also be removed. (you can only do this if you invited this user or are the relay admin)`,
    );
    if (confirmation) {
      try {
        // only publish to the relay in question, dont know why this needs so much code
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

        const newevent = new NDKEvent($ndk);
        newevent.kind = 20202;
        newevent.tags.push(["p", pk]);
        await newevent.publish(relaySet);
      } catch (error) {
        console.log("error while publishing", error);
      }
    }
  }
</script>

{#if show}
  <div class="rounded-lg border border-slate-500 bg-slate-50 w-full p-4 mt-8">
    <div class="columns-2 p-0 m-0">
      <div class="bg-white rounded-lg px-4 py-2">
        from
        {#await event.author?.fetchProfile()}
          <a href={`nostr:${event.author.npub}`}>...</a>
        {:then profile}
          <img
            class="h-7 w-7 m-0 p-0 rounded-full inline"
            src={profile &&
              JSON.parse(Array.from(profile)[0]?.content)?.picture}
            alt=""
          />
          <a class="hover:underline" href={`nostr:${event.author.npub}`}
            >{profile && JSON.parse(Array.from(profile)[0]?.content)?.name}</a
          >
        {/await}
      </div>
      <div class="bg-white rounded-lg px-4 py-2">
        {#if event?.tags.find((e) => e[0] == "p")?.[0] && event?.tags.find((e) => e[0] == "p")?.[1]}
          to
          {#await $ndk
            .getUser({ hexpubkey: event.tags.find((e) => e[0] == "p")?.[1] })
            .fetchProfile()}
            <img class="h-7 w-7 m-0 p-0 rounded-full inline" src="" alt="" />
            <a href={`nostr:${event.tags.find((e) => e[0] == "p")?.[1]}`}>...</a
            >
          {:then profile}
            <img
              class="h-7 w-7 m-0 p-0 rounded-full inline"
              src={profile &&
                JSON.parse(Array.from(profile)[0]?.content)?.picture}
              alt=""
            />
            <a
              class="hover:underline"
              href={`nostr:${event.tags.find((e) => e[0] == "p")?.[1]}`}
              >{profile &&
              JSON.parse(Array.from(profile)[0]?.content)?.name.length <= 16
                ? JSON.parse(Array.from(profile)[0]?.content)?.name
                : JSON.parse(Array.from(profile)[0]?.content)?.name.substring(
                    0,
                    13,
                  ) + "..."}</a
            >
          {/await}
        {:else}
          huh, nothing here
        {/if}
      </div>
    </div>
    <div class="bg-white max-h-64 scroll-auto rounded-lg m-0 p-4 mt-2">
      {#if event?.tags.find((e) => e[0] == "e")?.[0] && event?.tags.find((e) => e[0] == "e")?.[1] && event?.tags.find((e) => e[0] == "e")?.[2]}
        {event?.tags.find((e) => e[0] == "e")?.[2]}
      {:else if event?.tags.find((e) => e[0] == "p")?.[0] && event?.tags.find((e) => e[0] == "p")?.[1] && event?.tags.find((e) => e[0] == "p")?.[2]}
        {event?.tags.find((e) => e[0] == "p")?.[2]}
      {/if}
      {#if event?.content}
        <p class="text-gray-500 m-0 p-0">{event.content}</p>
      {/if}
    </div>
    {#if event?.tags.find((e) => e[0] == "e")?.[0] && event?.tags.find((e) => e[0] == "e")?.[1]}
      <div
        class="bg-white max-h-64 max-w-full overflow-y-scroll rounded-lg m-0 p-4 mt-2"
      >
        {#await $ndk.fetchEvent( { ids: [event?.tags.find((e) => e[0] == "e")?.[1]] }, )}
          <p class="text-gray-500 p-0 m-0">... (can we find this event?)</p>
        {:then theevent}
          <p class="m-0 p-0">{theevent.content}</p>
          <p class="text-gray-500 m-0 p-0">
            kind: <span class="text-black inline">{theevent.kind}</span>
          </p>
          <p class="text-gray-500 m-0 p-0">
            tags: {#each theevent.tags as tag}<p>
                {tag[0]}: <span class="text-black">{tag[1]}</span>
              </p>{/each}
          </p>
        {/await}
      </div>
    {/if}
    <div class="p-0 mb-0 mx-0 mt-2">
      <div class="columns-4 inline">
        <a
          href={`nostr:${
            event?.tags.find((e) => e[0] == "e")?.[0] &&
            event?.tags.find((e) => e[0] == "e")?.[1]
              ? event?.tags.find((e) => e[0] == "e")?.[1]
              : event?.tags.find((e) => e[0] == "p")?.[1]
          }`}
          class="rounded-lg inline bg-slate-100 p-2 cursor-pointer hover:bg-white"
        >
          Open in client
        </a>
        <button
          on:click={dismissReportEvent}
          class="rounded-lg inline bg-green-500 p-2 cursor-pointer hover:bg-green-400"
        >
          Remove This Report Event
        </button>
        <button
          on:click={() =>
            removeUser("user", event.tags.find((e) => e[0] == "p")?.[1])}
          class="rounded-lg inline bg-red-500 p-2 cursor-pointer hover:bg-red-400"
        >
          Exlude reported user
        </button>
      </div>
    </div>
  </div>
{/if}
