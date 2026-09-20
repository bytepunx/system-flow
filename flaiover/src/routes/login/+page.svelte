<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';

	let token = $state('');
	let error = $state<string | null>(null);
	let busy = $state(false);

	function nextPath(): string {
		const n = page.url.searchParams.get('next') ?? '/';
		return n.startsWith('/') && !n.startsWith('//') ? n : '/';
	}

	async function submit(value: string) {
		busy = true;
		error = null;
		const r = await fetch('/api/login', {
			method: 'POST',
			headers: { 'content-type': 'application/json', 'x-requested-with': 'flaiover' },
			body: JSON.stringify({ token: value })
		});
		if (!r.ok) {
			busy = false;
			error = 'That token was not accepted.';
			return;
		}
		// The token was right. Whether the browser kept the session cookie is another matter: ask
		// for something that needs it before going on, or a dropped cookie shows up as this page
		// asking for the token again, with no reason given (S-0083).
		const kept = await fetch('/api/agent').catch(() => null);
		busy = false;
		if (kept && kept.status === 401) {
			error =
				'The token is right, but this browser did not keep the session cookie, so you are not logged in. ' +
				(location.protocol === 'http:'
					? 'This page is served over plain HTTP: a proxy or tunnel in front of the dashboard that reports HTTPS (X-Forwarded-Proto) makes the cookie HTTPS-only. '
					: '') +
				'Check that cookies are allowed for this address, or use the address the dashboard is really served at.';
			return;
		}
		// The token came in the fragment; leave no trace of it in history.
		history.replaceState(null, '', '/login');
		await goto(resolve(nextPath() as '/'));
	}

	onMount(() => {
		const m = /(?:^#|&)token=([^&]+)/.exec(location.hash);
		if (m) {
			const value = decodeURIComponent(m[1]);
			history.replaceState(null, '', location.pathname + location.search);
			void submit(value);
		}
	});
</script>

<div class="mx-auto mt-24 max-w-md rounded-lg border border-line bg-surface p-6 shadow-sm">
	<h1 class="text-lg font-semibold">Log in to flaiover</h1>
	<p class="mt-2 text-sm text-ink-soft">
		Paste the project token, or open the login link that <code>flai dashboard token</code> prints.
	</p>
	<form
		class="mt-4 flex gap-2"
		onsubmit={(e) => {
			e.preventDefault();
			void submit(token);
		}}
	>
		<input
			class="flex-1 rounded border border-line-strong px-3 py-2 font-mono text-sm"
			type="password"
			placeholder="token"
			bind:value={token}
			autocomplete="off"
		/>
		<button
			class="rounded bg-primary px-4 py-2 text-sm text-on-primary disabled:opacity-50"
			disabled={busy || !token}
		>
			Log in
		</button>
	</form>
	{#if error}<p class="mt-3 text-sm text-danger" role="alert" data-testid="login-error">
			{error}
		</p>{/if}
</div>
