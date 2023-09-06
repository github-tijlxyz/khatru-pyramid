<script lang="js">
  import { onMount } from "svelte";
  import { ndk, userPublickey } from "../lib/nostr";
  import TreeNode from "./TreeNode.svelte";
  import { nip19 } from "nostr-tools";
  import { NDKEvent, NDKRelay, NDKRelaySet, NDKUser } from "@nostr-dev-kit/ndk";
  import { relayUrl } from "../lib/consts";
  export let person;
  export let reload;
  let username = "...";

  onMount(async () => {
    try {
      let user = $ndk.getUser({
        hexpubkey: person.pk,
      });
      await user.fetchProfile();
      username = user.profile.name;
    } catch (error) {
      console.log(error);
    }
  });

  async function removeThisUser() {
    let confirmation = confirm(
      `Are you sure you want to remove ${
        username ? username : person.pk
      }? All people they invited will also be removed.`,
    );
    if (confirmation) {
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
            console.log("connected");
          });
        });

        const event = new NDKEvent($ndk);
        event.kind = 20202;
        event.tags.push(["p", person.pk]);
        await event.publish(relaySet).then(() => reload());
      } catch (error) {
        console.log("error while publishing", error);
      }
    }
  }
</script>

<li>
  <span>
    <a
      class="inline hover:underline"
      href={`nostr:${nip19.npubEncode(person.pk)}`}
      >{#if username}{username}{:else}{nip19.npubEncode(person.pk)}{/if}</a
    >
    {#if $userPublickey == nip19.npubEncode(person.invited_by)}<button
        on:click={removeThisUser}
        class="inline cusor-pointer font-semibold text-red-500">[-]</button
      >{/if}{#if $userPublickey == nip19.npubEncode(person.pk)}
      <span> (you)</span>
    {/if}
  </span>
  {#if person.children && person.children.length > 0}
    <ul class="list-disc ml-2">
      {#each person.children as child (child.pk)}
        <TreeNode {reload} person={child} />
      {/each}
    </ul>
  {/if}
</li>
