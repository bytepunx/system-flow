<script module>
	// The fixture of a component the page's answers answer (S-0476): its
	// module asks once as it is imported, and the component, as it mounts,
	// fetches a path with one query and with another, a POST, a 204, the path
	// of an EventSource's answer, and a path no answer has, and opens an
	// EventSource with an answer and one without. It shows what each was
	// answered as JSON in .answered, once all are.
	const early = fetch('/api/module').then(answer);

	async function answer(r) {
		return { status: r.status, type: r.headers.get('content-type'), flai: r.headers.get('x-flai'), body: await r.text() };
	}
</script>

<script>
	let answers = $state(null);

	// stream is what the EventSource for path heard, in order, each as its
	// type, or onmessage for its handler, its data, and its lastEventId, or
	// open and error as their state, until a change event, when it closes it,
	// or an error; and its state a moment later, when an event after the close
	// would have come.
	function stream(path) {
		return new Promise((resolve) => {
			const source = new EventSource(path);
			const heard = [];
			const end = () => {
				source.close();
				setTimeout(() => resolve({ heard, state: source.readyState }), 100);
			};
			const message = (kind) => (e) => heard.push(`${kind}|${e.data}|${e.lastEventId}`);
			source.onopen = () => heard.push(`open|${source.readyState}`);
			source.onmessage = message('onmessage');
			source.addEventListener('message', message('message'));
			source.addEventListener('change', (e) => {
				message('change')(e);
				end();
			});
			source.onerror = () => {
				heard.push(`error|${source.readyState}`);
				end();
			};
		});
	}

	$effect(() => {
		Promise.all([
			fetch('/api/items?id=1'),
			fetch('/api/items?id=2'),
			fetch('/api/accept', { method: 'POST', body: '{}' }),
			fetch('/api/gone', { method: 'DELETE' }),
			fetch('/api/events?project=x'),
			fetch('/api/missing')
		])
			.then((rs) => Promise.all(rs.map(answer)))
			.then(async ([query, other, post, empty, events, missing]) => {
				const module = await early;
				answers = { module, query, other, post, empty, events, missing, stream: await stream('/api/events?project=x'), silent: await stream('/api/silent') };
			});
	});
</script>

{#if answers}
	<pre class="answered">{JSON.stringify(answers)}</pre>
{/if}
