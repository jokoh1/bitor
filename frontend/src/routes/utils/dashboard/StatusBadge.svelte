<script lang="ts">
	import { Badge, Spinner } from 'flowbite-svelte';
	import type { BadgeProps } from 'flowbite-svelte/Badge.svelte';
	import { onDestroy } from 'svelte';

	export const dark: boolean = false;
	export let state: string;
	export let destroyed: boolean = false;
	export let end_time: string | undefined = undefined;

	const states = {
		failed: 'Failed',
		started: 'Started',
		generating: 'Generating',
		deploying: 'Deploying',
		running: 'Running',
		finished: 'Finished',
		created: 'Created',
		stopped: 'Stopped',
		manual: 'Manual'
	} as const;

	const colors: Record<string, BadgeProps['color']> = {
		failed: 'red',
		started: 'blue',
		generating: 'yellow',
		deploying: 'blue',
		running: 'indigo',
		finished: 'green',
		created: 'dark',
		stopped: 'yellow',
		manual: 'purple'
	};

	let progress = 0;
	let interval: ReturnType<typeof setInterval> | null = null;

	// Reactive statement to update normalizedState whenever 'state' changes
	$: normalizedState = state.toLowerCase() as keyof typeof states;

	// Reactive statement to manage the interval based on normalizedState
	$: {
		if (['started', 'generating', 'deploying', 'running'].includes(normalizedState)) {
			if (!interval) {
				interval = setInterval(() => {
					progress = (progress + 10) % 100; // Loop progress from 0 to 100
				}, 1000); // Update every second
			}
		} else {
			if (interval) {
				clearInterval(interval);
				interval = null;
			}
		}
	}

	onDestroy(() => {
		if (interval) {
			clearInterval(interval);
			interval = null;
		}
	});
</script>

<div class="flex items-center gap-2">
	<Badge color={colors[normalizedState] || 'dark'} class="flex items-center">
		{#if ['started', 'generating', 'deploying', 'running'].includes(normalizedState)}
			<Spinner class="mr-2" style="color: {colors[normalizedState] || 'dark'}" />
		{/if}
		{states[normalizedState] || 'Unknown Status'}
	</Badge>
	{#if destroyed}
		<Badge color="dark">Destroyed</Badge>
	{:else if ['started', 'generating', 'deploying', 'running'].includes(normalizedState)}
		<Badge color="green" class="animate-pulse">Online</Badge>
	{:else if normalizedState === 'failed' && !end_time}
		<Badge color="red" class="animate-pulse">Failed but Running</Badge>
	{/if}
</div>
