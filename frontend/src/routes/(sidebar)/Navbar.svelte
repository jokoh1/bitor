<script>
	import UserMenu from '../utils/widgets/UserMenu.svelte';
	import {
		DarkMode,
		Dropdown,
		DropdownItem,
		NavBrand,
		NavHamburger,
		NavLi,
		NavUl,
		Navbar,
		Search
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import '../../app.pcss';
	import { pocketbase } from '$lib/stores/pocketbase'; // Import the pocketbase store

	import { onMount, onDestroy } from 'svelte'; // Import onMount and onDestroy

	export let fluid = true;
	export let drawerHidden = false;
	export let list = false;

	/** @type {import("pocketbase").AuthModel | null} */
	let currentUser = null;

	// Reactive statement to safely assign the current user
	$: if ($pocketbase.authStore.model) {
		currentUser = $pocketbase.authStore.model;
	} else {
		currentUser = null;
	}

</script>

<Navbar {fluid} class="text-black" color="default" let:NavContainer>
	<NavHamburger
		onClick={() => (drawerHidden = !drawerHidden)}
		class="m-0 me-3 md:block lg:hidden"
	/>
	<NavBrand href="/" class={list ? 'w-40' : 'lg:w-60'}>
		<!-- Logo for Light Mode -->
		<img
			src="/images/Orbit-Main-Logo.png"
			class="me-2.5 h-6 sm:h-8 block dark:hidden"
			alt="Orbit Logo Light"
		/>

		<!-- Logo for Dark Mode -->
		<img
			src="/images/Orbit_White_Logo.png"
			class="me-2.5 h-6 sm:h-8 hidden dark:block"
			alt="Orbit Logo Dark"
		/>
	</NavBrand>

	<div class="ms-auto flex items-center text-gray-500 dark:text-gray-400 sm:order-2">
		<DarkMode />
		{#if currentUser}
			<UserMenu user={currentUser} />
		{:else}
			<p>Loading user...</p>
		{/if}
	</div>
</Navbar>
