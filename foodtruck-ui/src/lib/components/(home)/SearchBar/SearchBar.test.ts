import { render, screen } from "@testing-library/svelte";
import MockComponent from "./SearchBar.test.svelte";

describe("test SearchBar", () => {
	it("searchbar should render", () => {
		render(MockComponent);
		expect(screen.getByPlaceholderText("Address")).toBeInTheDocument();

		const inputElem = screen.getByTestId("searchbar");
		expect(inputElem).toBeInTheDocument();
	});
});
