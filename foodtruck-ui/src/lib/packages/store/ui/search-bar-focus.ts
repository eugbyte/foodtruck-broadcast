import { Subject } from "rxjs";

type callback = (val: boolean) => void;

export type SearchBarFocus$ = Pick<Subject<boolean>, "subscribe" | "next">;

function createSearchBarFocus(): SearchBarFocus$ {
	const subject = new Subject<boolean>();

	return {
		subscribe: (cb: callback) => subject.subscribe(cb),
		next: (val: boolean) => subject.next(val)
	};
}

export const searchBarFocus$ = createSearchBarFocus();
