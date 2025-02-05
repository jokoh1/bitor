<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Input, Select, Label } from 'flowbite-svelte';

    export let settings: {
        smtp_host: string;
        smtp_port: number;
        smtp_username: string;
        smtp_password: string;
        from_address: string;
        to_addresses: string[];
        smtp_encryption: 'tls' | 'ssl' | 'none';
    };

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change', settings);
    }

    function handleToAddressesChange(event: Event) {
        const input = event.target as HTMLInputElement;
        settings.to_addresses = input.value.split(',').map(addr => addr.trim());
        handleChange();
    }
</script>

<div class="space-y-4">
    <div>
        <Label for="smtp_host">SMTP Host</Label>
        <Input
            id="smtp_host"
            type="text"
            placeholder="smtp.example.com"
            bind:value={settings.smtp_host}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="smtp_port">SMTP Port</Label>
        <Input
            id="smtp_port"
            type="number"
            placeholder="587"
            bind:value={settings.smtp_port}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="smtp_username">SMTP Username</Label>
        <Input
            id="smtp_username"
            type="text"
            placeholder="username@example.com"
            bind:value={settings.smtp_username}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="smtp_password">SMTP Password</Label>
        <Input
            id="smtp_password"
            type="password"
            placeholder="••••••••"
            bind:value={settings.smtp_password}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="from_address">From Address</Label>
        <Input
            id="from_address"
            type="email"
            placeholder="notifications@example.com"
            bind:value={settings.from_address}
            on:change={handleChange}
        />
    </div>
    <div>
        <Label for="to_addresses">To Addresses (comma-separated)</Label>
        <Input
            id="to_addresses"
            type="text"
            placeholder="user1@example.com, user2@example.com"
            value={settings.to_addresses.join(', ')}
            on:change={handleToAddressesChange}
        />
    </div>
    <div>
        <Label for="smtp_encryption">SMTP Encryption</Label>
        <Select
            id="smtp_encryption"
            bind:value={settings.smtp_encryption}
            on:change={handleChange}
        >
            <option value="tls">TLS</option>
            <option value="ssl">SSL</option>
            <option value="none">None</option>
        </Select>
    </div>
</div> 