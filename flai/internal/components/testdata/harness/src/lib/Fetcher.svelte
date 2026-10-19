<script module>
	// The fixture of a component that asks the app's own server (S-0449): its
	// module asks once as it is imported, and the component, as it mounts,
	// fetches by a string, a URL, and a Request, fetches a module of the dev
	// server's, and opens an EventSource. It shows what each was answered as
	// JSON in .fetched, once all are.
	const early = fetch('/api/module').then((r) => r.status);
</script>

<script>
	let answers = $state(null);

	async function answer(r) {
		return { status: r.status, statusText: r.statusText, type: r.headers.get('content-type'), body: await r.text() };
	}

	function stream() {
		return new Promise((resolve) => {
			const source = new EventSource('/api/events?project=x');
			source.onerror = () => {
				const state = source.readyState;
				source.close();
				resolve(state);
			};
		});
	}

	$effect(() => {
		Promise.all([
			fetch('/api/x'),
			fetch(new URL('/api/y?z=1', location.href), { method: 'put' }),
			fetch(new Request('/api/r', { method: 'POST', body: '{}' })),
			fetch('/src/app.css')
		])
			.then((rs) => Promise.all(rs.map(answer)))
			.then(async ([string, url, request, asset]) => {
				answers = { module: await early, string, url, request, asset, stream: await stream() };
			});
	});
</script>

{#if answers}
	<pre class="fetched">{JSON.stringify(answers)}</pre>
{/if}
