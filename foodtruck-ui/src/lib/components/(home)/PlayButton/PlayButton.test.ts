import { fireEvent, render, screen } from "@testing-library/svelte";
import MockComponent from "./PlayButton.test.svelte";

describe("test playbutton", () => {
	it("icon and text should change on click", async () => {
		const { getByRole } = render(MockComponent);

		const button = getByRole("button");
		expect(button).toBeInTheDocument();

		expect(screen.getByText("Pause")).toBeInTheDocument();
		await fireEvent.click(button);

		const element = await screen.findByText("Resume");
		expect(element).toBeInTheDocument();
	});
});
