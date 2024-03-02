import { render, screen } from "@testing-library/svelte";
import Alert from "./PermissionPrompt.svelte";

describe("test Alert", () => {
	it("Alert should render", () => {
		render(Alert, {
			props: {
				text: "Permission for notification",
				onAllow: () => undefined,
				onDismiss: () => undefined
			}
		});
		expect(screen.getByText("Permission for notification")).toBeInTheDocument();
	});
});
