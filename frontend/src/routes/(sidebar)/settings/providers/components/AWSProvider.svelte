<script lang="ts">
    import { Label, Input } from 'flowbite-svelte';
    import { pocketbase } from '$lib/stores/pocketbase';
    import type { Provider } from '../types';

    export let provider: Provider;
    export let onSave: (provider: Provider) => void;

    let error = '';

    async function saveSettings() {
        try {
            await $pocketbase.collection('providers').update(provider.id, {
                settings: provider.settings,
                updated: new Date().toISOString()
            });
            onSave(provider);
        } catch (e: any) {
            console.error('Error saving settings:', e);
            error = e.message || 'Failed to save settings';
        }
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="aws_region">Region</Label>
        <Input
            id="aws_region"
            bind:value={provider.settings.region}
            on:blur={saveSettings}
            placeholder="Enter AWS region"
        />
    </div>

    <div>
        <Label for="aws_account_id">Account ID</Label>
        <Input
            id="aws_account_id"
            bind:value={provider.settings.account_id}
            on:blur={saveSettings}
            placeholder="Enter AWS account ID"
        />
    </div>

    {#if error}
        <p class="text-red-500 text-sm">{error}</p>
    {/if}
</div> 