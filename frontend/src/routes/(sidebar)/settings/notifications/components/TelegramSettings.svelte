<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Input, Label } from 'flowbite-svelte';

    export let settings: {
        token: string;
        chat_ids: string[];
    };

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change', settings);
    }

    function handleChatIdsChange(event: Event) {
        const input = event.target as HTMLInputElement;
        settings.chat_ids = input.value.split(',').map(id => id.trim());
        handleChange();
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="token">Bot Token</Label>
        <Input
            id="token"
            type="password"
            placeholder="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
            bind:value={settings.token}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="chat_ids">Chat IDs (comma-separated)</Label>
        <Input
            id="chat_ids"
            type="text"
            placeholder="-1001234567890, -1009876543210"
            value={settings.chat_ids.join(', ')}
            on:change={handleChatIdsChange}
        />
    </div>
</div> 