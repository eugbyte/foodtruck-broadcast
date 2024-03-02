import { render, screen } from "@testing-library/svelte";
import Layout from "./+layout.svelte";
import html from "svelte-htm";

describe("test home page", () => {
	it("page should render within layout", () => {
		render(html`
        <${Layout}>
            <h2>Hello!</h2>
        </${Layout}>`);
		expect(screen.getByText("Hello!")).toBeInTheDocument();
	});
});
