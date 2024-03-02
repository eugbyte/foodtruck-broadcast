<script lang="ts">
	// Most of your app wide CSS should be put in this file
	import { getDrawerStore, type AutocompleteOption } from "@skeletonlabs/skeleton";
	import { AppShell, Drawer } from "@skeletonlabs/skeleton";
	import { SearchBar, TabBar } from "$lib/components/(home)";
	import { zIndexes } from "$lib/packages/config";
	import { geocode$, searchBarFocus$ } from "$lib/packages/store";
	import debounce from "debounce-promise";
	import { getSuggestions, type AddressAutofillSuggestion } from "$lib/packages/geocode";
	import { base as _base } from "$app/paths";

	let suggestions: AddressAutofillSuggestion[] = [];
	let autoCompleteOptions: AutocompleteOption[] = [];
	let searchValue = "";

	const debounceSuggest = debounce((text: string) => getSuggestions(text), 2000);

	const handleChange = async (text: string) => {
		suggestions = await debounceSuggest(text);
		autoCompleteOptions = suggestions.map((sug) => ({
			label: sug.full_address ?? "",
			value: sug.full_address
		}));
		console.log({ autoCompleteOptions });
	};

	const clearAddresses = () => {
		autoCompleteOptions = [];
	};

	const onSubmit = () => {
		if (searchValue.trim() != "") {
			geocode$.next(searchValue);
		}
	};

	$: {
		if (searchValue.trim() != "") {
			handleChange(searchValue);
		}
	}

	$: {
		// prevent race condition with debounce delay
		// when user backspace to empty string, debounce function will be return after with one character left.
		if (autoCompleteOptions.length > 0 && searchValue == "") {
			clearAddresses();
		} else if (searchValue == "") {
			clearAddresses();
		}
	}
	const base = _base.trim() == "" ? "." : _base;
</script>

<!-- App Shell -->
<AppShell slotHeader={`${zIndexes.appbar}`}>
	<!-- (header) -->
	<SearchBar
		slot="header"
		bind:value={searchValue}
		placeholder="Address"
		{searchBarFocus$}
		drawerStore$={getDrawerStore()}
		{autoCompleteOptions}
		{onSubmit}
	/>
	<!-- Page Route Content -->
	<slot />
	<!-- (footer) -->
	<TabBar slot="footer" {searchBarFocus$} />
	<!-- Need to make sure z index of navbar, i.e. slotHeader, is greater than the drawer for the nav bar to overlay the drawer -->
	<Drawer position="top" zIndex={`${zIndexes.drawer}`} />
</AppShell>

<svelte:head>
	<title>Foodtruck App</title>
	<link rel="manifest" href="{base}/manifest.webmanifest" crossorigin="use-credentials" />
</svelte:head>
