import { render } from "@testing-library/svelte";
import TabBar from "./TabBar.svelte";
import { searchBarFocus$ } from "$lib/packages/store";

describe("test tab bar", () => {
	it("tab bar should render", () => {
		render(TabBar, {
			props: {
				searchBarFocus$
			}
		});
		const icons = document.getElementsByTagName("svg");
		expect(icons).not.toBeNull();
		expect(icons.length).toBe(3);
	});
});
