<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Input, Label } from 'flowbite-svelte';

    export let settings: {
        webhook_id: string;
        token: string;
    };

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change', settings);
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="webhook_id">Webhook ID</Label>
        <Input
            id="webhook_id"
            type="text"
            placeholder="Discord webhook ID"
            bind:value={settings.webhook_id}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="webhook_token">Webhook Token</Label>
        <Input
            id="webhook_token"
            type="password"
            placeholder="Discord webhook token"
            bind:value={settings.token}
            on:change={handleChange}
        />
    </div>
</div> 