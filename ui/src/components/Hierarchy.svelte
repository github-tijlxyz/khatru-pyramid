<script lang="js">
  import TreeNode from "./TreeNode.svelte";

  export let reload;
  export let hierarchy;
</script>

<div class="px-4">
  <ul class="list-disc ml-4">
    {#each hierarchy as person (person.pk)}
      <TreeNode {reload} {person} />
    {/each}
  </ul>
</div>
