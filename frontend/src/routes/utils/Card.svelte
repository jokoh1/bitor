<!-- File: src/components/Card.svelte -->
<script>
  export let title = '';
</script>

<div class="max-w-full bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700">
  {#if title}
    <div class="p-4 border-b border-gray-200 dark:border-gray-700">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{title}</h3>
    </div>
  {/if}
  <div class="p-4">
    <slot></slot>
  </div>
</div>
