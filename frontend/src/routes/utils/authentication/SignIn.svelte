<script lang="ts">
	import { Label, Input, Checkbox, A, Button, Card } from 'flowbite-svelte';
	export let title = 'Sign in to platform';
	export let site = {
		name: 'Bitor',
		img: '/images/Bitor_Main_Logo.png',
		link: '/',
		imgAlt: 'Bitor Logo'
	};
	export let rememberMe = true;
	export let lostPassword = true;
	export let createAccount = true;
	export let lostPasswordLink = '';
	export let loginTitle = 'Login to your account';
	export let registerLink = '';
	export let createAccountTitle = 'Create account';

	export let mainClass = 'bg-gray-50 dark:bg-gray-900 w-full';
	export let mainDivClass =
		'flex flex-col items-center justify-center px-6 pt-8 mx-auto md:h-screen pt:mt-0 dark:bg-gray-900';
	export let siteLinkClass =
		'flex items-center justify-center mb-8 text-2xl font-semibold lg:mb-10 dark:text-white';
	export let siteImgClass = 'mr-4 h-11';
	export let cardH1Class = 'text-2xl font-bold text-gray-900 dark:text-white';
</script>

<main class={mainClass}>
	<div class={mainDivClass}>
		<a href={site.link} class={siteLinkClass}>
			<img src={site.img} class={siteImgClass} alt={site.imgAlt} />
		</a>
		<!-- Card -->
		<Card class="w-full" size="md" border={false}>
			<h1 class={cardH1Class}>
				{title}
			</h1>
			<form class="mt-8 space-y-6" on:submit|preventDefault>
				<slot />
				{#if rememberMe || lostPassword}
					<div class="flex items-start">
						{#if rememberMe}
							<Checkbox class="accent-primary-600" name="remember">Remember me</Checkbox>
						{/if}
						{#if lostPassword}
							<A href={lostPasswordLink} aClass="ml-auto text-sm">Lost Password?</A>
						{/if}
					</div>
				{/if}
				<Button type="submit" size="lg">{loginTitle}</Button>
				{#if createAccount}
					<div class="text-sm font-medium text-gray-500 dark:text-gray-400">
						Not registered? <A href={registerLink}>{createAccountTitle}</A>
					</div>
				{/if}
			</form>
		</Card>
	</div>
</main>