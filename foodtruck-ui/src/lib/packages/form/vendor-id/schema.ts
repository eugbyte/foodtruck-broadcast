import { z } from "zod";
import type { SuperForm } from "sveltekit-superforms/client";

export const schema = z.object({
	vendorID: z.string().min(1)
});

export type Schema = typeof schema;

export type FormHook = SuperForm<Schema>;
