<script lang="ts">
    import { page } from '$app/stores';
    import { Sidebar, SidebarGroup, SidebarItem } from 'flowbite-svelte';
    import {
        HomeSolid,
        ChartPieSolid,
        CogSolid,
        CloudArrowUpSolid,
        BellSolid,
        UsersSolid,
        ShieldCheckSolid,
        ServerSolid,
        FileSolid,
        FolderOpenSolid
    } from 'flowbite-svelte-icons';
</script>

<Sidebar>
    <SidebarGroup>
        <SidebarItem label="Dashboard" href="/" active={$page.url.pathname === '/'}>
            <svelte:fragment slot="icon">
                <HomeSolid class="w-5 h-5" />
            </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Findings" href="/findings" active={$page.url.pathname.startsWith('/findings')}>
            <svelte:fragment slot="icon">
                <ShieldCheckSolid class="w-5 h-5" />
            </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Assets" href="/assets" active={$page.url.pathname.startsWith('/assets')}>
            <svelte:fragment slot="icon">
                <ServerSolid class="w-5 h-5" />
            </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Reports" href="/reports" active={$page.url.pathname.startsWith('/reports')}>
            <svelte:fragment slot="icon">
                <FileSolid class="w-5 h-5" />
            </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Projects" href="/projects" active={$page.url.pathname.startsWith('/projects')}>
            <svelte:fragment slot="icon">
                <FolderOpenSolid class="w-5 h-5" />
            </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Analytics" href="/analytics" active={$page.url.pathname.startsWith('/analytics')}>
            <svelte:fragment slot="icon">
                <ChartPieSolid class="w-5 h-5" />
            </svelte:fragment>
        </SidebarItem>
    </SidebarGroup>

    {#if $page.url.pathname.startsWith('/settings')}
        <SidebarGroup>
            <SidebarItem label="General" href="/settings" activeClass={$page.url.pathname === '/settings' ? 'active' : ''}>
                <svelte:fragment slot="icon">
                    <CogSolid class="w-5 h-5" />
                </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Providers" href="/settings/providers" active={$page.url.pathname === '/settings/providers'}>
                <svelte:fragment slot="icon">
                    <CloudArrowUpSolid class="w-5 h-5" />
                </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Notifications" href="/settings/notifications" active={$page.url.pathname === '/settings/notifications'}>
                <svelte:fragment slot="icon">
                    <BellSolid class="w-5 h-5" />
                </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Users" href="/settings/users" active={$page.url.pathname === '/settings/users'}>
                <svelte:fragment slot="icon">
                    <UsersSolid class="w-5 h-5" />
                </svelte:fragment>
            </SidebarItem>
        </SidebarGroup>
    {/if}
</Sidebar> 