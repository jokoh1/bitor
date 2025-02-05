<script lang="ts">
	import { Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import Icon from '@iconify/svelte';
	import ChangePassword from './ChangePassword.svelte';
	import type { Record } from 'pocketbase';

	export let user: Record;
	let showPasswordModal = false;
</script>

<div class="flex items-center justify-between p-4 bg-white rounded-lg shadow">
	<div class="flex items-center space-x-4">
		<div class="flex-shrink-0">
			<img class="w-8 h-8 rounded-full" src={user.avatar || '/default-avatar.png'} alt="User avatar" />
		</div>
		<div>
			<p class="font-medium text-gray-900">{user.name || 'Unnamed User'}</p>
			<p class="text-sm text-gray-500">{user.email}</p>
		</div>
	</div>
	
	<div class="flex items-center space-x-2">
		<Button size="sm" class="gap-2 px-3" color="blue" on:click={() => showPasswordModal = true}>
			<Icon icon="mdi:key" class="w-4 h-4" /> Reset Password
		</Button>
		<Dropdown>
			<Button slot="trigger" size="sm">
				<Icon icon="mdi:dots-vertical" class="w-4 h-4" />
			</Button>
			<DropdownItem>Edit User</DropdownItem>
			<DropdownItem class="text-red-600">Delete User</DropdownItem>
		</Dropdown>
	</div>
</div>

<ChangePassword bind:open={showPasswordModal} userId={user.id} />
