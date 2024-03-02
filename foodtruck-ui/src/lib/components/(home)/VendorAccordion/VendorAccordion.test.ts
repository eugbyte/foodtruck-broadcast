import { render, screen } from "@testing-library/svelte";
import Accordion from "./VendorAccordion.svelte";

describe("test Accordion", () => {
	it("Accordion should render", () => {
		render(Accordion, {
			props: {
				header: "Vendor Info"
			}
		});
		expect(screen.getByText("Vendor Info")).toBeInTheDocument();
	});
});
