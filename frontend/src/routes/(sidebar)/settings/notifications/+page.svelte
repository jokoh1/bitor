<script lang="ts">
    import { Tabs, TabItem } from 'flowbite-svelte';
    import NotificationRules from './components/NotificationRules.svelte';
    import NotificationGroups from './components/NotificationGroups.svelte';
    import { onMount } from 'svelte';
    import { pocketbase } from '$lib/stores/pocketbase';

    interface NotificationRule {
        id: string;
        type: string;
        severity: string[];
        channels: string[];
        enabled: boolean;
        target_type: 'all' | 'user' | 'group';
        target_users?: string[];
        target_groups?: string[];
    }

    interface NotificationSettings {
        id: string;
        providers: string[];
        rules: NotificationRule[];
    }

    let settings: NotificationSettings | null = null;
    let rules: NotificationRule[] = [];

    onMount(async () => {
        try {
            // Fetch settings
            const settingsRecords = await $pocketbase.collection('notification_settings').getFullList();
            if (settingsRecords.length > 0) {
                const settingsRecord = settingsRecords[0];
                settings = {
                    id: settingsRecord.id,
                    providers: settingsRecord.providers || [],
                    rules: settingsRecord.rules ? JSON.parse(settingsRecord.rules) : []
                };
                rules = settings.rules;
            }
        } catch (error) {
            console.error('Error initializing notification settings:', error);
        }
    });
</script>

<div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Notification Settings</h1>

    <Tabs style="underline">
        <TabItem open title="Rules">
            <NotificationRules {rules} />
        </TabItem>
        <TabItem title="Groups">
            <NotificationGroups />
        </TabItem>
    </Tabs>
</div>