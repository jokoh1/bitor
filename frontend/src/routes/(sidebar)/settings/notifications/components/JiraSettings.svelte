<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Input, Select, Textarea, Label } from 'flowbite-svelte';

    export let settings: {
        url: string;
        username: string;
        api_token: string;
        project_key: string;
        issue_type: string;
        ticket_template: string;
    };

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change', settings);
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="url">Jira URL</Label>
        <Input
            id="url"
            type="url"
            placeholder="https://your-domain.atlassian.net"
            bind:value={settings.url}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="username">Username (Email)</Label>
        <Input
            id="username"
            type="text"
            placeholder="user@example.com"
            bind:value={settings.username}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="api_token">API Token</Label>
        <Input
            id="api_token"
            type="password"
            placeholder="Jira API token"
            bind:value={settings.api_token}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="project_key">Project Key</Label>
        <Input
            id="project_key"
            type="text"
            placeholder="PROJECT"
            bind:value={settings.project_key}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="issue_type">Issue Type</Label>
        <Select
            id="issue_type"
            bind:value={settings.issue_type}
            on:change={handleChange}
        >
            <option value="Task">Task</option>
            <option value="Bug">Bug</option>
            <option value="Story">Story</option>
            <option value="Epic">Epic</option>
        </Select>
    </div>
    <div>
        <Label for="ticket_template">Ticket Template</Label>
        <Textarea
            id="ticket_template"
            placeholder="Enter your ticket template with placeholders like {severity}, {type}, {target}, etc."
            rows={4}
            bind:value={settings.ticket_template}
            on:change={handleChange}
        />
    </div>
</div> 