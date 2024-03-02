<script lang="ts">
	import { Autocomplete, type AutocompleteOption } from "@skeletonlabs/skeleton";

	export let inputValue: string;
	export let autoCompleteOptions: AutocompleteOption[] = [];
	export let onSelection: (event: CustomEvent<AutocompleteOption>) => void;
	/**
	 * Classname to be used with [Skeleton popup utility](https://www.skeleton.dev/components/autocomplete#popup).
	 *
	 * To enable the list to close upon selecting an option.
	 */
</script>

<!--
  @component

  important to (1) not wrap DataList in more than one nested div
-->
<!-- (2) use the tailwind `truncate` class, (3) set a max width -->
<Autocomplete
	bind:input={inputValue}
	options={autoCompleteOptions}
	on:selection={onSelection}
	emptyState=""
	regionNav="truncate max-w-xs sm:max-w-none"
/>
