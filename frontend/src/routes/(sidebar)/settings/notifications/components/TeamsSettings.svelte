<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Input, Label } from 'flowbite-svelte';

    export let settings: {
        webhook_url: string;
    };

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change', settings);
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="webhook_url">Webhook URL</Label>
        <Input
            id="webhook_url"
            type="url"
            placeholder="https://outlook.office.com/webhook/..."
            bind:value={settings.webhook_url}
            on:change={handleChange}
        />
    </div>
</div> 