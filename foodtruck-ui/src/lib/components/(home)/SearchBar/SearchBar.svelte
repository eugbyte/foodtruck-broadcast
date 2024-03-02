<script lang="ts">
	import {
		type AutocompleteOption,
		type DrawerStore,
		type PopupSettings,
		popup
	} from "@skeletonlabs/skeleton";
	import type { SearchBarFocus$ } from "$lib/packages/store";
	import { tick } from "svelte";
	import { search } from "svelte-awesome/icons";
	import { Icon } from "svelte-awesome";
	import { DataList } from "$lib/components/shared";

	export let value: string;
	export let placeholder = "";

	export let searchBarFocus$: SearchBarFocus$;
	export let drawerStore$: DrawerStore;

	export let autoCompleteOptions: AutocompleteOption[] = [];
	type OnSubmit = () => void;
	export let onSubmit: OnSubmit = () => undefined;

	/**
	 * 1 way binding for easier tracking and debugging,
	 * since the only way for the drawer to open is through the change in searchBarFocus store state.
	 * When the input component is mounted, use handleFocus()
	 */
	const handleFocus = (node: HTMLInputElement) => {
		searchBarFocus$.subscribe((shouldOpen) => {
			if (shouldOpen) {
				drawerStore$.open();
				node.focus();
			}
		});
		return {
			destroy() {
				drawerStore$.close();
			}
		};
	};

	const onSelection = async (event: CustomEvent<AutocompleteOption>) => {
		value = event.detail.value as string;
		await tick();
		drawerStore$.close();
		onSubmit();
	};

	/**
	 * Required to force the auto select list to be visibility:hidden when an option is selected
	 */
	const popupSettings: PopupSettings = {
		event: "focus-click",
		target: "popupAutocomplete",
		placement: "bottom"
	};
</script>

<!-- Cannot use AppBar, which conflicts with AutoComplete -->
<div class="w-full flex flex-row justify-center py-4 bg-surface-800 px-1">
	<form on:submit={() => onSubmit()} class="relative w-full sm:w-2/3">
		<input
			type="text"
			use:handleFocus
			class="input autocomplete"
			{placeholder}
			bind:value
			on:click={() => searchBarFocus$.next(true)}
			data-testid="searchbar"
			use:popup={popupSettings}
		/>
		<button type="button" class="btn-icon variant-filled inline-block absolute right-0">
			<Icon slot="lead" data={search} scale={1} on:click={onSubmit} />
		</button>
	</form>
	<!-- important to (1) not wrap the AutoComplete in another div -->
	<div class="sm:w-2/3 shadow-2xl bg-surface-700" data-popup="popupAutocomplete">
		<DataList bind:inputValue={value} {autoCompleteOptions} {onSelection} />
	</div>
</div>
