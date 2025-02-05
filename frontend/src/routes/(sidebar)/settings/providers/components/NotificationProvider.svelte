<script lang="ts">
    import { Label, Input, Select } from 'flowbite-svelte';
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
    {#if provider.provider_type === 'email'}
        <div>
            <Label for="smtp_host">SMTP Host</Label>
            <Input
                id="smtp_host"
                bind:value={provider.settings.smtp_host}
                on:blur={saveSettings}
                placeholder="Enter SMTP host"
            />
        </div>

        <div>
            <Label for="smtp_port">SMTP Port</Label>
            <Input
                id="smtp_port"
                type="number"
                bind:value={provider.settings.smtp_port}
                on:blur={saveSettings}
                placeholder="Enter SMTP port"
            />
        </div>

        <div>
            <Label for="from_address">From Address</Label>
            <Input
                id="from_address"
                type="email"
                bind:value={provider.settings.from_address}
                on:blur={saveSettings}
                placeholder="Enter from address"
            />
        </div>

        <div>
            <Label for="encryption">Encryption</Label>
            <Select
                id="encryption"
                bind:value={provider.settings.encryption}
                on:change={saveSettings}
            >
                <option value="none">None</option>
                <option value="tls">TLS</option>
                <option value="starttls">STARTTLS</option>
            </Select>
        </div>
    {:else if ['slack', 'teams', 'discord'].includes(provider.provider_type)}
        <div>
            <Label for="webhook_url">Webhook URL</Label>
            <Input
                id="webhook_url"
                bind:value={provider.settings.webhook_url}
                on:blur={saveSettings}
                placeholder="Enter webhook URL"
            />
        </div>
    {:else if provider.provider_type === 'telegram'}
        <div>
            <Label for="bot_token">Bot Token</Label>
            <Input
                id="bot_token"
                bind:value={provider.settings.bot_token}
                on:blur={saveSettings}
                placeholder="Enter bot token"
            />
        </div>

        <div>
            <Label for="chat_id">Chat ID</Label>
            <Input
                id="chat_id"
                bind:value={provider.settings.chat_id}
                on:blur={saveSettings}
                placeholder="Enter chat ID"
            />
        </div>
    {:else if provider.provider_type === 'jira'}
        <div>
            <Label for="jira_url">Jira URL</Label>
            <Input
                id="jira_url"
                bind:value={provider.settings.jira_url}
                on:blur={saveSettings}
                placeholder="Enter Jira URL"
            />
        </div>

        <div>
            <Label for="username">Username</Label>
            <Input
                id="username"
                bind:value={provider.settings.username}
                on:blur={saveSettings}
                placeholder="Enter username"
            />
        </div>

        <div>
            <Label for="project">Project</Label>
            <Input
                id="project"
                bind:value={provider.settings.project}
                on:blur={saveSettings}
                placeholder="Enter project"
            />
        </div>
    {/if}

    {#if error}
        <p class="text-red-500 text-sm">{error}</p>
    {/if}
</div> 