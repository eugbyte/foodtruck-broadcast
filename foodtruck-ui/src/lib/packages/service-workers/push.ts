interface Notification {
	title: string;
	body: string;
	icon: string;
}

const sw = self as unknown as ServiceWorkerGlobalScope;

export const onPush = (event: PushEvent, onComplete = () => undefined) => {
	if (event == null || event.data == null) {
		return;
	}
	const data: Notification = event.data.json();
	const displayPromise: Promise<void> = sw.registration.showNotification(data.title, {
		body: data.body
	});
	const promiseChain = Promise.all([displayPromise])
		.then(() => onComplete())
		.catch((err) => console.error(err));
	event.waitUntil(promiseChain);
};
