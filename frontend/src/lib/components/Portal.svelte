<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    let target: HTMLElement;
    const portal = document.createElement('div');

    onMount(() => {
        document.body.appendChild(portal);
    });

    onDestroy(() => {
        portal.remove();
    });
</script>

<div bind:this={target} class="hidden">
    <slot />
</div>

{#if target}
    <svelte:self>{@html target.innerHTML}</svelte:self>
{/if} 