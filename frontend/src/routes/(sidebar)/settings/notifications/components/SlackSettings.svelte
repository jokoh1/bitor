<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Input, Label } from 'flowbite-svelte';

    export let settings: {
        token: string;
        channel: string;
    };

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change', settings);
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="token">Bot Token</Label>
        <Input
            id="token"
            type="password"
            placeholder="xoxb-..."
            bind:value={settings.token}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="channel">Default Channel</Label>
        <Input
            id="channel"
            type="text"
            placeholder="#general"
            bind:value={settings.channel}
            on:change={handleChange}
        />
    </div>
</div> 