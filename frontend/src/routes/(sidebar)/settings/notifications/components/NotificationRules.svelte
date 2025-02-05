<script lang="ts">
    import { Button, Card, Toggle, Select, MultiSelect } from 'flowbite-svelte';
    import { pocketbase } from '$lib/stores/pocketbase';
    import { onMount } from 'svelte';

    export let rules: Array<{
        id: string;
        type: string;
        severity: string[];
        channels: string[];
        enabled: boolean;
        target_type: 'all' | 'user' | 'group';
        target_users?: string[];
        target_groups?: string[];
    }> = [];

    let users: Array<{ value: string; name: string }> = [];
    let groups: Array<{ value: string; name: string }> = [];
    let providers: Array<{ value: string; name: string }> = [];

    const notificationTypes = [
        { value: 'scan_finished', name: 'Scan Finished' },
        { value: 'scan_started', name: 'Scan Started' },
        { value: 'high_finding', name: 'High Severity Finding' },
        { value: 'critical_finding', name: 'Critical Finding' }
    ];

    const severityLevels = [
        { value: 'low', name: 'Low' },
        { value: 'medium', name: 'Medium' },
        { value: 'high', name: 'High' },
        { value: 'critical', name: 'Critical' }
    ];

    const targetTypes = [
        { value: 'all', name: 'All Users' },
        { value: 'user', name: 'Specific Users' },
        { value: 'group', name: 'User Groups' }
    ];

    async function loadUsersAndGroups() {
        try {
            // Fetch users
            const userRecords = await $pocketbase.collection('users').getFullList({
                sort: 'created',
                fields: 'id,email,first_name,last_name'
            });
            users = userRecords.map(user => ({
                value: user.id,
                name: `${user.first_name} ${user.last_name} (${user.email})`
            }));

            // Fetch groups
            const groupRecords = await $pocketbase.collection('notification_groups').getFullList({
                sort: 'name'
            });
            groups = groupRecords.map(group => ({
                value: group.id,
                name: group.name
            }));

            // Fetch enabled notification providers
            console.log('Fetching notification providers...');
            const providerRecords = await $pocketbase.collection('providers')
                .getFullList({
                    filter: 'enabled = true && use ~ "notification"',
                    sort: 'name'
                });
            console.log('Provider records:', providerRecords);
            
            providers = providerRecords.map(provider => ({
                value: provider.id,
                name: provider.name
            }));
            console.log('Mapped providers:', providers);
        } catch (error) {
            console.error('Error fetching users, groups, and providers:', error);
            if (error.response?.data) {
                console.error('Response data:', error.response.data);
            }
        }
    }

    async function saveRules() {
        try {
            const records = await $pocketbase.collection('notification_settings').getFullList();
            if (records.length > 0) {
                const record = records[0];
                await $pocketbase.collection('notification_settings').update(record.id, {
                    rules: JSON.stringify(rules)
                });
            }
        } catch (error) {
            console.error('Error saving rules:', error);
        }
    }

    function addRule() {
        rules = [...rules, {
            id: crypto.randomUUID(),
            type: 'scan_finished',
            severity: [],
            channels: [],
            enabled: true,
            target_type: 'all'
        }];
        saveRules();
    }

    async function deleteRule(id: string) {
        rules = rules.filter(rule => rule.id !== id);
        await saveRules();
    }

    onMount(loadUsersAndGroups);
</script>

<div class="space-y-4">
    <div class="flex justify-end">
        <Button color="primary" on:click={addRule}>Add Rule</Button>
    </div>

    {#each rules as rule (rule.id)}
        <Card>
            <div class="flex justify-between items-center mb-4">
                <Select 
                    items={notificationTypes} 
                    bind:value={rule.type}
                    on:change={saveRules}
                />
                <div class="flex items-center gap-2">
                    <Toggle bind:checked={rule.enabled} on:change={saveRules} />
                    <Button color="red" size="xs" on:click={() => deleteRule(rule.id)}>Delete</Button>
                </div>
            </div>
            <div class="space-y-4">
                <div>
                    <label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
                        Severity Levels
                    </label>
                    <MultiSelect 
                        items={severityLevels} 
                        bind:value={rule.severity}
                        on:change={saveRules}
                    />
                </div>
                <div>
                    <label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
                        Notification Channels
                    </label>
                    <MultiSelect 
                        items={providers} 
                        bind:value={rule.channels}
                        on:change={saveRules}
                    />
                </div>
                {#if rule.channels.some(channel => providers.find(p => p.value === channel)?.name === 'Email')}
                    <div>
                        <label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
                            Target Recipients
                        </label>
                        <Select 
                            items={targetTypes} 
                            bind:value={rule.target_type} 
                            class="mb-2"
                            on:change={saveRules}
                        />
                        {#if rule.target_type === 'user'}
                            <div class="mt-2">
                                <label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
                                    Select Users
                                </label>
                                <MultiSelect 
                                    items={users} 
                                    bind:value={rule.target_users}
                                    on:change={saveRules}
                                />
                            </div>
                        {:else if rule.target_type === 'group'}
                            <div class="mt-2">
                                <label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
                                    Select Groups
                                </label>
                                <MultiSelect 
                                    items={groups} 
                                    bind:value={rule.target_groups}
                                    on:change={saveRules}
                                />
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        </Card>
    {/each}
</div> 