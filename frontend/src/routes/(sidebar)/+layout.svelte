<script>
	import '../../app.pcss';
	import Navbar from './Navbar.svelte';
	import Sidebar from './Sidebar.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { pocketbase } from '$lib/stores/pocketbase';
	import { page } from '$app/stores';

	let drawerHidden = false;

	onMount(() => {
		// Check if the user is authenticated
		if (!$pocketbase.authStore.isValid) {
			// Redirect to the login page if not authenticated
			goto('/authentication/sign-in');
		}
	});

	// Determine if the current page is an authentication page
	const isAuthPage = $page.url.pathname.startsWith('/authentication');
</script>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900">
	<header
		class="fixed top-0 z-40 w-full border-b border-gray-200 bg-white dark:border-gray-600 dark:bg-gray-800"
	>
		<Navbar bind:drawerHidden />
	</header>

	<div class="flex min-h-screen pt-[64px]">
		<Sidebar bind:drawerHidden />
		<main class="flex-1 lg:ml-64 bg-gray-50 dark:bg-gray-900">
			<div class="min-h-screen">
				<slot />
			</div>
		</main>
	</div>
</div>
