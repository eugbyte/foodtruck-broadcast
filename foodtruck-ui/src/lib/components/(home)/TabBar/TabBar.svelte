<script lang="ts">
	import { Tab, TabGroup } from "@skeletonlabs/skeleton";
	import { userCircle, search, home } from "svelte-awesome/icons";
	import Icon from "svelte-awesome";
	import type { SearchBarFocus$ } from "$lib/packages/store";
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { routeToPage } from "$lib/packages/route";

	let tabSet = "";
	export let searchBarFocus$: SearchBarFocus$;
	const focusSearchBar = () => searchBarFocus$.next(true);
	onMount(() => {
		console.log("tab bar mounted");
		page.subscribe((route) => (tabSet = route.url.pathname));
	});
</script>

<TabGroup
	justify="justify-center"
	active="variant-filled-primary"
	hover="hover:variant-soft-primary"
	flex="flex-1 lg:flex-none"
	rounded="rounded"
	border=""
	class="bg-surface-100-800-token w-full"
	data-testid="tab-group"
>
	<Tab bind:group={tabSet} value={"/"} name="user" on:click={() => routeToPage("/", false)}>
		<Icon slot="lead" data={home} scale={2} />
	</Tab>
	<Tab bind:group={tabSet} value={"_"} name="search" on:click={focusSearchBar} id="search_btn_icon">
		<Icon slot="lead" data={search} scale={2} data-testid="search_btn" />
	</Tab>
	<Tab
		bind:group={tabSet}
		value={"/profile/"}
		name="current_location"
		on:click={() => routeToPage("/profile/", false)}
	>
		<Icon slot="lead" data={userCircle} scale={2} />
	</Tab>
</TabGroup>
