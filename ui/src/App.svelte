<script lang="ts">
  import "./App.css";
  import Hieracy from "./components/Hierarchy.svelte";
  import Invite from "./components/Invite.svelte";
  import { onMount } from "svelte";
  import { ndk, userPublickey } from "./lib/nostr";
  import { NDKNip07Signer } from "@nostr-dev-kit/ndk";
  import { buildHierarchy } from "./lib/utils";
  import { nip19 } from "nostr-tools";
  import AdminView from "./components/AdminView.svelte";

  let adminView = false;

  async function login() {
    const signer = new NDKNip07Signer();
    $ndk.signer = signer;
    ndk.set($ndk);
    $userPublickey = (await $ndk.signer.user()).npub;
    userPublickey.set($userPublickey);
  }

  async function fetchData() {
    const response = await fetch("/invitedata");
    invitedata = Object.values(await response.json());
    hierarchy = buildHierarchy(invitedata, { pk: "", invited_by: "" });
  }

  let invitedata = [];
  let hierarchy = [];

  onMount(() => {
    addEventListener("load", (e) => {
      login();
    });
    fetchData();
  });
</script>

<article class="font-sans px-4 py-6 lg:max-w-7xl lg:pt-6 lg:pb-28">
  <h1>Invite Relay</h1>
  {#if adminView == true}
  <button
  on:click={() => adminView = false}
  type="button"
  class="inline-flex mr-2 items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
  >Leave Admin View</button
>
  <AdminView />
  {:else if adminView == false}
  <div>
    {#if $userPublickey === undefined}
      <button
        on:click={login}
        type="button"
        class="inline-flex mr-2 items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
        >Login with NIP07</button
      >
    {/if}
  </div>
  <div>
    <button
    on:click={() => adminView = true}
    type="button"
    class="inline-flex mr-2 items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
    >Open Admin View</button
  >
  </div>
  {#if invitedata.find((p) => $userPublickey == nip19.npubEncode(p.pk))}
    <div>
      <h3>Invite Someone</h3>
      <div>
        <Invite reload={fetchData} />
      </div>
    </div>
  {/if}
  <div>
    <h3>Current hieracy</h3>
    <div>
      <Hieracy {hierarchy} reload={fetchData} />
    </div>
  </div>
  {/if}
</article>
